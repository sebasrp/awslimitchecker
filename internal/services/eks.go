package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/eks"
)

type EksClientInterface interface {
	ListClustersPages(input *eks.ListClustersInput, fn func(*eks.ListClustersOutput, bool) bool) error
}

func NewEksChecker() Svcquota {
	serviceCode := "eks"
	supportedQuotas := map[string]func(ServiceChecker) (ret AWSQuotaInfo){
		"Clusters": ServiceChecker.getEKSClusterUsage,
	}
	requiredPermissions := []string{"dynamodb:ListClusters"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getEKSClusterUsage() (ret AWSQuotaInfo) {
	clusterNames := []*string{}
	err := conf.Eks.ListClustersPages(&eks.ListClustersInput{}, func(o *eks.ListClustersOutput, lastPage bool) bool {
		clusterNames = append(clusterNames, o.Clusters...)
		return true // continue paging
	})
	ret = c.GetAllDefaultQuotas()["Clusters"]

	if err != nil {
		fmt.Printf("failed to retrieve eks clusters, %v", err)
		return
	}

	ret.UsageValue = float64(len(clusterNames))
	return
}
