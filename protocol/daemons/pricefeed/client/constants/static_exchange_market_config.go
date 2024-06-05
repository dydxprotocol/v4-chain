package constants

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
)

const (
	// MinimumRequiredExchangesPerMarket is the minimum number of markets required for a market to be reliably priced
	// by the pricefeed daemon. This number was chosen to supply the minimum number of prices required to
	// compute an index price for a market, given exchange unavailability due to exchange geo-fencing,
	// downtime, etc.
	MinimumRequiredExchangesPerMarket = 5
)

// GenerateExchangeConfigJson generates human-readable exchange config json for each market based on the contents
// of an exchangeToExchangeConfig map. For the default exchange configs, pass in the
// `StaticExchangeMarketConfig` map above as the argument.
func GenerateExchangeConfigJson(
	exchangeToExchangeConfig map[types.ExchangeId]*types.MutableExchangeMarketConfig,
) (
	marketToExchangeConfigJson map[types.MarketId]string,
) {
	// marketToExchangeConfigJson maps markets to a map of exchange, exchange market config. This
	// is used to generate the exchange config json from the above map that is keyed by exchange id.
	// We keep this intermediate map so that we can sort the exchange configs for each market by exchange name
	// in order to make the output deterministic.
	marketToExchangeMarketConfigs := make(map[types.MarketId]map[string]types.ExchangeMarketConfigJson)

	// Generate the market-specific exchange config for each market, exchange.
	for id, exchangeConfig := range exchangeToExchangeConfig {
		// Skip config for the test exchange.
		if id == exchange_common.EXCHANGE_ID_TEST_EXCHANGE {
			continue
		}
		if id == exchange_common.EXCHANGE_ID_TEST_FIXED_PRICE_EXCHANGE {
			continue
		}
		for marketId, config := range exchangeConfig.MarketToMarketConfig {
			marketExchangeConfigs, ok := marketToExchangeMarketConfigs[marketId]
			if !ok {
				marketToExchangeMarketConfigs[marketId] = map[string]types.ExchangeMarketConfigJson{}
				marketExchangeConfigs = marketToExchangeMarketConfigs[marketId]
			}

			exchangeMarketConfigJson := types.ExchangeMarketConfigJson{
				ExchangeName: id,
				Ticker:       config.Ticker,
				Invert:       config.Invert,
			}

			// Convert adjust-by market id to name if specified.
			if config.AdjustByMarket != nil {
				adjustByMarketName, ok := exchange_config.StaticMarketNames[*config.AdjustByMarket]
				if !ok {
					panic(fmt.Sprintf("no name for adjust-by market %v", *config.AdjustByMarket))
				}
				exchangeMarketConfigJson.AdjustByMarket = adjustByMarketName
			}

			marketExchangeConfigs[id] = exchangeMarketConfigJson
		}
	}

	// Initialize the output map.
	marketToExchangeConfigJson = make(map[types.MarketId]string, len(marketToExchangeMarketConfigs))

	// Generate the output map of market to exchange config json.
	for marketId, exchangeToConfigs := range marketToExchangeMarketConfigs {
		// Sort output exchange configs by exchange name in order to make output deterministic.
		exchangeNames := make([]string, 0, len(exchangeToConfigs))

		// 1. Generate sorted list of exchange names.
		for name := range exchangeToConfigs {
			exchangeNames = append(exchangeNames, name)
		}
		sort.Strings(exchangeNames)

		// 2. Generate sorted list of exchange configs by exchange name.
		sortedExchangeConfigs := make([]types.ExchangeMarketConfigJson, 0, len(exchangeNames))
		for _, exchangeName := range exchangeNames {
			sortedExchangeConfigs = append(sortedExchangeConfigs, exchangeToConfigs[exchangeName])
		}
		exchangeConfigJson := types.ExchangeConfigJson{
			Exchanges: sortedExchangeConfigs,
		}

		// 3. Generate human-readable formatted output json for the market, sorted by exchange name.
		bytes, err := json.MarshalIndent(exchangeConfigJson, "", "  ")
		if err != nil {
			panic(err)
		}
		marketToExchangeConfigJson[marketId] = string(bytes)
	}
	return marketToExchangeConfigJson
}
