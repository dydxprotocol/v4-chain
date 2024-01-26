package types

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgSlashValidator) ValidateBasic() error {
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

	_, err := sdk.ConsAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return errorsmod.Wrap(
			ErrValidatorAddress,
			fmt.Sprintf(
				"got error when converting consensus address %s from bech32 '%v'",
				msg.ValidatorAddress,
				err.Error(),
			),
		)
	}

	if msg.PowerAtInfractionHeight.BigInt().Cmp(big.NewInt(0)) != 1 {
		return ErrInvalidPowerAtInfractionHeight
	}

	if msg.SlashFactor.IsNegative() || msg.SlashFactor.GT(math.LegacyNewDec(1)) {
		return ErrInvalidSlashFactor
	}

	return nil
}
