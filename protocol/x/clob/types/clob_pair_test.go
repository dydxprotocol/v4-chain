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
