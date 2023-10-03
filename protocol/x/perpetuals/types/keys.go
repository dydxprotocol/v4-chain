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

	// ParamsKey is the key to retrieve all params for the module.
	ParamsKey = "params"
)
