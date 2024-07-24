package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/dydx"
	dydxtypes "github.com/skip-mev/slinky/providers/apis/dydx/types"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	"go.uber.org/zap"
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
	for existingMarketNames[marketName] {
		marketName = simtypes.RandStringOfLength(r, simtypes.RandIntBetween(r, 3, 6)) + "-USD"
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
	marketMap := ConstructMarketMapFromParams(marketParams)
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

func ConstructMarketMapFromParams(
	allMarketParams []types.MarketParam,
) marketmaptypes.MarketMap {
	// fill out config with dummy variables to pass validation.  This handler is only used to run the
	// ConvertMarketParamsToMarketMap member function.
	h, err := dydx.NewAPIHandler(zap.NewNop(), config.APIConfig{
		Enabled:          true,
		Timeout:          1,
		Interval:         1,
		ReconnectTimeout: 1,
		MaxQueries:       1,
		Atomic:           false,
		Endpoints:        []config.Endpoint{{URL: "upgrade"}},
		BatchSize:        0,
		Name:             dydx.Name,
	})
	if err != nil {
		panic(err) // panic to halt/fail simulation.
	}

	var mpr dydxtypes.QueryAllMarketParamsResponse
	for _, mp := range allMarketParams {
		mpr.MarketParams = append(mpr.MarketParams, dydxtypes.MarketParam{
			Id:                 mp.Id,
			Pair:               mp.Pair,
			Exponent:           mp.Exponent,
			MinExchanges:       mp.MinExchanges,
			MinPriceChangePpm:  mp.MinPriceChangePpm,
			ExchangeConfigJson: mp.ExchangeConfigJson,
		})
	}
	mm, err := h.ConvertMarketParamsToMarketMap(mpr)
	if err != nil {
		panic(err) // panic to halt/fail simulation.
	}

	return mm.MarketMap
}
