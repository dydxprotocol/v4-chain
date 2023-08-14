package types

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	DefaultOpenInterest = 0
	// Clamp factor for 8-hour funding rate is by default 600%.
	DefaultFundingRateClampFactorPpm = 6 * lib.OneMillion
	// Clamp factor for premium vote is by default 6_000%.
	DefaultPremiumVoteClampFactorPpm = 60 * lib.OneMillion
	// Minimum number of votes per sample is by default 15.
	DefaultMinNumVotesPerSample = 15

	MaxDefaultFundingPpmAbs   = lib.OneMillion
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
	// 2. IDs are sequential
	// 3. `Ticker` is non-empty
	perpKeyMap := make(map[uint32]struct{})
	expectedPerpId := uint32(0)

	for _, perp := range gs.Perpetuals {
		if _, exists := perpKeyMap[perp.Id]; exists {
			return fmt.Errorf("duplicated perpetual id")
		}
		perpKeyMap[perp.Id] = struct{}{}

		if perp.Id != expectedPerpId {
			return fmt.Errorf("found a gap in perpetual id")
		}
		expectedPerpId = expectedPerpId + 1

		if len(perp.Ticker) == 0 {
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
			return fmt.Errorf("duplicated liquidity tier id")
		}
		liquidityTierKeyMap[liquidityTier.Id] = struct{}{}

		if liquidityTier.Id != expectedLiquidityTierId {
			return fmt.Errorf("found a gap in liquidity tier id")
		}
		expectedLiquidityTierId = expectedLiquidityTierId + 1

		if err := liquidityTier.Validate(); err != nil {
			return err
		}
	}
	return nil
}
