package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
)

type ServiceChecker struct {
	// ServiceCode is the name of the service this checker verifies
	ServiceCode string
	// Region the checker will run against
	Region string
	// the default quotas of the service
	DefaultQuotas map[string]AWSQuotaInfo
	// SupportedQuotas contains the service quota name and the func used to retrieve its usage
	SupportedQuotas map[string]func(ServiceChecker) (ret []AWSQuotaInfo)
	// Permissions required to get usage
	RequiredPermissions []string
}

func NewServiceChecker(
	serviceCode string,
	quotas map[string]func(ServiceChecker) (ret []AWSQuotaInfo),
	permissions []string,

) Svcquota {

	region := ""
	if conf.Session != nil {
		region = *conf.Session.Config.Region
	}

	c := &ServiceChecker{
		ServiceCode:         serviceCode,
		Region:              region,
		DefaultQuotas:       map[string]AWSQuotaInfo{},
		SupportedQuotas:     quotas,
		RequiredPermissions: permissions,
	}
	return c
}

func (c ServiceChecker) GetUsage() (ret []AWSQuotaInfo) {
	for _, q := range c.SupportedQuotas {
		quotaInfo := q(c)
		ret = append(ret, quotaInfo...)
	}
	return
}

func (c ServiceChecker) GetAllDefaultQuotas() map[string]AWSQuotaInfo {
	if len(c.DefaultQuotas) == 0 {
		c.DefaultQuotas = c.getServiceDefaultQuotas()
	}
	return c.DefaultQuotas
}

func (c ServiceChecker) getServiceDefaultQuotas() (ret map[string]AWSQuotaInfo) {
	ret = map[string]AWSQuotaInfo{}
	serviceQuotas := []*servicequotas.ServiceQuota{}
	err := conf.ServiceQuotas.ListAWSDefaultServiceQuotasPages(&servicequotas.ListAWSDefaultServiceQuotasInput{
		ServiceCode: &c.ServiceCode,
	}, func(p *servicequotas.ListAWSDefaultServiceQuotasOutput, lastPage bool) bool {
		serviceQuotas = append(serviceQuotas, p.Quotas...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve quotas for service %s, %v", c.ServiceCode, err)
		return
	}

	// we then convert to our data model
	for _, q := range serviceQuotas {
		quota := AWSQuotaInfo{
			Service:    c.ServiceCode,
			Name:       aws.StringValue(q.QuotaName),
			Region:     c.Region,
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
	return c.RequiredPermissions
}
