package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "perpetuals"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	TransientStoreKey = "tmp_" + ModuleName
)

// State
const (
	// PerpetualKeyPrefix is the prefix to retrieve all Perpetual
	PerpetualKeyPrefix = "Perp:"

	// PremiumVotesKey is the key to retrieve `PremiumStore` object
	// that represents existing premium sample votes during the current
	// `funding-sample` epoch.
	PremiumVotesKey = "PremVotes"

	// PremiumSamplesKey is the key to retrieve `PremiumStore` object
	// that represents existing premium samples during the current
	// `funding-tick` epoch.
	PremiumSamplesKey = "PremSamples"

	// LiquidityTierKeyPrefix is the prefix to retrieve all `LiquidityTier`s.
	LiquidityTierKeyPrefix = "LiqTier:"

	// ParamsKey is the key to retrieve all params for the module.
	ParamsKey = "Params"

	// UpdatedOIKeyPrefix is the key to retrieve the updated OI for the module.
	UpdatedOIKeyPrefix = "UpdatedOI"

	// NextPerpetualIDKey is the key to retrieve the next perpetual id to be used
	NextPerpetualIDKey = "NextPerpetualID"
)

// Module Accounts
const (
	// InsuranceFundName defines the root string for the insurance fund account address
	InsuranceFundName = "insurance_fund"
)
