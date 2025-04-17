package cli

import (
	"fmt"
	"testing"

	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/marketmap"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func NetworkWithMarketObjects(t testing.TB, n int) (*network.Network, []types.MarketParam, []types.MarketPrice) {
	t.Helper()
	cfg := network.DefaultConfig(nil)
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	// Overwrite market params and prices in default genesis state.
	state.MarketParams = []types.MarketParam{}
	state.MarketPrices = []types.MarketPrice{}

	// Market params
	for i := 0; i < n; i++ {
		marketName := fmt.Sprint(constants.BtcUsdPair, i)
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

		marketParam := types.MarketParam{
			Id:                uint32(i),
			Pair:              marketName,
			Exponent:          int32(-5),
			MinExchanges:      uint32(1),
			MinPriceChangePpm: uint32((i + 1) * 2),
			// x/marketmap expects at least as many valid exchange names defined as the value of MinExchanges.
			ExchangeConfigJson: exchangeJson,
		}
		state.MarketParams = append(state.MarketParams, marketParam)
	}

	// Market prices
	for i := 0; i < n; i++ {
		marketPrice := types.MarketPrice{
			Id:       uint32(i),
			Exponent: int32(-5),
			Price:    constants.FiveBillion,
		}
		state.MarketPrices = append(state.MarketPrices, marketPrice)
	}

	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf

	// Inject marketmap genesis
	marketMap, err := marketmap.ConstructMarketMapFromParams(state.MarketParams)
	require.NoError(t, err)
	marketmapGenesis := marketmaptypes.GenesisState{
		MarketMap: marketMap,
		Params:    marketmaptypes.DefaultParams(),
	}
	cfg.GenesisState[marketmaptypes.ModuleName] = cfg.Codec.MustMarshalJSON(&marketmapGenesis)

	return network.New(t, cfg), state.MarketParams, state.MarketPrices
}
