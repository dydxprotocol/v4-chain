package types

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/types"
)

var _ types.Msg = &MsgWithdrawFromVault{}

// ValidateBasic performs stateless validation on a MsgWithdrawFromVault.
func (msg *MsgWithdrawFromVault) ValidateBasic() error {
	// Note: msg signer must be the owner of the subaccount.
	// This is enforced by the following notatino on the msg proto:
	//    option (cosmos.msg.v1.signer) = "subaccount_id"

	// Validate subaccount to withdraw to.
	if err := msg.SubaccountId.Validate(); err != nil {
		return err
	}

	// Validate that `shares` is positive.
	if msg.Shares.NumShares.Sign() <= 0 {
		return errors.Wrap(ErrInvalidWithdrawalAmount, "shares must be strictly positive")
	}

	return nil
}
