package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"

	"github.com/stretchr/testify/require"
)

func TestMsgUpdateMarketPrices(t *testing.T) {
	update := types.NewMarketPriceUpdate(uint32(0), uint64(1))
	updates := []*types.MsgUpdateMarketPrices_MarketPrice{update}
	msg := types.NewMsgUpdateMarketPrices(updates)

	require.Equal(t, uint32(0), update.MarketId)
	require.Equal(t, uint64(1), update.Price)
	require.Equal(t, updates, msg.MarketPriceUpdates)
}

func TestMsgUpdateMarketPrices_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		updates []*types.MsgUpdateMarketPrices_MarketPrice

		expectedErr string
	}{
		"Error: price cannot be zero": {
			updates: []*types.MsgUpdateMarketPrices_MarketPrice{
				{Price: 0, MarketId: 1},
			},
			expectedErr: "price cannot be 0 for market id (1): Market price update is invalid: stateless.",
		},
		"Error: duplicate market ids": {
			updates: []*types.MsgUpdateMarketPrices_MarketPrice{
				{Price: 1_000, MarketId: 1},
				{Price: 2_000, MarketId: 2},
				{Price: 3_000, MarketId: 1}, // duplicate
			},
			expectedErr: "market price updates must be sorted by market id in ascending" +
				" order and cannot contain duplicates: Market price update is invalid: stateless.",
		},
		"Error: descending market ids": {
			updates: []*types.MsgUpdateMarketPrices_MarketPrice{
				{Price: 2_000, MarketId: 2},
				{Price: 3_000, MarketId: 3},
				{Price: 4_000, MarketId: 4},
				{Price: 1_000, MarketId: 1}, // descending
			},
			expectedErr: "market price updates must be sorted by market id in ascending" +
				" order and cannot contain duplicates: Market price update is invalid: stateless.",
		},
		"No error: empty price updates": {
			updates: []*types.MsgUpdateMarketPrices_MarketPrice{},
		},
		"No error: valid ordering": {
			updates: []*types.MsgUpdateMarketPrices_MarketPrice{
				{Price: 1_000, MarketId: 1},
				{Price: 3_000, MarketId: 3},
				{Price: 99_000, MarketId: 99},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := types.NewMsgUpdateMarketPrices(tc.updates)
			err := msg.ValidateBasic()
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
