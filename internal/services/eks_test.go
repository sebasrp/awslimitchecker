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

type mockedEksClient struct {
	eksiface.EKSAPI
	ListClustersPagesResp    eks.ListClustersOutput
	ListClustersPagesError   error
	ListNodegroupsPagesResp  eks.ListNodegroupsOutput
	ListNodegroupsPagesError error
}

func (m mockedEksClient) ListClustersPages(
	input *eks.ListClustersInput,
	fn func(*eks.ListClustersOutput, bool) bool) error {
	fn(&m.ListClustersPagesResp, false)
	return m.ListClustersPagesError
}

func (m mockedEksClient) ListNodegroupsPages(
	input *eks.ListNodegroupsInput,
	fn func(*eks.ListNodegroupsOutput, bool) bool) error {
	fn(&m.ListNodegroupsPagesResp, false)
	return m.ListNodegroupsPagesError
}

func TestNewEksCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewEksChecker())
}

func TestGetEKSClusterUsage(t *testing.T) {
	mockedOutput := eks.ListClustersOutput{
		Clusters: []*string{aws.String("foo"), aws.String("bar")},
	}
	conf.Eks = mockedEksClient{ListClustersPagesResp: mockedOutput}

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("eks", "Clusters", float64(10), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	eksChecker := NewEksChecker()
	svcChecker := eksChecker.(*ServiceChecker)
	actual := svcChecker.getEKSClusterUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "eks", quota.Service)
	assert.Equal(t, float64(10), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.Clusters)), quota.UsageValue)
}

func TestGetEKSClusterUsageError(t *testing.T) {
	mockedOutput := eks.ListClustersOutput{
		Clusters: []*string{aws.String("foo"), aws.String("bar")},
	}
	conf.Eks = mockedEksClient{ListClustersPagesResp: mockedOutput, ListClustersPagesError: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("eks", "Clusters", float64(10), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	eksChecker := NewEksChecker()
	svcChecker := eksChecker.(*ServiceChecker)
	actual := svcChecker.getEKSClusterUsage()
	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}

func TestGetEKSNodeGroupsPerClusterUsage(t *testing.T) {
	mockedListClustersOutput := eks.ListClustersOutput{
		Clusters: []*string{aws.String("foo"), aws.String("bar")},
	}
	mockedListNodegroupsOutput := eks.ListNodegroupsOutput{
		Nodegroups: []*string{aws.String("baz"), aws.String("qux")},
	}
	conf.Eks = mockedEksClient{ListClustersPagesResp: mockedListClustersOutput, ListNodegroupsPagesResp: mockedListNodegroupsOutput}

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("eks", "Managed node groups per cluster", float64(10), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	eksChecker := NewEksChecker()
	svcChecker := eksChecker.(*ServiceChecker)
	actual := svcChecker.getEKSNodeGroupsPerClusterUsage()

	assert.Len(t, actual, 2)
	firstQuota := actual[0]
	assert.Equal(t, "eks", firstQuota.Service)
	assert.Equal(t, "AWS::EKS::Cluster::foo", firstQuota.ResourceId)
	assert.Equal(t, float64(10), firstQuota.QuotaValue)
	assert.Equal(t, float64(len(mockedListNodegroupsOutput.Nodegroups)), firstQuota.UsageValue)
}

func TestGetEKSNodeGroupsPerClusterUsageErrorCluster(t *testing.T) {
	mockedListClustersOutput := eks.ListClustersOutput{
		Clusters: []*string{aws.String("foo"), aws.String("bar")},
	}
	mockedListNodegroupsOutput := eks.ListNodegroupsOutput{
		Nodegroups: []*string{aws.String("baz"), aws.String("qux")},
	}
	conf.Eks = mockedEksClient{
		ListClustersPagesResp: mockedListClustersOutput, ListClustersPagesError: errors.New("test error"),
		ListNodegroupsPagesResp: mockedListNodegroupsOutput}

	eksChecker := NewEksChecker()
	svcChecker := eksChecker.(*ServiceChecker)
	actual := svcChecker.getEKSNodeGroupsPerClusterUsage()

	assert.Len(t, actual, 0)
}

func TestGetEKSNodeGroupsPerClusterUsageErrorNodeGroup(t *testing.T) {
	mockedListClustersOutput := eks.ListClustersOutput{
		Clusters: []*string{aws.String("foo"), aws.String("bar")},
	}
	mockedListNodegroupsOutput := eks.ListNodegroupsOutput{
		Nodegroups: []*string{aws.String("baz"), aws.String("qux")},
	}
	conf.Eks = mockedEksClient{
		ListClustersPagesResp:   mockedListClustersOutput,
		ListNodegroupsPagesResp: mockedListNodegroupsOutput, ListNodegroupsPagesError: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("eks", "Managed node groups per cluster", float64(10), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	eksChecker := NewEksChecker()
	svcChecker := eksChecker.(*ServiceChecker)
	actual := svcChecker.getEKSNodeGroupsPerClusterUsage()

	assert.Len(t, actual, 0)
}
