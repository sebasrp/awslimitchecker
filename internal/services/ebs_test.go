package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEbsCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewEbsChecker())
}

func TestGetEbsSnapshotsUsage(t *testing.T) {
	mockedOutput := ec2.DescribeSnapshotsOutput{
		Snapshots: []*ec2.Snapshot{{SnapshotId: aws.String("foo")}},
	}
	conf.Ec2 = mockedEc2Client{DescribeSnapshotsPagesResp: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("ebs", "Snapshots per Region", float64(100), false)},
		nil)

	ebsChecker := NewEbsChecker()
	svcChecker := ebsChecker.(*ServiceChecker)
	actual := svcChecker.getEbsSnapshotsUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "ebs", quota.Service)
	assert.Equal(t, float64(100), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.Snapshots)), quota.UsageValue)
}

func TestGetEbsSnapshotsUsagerror(t *testing.T) {
	mockedOutput := ec2.DescribeSnapshotsOutput{
		Snapshots: []*ec2.Snapshot{{SnapshotId: aws.String("foo")}},
	}
	conf.Ec2 = mockedEc2Client{DescribeSnapshotsPagesResp: mockedOutput, DescribeSnapshotsPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("ebs", "Snapshots per Region", float64(100), false)},
		nil)

	ebsChecker := NewEbsChecker()
	svcChecker := ebsChecker.(*ServiceChecker)
	actual := svcChecker.getEbsSnapshotsUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}

func TestGetEbsIo1IopsUsage(t *testing.T) {
	mockedOutput := ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			{VolumeId: aws.String("foo"), Iops: aws.Int64(1000), Size: aws.Int64(1024)},
			{VolumeId: aws.String("bar"), Iops: aws.Int64(1000), Size: aws.Int64(2048)},
		},
	}
	conf.Ec2 = mockedEc2Client{DescribeVolumesPagesRes: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("ebs", "IOPS for Provisioned IOPS SSD (io1) volumes", float64(10000), false)},
		nil)

	ebsChecker := NewEbsChecker()
	svcChecker := ebsChecker.(*ServiceChecker)
	actual := svcChecker.getEbsIo1IopsUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "ebs", quota.Service)
	assert.Equal(t, float64(10000), quota.QuotaValue)
	assert.Equal(t, float64(2000), quota.UsageValue)
}
func TestGetEbsIo1IopsUsageError(t *testing.T) {
	mockedOutput := ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			{VolumeId: aws.String("foo"), Iops: aws.Int64(1000), Size: aws.Int64(1024)},
			{VolumeId: aws.String("bar"), Iops: aws.Int64(1000), Size: aws.Int64(2048)},
		},
	}
	conf.Ec2 = mockedEc2Client{DescribeVolumesPagesRes: mockedOutput, DescribeVolumesPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("ebs", "IOPS for Provisioned IOPS SSD (io1) volumes", float64(10000), false)},
		nil)

	ebsChecker := NewEbsChecker()
	svcChecker := ebsChecker.(*ServiceChecker)
	actual := svcChecker.getEbsIo1IopsUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}

func TestGetEbsIo1SizeUsage(t *testing.T) {
	mockedOutput := ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			{VolumeId: aws.String("foo"), Iops: aws.Int64(1000), Size: aws.Int64(1024)},
			{VolumeId: aws.String("bar"), Iops: aws.Int64(1000), Size: aws.Int64(2048)},
		},
	}
	conf.Ec2 = mockedEc2Client{DescribeVolumesPagesRes: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("ebs", "Storage for Provisioned IOPS SSD (io1) volumes, in TiB", float64(50), false)},
		nil)

	ebsChecker := NewEbsChecker()
	svcChecker := ebsChecker.(*ServiceChecker)
	actual := svcChecker.getEbsIo1SizeUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "ebs", quota.Service)
	assert.Equal(t, float64(50), quota.QuotaValue)
	assert.Equal(t, float64(3), quota.UsageValue)
}
func TestGetEbsIo1SizeUsageError(t *testing.T) {
	mockedOutput := ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			{VolumeId: aws.String("foo"), Iops: aws.Int64(1000), Size: aws.Int64(1024)},
			{VolumeId: aws.String("bar"), Iops: aws.Int64(1000), Size: aws.Int64(2048)},
		},
	}
	conf.Ec2 = mockedEc2Client{DescribeVolumesPagesRes: mockedOutput, DescribeVolumesPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("ebs", "Storage for Provisioned IOPS SSD (io1) volumes, in TiB", float64(50), false)},
		nil)

	ebsChecker := NewEbsChecker()
	svcChecker := ebsChecker.(*ServiceChecker)
	actual := svcChecker.getEbsIo1SizeUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}

func TestGetEbsVolumeDetails(t *testing.T) {
	mockedOutput := ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			{VolumeId: aws.String("foo"), Iops: aws.Int64(1000), Size: aws.Int64(1024)},
			{VolumeId: aws.String("bar"), Iops: aws.Int64(1000), Size: aws.Int64(2048)},
		},
	}
	conf.Ec2 = mockedEc2Client{DescribeVolumesPagesRes: mockedOutput}

	iops, size, err := getEbsVolumeDetails("standard")
	assert.Equal(t, 2000, iops)
	assert.Equal(t, 3072, size)
	assert.Nil(t, err)
}

func TestGetEbsVolumeDetailsError(t *testing.T) {
	mockedOutput := ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			{VolumeId: aws.String("foo"), Iops: aws.Int64(1000), Size: aws.Int64(1024)},
			{VolumeId: aws.String("bar"), Iops: aws.Int64(1000), Size: aws.Int64(2048)},
		},
	}
	conf.Ec2 = mockedEc2Client{DescribeVolumesPagesRes: mockedOutput, DescribeVolumesPagesError: errors.New("test error")}

	iops, size, err := getEbsVolumeDetails("standard")
	assert.Equal(t, 0, iops)
	assert.Equal(t, 0, size)
	assert.NotNil(t, err)
}

func TestGiBtoTiB(t *testing.T) {
	tib := GiBtoTiB(float64(1024))
	assert.Equal(t, float64(1), tib)
}
