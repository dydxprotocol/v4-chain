package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateBlockLimitsConfig{}

// ValidateBasic validates the message's BlockLimitsConfig. Returns an error if the authority
// is empty or if the BlockLimitsConfig is invalid.
func (msg *MsgUpdateBlockLimitsConfig) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(
			ErrInvalidAuthority,
			fmt.Sprintf(
				"authority '%s' must be a valid bech32 address, but got error '%v'",
				msg.Authority,
				err.Error(),
			),
		)
	}

	return msg.BlockLimitsConfig.Validate()
}
