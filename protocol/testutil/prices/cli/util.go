package cli

import (
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

func NetworkWithMarketObjects(t *testing.T, n int) (*network.Network, []types.MarketParam, []types.MarketPrice) {
	t.Helper()
	state := types.GenesisState{}

	state.MarketParams = []types.MarketParam{}
	state.MarketPrices = []types.MarketPrice{}

	// Market params
	for i := 0; i < n; i++ {
		marketParam := types.MarketParam{
			Id:                 uint32(i),
			Pair:               fmt.Sprint(constants.BtcUsdPair, i),
			MinExchanges:       uint32(1),
			MinPriceChangePpm:  uint32((i + 1) * 2),
			ExchangeConfigJson: "{}",
		}
		state.MarketParams = append(state.MarketParams, marketParam)
	}

	// Market prices
	for i := 0; i < n; i++ {
		marketPrice := types.MarketPrice{
			Id:    uint32(i),
			Price: constants.FiveBillion,
		}
		state.MarketPrices = append(state.MarketPrices, marketPrice)
	}

	genesis := getFullGenesisForMarketObjects(n)
	network.DeployCustomNetwork(genesis)
	return nil, state.MarketParams, state.MarketPrices
}

func getFullGenesisForMarketObjects(n int) string {
	fullGenesisTwo := "\".app_state.prices.market_prices = [{\\\"exponent\\\": \\\"0\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"0\\\", \\\"price\\\": \\\"5000000000\\\"}] | .app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD0\\\", \\\"exponent\\\": \\\"0\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"2\\\", \\\"exchange_config_json\\\": \\\"{}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"BTC-USD1\\\", \\\"exponent\\\": \\\"0\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"4\\\", \\\"exchange_config_json\\\": \\\"{}\\\"}]\" \"\""
	fullGenesisFive := "\".app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD0\\\", \\\"exponent\\\": \\\"0\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"2\\\", \\\"exchange_config_json\\\": \\\"{}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"BTC-USD1\\\", \\\"exponent\\\": \\\"0\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"4\\\", \\\"exchange_config_json\\\": \\\"{}\\\"}, {\\\"id\\\": \\\"2\\\", \\\"pair\\\": \\\"BTC-USD2\\\", \\\"exponent\\\": \\\"0\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"6\\\", \\\"exchange_config_json\\\": \\\"{}\\\"}, {\\\"id\\\": \\\"3\\\", \\\"pair\\\": \\\"BTC-USD3\\\", \\\"exponent\\\": \\\"0\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"8\\\", \\\"exchange_config_json\\\": \\\"{}\\\"}, {\\\"id\\\": \\\"4\\\", \\\"pair\\\": \\\"BTC-USD4\\\", \\\"exponent\\\": \\\"0\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"10\\\", \\\"exchange_config_json\\\": \\\"{}\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"0\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"0\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"exponent\\\": \\\"0\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"3\\\", \\\"exponent\\\": \\\"0\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"4\\\", \\\"exponent\\\": \\\"0\\\", \\\"price\\\": \\\"5000000000\\\"}]\" \"\""

	var genesis string
	if n == 2 {
		genesis = fullGenesisTwo
	} else if n == 5 {
		genesis = fullGenesisFive
	}
	return genesis
}
