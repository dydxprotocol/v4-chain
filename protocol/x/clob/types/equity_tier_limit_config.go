package types

import (
	errorsmod "cosmossdk.io/errors"
	"sort"
)

const (
	MaxShortTermOrdersForEquityTier = 10_000_000
	MaxStatefulOrdersForEquityTier  = 10_000_000
)

// Validate validates each individual EquityTierLimit.
// It returns an error if any of the equity tier limits fail the following validations:
//   - `Limit > MaxShortTermOrdersForEquityTier` for short term order equity tier limits.
//   - `Limit > MaxStatefulOrdersPerEquityTier` for stateful order equity tier limits.
//   - There are multiple equity tier limits for the same `UsdTncRequired` in `ShortTermOrderEquityTiers`,
//     or `StatefulOrderEquityTiers`.
func (lc EquityTierLimitConfiguration) Validate() error {
	if err := (equityTierLimits)(lc.ShortTermOrderEquityTiers).validate(
		"ShortTermOrderEquityTiers",
		MaxShortTermOrdersForEquityTier,
	); err != nil {
		return err
	}
	if err := (equityTierLimits)(lc.StatefulOrderEquityTiers).validate(
		"StatefulOrderEquityTiers",
		MaxStatefulOrdersForEquityTier,
	); err != nil {
		return err
	}
	return nil
}

// Initialize ensures the fields are ordered by the application requirements:
//   - ShortTermOrderEquityTiers by UsdTncRequired in ascending order.
//   - StatefulOrderEquityTiers by UsdTncRequired in ascending order.
func (lc EquityTierLimitConfiguration) Initialize() {
	(equityTierLimits)(lc.ShortTermOrderEquityTiers).sortByUsdTncRequiredAsc()
	(equityTierLimits)(lc.StatefulOrderEquityTiers).sortByUsdTncRequiredAsc()
}

type equityTierLimits []EquityTierLimit

// Sorts the list by UsdTncRequired in ascending order.
func (l equityTierLimits) sortByUsdTncRequiredAsc() {
	sort.Slice(l, func(i, j int) bool {
		return l[i].UsdTncRequired.Cmp(l[j].UsdTncRequired) <= 0
	})
}

func (l equityTierLimits) validate(field string, maxOrders uint32) error {
	// Work on a copy to not modify the original slice.
	sortSlice := make([]EquityTierLimit, len(l))
	copy(sortSlice, l)
	sort.Slice(sortSlice, func(i, j int) bool {
		return sortSlice[i].UsdTncRequired.Cmp(sortSlice[j].UsdTncRequired) <= 0
	})

	for i, limit := range sortSlice {
		if err := limit.validate(field, maxOrders); err != nil {
			return err
		}

		if i > 0 && l[i-1].UsdTncRequired.Cmp(l[i].UsdTncRequired) >= 0 {
			return errorsmod.Wrapf(
				ErrInvalidEquityTierLimitConfig,
				"Expected %s equity tier UsdTncRequired to be strictly ascending. "+
					"Found %+v and %+v out of order",
				field,
				l[i-1],
				l[i],
			)
		}
	}
	return nil
}

func (l EquityTierLimit) validate(field string, maxOrders uint32) error {
	if l.Limit > maxOrders {
		return errorsmod.Wrapf(
			ErrInvalidEquityTierLimitConfig,
			"%d is not a valid Limit for %s equity tier limit %+v",
			l.Limit,
			field,
			l,
		)
	}
	if l.UsdTncRequired.IsNil() || l.UsdTncRequired.Sign() < 0 {
		return errorsmod.Wrapf(
			ErrInvalidEquityTierLimitConfig,
			"%d is not a valid UsdTncRequired for %s equity tier limit %+v",
			l.UsdTncRequired.BigInt(),
			field,
			l,
		)
	}
	return nil
}
