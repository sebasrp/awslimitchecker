package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedElastiCacheClient struct {
	ElastiCacheClientInterface
	DescribeCacheClustersPagesResp  elasticache.DescribeCacheClustersOutput
	DescribeCacheClustersPagesError error
}

func (m mockedElastiCacheClient) DescribeCacheClustersPages(input *elasticache.DescribeCacheClustersInput, fn func(*elasticache.DescribeCacheClustersOutput, bool) bool) error {
	fn(&m.DescribeCacheClustersPagesResp, false)
	return m.DescribeCacheClustersPagesError
}

func TestNewElastiCacheCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewElastiCacheChecker())
}

func TestGetElastiCacheNodesUsage(t *testing.T) {
	mockedOutput := elasticache.DescribeCacheClustersOutput{
		CacheClusters: []*elasticache.CacheCluster{{ARN: aws.String("foo")}},
	}
	conf.ElastiCache = mockedElastiCacheClient{DescribeCacheClustersPagesResp: mockedOutput}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("elasticache", "Nodes per Region", float64(100), false)},
		nil)

	elastiCacheChecker := NewElastiCacheChecker()
	svcChecker := elastiCacheChecker.(*ServiceChecker)
	actual := svcChecker.getElastiCacheNodesUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "elasticache", quota.Service)
	assert.Equal(t, float64(100), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.CacheClusters)), quota.UsageValue)
}

func TestGetElastiCacheNodesUsageError(t *testing.T) {
	mockedOutput := elasticache.DescribeCacheClustersOutput{
		CacheClusters: []*elasticache.CacheCluster{{ARN: aws.String("foo")}},
	}
	conf.ElastiCache = mockedElastiCacheClient{DescribeCacheClustersPagesResp: mockedOutput, DescribeCacheClustersPagesError: errors.New("test error")}

	conf.ServiceQuotas = NewSvcQuotaMockListServiceQuotas(
		[]*servicequotas.ServiceQuota{NewQuota("elasticache", "Nodes per Region", float64(100), false)},
		nil)

	elastiCacheChecker := NewElastiCacheChecker()
	svcChecker := elastiCacheChecker.(*ServiceChecker)
	actual := svcChecker.getElastiCacheNodesUsage()

	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}
