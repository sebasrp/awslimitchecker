package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/acm"
)

type AcmClientInterface interface {
	ListCertificatesPages(input *acm.ListCertificatesInput, fn func(*acm.ListCertificatesOutput, bool) bool) error
}

func NewAcmChecker() ServiceQuota {
	serviceCode := "acm"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"ACM certificates": ServiceChecker.getAcmCertificatesUsage,
	}
	requiredPermissions := []string{"acm:ListCertificates"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getAcmCertificatesUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	quotaInfo := c.GetAllAppliedQuotas()["ACM certificates"]

	certificates := []*acm.CertificateSummary{}
	err := conf.Acm.ListCertificatesPages(&acm.ListCertificatesInput{}, func(p *acm.ListCertificatesOutput, lastPage bool) bool {
		certificates = append(certificates, p.CertificateSummaryList...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve acm certificates, %v", err)
		return
	}

	quotaInfo.UsageValue = float64(len(certificates))
	ret = append(ret, quotaInfo)
	return
}
