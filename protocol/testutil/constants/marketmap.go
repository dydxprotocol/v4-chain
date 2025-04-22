package constants

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"
	marketmapmoduletypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
)

var GovAuthority = authtypes.NewModuleAddress(govtypes.ModuleName).String()

var (
	BtcUsdCurrencyPair, _ = slinky.MarketPairToCurrencyPair(BtcUsdPair)
	EthUsdCurrencyPair, _ = slinky.MarketPairToCurrencyPair(EthUsdPair)
	UsdtUsdCurrencyPair   = slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}
)

var MarketMap_DefaultGenesisState = marketmapmoduletypes.GenesisState{
	MarketMap: marketmapmoduletypes.MarketMap{
		Markets: map[string]marketmapmoduletypes.Market{
			BtcUsdCurrencyPair.String(): {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     BtcUsdCurrencyPair,
					Decimals:         uint64(5),
					MinProviderCount: uint64(2),
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "BTC-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "btcusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XXBTZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "BTCUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			EthUsdCurrencyPair.String(): {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     EthUsdCurrencyPair,
					Decimals:         uint64(6),
					MinProviderCount: uint64(1),
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "ETH-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "ethusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "XETHZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "ETH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "mexc_ws", OffChainTicker: "ETHUSDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "ETH-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "USDT", Quote: "USD"}, Invert: false, Metadata_JSON: ""},
				},
			},
			UsdtUsdCurrencyPair.String(): {
				Ticker: marketmapmoduletypes.Ticker{
					CurrencyPair:     UsdtUsdCurrencyPair,
					Decimals:         uint64(9),
					MinProviderCount: uint64(3),
					Enabled:          true,
					Metadata_JSON:    "",
				},
				ProviderConfigs: []marketmapmoduletypes.ProviderConfig{
					{Name: "binance_ws", OffChainTicker: "USDCUSDT",
						NormalizeByPair: nil, Invert: true, Metadata_JSON: ""},
					{Name: "bybit_ws", OffChainTicker: "USDCUSDT",
						NormalizeByPair: nil, Invert: true, Metadata_JSON: ""},
					{Name: "coinbase_ws", OffChainTicker: "USDT-USD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "huobi_ws", OffChainTicker: "ethusdt",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "ETH", Quote: "USD"}, Invert: true, Metadata_JSON: ""},
					{Name: "kraken_api", OffChainTicker: "USDTZUSD",
						NormalizeByPair: nil, Invert: false, Metadata_JSON: ""},
					{Name: "kucoin_ws", OffChainTicker: "BTC-USDT",
						NormalizeByPair: &slinkytypes.CurrencyPair{Base: "BTC", Quote: "USD"}, Invert: true, Metadata_JSON: ""},
					{Name: "okx_ws", OffChainTicker: "USDC-USDT",
						NormalizeByPair: nil, Invert: true, Metadata_JSON: ""},
				},
			},
		},
	},
	Params: marketmapmoduletypes.Params{
		MarketAuthorities: []string{GovAuthority},
		Admin:             GovAuthority,
	},
}
