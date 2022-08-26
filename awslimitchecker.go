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

var SupportedAwsServices = map[string]bool{
	"s3": true,
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
	s3checker := services.NewS3Checker(session, quotaClient)
	usage := s3checker.GetUsage()
	for _, u := range usage {
		fmt.Printf("* %s %g/%g\n",
			u.Name, u.UsageValue, u.QuotaValue)
	}
}

func IsValidAwsService(service string) bool {
	return SupportedAwsServices[service] || service == "all"
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
