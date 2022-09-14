package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/kinesis"
)

type KinesisClientInterface interface {
	DescribeLimits(input *kinesis.DescribeLimitsInput) (*kinesis.DescribeLimitsOutput, error)
}

func NewKinesisChecker() Svcquota {
	serviceCode := "kinesis"
	supportedQuotas := map[string]func(ServiceChecker) (ret AWSQuotaInfo){
		"Shards per Region":                  ServiceChecker.getKinesisShardUsage,
		"On-demand Data Streams per account": ServiceChecker.getKinesisOnDemandStreamCountUsage,
	}
	requiredPermissions := []string{"kinesis:DescribeLimits"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getKinesisShardUsage() (ret AWSQuotaInfo) {
	result, err := conf.Kinesis.DescribeLimits(nil)
	if err != nil {
		fmt.Printf("Unable to retrieve kinesis limits, %v", err)
		return
	}

	ret = c.GetAllDefaultQuotas()["Shards per Region"]
	ret.UsageValue = float64(*result.OpenShardCount)
	return
}

func (c ServiceChecker) getKinesisOnDemandStreamCountUsage() (ret AWSQuotaInfo) {
	result, err := conf.Kinesis.DescribeLimits(nil)
	ret = AWSQuotaInfo{
		Service:   c.ServiceCode,
		Name:      "On-demand Data Streams per account",
		Region:    c.Region,
		Quotacode: "",
		Unit:      "",
		Global:    true,
	}
	if err != nil {
		fmt.Printf("Unable to retrieve kinesis limits, %v", err)
		return
	}

	// On-demand Data Streams per account is not in service quotas, so we will
	// need to create its entry in the quota list
	ret.QuotaValue = float64(*result.OnDemandStreamCountLimit)
	ret.UsageValue = float64(*result.OnDemandStreamCount)

	c.GetAllDefaultQuotas()[ret.Name] = ret
	return
}
