package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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

// GetSigners specifies that the sender of the message must sign.
func (msg *MsgWithdrawFromSubaccount) GetSigners() []sdk.AccAddress {
	// Get address of sender subaccount's account address.
	sender, err := sdk.AccAddressFromBech32(msg.Sender.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
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

	// Validate that asset is USDC.
	if msg.AssetId != assettypes.AssetUsdc.Id {
		return ErrNonUsdcAssetTransferNotImplemented
	}

	// Validate that quantums is not zero.
	if msg.Quantums == lib.ZeroUint64 {
		return ErrInvalidTransferAmount
	}

	return nil
}
