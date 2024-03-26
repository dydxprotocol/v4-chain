package constants

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

func init() {
	_ = TestTxBuilder.SetMsgs(EmptyMsgAddPremiumVotes)
	EmptyMsgAddPremiumVotesTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(ValidMsgAddPremiumVotes)
	ValidMsgAddPremiumVotesTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())

	_ = TestTxBuilder.SetMsgs(InvalidMsgAddPremiumVotes)
	InvalidMsgAddPremiumVotesTxBytes, _ = TestEncodingCfg.TxConfig.TxEncoder()(TestTxBuilder.GetTx())
}

// LiquidityTier objects.
var LiquidityTiers = []perptypes.LiquidityTier{
	{
		Id:                     0,
		Name:                   "0",
		InitialMarginPpm:       1_000_000,
		MaintenanceFractionPpm: 1_000_000,
		ImpactNotional:         500_000_000,
	},
	{
		Id:                     1,
		Name:                   "1",
		InitialMarginPpm:       1_000_000,
		MaintenanceFractionPpm: 750_000,
		ImpactNotional:         500_000_000,
	},
	{
		Id:                     2,
		Name:                   "2",
		InitialMarginPpm:       1_000_000,
		MaintenanceFractionPpm: 0,
		ImpactNotional:         500_000_000,
	},
	{
		Id:                     3,
		Name:                   "3",
		InitialMarginPpm:       200_000,
		MaintenanceFractionPpm: 500_000,
		ImpactNotional:         2_500_000_000,
	},
	{
		Id:                     4,
		Name:                   "4",
		InitialMarginPpm:       500_000,
		MaintenanceFractionPpm: 800_000,
		ImpactNotional:         1_000_000_000,
	},
	{
		Id:                     5,
		Name:                   "5",
		InitialMarginPpm:       500_000,
		MaintenanceFractionPpm: 600_000,
		ImpactNotional:         1_000_000_000,
	},
	{
		Id:                     6,
		Name:                   "6",
		InitialMarginPpm:       200_000,
		MaintenanceFractionPpm: 900_000,
		ImpactNotional:         2_500_000_000,
	},
	{
		Id:                     7,
		Name:                   "7",
		InitialMarginPpm:       0,
		MaintenanceFractionPpm: 0,
		ImpactNotional:         1_000_000_000,
	},
	{
		Id:                     8,
		Name:                   "8",
		InitialMarginPpm:       9_910, // 0.9910%
		MaintenanceFractionPpm: 1_000_000,
		ImpactNotional:         50_454_000_000,
	},
	{
		Id:                     9,
		Name:                   "9",
		InitialMarginPpm:       200_000, // 20%
		MaintenanceFractionPpm: 500_000, // 20% * 0.5 = 10%
		ImpactNotional:         2_500_000_000,
		OpenInterestUpperCap:   50_000_000_000_000, // 50mm USDC
		OpenInterestLowerCap:   25_000_000_000_000, // 25mm USDC
	},
	{
		Id:                     101,
		Name:                   "101",
		InitialMarginPpm:       200_000,
		MaintenanceFractionPpm: 500_000,
		ImpactNotional:         2_500_000_000,
	},
}

// Perpetual OI setup in tests
var (
	BtcUsd_OpenInterest1_AtomicRes8 = perptypes.OpenInterestDelta{
		PerpetualId:  0,
		BaseQuantums: big.NewInt(100_000_000),
	}
	EthUsd_OpenInterest1_AtomicRes9 = perptypes.OpenInterestDelta{
		PerpetualId:  1,
		BaseQuantums: big.NewInt(1_000_000_000),
	}
	DefaultTestPerpOIs = []perptypes.OpenInterestDelta{
		BtcUsd_OpenInterest1_AtomicRes8,
		EthUsd_OpenInterest1_AtomicRes9,
	}
)

// Perpetual genesis parameters.
const TestFundingRateClampFactorPpm = 6_000_000
const TestPremiumVoteClampFactorPpm = 60_000_000
const TestMinNumVotesPerSample = 15

var PerpetualsGenesisParams = perptypes.Params{
	FundingRateClampFactorPpm: TestFundingRateClampFactorPpm,
	PremiumVoteClampFactorPpm: TestPremiumVoteClampFactorPpm,
	MinNumVotesPerSample:      TestMinNumVotesPerSample,
}

