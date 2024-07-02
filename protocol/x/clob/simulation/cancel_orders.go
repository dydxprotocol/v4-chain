package simulation

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const (
	// Have some cancels submit an invalid cancel.
	invalidGoodTilBlockWeight = 5

	// Ensure most cancels use the correct GoodTilBlock, but have a lower percentage of cancels
	// use a GTB that is lower than the order's. These will not successfully cancel the order.
	uncancellableGoodTilBlockWeight = 5
)

var (
	typeMsgCancelOrder = sdk.MsgTypeURL(&types.MsgCancelOrder{})
)

func SimulateMsgCancelOrder(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.SubaccountsKeeper,
	k keeper.Keeper,
	memClob types.MemClob,
	cdc *codec.ProtoCodec,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get a random subaccount.
		subAccount, err := sk.GetRandomSubaccount(ctx, r)
		if err != nil {
			panic(fmt.Errorf("SimulateMsgCancelOrder: Simulation has no subaccounts available"))
		}
		subaccountId := *subAccount.GetId()

		// Get all clob pairs.
		clobPairs := k.GetAllClobPairs(ctx)
		if len(clobPairs) < 1 {
			panic(fmt.Errorf("SimulateMsgCancelOrder: Simulation has no CLOB pairs available"))
		}

		// Shuffle clob pair slice, so we randomize search for open orders.
		randomizedClobPairs := sim_helpers.RandSliceShuffle(r, clobPairs)

		// Shuffle order side slice, so we randomize search for open orders.
		randomizedOrderSide := sim_helpers.RandSliceShuffle(
			r,
			[]types.Order_Side{types.Order_SIDE_BUY, types.Order_SIDE_SELL},
		)

		// Search all clob pair + side for an order to cancel.
		// TODO(DEC-1643): Support long-term orders.
		openOrderToCancel, found := findOpenShortTermOrderToCancel(
			ctx,
			r,
			memClob,
			randomizedClobPairs,
			randomizedOrderSide,
			subaccountId,
		)

		if !found {
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgPlaceOrder,
				"No open, short-term orders found for subaccount",
			), nil, nil
		}

		proposer, _ := simtypes.FindAccount(accs, subaccountId.MustGetAccAddress())

		// By default, use a random GoodTilBlock that would successfully cancel the order.
		var goodTilBlock uint32 = uint32(simtypes.RandIntBetween(
			r,
			int(openOrderToCancel.GetGoodTilBlock()),
			int(openOrderToCancel.GetGoodTilBlock()+types.ShortBlockWindow),
		))

		// Determine if the cancel should be invalid.
		invalidGTBIndex := simtypes.RandIntBetween(r, 0, 100)
		if invalidGTBIndex < invalidGoodTilBlockWeight {
			goodTilBlock = uint32(ctx.BlockHeight()) + types.ShortBlockWindow + 1
		}

		// Determine if the cancel should use a GoodTilBlock that will not cancel the order.
		uncancellableGTBIndex := simtypes.RandIntBetween(r, 0, 100)
		if uncancellableGTBIndex < uncancellableGoodTilBlockWeight {
			goodTilBlock = openOrderToCancel.GetGoodTilBlock() - 1
		}

		// Generate a cancel that matches the corresponding open order.
		msg := &types.MsgCancelOrder{
			OrderId:      openOrderToCancel.OrderId,
			GoodTilOneof: &types.MsgCancelOrder_GoodTilBlock{GoodTilBlock: goodTilBlock},
		}

		opMsg, err := sim_helpers.GenerateAndCheckTx(
			r,
			app,
			ctx,
			chainID,
			cdc,
			ak,
			bk,
			proposer,
			types.ModuleName,
			msg,
			typeMsgCancelOrder,
			true, // MsgCancelOrder should be zero fees.
		)
		if err != nil {
			switch {
			case errors.Is(err, types.ErrMemClobCancelAlreadyExists),
				errors.Is(err, types.ErrHeightExceedsGoodTilBlock),
				errors.Is(err, types.ErrGoodTilBlockExceedsShortBlockWindow):
				// These errors are expected, and can occur during normal operation. We shouldn't panic on them.
			default:
				panic(err)
			}
		}

		return opMsg, nil, nil
	}
}

// Searches each memclob clob pair + order side and returns a random open, short-term order found for a subaccount.
func findOpenShortTermOrderToCancel(
	ctx sdk.Context,
	r *rand.Rand,
	memClob types.MemClob,
	clobPairs []types.ClobPair,
	orderSides []types.Order_Side,
	subaccountId satypes.SubaccountId,
) (order types.Order, found bool) {
	for _, clobPair := range clobPairs {
		for _, orderSide := range orderSides {
			openOrders, err := memClob.GetSubaccountOrders(
				clobPair.GetClobPairId(),
				subaccountId,
				orderSide,
			)
			if err != nil {
				panic(err)
			}

			// Filter out any long-term orders
			shortTermOpenOrders := make([]types.Order, 0)
			for _, openOrder := range openOrders {
				if openOrder.OrderId.IsShortTermOrder() {
					shortTermOpenOrders = append(shortTermOpenOrders, openOrder)
				}
			}

			// Pick a random open order to cancel
			if len(shortTermOpenOrders) > 0 {
				shortTermOpenOrderIndex := simtypes.RandIntBetween(r, 0, len(shortTermOpenOrders))
				return shortTermOpenOrders[shortTermOpenOrderIndex], true
			}
		}
	}

	return types.Order{}, false
}
