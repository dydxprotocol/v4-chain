package types

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	// Clamp factor for 8-hour funding rate is by default 600%.
	DefaultFundingRateClampFactorPpm = 6 * lib.OneMillion
	// Clamp factor for premium vote is by default 6_000%.
	DefaultPremiumVoteClampFactorPpm = 60 * lib.OneMillion
	// Minimum number of votes per sample is by default 15.
	DefaultMinNumVotesPerSample = 15

	// Maximum default funding rate magnitude is 100%.
	MaxDefaultFundingPpmAbs = lib.OneMillion

	// Liquidity-tier related constants
	MaxInitialMarginPpm       = lib.OneMillion
	MaxMaintenanceFractionPpm = lib.OneMillion
)

// DefaultGenesis returns the default Perpetual genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Perpetuals:     []Perpetual{},
		LiquidityTiers: []LiquidityTier{},
		Params: Params{
			FundingRateClampFactorPpm: DefaultFundingRateClampFactorPpm,
			PremiumVoteClampFactorPpm: DefaultPremiumVoteClampFactorPpm,
			MinNumVotesPerSample:      DefaultMinNumVotesPerSample,
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate parameters.
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate perpetuals
	// 1. keys are unique
	// 2. `Ticker` is non-empty
	perpKeyMap := make(map[uint32]struct{})

	for _, perp := range gs.Perpetuals {
		if _, exists := perpKeyMap[perp.Params.Id]; exists {
			return fmt.Errorf("duplicated perpetual id")
		}
		perpKeyMap[perp.Params.Id] = struct{}{}

		if len(perp.Params.Ticker) == 0 {
			return ErrTickerEmptyString
		}
	}

	// Validate liquidity tiers.
	// 1. keys are unique.
	// 2. IDs are sequential.
	// 3. initial margin does not exceed its max value.
	// 4. maintenance margin does not exceed its max value.
	// 5. base position notional is not zero.
	liquidityTierKeyMap := make(map[uint32]struct{})
	expectedLiquidityTierId := uint32(0)
	for _, liquidityTier := range gs.LiquidityTiers {
		if _, exists := liquidityTierKeyMap[liquidityTier.Id]; exists {
			return fmt.Errorf("duplicated liquidity tier id: %d", liquidityTier.Id)
		}
		liquidityTierKeyMap[liquidityTier.Id] = struct{}{}

		if liquidityTier.Id != expectedLiquidityTierId {
			return fmt.Errorf("found a gap in liquidity tier id. Expected %d, got %d", expectedLiquidityTierId, liquidityTier.Id)
		}
		expectedLiquidityTierId = expectedLiquidityTierId + 1

		if err := liquidityTier.Validate(); err != nil {
			return err
		}
	}
	return nil
}
