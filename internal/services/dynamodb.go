package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var ddbClient DynamodbClientInterface

type DynamodbClientInterface interface {
	ListTablesPages(input *dynamodb.ListTablesInput, fn func(*dynamodb.ListTablesOutput, bool) bool) error
}

func NewDynamoDbChecker(session *session.Session, svcQuotaClient SvcQuotaClientInterface) Svcquota {
	serviceCode := "dynamodb"
	supportedQuotas := map[string]func(ServiceChecker) (ret AWSQuotaInfo){
		"Maximum number of tables": ServiceChecker.getDynanoDBTableUsage,
	}
	requiredPermissions := []string{"dynamodb:ListTables"}

	return NewServiceChecker(serviceCode, session, svcQuotaClient, supportedQuotas, requiredPermissions)
}

func (c ServiceChecker) getDynanoDBTableUsage() (ret AWSQuotaInfo) {
	tableNames := []*string{}

	if ddbClient == nil && c.session != nil {
		ddbClient = dynamodb.New(c.session)
	}

	err := ddbClient.ListTablesPages(&dynamodb.ListTablesInput{}, func(p *dynamodb.ListTablesOutput, lastPage bool) bool {
		tableNames = append(tableNames, p.TableNames...)
		return true // continue paging
	})
	if err != nil {
		fmt.Printf("failed to retrieve dynamodb tables, %v", err)
		return
	}
	ret = c.GetAllDefaultQuotas()["Maximum number of tables"]
	ret.UsageValue = float64(len(tableNames))
	return
}
