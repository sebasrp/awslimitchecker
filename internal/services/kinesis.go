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

func NewKinesisChecker(session *session.Session, svcQuota *servicequotas.ServiceQuotas) Svcquota {
	c := &KinesisChecker{
		serviceCode:    "kinesis",
		region:         aws.StringValue(session.Config.Region),
		client:         kinesis.New(session),
		svcQuotaClient: svcQuota,
		defaultQuotas:  map[string]AWSQuotaInfo{},
		supportedQuotas: map[string]func(KinesisChecker) (ret AWSQuotaInfo){
			"Shards per Region": KinesisChecker.GetShardUsage,
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

func (c KinesisChecker) GetShardUsage() (ret AWSQuotaInfo) {
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
		// we first retrieve all default quotas from servicequotas
		serviceQuotas := []*servicequotas.ServiceQuota{}
		err := c.svcQuotaClient.ListAWSDefaultServiceQuotasPages(&servicequotas.ListAWSDefaultServiceQuotasInput{
			ServiceCode: &c.serviceCode,
		}, func(p *servicequotas.ListAWSDefaultServiceQuotasOutput, lastPage bool) bool {
			serviceQuotas = append(serviceQuotas, p.Quotas...)
			return true // continue paging
		})
		if err != nil {
			fmt.Printf("failed to retrieve quotas for service %s, %v", c.serviceCode, err)
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
			c.defaultQuotas[aws.StringValue(q.QuotaName)] = quota
		}
	}
	return c.defaultQuotas
}

func (c KinesisChecker) GetRequiredPermissions() []string {
	return []string{"kinesis:DescribeLimits"}
}
