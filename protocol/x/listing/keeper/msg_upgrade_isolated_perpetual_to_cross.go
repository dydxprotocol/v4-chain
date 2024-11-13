package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

var _ sdk.Msg = &types.MsgUpgradeIsolatedPerpetualToCross{}

/*
func (msg *types.MsgUpgradeIsolatedPerpetualToCross) ValidateBasic() error {
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
	// TODO Validation? Do we need to check if the PerpetualId is valid?
	return nil
}
*/
