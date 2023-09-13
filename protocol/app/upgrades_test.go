package app_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/stretchr/testify/require"
)

func TestDefaultUpgradesAndForks(t *testing.T) {
	require.Len(t, app.Upgrades, 1, "Expected 1 upgrade")
	require.Len(t, app.Forks, 1, "Expected 1 fork")
}
