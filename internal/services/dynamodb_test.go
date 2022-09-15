package services

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockedListTablesPagesMsgs struct {
	dynamodbiface.DynamoDBAPI
	Resp  dynamodb.ListTablesOutput
	Error error
}

func (m mockedListTablesPagesMsgs) ListTablesPages(
	input *dynamodb.ListTablesInput,
	fn func(*dynamodb.ListTablesOutput, bool) bool) error {
	fn(&m.Resp, false)
	return m.Error
}

func TestNewDynamoDbCheckerImpl(t *testing.T) {
	require.Implements(t, (*Svcquota)(nil), NewDynamoDbChecker())
}

func TestGetDynanoDBTableUsage(t *testing.T) {
	mockedOutput := dynamodb.ListTablesOutput{
		TableNames: []*string{aws.String("table1"), aws.String("table2")},
	}
	conf.DynamoDb = mockedListTablesPagesMsgs{Resp: mockedOutput}

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("dynamodb", "Maximum number of tables", float64(2500), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	ddbChecker := NewDynamoDbChecker()
	svcChecker := ddbChecker.(*ServiceChecker)
	actual := svcChecker.getDynanoDBTableUsage()

	assert.Len(t, actual, 1)
	quota := actual[0]
	assert.Equal(t, "dynamodb", quota.Service)
	assert.Equal(t, float64(2500), quota.QuotaValue)
	assert.Equal(t, float64(len(mockedOutput.TableNames)), quota.UsageValue)
}

func TestGetDynanoDBTableUsageError(t *testing.T) {
	mockedOutput := dynamodb.ListTablesOutput{
		TableNames: []*string{aws.String("table1"), aws.String("table2")},
	}
	conf.DynamoDb = mockedListTablesPagesMsgs{Resp: mockedOutput, Error: errors.New("test error")}

	mockedSvcQuotaOutput := servicequotas.ListServiceQuotasOutput{
		Quotas: []*servicequotas.ServiceQuota{
			NewQuota("dynamodb", "Maximum number of tables", float64(2500), false),
		},
	}
	conf.ServiceQuotas = mockedScvQuotaClient{
		ListServiceQuotasOutputResp: mockedSvcQuotaOutput,
	}

	ddbChecker := NewDynamoDbChecker()
	svcChecker := ddbChecker.(*ServiceChecker)
	actual := svcChecker.getDynanoDBTableUsage()
	expected := []AWSQuotaInfo{}
	assert.Equal(t, expected, actual)
}
