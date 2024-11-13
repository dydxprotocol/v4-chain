package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	types "github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

var (
	// validAuthority is a valid bech32 address.
	validAuthority = constants.AliceAccAddress.String()
)

func TestMsgUpgradeIsolatedPerpetualToCross_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg         types.MsgUpgradeIsolatedPerpetualToCross
		expectedErr string
	}{
		"Success": {
			msg: types.MsgUpgradeIsolatedPerpetualToCross{
				Authority:   validAuthority,
				PerpetualId: 1,
			},
		},
		"Failure: Invalid authority": {
			msg: types.MsgUpgradeIsolatedPerpetualToCross{
				Authority: "",
			},
			expectedErr: "Authority is invalid",
		},
	}

	for name, _ := range tests {
		t.Run(name, func(t *testing.T) {
			/*
				err := tc.msg.ValidateBasic()
				if tc.expectedErr == "" {
					require.NoError(t, err)
				} else {
					require.ErrorContains(t, err, tc.expectedErr)
				}
			*/
		})
	}
}

/*
var _ sdk.Msg = &types.MsgUpgradeIsolatedPerpetualToCross{}

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
