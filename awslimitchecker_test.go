package awslimitchecker_test

import (
	"testing"

	"github.com/sebasrp/awslimitchecker"
)

func TestValidateAwsServiceSuccess(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]bool{
		"foo": true,
	}
	var input = "foo"
	var expected = true
	var actual = awslimitchecker.IsValidAwsService(input)

	if actual != expected {
		t.Fatalf(`IsValidAwsService("%s") = %t, expected %t, nil`, input, actual, expected)
	}
}

func TestValidateAwsServiceFailure(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]bool{
		"foo": true,
	}
	var input = "bar"
	var expected = false
	var actual = awslimitchecker.IsValidAwsService(input)

	if actual != expected {
		t.Fatalf(`IsValidAwsService("%s") = %t, expected %t, nil`, input, actual, expected)
	}
}

func TestValidateAwsServiceAll(t *testing.T) {
	awslimitchecker.SupportedAwsServices = map[string]bool{
		"foo": true,
	}
	var input = "all"
	var expected = true
	var actual = awslimitchecker.IsValidAwsService(input)

	if actual != expected {
		t.Fatalf(`IsValidAwsService("%s") = %t, expected %t, nil`, input, actual, expected)
	}
}
