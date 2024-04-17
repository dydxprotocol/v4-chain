package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "subaccounts", types.ModuleName)
	require.Equal(t, "subaccounts", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "SA:", types.SubaccountKeyPrefix)
}
