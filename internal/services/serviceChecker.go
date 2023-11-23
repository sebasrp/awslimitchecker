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
	// the applied quotas for the service. For some quotas, only default values are available
	AppliedQuotas map[string]AWSQuotaInfo
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

) ServiceQuota {

	region := ""
	if conf.Session != nil {
		region = *conf.Session.Config.Region
	}

	c := &ServiceChecker{
		ServiceCode:         serviceCode,
		Region:              region,
		AppliedQuotas:       map[string]AWSQuotaInfo{},
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

func (c ServiceChecker) GetAllAppliedQuotas() map[string]AWSQuotaInfo {
	if len(c.AppliedQuotas) == 0 {
		temp := c.getServiceAppliedQuotas()
		// sometimes applied quotas does not include all default quotas, so we need
		// to make a union between applied and default - taking applied as source
		// of truth
		for name, quota := range c.GetAllDefaultQuotas() {
			if _, ok := temp[name]; !ok {
				temp[name] = quota
			}
		}
		for key, value := range temp {
			c.AppliedQuotas[key] = value
		}
	}

	return c.AppliedQuotas
}

func (c ServiceChecker) getServiceAppliedQuotas() (ret map[string]AWSQuotaInfo) {
	ret = map[string]AWSQuotaInfo{}
	serviceQuotas := []*servicequotas.ServiceQuota{}
	err := conf.ServiceQuotas.ListServiceQuotasPages(&servicequotas.ListServiceQuotasInput{
		ServiceCode: &c.ServiceCode,
	}, func(p *servicequotas.ListServiceQuotasOutput, lastPage bool) bool {
		serviceQuotas = append(serviceQuotas, p.Quotas...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve applied quotas for service %s, %v", c.ServiceCode, err)
		return
	}

	// we then convert to our data model
	for _, q := range serviceQuotas {
		ret[aws.StringValue(q.QuotaName)] = svcQuotaToQuotaInfo(q)
	}
	return
}

func (c ServiceChecker) GetAllDefaultQuotas() map[string]AWSQuotaInfo {
	if len(c.DefaultQuotas) == 0 {
		temp := c.getServiceDefaultQuotas()
		for key, value := range temp {
			c.DefaultQuotas[key] = value
		}
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
		ret[aws.StringValue(q.QuotaName)] = svcQuotaToQuotaInfo(q)
	}
	return
}

func svcQuotaToQuotaInfo(i *servicequotas.ServiceQuota) (ret AWSQuotaInfo) {
	ret = AWSQuotaInfo{
		Service:    aws.StringValue(i.ServiceCode),
		QuotaName:  aws.StringValue(i.QuotaName),
		Quotacode:  aws.StringValue(i.QuotaCode),
		QuotaValue: aws.Float64Value(i.Value),
		UsageValue: 0.0,
		Unit:       aws.StringValue(i.Unit),
		Global:     aws.BoolValue(i.GlobalQuota),
	}
	return
}

func (c ServiceChecker) SetQuotasOverride(quotasOverride []AWSQuotaOverride) {
	for _, override := range quotasOverride {
		if c.ServiceCode != override.Service {
			return
		}
		if quota, ok := c.GetAllAppliedQuotas()[override.QuotaName]; ok {
			quota.QuotaValue = override.QuotaValue
			c.AppliedQuotas[override.QuotaName] = quota
		}
	}
}

func (c ServiceChecker) GetRequiredPermissions() []string {
	return c.RequiredPermissions
}
