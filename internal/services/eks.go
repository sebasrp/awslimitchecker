package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
)

type EksClientInterface interface {
	ListClustersPages(input *eks.ListClustersInput, fn func(*eks.ListClustersOutput, bool) bool) error
	ListNodegroupsPages(input *eks.ListNodegroupsInput, fn func(*eks.ListNodegroupsOutput, bool) bool) error
}

func NewEksChecker() ServiceQuota {
	serviceCode := "eks"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Clusters":                        ServiceChecker.getEKSClusterUsage,
		"Managed node groups per cluster": ServiceChecker.getEKSNodeGroupsPerClusterUsage,
	}
	requiredPermissions := []string{"eks:ListClusters", "eks:ListNodegroups"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getEKSClusterUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	clusterNames := []*string{}
	err := conf.Eks.ListClustersPages(&eks.ListClustersInput{}, func(o *eks.ListClustersOutput, lastPage bool) bool {
		clusterNames = append(clusterNames, o.Clusters...)
		return true // continue paging
	})
	quotaInfo := c.GetAllAppliedQuotas()["Clusters"]

	if err != nil {
		fmt.Printf("failed to retrieve eks clusters, %v", err)
		return
	}

	quotaInfo.UsageValue = float64(len(clusterNames))
	ret = append(ret, quotaInfo)
	return
}

func (c ServiceChecker) getEKSNodeGroupsPerClusterUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	clusterNames := []*string{}
	errListClusters := conf.Eks.ListClustersPages(&eks.ListClustersInput{}, func(o *eks.ListClustersOutput, lastPage bool) bool {
		clusterNames = append(clusterNames, o.Clusters...)
		return true // continue paging
	})
	if errListClusters != nil {
		fmt.Printf("failed to retrieve eks clusters, %v", errListClusters)
		return
	}

	for _, cluster := range clusterNames {
		nodegroups := []*string{}
		quotaInfo := c.GetAllAppliedQuotas()["Managed node groups per cluster"]
		errListNodeGroups := conf.Eks.ListNodegroupsPages(&eks.ListNodegroupsInput{ClusterName: cluster}, func(o *eks.ListNodegroupsOutput, lastPage bool) bool {
			nodegroups = append(nodegroups, o.Nodegroups...)
			return true // continue paging
		})
		if errListNodeGroups != nil {
			fmt.Printf("failed to retrieve nodegroups for cluster %s, %v", *cluster, errListClusters)
			continue
		}

		quotaInfo.UsageValue = float64(len(nodegroups))
		quotaInfo.ResourceId = fmt.Sprintf("AWS::EKS::Cluster::%s", *cluster)
		ret = append(ret, quotaInfo)
	}
	return
}
