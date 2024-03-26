package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var _ sdk.Msg = &MsgBatchCancel{}

// NewMsgBatchCancel constructs a MsgBatchCancel.
func NewMsgBatchCancel(
	subaccountId satypes.SubaccountId,
	cancelBatch []OrderBatch,
	goodTilBlock uint32,
) *MsgBatchCancel {
	return &MsgBatchCancel{
		SubaccountId:     subaccountId,
		ShortTermCancels: cancelBatch,
		GoodTilBlock:     goodTilBlock,
	}
}

// ValidateBasic performs stateless validation for the `MsgBatchCancel` msg.
func (msg *MsgBatchCancel) ValidateBasic() (err error) {
	subaccountId := msg.GetSubaccountId()
	if err := subaccountId.Validate(); err != nil {
		return err
	}

	cancelBatches := msg.GetShortTermCancels()
	if len(cancelBatches) == 0 {
		return errorsmod.Wrapf(
			ErrInvalidBatchCancel,
			"Batch cancel cannot have zero orders specified.",
		)
	}
	totalNumberCancels := 0
	seenClobPairIds := map[uint32]struct{}{}
	for _, cancelBatch := range cancelBatches {
		// Check for duplicate clob pair ids across all cancel batches
		clobPairId := cancelBatch.GetClobPairId()
		if _, exists := seenClobPairIds[clobPairId]; exists {
			return errorsmod.Wrapf(
				ErrInvalidBatchCancel,
				"Batch cancel cannot have 2 order batches with the same clob pair id: %+v",
				clobPairId,
			)
		}
		seenClobPairIds[clobPairId] = struct{}{}

		numClientIds := len(cancelBatch.GetClientIds())
		if numClientIds == 0 {
			return errorsmod.Wrapf(
				ErrInvalidBatchCancel,
				"Order Batch cannot have zero client ids.",
			)
		}
		totalNumberCancels += numClientIds

		// Check for duplicate clob client ids across this cancel batch
		seenClientIds := map[uint32]struct{}{}
		for _, clientId := range cancelBatch.GetClientIds() {
			if _, seen := seenClientIds[clientId]; seen {
				return errorsmod.Wrapf(
					ErrInvalidBatchCancel,
					"Batch Cancel Clob pair id %+v has duplicate client ids: %+v",
					cancelBatch.GetClobPairId(),
					clientId,
				)
			}
			seenClientIds[clientId] = struct{}{}
		}
	}
	if uint32(totalNumberCancels) > MaxMsgBatchCancelBatchSize {
		return errorsmod.Wrapf(
			ErrInvalidBatchCancel,
			"Batch cancel cannot have over %+v orders. Order count: %+v",
			MaxMsgBatchCancelBatchSize,
			totalNumberCancels,
		)
	}
	return nil
}
