package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
)

type AWSQuotaInfo struct {
	Service    string  // service the quota applies to
	Name       string  // the name of the aws service resource the usage is for
	Region     string  // the region this quota applies to
	Quotacode  string  // servicequota code
	QuotaValue float64 // the quota value
	UsageValue float64 // the usage value
	Unit       string  // unit of the quota/usage
	Global     bool    // whether the quota is global or not
}

type Svcquota interface {
	// Get Usage retrieve the quotas and usage for the given service
	GetUsage() []AWSQuotaInfo

	// GetAllDefaultQuotas retrieves all the default quotas for the given
	// service. Usage of those resources are not retrieved/calculated
	GetAllDefaultQuotas() map[string]AWSQuotaInfo

	// GetRequiredPermissions returns a list of the IAM permissions required
	// to retrieve the usage for this service.
	GetRequiredPermissions() []string
}

func GetServiceDefaultQuotas(serviceCode string, region string, svcQuota servicequotasiface.ServiceQuotasAPI) (ret map[string]AWSQuotaInfo) {
	ret = map[string]AWSQuotaInfo{}
	serviceQuotas := []*servicequotas.ServiceQuota{}
	err := svcQuota.ListAWSDefaultServiceQuotasPages(&servicequotas.ListAWSDefaultServiceQuotasInput{
		ServiceCode: &serviceCode,
	}, func(p *servicequotas.ListAWSDefaultServiceQuotasOutput, lastPage bool) bool {
		serviceQuotas = append(serviceQuotas, p.Quotas...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve quotas for service %s, %v", serviceCode, err)
	}

	// we then convert to our data model
	for _, q := range serviceQuotas {
		quota := AWSQuotaInfo{
			Service:    serviceCode,
			Name:       aws.StringValue(q.QuotaName),
			Region:     region,
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
