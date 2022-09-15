package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedRdsClient struct {
	RdsClientInterface
	DescribeAccountAttributesResp  rds.DescribeAccountAttributesOutput
	DescribeAccountAttributesError error
}

func (m mockedRdsClient) DescribeAccountAttributes(input *rds.DescribeAccountAttributesInput) (*rds.DescribeAccountAttributesOutput, error) {
	return &m.DescribeAccountAttributesResp, m.DescribeAccountAttributesError
}

func TestNewRdsCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewRdsChecker())
}

func TestGetRdsAccountQuotas(t *testing.T) {
	mockedDescribeAccountAttributesOutput := rds.DescribeAccountAttributesOutput{
		AccountQuotas: []*rds.AccountQuota{{AccountQuotaName: aws.String("foo"), Max: aws.Int64(10), Used: aws.Int64(1)}},
	}
	conf.Rds = mockedRdsClient{DescribeAccountAttributesResp: mockedDescribeAccountAttributesOutput, DescribeAccountAttributesError: nil}

	rdsChecker := NewRdsChecker()
	svcChecker := rdsChecker.(*ServiceChecker)
	actual := svcChecker.getRdsAccountQuotas()
	assert.Len(t, actual, 1)
	accQuota := actual["foo"]
	assert.NotNil(t, accQuota)
	assert.Equal(t, aws.Int64(10), accQuota.Max)
	assert.Equal(t, aws.Int64(1), accQuota.Used)
	t.Cleanup(func() { rdsAccountQuota = map[string]*rds.AccountQuota{} })
}

func TestGetRdsAccountQuotasError(t *testing.T) {
	mockedDescribeAccountAttributesOutput := rds.DescribeAccountAttributesOutput{
		AccountQuotas: []*rds.AccountQuota{{AccountQuotaName: aws.String("foo"), Max: aws.Int64(10), Used: aws.Int64(1)}},
	}
	conf.Rds = mockedRdsClient{DescribeAccountAttributesResp: mockedDescribeAccountAttributesOutput, DescribeAccountAttributesError: errors.New("test error")}

	rdsChecker := NewRdsChecker()
	svcChecker := rdsChecker.(*ServiceChecker)
	actual := svcChecker.getRdsAccountQuotas()
	assert.Len(t, actual, 0)
	t.Cleanup(func() { rdsAccountQuota = map[string]*rds.AccountQuota{} })
}

func TestGetRdsAccountQuotasExists(t *testing.T) {
	rdsAccountQuota = map[string]*rds.AccountQuota{
		"foo": {AccountQuotaName: aws.String("foo"), Max: aws.Int64(10), Used: aws.Int64(1)},
		"bar": {AccountQuotaName: aws.String("bar"), Max: aws.Int64(100), Used: aws.Int64(5)},
	}
	rdsChecker := NewRdsChecker()
	svcChecker := rdsChecker.(*ServiceChecker)
	actual := svcChecker.getRdsAccountQuotas()
	assert.Len(t, actual, 2)
	t.Cleanup(func() { rdsAccountQuota = map[string]*rds.AccountQuota{} })
}

func TestGetRdsInstancesCountUsage(t *testing.T) {
	mockedDescribeAccountAttributesOutput := rds.DescribeAccountAttributesOutput{
		AccountQuotas: []*rds.AccountQuota{{AccountQuotaName: aws.String("DBInstances"), Max: aws.Int64(10), Used: aws.Int64(1)}},
	}
	conf.Rds = mockedRdsClient{DescribeAccountAttributesResp: mockedDescribeAccountAttributesOutput, DescribeAccountAttributesError: nil}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("rds", "DB instances", float64(20), false)},
		nil)

	rdsChecker := NewRdsChecker()
	svcChecker := rdsChecker.(*ServiceChecker)
	actual := svcChecker.getRdsInstancesCountUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "rds", quota.Service)
	assert.Equal(t, float64(20), quota.QuotaValue)
	assert.Equal(t, float64(1), quota.UsageValue)
}
