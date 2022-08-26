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

var SupportedAwsServices = map[string]func(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo){
	"s3":      GetS3Usage,
	"kinesis": GetKinesisUsage,
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

func GetLimits(awsService string, awsprofile string, region string) {
	fmt.Printf("AWS profile: %s | AWS region: %s | service: %s\n", awsprofile, region, awsService)
	session := createAwsSession(awsprofile, region)
	quotaClient := servicequotas.New(&session)

	usage := []services.AWSQuotaInfo{}
	for _, checker := range SupportedAwsServices {
		usage = append(usage, checker(session, quotaClient)...)
	}
	for _, u := range usage {
		fmt.Printf("* %s %g/%g\n",
			u.Name, u.UsageValue, u.QuotaValue)
	}
}

func GetS3Usage(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo) {
	s3checker := services.NewS3Checker(session, quotaClient)
	ret = s3checker.GetUsage()
	return
}

func GetKinesisUsage(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo) {
	kinesischecker := services.NewKinesisChecker(session, quotaClient)
	ret = kinesischecker.GetUsage()
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
