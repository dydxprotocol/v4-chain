package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateLiquidationsConfig{}

// ValidateBasic validates the message's LiquidationConfig. Returns an error if the authority
// is empty or if the LiquidationsConfig is invalid.
func (msg *MsgUpdateLiquidationsConfig) ValidateBasic() error {
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

	return msg.LiquidationsConfig.Validate()
}
