package types

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgWithdrawFromSubaccount{}

// NewMsgWithdrawFromSubaccount constructs a `MsgWithdrawFromSubaccount` from an
// `x/subaccounts` subaccount sender, an `x/bank` account recipient, an asset ID,
// and a number of quantums.
func NewMsgWithdrawFromSubaccount(
	sender satypes.SubaccountId,
	recipient string,
	assetId uint32,
	quantums uint64,
) *MsgWithdrawFromSubaccount {
	return &MsgWithdrawFromSubaccount{
		Sender:    sender,
		Recipient: recipient,
		AssetId:   assetId,
		Quantums:  quantums,
	}
}

// ValidateBasic runs validation on the fields of a MsgWithdrawFromSubaccount.
func (msg *MsgWithdrawFromSubaccount) ValidateBasic() error {
	// Validate subaccount sender.
	if err := msg.Sender.Validate(); err != nil {
		return err
	}

	// Validate account recipient.
	_, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return ErrInvalidAccountAddress
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
