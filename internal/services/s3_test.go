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
	require.Implements(t, (*Svcquota)(nil), NewS3Checker(nil, nil))
}

func TestGetS3BucketUsage(t *testing.T) {
	mockedS3Output := s3.ListBucketsOutput{
		Buckets: []*s3.Bucket{},
	}
	s3Client = mockedS3ClientListBucketsMsg{Resp: mockedS3Output, Error: nil}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("s3", "Buckets", float64(300), false),
		},
	}
	mockedSvcQuotaClient := mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	s3Checker := NewS3Checker(nil, mockedSvcQuotaClient)
	actual := s3Checker.GetUsage()
	firstQuota := actual[0]
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, "s3", firstQuota.Service)
	assert.Equal(t, float64(300), firstQuota.QuotaValue)
	assert.Equal(t, float64(len(mockedS3Output.Buckets)), firstQuota.UsageValue)
}

func TestGetS3BucketUsageError(t *testing.T) {
	mockedS3Output := s3.ListBucketsOutput{
		Buckets: []*s3.Bucket{},
	}
	s3Client = mockedS3ClientListBucketsMsg{Resp: mockedS3Output, Error: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("s3", "Buckets", float64(300), false),
		},
	}
	mockedSvcQuotaClient := mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	s3Checker := NewS3Checker(nil, mockedSvcQuotaClient)
	actual := s3Checker.GetUsage()
	expected := []AWSQuotaInfo([]AWSQuotaInfo{{Service: "", Name: "", Region: "", Quotacode: "", QuotaValue: 0, UsageValue: 0, Unit: "", Global: false}})
	assert.Equal(t, expected, actual)
}
