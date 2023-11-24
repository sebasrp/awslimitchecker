package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedSnsClient struct {
	SnsClientInterface
	ListTopicPagesResp         sns.ListTopicsOutput
	ListTopicPagesError        error
	ListSubscriptionsPagesRest sns.ListSubscriptionsOutput
	ListSubscriptionsPagesErr  error
}

func (m mockedSnsClient) ListTopicsPages(input *sns.ListTopicsInput, fn func(*sns.ListTopicsOutput, bool) bool) error {
	fn(&m.ListTopicPagesResp, false)
	return m.ListTopicPagesError
}

func (m mockedSnsClient) ListSubscriptionsPages(
	input *sns.ListSubscriptionsInput, fn func(*sns.ListSubscriptionsOutput, bool) bool) error {
	fn(&m.ListSubscriptionsPagesRest, false)
	return m.ListSubscriptionsPagesErr
}

func TestNewSnsCheckerImpl(t *testing.T) {
	require.Implements(t, (*ServiceQuota)(nil), NewSnsChecker())
}

func TestGetSnsTopicsUsage(t *testing.T) {
	mockedOutput := sns.ListTopicsOutput{
		Topics: []*sns.Topic{{TopicArn: aws.String("foo")}},
	}
	conf.Sns = mockedSnsClient{ListTopicPagesResp: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("sns", "Topics per Account", float64(100), false)},
		nil)

	snseChecker := NewSnsChecker()
	svcChecker := snseChecker.(*ServiceChecker)
	actual := svcChecker.getSnsTopicsUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "sns", quota.Service)
	assert.Equal(t, float64(100), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.Topics)), quota.UsageValue)
}

func TestGetSnsTopicsUsageError(t *testing.T) {
	mockedOutput := sns.ListTopicsOutput{
		Topics: []*sns.Topic{{TopicArn: aws.String("foo")}},
	}
	conf.Sns = mockedSnsClient{ListTopicPagesResp: mockedOutput, ListTopicPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("sns", "Topics per Account", float64(100), false)},
		nil)

	snseChecker := NewSnsChecker()
	svcChecker := snseChecker.(*ServiceChecker)
	actual := svcChecker.getSnsTopicsUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}

func TestGetSnsPendingSubsUsage(t *testing.T) {
	mockedOutput := sns.ListSubscriptionsOutput{
		Subscriptions: []*sns.Subscription{{SubscriptionArn: aws.String("foo")}},
	}
	conf.Sns = mockedSnsClient{ListSubscriptionsPagesRest: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("sns", "Pending Subscriptions per Account", float64(10), false)},
		nil)

	snseChecker := NewSnsChecker()
	svcChecker := snseChecker.(*ServiceChecker)
	actual := svcChecker.getSnsPendingSubsUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "sns", quota.Service)
	assert.Equal(t, float64(10), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.Subscriptions)), quota.UsageValue)
}

func TestGetSnsPendingSubsUsageError(t *testing.T) {
	mockedOutput := sns.ListSubscriptionsOutput{
		Subscriptions: []*sns.Subscription{{SubscriptionArn: aws.String("foo")}},
	}
	conf.Sns = mockedSnsClient{ListSubscriptionsPagesRest: mockedOutput, ListSubscriptionsPagesErr: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("sns", "Pending Subscriptions per Account", float64(10), false)},
		nil)

	snseChecker := NewSnsChecker()
	svcChecker := snseChecker.(*ServiceChecker)
	actual := svcChecker.getSnsPendingSubsUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}
