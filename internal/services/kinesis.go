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
		"Shards per Region": ServiceChecker.getKinesisShardUsage,
	}
	requiredPermissions := []string{"kinesis:DescribeLimits"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getKinesisShardUsage() (ret AWSQuotaInfo) {
	result, err := conf.Kinesis.DescribeLimits(nil)
	if err != nil {
		fmt.Printf("Unable to list shards, %v", err)
		return
	}

	ret = c.GetAllDefaultQuotas()["Shards per Region"]
	ret.UsageValue = float64(*result.OpenShardCount)
	return
}
