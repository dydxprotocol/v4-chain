package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/mock"
)

func TestAddSeenPlaceOrder_PruneExpiredSeenPlaceOrders(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup              func(ctx sdk.Context, k keeper.Keeper)
		blockHeightToPrune uint32

		// Expectations before prune
		expectedSeenPlaceOrderIds           map[types.Order]bool
		expectedSeenPlaceOrderIdsAfterPrune map[types.Order]bool
	}{
		"Reads an empty state": {
			setup:              func(ctx sdk.Context, k keeper.Keeper) {},
			blockHeightToPrune: 0,

			expectedSeenPlaceOrderIds:           map[types.Order]bool{},
			expectedSeenPlaceOrderIdsAfterPrune: map[types.Order]bool{},
		},
		"Adds a single seen short-term place order and prunes it": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					},
				)
			},
			blockHeightToPrune: 15,

			expectedSeenPlaceOrderIds: map[types.Order]bool{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15: true,
			},
			expectedSeenPlaceOrderIdsAfterPrune: map[types.Order]bool{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15: false,
			},
		},
		"Adds multiple seen short-term place orders and prunes some": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					},
				)
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15,
					},
				)
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
					},
				)
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31,
					},
				)
			},
			blockHeightToPrune: 15,

			expectedSeenPlaceOrderIds: map[types.Order]bool{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15:  true,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20:  true,
				constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15: true,
				constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15: false,
				constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31: true,
				constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32: false,
			},
			expectedSeenPlaceOrderIdsAfterPrune: map[types.Order]bool{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15:  true,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20:  true,
				constants.Order_Alice_Num0_Id1_Clob0_Sell5_Price15_GTB15: false,
				constants.Order_Alice_Num0_Id2_Clob1_Sell5_Price10_GTB15: false,
				constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB31: true,
				constants.Order_Alice_Num1_Id10_Clob0_Buy5_Price30_GTB32: false,
			},
		},
		"Adds a short-term order, adds again with a lower GTB, then prunes the lower GTB": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
					},
				)
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15,
					},
				)
			},
			blockHeightToPrune: 15,

			expectedSeenPlaceOrderIds: map[types.Order]bool{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20: true,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15: true,
			},
			expectedSeenPlaceOrderIdsAfterPrune: map[types.Order]bool{
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20: true,
				constants.Order_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB15: true,
			},
		},
		"Adds a long-term seen place orders (which is not supported)": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.AddSeenPlaceOrder(
					ctx,
					types.MsgPlaceOrder{
						Order: constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					},
				)
			},
			blockHeightToPrune: 15,

			expectedSeenPlaceOrderIds:           map[types.Order]bool{},
			expectedSeenPlaceOrderIdsAfterPrune: map[types.Order]bool{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ctx,
				keeper,
				_, _, _, _, _, _ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			tc.setup(ctx, *keeper)

			for order, expectedIsSeen := range tc.expectedSeenPlaceOrderIds {
				require.Equal(
					t,
					expectedIsSeen,
					keeper.HasSeenPlaceOrder(
						ctx,
						types.MsgPlaceOrder{Order: order},
					),
				)
			}

			// Prune expired place orders
			keeper.PruneExpiredSeenPlaceOrders(ctx, tc.blockHeightToPrune)

			for order, expectedIsSeen := range tc.expectedSeenPlaceOrderIdsAfterPrune {
				require.Equal(
					t,
					expectedIsSeen,
					keeper.HasSeenPlaceOrder(
						ctx,
						types.MsgPlaceOrder{Order: order},
					),
				)
			}
		})
	}
}
