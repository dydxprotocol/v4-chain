package constants

import (
	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	daemonClientTypes "github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4/testutil/client"
	"github.com/dydxprotocol/v4/x/prices/types"
)

var (
	// Markets
	MarketId0 = uint32(0)
	MarketId1 = uint32(1)
	MarketId2 = uint32(2)

	MarketId7  = uint32(7)
	MarketId8  = uint32(8)
	MarketId9  = uint32(9)
	MarketId10 = uint32(10)
	MarketId11 = uint32(11)

	// Exchanges
	ExchangeFeedId0 = uint32(0)
	ExchangeFeedId1 = uint32(1)
	ExchangeFeedId2 = uint32(2)
	ExchangeFeedId3 = uint32(3)

	// ExchangeArray
	Exchange1Exchange2Array = []uint32{
		ExchangeFeedId1,
		ExchangeFeedId2,
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

	// Valid Exchanges
	ValidExchanges1   = map[uint32]bool{ExchangeFeedId1: true}
	ValidExchangesAll = map[uint32]bool{
		ExchangeFeedId1: true,
		ExchangeFeedId2: true,
		ExchangeFeedId3: true,
	}

	// Price module's Valid Exchanges
	PriceModuleValidExchanges = []uint32{ExchangeFeedId1, ExchangeFeedId2, ExchangeFeedId3}

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
		ExchangeFeedId: ExchangeFeedId0,
		Price:          Price4,
		LastUpdateTime: &TimeT,
	}

	// Exchange 1 prices
	Exchange1_Price1_TimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId1,
		Price:          Price1,
		LastUpdateTime: &TimeT,
	}
	Exchange1_Price2_AfterTimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId1,
		Price:          Price2,
		LastUpdateTime: &TimeTPlusThreshold,
	}
	Exchange1_Price3_BeforeTimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId1,
		Price:          Price3,
		LastUpdateTime: &TimeTMinusThreshold,
	}

	// Exchange 2 prices
	Exchange2_Price2_TimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId2,
		Price:          Price2,
		LastUpdateTime: &TimeT,
	}
	Exchange2_Price3_AfterTimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId2,
		Price:          Price3,
		LastUpdateTime: &TimeTPlusThreshold,
	}
	Exchange2_Price1_BeforeTimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId2,
		Price:          Price1,
		LastUpdateTime: &TimeTMinusThreshold,
	}

	// Exchange 3 prices
	Exchange3_Price3_TimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId3,
		Price:          Price3,
		LastUpdateTime: &TimeT,
	}
	Exchange3_Price4_AfterTimeT = &api.ExchangePrice{
		ExchangeFeedId: ExchangeFeedId3,
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
	}
	AtTimeTSingleExchangeSmoothedPrices = map[uint32]uint64{
		MarketId0: Exchange0_Price4_TimeT.Price,
		MarketId1: Exchange1_Price1_TimeT.Price,
		MarketId2: Exchange2_Price2_TimeT.Price,
	}

	AtTimeTSingleExchangeSmoothedPricesPlus10 = map[uint32]uint64{
		MarketId0: Exchange0_Price4_TimeT.Price + 10,
		MarketId1: Exchange1_Price1_TimeT.Price + 10,
		MarketId2: Exchange2_Price2_TimeT.Price + 10,
	}

	AtTimeTSingleExchangeSmoothedPricesPlus7 = map[uint32]uint64{
		MarketId0: Exchange0_Price4_TimeT.Price + 7,
		MarketId1: Exchange1_Price1_TimeT.Price + 7,
		MarketId2: Exchange2_Price2_TimeT.Price + 7,
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
	AllMarketsMinExchanges2 = []types.Market{
		{
			Id:           MarketId9,
			Exchanges:    PriceModuleValidExchanges,
			MinExchanges: 2,
		},
		{
			Id:           MarketId8,
			Exchanges:    PriceModuleValidExchanges,
			MinExchanges: 2,
		},
		{
			Id:           MarketId7,
			Exchanges:    PriceModuleValidExchanges,
			MinExchanges: 2,
		},
	}
	AllMarketsMinExchanges3 = []types.Market{
		{
			Id:           MarketId9,
			Exchanges:    PriceModuleValidExchanges,
			MinExchanges: 3,
		},
		{
			Id:           MarketId8,
			Exchanges:    PriceModuleValidExchanges,
			MinExchanges: 3,
		},
		{
			Id:           MarketId7,
			Exchanges:    PriceModuleValidExchanges,
			MinExchanges: 3,
		},
	}

	// ExchangeConfig
	Exchange1_NoMarkets_0MaxQueries_Config = daemonClientTypes.ExchangeConfig{
		Markets: []uint32{},
		ExchangeStartupConfig: daemonClientTypes.ExchangeStartupConfig{
			ExchangeFeedId: ExchangeFeedId1,
			IntervalMs:     100,
			TimeoutMs:      3_000,
			MaxQueries:     0,
		},
	}
	Exchange1_1Markets_1MaxQueries_Config = daemonClientTypes.ExchangeConfig{
		Markets: []uint32{
			MarketId7,
		},
		ExchangeStartupConfig: daemonClientTypes.ExchangeStartupConfig{
			ExchangeFeedId: ExchangeFeedId1,
			IntervalMs:     100,
			TimeoutMs:      3_000,
			MaxQueries:     1,
		},
	}
	Exchange1_1Markets_2MaxQueries_Config = daemonClientTypes.ExchangeConfig{
		Markets: []uint32{
			MarketId7,
		},
		ExchangeStartupConfig: daemonClientTypes.ExchangeStartupConfig{
			ExchangeFeedId: ExchangeFeedId1,
			IntervalMs:     100,
			TimeoutMs:      3_000,
			MaxQueries:     2,
		},
	}
	Exchange1_2Markets_2MaxQueries_Config = daemonClientTypes.ExchangeConfig{
		Markets: []uint32{
			MarketId7,
			MarketId8,
		},
		ExchangeStartupConfig: daemonClientTypes.ExchangeStartupConfig{
			ExchangeFeedId: ExchangeFeedId1,
			IntervalMs:     100,
			TimeoutMs:      3_000,
			MaxQueries:     2,
		},
	}
	Exchange1_2Markets_Multimarket_Config = daemonClientTypes.ExchangeConfig{
		Markets: []uint32{
			MarketId7,
			MarketId8,
		},
		ExchangeStartupConfig: daemonClientTypes.ExchangeStartupConfig{
			ExchangeFeedId: ExchangeFeedId1,
			IntervalMs:     100,
			TimeoutMs:      3_000,
			MaxQueries:     1,
		},
		IsMultiMarket: true,
	}

	Exchange1_3Markets_2MaxQueries_Config = daemonClientTypes.ExchangeConfig{
		Markets: []uint32{
			MarketId7,
			MarketId8,
			MarketId9,
		},
		ExchangeStartupConfig: daemonClientTypes.ExchangeStartupConfig{
			ExchangeFeedId: ExchangeFeedId1,
			IntervalMs:     100,
			TimeoutMs:      3_000,
			MaxQueries:     2,
		},
	}

	Exchange1_5Markets_Multimarket_Config = daemonClientTypes.ExchangeConfig{
		Markets: []uint32{
			MarketId7,
			MarketId8,
			MarketId9,
			MarketId10,
			MarketId11,
		},
		ExchangeStartupConfig: daemonClientTypes.ExchangeStartupConfig{
			ExchangeFeedId: ExchangeFeedId1,
			IntervalMs:     100,
			TimeoutMs:      3_000,
			MaxQueries:     2,
		},
		IsMultiMarket: true,
	}

	// ExchangeStartupConfig
	Exchange1 = daemonClientTypes.ExchangeStartupConfig{
		ExchangeFeedId: 1,
		IntervalMs:     3000,
		TimeoutMs:      3_000,
		MaxQueries:     3,
	}

	CanonicalMarketPriceTimestampResponses = map[uint32]*daemonClientTypes.MarketPriceTimestamp{
		MarketId7:  Market7_TimeTPlusThreshold_Price3,
		MarketId8:  Market8_TimeTMinusThreshold_Price2,
		MarketId9:  Market9_TimeT_Price1,
		MarketId10: Market10_TimeT_Price4,
		MarketId11: Market11_TimeT_Price5,
	}

	// ExchangeFeedIdMarketPriceTimestamps
	ExchangeId1_Market9_TimeT_Price1 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId1,
		MarketPriceTimestamp: Market9_TimeT_Price1,
	}
	ExchangeId2_Market9_TimeT_Price2 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId2,
		MarketPriceTimestamp: Market9_TimeT_Price2,
	}
	ExchangeId3_Market9_TimeT_Price3 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId3,
		MarketPriceTimestamp: Market9_TimeT_Price3,
	}
	ExchangeId1_Market8_BeforeTimeT_Price3 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId1,
		MarketPriceTimestamp: Market8_TimeTMinusThreshold_Price3,
	}
	ExchangeId2_Market8_TimeT_Price2 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId2,
		MarketPriceTimestamp: Market8_TimeT_Price2,
	}
	ExchangeId3_Market8_TimeT_Price3 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId3,
		MarketPriceTimestamp: Market8_TimeT_Price3,
	}
	ExchangeId1_Market7_BeforeTimeT_Price3 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId1,
		MarketPriceTimestamp: Market7_BeforeTimeT_Price3,
	}
	ExchangeId2_Market7_BeforeTimeT_Price1 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId2,
		MarketPriceTimestamp: Market7_BeforeTimeT_Price1,
	}
	ExchangeId1_Market7_TimeT_Price1 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId1,
		MarketPriceTimestamp: Market7_TimeT_Price1,
	}
	ExchangeId3_Market7_TimeT_Price3 = &client.ExchangeFeedIdMarketPriceTimestamp{
		ExchangeFeedId:       ExchangeFeedId3,
		MarketPriceTimestamp: Market7_TimeT_Price3,
	}
)
