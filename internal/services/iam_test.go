package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedIamClient struct {
	IamClientInterface
	GetAccountSummaryResp  iam.GetAccountSummaryOutput
	GetAccountSummaryError error
}

func (m mockedIamClient) GetAccountSummary(input *iam.GetAccountSummaryInput) (*iam.GetAccountSummaryOutput, error) {
	return &m.GetAccountSummaryResp, m.GetAccountSummaryError
}

func TestNewIamCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewIamChecker())
}

func TestGetIamAccountQuotas(t *testing.T) {
	mockedGetAccountSummaryOutput := iam.GetAccountSummaryOutput{
		SummaryMap: map[string]*int64{
			"foo": aws.Int64(100),
		},
	}
	conf.Iam = mockedIamClient{GetAccountSummaryResp: mockedGetAccountSummaryOutput}

	iamChecker := NewIamChecker()
	svcChecker := iamChecker.(*ServiceChecker)
	actual := svcChecker.getIamAccountQuotas()
	assert.Len(t, actual, 1)
	fooQuota := actual["foo"]
	assert.NotNil(t, fooQuota)
	assert.Equal(t, aws.Int64(100), fooQuota)
	t.Cleanup(func() { iamAccountQuota = map[string]*int64{} })
}

func TestGetIamAccountQuotasError(t *testing.T) {
	mockedGetAccountSummaryOutput := iam.GetAccountSummaryOutput{
		SummaryMap: map[string]*int64{
			"foo": aws.Int64(100),
		},
	}
	conf.Iam = mockedIamClient{GetAccountSummaryResp: mockedGetAccountSummaryOutput, GetAccountSummaryError: errors.New("test error")}

	iamChecker := NewIamChecker()
	svcChecker := iamChecker.(*ServiceChecker)
	actual := svcChecker.getIamAccountQuotas()
	assert.Len(t, actual, 0)
	t.Cleanup(func() { iamAccountQuota = map[string]*int64{} })
}

func TestGetIamAccountQuotasExists(t *testing.T) {
	iamAccountQuota = map[string]*int64{
		"foo": aws.Int64(100),
		"bar": aws.Int64(200),
	}
	iamChecker := NewIamChecker()
	svcChecker := iamChecker.(*ServiceChecker)
	actual := svcChecker.getIamAccountQuotas()
	assert.Len(t, actual, 2)
	t.Cleanup(func() { iamAccountQuota = map[string]*int64{} })
}

func TestGetIamRolesUsage(t *testing.T) {
	mockedGetAccountSummaryOutput := iam.GetAccountSummaryOutput{
		SummaryMap: map[string]*int64{
			"Roles":      aws.Int64(100),
			"RolesQuota": aws.Int64(1000),
		},
	}
	conf.Iam = mockedIamClient{GetAccountSummaryResp: mockedGetAccountSummaryOutput}

	iamChecker := NewIamChecker()
	svcChecker := iamChecker.(*ServiceChecker)
	actual := svcChecker.getIamRolesUsage()
	assert.Len(t, actual, 1)
	usage := actual[0]
	assert.Equal(t, "iam", usage.Service)
	assert.Equal(t, "Roles per Account", usage.Name)
	assert.Equal(t, float64(1000), usage.QuotaValue)
	assert.Equal(t, float64(100), usage.UsageValue)
	t.Cleanup(func() { iamAccountQuota = map[string]*int64{} })
}

func TestGetIamRolesUsageError(t *testing.T) {
	mockedGetAccountSummaryOutput := iam.GetAccountSummaryOutput{
		SummaryMap: map[string]*int64{
			"Roles":      aws.Int64(100),
			"RolesQuota": aws.Int64(1000),
		},
	}
	conf.Iam = mockedIamClient{GetAccountSummaryResp: mockedGetAccountSummaryOutput, GetAccountSummaryError: errors.New("test error")}

	iamChecker := NewIamChecker()
	svcChecker := iamChecker.(*ServiceChecker)
	actual := svcChecker.getIamRolesUsage()
	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
	t.Cleanup(func() { iamAccountQuota = map[string]*int64{} })
}