package types

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

const (
	// PerpetualKeyPrefix is the prefix to retrieve all Perpetual
	PerpetualKeyPrefix = "Perpetual/value/"
	// PerpetualKeyPrefix is the prefix to retrieve all funding premium samples
	FundingSamplesKeyPrefix = "Perpetual/funding/"
	// FundingRateClampFactorKey is the key to retrieve funding rate clamp factor
	// in parts-per-million
	// |R| <= clamp_factor * (initial margin - maintenance margin)
	FundingRateClampFactorPpmKey = "Perpetual/funding_rate_clamp_factor_ppm"
	// PremiumVoteClampFactorPpmKey is the key to retrieve premium vote clamp factor
	// in parts-per-million
	// |V| <= clamp_factor * (initial margin - maintenance margin)
	PremiumVoteClampFactorPpmKey = "Perpetual/premium_vote_clamp_factor_ppm"
	// PremiumVotesKey is the key to retrieve `PremiumStore` object
	// that represents existing premium sample votes during the current
	// `funding-sample` epoch.
	PremiumVotesKey = "Perpetual/premium_sample_votes"
	// PremiumSamplesKey is the key to retrieve `PremiumStore` object
	// that represents existing premium samples during the current
	// `funding-tick` epoch.
	PremiumSamplesKey = "Perpetual/premium_samples"
	// MinNumVotesPerSampleKey is the key to retrieve `min_num_votes_per_sample`,
	// the minimum number of votes needed to calculate a premium sample.
	MinNumVotesPerSampleKey = "Perpetual/min_num_votes_per_sample"
	// LiquidityTierKeyPrefix is the prefix to retrieve all `LiquidityTier`s.
	LiquidityTierKeyPrefix = "Perpetual/liquidity_tier/"
	// NumLiquidityTiersKey is the prefix to retrieve the cardinality of `LiquidityTier`s.
	NumLiquidityTiersKey = "Perpetual/liquidity_tier/num"
)

// PerpetualKey returns the store key to retrieve a Perpetual by its id.
func PerpetualKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytesForState(id)
}

// LiquidityTierKey returns the store key to retrieve a LiquidityTier by its id.
func LiquidityTierKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytesForState(id)
}
