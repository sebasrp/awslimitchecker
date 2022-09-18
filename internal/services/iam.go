package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/iam"
)

type IamClientInterface interface {
	GetAccountSummary(input *iam.GetAccountSummaryInput) (*iam.GetAccountSummaryOutput, error)
}

func NewIamChecker() Svcquota {
	serviceCode := "iam"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Roles per Account": ServiceChecker.getIamRolesUsage,
		"Users per Account": ServiceChecker.getIamUsersUsage,
	}
	requiredPermissions := []string{"iam:GetAccountSummary"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

var iamAccountQuota map[string]*int64 = map[string]*int64{}

func getIamAccountQuotas() (ret map[string]*int64, err error) {
	ret = iamAccountQuota
	if len(iamAccountQuota) != 0 {
		return
	}

	result, err := conf.Iam.GetAccountSummary(nil)
	if err != nil {
		fmt.Printf("Unable to retrieve iam account summary, %v", err)
		return
	}

	for quotaName, value := range result.SummaryMap {
		iamAccountQuota[quotaName] = value
	}
	return
}

func IamSummaryToAWSQuotaInfo(summaryName string, quotaName string) (ret AWSQuotaInfo, err error) {
	ret = AWSQuotaInfo{}
	quotas, err := getIamAccountQuotas()
	if len(quotas) == 0 || err != nil {
		return ret, err
	}

	ret = AWSQuotaInfo{
		Service: "iam",
		Name:    quotaName,
		Global:  true,
	}
	if val, ok := quotas[summaryName+"Quota"]; ok {
		ret.QuotaValue = float64(*val)
	}
	if val, ok := quotas[summaryName]; ok {
		ret.UsageValue = float64(*val)
	}
	return
}

func (c ServiceChecker) getIamRolesUsage() (ret []AWSQuotaInfo) {
	if quotaInfo, err := IamSummaryToAWSQuotaInfo("Roles", "Roles per Account"); err != nil {
		return []AWSQuotaInfo{}
	} else {
		return []AWSQuotaInfo{quotaInfo}
	}
}

func (c ServiceChecker) getIamUsersUsage() (ret []AWSQuotaInfo) {
	if quotaInfo, err := IamSummaryToAWSQuotaInfo("Users", "Users per Account"); err != nil {
		return []AWSQuotaInfo{}
	} else {
		return []AWSQuotaInfo{quotaInfo}
	}
}
