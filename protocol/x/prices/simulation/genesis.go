package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/marketmap"
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
func genMarketName(r *rand.Rand, existingMarketNames map[string]bool) string {
	marketName := simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 3, 6)) + "-USD"
	marketName = strings.ToUpper(marketName)
	for existingMarketNames[marketName] {
		marketName = simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 3, 6)) + "-USD"
		marketName = strings.ToUpper(marketName)
	}
	return marketName
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

	marketNames := make(map[string]bool)

	for i := 0; i < numMarkets; i++ {
		marketName := genMarketName(r, marketNames)
		marketNames[marketName] = true

		var minExchanges = genMinExchanges(r, isReasonableGenesis)
		if minExchanges < minExchangesPerMarket {
			minExchanges = minExchangesPerMarket
		}

		marketExponent := genMarketExponent(r, isReasonableGenesis)

		exchangeJsonTemplate := `{"exchanges":[
			{"exchangeName":"Binance","ticker":"\"%[1]s\""},
			{"exchangeName":"Bitfinex","ticker":"%[1]s"},
			{"exchangeName":"CoinbasePro","ticker":"%[1]s"},
			{"exchangeName":"Gate","ticker":"%[1]s"},
			{"exchangeName":"Huobi","ticker":"%[1]s"},
			{"exchangeName":"Kraken","ticker":"%[1]s"},
			{"exchangeName":"Okx","ticker":"%[1]s"}
		]}`
		exchangeJson := fmt.Sprintf(exchangeJsonTemplate, marketName)

		marketParams[i] = types.MarketParam{
			Id:                uint32(i),
			Pair:              marketName,
			Exponent:          int32(marketExponent),
			MinExchanges:      uint32(minExchanges),
			MinPriceChangePpm: uint32(simtypes.RandIntBetween(r, 1, int(lib.MaxPriceChangePpm))),
			// x/marketmap expects at least as many valid exchange names defined as the value of MinExchanges.
			ExchangeConfigJson: exchangeJson,
		}
		marketPrices[i] = types.MarketPrice{
			Id:       uint32(i),
			Exponent: int32(marketExponent),
			Price:    genMarketPrice(r, isReasonableGenesis),
		}
	}

	var GovAuthority = authtypes.NewModuleAddress(govtypes.ModuleName).String()
	marketMap, _ := marketmap.ConstructMarketMapFromParams(marketParams)
	marketmapGenesis := marketmaptypes.GenesisState{
		MarketMap: marketMap,
		Params: marketmaptypes.Params{
			MarketAuthorities: []string{GovAuthority},
			Admin:             GovAuthority,
		},
	}
	pricesGenesis := types.GenesisState{
		MarketParams: marketParams,
		MarketPrices: marketPrices,
	}

	simState.GenState[marketmaptypes.ModuleName] = simState.Cdc.MustMarshalJSON(&marketmapGenesis)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&pricesGenesis)
}
