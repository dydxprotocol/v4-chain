package types

import (
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
	seenOffsettingSubacountIds := map[satypes.SubaccountId]struct{}{}
	for _, fill := range fills {
		offsettingSubaccountId := fill.GetOffsettingSubaccountId()
		if err := offsettingSubaccountId.Validate(); err != nil {
			return err
		}

		if offsettingSubaccountId == liquidatedSubaccountId {
			return ErrDeleveragingAgainstSelf
		}
		if _, exists := seenOffsettingSubacountIds[offsettingSubaccountId]; exists {
			return ErrDuplicateDeleveragingFillSubaccounts
		}
		if fill.GetFillAmount() == 0 {
			return ErrZeroDeleveragingFillAmount
		}
		seenOffsettingSubacountIds[offsettingSubaccountId] = struct{}{}
	}
	return nil
}