var Perpetuals_GenesisState_ParamsOnly = perptypes.GenesisState{
	Params: PerpetualsGenesisParams,
}

// Perpetual objects.
var (
	BtcUsd_InvalidMarketId = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD invalid market Id",
			MarketId:          uint32(9999),
			AtomicResolution:  int32(-10),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(0),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_0DefaultFunding_0AtomicResolution = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 0 percent default funding, 0 atomic resolution",
			MarketId:          uint32(0),
			AtomicResolution:  int32(0),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(2),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_NegativeDefaultFunding_10AtomicResolution = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD -0.001 percent percent default funding",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-10),
			DefaultFundingPpm: int32(-1_000),
			LiquidityTier:     uint32(1),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_0DefaultFunding_10AtomicResolution = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 0 percent default funding",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-10),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(1),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 0 percent default funding, 20% IM, 18% MM",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-10),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(6),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_0_001Percent_DefaultFunding_10AtomicResolution = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                10000,
			Ticker:            "BTC-USD 0.001 percent default funding",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-10),
			DefaultFundingPpm: int32(1000), // 0.001%
			LiquidityTier:     uint32(1),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_SmallMarginRequirement = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD small margin requirement",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(8),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_100PercentMarginRequirement = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 100% margin requirement",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(0),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_50PercentInitial_40PercentMaintenance = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 50/40 margin requirements",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(4),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_20PercentInitial_10PercentMaintenance = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 20/10 margin requirements",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest1 = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 20/10 margin requirements",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.NewInt(100_000_000),
	}
	BtcUsd_20PercentInitial_10PercentMaintenance_OpenInterest2 = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 20/10 margin requirements",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.NewInt(200_000_000),
	}
	BtcUsd_20PercentInitial_10PercentMaintenance_25mmLowerCap_50mmUpperCap = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD 20/10 margin requirements",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(9),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	BtcUsd_NoMarginRequirement = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD no margin requirement",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-8),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(7),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	EthUsd_0DefaultFunding_9AtomicResolution = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                1,
			Ticker:            "ETH-USD default fundingm, -9 atomic resolution",
			MarketId:          uint32(1),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(5),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	EthUsd_NoMarginRequirement = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                1,
			Ticker:            "ETH-USD no margin requirement",
			MarketId:          uint32(1),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(7),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	EthUsd_20PercentInitial_10PercentMaintenance = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                1,
			Ticker:            "ETH-USD 20/10 margin requirements",
			MarketId:          uint32(1),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	EthUsd_100PercentMarginRequirement = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                1,
			Ticker:            "ETH-USD 100/100 margin requirements",
			MarketId:          uint32(1),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(0),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	SolUsd_20PercentInitial_10PercentMaintenance = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                2,
			Ticker:            "SOL-USD 20/10 margin requirements",
			MarketId:          uint32(2),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	IsoUsd_IsolatedMarket = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                3,
			Ticker:            "ISO-USD",
			MarketId:          uint32(3),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
	Iso2Usd_IsolatedMarket = perptypes.Perpetual{
		Params: perptypes.PerpetualParams{
			Id:                4,
			Ticker:            "ISO2-USD",
			MarketId:          uint32(4),
			AtomicResolution:  int32(-7),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
		},
		FundingIndex: dtypes.ZeroInt(),
		OpenInterest: dtypes.ZeroInt(),
	}
)

