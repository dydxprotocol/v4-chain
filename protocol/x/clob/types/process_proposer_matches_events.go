package types

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// ValidateProcessProposerMatchesEvents performs basic stateless validation on ProcessProposerMatchesEvents.
// It returns an error if:
// - Any of the fields have duplicate OrderIds.
func (ppme *ProcessProposerMatchesEvents) ValidateProcessProposerMatchesEvents() error {
	if lib.ContainsDuplicates(ppme.PlacedLongTermOrderIds) {
		return fmt.Errorf(
			"ProcessProposerMatchesEvents contains duplicate PlacedLongTermOrderIds: %+v",
			ppme.PlacedLongTermOrderIds,
		)
	}
	if lib.ContainsDuplicates(ppme.ExpiredStatefulOrderIds) {
		return fmt.Errorf(
			"ProcessProposerMatchesEvents contains duplicate ExpiredStatefulOrderIds: %+v",
			ppme.ExpiredStatefulOrderIds,
		)
	}
	if lib.ContainsDuplicates(ppme.OrderIdsFilledInLastBlock) {
		return fmt.Errorf(
			"ProcessProposerMatchesEvents contains duplicate OrderIdsFilledInLastBlock: %+v",
			ppme.OrderIdsFilledInLastBlock,
		)
	}
	if lib.ContainsDuplicates(ppme.PlacedStatefulCancellationOrderIds) {
		return fmt.Errorf(
			"ProcessProposerMatchesEvents contains duplicate PlacedStatefulCancellationOrderIds: %+v",
			ppme.PlacedStatefulCancellationOrderIds,
		)
	}
	if lib.ContainsDuplicates(ppme.RemovedStatefulOrderIds) {
		return fmt.Errorf(
			"ProcessProposerMatchesEvents contains duplicate RemovedStatefulOrderIds: %+v",
			ppme.RemovedStatefulOrderIds,
		)
	}
	if lib.ContainsDuplicates(ppme.ConditionalOrderIdsTriggeredInLastBlock) {
		return fmt.Errorf(
			"ProcessProposerMatchesEvents contains duplicate ConditionalOrderIdsTriggeredInLastBlock: %+v",
			ppme.ConditionalOrderIdsTriggeredInLastBlock,
		)
	}
	if lib.ContainsDuplicates(ppme.PlacedConditionalOrderIds) {
		return fmt.Errorf(
			"ProcessProposerMatchesEvents contains duplicate PlacedConditionalOrderIds: %+v",
			ppme.PlacedConditionalOrderIds,
		)
	}
	return nil
}
