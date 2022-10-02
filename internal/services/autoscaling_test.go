package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedAutoscalingDescribeAccountLimitsMsg struct {
	AutoscalingClientInterface
	Resp  autoscaling.DescribeAccountLimitsOutput
	Error error
}

func (m mockedAutoscalingDescribeAccountLimitsMsg) DescribeAccountLimits(input *autoscaling.DescribeAccountLimitsInput) (*autoscaling.DescribeAccountLimitsOutput, error) {
	return &m.Resp, m.Error
}

func TestNewAutoscalingCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewAutoscalingChecker())
}

func TestGetAutoscalingGroupsUsage(t *testing.T) {
	mockedOutput := autoscaling.DescribeAccountLimitsOutput{
		MaxNumberOfAutoScalingGroups: aws.Int64(10),
		NumberOfAutoScalingGroups:    aws.Int64(1),
	}
	conf.Autoscaling = mockedAutoscalingDescribeAccountLimitsMsg{Resp: mockedOutput, Error: nil}

	autoscalingChecker := NewAutoscalingChecker()
	svcChecker := autoscalingChecker.(*ServiceChecker)
	actual := svcChecker.getAutoscalingGroupsUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "autoscaling", quota.Service)
	assert.Equal(t, float64(10), quota.QuotaValue)
	assert.Equal(t, float64(1), quota.UsageValue)
}

func TestGetAutoscalingGroupsUsageError(t *testing.T) {
	mockedOutput := autoscaling.DescribeAccountLimitsOutput{
		MaxNumberOfAutoScalingGroups: aws.Int64(10),
		NumberOfAutoScalingGroups:    aws.Int64(1),
	}
	conf.Autoscaling = mockedAutoscalingDescribeAccountLimitsMsg{Resp: mockedOutput, Error: errors.New("test error")}

	autoscalingChecker := NewAutoscalingChecker()
	svcChecker := autoscalingChecker.(*ServiceChecker)
	actual := svcChecker.getAutoscalingGroupsUsage()
	expected := []AWSQuotaInfo{}

	assert.Equal(t, expected, actual)
}

func TestGetAutoscalingLaunchConfigsUsage(t *testing.T) {
	mockedOutput := autoscaling.DescribeAccountLimitsOutput{
		MaxNumberOfLaunchConfigurations: aws.Int64(10),
		NumberOfLaunchConfigurations:    aws.Int64(1),
	}
	conf.Autoscaling = mockedAutoscalingDescribeAccountLimitsMsg{Resp: mockedOutput, Error: nil}

	autoscalingChecker := NewAutoscalingChecker()
	svcChecker := autoscalingChecker.(*ServiceChecker)
	actual := svcChecker.getAutoscalingLaunchConfigsUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "autoscaling", quota.Service)
	assert.Equal(t, float64(10), quota.QuotaValue)
	assert.Equal(t, float64(1), quota.UsageValue)
}

func TestGetAutoscalingLaunchConfigsUsageError(t *testing.T) {
	mockedOutput := autoscaling.DescribeAccountLimitsOutput{
		MaxNumberOfLaunchConfigurations: aws.Int64(10),
		NumberOfLaunchConfigurations:    aws.Int64(1),
	}
	conf.Autoscaling = mockedAutoscalingDescribeAccountLimitsMsg{Resp: mockedOutput, Error: errors.New("test error")}

	autoscalingChecker := NewAutoscalingChecker()
	svcChecker := autoscalingChecker.(*ServiceChecker)
	actual := svcChecker.getAutoscalingLaunchConfigsUsage()
	expected := []AWSQuotaInfo{}

	assert.Equal(t, expected, actual)
}
