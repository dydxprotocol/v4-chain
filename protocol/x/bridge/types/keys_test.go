package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "bridge", types.ModuleName)
	require.Equal(t, "bridge", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "AckEventInfo", types.AcknowledgedEventInfoKey)
	require.Equal(t, "EventParams", types.EventParamsKey)
	require.Equal(t, "ProposeParams", types.ProposeParamsKey)
	require.Equal(t, "SafetyParams", types.SafetyParamsKey)
}
