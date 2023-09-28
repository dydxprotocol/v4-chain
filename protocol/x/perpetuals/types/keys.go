package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "perpetuals"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// PerpetualKeyPrefix is the prefix to retrieve all Perpetual
	PerpetualKeyPrefix = "perpetual/"

	// PerpetualKeyPrefix is the prefix to retrieve all funding premium samples
	FundingSamplesKeyPrefix = "funding_samples/"

	// PremiumVotesKey is the key to retrieve `PremiumStore` object
	// that represents existing premium sample votes during the current
	// `funding-sample` epoch.
	PremiumVotesKey = "premium_votes"

	// PremiumSamplesKey is the key to retrieve `PremiumStore` object
	// that represents existing premium samples during the current
	// `funding-tick` epoch.
	PremiumSamplesKey = "premium_samples"

	// LiquidityTierKeyPrefix is the prefix to retrieve all `LiquidityTier`s.
	LiquidityTierKeyPrefix = "liquidity_tier/"
)

// Params
const (
	// FundingRateClampFactorKey is the key to retrieve funding rate clamp factor
	// in parts-per-million
	// |R| <= clamp_factor * (initial margin - maintenance margin)
	FundingRateClampFactorPpmKey = "funding_rate_clamp_factor_ppm"

	// PremiumVoteClampFactorPpmKey is the key to retrieve premium vote clamp factor
	// in parts-per-million
	// |V| <= clamp_factor * (initial margin - maintenance margin)
	PremiumVoteClampFactorPpmKey = "premium_vote_clamp_factor_ppm"

	// MinNumVotesPerSampleKey is the key to retrieve `min_num_votes_per_sample`,
	// the minimum number of votes needed to calculate a premium sample.
	MinNumVotesPerSampleKey = "min_num_votes_per_sample"
)
