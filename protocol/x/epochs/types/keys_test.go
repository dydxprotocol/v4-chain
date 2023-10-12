package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "epochs", types.ModuleName)
	require.Equal(t, "epochs", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "Info:", types.EpochInfoKeyPrefix)
}
