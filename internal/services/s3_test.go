package services_test

import (
	"testing"

	"github.com/sebasrp/awslimitchecker/internal/services"
	"github.com/stretchr/testify/require"
)

func TestS3CheckerImpl(t *testing.T) {
	require.Implements(t, (*services.Svcquota)(nil), new(services.S3Checker))
}
