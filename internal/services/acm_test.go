package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedAcmClient struct {
	AcmClientInterface
	ListCertificatesPagesResp  acm.ListCertificatesOutput
	ListCertificatesPagesError error
}

func (m mockedAcmClient) ListCertificatesPages(input *acm.ListCertificatesInput, fn func(*acm.ListCertificatesOutput, bool) bool) error {
	fn(&m.ListCertificatesPagesResp, false)
	return m.ListCertificatesPagesError
}

func TestNewAcmCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewAcmChecker())
}

func TestGetAcmCertificatesUsage(t *testing.T) {
	mockedOutput := acm.ListCertificatesOutput{
		CertificateSummaryList: []*acm.CertificateSummary{{CertificateArn: aws.String("foo")}},
	}
	conf.Acm = mockedAcmClient{ListCertificatesPagesResp: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("acm", "ACM certificates", float64(100), false)},
		nil)

	acmChecker := NewAcmChecker()
	svcChecker := acmChecker.(*ServiceChecker)
	actual := svcChecker.getAcmCertificatesUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "acm", quota.Service)
	assert.Equal(t, float64(100), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.CertificateSummaryList)), quota.UsageValue)
}

func TestGetAcmCertificatesUsageError(t *testing.T) {
	mockedOutput := acm.ListCertificatesOutput{
		CertificateSummaryList: []*acm.CertificateSummary{{CertificateArn: aws.String("foo")}},
	}
	conf.Acm = mockedAcmClient{ListCertificatesPagesResp: mockedOutput, ListCertificatesPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("acm", "ACM certificates", float64(100), false)},
		nil)

	acmChecker := NewAcmChecker()
	svcChecker := acmChecker.(*ServiceChecker)
	actual := svcChecker.getAcmCertificatesUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}
