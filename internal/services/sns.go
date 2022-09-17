package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsClientInterface interface {
	ListTopicsPages(input *sns.ListTopicsInput, fn func(*sns.ListTopicsOutput, bool) bool) error
}

func NewSnsChecker() Svcquota {
	serviceCode := "sns"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Topics per Account": ServiceChecker.getSnsTopicsUsage,
	}
	requiredPermissions := []string{"sns:ListTopics"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getSnsTopicsUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["Topics per Account"]

	topics := []*sns.Topic{}
	err := conf.Sns.ListTopicsPages(&sns.ListTopicsInput{}, func(p *sns.ListTopicsOutput, lastPage bool) bool {
		topics = append(topics, p.Topics...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve sns topics, %v", err)
		return
	}

	quotaInfo.UsageValue = float64(len(topics))
	ret = append(ret, quotaInfo)
	return
}
