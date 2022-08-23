package awslimitchecker

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var SupportedAwsServices = map[string]bool{
	"s3": true,
}

func GetLimits(awsService string, awsprofile string, region string) {
	fmt.Printf("AWS profile: %s | AWS region: %s | service: %s", awsprofile, region, awsService)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", awsprofile)},
	)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}

func IsValidAwsService(service string) bool {
	return SupportedAwsServices[service] || service == "all"
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
