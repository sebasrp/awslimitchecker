package services

import (
	"errors"
	"testing"

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
