package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "assets", types.ModuleName)
	require.Equal(t, "assets", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "Asset:", types.AssetKeyPrefix)
}
