package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "vest", types.ModuleName)
	require.Equal(t, "vest", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "Entry:", types.VestEntryKeyPrefix)
}

func TestModuleAccountKeys(t *testing.T) {
	require.Equal(t, "community_treasury", types.CommunityTreasuryAccountName)
	require.Equal(t, "community_vester", types.CommunityVesterAccountName)
}
