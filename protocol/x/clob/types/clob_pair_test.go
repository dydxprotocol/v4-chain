package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"

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
		MinOrderBaseQuantums: uint64(100),
	}

	minOrderBaseQuantums := clobPair.GetClobPairMinOrderBaseQuantums()
	require.Equal(t, satypes.BaseQuantums(100), minOrderBaseQuantums)
}

func TestGetFeePpm(t *testing.T) {
	makerFeePpm := uint32(500)
	takerFeePpm := uint32(1000)

	clobPair := types.ClobPair{
		MakerFeePpm: makerFeePpm,
		TakerFeePpm: takerFeePpm,
	}

	require.Equal(t, takerFeePpm, clobPair.GetFeePpm(true))
	require.Equal(t, makerFeePpm, clobPair.GetFeePpm(false))
}

func TestGetPerpetualId(t *testing.T) {
	perpetualId, err := constants.ClobPair_Eth.GetPerpetualId()
	require.Equal(t, uint32(1), perpetualId)
	require.NoError(t, err)

	perpetualId, err = constants.ClobPair_Asset.GetPerpetualId()
	require.Equal(t, uint32(0), perpetualId)
	require.ErrorIs(t, types.ErrAssetOrdersNotImplemented, err)
}
