package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedKinesisDescribeLimitsMsg struct {
	KinesisClientInterface
	Resp  kinesis.DescribeLimitsOutput
	Error error
}

func (m mockedKinesisDescribeLimitsMsg) DescribeLimits(input *kinesis.DescribeLimitsInput) (*kinesis.DescribeLimitsOutput, error) {
	return &m.Resp, m.Error
}

func TestNewKinesisCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewKinesisChecker(nil, nil))
}

func TestGetKinesisShardUsage(t *testing.T) {
	mockedkinesisOutput := kinesis.DescribeLimitsOutput{
		OpenShardCount: aws.Int64(2),
	}
	kinesisClient = mockedKinesisDescribeLimitsMsg{Resp: mockedkinesisOutput, Error: nil}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("kinesis", "Shards per Region", float64(10), false),
		},
	}
	mockedSvcQuotaClient := mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	kinesisChecker := NewKinesisChecker(nil, mockedSvcQuotaClient)
	actual := kinesisChecker.GetUsage()
	assert.Equal(t, 1, len(actual))
	firstQuota := actual[0]
	assert.Equal(t, "kinesis", firstQuota.Service)
	assert.Equal(t, float64(10), firstQuota.QuotaValue)
	assert.Equal(t, float64(2), firstQuota.UsageValue)
}

func TestGetKinesisShardUsageError(t *testing.T) {
	mockedkinesisOutput := kinesis.DescribeLimitsOutput{
		OpenShardCount: aws.Int64(2),
	}
	kinesisClient = mockedKinesisDescribeLimitsMsg{Resp: mockedkinesisOutput, Error: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("kinesis", "Shards per Region", float64(10), false),
		},
	}
	mockedSvcQuotaClient := mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	kinesisChecker := NewKinesisChecker(nil, mockedSvcQuotaClient)
	actual := kinesisChecker.GetUsage()
	expected := []AWSQuotaInfo([]AWSQuotaInfo{{Service: "", Name: "", Region: "", Quotacode: "", QuotaValue: 0, UsageValue: 0, Unit: "", Global: false}})
	assert.Equal(t, expected, actual)
}
