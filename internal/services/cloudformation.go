package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type CloudformationClientInterface interface {
	ListStacksPages(input *cloudformation.ListStacksInput, fn func(*cloudformation.ListStacksOutput, bool) bool) error
}

func NewCloudformationChecker() Svcquota {
	serviceCode := "cloudformation"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Stack count": ServiceChecker.getCloudformationStackUsage,
	}
	requiredPermissions := []string{"cloudformation:ListStacks"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getCloudformationStackUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["Stack count"]

	stacks := []*cloudformation.StackSummary{}

	// we only count active stacks - see https://docs.aws.amazon.com/sdk-for-go/api/service/cloudformation/#pkg-constants
	validStackStatuses := []*string{
		aws.String("CREATE_IN_PROGRESS"), aws.String("CREATE_COMPLETE"), aws.String("ROLLBACK_IN_PROGRESS"), aws.String("ROLLBACK_FAILED"),
		aws.String("ROLLBACK_COMPLETE"), aws.String("DELETE_IN_PROGRESS"), aws.String("DELETE_FAILED"), aws.String("UPDATE_IN_PROGRESS"),
		aws.String("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"), aws.String("UPDATE_COMPLETE"), aws.String("UPDATE_FAILED"),
		aws.String("UPDATE_ROLLBACK_IN_PROGRESS"), aws.String("UPDATE_ROLLBACK_FAILED"),
		aws.String("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"), aws.String("UPDATE_ROLLBACK_COMPLETE"),
		aws.String("REVIEW_IN_PROGRESS"), aws.String("IMPORT_IN_PROGRESS"), aws.String("IMPORT_COMPLETE"),
		aws.String("IMPORT_ROLLBACK_IN_PROGRESS"), aws.String("IMPORT_ROLLBACK_FAILED"), aws.String("IMPORT_ROLLBACK_COMPLETE")}

	err := conf.Cloudformation.ListStacksPages(&cloudformation.ListStacksInput{StackStatusFilter: validStackStatuses}, func(p *cloudformation.ListStacksOutput, lastPage bool) bool {
		stacks = append(stacks, p.StackSummaries...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve cloudformation stacks, %v", err)
		return
	}

	quotaInfo.UsageValue = float64(len(stacks))

	ret = append(ret, quotaInfo)
	return
}
