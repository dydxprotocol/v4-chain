package constants

import (
	"github.com/dydxprotocol/v4/x/prices/types"
)

func init() {
	_ = TestTxBuilder.SetMsgs(ValidMsgUpdateMarketPrices)
	ValidMsgUpdateMarketPricesTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidMsgUpdateMarketPricesStateless)
	InvalidMsgUpdateMarketPricesStatelessTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidMsgUpdateMarketPricesStateful)
	InvalidMsgUpdateMarketPricesStatefulTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())
}

const (
	BtcUsdPair = "BTC-USD"
	EthUsdPair = "ETH-USD"
	SolUsdPair = "SOL-USD"

	BtcUsdExponent = -5
	EthUsdExponent = -6
	SolUsdExponent = -8

	CoinbaseExchangeName  = "Coinbase"
	BinanceExchangeName   = "Binance"
	BinanceUSExchangeName = "BinanceUS"
	BitfinexExchangeName  = "Bitfinex"

	FiveBillion  = uint64(5_000_000_000)
	ThreeBillion = uint64(3_000_000_000)
	FiveMillion  = uint64(5_000_000)
	OneMillion   = uint64(1_000_000)
)

var TestExchangeFeeds = []types.ExchangeFeed{
	{
		Id:   0,
		Name: CoinbaseExchangeName,
		Memo: "test memo 0",
	},
	{
		Id:   1,
		Name: BinanceExchangeName,
		Memo: "test memo 1",
	},
	{
		Id:   2,
		Name: BitfinexExchangeName,
		Memo: "test memo 2",
	},
}

var TestMarkets = []types.Market{
	{
		Id:                0,
		Pair:              BtcUsdPair,
		Exponent:          BtcUsdExponent,
		Exchanges:         []uint32{0, 1},
		MinExchanges:      1,
		MinPriceChangePpm: 50,
		Price:             FiveBillion, // $50,000 == 1 BTC.
	},
	{
		Id:                1,
		Pair:              EthUsdPair,
		Exponent:          EthUsdExponent,
		Exchanges:         []uint32{1, 2},
		MinExchanges:      1,
		MinPriceChangePpm: 50,
		Price:             ThreeBillion, // $3,000 == 1 ETH.
	},
	{
		Id:                2,
		Pair:              SolUsdPair,
		Exponent:          SolUsdExponent,
		Exchanges:         []uint32{0, 2},
		MinExchanges:      1,
		MinPriceChangePpm: 50,
		Price:             FiveBillion, // $50 == 1 SOL.
	},
}

var TestPricesGenesisState = types.GenesisState{
	ExchangeFeeds: TestExchangeFeeds,
	Markets:       TestMarkets,
}

var (
	ValidMarketPriceUpdates = []*types.MsgUpdateMarketPrices_MarketPrice{
		types.NewMarketPriceUpdate(MarketId0, Price5),
		types.NewMarketPriceUpdate(MarketId1, Price6),
		types.NewMarketPriceUpdate(MarketId2, Price7),
	}

	// `MsgUpdateMarketPrices`.
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
		ExchangeFeeds: []types.ExchangeFeed{
			{
				Id:   uint32(0),
				Name: CoinbaseExchangeName,
				Memo: "test memo 0",
			},
			{
				Id:   uint32(1),
				Name: BinanceExchangeName,
				Memo: "test memo 1",
			},
		},
		Markets: []types.Market{{
			Pair:              BtcUsdPair,
			Exchanges:         []uint32{0, 1},
			Exponent:          BtcUsdExponent,
			MinExchanges:      uint32(2),
			MinPriceChangePpm: uint32(50),
		}},
	}

	Prices_MultiExchangeMarketGenesisState = types.GenesisState{
		ExchangeFeeds: []types.ExchangeFeed{
			{ // Binance
				Id:   uint32(0),
				Name: BinanceExchangeName,
				Memo: "test memo 0",
			},
			{ // BinanceUS
				Id:   uint32(1),
				Name: BinanceUSExchangeName,
				Memo: "test memo 1",
			},
			{ // Bitfinex
				Id:   uint32(2),
				Name: BitfinexExchangeName,
				Memo: "test memo 2",
			},
		},
		Markets: []types.Market{
			{ // BTC-USD
				Id:                uint32(0),
				Pair:              BtcUsdPair,
				Exchanges:         []uint32{0, 1, 2},
				Exponent:          BtcUsdExponent,
				MinExchanges:      uint32(2),
				MinPriceChangePpm: uint32(50),
			},
			{ // ETH-USD
				Id:                uint32(1),
				Pair:              EthUsdPair,
				Exchanges:         []uint32{0, 1, 2},
				Exponent:          EthUsdExponent,
				MinExchanges:      uint32(2),
				MinPriceChangePpm: uint32(50),
			},
		},
	}
)
