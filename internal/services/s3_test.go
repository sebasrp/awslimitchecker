package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedS3ClientListBucketsMsg struct {
	S3ClientInterface
	Resp  s3.ListBucketsOutput
	Error error
}

func (m mockedS3ClientListBucketsMsg) ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {
	return &m.Resp, m.Error
}

func TestNewS3CheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewS3Checker())
}

func TestGetS3BucketUsage(t *testing.T) {
	mockedS3Output := s3.ListBucketsOutput{
		Buckets: []*s3.Bucket{},
	}
	conf.S3 = mockedS3ClientListBucketsMsg{Resp: mockedS3Output, Error: nil}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("s3", "Buckets", float64(300), false),
		},
	}
	conf.ServiceQuotas = mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	s3Checker := NewS3Checker()
	svcChecker := s3Checker.(*ServiceChecker)
	actual := svcChecker.getS3BucketUsage()
	assert.Equal(t, "s3", actual.Service)
	assert.Equal(t, float64(300), actual.QuotaValue)
	assert.Equal(t, float64(len(mockedS3Output.Buckets)), actual.UsageValue)
}

func TestGetS3BucketUsageError(t *testing.T) {
	mockedS3Output := s3.ListBucketsOutput{
		Buckets: []*s3.Bucket{},
	}
	conf.S3 = mockedS3ClientListBucketsMsg{Resp: mockedS3Output, Error: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("s3", "Buckets", float64(300), false),
		},
	}
	conf.ServiceQuotas = mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	s3Checker := NewS3Checker()
	svcChecker := s3Checker.(*ServiceChecker)
	actual := svcChecker.getS3BucketUsage()
	expected := AWSQuotaInfo{Service: "s3", Name: "Buckets", Region: "", Quotacode: "", QuotaValue: 300, UsageValue: 0, Unit: "", Global: false}
	assert.Equal(t, expected, actual)
}
