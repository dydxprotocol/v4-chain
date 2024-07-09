package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/tracer"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
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
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						OrderIdsFilledInLastBlock: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
						},
					},
				)
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
				},
				OrderIdsFilledInLastBlock: []types.OrderId{
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
						ExpiredStatefulOrderIds:   []types.OrderId{},
						OrderIdsFilledInLastBlock: []types.OrderId{},
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
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
						},
						OrderIdsFilledInLastBlock: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT20.OrderId,
						},
					},
				)

				k.MustSetProcessProposerMatchesEvents(
					ctx,
					types.ProcessProposerMatchesEvents{
						ExpiredStatefulOrderIds: []types.OrderId{
							constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
							constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
						},
						OrderIdsFilledInLastBlock: []types.OrderId{
							constants.LongTermOrder_Alice_Num0_Id1_Clob1_Sell65_Price15_GTBT25.OrderId,
						},
					},
				)
			},

			expectedProcessProposerMatchesEvents: types.ProcessProposerMatchesEvents{
				ExpiredStatefulOrderIds: []types.OrderId{
					constants.LongTermOrder_Bob_Num0_Id0_Clob0_Buy25_Price30_GTBT10.OrderId,
					constants.ConditionalOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTBT15_StopLoss20.OrderId,
				},
				OrderIdsFilledInLastBlock: []types.OrderId{
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
				expectedMultistoreWrites[i] = types.ProcessProposerMatchesEventsKey
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

func TestOrderedOrderIdList(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	prefixKey := "p"
	indexKey := "i"
	orderIds := []types.OrderId{
		{
			ClientId: 0,
		},
		{
			ClientId: 1,
		},
		{
			ClientId: 2,
		},
		{
			ClientId: 3,
		},
	}

	memStore := ks.Ctx.KVStore(ks.MemKey)

	for _, orderId := range orderIds {
		ks.ClobKeeper.AppendOrderedOrderId(ks.Ctx, memStore, prefixKey, indexKey, orderId)
	}
	require.Equal(t, orderIds, ks.ClobKeeper.GetOrderIds(ks.Ctx, memStore, prefixKey))

	ks.ClobKeeper.ResetOrderedOrderIds(ks.Ctx, memStore, prefixKey, indexKey)
	for _, orderId := range orderIds {
		ks.ClobKeeper.AppendOrderedOrderId(ks.Ctx, memStore, prefixKey, indexKey, orderId)
	}
	require.Equal(t, orderIds, ks.ClobKeeper.GetOrderIds(ks.Ctx, memStore, prefixKey))
}

func TestUnorderedOrderIdList(t *testing.T) {
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
	)

	prefixKey := "p"
	indexKey := "i"
	orderIds := []types.OrderId{
		{
			ClientId: 0,
		},
		{
			ClientId: 1,
		},
		{
			ClientId: 2,
		},
		{
			ClientId: 3,
		},
	}

	memStore := ks.Ctx.KVStore(ks.MemKey)

	for _, orderId := range orderIds {
		ks.ClobKeeper.SetUnorderedOrderId(ks.Ctx, memStore, prefixKey, orderId)
	}
	require.Equal(t, orderIds, ks.ClobKeeper.GetOrderIds(ks.Ctx, memStore, prefixKey))

	ks.ClobKeeper.ResetOrderedOrderIds(ks.Ctx, memStore, prefixKey, indexKey)
	for _, orderId := range orderIds {
		ks.ClobKeeper.SetUnorderedOrderId(ks.Ctx, memStore, prefixKey, orderId)
	}
	require.Equal(t, orderIds, ks.ClobKeeper.GetOrderIds(ks.Ctx, memStore, prefixKey))
}
