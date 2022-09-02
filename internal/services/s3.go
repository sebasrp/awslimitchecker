package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
)

type S3Checker struct {
	// serviceCode is the name of the service this checker verifies
	serviceCode string
	// region the checker will run against
	region string
	// aws client used to call kinesis service
	client s3iface.S3API
	// aws client used to call service quotas service
	svcQuotaClient servicequotasiface.ServiceQuotasAPI
	// the default quotas of the service
	defaultQuotas map[string]AWSQuotaInfo
	// supportedQuotas contains the service quota name and the func used to retrieve its usage
	supportedQuotas map[string]func(S3Checker) (ret AWSQuotaInfo)
}

func NewS3Checker(session *session.Session, svcQuota *servicequotas.ServiceQuotas) Svcquota {
	serviceCode := "s3"
	region := ""
	var client s3iface.S3API
	if session != nil {
		region = aws.StringValue(session.Config.Region)
		client = s3.New(session)
	}
	c := &S3Checker{
		serviceCode:    serviceCode,
		region:         region,
		client:         client,
		svcQuotaClient: svcQuota,
		defaultQuotas:  map[string]AWSQuotaInfo{},
		supportedQuotas: map[string]func(S3Checker) (ret AWSQuotaInfo){
			"Buckets": S3Checker.getBucketUsage},
	}
	return c
}

func (c S3Checker) GetUsage() (ret []AWSQuotaInfo) {
	for _, q := range c.supportedQuotas {
		quotaInfo := q(c)
		ret = append(ret, quotaInfo)
	}
	return
}

func (c S3Checker) getBucketUsage() (ret AWSQuotaInfo) {
	result, err := c.client.ListBuckets(nil)
	if err != nil {
		fmt.Printf("Unable to list buckets, %v", err)
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

	ret = c.GetAllDefaultQuotas()["Buckets"]
	ret.UsageValue = float64(len(result.Buckets))
	return
}

func (c S3Checker) GetAllDefaultQuotas() map[string]AWSQuotaInfo {
	if len(c.defaultQuotas) == 0 {
		c.defaultQuotas = GetServiceDefaultQuotas(c.serviceCode, c.region, c.svcQuotaClient)
	}
	return c.defaultQuotas
}

func (c S3Checker) GetRequiredPermissions() []string {
	return []string{"s3:ListAllMyBuckets"}
}
