package marketmap_test

import (
	"testing"

	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"
	marketmaptypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/marketmap"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestConstructMarketMapFromParams(t *testing.T) {
	marketParams := []pricestypes.MarketParam{
		{
			Pair:               "BTC-USD",
			Exponent:           -8,
			MinExchanges:       1,
			MinPriceChangePpm:  10,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Binance","ticker":"BTCUSDT"}]}`,
		},
	}

	expectedMarketMap := marketmaptypes.MarketMap{
		Markets: map[string]marketmaptypes.Market{
			"BTC/USD": {
				Ticker: marketmaptypes.Ticker{
					CurrencyPair:     slinkytypes.CurrencyPair{Base: "BTC", Quote: "USD"},
					Decimals:         8,
					MinProviderCount: 1,
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmaptypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "BTCUSDT"},
				},
			},
		},
	}
	marketMap, err := marketmap.ConstructMarketMapFromParams(marketParams)
	require.NoError(t, err)
	require.NotNil(t, marketMap)
	require.Equal(t, expectedMarketMap, marketMap)
}
