package events

// NewLiquidityTierUpsertEvent creates a LiquidityTierUpsertEvent representing
// upsert of a liquidity tier.
func NewLiquidityTierUpsertEvent(
	id uint32,
	name string,
	initialMarginPpm uint32,
	maintenanceFractionPpm uint32,
	openInterestLowerCap uint64,
	openInterestUpperCap uint64,
) *LiquidityTierUpsertEventV2 {
	return &LiquidityTierUpsertEventV2{
		Id:                     id,
		Name:                   name,
		InitialMarginPpm:       initialMarginPpm,
		MaintenanceFractionPpm: maintenanceFractionPpm,
		OpenInterestLowerCap:   openInterestLowerCap,
		OpenInterestUpperCap:   openInterestUpperCap,
	}
}
