package simulation

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const (
	clientIdMaxInt = math.MaxUint32

	// Ensure most orders pass validation and can be placed, but also
	// test undercollateralization sometimes.
	orderProperlyCollateralizedWeight = 98
)

var (
	maxCurrentPositionQuantums = big.NewInt(math.MaxUint32)
	orderLevels                = big.NewInt(4)

	weightedSupportedTimeInForces = map[types.Order_TimeInForce]int{
		types.Order_TIME_IN_FORCE_UNSPECIFIED: 3,
		types.Order_TIME_IN_FORCE_IOC:         1,
		types.Order_TIME_IN_FORCE_POST_ONLY:   1,
	}

	weightedReduceOnly = map[bool]int{
		false: 3,
		true:  1,
	}

	supportedOrderSides = []types.Order_Side{
		types.Order_SIDE_BUY,
		types.Order_SIDE_SELL,
	}
)

var (
	// Maximum position size must fit into a Uint64, so we should upper bound order quote quantums to
	// to prevent potential overflows
	// TODO(DEC-1214): Remove this maximum value once position sizes are stored as arbitrary precision integers in state.
	maxNonOverflowOrderQuoteQuantums = big.NewInt(1000000)
)

var (
	typeMsgPlaceOrder = sdk.MsgTypeURL(&types.MsgPlaceOrder{})
)

func SimulateMsgPlaceOrder(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	sk types.SubaccountsKeeper,
	k keeper.Keeper,
	cdc *codec.ProtoCodec,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// Get a random subaccount.
		subAccount, err := sk.GetRandomSubaccount(ctx, r)
		if err != nil {
			panic(fmt.Errorf("SimulateMsgPlaceOrder: Simulation has no subaccounts available"))
		}
		subaccountId := *subAccount.GetId()

		// Get all clob pairs.
		clobPairs := k.GetAllClobPairs(ctx)
		if len(clobPairs) < 1 {
			panic(fmt.Errorf("SimulateMsgPlaceOrder: Simulation has no CLOB pairs available"))
		}

		// Get a random clob pair.
		clobPairIndex := simtypes.RandIntBetween(r, 0, len(clobPairs))
		clobPair := clobPairs[clobPairIndex]

		// Get subaccount position for clob pair.
		currentPositionSizeQuantums := k.GetStatePosition(
			ctx,
			subaccountId,
			clobPair.GetClobPairId(),
		)
		if currentPositionSizeQuantums.Cmp(maxCurrentPositionQuantums) == 1 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgPlaceOrder,
				"Subaccount position size is already unreasonably high, new orders should not be placed",
			), nil, nil
		}

		bigSubaccountMaxOrderQuoteQuantums := getMaxSubaccountOrderQuoteQuantums(
			ctx,
			sk,
			subaccountId,
		)

		bigMinOrderQuoteQuantums := types.FillAmountToQuoteQuantums(
			types.Subticks(clobPair.SubticksPerTick),
			satypes.BaseQuantums(clobPair.StepBaseQuantums),
			clobPair.QuantumConversionExponent,
		)

		if bigMinOrderQuoteQuantums.Cmp(bigSubaccountMaxOrderQuoteQuantums) == 1 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgPlaceOrder,
				"Subaccount does not have enough free collateral to place minimum order",
			), nil, nil
		}
		if bigMinOrderQuoteQuantums.BitLen() == 0 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgPlaceOrder,
				"Clob pair is unreasonable and has a minimum quote quantums of 0",
			), nil, nil
		}

		// Prevent integer overflow for unreasonable clob pairs.
		if bigMinOrderQuoteQuantums.Cmp(maxNonOverflowOrderQuoteQuantums) == 1 {
			return simtypes.NoOpMsg(
				types.ModuleName,
				typeMsgPlaceOrder,
				"Clob pair minimum order quote quantums may cause integer overflow",
			), nil, nil
		}

		proposer, _ := simtypes.FindAccount(accs, subaccountId.MustGetAccAddress())

		msg := generateValidPlaceOrder(
			r,
			ctx,
			clobPair,
			proposer,
			subaccountId,
			currentPositionSizeQuantums,
			uint64(clobPair.SubticksPerTick),
			clobPair.StepBaseQuantums,
			bigMinOrderQuoteQuantums,
			bigSubaccountMaxOrderQuoteQuantums,
		)

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
			typeMsgPlaceOrder,
			true, // MsgPlaceOrder should be zero fees
		)

		if err != nil {
			switch {
			case errors.Is(err, satypes.ErrIntegerOverflow),
				errors.Is(err, types.ErrPostOnlyWouldCrossMakerOrder),
				errors.Is(err, types.ErrWouldViolateIsolatedSubaccountConstraints):
				// These errors are expected, and can occur during normal operation. We shouldn't panic on them.
			default:
				panic(err)
			}
		}

		return opMsg, nil, nil
	}
}

