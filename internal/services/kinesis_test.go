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
	require.Implements(t, (*ServiceQuota)(nil), NewKinesisChecker())
}

func TestGetKinesisShardUsage(t *testing.T) {
	mockedkinesisOutput := kinesis.DescribeLimitsOutput{
		OpenShardCount: aws.Int64(2),
	}
	conf.Kinesis = mockedKinesisDescribeLimitsMsg{Resp: mockedkinesisOutput, Error: nil}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("kinesis", "Shards per Region", float64(10), false)},
		nil)
	kinesisChecker := NewKinesisChecker()
	svcChecker := kinesisChecker.(*ServiceChecker)
	actual := svcChecker.getKinesisShardUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "kinesis", quota.Service)
	assert.Equal(t, float64(10), quota.QuotaValue)
	assert.Equal(t, float64(2), quota.UsageValue)
}

func TestGetKinesisShardUsageError(t *testing.T) {
	mockedkinesisOutput := kinesis.DescribeLimitsOutput{
		OpenShardCount: aws.Int64(2),
	}
	conf.Kinesis = mockedKinesisDescribeLimitsMsg{Resp: mockedkinesisOutput, Error: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("kinesis", "Shards per Region", float64(10), false)},
		nil)

	kinesisChecker := NewKinesisChecker()
	svcChecker := kinesisChecker.(*ServiceChecker)
	actual := svcChecker.getKinesisShardUsage()
	expected := []AWSQuotaInfo{}

	assert.Equal(t, expected, actual)
}

func TestGetKinesisOnDemandStreamCountUsage(t *testing.T) {
	mockedkinesisOutput := kinesis.DescribeLimitsOutput{
		OnDemandStreamCount:      aws.Int64(10),
		OnDemandStreamCountLimit: aws.Int64(200),
	}
	conf.Kinesis = mockedKinesisDescribeLimitsMsg{Resp: mockedkinesisOutput, Error: nil}

	kinesisChecker := NewKinesisChecker()
	svcChecker := kinesisChecker.(*ServiceChecker)
	actual := svcChecker.getKinesisOnDemandStreamCountUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "kinesis", quota.Service)
	assert.Equal(t, float64(200), quota.QuotaValue)
	assert.Equal(t, float64(10), quota.UsageValue)
	assert.True(t, quota.Global)
}

func TestGetKinesisOnDemandStreamCountUsageError(t *testing.T) {
	mockedkinesisOutput := kinesis.DescribeLimitsOutput{
		OnDemandStreamCount:      aws.Int64(10),
		OnDemandStreamCountLimit: aws.Int64(200),
	}
	conf.Kinesis = mockedKinesisDescribeLimitsMsg{Resp: mockedkinesisOutput, Error: errors.New("test error")}

	kinesisChecker := NewKinesisChecker()
	svcChecker := kinesisChecker.(*ServiceChecker)
	actual := svcChecker.getKinesisOnDemandStreamCountUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}
