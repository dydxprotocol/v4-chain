package simulation

// DONTCOVER

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// genNumMarkets returns randomized num markets.
func genNumMarkets(r *rand.Rand, isReasonableGenesis bool) int {
	return simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinMarkets, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxMarkets, isReasonableGenesis)+1,
	)
}

// genMinExchanges returns randomized min exchanges per market.
func genMinExchanges(r *rand.Rand, isReasonableGenesis bool) int {
	return simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinMinExchangesPerMarket, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxMinExchangesPerMarket, isReasonableGenesis)+1,
	)
}

// genMarketExponent returns randomized market exponent.
func genMarketExponent(r *rand.Rand, isReasonableGenesis bool) int {
	return simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinMarketExponent, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxMarketExponent, isReasonableGenesis)+1,
	)
}

// genMarketName return randomized market name.
func genMarketName(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 3, 6)) + "-USD"
}

// genMarketPrice returns randomized market price.
func genMarketPrice(r *rand.Rand, isReasonableGenesis bool) uint64 {
	return uint64(simtypes.RandIntBetween(r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinMarketPrice, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxMarketPrice, isReasonableGenesis),
	))
}

// RandomizedGenState generates a random GenesisState for `Prices`.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	isReasonableGenesis := sim_helpers.ShouldGenerateReasonableGenesis(r, simState.GenTimestamp)

	minExchangesPerMarket := sim_helpers.MinExchangesPerMarket
	numMarkets := genNumMarkets(r, isReasonableGenesis)

	marketParams := make([]types.MarketParam, numMarkets)
	marketPrices := make([]types.MarketPrice, numMarkets)

	for i := 0; i < numMarkets; i++ {
		var minExchanges = genMinExchanges(r, isReasonableGenesis)
		if minExchanges < minExchangesPerMarket {
			minExchanges = minExchangesPerMarket
		}

		marketExponent := genMarketExponent(r, isReasonableGenesis)

		marketParams[i] = types.MarketParam{
			Id:                uint32(i),
			Pair:              genMarketName(r),
			Exponent:          int32(marketExponent),
			MinExchanges:      uint32(minExchanges),
			MinPriceChangePpm: uint32(simtypes.RandIntBetween(r, 1, int(lib.MaxPriceChangePpm))),
			// The simulation tests don't run the daemon currently so we pass in empty exchange config.
			ExchangeConfigJson: "{}",
		}
		marketPrices[i] = types.MarketPrice{
			Id:       uint32(i),
			Exponent: int32(marketExponent),
			Price:    genMarketPrice(r, isReasonableGenesis),
		}
	}

	pricesGenesis := types.GenesisState{
		MarketParams: marketParams,
		MarketPrices: marketPrices,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pricesGenesis)
}
