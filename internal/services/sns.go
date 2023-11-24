package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/sns"
)

type SnsClientInterface interface {
	ListTopicsPages(input *sns.ListTopicsInput, fn func(*sns.ListTopicsOutput, bool) bool) error
	ListSubscriptionsPages(input *sns.ListSubscriptionsInput, fn func(*sns.ListSubscriptionsOutput, bool) bool) error
}

func NewSnsChecker() ServiceQuota {
	serviceCode := "sns"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Topics per Account":                ServiceChecker.getSnsTopicsUsage,
		"Pending Subscriptions per Account": ServiceChecker.getSnsPendingSubsUsage,
	}
	requiredPermissions := []string{"sns:ListTopics", "sns:ListSubscriptions"}

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

func (c ServiceChecker) getSnsPendingSubsUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["Pending Subscriptions per Account"]

	subscriptions := []*sns.Subscription{}
	err := conf.Sns.ListSubscriptionsPages(&sns.ListSubscriptionsInput{}, func(p *sns.ListSubscriptionsOutput, lastPage bool) bool {
		subscriptions = append(subscriptions, p.Subscriptions...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve sns subscriptions, %v", err)
		return
	}

	quotaInfo.UsageValue = float64(len(subscriptions))
	ret = append(ret, quotaInfo)
	return
}
