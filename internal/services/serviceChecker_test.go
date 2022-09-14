package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestChecker(svcQuotaMockClient SvcQuotaClientInterface, supportedQuotas map[string]func(ServiceChecker) (ret []AWSQuotaInfo)) Svcquota {
	serviceCode := "testService"
	requiredPermissions := []string{"test:ListTestIAM"}
	conf.ServiceQuotas = svcQuotaMockClient

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func NewQuota(svcName string, quotaName string, quotaValue float64, isGlobal bool) *servicequotas.ServiceQuota {
	return &servicequotas.ServiceQuota{
		ServiceName: &svcName,
		QuotaName:   &quotaName,
		Value:       &quotaValue,
		GlobalQuota: &isGlobal,
	}
}

type mockedListAWSDefaultServiceQuotasPagesMsgs struct {
	servicequotasiface.ServiceQuotasAPI
	Resp  servicequotas.ListAWSDefaultServiceQuotasOutput
	Error error
}

func (m mockedListAWSDefaultServiceQuotasPagesMsgs) ListAWSDefaultServiceQuotasPages(
	input *servicequotas.ListAWSDefaultServiceQuotasInput,
	fn func(*servicequotas.ListAWSDefaultServiceQuotasOutput, bool) bool) error {
	fn(&m.Resp, false)
	return m.Error
}

func TestNewServiceCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewTestChecker(nil, nil))
}

func TestGetUsage(t *testing.T) {
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"testQuotaName": func(c ServiceChecker) (ret []AWSQuotaInfo) {
			quota := c.GetAllDefaultQuotas()["testQuotaName"]
			quota.UsageValue = float64(100)
			ret = append(ret, quota)
			return
		},
	}
	mockedDefaultQuotasOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("testServiceName", "testQuotaName", float64(100), false),
		},
	}
	mockedSvcQuotaClient := mockedListAWSDefaultServiceQuotasPagesMsgs{Resp: mockedDefaultQuotasOutput}
	testChecker := NewTestChecker(mockedSvcQuotaClient, supportedQuotas)
	assert.Equal(t, 1, len(testChecker.GetUsage()))
}

func TestGetAllDefaultQuotas(t *testing.T) {
	mockedOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("testServiceName", "testQuotaName", float64(100), false),
		},
	}
	mockedSvcQuotaClient := mockedListAWSDefaultServiceQuotasPagesMsgs{Resp: mockedOutput}
	testChecker := NewTestChecker(mockedSvcQuotaClient, nil)
	assert.Equal(t, 1, len(testChecker.GetAllDefaultQuotas()))
}

func TestGetAllDefaultQuotasError(t *testing.T) {
	mockedOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("testServiceNam2e", "testQuotaName2", float64(100), false),
		},
	}
	mockedSvcQuotaClient := mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp:  mockedOutput,
		Error: errors.New("test error"),
	}
	testChecker := NewTestChecker(mockedSvcQuotaClient, nil)
	assert.Empty(t, testChecker.GetAllDefaultQuotas())
}

func TestServiceCheckerGetRequiredPermissions(t *testing.T) {
	testChecker := NewTestChecker(nil, nil)
	assert.Equal(t, 1, len(testChecker.GetRequiredPermissions()))
}
