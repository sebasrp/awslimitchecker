package services

import "github.com/aws/aws-sdk-go/service/servicequotas"

type SvcQuotaClientInterface interface {
	ListAWSDefaultServiceQuotasPages(*servicequotas.ListAWSDefaultServiceQuotasInput, func(*servicequotas.ListAWSDefaultServiceQuotasOutput, bool) bool) error
}
