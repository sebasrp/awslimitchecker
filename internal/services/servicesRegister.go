package services

import (
	"log"
)

var SupportedAwsServices = map[string]func() ServiceQuota{
	"acm":            NewAcmChecker,
	"autoscaling":    NewAutoscalingChecker,
	"cloudformation": NewCloudformationChecker,
	"dynamodb":       NewDynamoDbChecker,
	"ebs":            NewEbsChecker,
	"eks":            NewEksChecker,
	"elasticache":    NewElastiCacheChecker,
	"elb":            NewElbChecker,
	"iam":            NewIamChecker,
	"kinesis":        NewKinesisChecker,
	"rds":            NewRdsChecker,
	"s3":             NewS3Checker,
	"sns":            NewSnsChecker,
}

// GetUsage is a function that returns the usage information of a given AWS service in a given region.
// It takes three parameters: awsService, region, and overrides.
// awsService is a string that specifies the name of the AWS service to query.
// region is a string that specifies the AWS region to use.
// overrides is a slice of AWSQuotaOverride structs that defines the custom quotas to apply.
// It returns a slice of AWSQuotaInfo structs that contains the usage data for the service.

func GetUsage(service string, region string, overrides []AWSQuotaOverride) (ret []AWSQuotaInfo) {
	// Initialize the AWS session with the given region
	InitializeConfig(region)
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
