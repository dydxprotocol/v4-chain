package types

import (
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// Validate performs stateless validation on a `MatchPerpetualDeleveraging` object.
// It checks the following conditions to be true:
// - Validation for all subaccount Ids
// - length of fills to be greater than zero
// - For each fill, fill amount must be greater than zero
// - Subaccount ids in fills are all unique
// - Subaccount ids in fills cannot be the same as the liquidated subaccount id
func (match *MatchPerpetualDeleveraging) Validate() error {
	liquidatedSubaccountId := match.GetLiquidated()
	if err := liquidatedSubaccountId.Validate(); err != nil {
		return err
	}

	fills := match.GetFills()
	if len(fills) == 0 {
		return ErrEmptyDeleveragingFills
	}
	seenDeleveragedSubacountIds := map[satypes.SubaccountId]struct{}{}
	for _, fill := range fills {
		deleveragedSubaccountId := fill.GetDeleveraged()
		if err := deleveragedSubaccountId.Validate(); err != nil {
			return err
		}

		if deleveragedSubaccountId == liquidatedSubaccountId {
			return ErrDeleveragingAgainstSelf
		}
		if _, exists := seenDeleveragedSubacountIds[deleveragedSubaccountId]; exists {
			return ErrDuplicateDeleveragingFillSubaccounts
		}
		if fill.GetFillAmount() == 0 {
			return ErrZeroDeleveragingFillAmount
		}
		seenDeleveragedSubacountIds[deleveragedSubaccountId] = struct{}{}
	}
	return nil
}