// Generates a valid short-term place order with randomized order parameters. The generated place
// order is not guaranteed to be accepted by the memclob.
// Please note that all integer parameters passed into this function should fit inside a uint32
// to prevent any potential integer overflow in the order, fill, and position updates.
func generateValidPlaceOrder(
	r *rand.Rand,
	ctx sdk.Context,
	clobPair types.ClobPair,
	proposer simtypes.Account,
	subaccountId satypes.SubaccountId,
	currentPositionSizeQuantums *big.Int,
	minOrderSubticks uint64,
	minOrderQuantums uint64,
	bigMinOrderQuoteQuantums *big.Int,
	bigSubaccountMaxOrderQuoteQuantums *big.Int,
) *types.MsgPlaceOrder {
	// Generate a random, valid clientId.
	clientId := uint32(simtypes.RandIntBetween(r, 0, clientIdMaxInt))

	// Generate a random, valid GoodTilBlock.
	goodTilBlock := uint32(simtypes.RandIntBetween(
		r,
		int(ctx.BlockHeight()),
		int(lib.MustConvertIntegerToUint32(ctx.BlockHeight())+types.ShortBlockWindow)+1,
	))

	// Generate a random, valid reduceOnly value.
	reduceOnly := sim_helpers.RandWithWeight(r, weightedReduceOnly)

	// Generate a random Order_Side.
	orderSide := sim_helpers.RandSliceShuffle(r, supportedOrderSides)[0]

	// Generate a random, valid TimeInForce.
	timeInForce := sim_helpers.RandWithWeight(r, weightedSupportedTimeInForces)

	// Determine maxOrderQuoteQuantums based on if the order should be collateralized properly.
	maxOrderQuoteQuantums := bigSubaccountMaxOrderQuoteQuantums
	collatIndex := simtypes.RandIntBetween(r, 0, 100)
	if collatIndex >= orderProperlyCollateralizedWeight {
		// Order can exceed collateralization.
		maxOrderQuoteQuantums = maxNonOverflowOrderQuoteQuantums
	}

	// Generate random, valid subticks and quantums.
	// Subticks and quantums will be multiples of the minimums of each respective value.
	// The bounds are determined by:
	// (subtickMultiple * quantumMultiple) <= (maxOrderQuoteQuantums / minOrderQuoteQuantums)
	maxOrderMultiple := new(big.Int).Div(maxOrderQuoteQuantums, bigMinOrderQuoteQuantums)

	// Subticks are bounded by both maxOrderMultiple and orderLevels.
	bigMaxSubtickMultiple := lib.BigMin(maxOrderMultiple, orderLevels)
	maxSubtickMultiple := int(bigMaxSubtickMultiple.Uint64())

	// Generate random, valid subticks.
	subtickMultiple := uint64(simtypes.RandIntBetween(r, 1, maxSubtickMultiple+1))
	bigSubtickMultiple := new(big.Int).SetUint64(subtickMultiple)
	subticks := minOrderSubticks * subtickMultiple

	// Quantums are bounded by the new maxOrderMultiple.
	maxQuantumsMultiple := new(big.Int).Div(
		maxOrderQuoteQuantums,
		new(big.Int).Mul(bigMinOrderQuoteQuantums, bigSubtickMultiple),
	)

	// Generate random, valid quantums.
	quantumsMultiple := uint64(simtypes.RandIntBetween(r, 1, int(maxQuantumsMultiple.Int64())+1))
	bigQuantums := new(big.Int).Mul(
		big.NewInt(int64(minOrderQuantums)),
		big.NewInt(int64(quantumsMultiple)),
	)
	// Default to minOrderQuantums if theres a potential overflow.
	quantums := minOrderQuantums
	if bigQuantums.IsUint64() {
		quantums = minOrderQuantums * quantumsMultiple
	}

	// Handle special order conditions.
	if reduceOnly {
		// Reduce only must be opposite of current positions in clob pair.
		curPositionSign := currentPositionSizeQuantums.Sign()
		if curPositionSign < 0 {
			// currently short, order should go long
			orderSide = types.Order_SIDE_BUY
		} else if curPositionSign == 0 {
			// no current position, cannot place a reduce-only order
			reduceOnly = false
		} else {
			// currently long, order should go short
			orderSide = types.Order_SIDE_SELL
		}
	}

	return &types.MsgPlaceOrder{
		Order: types.Order{
			OrderId: types.OrderId{
				SubaccountId: subaccountId,
				ClientId:     clientId,
				ClobPairId:   clobPair.Id,
			},
			Side:     orderSide,
			Quantums: quantums,
			Subticks: subticks,
			GoodTilOneof: &types.Order_GoodTilBlock{
				GoodTilBlock: goodTilBlock,
			},
			TimeInForce: timeInForce,
			ReduceOnly:  reduceOnly,
		},
	}
}

// Gets the max order quote quantums that would allow the Subaccount stay collateralized and prevent
// integer overflow.
func getMaxSubaccountOrderQuoteQuantums(
	ctx sdk.Context,
	sk types.SubaccountsKeeper,
	subaccountId satypes.SubaccountId,
) *big.Int {
	risk, err := sk.GetNetCollateralAndMarginRequirements(
		ctx,
		satypes.Update{
			SubaccountId: subaccountId,
		},
	)
	if err != nil {
		panic(err)
	}

	maxQuoteQuantums := new(big.Int).Sub(risk.NC, risk.IMR)
	return lib.BigMin(maxQuoteQuantums, maxNonOverflowOrderQuoteQuantums)
}
