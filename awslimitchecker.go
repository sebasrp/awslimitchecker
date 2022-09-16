package awslimitchecker

import (
	"fmt"

	"github.com/sebasrp/awslimitchecker/internal/services"
)

var SupportedAwsServices = map[string]func() services.Svcquota{
	"dynamodb":    services.NewDynamoDbChecker,
	"eks":         services.NewEksChecker,
	"elasticache": services.NewElastiCacheChecker,
	"kinesis":     services.NewKinesisChecker,
	"rds":         services.NewRdsChecker,
	"s3":          services.NewS3Checker,
}

func GetLimits(awsService string, awsprofile string, region string) (ret []services.AWSQuotaInfo) {
	_, err := services.InitializeConfig(awsprofile, region)
	if err != nil {
		fmt.Printf("Unable to create AWS session, %v", err)
		return
	}

	if awsService == "all" {
		for _, checker := range SupportedAwsServices {
			service := checker()
			ret = append(ret, service.GetUsage()...)
		}
	} else if val, ok := SupportedAwsServices[awsService]; ok {
		service := val()
		ret = service.GetUsage()
	}
	return
}

func GetIamPolicies() (ret []string) {
	for _, checker := range SupportedAwsServices {
		service := checker()
		ret = append(ret, service.GetRequiredPermissions()...)
	}
	return
}

func IsValidAwsService(service string) bool {
	if _, ok := SupportedAwsServices[service]; ok || service == "all" {
		return true
	} else {
		return false
	}
}
