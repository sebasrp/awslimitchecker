package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type Elbv2ClientInterface interface {
	DescribeAccountLimits(input *elbv2.DescribeAccountLimitsInput) (*elbv2.DescribeAccountLimitsOutput, error)
	DescribeLoadBalancersPages(input *elbv2.DescribeLoadBalancersInput, fn func(*elbv2.DescribeLoadBalancersOutput, bool) bool) error
}

func NewElbChecker() Svcquota {
	serviceCode := "elasticloadbalancing"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Application Load Balancers per Region": ServiceChecker.getElbApplicationLoadBalancerUsage,
	}
	requiredPermissions := []string{
		"elasticloadbalancing:DescribeLoadBalancers",
		"elasticloadbalancing:DescribeAccountLimits",
	}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

var elbAccountQuota map[string]*elbv2.Limit = map[string]*elbv2.Limit{}

func (c ServiceChecker) getElbAccountQuotas() (ret map[string]*elbv2.Limit) {
	ret = elbAccountQuota
	if len(elbAccountQuota) != 0 {
		return
	}

	result, err := conf.Elbv2.DescribeAccountLimits(nil)
	if err != nil {
		fmt.Printf("Unable to retrieve elb account attributes, %v", err)
		return
	}

	for _, q := range result.Limits {
		elbAccountQuota[aws.StringValue(q.Name)] = q
	}
	return
}

func (c ServiceChecker) getElbApplicationLoadBalancerUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["Application Load Balancers per Region"]

	// we need to iterate through all LBs and check which ones are NLB vs ALB
	albs := []*elbv2.LoadBalancer{}
	err := conf.Elbv2.DescribeLoadBalancersPages(&elbv2.DescribeLoadBalancersInput{}, func(p *elbv2.DescribeLoadBalancersOutput, lastPage bool) bool {
		for _, q := range p.LoadBalancers {
			if *q.Type == "application" {
				albs = append(albs, q)
			}
		}
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve load balancers, %v", err)
		return
	}

	// we then get the quota info from the service itself (overwrites servicequotas')
	if val, ok := c.getElbAccountQuotas()["application-load-balancers"]; ok {
		quotaFloat, _ := strconv.ParseFloat(strings.TrimSpace(*val.Max), 64)
		quotaInfo.QuotaValue = quotaFloat
	}

	quotaInfo.UsageValue = float64(len(albs))
	ret = append(ret, quotaInfo)
	return
}
