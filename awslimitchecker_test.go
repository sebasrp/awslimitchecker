package awslimitchecker_test

import (
	"errors"
	"testing"

	"github.com/sebasrp/awslimitchecker"
	"github.com/sebasrp/awslimitchecker/internal/services"
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
		serviceCode:   "test",
		region:        "testRegion",
		appliedQuotas: map[string]services.AWSQuotaInfo{},
		defaultQuotas: map[string]services.AWSQuotaInfo{},
		supportedQuotas: map[string]func(TestChecker) (ret services.AWSQuotaInfo){
			"foo": TestChecker.GetTestUsage,
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
	ret = services.AWSQuotaInfo{
		Service:    "testService",
		QuotaName:  "testQuota",
		Region:     "testRegion",
		Quotacode:  "test-quota",
		QuotaValue: 100,
		UsageValue: 50,
		Unit:       "",
		Global:     true,
	}
	return
}

func (c TestChecker) GetAllAppliedQuotas() map[string]services.AWSQuotaInfo {
	c.appliedQuotas["testQuota"] = services.AWSQuotaInfo{
		Service:    "testService",
		QuotaName:  "testQuota",
		Region:     "testRegion",
		Quotacode:  "test-quota",
		QuotaValue: 200,
		UsageValue: 0.0,
		Unit:       "",
		Global:     true,
	}
	return c.appliedQuotas
}

func (c TestChecker) GetAllDefaultQuotas() map[string]services.AWSQuotaInfo {
	c.defaultQuotas["testQuota"] = services.AWSQuotaInfo{
		Service:    "testService",
		QuotaName:  "testQuota",
		Region:     "testRegion",
		Quotacode:  "test-quota",
		QuotaValue: 100,
		UsageValue: 0.0,
		Unit:       "",
		Global:     true,
	}
	return c.defaultQuotas
}

func (c TestChecker) SetQuotaOverride(serviceName string, quotaName string, value float64) {
	if c.serviceCode != serviceName {
		return
	}
	if quota, ok := c.GetAllAppliedQuotas()[quotaName]; !ok {
		quota.QuotaValue = value
		c.appliedQuotas[quotaName] = quota
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

func TestGetLimitsAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(awsprofile, region string) (*services.Config, error) {
		return &services.Config{}, nil
	}
	assert.Equal(t, 2, len(awslimitchecker.GetUsage("all", "testProfile", "testRegion")))
}

func TestGetLimitsSingle(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(awsprofile, region string) (*services.Config, error) {
		return &services.Config{}, nil
	}
	assert.Equal(t, 1, len(awslimitchecker.GetUsage("foo", "testProfile", "testRegion")))
}

func TestGetLimitsSingleWrong(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(awsprofile, region string) (*services.Config, error) {
		return &services.Config{}, nil
	}
	assert.Equal(t, 0, len(awslimitchecker.GetUsage("boz", "testProfile", "testRegion")))
}
func TestGetLimitsErrorInit(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func() services.Svcquota{
		"foo": NewTestChecker,
		"bar": NewTestChecker,
	}
	services.InitializeConfig = func(awsprofile, region string) (*services.Config, error) {
		return &services.Config{}, errors.New("test error")
	}
	assert.Equal(t, 0, len(awslimitchecker.GetUsage("all", "testProfile", "testRegion")))
}
