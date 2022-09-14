package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
)

type S3ClientInterface interface {
	ListBuckets(input *s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
}

func NewS3Checker() Svcquota {
	serviceCode := "s3"
	supportedQuotas := map[string]func(ServiceChecker) (ret AWSQuotaInfo){
		"Buckets": ServiceChecker.getS3BucketUsage,
	}
	requiredPermissions := []string{"s3:ListAllMyBuckets"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getS3BucketUsage() (ret AWSQuotaInfo) {
	result, err := conf.S3.ListBuckets(nil)
	ret = c.GetAllDefaultQuotas()["Buckets"]
	if err != nil {
		fmt.Printf("Unable to list buckets, %v", err)
		return
	}

	/* breakdown per region.
	// TODO: figure out how to make this more efficient / another API?
	regionalBucket := []s3.Bucket{}
	for _, b := range result.Buckets {
		loc, err := c.client.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: b.Name})
		if err != nil {
			fmt.Printf("Unable to retrieve location for bucket, %s", *b.Name)
		}
		locString := "us-east-1" // the location string is nil for us-east-1, see https://github.com/aws/aws-sdk-go-v2/blob/service/s3/v1.27.5/service/s3/api_op_GetBucketLocation.go#L72
		if loc.LocationConstraint != nil {
			locString = *loc.LocationConstraint
		}

		if locString == c.region {
			fmt.Printf("[%s] Bucket %s\n", locString, *b.Name)
			regionalBucket = append(regionalBucket, []s3.Bucket{*b}...)
		}

	}
	fmt.Printf("regionalBucket: %v\n", regionalBucket)
	*/
	ret.UsageValue = float64(len(result.Buckets))
	return
}
