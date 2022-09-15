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

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("s3", "Buckets", float64(300), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	s3Checker := NewS3Checker()
	svcChecker := s3Checker.(*ServiceChecker)
	actual := svcChecker.getS3BucketUsage()
	assert.Len(t, actual, 1)
	usage := actual[0]
	assert.Equal(t, "s3", usage.Service)
	assert.Equal(t, float64(300), usage.QuotaValue)
	assert.Equal(t, float64(len(mockedS3Output.Buckets)), usage.UsageValue)
}

func TestGetS3BucketUsageError(t *testing.T) {
	mockedS3Output := s3.ListBucketsOutput{
		Buckets: []*s3.Bucket{},
	}
	conf.S3 = mockedS3ClientListBucketsMsg{Resp: mockedS3Output, Error: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("s3", "Buckets", float64(300), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	s3Checker := NewS3Checker()
	svcChecker := s3Checker.(*ServiceChecker)
	actual := svcChecker.getS3BucketUsage()
	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}
