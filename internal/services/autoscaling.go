package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

type AutoscalingClientInterface interface {
	DescribeAccountLimits(input *autoscaling.DescribeAccountLimitsInput) (*autoscaling.DescribeAccountLimitsOutput, error)
}

func NewAutoscalingChecker() Svcquota {
	serviceCode := "autoscaling"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Auto Scaling groups per region":   ServiceChecker.getAutoscalingGroupsUsage,
		"Launch configurations per region": ServiceChecker.getAutoscalingLaunchConfigsUsage,
	}
	requiredPermissions := []string{"autoscaling:DescribeAccountLimits"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getAutoscalingGroupsUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	result, err := conf.Autoscaling.DescribeAccountLimits(nil)
	quotaInfo := AWSQuotaInfo{
		Service:   c.ServiceCode,
		QuotaName: "Auto Scaling groups per region",
		Region:    c.Region,
		Quotacode: "",
		Unit:      "",
		Global:    true,
	}
	if err != nil {
		fmt.Printf("Unable to retrieve Autoscaling limits, %v", err)
		return
	}

	quotaInfo.QuotaValue = float64(*result.MaxNumberOfAutoScalingGroups)
	quotaInfo.UsageValue = float64(*result.NumberOfAutoScalingGroups)

	ret = append(ret, quotaInfo)
	return
}

func (c ServiceChecker) getAutoscalingLaunchConfigsUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	result, err := conf.Autoscaling.DescribeAccountLimits(nil)
	quotaInfo := AWSQuotaInfo{
		Service:   c.ServiceCode,
		QuotaName: "Launch configurations per region",
		Region:    c.Region,
		Quotacode: "",
		Unit:      "",
		Global:    true,
	}
	if err != nil {
		fmt.Printf("Unable to retrieve Autoscaling limits, %v", err)
		return
	}

	quotaInfo.QuotaValue = float64(*result.MaxNumberOfLaunchConfigurations)
	quotaInfo.UsageValue = float64(*result.NumberOfLaunchConfigurations)

	ret = append(ret, quotaInfo)
	return
}
