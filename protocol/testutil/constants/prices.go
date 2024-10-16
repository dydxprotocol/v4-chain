package constants

import (
	pricefeedclient "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func init() {
	_ = TestTxBuilder.SetMsgs(EmptyMsgUpdateMarketPrices)
	EmptyMsgUpdateMarketPricesTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(ValidMsgUpdateMarketPrices)
	ValidMsgUpdateMarketPricesTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidMsgUpdateMarketPricesStateless)
	InvalidMsgUpdateMarketPricesStatelessTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidMsgUpdateMarketPricesStateful)
	InvalidMsgUpdateMarketPricesStatefulTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())
}

const (
	BtcUsdPair  = "BTC-USD"
	EthUsdPair  = "ETH-USD"
	PolUsdPair  = "POL-USD"
	SolUsdPair  = "SOL-USD"
	LtcUsdPair  = "LTC-USD"
	IsoUsdPair  = "ISO-USD"
	Iso2UsdPair = "ISO2-USD"

	BtcUsdExponent  = -5
	EthUsdExponent  = -6
	LinkUsdExponent = -8
	PolUsdExponent  = -9
	CrvUsdExponent  = -10
	SolUsdExponent  = -8
	LtcUsdExponent  = -7
	IsoUsdExponent  = -8
	Iso2UsdExponent = -7

	CoinbaseExchangeName  = "Coinbase"
	BinanceExchangeName   = "Binance"
	BinanceUSExchangeName = "BinanceUS"
	BitfinexExchangeName  = "Bitfinex"
	KrakenExchangeName    = "Kraken"

	FiveBillion  = uint64(5_000_000_000)
	ThreeBillion = uint64(3_000_000_000)
	FiveMillion  = uint64(5_000_000)
	OneMillion   = uint64(1_000_000)

	// Market param validation errors.
	ErrorMsgMarketPairCannotBeEmpty = "Pair cannot be empty"
	ErrorMsgInvalidMinPriceChange   = "Min price change in parts-per-million must be greater than 0 and less than 10000"
)

var TestMarketExchangeConfigs = map[pricefeedclient.MarketId]string{
	exchange_config.MARKET_BTC_USD: `{
		"exchanges": [
		  {
			"exchangeName": "Binance",
			"ticker": "BTCUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "BinanceUS",
			"ticker": "BTCUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Bitfinex",
			"ticker": "tBTCUSD"
		  },
		  {
			"exchangeName": "Bitstamp",
			"ticker": "BTC/USD"
		  },
		  {
			"exchangeName": "Bybit",
			"ticker": "BTCUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "CoinbasePro",
			"ticker": "BTC-USD"
		  },
		  {
			"exchangeName": "CryptoCom",
			"ticker": "BTC_USD"
		  },
		  {
			"exchangeName": "Kraken",
			"ticker": "XXBTZUSD"
		  },
		  {
			"exchangeName": "Mexc",
			"ticker": "BTC_USDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Okx",
			"ticker": "BTC-USDT",
			"adjustByMarket": "USDT-USD"
		  }
		]
	  }`,
	exchange_config.MARKET_ETH_USD: `{
		"exchanges": [
		  {
			"exchangeName": "Binance",
			"ticker": "ETHUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "BinanceUS",
			"ticker": "ETHUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Bitfinex",
			"ticker": "tETHUSD"
		  },
		  {
			"exchangeName": "Bitstamp",
			"ticker": "ETH/USD"
		  },
		  {
			"exchangeName": "Bybit",
			"ticker": "ETHUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "CoinbasePro",
			"ticker": "ETH-USD"
		  },
		  {
			"exchangeName": "CryptoCom",
			"ticker": "ETH_USD"
		  },
		  {
			"exchangeName": "Kraken",
			"ticker": "XETHZUSD"
		  },
		  {
			"exchangeName": "Mexc",
			"ticker": "ETH_USDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Okx",
			"ticker": "ETH-USDT",
			"adjustByMarket": "USDT-USD"
		  }
		]
	  }`,
	exchange_config.MARKET_SOL_USD: `{
		"exchanges": [
		  {
			"exchangeName": "Binance",
			"ticker": "SOLUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Bitfinex",
			"ticker": "tSOLUSD",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Bybit",
			"ticker": "SOLUSDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "CoinbasePro",
			"ticker": "SOL-USD"
		  },
		  {
			"exchangeName": "CryptoCom",
			"ticker": "SOL_USD"
		  },
		  {
			"exchangeName": "Huobi",
			"ticker": "solusdt",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Kraken",
			"ticker": "SOLUSD"
		  },
		  {
			"exchangeName": "Kucoin",
			"ticker": "SOL-USDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Mexc",
			"ticker": "SOL_USDT",
			"adjustByMarket": "USDT-USD"
		  },
		  {
			"exchangeName": "Okx",
			"ticker": "SOL-USDT",
			"adjustByMarket": "USDT-USD"
		  }
		]
	  }`,
	exchange_config.MARKET_ISO_USD: `{
		"exchanges": [
		  {
			"exchangeName": "Binance",
			"ticker": "ISOUSDT",
			"adjustByMarket": "USDT-USD"
		  }
		]
	  }`,
	exchange_config.MARKET_ISO2_USD: `{
			"exchanges": [
				{
				"exchangeName": "Binance",
				"ticker": "ISO2USDT",
				"adjustByMarket": "USDT-USD"
				}
			]
			}`,
}

