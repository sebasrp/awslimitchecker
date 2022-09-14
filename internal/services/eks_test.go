package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedListClustersPagesMsgs struct {
	eksiface.EKSAPI
	Resp  eks.ListClustersOutput
	Error error
}

func (m mockedListClustersPagesMsgs) ListClustersPages(
	input *eks.ListClustersInput,
	fn func(*eks.ListClustersOutput, bool) bool) error {
	fn(&m.Resp, false)
	return m.Error
}

func TestNewEksCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewEksChecker())
}

func TestGetEKSClusterUsage(t *testing.T) {
	mockedOutput := eks.ListClustersOutput{
		Clusters: []*string{aws.String("foo"), aws.String("bar")},
	}
	conf.Eks = mockedListClustersPagesMsgs{Resp: mockedOutput}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("eks", "Clusters", float64(10), false),
		},
	}
	conf.ServiceQuotas = mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	eksChecker := NewEksChecker()
	svcChecker := eksChecker.(*ServiceChecker)
	actual := svcChecker.getEKSClusterUsage()

	assert.Equal(t, "eks", actual.Service)
	assert.Equal(t, float64(10), actual.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.Clusters)), actual.UsageValue)
}

func TestGetEKSClusterUsageError(t *testing.T) {
	mockedOutput := eks.ListClustersOutput{
		Clusters: []*string{aws.String("foo"), aws.String("bar")},
	}
	conf.Eks = mockedListClustersPagesMsgs{Resp: mockedOutput, Error: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListAWSDefaultServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("eks", "Clusters", float64(10), false),
		},
	}
	conf.ServiceQuotas = mockedListAWSDefaultServiceQuotasPagesMsgs{
		Resp: mockedSvcQuotaOutput,
	}

	eksChecker := NewEksChecker()
	svcChecker := eksChecker.(*ServiceChecker)
	actual := svcChecker.getEKSClusterUsage()
	expected := AWSQuotaInfo{Service: "eks", Name: "Clusters", Region: "", Quotacode: "", QuotaValue: 10, UsageValue: 0, Unit: "", Global: false}
	assert.Equal(t, expected, actual)
}
