package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateParams{}

func (msg *MsgUpdateParams) ValidateBasic() error {
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
	return msg.Params.Validate()
}
