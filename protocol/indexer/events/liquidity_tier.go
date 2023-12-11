package events

// NewLiquidityTierUpsertEvent creates a LiquidityTierUpsertEvent representing
// upsert of a liquidity tier.
func NewLiquidityTierUpsertEvent(
	id uint32,
	name string,
	initialMarginPpm uint32,
	maintenanceFractionPpm uint32,
) *LiquidityTierUpsertEventV1 {
	return &LiquidityTierUpsertEventV1{
		Id:                     id,
		Name:                   name,
		InitialMarginPpm:       initialMarginPpm,
		MaintenanceFractionPpm: maintenanceFractionPpm,
	}
}
