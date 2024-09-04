package types_test

import (
	"testing"

	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"
)

func TestModuleKeys(t *testing.T) {
	require.Equal(t, "ratelimit", types.ModuleName)
	require.Equal(t, "ratelimit", types.StoreKey)
}

func TestStateKeys(t *testing.T) {
	require.Equal(t, "DenomCapacity:", types.DenomCapacityKeyPrefix)
	require.Equal(t, "LimitParams:", types.LimitParamsKeyPrefix)
}

func TestSplitPendingSendPacketKey(t *testing.T) {
	channelId := "channel-0"
	sequenceNumber := uint64(2)
	channelIdReceived, sequenceNumberReceived, err := types.SplitPendingSendPacketKey(
		types.GetPendingSendPacketKey(channelId, sequenceNumber),
	)
	require.NoError(t, err)
	require.Equal(t, channelId, channelIdReceived)
	require.Equal(t, sequenceNumber, sequenceNumberReceived)
}

func TestAccountConstants(t *testing.T) {
	require.Equal(t, "transfer", types.TDaiPoolAccount)
	require.Equal(t, "sDAIPoolAccount", types.SDaiPoolAccount)
}

func TestSDaiConstants(t *testing.T) {
	require.Equal(t, "ibc/DEEFE2DEFDC8EA8879923C4CCA42BB888C3CD03FF7ECFEFB1C2FEC27A732ACC8", types.SDaiDenom)
	require.Equal(t, "gsdai", types.SDaiBaseDenom)
	require.Equal(t, "transfer/channel-0", types.SDaiBaseDenomPathPrefix)
	require.Equal(t, "transfer/channel-0/gsdai", types.SDaiBaseDenomFullPath)
	require.Equal(t, -18, types.SDaiDenomExponent)
}

func TestTDaiConstants(t *testing.T) {
	require.Equal(t, assettypes.TDaiDenom, types.TDaiDenom)
}

func TestKeyPrefixes(t *testing.T) {
	require.Equal(t, "SDAIPrice:", types.SDaiKeyPrefix)
	require.Equal(t, "AssetYieldIndex:", types.AssetYieldIndexPrefix)
}
