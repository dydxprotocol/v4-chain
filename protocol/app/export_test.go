package app_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
)

func TestExportAppStateAndValidators_Panics(t *testing.T) {
	klyraApp := app.DefaultTestApp(nil)
	require.Panics(t, func() { klyraApp.ExportAppStateAndValidators(false, nil, nil) }) // nolint:errcheck
}
