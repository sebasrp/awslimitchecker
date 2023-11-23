package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamodbClientInterface interface {
	ListTablesPages(input *dynamodb.ListTablesInput, fn func(*dynamodb.ListTablesOutput, bool) bool) error
}

func NewDynamoDbChecker() ServiceQuota {
	serviceCode := "dynamodb"
	supportedQuotas := map[string]func(ServiceChecker) (ret []AWSQuotaInfo){
		"Maximum number of tables": ServiceChecker.getDynanoDBTableUsage,
	}
	requiredPermissions := []string{"dynamodb:ListTables"}

	return NewServiceChecker(serviceCode, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getDynanoDBTableUsage() (ret []AWSQuotaInfo) {
	ret = []AWSQuotaInfo{}
	tableNames := []*string{}
	err := conf.DynamoDb.ListTablesPages(&dynamodb.ListTablesInput{}, func(p *dynamodb.ListTablesOutput, lastPage bool) bool {
		tableNames = append(tableNames, p.TableNames...)
		return true // continue paging
	})
	quotaInfo := c.GetAllAppliedQuotas()["Maximum number of tables"]

	if err != nil {
		fmt.Printf("failed to retrieve dynamodb tables, %v", err)
		return
	}

	quotaInfo.UsageValue = float64(len(tableNames))
	ret = append(ret, quotaInfo)
	return
}
