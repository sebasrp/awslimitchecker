package services_test

import (
	"testing"

	"github.com/sebasrp/awslimitchecker/internal/services"
	"github.com/stretchr/testify/require"
)

func TestNewDynamoDbCheckerImpl(t *testing.T) {
	require.Implements(t, (*services.Svcquota)(nil), services.NewDynamoDbChecker(nil))
}
