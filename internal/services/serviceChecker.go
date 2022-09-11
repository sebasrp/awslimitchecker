package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

type ServiceChecker struct {
	// serviceCode is the name of the service this checker verifies
	serviceCode string
	// region the checker will run against
	region string
	// the default quotas of the service
	defaultQuotas map[string]AWSQuotaInfo
	// supportedQuotas contains the service quota name and the func used to retrieve its usage
	supportedQuotas map[string]func(ServiceChecker) (ret AWSQuotaInfo)
	// Permissions required to get usage
	requiredPermissions []string
}

func NewServiceChecker(
	serviceCode string,
	quotas map[string]func(ServiceChecker) (ret AWSQuotaInfo),
	permissions []string,

) Svcquota {

	region := ""
	if conf.Session != nil {
		region = *conf.Session.Config.Region
	}

	c := &ServiceChecker{
		serviceCode:         serviceCode,
		region:              region,
		defaultQuotas:       map[string]AWSQuotaInfo{},
		supportedQuotas:     quotas,
		requiredPermissions: permissions,
	}
	return c
}

func (c ServiceChecker) GetUsage() (ret []AWSQuotaInfo) {
	for _, q := range c.supportedQuotas {
		quotaInfo := q(c)
		ret = append(ret, quotaInfo)
	}
	return
}

func (c ServiceChecker) GetAllDefaultQuotas() map[string]AWSQuotaInfo {
	if len(c.defaultQuotas) == 0 {
		c.defaultQuotas = c.getServiceDefaultQuotas()
	}
	return c.defaultQuotas
}

func (c ServiceChecker) getServiceDefaultQuotas() (ret map[string]AWSQuotaInfo) {
	ret = map[string]AWSQuotaInfo{}
	serviceQuotas := []*servicequotas.ServiceQuota{}
	err := conf.ServiceQuotas.ListAWSDefaultServiceQuotasPages(&servicequotas.ListAWSDefaultServiceQuotasInput{
		ServiceCode: &c.serviceCode,
	}, func(p *servicequotas.ListAWSDefaultServiceQuotasOutput, lastPage bool) bool {
		serviceQuotas = append(serviceQuotas, p.Quotas...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve quotas for service %s, %v", c.serviceCode, err)
		return
	}

	// we then convert to our data model
	for _, q := range serviceQuotas {
		quota := AWSQuotaInfo{
			Service:    c.serviceCode,
			Name:       aws.StringValue(q.QuotaName),
			Region:     c.region,
			Quotacode:  aws.StringValue(q.QuotaCode),
			QuotaValue: aws.Float64Value(q.Value),
			UsageValue: 0.0,
			Unit:       aws.StringValue(q.Unit),
			Global:     aws.BoolValue(q.GlobalQuota),
		}
		ret[aws.StringValue(q.QuotaName)] = quota
	}
	return
}

func (c ServiceChecker) GetRequiredPermissions() []string {
	return c.requiredPermissions
}
