package types

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgDepositToSubaccount{}

// NewMsgDepositToSubaccount constructs a `MsgDepositToSubaccount` from an
// `x/bank` account sender, an `x/subaccounts` subaccount sender, an asset ID,
// and a number of quantums.
func NewMsgDepositToSubaccount(
	sender string,
	recipient satypes.SubaccountId,
	assetId uint32,
	quantums uint64,
) *MsgDepositToSubaccount {
	return &MsgDepositToSubaccount{
		Sender:    sender,
		Recipient: recipient,
		AssetId:   assetId,
		Quantums:  quantums,
	}
}

// ValidateBasic runs validation on the fields of a MsgDepositToSubaccount.
func (msg *MsgDepositToSubaccount) ValidateBasic() error {
	// Validate account sender.
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidAccountAddress
	}

	// Validate subaccount recipient.
	if err := msg.Recipient.Validate(); err != nil {
		return err
	}

	// Validate that asset is TDai.
	if msg.AssetId != assettypes.AssetTDai.Id {
		return ErrNonTDaiAssetTransferNotImplemented
	}

	// Validate that quantums is not zero.
	if msg.Quantums == lib.ZeroUint64 {
		return ErrInvalidTransferAmount
	}

	return nil
}
