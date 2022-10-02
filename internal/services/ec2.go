package services

import "github.com/aws/aws-sdk-go/service/ec2"

type Ec2ClientInterface interface {
	DescribeSnapshotsPages(input *ec2.DescribeSnapshotsInput, fn func(*ec2.DescribeSnapshotsOutput, bool) bool) error
}
