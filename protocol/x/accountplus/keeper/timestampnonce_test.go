package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/stretchr/testify/require"
)

func TestIsTimestampNonce(t *testing.T) {
	require.True(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff))
	require.True(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff+uint64(1)))
	require.False(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff-uint64(1)))

	require.False(t, keeper.IsTimestampNonce(0))
	require.True(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff+uint64(100000)))
	require.False(t, keeper.IsTimestampNonce(keeper.TimestampNonceSequenceCutoff-uint64(100000)))
}
