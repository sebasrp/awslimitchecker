package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

func NewKinesisChecker(session *session.Session, svcQuotaClient SvcQuotaClientInterface) Svcquota {
	serviceCode := "kinesis"
	supportedQuotas := map[string]func(ServiceChecker) (ret AWSQuotaInfo){
		"Shards per Region": ServiceChecker.getKinesisShardUsage,
	}
	requiredPermissions := []string{"kinesis:DescribeLimits"}

	return NewServiceChecker(serviceCode, session, svcQuotaClient, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getKinesisShardUsage() (ret AWSQuotaInfo) {
	kinesisClient := kinesis.New(c.session)
	result, err := kinesisClient.DescribeLimits(nil)
	if err != nil {
		fmt.Printf("Unable to list shards, %v", err)
	}

	ret = c.GetAllDefaultQuotas()["Shards per Region"]
	ret.UsageValue = float64(*result.OpenShardCount)
	return
}
