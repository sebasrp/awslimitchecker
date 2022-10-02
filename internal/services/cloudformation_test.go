package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedCloudformationClient struct {
	CloudformationClientInterface
	ListStacksPagesResp  cloudformation.ListStacksOutput
	ListStacksPagesError error
}

func (m mockedCloudformationClient) ListStacksPages(input *cloudformation.ListStacksInput, fn func(*cloudformation.ListStacksOutput, bool) bool) error {
	fn(&m.ListStacksPagesResp, false)
	return m.ListStacksPagesError
}

func TestNewCloudformationCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewCloudformationChecker())
}

func TestGetCloudformationStackUsage(t *testing.T) {
	mockedOutput := cloudformation.ListStacksOutput{
		StackSummaries: []*cloudformation.StackSummary{{StackId: aws.String("foo")}},
	}
	conf.Cloudformation = mockedCloudformationClient{ListStacksPagesResp: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("cloudformation", "Stack count", float64(100), false)},
		nil)

	cfChecker := NewCloudformationChecker()
	svcChecker := cfChecker.(*ServiceChecker)
	actual := svcChecker.getCloudformationStackUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "cloudformation", quota.Service)
	assert.Equal(t, float64(100), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.StackSummaries)), quota.UsageValue)
}

func TestGetCloudformationStackUsageError(t *testing.T) {
	mockedOutput := cloudformation.ListStacksOutput{
		StackSummaries: []*cloudformation.StackSummary{{StackId: aws.String("foo")}},
	}
	conf.Cloudformation = mockedCloudformationClient{ListStacksPagesResp: mockedOutput, ListStacksPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("cloudformation", "Stack count", float64(100), false)},
		nil)

	cfChecker := NewCloudformationChecker()
	svcChecker := cfChecker.(*ServiceChecker)
	actual := svcChecker.getCloudformationStackUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}
