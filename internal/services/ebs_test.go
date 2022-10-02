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
