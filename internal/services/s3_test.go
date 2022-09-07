package services_test

import (
	"testing"

	"github.com/sebasrp/awslimitchecker/internal/services"
	"github.com/stretchr/testify/require"
)

func TestNewS3CheckerImpl(t *testing.T) {
	require.Implements(t, (*services.Svcquota)(nil), services.NewS3Checker(nil, nil))
}
