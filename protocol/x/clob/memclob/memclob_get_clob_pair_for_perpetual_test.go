package memclob

import (
	"testing"

	"github.com/dydxprotocol/v4/testutil/constants"
	sdktest "github.com/dydxprotocol/v4/testutil/sdk"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGetClobPairForPerpetual_Success(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()

	memclob := NewMemClobPriceTimePriority(false)

	// Create the orderbook.
	clobPair := constants.ClobPair_Btc
	memclob.CreateOrderbook(ctx, clobPair)

	clobPairId, err := memclob.GetClobPairForPerpetual(ctx, 0)
	require.NoError(t, err)
	require.Equal(t, clobPair.GetClobPairId(), clobPairId)
}

func TestGetClobPairForPerpetual_SuccessMultipleClobPairIds(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()

	memclob := NewMemClobPriceTimePriority(false)

	// Create multiple orderbooks.
	clobPair_Btc2 := types.ClobPair{
		Id:                   100,
		SubticksPerTick:      120,
		MinOrderBaseQuantums: 1,
		Metadata: &types.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &types.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
	}

	clobPairs := []types.ClobPair{
		constants.ClobPair_Btc,
		clobPair_Btc2,
		constants.ClobPair_Eth,
	}
	for _, clobPair := range clobPairs {
		memclob.CreateOrderbook(ctx, clobPair)
	}

	clobPairId, err := memclob.GetClobPairForPerpetual(ctx, 0)
	require.NoError(t, err)
	// The first CLOB pair ID should be returned.
	require.Equal(t, types.ClobPairId(clobPairs[0].Id), clobPairId)
}

func TestGetClobPairForPerpetual_ErrorNoClobPair(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()

	memclob := NewMemClobPriceTimePriority(false)

	_, err := memclob.GetClobPairForPerpetual(ctx, 0)
	require.EqualError(
		t,
		err,
		"Perpetual ID 0 has no associated CLOB pairs: The provided perpetual ID "+
			"does not have any associated CLOB pairs",
	)
}

func TestGetClobPairForPerpetual_PanicsEmptyClobPair(t *testing.T) {
	ctx, _, _ := sdktest.NewSdkContextWithMultistore()

	memclob := NewMemClobPriceTimePriority(false)

	perpetualId := uint32(0)
	memclob.perpetualIdToClobPairId[perpetualId] = make([]types.ClobPairId, 0)

	require.Panics(t, func() {
		//nolint:errcheck
		memclob.GetClobPairForPerpetual(ctx, perpetualId)
	})
}
