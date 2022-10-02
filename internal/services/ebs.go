package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func NewEbsChecker() Svcquota {
	serviceCode := "ebs"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Snapshots per Region": ServiceChecker.getEbsSnapshotsUsage,
	}
	requiredPermissions := []string{"ec2:DescribeSnapshots"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getEbsSnapshotsUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	snapshots := []*ec2.Snapshot{}

	err := conf.Ec2.DescribeSnapshotsPages(&ec2.DescribeSnapshotsInput{OwnerIds: []*string{aws.String("self")}}, func(p *ec2.DescribeSnapshotsOutput, lastPage bool) bool {
		snapshots = append(snapshots, p.Snapshots...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve ec2 ebs snapshots, %v", err)
		return
	}
	quotaInfo := c.GetAllAppliedQuotas()["Snapshots per Region"]
	quotaInfo.UsageValue = float64(len(snapshots))

	ret = append(ret, quotaInfo)
	return
}
