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
	}
	requiredPermissions := []string{"iam:GetAccountSummary"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

var iamAccountQuota map[string]*int64 = map[string]*int64{}

func (c ServiceChecker) getIamAccountQuotas() (ret map[string]*int64) {
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

func (c ServiceChecker) getIamRolesUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	if len(c.getIamAccountQuotas()) == 0 {
		return
	}

	rolesQuota := AWSQuotaInfo{
		Service: c.ServiceCode,
		Name:    "Roles per Account",
		Global:  true,
	}
	if val, ok := c.getIamAccountQuotas()["RolesQuota"]; ok {
		rolesQuota.QuotaValue = float64(*val)
	}
	if val, ok := iamAccountQuota["Roles"]; ok {
		rolesQuota.UsageValue = float64(*val)
	}
	ret = append(ret, rolesQuota)
	return
}
