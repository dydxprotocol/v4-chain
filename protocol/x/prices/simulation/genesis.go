package simulation

// DONTCOVER

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/sim_helpers"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// genNumMarkets returns randomized num markets.
func genNumMarkets(r *rand.Rand, isReasonableGenesis bool) int {
	return simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinMarkets, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxMarkets, isReasonableGenesis)+1,
	)
}

// genNumExchangeFeeds returns randomized num exchange feeds.
func genNumExchangeFeeds(r *rand.Rand, isReasonableGenesis bool) int {
	return simtypes.RandIntBetween(
		r,
		sim_helpers.PickGenesisParameter(sim_helpers.MinExchangeFeeds, isReasonableGenesis),
		sim_helpers.PickGenesisParameter(sim_helpers.MaxExchangeFeeds, isReasonableGenesis)+1,
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

// genMarketName return randomized exchange feed name.
func genExchangeFeedName(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 5, 20))
}

// genExchangeFeedMemo return randomized exchange feed memo.
func genExchangeFeedMemo(r *rand.Rand) string {
	return simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 5, 20))
}

// RandomizedGenState generates a random GenesisState for `Prices`.
func RandomizedGenState(simState *module.SimulationState) {
	r := simState.Rand
	isReasonableGenesis := sim_helpers.ShouldGenerateReasonableGenesis(r, simState.GenTimestamp)

	minExchangeFeedsPerMarket := sim_helpers.MinExchangeFeedsPerMarket
	numMarkets := genNumMarkets(r, isReasonableGenesis)
	numExchangeFeeds := genNumExchangeFeeds(r, isReasonableGenesis)

	markets := make([]types.Market, numMarkets)
	exchangeFeeds := make([]types.ExchangeFeed, numExchangeFeeds)
	allExchangeIds := make([]uint32, numExchangeFeeds)
	for i := 0; i < numExchangeFeeds; i++ {
		allExchangeIds[i] = uint32(i)
		exchangeFeeds[i] = types.ExchangeFeed{
			Id:   uint32(i),
			Name: genExchangeFeedName(r),
			Memo: genExchangeFeedMemo(r),
		}
	}

	for i := 0; i < numMarkets; i++ {
		randomizedExchangeIds := sim_helpers.RandSliceShuffle(r, allExchangeIds)

		var numExchangesForMarket int
		// RandIntBetween panics if arguments are equal for some ungodly reason.
		if minExchangeFeedsPerMarket == numExchangeFeeds {
			numExchangesForMarket = minExchangeFeedsPerMarket
		} else {
			numExchangesForMarket = simtypes.RandIntBetween(r, minExchangeFeedsPerMarket, numExchangeFeeds)
		}

		var minExchanges int
		// RandIntBetween panics if arguments are equal for some ungodly reason.
		if minExchangeFeedsPerMarket == numExchangesForMarket {
			minExchanges = minExchangeFeedsPerMarket
		} else {
			minExchanges = simtypes.RandIntBetween(r, minExchangeFeedsPerMarket, numExchangesForMarket)
		}

		marketExchangeIds := randomizedExchangeIds[:numExchangesForMarket]
		marketExponent := genMarketExponent(r, isReasonableGenesis)

		markets[i] = types.Market{
			Pair:              genMarketName(r),
			Exponent:          int32(marketExponent),
			Exchanges:         marketExchangeIds,
			MinExchanges:      uint32(minExchanges),
			MinPriceChangePpm: uint32(simtypes.RandIntBetween(r, 1, int(lib.MaxPriceChangePpm))),
			Price:             r.Uint64(),
		}
	}

	pricesGenesis := types.GenesisState{
		ExchangeFeeds: exchangeFeeds,
		Markets:       markets,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pricesGenesis)
}