var TestMarketParams = []types.MarketParam{
	{
		Id:                 0,
		Pair:               BtcUsdPair,
		Exponent:           BtcUsdExponent,
		MinExchanges:       1,
		MinPriceChangePpm:  50,
		ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_BTC_USD],
	},
	{
		Id:                 1,
		Pair:               EthUsdPair,
		Exponent:           EthUsdExponent,
		MinExchanges:       1,
		MinPriceChangePpm:  50,
		ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_ETH_USD],
	},
	{
		Id:                 2,
		Pair:               SolUsdPair,
		Exponent:           SolUsdExponent,
		MinExchanges:       1,
		MinPriceChangePpm:  50,
		ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_SOL_USD],
	},
	{
		Id:                 3,
		Pair:               IsoUsdPair,
		Exponent:           IsoUsdExponent,
		MinExchanges:       1,
		MinPriceChangePpm:  50,
		ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_ISO_USD],
	},
	{
		Id:                 4,
		Pair:               Iso2UsdPair,
		Exponent:           Iso2UsdExponent,
		MinExchanges:       1,
		MinPriceChangePpm:  50,
		ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_ISO2_USD],
	},
}

var TestMarketPrices = []types.MarketPrice{
	{
		Id:       0,
		Exponent: BtcUsdExponent,
		Price:    FiveBillion, // $50,000 == 1 BTC
	},
	{
		Id:       1,
		Exponent: EthUsdExponent,
		Price:    ThreeBillion, // $3,000 == 1 ETH
	},
	{
		Id:       2,
		Exponent: SolUsdExponent,
		Price:    FiveBillion, // 50$ == 1 SOL
	},
	{
		Id:       3,
		Exponent: IsoUsdExponent,
		Price:    FiveBillion, // 50$ == 1 ISO
	},
	{
		Id:       4,
		Exponent: Iso2UsdExponent,
		Price:    ThreeBillion, // 300$ == 1 ISO2
	},
}

var TestMarketIdsToExponents = map[uint32]int32{
	0: BtcUsdExponent,
	1: EthUsdExponent,
	2: SolUsdExponent,
	3: IsoUsdExponent,
	4: Iso2UsdExponent,
}

var TestPricesGenesisState = types.GenesisState{
	MarketParams: TestMarketParams,
	MarketPrices: TestMarketPrices,
}

var (
	ValidMarketPriceUpdates = []*types.MsgUpdateMarketPrices_MarketPrice{
		types.NewMarketPriceUpdate(MarketId0, Price5),
		types.NewMarketPriceUpdate(MarketId1, Price6),
		types.NewMarketPriceUpdate(MarketId2, Price7),
		types.NewMarketPriceUpdate(MarketId3, Price4),
		types.NewMarketPriceUpdate(MarketId4, Price3),
	}

	// `MsgUpdateMarketPrices`.
	EmptyMsgUpdateMarketPrices        = &types.MsgUpdateMarketPrices{}
	EmptyMsgUpdateMarketPricesTxBytes []byte

	ValidMsgUpdateMarketPrices = &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: ValidMarketPriceUpdates,
	}
	ValidMsgUpdateMarketPricesTxBytes []byte

	InvalidMsgUpdateMarketPricesStateless = &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
			types.NewMarketPriceUpdate(MarketId0, 0), // 0 price value is invalid.
		},
	}
	InvalidMsgUpdateMarketPricesStatelessTxBytes []byte

	InvalidMsgUpdateMarketPricesStateful = &types.MsgUpdateMarketPrices{
		MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
			types.NewMarketPriceUpdate(MarketId0, Price5),
			types.NewMarketPriceUpdate(MarketId1, Price6),
			types.NewMarketPriceUpdate(99, Price3), // Market with id 99 does not exist.
		},
	}
	InvalidMsgUpdateMarketPricesStatefulTxBytes []byte

	Prices_DefaultGenesisState = types.GenesisState{
		MarketParams: []types.MarketParam{
			{
				Id:                 uint32(0),
				Pair:               BtcUsdPair,
				Exponent:           BtcUsdExponent,
				MinExchanges:       uint32(2),
				ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_BTC_USD],
				MinPriceChangePpm:  uint32(50),
			},
			{
				Id:                 uint32(1),
				Pair:               EthUsdPair,
				Exponent:           EthUsdExponent,
				MinExchanges:       uint32(1),
				ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_ETH_USD],
				MinPriceChangePpm:  uint32(50),
			},
		},
		MarketPrices: []types.MarketPrice{
			{
				Id:       uint32(0),
				Exponent: BtcUsdExponent,
				Price:    FiveBillion, // $50,000 == 1 BTC
			},
			{
				Id:       uint32(1),
				Exponent: EthUsdExponent,
				Price:    ThreeBillion, // $3,000 == 1 ETH
			},
		},
	}

	Prices_MultiExchangeMarketGenesisState = types.GenesisState{
		MarketParams: []types.MarketParam{
			{ // BTC-USD
				Id:                 uint32(0),
				Pair:               BtcUsdPair,
				Exponent:           BtcUsdExponent,
				ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_BTC_USD],
				MinExchanges:       uint32(2),
				MinPriceChangePpm:  uint32(50),
			},
			{ // ETH-USD
				Id:                 uint32(1),
				Pair:               EthUsdPair,
				Exponent:           EthUsdExponent,
				ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_ETH_USD],
				MinExchanges:       uint32(2),
				MinPriceChangePpm:  uint32(50),
			},
		},
		MarketPrices: []types.MarketPrice{
			{ // BTC-USD
				Id:       uint32(0),
				Exponent: BtcUsdExponent,
				Price:    FiveBillion, // $50,000 == 1 BTC
			},
			{ // ETH-USD
				Id:       uint32(1),
				Exponent: EthUsdExponent,
				Price:    ThreeBillion, // $3,000 == 1 ETH
			},
		},
	}
)
