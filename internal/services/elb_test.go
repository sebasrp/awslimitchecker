package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedElbv2Client struct {
	Elbv2ClientInterface
	DescribeAccountLimitsResp       elbv2.DescribeAccountLimitsOutput
	DescribeAccountLimitsError      error
	DescribeLoadBalancersPagesRest  elbv2.DescribeLoadBalancersOutput
	DescribeLoadBalancersPagesError error
}

func (m mockedElbv2Client) DescribeAccountLimits(input *elbv2.DescribeAccountLimitsInput) (*elbv2.DescribeAccountLimitsOutput, error) {
	return &m.DescribeAccountLimitsResp, m.DescribeAccountLimitsError
}

func (m mockedElbv2Client) DescribeLoadBalancersPages(input *elbv2.DescribeLoadBalancersInput, fn func(*elbv2.DescribeLoadBalancersOutput, bool) bool) error {
	fn(&m.DescribeLoadBalancersPagesRest, false)
	return m.DescribeLoadBalancersPagesError
}

func TestNewElbCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewElbChecker())
}

func TestGetElbAccountQuotas(t *testing.T) {
	mockedDescribeAccountLimitsOutput := elbv2.DescribeAccountLimitsOutput{
		Limits: []*elbv2.Limit{{Name: aws.String("foo"), Max: aws.String("100")}},
	}
	conf.Elbv2 = mockedElbv2Client{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutput}

	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbAccountQuotas()
	assert.Len(t, actual, 1)
	accQuota := actual["foo"]
	assert.NotNil(t, accQuota)
	assert.Equal(t, aws.String("100"), accQuota.Max)
	t.Cleanup(func() { elbAccountQuota = map[string]*elbv2.Limit{} })
}

func TestGetElbAccountQuotasError(t *testing.T) {
	mockedDescribeAccountLimitsOutput := elbv2.DescribeAccountLimitsOutput{
		Limits: []*elbv2.Limit{{Name: aws.String("foo"), Max: aws.String("100")}},
	}
	conf.Elbv2 = mockedElbv2Client{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutput, DescribeAccountLimitsError: errors.New("test error")}

	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbAccountQuotas()
	assert.Len(t, actual, 0)
	t.Cleanup(func() { elbAccountQuota = map[string]*elbv2.Limit{} })
}

func TestGetElbAccountQuotasExists(t *testing.T) {
	elbAccountQuota = map[string]*elbv2.Limit{
		"foo": {Name: aws.String("foo"), Max: aws.String("10")},
		"bar": {Name: aws.String("bar"), Max: aws.String("100")},
	}
	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbAccountQuotas()
	assert.Len(t, actual, 2)
	t.Cleanup(func() { elbAccountQuota = map[string]*elbv2.Limit{} })
}

func TestGetElbApplicationLoadBalancerUsage(t *testing.T) {
	mockedDescribeAccountLimitsOutput := elbv2.DescribeAccountLimitsOutput{
		Limits: []*elbv2.Limit{{Name: aws.String("application-load-balancers"), Max: aws.String("200")}},
	}
	mockedDescribeLoadBalancersOutput := elbv2.DescribeLoadBalancersOutput{
		LoadBalancers: []*elbv2.LoadBalancer{
			{LoadBalancerName: aws.String("foo"), Type: aws.String("network")},
			{LoadBalancerName: aws.String("foo"), Type: aws.String("application")},
			{LoadBalancerName: aws.String("foo"), Type: aws.String("gateway")},
		},
	}
	conf.Elbv2 = mockedElbv2Client{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutput, DescribeLoadBalancersPagesRest: mockedDescribeLoadBalancersOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("elasticloadbalancing", "Application Load Balancers per Region", float64(100), false)},
		nil)

	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbApplicationLoadBalancerUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "elasticloadbalancing", quota.Service)
	assert.Equal(t, float64(200), quota.QuotaValue)
	assert.Equal(t, float64(1), quota.UsageValue)
	t.Cleanup(func() { elbAccountQuota = map[string]*elbv2.Limit{} })
}

func TestGetElbApplicationLoadBalancerUsageError(t *testing.T) {
	mockedDescribeAccountLimitsOutput := elbv2.DescribeAccountLimitsOutput{
		Limits: []*elbv2.Limit{{Name: aws.String("application-load-balancers"), Max: aws.String("200")}},
	}
	mockedDescribeLoadBalancersOutput := elbv2.DescribeLoadBalancersOutput{
		LoadBalancers: []*elbv2.LoadBalancer{
			{LoadBalancerName: aws.String("foo"), Type: aws.String("network")},
			{LoadBalancerName: aws.String("foo"), Type: aws.String("application")},
			{LoadBalancerName: aws.String("foo"), Type: aws.String("gateway")},
		},
	}
	conf.Elbv2 = mockedElbv2Client{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutput, DescribeLoadBalancersPagesRest: mockedDescribeLoadBalancersOutput, DescribeLoadBalancersPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("elasticloadbalancing", "Application Load Balancers per Region", float64(100), false)},
		nil)

	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbApplicationLoadBalancerUsage()

	assert.Len(t, actual, 0)
	t.Cleanup(func() { elbAccountQuota = map[string]*elbv2.Limit{} })
}