var TestMarketPerpetuals = []perptypes.Perpetual{
	{
		Params: perptypes.PerpetualParams{
			Id:                0,
			Ticker:            "BTC-USD",
			MarketId:          uint32(0),
			AtomicResolution:  int32(-10),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(0),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
	},
	{
		Params: perptypes.PerpetualParams{
			Id:                1,
			Ticker:            "ETH-USD",
			MarketId:          uint32(1),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(0),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
	},
	{
		Params: perptypes.PerpetualParams{
			Id:                2,
			Ticker:            "SOL-USD",
			MarketId:          uint32(2),
			AtomicResolution:  int32(-9),
			DefaultFundingPpm: int32(0),
			LiquidityTier:     uint32(3),
			MarketType:        perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
		},
		FundingIndex: dtypes.ZeroInt(),
	},
	IsoUsd_IsolatedMarket,
	Iso2Usd_IsolatedMarket,
}

// AddPremiumVotes messages.
var (
	TestAddPremiumVotesMsg = &perptypes.MsgAddPremiumVotes{
		Votes: []perptypes.FundingPremium{
			{
				PerpetualId: 0,
				PremiumPpm:  1000,
			},
		},
	}

	EmptyMsgAddPremiumVotes        = &perptypes.MsgAddPremiumVotes{}
	EmptyMsgAddPremiumVotesTxBytes []byte

	ValidMsgAddPremiumVotes = &perptypes.MsgAddPremiumVotes{
		Votes: []perptypes.FundingPremium{
			{PerpetualId: 1, PremiumPpm: 1_000},
			{PerpetualId: 2, PremiumPpm: 2_000},
		},
	}
	ValidMsgAddPremiumVotesTxBytes []byte

	InvalidMsgAddPremiumVotes = &perptypes.MsgAddPremiumVotes{
		Votes: []perptypes.FundingPremium{
			{PerpetualId: 3, PremiumPpm: 3_000}, // descending order is incorrect.
			{PerpetualId: 2, PremiumPpm: 2_000},
		},
	}
	InvalidMsgAddPremiumVotesTxBytes []byte

	Perpetuals_DefaultGenesisState = perptypes.GenesisState{
		LiquidityTiers: []perptypes.LiquidityTier{
			{
				Id:                     uint32(0),
				Name:                   "Large-Cap",
				InitialMarginPpm:       200_000,
				MaintenanceFractionPpm: 500_000,
				ImpactNotional:         2_500_000_000,
				OpenInterestLowerCap:   0,
				OpenInterestUpperCap:   0,
			},
			{
				Id:                     uint32(1),
				Name:                   "Mid-Cap",
				InitialMarginPpm:       300_000,
				MaintenanceFractionPpm: 600_000,
				ImpactNotional:         1_667_000_000,
				OpenInterestLowerCap:   25_000_000_000_000,
				OpenInterestUpperCap:   50_000_000_000_000,
			},
			{
				Id:                     uint32(2),
				Name:                   "Small-Cap",
				InitialMarginPpm:       400_000,
				MaintenanceFractionPpm: 700_000,
				ImpactNotional:         1_250_000_000,
				OpenInterestLowerCap:   10_000_000_000_000,
				OpenInterestUpperCap:   20_000_000_000_000,
			},
		},
		Params: PerpetualsGenesisParams,
		Perpetuals: []perptypes.Perpetual{
			{
				Params: perptypes.PerpetualParams{
					Id:            uint32(0),
					Ticker:        "genesis_test_ticker_0",
					LiquidityTier: 0,
					MarketType:    perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
				FundingIndex: dtypes.ZeroInt(),
				OpenInterest: dtypes.ZeroInt(),
			},
			{
				Params: perptypes.PerpetualParams{
					Id:            uint32(1),
					Ticker:        "genesis_test_ticker_1",
					LiquidityTier: 1,
					MarketType:    perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				},
				FundingIndex: dtypes.ZeroInt(),
				OpenInterest: dtypes.ZeroInt(),
			},
		},
	}
)

// Return a list of `count` constant funding premiums equal to `value`.
func GenerateConstantFundingPremiums(
	value int32,
	count uint32,
) (
	result []int32,
) {
	result = make([]int32, count)
	for i := uint32(0); i < count; i += 1 {
		result[i] = value
	}
	return result
}

// Returns a funding sample list of length `n = sum(counts)`, where
// each `values[i]` appears `counts[i]` times.
func GenerateFundingSamplesWithValues(
	values []int32,
	counts []uint32,
) (
	result []int32,
) {
	for i, count := range counts {
		result = append(result, GenerateConstantFundingPremiums(values[i], count)...)
	}

	return result
}

// Return a list of 60 funding premium samples, with 30 equal to
// 0.001% or 1000 in ppm, 15 equal to -0.1% and 15 equal to 0.1%.
func FundingSamples_Constant_0_001_Percent_Length_60_With_Noises() (
	result []int32,
) {
	result = make([]int32, 60)
	for i := 0; i < 30; i += 1 {
		result[i] = 1000
	}

	for i := 30; i < 45; i += 1 {
		result[i] = 100000
	}

	for i := 45; i < 60; i += 1 {
		result[i] = -100000
	}
	return result
}
