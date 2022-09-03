package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
)

type KinesisChecker struct {
	// serviceCode is the name of the service this checker verifies
	serviceCode string
	// region the checker will run against
	region string
	// aws client used to call kinesis service
	client kinesisiface.KinesisAPI
	// aws client used to call service quotas service
	svcQuotaClient servicequotasiface.ServiceQuotasAPI
	// the default quotas of the service
	defaultQuotas map[string]AWSQuotaInfo
	// supportedQuotas contains the service quota name and the func used to retrieve its usage
	supportedQuotas map[string]func(KinesisChecker) (ret AWSQuotaInfo)
}

func NewKinesisChecker(session *session.Session) Svcquota {
	serviceCode := "kinesis"
	region := ""
	var client kinesisiface.KinesisAPI
	var svcQuota servicequotasiface.ServiceQuotasAPI

	if session != nil {
		region = aws.StringValue(session.Config.Region)
		client = kinesis.New(session)
		svcQuota = servicequotas.New(session)
	}
	c := &KinesisChecker{
		serviceCode:    serviceCode,
		region:         region,
		client:         client,
		svcQuotaClient: svcQuota,
		defaultQuotas:  map[string]AWSQuotaInfo{},
		supportedQuotas: map[string]func(KinesisChecker) (ret AWSQuotaInfo){
			"Shards per Region": KinesisChecker.getShardUsage,
		},
	}
	return c
}

func (c KinesisChecker) GetUsage() (ret []AWSQuotaInfo) {
	for _, q := range c.supportedQuotas {
		quotaInfo := q(c)
		ret = append(ret, quotaInfo)
	}
	return
}

func (c KinesisChecker) getShardUsage() (ret AWSQuotaInfo) {
	result, err := c.client.DescribeLimits(nil)
	if err != nil {
		fmt.Printf("Unable to list shards, %v", err)
	}

	ret = c.GetAllDefaultQuotas()["Shards per Region"]
	ret.UsageValue = float64(*result.OpenShardCount)
	return
}

func (c KinesisChecker) GetAllDefaultQuotas() map[string]AWSQuotaInfo {
	if len(c.defaultQuotas) == 0 {
		c.defaultQuotas = GetServiceDefaultQuotas(c.serviceCode, c.region, c.svcQuotaClient)
	}
	return c.defaultQuotas
}

func (c KinesisChecker) GetRequiredPermissions() []string {
	return []string{"kinesis:DescribeLimits"}
}
