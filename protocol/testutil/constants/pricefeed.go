package constants

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/constants"
	daemonClientTypes "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/client"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

var (
	// Markets
	MarketId0 = uint32(0)
	MarketId1 = uint32(1)
	MarketId2 = uint32(2)
	MarketId3 = uint32(3)
	MarketId4 = uint32(4)

	MarketId7  = uint32(7)
	MarketId8  = uint32(8)
	MarketId9  = uint32(9)
	MarketId10 = uint32(10)
	MarketId11 = uint32(11)

	// Exponents
	Exponent9 = int32(-9)
	Exponent8 = int32(-8)
	Exponent7 = int32(-7)

	// Exchanges
	ExchangeId0 = "Exchange0"
	ExchangeId1 = "Exchange1"
	ExchangeId2 = "Exchange2"
	ExchangeId3 = "Exchange3"

	// ExchangeArray
	Exchange1Exchange2Array = []string{
		ExchangeId1,
		ExchangeId2,
	}

	// MarketPriceTimestamps
	Market8_TimeTMinusThreshold_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId8,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price1,
	}
	Market8_TimeTMinusThreshold_Price2 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId8,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price2,
	}
	Market8_TimeTMinusThreshold_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId8,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price3,
	}
	Market8_TimeT_Price2 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId8,
		LastUpdatedAt: TimeT,
		Price:         Price2,
	}
	Market8_TimeT_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId8,
		LastUpdatedAt: TimeT,
		Price:         Price3,
	}
	Market8_TimeT_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId8,
		LastUpdatedAt: TimeT,
		Price:         Price1,
	}
	Market9_TimeTMinusThreshold_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price1,
	}
	Market9_TimeTMinusThreshold_Price2 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price2,
	}
	Market9_TimeTMinusThreshold_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price3,
	}
	Market9_TimeT_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeT,
		Price:         Price1,
	}
	Market9_TimeT_Price2 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeT,
		Price:         Price2,
	}
	Market9_TimeT_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeT,
		Price:         Price3,
	}
	Market9_TimeTPlusThreshold_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeTPlusThreshold,
		Price:         Price1,
	}
	Market9_TimeTPlusThreshold_Price2 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeTPlusThreshold,
		Price:         Price2,
	}
	Market9_TimeTPlusThreshold_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId9,
		LastUpdatedAt: TimeTPlusThreshold,
		Price:         Price3,
	}
	Market10_TimeT_Price4 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId10,
		LastUpdatedAt: TimeT,
		Price:         Price4,
	}
	Market11_TimeT_Price5 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId11,
		LastUpdatedAt: TimeT,
		Price:         Price5,
	}
	Market7_BeforeTimeT_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId7,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price1,
	}
	Market7_BeforeTimeT_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId7,
		LastUpdatedAt: TimeTMinusThreshold,
		Price:         Price3,
	}
	Market7_TimeT_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId7,
		LastUpdatedAt: TimeT,
		Price:         Price1,
	}
	Market7_TimeT_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId7,
		LastUpdatedAt: TimeT,
		Price:         Price3,
	}
	Market7_TimeTPlusThreshold_Price1 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId7,
		LastUpdatedAt: TimeTPlusThreshold,
		Price:         Price1,
	}
	Market7_TimeTPlusThreshold_Price3 = &daemonClientTypes.MarketPriceTimestamp{
		MarketId:      MarketId7,
		LastUpdatedAt: TimeTPlusThreshold,
		Price:         Price3,
	}

	// Prices
	InvalidPrice uint64 = 0
	Price1       uint64 = 1001
	Price2       uint64 = 2002
	Price3       uint64 = 3003
	Price4       uint64 = 4004
	Price5       uint64 = 500005
	Price6       uint64 = 60006
	Price7       uint64 = 7007

	// Exchange 0 prices
	Exchange0_Price4_TimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId1,
		Price:          Price4,
		LastUpdateTime: &TimeT,
	}

	// Exchange 1 prices
	Exchange1_Price1_TimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId1,
		Price:          Price1,
		LastUpdateTime: &TimeT,
	}
	Exchange1_Price2_AfterTimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId1,
		Price:          Price2,
		LastUpdateTime: &TimeTPlusThreshold,
	}
	Exchange1_Price3_BeforeTimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId1,
		Price:          Price3,
		LastUpdateTime: &TimeTMinusThreshold,
	}

	// Exchange 2 prices
	Exchange2_Price2_TimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId2,
		Price:          Price2,
		LastUpdateTime: &TimeT,
	}
	Exchange2_Price3_AfterTimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId2,
		Price:          Price3,
		LastUpdateTime: &TimeTPlusThreshold,
	}
	Exchange2_Price1_BeforeTimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId2,
		Price:          Price1,
		LastUpdateTime: &TimeTMinusThreshold,
	}

	// Exchange 3 prices
	Exchange3_Price3_TimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId3,
		Price:          Price3,
		LastUpdateTime: &TimeT,
	}
	Exchange3_Price4_AfterTimeT = &api.ExchangePrice{
		ExchangeId:     ExchangeId3,
		Price:          Price4,
		LastUpdateTime: &TimeTPlusThreshold,
	}

	// Price Updates
	Market9_SingleExchange_AtTimeUpdate = []*api.MarketPriceUpdate{
		{
			MarketId: MarketId9,
			ExchangePrices: []*api.ExchangePrice{
				Exchange1_Price1_TimeT,
			},
		},
	}
	AtTimeTPriceUpdate = []*api.MarketPriceUpdate{
		{
			MarketId: MarketId9,
			ExchangePrices: []*api.ExchangePrice{
				Exchange1_Price1_TimeT,
				Exchange2_Price2_TimeT,
			},
		},
		{
			MarketId: MarketId8,
			ExchangePrices: []*api.ExchangePrice{
				Exchange2_Price2_TimeT,
				Exchange3_Price3_TimeT,
			},
		},
		{
			MarketId: MarketId7,
			ExchangePrices: []*api.ExchangePrice{
				Exchange1_Price1_TimeT,
				Exchange3_Price3_TimeT,
			},
		},
	}

	AtTimeTSingleExchangePriceUpdate = []*api.MarketPriceUpdate{
		{
			MarketId: MarketId0,
			ExchangePrices: []*api.ExchangePrice{
				Exchange0_Price4_TimeT,
			},
		},
		{
			MarketId: MarketId1,
			ExchangePrices: []*api.ExchangePrice{
				Exchange1_Price1_TimeT,
			},
		},
		{
			MarketId: MarketId2,
			ExchangePrices: []*api.ExchangePrice{
				Exchange2_Price2_TimeT,
			},
		},
		{
			MarketId: MarketId3,
			ExchangePrices: []*api.ExchangePrice{
				Exchange3_Price3_TimeT,
			},
		},
		{
			MarketId: MarketId4,
			ExchangePrices: []*api.ExchangePrice{
				Exchange3_Price3_TimeT,
			},
		},
	}
	AtTimeTSingleExchangeSmoothedPrices = map[uint32]uint64{
		MarketId0: Exchange0_Price4_TimeT.Price,
		MarketId1: Exchange1_Price1_TimeT.Price,
		MarketId2: Exchange2_Price2_TimeT.Price,
		MarketId3: Exchange3_Price3_TimeT.Price,
		MarketId4: Exchange3_Price3_TimeT.Price,
	}

	AtTimeTSingleExchangeSmoothedPricesPlus10 = map[uint32]uint64{
		MarketId0: Exchange0_Price4_TimeT.Price + 10,
		MarketId1: Exchange1_Price1_TimeT.Price + 10,
		MarketId2: Exchange2_Price2_TimeT.Price + 10,
		MarketId3: Exchange3_Price3_TimeT.Price + 10,
		MarketId4: Exchange3_Price3_TimeT.Price + 10,
	}

	AtTimeTSingleExchangeSmoothedPricesPlus7 = map[uint32]uint64{
		MarketId0: Exchange0_Price4_TimeT.Price + 7,
		MarketId1: Exchange1_Price1_TimeT.Price + 7,
		MarketId2: Exchange2_Price2_TimeT.Price + 7,
		MarketId3: Exchange3_Price3_TimeT.Price + 7,
		MarketId4: Exchange3_Price3_TimeT.Price + 7,
	}

	MixedTimePriceUpdate = []*api.MarketPriceUpdate{
		{
			MarketId: MarketId9,
			ExchangePrices: []*api.ExchangePrice{
				Exchange1_Price1_TimeT,
				Exchange2_Price2_TimeT,
				Exchange3_Price3_TimeT,
			},
		},
		{
			MarketId: MarketId8,
			ExchangePrices: []*api.ExchangePrice{
				Exchange1_Price3_BeforeTimeT,
				Exchange2_Price2_TimeT,
				Exchange3_Price3_TimeT,
			},
		},
		{
			MarketId: MarketId7,
			ExchangePrices: []*api.ExchangePrice{
				Exchange1_Price3_BeforeTimeT,
				Exchange2_Price1_BeforeTimeT,
				Exchange3_Price3_TimeT,
			},
		},
	}

	// Markets
	AllMarketParamsMinExchanges2 = []types.MarketParam{
		{
			Id:           MarketId9,
			Exponent:     Exponent9,
			MinExchanges: 2,
		},
		{
			Id:           MarketId8,
			Exponent:     Exponent8,
			MinExchanges: 2,
		},
		{
			Id:           MarketId7,
			Exponent:     Exponent7,
			MinExchanges: 2,
		},
	}
	AllMarketParamsMinExchanges3 = []types.MarketParam{
		{
			Id:           MarketId9,
			MinExchanges: 3,
		},
		{
			Id:           MarketId8,
			MinExchanges: 3,
		},
		{
			Id:           MarketId7,
			MinExchanges: 3,
		},
	}

	// ExchangeConfig, MutableExchangeMarketConfig for various tests are defined below.

	SingleMarketExchangeQueryDetails = daemonClientTypes.ExchangeQueryDetails{IsMultiMarket: false}
	MultiMarketExchangeQueryDetails  = daemonClientTypes.ExchangeQueryDetails{IsMultiMarket: true}

	// ExchangeQueryConfigs.
	Exchange1_0MaxQueries_QueryConfig = daemonClientTypes.ExchangeQueryConfig{
		ExchangeId: ExchangeId1,
		IntervalMs: 100,
		TimeoutMs:  3_000,
		MaxQueries: 0,
	}

	Exchange1_1MaxQueries_QueryConfig = daemonClientTypes.ExchangeQueryConfig{
		ExchangeId: ExchangeId1,
		IntervalMs: 100,
		TimeoutMs:  3_000,
		MaxQueries: 1,
	}

	Exchange1_2MaxQueries_QueryConfig = daemonClientTypes.ExchangeQueryConfig{
		ExchangeId: ExchangeId1,
		IntervalMs: 100,
		TimeoutMs:  3_000,
		MaxQueries: 2,
	}

	// MutableExchangeMarketConfigs for 0, 1, 2, 3, and 5 markets.
	Exchange1_NoMarkets_MutableExchangeMarketConfig = daemonClientTypes.MutableExchangeMarketConfig{
		Id:                   ExchangeId1,
		MarketToMarketConfig: map[daemonClientTypes.MarketId]daemonClientTypes.MarketConfig{},
	}

	Exchange1_1Markets_MutableExchangeMarketConfig = daemonClientTypes.MutableExchangeMarketConfig{
		Id: ExchangeId1,
		MarketToMarketConfig: map[daemonClientTypes.MarketId]daemonClientTypes.MarketConfig{
			MarketId7: {
				Ticker: "BTC-USD",
			},
		},
	}

	Exchange1_2Markets_MutableExchangeMarketConfig = daemonClientTypes.MutableExchangeMarketConfig{
		Id: ExchangeId1,
		MarketToMarketConfig: map[daemonClientTypes.MarketId]daemonClientTypes.MarketConfig{
			MarketId7: {
				Ticker: "BTC-USD",
			},
			MarketId8: {
				Ticker: "ETH-USD",
			},
		},
	}

	Exchange1_3Markets_MutableExchangeMarketConfig = daemonClientTypes.MutableExchangeMarketConfig{
		Id: ExchangeId1,
		MarketToMarketConfig: map[daemonClientTypes.MarketId]daemonClientTypes.MarketConfig{
			MarketId7: {
				Ticker: "BTC-USD",
			},
			MarketId8: {
				Ticker: "ETH-USD",
			},
			MarketId9: {
				Ticker: "LTC-USD",
			},
		},
	}

	Exchange1_5Markets_MutableExchangeMarketConfig = daemonClientTypes.MutableExchangeMarketConfig{
		Id: ExchangeId1,
		MarketToMarketConfig: map[daemonClientTypes.MarketId]daemonClientTypes.MarketConfig{
			MarketId7: {
				Ticker: "BTC-USD",
			},
			MarketId8: {
				Ticker: "ETH-USD",
			},
			MarketId9: {
				Ticker: "LTC-USD",
			},
			MarketId10: {
				Ticker: "XRP-USD",
			},
			MarketId11: {
				Ticker: "BCH-USD",
			},
		},
	}

	CanonicalMarketExponents = map[daemonClientTypes.MarketId]daemonClientTypes.Exponent{
		MarketId7:  MutableMarketConfigs_5Markets[0].Exponent,
		MarketId8:  MutableMarketConfigs_5Markets[1].Exponent,
		MarketId9:  MutableMarketConfigs_5Markets[2].Exponent,
		MarketId10: MutableMarketConfigs_5Markets[3].Exponent,
		MarketId11: MutableMarketConfigs_5Markets[4].Exponent,
	}

	CanonicalMarketPriceTimestampResponses = map[uint32]*daemonClientTypes.MarketPriceTimestamp{
		MarketId7:  Market7_TimeTPlusThreshold_Price3,
		MarketId8:  Market8_TimeTMinusThreshold_Price2,
		MarketId9:  Market9_TimeT_Price1,
		MarketId10: Market10_TimeT_Price4,
		MarketId11: Market11_TimeT_Price5,
	}

	// ExchangeIdMarketPriceTimestamps
	ExchangeId1_Market9_TimeT_Price1 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId1,
		MarketPriceTimestamp: Market9_TimeT_Price1,
	}
	ExchangeId2_Market9_TimeT_Price2 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId2,
		MarketPriceTimestamp: Market9_TimeT_Price2,
	}
	ExchangeId3_Market9_TimeT_Price3 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId3,
		MarketPriceTimestamp: Market9_TimeT_Price3,
	}
	ExchangeId1_Market8_BeforeTimeT_Price3 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId1,
		MarketPriceTimestamp: Market8_TimeTMinusThreshold_Price3,
	}
	ExchangeId2_Market8_TimeT_Price2 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId2,
		MarketPriceTimestamp: Market8_TimeT_Price2,
	}
	ExchangeId3_Market8_TimeT_Price3 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId3,
		MarketPriceTimestamp: Market8_TimeT_Price3,
	}
	ExchangeId1_Market7_BeforeTimeT_Price3 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId1,
		MarketPriceTimestamp: Market7_BeforeTimeT_Price3,
	}
	ExchangeId2_Market7_BeforeTimeT_Price1 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId2,
		MarketPriceTimestamp: Market7_BeforeTimeT_Price1,
	}
	ExchangeId1_Market7_TimeT_Price1 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId1,
		MarketPriceTimestamp: Market7_TimeT_Price1,
	}
	ExchangeId3_Market7_TimeT_Price3 = &client.ExchangeIdMarketPriceTimestamp{
		ExchangeId:           ExchangeId3,
		MarketPriceTimestamp: Market7_TimeT_Price3,
	}

	CoinbaseMutableMarketConfig = &daemonClientTypes.MutableExchangeMarketConfig{
		Id: CoinbaseExchangeName,
		MarketToMarketConfig: map[daemonClientTypes.MarketId]daemonClientTypes.MarketConfig{
			MarketId7: {
				Ticker: "BTC-USD",
			},
			MarketId8: {
				Ticker: "ETH-USD",
			},
		},
	}
	BinanceMutableMarketConfig = &daemonClientTypes.MutableExchangeMarketConfig{
		Id: BinanceExchangeName,
		MarketToMarketConfig: map[daemonClientTypes.MarketId]daemonClientTypes.MarketConfig{
			MarketId7: {
				Ticker: "BTCUSDT",
			},
			MarketId8: {
				Ticker: "ETHUSDT",
			},
		},
	}

	TestCanonicalExchangeIds = []string{ExchangeId1, ExchangeId2}

	// Test constants for starting the daemon client.
	TestMutableExchangeMarketConfigs = map[string]*daemonClientTypes.MutableExchangeMarketConfig{
		CoinbaseExchangeName: CoinbaseMutableMarketConfig,
		BinanceExchangeName:  BinanceMutableMarketConfig,
	}

	TestMarket7And8Params = []types.MarketParam{
		{
			Id:       7,
			Pair:     BtcUsdPair,
			Exponent: BtcUsdExponent,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Binance","ticker":"BTCUSDT"},` +
				`{"exchangeName":"Coinbase","ticker":"BTC-USD"}]}`,
			MinExchanges:      1,
			MinPriceChangePpm: 50,
		},
		{
			Id:           8,
			Pair:         EthUsdPair,
			Exponent:     EthUsdExponent,
			MinExchanges: 1,
			ExchangeConfigJson: `{"exchanges":[{"exchangeName":"Binance","ticker":"ETHUSDT"},` +
				`{"exchangeName":"Coinbase","ticker":"ETH-USD"}]}`,
			MinPriceChangePpm: 50,
		},
	}

	TestMutableMarketConfigs = map[daemonClientTypes.MarketId]*daemonClientTypes.MutableMarketConfig{
		MarketId7: {
			Id:           MarketId7,
			Pair:         BtcUsdPair,
			Exponent:     BtcUsdExponent,
			MinExchanges: 1,
		},
		MarketId8: {
			Id:           MarketId8,
			Pair:         EthUsdPair,
			Exponent:     EthUsdExponent,
			MinExchanges: 1,
		},
	}

	// Pricefetcher MutableMarketConfigs for 0, 1, 2, 3 and 5 markets.
	MutableMarketConfigs_0Markets = []*daemonClientTypes.MutableMarketConfig{}

	MutableMarketConfigs_1Markets = []*daemonClientTypes.MutableMarketConfig{
		{
			Id:           MarketId7,
			Pair:         BtcUsdPair,
			Exponent:     BtcUsdExponent,
			MinExchanges: 1,
		},
	}

	MutableMarketConfigs_2Markets = []*daemonClientTypes.MutableMarketConfig{
		{
			Id:           MarketId7,
			Pair:         BtcUsdPair,
			Exponent:     BtcUsdExponent,
			MinExchanges: 1,
		},
		{
			Id:           MarketId8,
			Pair:         EthUsdPair,
			Exponent:     EthUsdExponent,
			MinExchanges: 1,
		},
	}

	MutableMarketConfigs_3Markets = []*daemonClientTypes.MutableMarketConfig{
		{
			Id:           MarketId7,
			Pair:         BtcUsdPair,
			Exponent:     BtcUsdExponent,
			MinExchanges: 1,
		},
		{
			Id:           MarketId8,
			Pair:         EthUsdPair,
			Exponent:     EthUsdExponent,
			MinExchanges: 1,
		},
		{
			Id:           MarketId9,
			Pair:         LtcUsdPair,
			Exponent:     LtcUsdExponent,
			MinExchanges: 1,
		},
	}

	MutableMarketConfigs_5Markets = []*daemonClientTypes.MutableMarketConfig{
		{
			Id:           MarketId7,
			Pair:         BtcUsdPair,
			Exponent:     BtcUsdExponent,
			MinExchanges: 1,
		},
		{
			Id:           MarketId8,
			Pair:         EthUsdPair,
			Exponent:     EthUsdExponent,
			MinExchanges: 1,
		},
		{
			Id:           MarketId9,
			Pair:         LtcUsdPair,
			Exponent:     LtcUsdExponent,
			MinExchanges: 1,
		},
		{
			Id:           MarketId10,
			Pair:         SolUsdPair,
			Exponent:     SolUsdExponent,
			MinExchanges: 1,
		},
		{
			Id:           MarketId11,
			Pair:         PolUsdPair,
			Exponent:     PolUsdExponent,
			MinExchanges: 1,
		},
	}

	MarketToMutableMarketConfigs_5Markets = map[daemonClientTypes.MarketId]*daemonClientTypes.MutableMarketConfig{
		MarketId7:  MutableMarketConfigs_5Markets[0],
		MarketId8:  MutableMarketConfigs_5Markets[1],
		MarketId9:  MutableMarketConfigs_5Markets[2],
		MarketId10: MutableMarketConfigs_5Markets[3],
		MarketId11: MutableMarketConfigs_5Markets[4],
	}

	// Expected exponents for above configs.
	MutableMarketConfigs_3Markets_ExpectedExponents = map[daemonClientTypes.MarketId]daemonClientTypes.Exponent{
		MarketId7: BtcUsdExponent,
		MarketId8: EthUsdExponent,
		MarketId9: LtcUsdExponent,
	}

	MutableMarketConfigs_5Markets_ExpectedExponents = map[daemonClientTypes.MarketId]daemonClientTypes.Exponent{
		MarketId7:  BtcUsdExponent,
		MarketId8:  EthUsdExponent,
		MarketId9:  LtcUsdExponent,
		MarketId10: SolUsdExponent,
		MarketId11: PolUsdExponent,
	}

	TestExchangeQueryConfigs = map[string]*daemonClientTypes.ExchangeQueryConfig{
		ExchangeId1: {
			ExchangeId: ExchangeId1,
			IntervalMs: 100,
			TimeoutMs:  3_000,
			MaxQueries: 2,
		},
		ExchangeId2: {
			ExchangeId: ExchangeId2,
			IntervalMs: 100,
			TimeoutMs:  3_000,
			MaxQueries: 2,
		},
	}
	TestExchangeIdToExchangeQueryDetails = map[string]daemonClientTypes.ExchangeQueryDetails{
		ExchangeId1: constants.StaticExchangeDetails[ExchangeId1],
		ExchangeId2: constants.StaticExchangeDetails[ExchangeId2],
	}
)
