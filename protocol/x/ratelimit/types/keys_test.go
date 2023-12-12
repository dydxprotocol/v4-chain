package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "ratelimit", types.ModuleName)
	require.Equal(t, "ratelimit", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	// TODO(CORE-824): test state keys
}
