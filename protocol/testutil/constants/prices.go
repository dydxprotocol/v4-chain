package constants

import (
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	pricefeedclient "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

const (
	BtcUsdPair   = "BTC-USD"
	EthUsdPair   = "ETH-USD"
	MaticUsdPair = "MATIC-USD"
	SolUsdPair   = "SOL-USD"
	LtcUsdPair   = "LTC-USD"
	IsoUsdPair   = "ISO-USD"
	Iso2UsdPair  = "ISO2-USD"

	BtcUsdExponent   = -5
	EthUsdExponent   = -6
	LinkUsdExponent  = -8
	MaticUsdExponent = -9
	CrvUsdExponent   = -10
	SolUsdExponent   = -8
	LtcUsdExponent   = -7
	IsoUsdExponent   = -8
	Iso2UsdExponent  = -7

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

var TestSingleMarketParam = types.MarketParam{
	Id:                 0,
	Pair:               BtcUsdPair,
	Exponent:           BtcUsdExponent,
	MinExchanges:       1,
	MinPriceChangePpm:  50,
	ExchangeConfigJson: TestMarketExchangeConfigs[exchange_config.MARKET_BTC_USD],
}

var TestMarketPrices = []types.MarketPrice{
	{
		Id:        0,
		Exponent:  BtcUsdExponent,
		SpotPrice: FiveBillion, // $50,000 == 1 BTC
		PnlPrice:  FiveBillion, // $50,000 == 1 BTC
	},
	{
		Id:        1,
		Exponent:  EthUsdExponent,
		SpotPrice: ThreeBillion, // $3,000 == 1 ETH
		PnlPrice:  ThreeBillion, // $3,000 == 1 ETH
	},
	{
		Id:        2,
		Exponent:  SolUsdExponent,
		SpotPrice: FiveBillion, // 50$ == 1 SOL
		PnlPrice:  FiveBillion, // 50$ == 1 SOL
	},
	{
		Id:        3,
		Exponent:  IsoUsdExponent,
		SpotPrice: FiveBillion, // 50$ == 1 ISO
		PnlPrice:  FiveBillion, // 50$ == 1 ISO
	},
	{
		Id:        4,
		Exponent:  Iso2UsdExponent,
		SpotPrice: ThreeBillion, // 300$ == 1 ISO2
		PnlPrice:  ThreeBillion, // 300$ == 1 ISO2
	},
}
var IdsToPairs = map[uint32]string{
	0: BtcUsdPair,
	1: EthUsdPair,
	2: SolUsdPair,
	3: IsoUsdPair,
	4: Iso2UsdPair,
}

var TestMarketIdsToExponents = map[uint32]int32{
	0: BtcUsdExponent,
	1: EthUsdExponent,
	2: SolUsdExponent,
}

var TestPricesGenesisState = types.GenesisState{
	MarketParams: TestMarketParams,
	MarketPrices: TestMarketPrices,
}

var (
	ValidMultiMarketSpotPriceUpdates = []*types.MarketSpotPriceUpdate{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5,
		},
		{
			MarketId:  MarketId1,
			SpotPrice: Price6,
		},
		{
			MarketId:  MarketId2,
			SpotPrice: Price7,
		},
		{
			MarketId:  MarketId3,
			SpotPrice: Price4,
		},
		{
			MarketId:  MarketId4,
			SpotPrice: Price3,
		},
	}
	ValidMarketPriceUpdates = []*types.MarketPriceUpdate{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5,
			PnlPrice:  Price5,
		},
		{
			MarketId:  MarketId1,
			SpotPrice: Price6,
			PnlPrice:  Price6,
		},
		{
			MarketId:  MarketId2,
			SpotPrice: Price7,
			PnlPrice:  Price7,
		},
		{
			MarketId:  MarketId3,
			SpotPrice: Price4,
			PnlPrice:  Price4,
		},
		{
			MarketId:  MarketId4,
			SpotPrice: Price3,
			PnlPrice:  Price3,
		},
	}

	ValidSingleSpotMarketPriceUpdate = []*types.MarketSpotPriceUpdate{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5,
		},
	}

	ValidSingleMarketPriceUpdate = []*types.MarketPriceUpdate{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5,
			PnlPrice:  Price5,
		},
	}

	ValidSingleVEPrice = []vetypes.PricePair{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5Bytes,
			PnlPrice:  Price5Bytes,
		},
	}

	ValidVEPricesWithOneInvalid = []vetypes.PricePair{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5Bytes,
			PnlPrice:  Price5Bytes,
		},
		{
			MarketId:  MarketId1,
			SpotPrice: Price6Bytes,
			PnlPrice:  Price6Bytes,
		},
		{
			MarketId:  MarketId2,
			SpotPrice: []byte("invalid"),
			PnlPrice:  []byte("invalid"),
		},
	}

	ValidVEPrices = []vetypes.PricePair{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5Bytes,
			PnlPrice:  Price5Bytes,
		},
		{
			MarketId:  MarketId1,
			SpotPrice: Price6Bytes,
			PnlPrice:  Price6Bytes,
		},
		{
			MarketId:  MarketId2,
			SpotPrice: Price7Bytes,
			PnlPrice:  Price7Bytes,
		},
	}

	InvalidVEPriceBytes = []vetypes.PricePair{
		{
			MarketId:  MarketId0,
			SpotPrice: Price5NegativeBytes,
			PnlPrice:  Price5NegativeBytes,
		},
		{
			MarketId:  MarketId1,
			SpotPrice: Price6NegativeBytes,
			PnlPrice:  Price6NegativeBytes,
		},
		{
			MarketId:  MarketId2,
			SpotPrice: Price7NegativeBytes,
			PnlPrice:  Price7NegativeBytes,
		},
	}

	InvalidVePricesMarketIds = []vetypes.PricePair{
		{
			MarketId:  99,
			SpotPrice: Price5Bytes,
			PnlPrice:  Price5Bytes,
		},
		{
			MarketId:  101,
			SpotPrice: Price6Bytes,
			PnlPrice:  Price6Bytes,
		},
		{
			MarketId:  102,
			SpotPrice: Price7Bytes,
			PnlPrice:  Price7Bytes,
		},
	}

	ValidEmptyMarketParams         = []types.MarketParam{}
	EmptyUpdateMarketPrices        = &types.MarketPriceUpdates{}
	EmptyUpdateMarketPricesTxBytes []byte

	ValidUpdateMarketPrices = &types.MarketPriceUpdates{
		MarketPriceUpdates: ValidMarketPriceUpdates,
	}
	ValidUpdateMarketPricesTxBytes []byte

	InvalidUpdateMarketPricesStateless = &types.MarketPriceUpdates{
		MarketPriceUpdates: []*types.MarketPriceUpdate{
			{
				MarketId:  MarketId0,
				SpotPrice: 0, // 0 price value is invalid.
				PnlPrice:  0, // 0 price value is invalid.
			},
		},
	}
	InvalidUpdateMarketPricesStatelessTxBytes []byte

	InvalidUpdateMarketPricesStateful = &types.MarketPriceUpdates{
		MarketPriceUpdates: []*types.MarketPriceUpdate{
			{
				MarketId:  MarketId0,
				SpotPrice: Price5,
				PnlPrice:  Price5,
			},
			{
				MarketId:  MarketId1,
				SpotPrice: Price6,
				PnlPrice:  Price6,
			},
			{
				MarketId:  99,
				SpotPrice: Price3, // Market with id 99 does not exist.
				PnlPrice:  Price3, // Market with id 99 does not exist.
			},
		},
	}
	InvalidUpdateMarketPricesStatefulTxBytes []byte

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
				Id:        uint32(0),
				Exponent:  BtcUsdExponent,
				SpotPrice: FiveBillion, // $50,000 == 1 BTC
				PnlPrice:  FiveBillion, // $50,000 == 1 BTC
			},
			{
				Id:        uint32(1),
				Exponent:  EthUsdExponent,
				SpotPrice: ThreeBillion, // $3,000 == 1 ETH
				PnlPrice:  ThreeBillion, // $3,000 == 1 ETH
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
				Id:        uint32(0),
				Exponent:  BtcUsdExponent,
				SpotPrice: FiveBillion, // $50,000 == 1 BTC
				PnlPrice:  FiveBillion, // $50,000 == 1 BTC
			},
			{ // ETH-USD
				Id:        uint32(1),
				Exponent:  EthUsdExponent,
				SpotPrice: ThreeBillion, // $3,000 == 1 ETH
				PnlPrice:  ThreeBillion, // $3,000 == 1 ETH
			},
		},
	}
)
