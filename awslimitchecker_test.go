package awslimitchecker_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/sebasrp/awslimitchecker"
	"github.com/sebasrp/awslimitchecker/internal/services"
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
	var expected = true
	var actual = awslimitchecker.IsValidAwsService(input)

	if actual != expected {
		t.Fatalf(`IsValidAwsService("%s") = %t, expected %t, nil`, input, actual, expected)
	}
}

func TestValidateAwsServiceFailure(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo){
		"foo": GetTestUsage,
	}
	var input = "bar"
	var expected = false
	var actual = awslimitchecker.IsValidAwsService(input)

	if actual != expected {
		t.Fatalf(`IsValidAwsService("%s") = %t, expected %t, nil`, input, actual, expected)
	}
}

func TestValidateAwsServiceAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]func(session session.Session, quotaClient *servicequotas.ServiceQuotas) (ret []services.AWSQuotaInfo){
		"foo": GetTestUsage,
	}
	var input = "all"
	var expected = true
	var actual = awslimitchecker.IsValidAwsService(input)

	if actual != expected {
		t.Fatalf(`IsValidAwsService("%s") = %t, expected %t, nil`, input, actual, expected)
	}
}
