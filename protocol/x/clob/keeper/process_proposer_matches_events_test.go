package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	"github.com/dydxprotocol/v4/testutil/tracer"
	"github.com/dydxprotocol/v4/x/clob/keeper"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSetProcessProposerMatchesEvents(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		setup func(ctx sdk.Context, k keeper.Keeper)

		// Expectations.
		expectedProcessProposerMatchesEvents types.ProcessProposerMatchesEvents
		expectedNumWrites                    uint
	}{
		"Can set and get process proposer matches events": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders: []types.Order{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						OrdersIdsFilledInLastBlock: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
						},
					},
				)
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
				},
				ExpiredStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
				},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
				},
			},
			expectedNumWrites: 1,
		},
		"Can get empty slice for process proposer matches events": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{},
			expectedNumWrites:                    0,
		},
		"Can set and get empty slice for process proposer matches events": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders:       []types.Order{},
						ExpiredStatefulOrderIds:    []types.OrderId{},
						OrdersIdsFilledInLastBlock: []types.OrderId{},
					},
				)
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{},
			expectedNumWrites:                    1,
		},
		"Can overwrite process proposer matches events": {
			setup: func(ctx sdk.Context, k keeper.Keeper) {
				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders: []types.Order{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						OrdersIdsFilledInLastBlock: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
						},
					},
				)

				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						PlacedStatefulOrders: []types.Order{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
						},
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
						},
						OrdersIdsFilledInLastBlock: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
						},
					},
				)
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				PlacedStatefulOrders: []types.Order{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10,
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15,
				},
				ExpiredStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				},
				OrdersIdsFilledInLastBlock: []types.OrderId{
					constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
				},
			},
			expectedNumWrites: 2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := &mocks.MemClob{}
			memClob.On("SetClobKeeper", mock.Anything).Return()
			ks := keepertest.NewClobKeepersTestContext(
				t,
				memClob,
				&mocks.BankKeeper{},
				&mocks.IndexerEventManager{},
			)

			// Set the tracer on the multistore to verify the performed writes are correct.
			traceDecoder := &tracer.TraceDecoder{}
			ks.Ctx.MultiStore().SetTracer(traceDecoder)

			tc.setup(ks.Ctx, *ks.ClobKeeper)

			processProposerMatchesEvents := ks.ClobKeeper.GetProcessProposerMatchesEvents(ks.Ctx)

			require.Equal(t, tc.expectedProcessProposerMatchesEvents, processProposerMatchesEvents)

			expectedMultistoreWrites := make([]string, tc.expectedNumWrites)
			for i := 0; i < int(tc.expectedNumWrites); i++ {
				expectedMultistoreWrites[i] = "ProcessProposerMatchesEvents/value"
			}

			traceDecoder.RequireKeyPrefixWrittenInSequence(
				t,
				expectedMultistoreWrites,
			)
		})
	}
}

func TestSetProcessProposerMatchesEvents_BadBlockHeight(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)
	ctx := ks.Ctx.WithBlockHeight(1)
	require.Panics(t, func() {
		ks.ClobKeeper.MustSetProcessProposerMatchesEvents(ctx, types.ProcessProposerMatchesEvents{
			BlockHeight: 5,
		})
	})
}
