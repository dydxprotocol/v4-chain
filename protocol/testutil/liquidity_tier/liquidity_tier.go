package liquidity_tier

import (
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

type LtModifierOption func(cp *perptypes.LiquidityTier)

func WithId(id uint32) LtModifierOption {
	return func(lt *perptypes.LiquidityTier) {
		lt.Id = id
	}
}

func WithName(name string) LtModifierOption {
	return func(lt *perptypes.LiquidityTier) {
		lt.Name = name
	}
}

func WithInitialMarginPpm(initialMarginPpm uint32) LtModifierOption {
	return func(lt *perptypes.LiquidityTier) {
		lt.InitialMarginPpm = initialMarginPpm
	}
}

func WithMaintenanceFractionPpm(maintenanceFractionPpm uint32) LtModifierOption {
	return func(lt *perptypes.LiquidityTier) {
		lt.MaintenanceFractionPpm = maintenanceFractionPpm
	}
}

func WithImpactNotional(impactNotional uint64) LtModifierOption {
	return func(lt *perptypes.LiquidityTier) {
		lt.ImpactNotional = impactNotional
	}
}

// GenerateLiquidityTier returns a `LiquidityTier` object set to default values.
// Passing in `LtModifierOption` methods alters the value of the `LiquidityTier` returned.
// It will start with the default, valid `LiquidityTier` value defined within the method
// and make the requested modifications before returning the object.
//
// Example usage:
// `GenerateLiquidityTier(WithId(7))`
// This will start with the default `LiquidityTier` object defined within the method and
// return the newly-created object after overriding the values of `Id` to 7.
func GenerateLiquidityTier(optionalModifications ...LtModifierOption) *perptypes.LiquidityTier {
	lt := &perptypes.LiquidityTier{
		Id:                     0,
		Name:                   "Large-Cap",
		InitialMarginPpm:       1_000_000,
		MaintenanceFractionPpm: 1_000_000,
		ImpactNotional:         500_000_000,
	}

	for _, opt := range optionalModifications {
		opt(lt)
	}

	return lt
}
