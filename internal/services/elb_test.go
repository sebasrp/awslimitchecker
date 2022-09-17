package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedElbClient struct {
	ElbClientInterface
	DescribeAccountLimitsResp       elb.DescribeAccountLimitsOutput
	DescribeAccountLimitsError      error
	DescribeLoadBalancersPagesRest  elb.DescribeLoadBalancersOutput
	DescribeLoadBalancersPagesError error
}

func (m mockedElbClient) DescribeAccountLimits(input *elb.DescribeAccountLimitsInput) (*elb.DescribeAccountLimitsOutput, error) {
	return &m.DescribeAccountLimitsResp, m.DescribeAccountLimitsError
}

func (m mockedElbClient) DescribeLoadBalancersPages(input *elb.DescribeLoadBalancersInput, fn func(*elb.DescribeLoadBalancersOutput, bool) bool) error {
	fn(&m.DescribeLoadBalancersPagesRest, false)
	return m.DescribeLoadBalancersPagesError
}

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
	mockedDescribeAccountLimitsOutputv2 := elbv2.DescribeAccountLimitsOutput{
		Limits: []*elbv2.Limit{{Name: aws.String("foo"), Max: aws.String("100")}},
	}
	conf.Elbv2 = mockedElbv2Client{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutputv2}
	mockedDescribeAccountLimitsOutput := elb.DescribeAccountLimitsOutput{
		Limits: []*elb.Limit{{Name: aws.String("bar"), Max: aws.String("1000")}},
	}
	conf.Elb = mockedElbClient{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutput}

	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbAccountQuotas()
	assert.Len(t, actual, 2)
	foo := actual["foo"]
	assert.Equal(t, float64(100), foo)
	bar := actual["bar"]
	assert.Equal(t, float64(1000), bar)
	t.Cleanup(func() { elbAccountQuota = map[string]float64{} })
}

func TestGetElbAccountQuotasErrorv2(t *testing.T) {
	mockedDescribeAccountLimitsOutputv2 := elbv2.DescribeAccountLimitsOutput{
		Limits: []*elbv2.Limit{{Name: aws.String("foo"), Max: aws.String("100")}},
	}
	conf.Elbv2 = mockedElbv2Client{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutputv2, DescribeAccountLimitsError: errors.New("test error")}
	mockedDescribeAccountLimitsOutput := elb.DescribeAccountLimitsOutput{
		Limits: []*elb.Limit{{Name: aws.String("bar"), Max: aws.String("1000")}},
	}
	conf.Elb = mockedElbClient{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutput}

	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbAccountQuotas()
	assert.Len(t, actual, 0)
	t.Cleanup(func() { elbAccountQuota = map[string]float64{} })
}

func TestGetElbAccountQuotasErrorClassic(t *testing.T) {
	mockedDescribeAccountLimitsOutputv2 := elbv2.DescribeAccountLimitsOutput{
		Limits: []*elbv2.Limit{{Name: aws.String("foo"), Max: aws.String("100")}},
	}
	conf.Elbv2 = mockedElbv2Client{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutputv2}
	mockedDescribeAccountLimitsOutput := elb.DescribeAccountLimitsOutput{
		Limits: []*elb.Limit{{Name: aws.String("bar"), Max: aws.String("1000")}},
	}
	conf.Elb = mockedElbClient{DescribeAccountLimitsResp: mockedDescribeAccountLimitsOutput, DescribeAccountLimitsError: errors.New("test error")}

	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbAccountQuotas()
	assert.Len(t, actual, 0)
	t.Cleanup(func() { elbAccountQuota = map[string]float64{} })
}

func TestGetElbAccountQuotasExists(t *testing.T) {
	elbAccountQuota = map[string]float64{
		"foo": float64(10),
		"bar": float64(100),
	}
	elbChecker := NewElbChecker()
	svcChecker := elbChecker.(*ServiceChecker)
	actual := svcChecker.getElbAccountQuotas()
	assert.Len(t, actual, 2)
	t.Cleanup(func() { elbAccountQuota = map[string]float64{} })
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
	conf.Elb = mockedElbClient{DescribeAccountLimitsResp: elb.DescribeAccountLimitsOutput{}}

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
	t.Cleanup(func() { elbAccountQuota = map[string]float64{} })
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
	t.Cleanup(func() { elbAccountQuota = map[string]float64{} })
}
