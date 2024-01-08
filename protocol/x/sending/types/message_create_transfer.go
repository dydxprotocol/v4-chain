package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

var _ sdk.Msg = &MsgCreateTransfer{}

func NewMsgCreateTransfer(transfer *Transfer) *MsgCreateTransfer {
	return &MsgCreateTransfer{
		Transfer: transfer,
	}
}

func (msg *MsgCreateTransfer) ValidateBasic() error {
	err := msg.Transfer.Sender.Validate()
	if err != nil {
		return err
	}

	err = msg.Transfer.Recipient.Validate()
	if err != nil {
		return err
	}

	if msg.Transfer.Sender == msg.Transfer.Recipient {
		return errorsmod.Wrapf(ErrSenderSameAsRecipient, "Sender is the same as recipient (%s)", &msg.Transfer.Sender)
	}

	if msg.Transfer.AssetId != assettypes.AssetUsdc.Id {
		return ErrNonUsdcAssetTransferNotImplemented
	}

	if msg.Transfer.Amount == uint64(0) {
		return ErrInvalidTransferAmount
	}

	return nil
}
