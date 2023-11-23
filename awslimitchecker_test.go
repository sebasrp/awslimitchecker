package awslimitchecker_test

import (
	"errors"
	"fmt"
	"testing"

	awslimitchecker "github.com/nyambati/aws-service-limits-exporter"
	"github.com/nyambati/aws-service-limits-exporter/internal/services"
	"github.com/stretchr/testify/assert"
)

type TestChecker struct {
	// serviceCode is the name of the service this checker verifies
	serviceCode string
	// region the checker will run against
	region string
	// the applied quotas for the service. For some quotas, only default values are available
	appliedQuotas map[string]services.AWSQuotaInfo
	// the default quotas of the service
	defaultQuotas map[string]services.AWSQuotaInfo
	// supportedQuotas contains the service quota name and the func used to retrieve its usage
	supportedQuotas map[string]func(TestChecker) (ret services.AWSQuotaInfo)
	// Permissions required to get usage
	requiredPermissions []string
}

func NewTestChecker() services.Svcquota {
	c := &TestChecker{
		serviceCode: "testService",
		region:      "testRegion",
		appliedQuotas: map[string]services.AWSQuotaInfo{"testQuota": {
			Service:    "testService",
			QuotaName:  "testQuota",
			Region:     "testRegion",
			Quotacode:  "test-quota",
			QuotaValue: 200,
			UsageValue: 0.0,
			Unit:       "",
			Global:     true,
		}},
		defaultQuotas: map[string]services.AWSQuotaInfo{"testQuota": {
			Service:    "testService",
			QuotaName:  "testQuota",
			Region:     "testRegion",
			Quotacode:  "test-quota",
			QuotaValue: 100,
			UsageValue: 0.0,
			Unit:       "",
			Global:     true,
		}},
		supportedQuotas: map[string]func(TestChecker) (ret services.AWSQuotaInfo){
			"testQuota": TestChecker.GetTestUsage,
		},
		requiredPermissions: []string{"testService:testIAMPolicy"},
	}
	return c
}

func (c TestChecker) GetUsage() (ret []services.AWSQuotaInfo) {
	for _, q := range c.supportedQuotas {
		quotaInfo := q(c)
		ret = append(ret, quotaInfo)
	}
	return
}

func (c TestChecker) GetTestUsage() (ret services.AWSQuotaInfo) {
	ret = c.GetAllAppliedQuotas()["testQuota"]
	ret.UsageValue = float64(50)
	return
}

func (c TestChecker) GetAllAppliedQuotas() map[string]services.AWSQuotaInfo {
	return c.appliedQuotas
}

func (c TestChecker) GetAllDefaultQuotas() map[string]services.AWSQuotaInfo {
	return c.defaultQuotas
}

func (c TestChecker) SetQuotasOverride(quotasOverride []services.AWSQuotaOverride) {
	for _, override := range quotasOverride {
		if c.serviceCode != override.Service {
			fmt.Print("not same serviceCode. returning \n")
			return
		}
		if quota, ok := c.GetAllAppliedQuotas()[override.QuotaName]; ok {
			fmt.Print("Applyng override \n")

			quota.QuotaValue = override.QuotaValue
			c.appliedQuotas[override.QuotaName] = quota
		}
	}
}

func (c TestChecker) GetRequiredPermissions() []string {
	return c.requiredPermissions
}

func TestValidateAwsServiceSuccess(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
	}
	var input = "foo"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Truef(t, actual, "%s should be valid service", input)
}

func TestValidateAwsServiceFailure(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
	}
	var input = "bar"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Falsef(t, actual, "%s should not be valid service", input)
}

func TestValidateAwsServiceAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
	}
	var input = "all"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Truef(t, actual, "%s should be valid service", input)
}

func TestGetIamPolicies(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	assert.Equal(t, 2, len(awslimitchecker.GetIamPolicies()))
}

func TestGetUsageAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(region string) error {
		return nil
	}
	assert.Equal(t, 2, len(awslimitchecker.GetUsage("all", "testRegion", nil)))
}

func TestGetUsageSingle(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(region string) error {
		return nil
	}
	assert.Equal(t, 1, len(awslimitchecker.GetUsage("foo", "testProfile", "testRegion", nil)))
}

func TestGetUsageSingleWrong(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(region string) error {
		return nil
	}
	assert.Equal(t, 0, len(awslimitchecker.GetUsage("boz", "testProfile", "testRegion", nil)))
}
func TestGetUsageErrorInit(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(region string) error {
		return errors.New("test error")
	}
	assert.Equal(t, 0, len(awslimitchecker.GetUsage("all", "testRegion", nil)))
}

func TestGetUsageOverride(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"testService":  NewTestChecker,
		"testService2": NewTestChecker,
	}
	services.InitializeConfig = func(region string) error {
		return nil
	}
	actual := awslimitchecker.GetUsage("testService", "testRegion", []services.AWSQuotaOverride{
		{Service: "testService", QuotaName: "testQuota", QuotaValue: float64(300)}})
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, float64(300), actual[0].QuotaValue)
}

func TestGetUsageOverrideAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"testService":  NewTestChecker,
		"testService2": NewTestChecker,
	}
	services.InitializeConfig = func(region string) error {
		return nil
	}
	actual := awslimitchecker.GetUsage("all", "testRegion", []services.AWSQuotaOverride{
		{Service: "testService", QuotaName: "testQuota", QuotaValue: float64(300)}})
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, float64(300), actual[0].QuotaValue)
	assert.Equal(t, float64(300), actual[1].QuotaValue) // because both services have same name
}
