package simulation

// DONTCOVER

import (
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"math/big"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// Simulation operation weights constants
const (
	opWeightMsgUpdateMarketPrices = "op_weight_msg_update_market_prices"

	defaultWeightMsgUpdateMarketPrices int = 100
)

var (
	maxPriceChangePpm = int(0.5 * 1_000_000) // 50% in ppm

	typeMsgUpdateMarketPrices = sdk.MsgTypeURL(&types.MsgUpdateMarketPrices{})
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {
	protoCdc := codec.NewProtoCodec(module.InterfaceRegistry)

	operations := make([]simtypes.WeightedOperation, 0)

	// MsgUpdateMarketPrices
	var weightMsgUpdateMarketPrices int
	appParams.GetOrGenerate(opWeightMsgUpdateMarketPrices, &weightMsgUpdateMarketPrices, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateMarketPrices = defaultWeightMsgUpdateMarketPrices
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateMarketPrices,
		SimulateMsgUpdateMarketPrices(protoCdc, k, ak, bk),
	))

	return operations
}

// SimulateMsgUpdateMarketPrices generates a random MsgUpdateMarketPrices.
func SimulateMsgUpdateMarketPrices(
	cdc *codec.ProtoCodec,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simtypes.Operation {
	return func(
		r *rand.Rand,
		app *baseapp.BaseApp,
		ctx sdk.Context,
		accs []simtypes.Account,
		chainId string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		proposer, _ := simtypes.RandomAcc(r, accs)

		allMarketParamPrices, _ := k.GetAllMarketParamPrices(ctx)

		priceUpdates := make([]*types.MsgUpdateMarketPrices_MarketPrice, 0)
		for _, marketParamPrice := range allMarketParamPrices {
			// 50% chance of updating the price.
			if sim_helpers.RandBool(r) {
				newPrice := getRandomlyUpdatedPrice(r, marketParamPrice)

				// only update if the new price is not 0
				if newPrice != 0 {
					priceUpdates = append(
						priceUpdates,
						types.NewMarketPriceUpdate(marketParamPrice.Param.Id, newPrice),
					)
				}
			}
		}

		if len(priceUpdates) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, typeMsgUpdateMarketPrices, "empty: market price updates"), nil, nil
		}

		msg := &types.MsgUpdateMarketPrices{
			MarketPriceUpdates: priceUpdates,
		}

		opMsg, err := sim_helpers.GenerateAndDeliverTx(
			r,
			app,
			ctx,
			chainId,
			cdc,
			ak,
			bk,
			proposer,
			types.ModuleName,
			msg,
			typeMsgUpdateMarketPrices,
			true, // fee does not apply when updating market prices.
		)
		if err != nil {
			panic(err) // panic to halt/fail simulation.
		}

		return opMsg, nil, nil
	}
}

// getRandomlyUpdatedPrice returns a valid, random new market price.
func getRandomlyUpdatedPrice(r *rand.Rand, marketParamPrice types.MarketParamPrice) uint64 {
	randomValidChangePpm := uint32(
		simtypes.RandIntBetween(r, int(marketParamPrice.Param.MinPriceChangePpm), maxPriceChangePpm+1),
	)
	bigPrice := new(big.Int).SetUint64(marketParamPrice.Price.Price)

	bigPriceChange := lib.BigIntMulPpm(bigPrice, randomValidChangePpm)
	// 50% chance that the change is in the negative direction.
	if sim_helpers.RandBool(r) {
		bigPriceChange = new(big.Int).Neg(bigPriceChange)
	}

	newBigPrice := new(big.Int).Add(bigPrice, bigPriceChange)
	return newBigPrice.Uint64()
}
