package constants

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	"github.com/dydxprotocol/v4/x/prices/types"
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
	BtcUsdPair   = "BTC-USD"
	EthUsdPair   = "ETH-USD"
	MaticUsdPair = "MATIC-USD"
	SolUsdPair   = "SOL-USD"
	LtcUsdPair   = "LTC-USD"

	BtcUsdExponent   = -5
	EthUsdExponent   = -6
	LinkUsdExponent  = -8
	MaticUsdExponent = -9
	CrvUsdExponent   = -10
	SolUsdExponent   = -8
	LtcUsdExponent   = -7

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

// The `MarketParam.ExchangeConfigJson` field is left unset as it is not used by the server.
var TestMarketParams = []types.MarketParam{
	{
		Id:                0,
		Pair:              BtcUsdPair,
		Exponent:          BtcUsdExponent,
		MinExchanges:      1,
		MinPriceChangePpm: 50,
	},
	{
		Id:                1,
		Pair:              EthUsdPair,
		Exponent:          EthUsdExponent,
		MinExchanges:      1,
		MinPriceChangePpm: 50,
	},
	{
		Id:                2,
		Pair:              SolUsdPair,
		Exponent:          SolUsdExponent,
		MinExchanges:      1,
		MinPriceChangePpm: 50,
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

	marketExchangeConfigs = constants.GenerateExchangeConfigJson(constants.StaticExchangeMarketConfig)

	Prices_DefaultGenesisState = types.GenesisState{
		// `ExchangeConfigJson` is left unset as it is not used by the server.
		MarketParams: []types.MarketParam{
			{
				Id:                 uint32(0),
				Pair:               BtcUsdPair,
				Exponent:           BtcUsdExponent,
				MinExchanges:       uint32(2),
				ExchangeConfigJson: marketExchangeConfigs[exchange_common.MARKET_BTC_USD],
				MinPriceChangePpm:  uint32(50),
			},
			{
				Id:                 uint32(1),
				Pair:               EthUsdPair,
				Exponent:           EthUsdExponent,
				MinExchanges:       uint32(1),
				ExchangeConfigJson: marketExchangeConfigs[exchange_common.MARKET_ETH_USD],
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
		// `ExchangeConfigJson` is left unset as it is unused by the server.
		MarketParams: []types.MarketParam{
			{ // BTC-USD
				Id:                 uint32(0),
				Pair:               BtcUsdPair,
				Exponent:           BtcUsdExponent,
				ExchangeConfigJson: marketExchangeConfigs[exchange_common.MARKET_BTC_USD],
				MinExchanges:       uint32(2),
				MinPriceChangePpm:  uint32(50),
			},
			{ // ETH-USD
				Id:                 uint32(1),
				Pair:               EthUsdPair,
				Exponent:           EthUsdExponent,
				ExchangeConfigJson: marketExchangeConfigs[exchange_common.MARKET_ETH_USD],
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
