package awslimitchecker_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/sebasrp/awslimitchecker"
	"github.com/sebasrp/awslimitchecker/internal/services"
	"github.com/stretchr/testify/assert"
)

func GetTestUsage(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo) {
	s3checker := services.NewS3Checker(session, quotaClient)
	ret = s3checker.GetUsage()
	return
}

func TestValidateAwsServiceSuccess(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo){
		"foo": GetTestUsage,
	}
	var input = "foo"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Truef(t, actual, "%s should be valid service", input)
}

func TestValidateAwsServiceFailure(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo){
		"foo": GetTestUsage,
	}
	var input = "bar"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Falsef(t, actual, "%s should not be valid service", input)
}

func TestValidateAwsServiceAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo){
		"foo": GetTestUsage,
	}
	var input = "all"
	var actual = awslimitchecker.IsValidAwsService(input)
	assert.Truef(t, actual, "%s should be valid service", input)

}
