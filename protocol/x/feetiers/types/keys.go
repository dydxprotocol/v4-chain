package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "feetiers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

// State
const (
	// PerpetualFeeParamsKey defines the key for the PerpetualFeeParams
	PerpetualFeeParamsKey = "PerpParams"

	// MarketFeeDiscountPrefix is the prefix for storing market fee discount
	MarketFeeDiscountPrefix = "MarketFeeDiscount:"

	// StakingTierKeyPrefix is the prefix for staking tier store
	StakingTierKeyPrefix = "StakingTier:"
)

// StakingTierKey returns the store key for a staking tier
func StakingTierKey(tierName string) []byte {
	return append([]byte(StakingTierKeyPrefix), []byte(tierName)...)
}
