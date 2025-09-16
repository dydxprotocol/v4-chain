package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgUpdateLeverage{}

// ValidateBasic validates the message's subaccount ID and leverage data.
func (msg *MsgUpdateLeverage) ValidateBasic() error {
	if msg.SubaccountId == nil {
		return errorsmod.Wrap(ErrInvalidAddress, "subaccount ID cannot be nil")
	}

	if msg.SubaccountId.Owner == "" {
		return errorsmod.Wrap(ErrInvalidAddress, "subaccount owner cannot be empty")
	}

	if _, err := sdk.AccAddressFromBech32(msg.SubaccountId.Owner); err != nil {
		return errorsmod.Wrap(
			ErrInvalidAddress,
			fmt.Sprintf(
				"subaccount owner '%s' must be a valid bech32 address, but got error '%v'",
				msg.SubaccountId.Owner,
				err.Error(),
			),
		)
	}

	// Validate that leverage entries are not empty
	if len(msg.PerpetualLeverage) == 0 {
		return errorsmod.Wrap(ErrInvalidLeverage, "perpetual leverage entries cannot be empty")
	}

	// Validate leverage values are positive and perpetual IDs are unique
	perpetualIds := make(map[uint32]bool)
	for _, entry := range msg.PerpetualLeverage {
		if entry == nil {
			return errorsmod.Wrap(ErrInvalidLeverage, "leverage entry cannot be nil")
		}

		if entry.Leverage == 0 {
			return errorsmod.Wrap(
				ErrInvalidLeverage,
				fmt.Sprintf("leverage for perpetual %d cannot be zero", entry.PerpetualId),
			)
		}

		if perpetualIds[entry.PerpetualId] {
			return errorsmod.Wrap(
				ErrInvalidLeverage,
				fmt.Sprintf("duplicate perpetual ID %d", entry.PerpetualId),
			)
		}
		perpetualIds[entry.PerpetualId] = true
	}

	return nil
}
