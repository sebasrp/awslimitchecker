package awslimitchecker_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
	"github.com/sebasrp/awslimitchecker"
	"github.com/sebasrp/awslimitchecker/internal/services"
	"github.com/stretchr/testify/assert"
)

type TestChecker struct {
	// serviceCode is the name of the service this checker verifies
	serviceCode string
	// region the checker will run against
	region string
	// aws client used to call kinesis service
	client kinesisiface.KinesisAPI
	// aws client used to call service quotas service
	svcQuotaClient servicequotasiface.ServiceQuotasAPI
	// the default quotas of the service
	defaultQuotas map[string]services.AWSQuotaInfo
	// supportedQuotas contains the service quota name and the func used to retrieve its usage
	supportedQuotas map[string]func(TestChecker) (ret services.AWSQuotaInfo)
}

func NewTestChecker(session *session.Session, svcQuota *servicequotas.ServiceQuotas) services.Svcquota {
	c := &TestChecker{
		serviceCode:    "test",
		region:         *session.Config.Region,
		client:         kinesis.New(session),
		svcQuotaClient: svcQuota,
		defaultQuotas:  map[string]services.AWSQuotaInfo{},
		supportedQuotas: map[string]func(TestChecker) (ret services.AWSQuotaInfo){
			"foo": TestChecker.GetTestUsage,
		},
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
		Name:       "testQuota",
		Region:     "testRegion",
		Quotacode:  "test-quota",
		QuotaValue: 100,
		UsageValue: 50,
		Unit:       "",
		Global:     true,
	}
	return
}

func (c TestChecker) GetAllDefaultQuotas() map[string]services.AWSQuotaInfo {
	c.defaultQuotas["testQuota"] = services.AWSQuotaInfo{
		Service:    "testService",
		Name:       "testQuota",
		Region:     "testRegion",
		Quotacode:  "test-quota",
		QuotaValue: 100,
		UsageValue: 0.0,
		Unit:       "",
		Global:     true,
	}
	return c.defaultQuotas
}

func (c TestChecker) GetRequiredPermissions() []string {
	return []string{"testService:testIAMPolicy"}
}

func TestValidateAwsServiceSuccess(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session *session.Session, quotaClient *servicequotas.ServiceQuotas) services.Svcquota{
		"foo": NewTestChecker,
	}
	var input = "foo"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Truef(t, actual, "%s should be valid service", input)
}

func TestValidateAwsServiceFailure(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session *session.Session, quotaClient *servicequotas.ServiceQuotas) services.Svcquota{
		"foo": NewTestChecker,
	}
	var input = "bar"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Falsef(t, actual, "%s should not be valid service", input)
}

func TestValidateAwsServiceAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session *session.Session, quotaClient *servicequotas.ServiceQuotas) services.Svcquota{
		"foo": NewTestChecker,
	}
	var input = "all"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Truef(t, actual, "%s should be valid service", input)
}
