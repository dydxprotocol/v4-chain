package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// ValidateProcessProposerMatchesEvents performs basic stateless validation on ProcessProposerMatchesEvents.
// It returns an error if:
//   - Block height does not equal current block height.
//   - Any of the fields have duplicate OrderIds. Note that this is currently invalid since
//     stateful order replacements are not enabled.
func (ppme *ProcessProposerMatchesEvents) ValidateProcessProposerMatchesEvents(
	ctx sdk.Context,
) error {
	if ctx.BlockHeight() != int64(ppme.BlockHeight) {
		return fmt.Errorf(
			"block height %d for ProcessProposerMatchesEvents does not equal current block height %d",
			ppme.BlockHeight,
			ctx.BlockHeight(),
		)
	}

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
