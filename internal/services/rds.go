package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

type RdsClientInterface interface {
	DescribeAccountAttributes(input *rds.DescribeAccountAttributesInput) (*rds.DescribeAccountAttributesOutput, error)
}

func NewRdsChecker() ServiceQuota {
	serviceCode := "rds"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"DB instances":          ServiceChecker.getRdsInstancesCountUsage,
		"DB clusters":           ServiceChecker.getRdsClusterCountUsage,
		"Reserved DB instances": ServiceChecker.getRdsReservedDbCountUsage,
	}
	requiredPermissions := []string{"rds:DescribeAccountAttributes"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

var rdsAccountQuota map[string]*rds.AccountQuota = map[string]*rds.AccountQuota{}

func (c ServiceChecker) getRdsAccountQuotas() (ret map[string]*rds.AccountQuota) {
	ret = rdsAccountQuota
	if len(rdsAccountQuota) != 0 {
		return
	}

	result, err := conf.Rds.DescribeAccountAttributes(nil)
	if err != nil {
		fmt.Printf("Unable to retrieve account attributes, %v", err)
		return
	}

	for _, q := range result.AccountQuotas {
		rdsAccountQuota[aws.StringValue(q.AccountQuotaName)] = q
	}
	return
}

func (c ServiceChecker) getRdsInstancesCountUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["DB instances"]
	if val, ok := c.getRdsAccountQuotas()["DBInstances"]; ok {
		quotaInfo.UsageValue = float64(*val.Used)
	}
	ret = append(ret, quotaInfo)
	return
}

func (c ServiceChecker) getRdsClusterCountUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["DB clusters"]
	if val, ok := c.getRdsAccountQuotas()["DBClusters"]; ok {
		quotaInfo.UsageValue = float64(*val.Used)
	}
	ret = append(ret, quotaInfo)
	return
}

func (c ServiceChecker) getRdsReservedDbCountUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["Reserved DB instances"]
	if val, ok := c.getRdsAccountQuotas()["ReservedDBInstances"]; ok {
		quotaInfo.UsageValue = float64(*val.Used)
	}
	ret = append(ret, quotaInfo)
	return
}
