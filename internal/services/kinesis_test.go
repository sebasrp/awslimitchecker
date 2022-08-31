package services_test

import (
	"testing"

	"github.com/sebasrp/awslimitchecker/internal/services"
	"github.com/stretchr/testify/require"
)

func TestKinesisCheckerImpl(t *testing.T) {
	require.Implements(t, (*services.Svcquota)(nil), new(services.KinesisChecker))
}
