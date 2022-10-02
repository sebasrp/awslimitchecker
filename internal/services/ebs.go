package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func NewEbsChecker() Svcquota {
	serviceCode := "ebs"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Snapshots per Region":                                   ServiceChecker.getEbsSnapshotsUsage,
		"IOPS for Provisioned IOPS SSD (io1) volumes":            ServiceChecker.getEbsIo1IopsUsage,
		"Storage for Provisioned IOPS SSD (io1) volumes, in TiB": ServiceChecker.getEbsIo1SizeUsage,
	}
	requiredPermissions := []string{"ec2:DescribeSnapshots", "ec2:DescribeVolumes"}

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

func (c ServiceChecker) getEbsIo1IopsUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}

	iops, _, err := getEbsVolumeDetails("io1")
	if err != nil {
		fmt.Printf("failed to retrieve ec2 ebs io1 volumes, %v", err)
		return
	}
	quotaInfo := c.GetAllAppliedQuotas()["IOPS for Provisioned IOPS SSD (io1) volumes"]
	quotaInfo.UsageValue = float64(iops)

	ret = append(ret, quotaInfo)
	return
}

func (c ServiceChecker) getEbsIo1SizeUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}

	_, size, err := getEbsVolumeDetails("io1")
	if err != nil {
		fmt.Printf("failed to retrieve ec2 ebs io1 volumes, %v", err)
		return
	}
	quotaInfo := c.GetAllAppliedQuotas()["Storage for Provisioned IOPS SSD (io1) volumes, in TiB"]
	quotaInfo.UsageValue = GiBtoTiB(float64(size))

	ret = append(ret, quotaInfo)
	return
}

func getEbsVolumeDetails(volumeType string) (iops int, size int, err error) {
	iops = 0
	size = 0

	volumes := []*ec2.Volume{}
	err = conf.Ec2.DescribeVolumesPages(
		&ec2.DescribeVolumesInput{
			Filters: []*ec2.Filter{
				{Name: aws.String("volume-type"),
					Values: []*string{aws.String(volumeType)},
				}}},
		func(p *ec2.DescribeVolumesOutput, lastPage bool) bool {
			volumes = append(volumes, p.Volumes...)
			return true // continue paging
		})
	if err != nil {
		fmt.Printf("failed to retrieve ec2 ebs volumes, %v", err)
		return
	}

	for _, v := range volumes {
		iops += int(*v.Iops)
		size += int(*v.Size)
	}
	return
}

func GiBtoTiB(gib float64) (ret float64) {
	return gib / 1024
}
