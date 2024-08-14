package sim_helpers

import (
	"math"

	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// All generated genesis parameters should be defined here.

// Clob genesis parameters.
var (
	MinValidClobPairs = MinNumPerpetuals.Valid
	MaxValidClobPairs = MaxNumPerpetuals.Valid * 2

	MinQuantumConversionExponent = GenesisParameters[int]{
		Reasonable: -9,
		Valid:      -25,
	}
	MaxQuantumConversionExponent = GenesisParameters[int]{
		Reasonable: 0,
		Valid:      25,
	}

	MinStepBaseQuantums = GenesisParameters[int]{
		Reasonable: 10,
		Valid:      1,
	}
	MaxStepBaseQuantums = GenesisParameters[int]{
		Reasonable: 100_000,
		Valid:      math.MaxUint32,
	}

	MaxOrderBaseQuantums = GenesisParameters[int]{
		Reasonable: 10_000_000,
		Valid:      math.MaxUint32,
	}

	MinSubticksPerTick = GenesisParameters[int]{
		Reasonable: 10,
		Valid:      1,
	}
	MaxSubticksPerTick = GenesisParameters[int]{
		Reasonable: 100_000,
		Valid:      math.MaxUint32,
	}

	MinPositionNotionalBuckets = []int{
		1, // min
		1_000_000,
		100_000_000,
		1_000_000_000_000, // $1,000,000
	}
	SubaccountBlockLimitsBuckets = []int{
		1_000_000_000, // $1,000
		10_000_000_000,
		100_000_000_000,
		1_000_000_000_000, // $1,000,000
	}
	PpmSkewedTowardLowerBuckets = []int{
		1,
		1_000,
		10_000,
		100_000,
		1_000_000,
	}
	PpmSkewedTowardLargerBuckets = []int{
		1,
		500_000,
		750_000,
		900_000,
		1_000_000,
	}
)

// Prices genesis parameters.
var (
	MinMarkets = GenesisParameters[int]{
		Reasonable: 10,
		Valid:      1,
	}
	MaxMarkets = GenesisParameters[int]{
		Reasonable: 200,
		Valid:      1000,
	}

	MinMinExchangesPerMarket = GenesisParameters[int]{
		Reasonable: 2,
		Valid:      1,
	}
	MaxMinExchangesPerMarket = GenesisParameters[int]{
		Reasonable: 3,
		Valid:      6,
	}

	MinMarketExponent = GenesisParameters[int]{
		Reasonable: -10,
		Valid:      -10,
	}
	MaxMarketExponent = GenesisParameters[int]{
		Reasonable: -5,
		Valid:      -1,
	}
	MinMarketPrice = GenesisParameters[int]{
		Reasonable: 10_000,
		Valid:      1,
	}
	MaxMarketPrice = GenesisParameters[int]{
		Reasonable: 500_000_000_000,
		Valid:      math.MaxInt,
	}

	MinExchangesPerMarket = 1
)

// Perpetuals genesis parameters.
var (
	MinNumPerpetuals = GenesisParameters[int]{
		Reasonable: MinMarkets.Reasonable * 2,
		Valid:      MinMarkets.Valid * 2,
	}
	MaxNumPerpetuals = GenesisParameters[int]{
		Reasonable: MaxMarkets.Reasonable * 2,
		Valid:      MaxMarkets.Valid * 2,
	}

	MinNumLiquidityTiers = GenesisParameters[int]{
		Reasonable: 1,
		Valid:      2,
	}
	MaxNumLiquidityTiers = GenesisParameters[int]{
		Reasonable: 4,
		Valid:      MaxNumPerpetuals.Valid,
	}

	MinFundingRateClampFactorPpm = GenesisParameters[int]{
		Reasonable: 4_000_000, // 400%
		Valid:      1_000_000, // 100%
	}
	MaxFundingRateClampFactorPpm = GenesisParameters[int]{
		Reasonable: 8_000_000,  // 400%
		Valid:      12_000_000, // 1200%
	}

	MinPremiumVoteClampFactorPpm = GenesisParameters[int]{
		Reasonable: 40_000_000, // 4_000%
		Valid:      10_000_000, // 1_000%
	}
	MaxPremiumVoteClampFactorPpm = GenesisParameters[int]{
		Reasonable: 80_000_000,  // 8_000%
		Valid:      120_000_000, // 12_000%
	}

	MinMinNumVotesPerSample = GenesisParameters[int]{
		Reasonable: 10,
		Valid:      0,
	}
	MaxMinNumVotesPerSample = GenesisParameters[int]{
		Reasonable: 25,
		Valid:      80,
	}

	MinAtomicResolution = GenesisParameters[int]{
		Reasonable: -10,
		Valid:      -10,
	}
	MaxAtomicResolution = GenesisParameters[int]{
		Reasonable: 0,
		Valid:      10,
	}

	DefaultFundingPpmAbsBuckets = []int{
		0, // min
		100,
		1_000,
		10_000,
		100_000,
		int(perptypes.MaxDefaultFundingPpmAbs), // max
	}

	InitialMarginBuckets = []int{
		0, // min
		100,
		1_000,
		10_000,
		100_000,
		int(perptypes.MaxInitialMarginPpm) + 1, // max
	}
)

// Subaccounts genesis parameters.
var (
	MinNumSubaccount = 1
	MaxNumSubaccount = 128
)
