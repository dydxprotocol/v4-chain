package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestGetClobPairSubticksPerTick(t *testing.T) {
	clobPair := types.ClobPair{
		SubticksPerTick: uint32(100),
	}

	subticksPerTick := clobPair.GetClobPairSubticksPerTick()
	require.Equal(t, types.SubticksPerTick(100), subticksPerTick)
}

func TestGetClobPairMinOrderBaseQuantums(t *testing.T) {
	clobPair := types.ClobPair{
		StepBaseQuantums: uint64(100),
	}

	minOrderBaseQuantums := clobPair.GetClobPairMinOrderBaseQuantums()
	require.Equal(t, satypes.BaseQuantums(100), minOrderBaseQuantums)
}

func TestGetPerpetualId(t *testing.T) {
	perpetualId, err := constants.ClobPair_Eth.GetPerpetualId()
	require.Equal(t, uint32(1), perpetualId)
	require.NoError(t, err)

	perpetualId, err = constants.ClobPair_Asset.GetPerpetualId()
	require.Equal(t, uint32(0), perpetualId)
	require.ErrorIs(t, types.ErrAssetOrdersNotImplemented, err)
}

func TestIsSupportedClobPairStatus_Supported(t *testing.T) {
	// these are the only two supported statuses
	require.True(t, types.IsSupportedClobPairStatus(types.ClobPair_STATUS_ACTIVE))
	require.True(t, types.IsSupportedClobPairStatus(types.ClobPair_STATUS_INITIALIZING))
}

func TestIsSupportedClobPairStatus_Unsupported(t *testing.T) {
	// out of bounds of the clob pair status enum
	require.False(t, types.IsSupportedClobPairStatus(types.ClobPair_Status(100)))

	// these are part of the ClobPair_Status enum but are not supported
	require.False(t, types.IsSupportedClobPairStatus(types.ClobPair_STATUS_UNSPECIFIED))
	require.False(t, types.IsSupportedClobPairStatus(types.ClobPair_STATUS_PAUSED))
	require.False(t, types.IsSupportedClobPairStatus(types.ClobPair_STATUS_CANCEL_ONLY))
	require.False(t, types.IsSupportedClobPairStatus(types.ClobPair_STATUS_POST_ONLY))
}

func TestIsSupportedClobPairStatusTransition_Supported(t *testing.T) {
	// only supported transition
	require.True(t, types.IsSupportedClobPairStatusTransition(
		types.ClobPair_STATUS_INITIALIZING, types.ClobPair_STATUS_ACTIVE,
	))
}

func TestIsSupportedClobPairStatusTransition_Unsupported(t *testing.T) {
	// out of bounds of the clob pair status enum
	require.False(t, types.IsSupportedClobPairStatusTransition(
		types.ClobPair_Status(100), types.ClobPair_Status(100),
	))
	require.False(t, types.IsSupportedClobPairStatusTransition(
		types.ClobPair_STATUS_ACTIVE, types.ClobPair_Status(100),
	))
	require.False(t, types.IsSupportedClobPairStatusTransition(
		types.ClobPair_Status(100), types.ClobPair_STATUS_ACTIVE,
	))

	// iterate over all permutations of clob pair statuses
	for _, fromClobPairStatus := range types.ClobPair_Status_value {
		for _, toClobPairStatus := range types.ClobPair_Status_value {
			if fromClobPairStatus == int32(types.ClobPair_STATUS_INITIALIZING) &&
				toClobPairStatus == int32(types.ClobPair_STATUS_ACTIVE) {
				continue
			} else {
				require.False(
					t,
					types.IsSupportedClobPairStatusTransition(
						types.ClobPair_Status(fromClobPairStatus),
						types.ClobPair_Status(toClobPairStatus),
					),
				)
			}
		}
	}
}
