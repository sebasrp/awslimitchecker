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
		ServiceCode: &svcName,
		QuotaName:   &quotaName,
		Value:       &quotaValue,
		GlobalQuota: &isGlobal,
	}
}

type mockedScvQuotaClient struct {
	servicequotasiface.ServiceQuotasAPI
	ListAWSDefaultServiceQuotasOutputResp  servicequotas.ListAWSDefaultServiceQuotasOutput
	ListAWSDefaultServiceQuotasOutputError error
	ListServiceQuotasOutputResp            servicequotas.ListServiceQuotasOutput
	ListServiceQuotasOutputError           error
}

func (m mockedScvQuotaClient) ListAWSDefaultServiceQuotasPages(
	input *servicequotas.ListAWSDefaultServiceQuotasInput,
	fn func(*servicequotas.ListAWSDefaultServiceQuotasOutput, bool) bool) error {
	fn(&m.ListAWSDefaultServiceQuotasOutputResp, false)
	return m.ListAWSDefaultServiceQuotasOutputError
}

func (m mockedScvQuotaClient) ListServiceQuotasPages(
	input *servicequotas.ListServiceQuotasInput,
	fn func(*servicequotas.ListServiceQuotasOutput, bool) bool) error {
	fn(&m.ListServiceQuotasOutputResp, false)
	return m.ListServiceQuotasOutputError
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
	mockedDefaultQuotasOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("testServiceName", "testQuotaName", float64(100), false),
		},
	}
	mockedSvcQuotaClient := mockedScvQuotaClient{ListServiceQuotasOutputResp: mockedDefaultQuotasOutput}
	testChecker := NewTestChecker(mockedSvcQuotaClient, supportedQuotas)
	assert.Equal(t, 1, len(testChecker.GetUsage()))
}

func TestGetAllDefaultQuotas(t *testing.T) {
	mockedOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("testServiceName", "testQuotaName", float64(100), false),
		},
	}
	mockedSvcQuotaClient := mockedScvQuotaClient{ListAWSDefaultServiceQuotasOutputResp: mockedOutput}
	testChecker := NewTestChecker(mockedSvcQuotaClient, nil)
	assert.Equal(t, 1, len(testChecker.GetAllDefaultQuotas()))
}

func TestGetAllDefaultQuotasError(t *testing.T) {
	mockedOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("testServiceNam2e", "testQuotaName2", float64(100), false),
		},
	}
	mockedSvcQuotaClient := mockedScvQuotaClient{
		ListAWSDefaultServiceQuotasOutputResp:  mockedOutput,
		ListAWSDefaultServiceQuotasOutputError: errors.New("test error"),
	}
	testChecker := NewTestChecker(mockedSvcQuotaClient, nil)
	assert.Empty(t, testChecker.GetAllDefaultQuotas())
}

func TestServiceCheckerGetRequiredPermissions(t *testing.T) {
	testChecker := NewTestChecker(nil, nil)
	assert.Equal(t, 1, len(testChecker.GetRequiredPermissions()))
}

func TestSvcQuotaToQuotaInfo(t *testing.T) {
	svcQuotaServiceCode := "testService"
	svqQuotaServiceName := "my verbose testService name"
	svcQuotaQuotaName := "quotaName"
	svcQuotaQuotaCode := "quotaCode"
	svqQuotaQuotaValue := float64(10)
	svcQuotaUnit := "myUnit"
	svcQuotaGlobal := true

	svcQuota := servicequotas.ServiceQuota{
		ServiceCode: &svcQuotaServiceCode,
		ServiceName: &svqQuotaServiceName,
		QuotaName:   &svcQuotaQuotaName,
		QuotaCode:   &svcQuotaQuotaCode,
		Value:       &svqQuotaQuotaValue,
		Unit:        &svcQuotaUnit,
		GlobalQuota: &svcQuotaGlobal,
	}

	expected := AWSQuotaInfo{
		Service:    "testService",
		Name:       "quotaName",
		Quotacode:  "quotaCode",
		QuotaValue: float64(10),
		UsageValue: 0.0,
		Unit:       "myUnit",
		Global:     true,
	}
	assert.Equal(t, expected, svcQuotaToQuotaInfo(&svcQuota))
}
