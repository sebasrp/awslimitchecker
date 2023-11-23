package awslimitchecker

import (
	"log"

	"github.com/nyambati/aws-service-limits-exporter/internal/services"
)

var SupportedAwsServices = map[string]func() services.ServiceQuota{
	"acm":            services.NewAcmChecker,
	"autoscaling":    services.NewAutoscalingChecker,
	"cloudformation": services.NewCloudformationChecker,
	"dynamodb":       services.NewDynamoDbChecker,
	"ebs":            services.NewEbsChecker,
	"eks":            services.NewEksChecker,
	"elasticache":    services.NewElastiCacheChecker,
	"elb":            services.NewElbChecker,
	"iam":            services.NewIamChecker,
	"kinesis":        services.NewKinesisChecker,
	"rds":            services.NewRdsChecker,
	"s3":             services.NewS3Checker,
	"sns":            services.NewSnsChecker,
}

// GetUsage is a function that returns the usage information of a given AWS service in a given region.
// It takes three parameters: awsService, region, and overrides.
// awsService is a string that specifies the name of the AWS service to query.
// region is a string that specifies the AWS region to use.
// overrides is a slice of AWSQuotaOverride structs that defines the custom quotas to apply.
// It returns a slice of AWSQuotaInfo structs that contains the usage data for the service.

func GetUsage(service string, region string, overrides []services.AWSQuotaOverride) (ret []services.AWSQuotaInfo) {
	// Initialize the AWS session with the given region
	err := services.InitializeConfig(region)
	if err != nil {
		// Log the error and return
		log.Printf("Unable to create AWS session, %v", err)
		return
	}
	// Check the value of awsService and create the corresponding service instance
	switch getServiceChecker, ok := SupportedAwsServices[service]; {
	case ok:
		// Create the service instance using the value function
		serviceQuota := getServiceChecker()
		// Set the quotas override for the service
		serviceQuota.SetQuotasOverride(overrides)
		// Get the usage data for the service
		ret = serviceQuota.GetUsage()
	default:
		// Log the error and return
		log.Printf("Unsupported AWS service, %v", service)
		return
	}
	// Return the usage data
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
	_, ok := SupportedAwsServices[service]
	return ok
}
