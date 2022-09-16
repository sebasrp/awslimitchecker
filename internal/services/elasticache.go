package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/elasticache"
)

type ElastiCacheClientInterface interface {
	DescribeCacheClustersPages(input *elasticache.DescribeCacheClustersInput, fn func(*elasticache.DescribeCacheClustersOutput, bool) bool) error
}

func NewElastiCacheChecker() Svcquota {
	serviceCode := "elasticache"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Nodes per Region": ServiceChecker.getElastiCacheNodesUsage,
	}
	requiredPermissions := []string{"elasticache:DescribeCacheClusters"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getElastiCacheNodesUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	nodeNames := []*elasticache.CacheCluster{}
	err := conf.ElastiCache.DescribeCacheClustersPages(&elasticache.DescribeCacheClustersInput{}, func(p *elasticache.DescribeCacheClustersOutput, lastPage bool) bool {
		nodeNames = append(nodeNames, p.CacheClusters...)
		return true // continue paging
	})
	quotaInfo := c.GetAllAppliedQuotas()["Nodes per Region"]

	if err != nil {
		fmt.Printf("failed to retrieve elasticache nodes, %v", err)
		return
	}

	quotaInfo.UsageValue = float64(len(nodeNames))
	ret = append(ret, quotaInfo)
	return
}
