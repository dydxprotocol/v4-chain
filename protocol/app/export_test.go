package app_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExportAppStateAndValidators_Panics(t *testing.T) {
	dydxApp := app.DefaultTestApp(nil)
	require.Panics(t, func() { dydxApp.ExportAppStateAndValidators(false, nil, nil) }) // nolint:errcheck
}
