package awslimitchecker

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/sebasrp/awslimitchecker/internal/services"
)

var SupportedAwsServices = map[string]func(session *session.Session, svcQuotaClient services.SvcQuotaClientInterface) services.Svcquota{
	"s3":       services.NewS3Checker,
	"kinesis":  services.NewKinesisChecker,
	"dynamodb": services.NewDynamoDbChecker,
}

func createAwsSession(awsprofile string, region string) session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", awsprofile)},
	)
	if err != nil {
		exitErrorf("Unable to create AWS session, %v", err)
	}
	return *sess
}

func GetLimits(awsService string, awsprofile string, region string) (ret []services.AWSQuotaInfo) {
	session := createAwsSession(awsprofile, region)
	quotaClient := servicequotas.New(&session)

	if awsService == "all" {
		for _, checker := range SupportedAwsServices {
			service := checker(&session, quotaClient)
			ret = append(ret, service.GetUsage()...)
		}
	} else {
		service := SupportedAwsServices[awsService](&session, quotaClient)
		ret = service.GetUsage()
	}
	return
}

func GetIamPolicies() (ret []string) {
	for _, checker := range SupportedAwsServices {
		service := checker(nil, nil)
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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
