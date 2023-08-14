package app_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultUpgradesAndForks(t *testing.T) {
	require.Empty(t, app.Upgrades, "Expected empty upgrades list")
	require.Empty(t, app.Forks, "Expected empty forks list")
}
