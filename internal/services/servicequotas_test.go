package services

import (
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"
)

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

func NewQuota(svcName string, quotaName string, quotaValue float64, isGlobal bool) *servicequotas.ServiceQuota {
	return &servicequotas.ServiceQuota{
		ServiceCode: &svcName,
		QuotaName:   &quotaName,
		Value:       &quotaValue,
		GlobalQuota: &isGlobal,
	}
}

func NewSvcQuotaMockListServiceQuotas(quotas []*servicequotas.ServiceQuota, err error) (ret mockedScvQuotaClient) {
	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: quotas,
	}
	ret = mockedScvQuotaClient{
		ListServiceQuotasOutputResp:  mockedSvcQuotaOutput,
		ListServiceQuotasOutputError: err,
	}
	return
}
func NewSvcQuotaMockListAWSDefaultServiceQuotas(quotas []*servicequotas.ServiceQuota, err error) (ret mockedScvQuotaClient) {
	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: quotas,
	}
	ret = mockedScvQuotaClient{
		ListAWSDefaultServiceQuotasOutputResp:  mockedSvcQuotaOutput,
		ListAWSDefaultServiceQuotasOutputError: err,
	}
	return
}
