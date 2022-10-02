package services

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type mockedEc2Client struct {
	Ec2ClientInterface
	DescribeSnapshotsPagesResp  ec2.DescribeSnapshotsOutput
	DescribeSnapshotsPagesError error
}

func (m mockedEc2Client) DescribeSnapshotsPages(input *ec2.DescribeSnapshotsInput, fn func(*ec2.DescribeSnapshotsOutput, bool) bool) error {
	fn(&m.DescribeSnapshotsPagesResp, false)
	return m.DescribeSnapshotsPagesError
}
