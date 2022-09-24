package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewTestChecker(supportedQuotas map[string]func(ServiceChecker) (ret []AWSQuotaInfo)) Svcquota {
	serviceCode := "testService"
	requiredPermissions := []string{"test:ListTestIAM"}
	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func TestNewServiceCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewTestChecker(nil))
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
	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testServiceName", "testQuotaName", float64(100), false)},
		nil)
	testChecker := NewTestChecker(supportedQuotas)
	assert.Equal(t, 1, len(testChecker.GetUsage()))
}

func TestGetAllAppliedQuotas(t *testing.T) {
	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testServiceName", "testQuotaName", float64(100), false)},
		nil)
	testChecker := NewTestChecker(nil)
	assert.Equal(t, 1, len(testChecker.GetAllAppliedQuotas()))
}

func TestGetAllAppliedQuotasFallback(t *testing.T) {
	mockedListServiceQuotasOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{NewQuota("servicename1", "testQuotaName1", float64(100), false)},
	}
	mockedListAWSDefaultServiceQuotasOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{NewQuota("servicename2", "testQuotaName2", float64(100), false)},
	}

	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp:           mockedListServiceQuotasOutput,
		ListAWSDefaultServiceQuotasOutputResp: mockedListAWSDefaultServiceQuotasOutput,
	}

	testChecker := NewTestChecker(nil)
	appliedQuotas := testChecker.GetAllAppliedQuotas()
	assert.Equal(t, 2, len(appliedQuotas))
	assert.Contains(t, appliedQuotas, "testQuotaName1")
	assert.Contains(t, appliedQuotas, "testQuotaName2")
}

func TestGetAllAppliedQuotasError(t *testing.T) {
	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testServiceNam2e", "testQuotaName2", float64(100), false)},
		errors.New("test error"))
	testChecker := NewTestChecker(nil)
	assert.Empty(t, testChecker.GetAllAppliedQuotas())
}

func TestGetAllDefaultQuotas(t *testing.T) {
	conf.ServiceQuotas = NewSvcQuotaMockListAWSDefaultServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testServiceName", "testQuotaName", float64(100), false)},
		nil)
	testChecker := NewTestChecker(nil)
	assert.Equal(t, 1, len(testChecker.GetAllDefaultQuotas()))
}

func TestGetAllDefaultQuotasError(t *testing.T) {
	conf.ServiceQuotas = NewSvcQuotaMockListAWSDefaultServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testServiceNam2e", "testQuotaName2", float64(100), false)},
		errors.New("test error"))
	testChecker := NewTestChecker(nil)
	assert.Empty(t, testChecker.GetAllDefaultQuotas())
}

func TestServiceCheckerGetRequiredPermissions(t *testing.T) {
	testChecker := NewTestChecker(nil)
	assert.Equal(t, 1, len(testChecker.GetRequiredPermissions()))
}

func TestSvcQuotaToQuotaInfo(t *testing.T) {
	svcQuota := servicequotas.ServiceQuota{
		ServiceCode: aws.String("testService"),
		ServiceName: aws.String("my verbose testService name"),
		QuotaName:   aws.String("quotaName"),
		QuotaCode:   aws.String("quotaCode"),
		Value:       aws.Float64(float64(10)),
		Unit:        aws.String("myUnit"),
		GlobalQuota: aws.Bool(true),
	}

	expected := AWSQuotaInfo{
		Service:    "testService",
		QuotaName:  "quotaName",
		Quotacode:  "quotaCode",
		QuotaValue: float64(10),
		UsageValue: 0.0,
		Unit:       "myUnit",
		Global:     true,
	}
	assert.Equal(t, expected, svcQuotaToQuotaInfo(&svcQuota))
}

func TestSetQuotaOverride(t *testing.T) {
	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testService", "testQuotaName", float64(100), false)},
		nil)
	testChecker := NewTestChecker(nil)
	testChecker.SetQuotasOverride([]AWSQuotaOverride{{Service: "testService", QuotaName: "testQuotaName", QuotaValue: float64(200)}})
	appliedQuotas := testChecker.GetAllAppliedQuotas()
	assert.Equal(t, float64(200), appliedQuotas["testQuotaName"].QuotaValue)
}

func TestSetQuotaOverrideWrongService(t *testing.T) {
	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testService", "testQuotaName", float64(100), false)},
		nil)
	testChecker := NewTestChecker(nil)
	testChecker.SetQuotasOverride([]AWSQuotaOverride{{Service: "testServiceWrong", QuotaName: "testQuotaName", QuotaValue: float64(200)}})
	appliedQuotas := testChecker.GetAllAppliedQuotas()
	assert.Equal(t, float64(100), appliedQuotas["testQuotaName"].QuotaValue)
}

func TestSetQuotaOverrideWrongQuotaName(t *testing.T) {
	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("testService", "testQuotaName", float64(100), false)},
		nil)
	testChecker := NewTestChecker(nil)
	testChecker.SetQuotasOverride([]AWSQuotaOverride{{Service: "testService", QuotaName: "testQuotaNameWrong", QuotaValue: float64(200)}})
	appliedQuotas := testChecker.GetAllAppliedQuotas()
	assert.Equal(t, float64(100), appliedQuotas["testQuotaName"].QuotaValue)
}
