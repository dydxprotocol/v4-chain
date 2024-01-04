package types_test

import (
	sdkmath "cosmossdk.io/math"
	fmt "fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/stretchr/testify/require"
)

var (
	validAuthority = constants.AliceAccAddress.String()
)

func TestMsgSendFromModuleToAccount_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg types.MsgSendFromModuleToAccount
		err error
	}{
		"Valid": {
			msg: types.MsgSendFromModuleToAccount{
				Authority:        validAuthority,
				SenderModuleName: "gov",
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
			},
		},
		"Valid - module name has underscore": {
			msg: types.MsgSendFromModuleToAccount{
				Authority:        validAuthority,
				SenderModuleName: "insurance_fund",
				Recipient:        constants.AliceAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
			},
		},
		"Invalid authority": {
			msg: types.MsgSendFromModuleToAccount{
				Authority: "",
			},
			err: types.ErrInvalidAuthority,
		},
		"Invalid sender module name": {
			msg: types.MsgSendFromModuleToAccount{
				Authority:        validAuthority,
				SenderModuleName: "", // empty module name
				Recipient:        constants.BobAccAddress.String(),
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
			},
			err: types.ErrEmptyModuleName,
		},
		"Invalid recipient address": {
			msg: types.MsgSendFromModuleToAccount{
				Authority:        validAuthority,
				SenderModuleName: "bridge",
				Recipient:        "invalid_address",
				Coin:             sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
			},
			err: types.ErrInvalidAccountAddress,
		},
		"Invalid coin denom": {
			msg: types.MsgSendFromModuleToAccount{
				Authority:        validAuthority,
				SenderModuleName: "bridge",
				Recipient:        constants.CarlAccAddress.String(),
				Coin: sdk.Coin{
					Denom:  "7coin",
					Amount: sdkmath.NewInt(100),
				},
			},
			err: fmt.Errorf("invalid denom: %s", "7coin"),
		},
		"Invalid coin amount": {
			msg: types.MsgSendFromModuleToAccount{
				Authority:        validAuthority,
				SenderModuleName: "rewards",
				Recipient:        constants.CarlAccAddress.String(),
				Coin: sdk.Coin{
					Denom:  "random/coin",
					Amount: sdkmath.NewInt(-1),
				},
			},
			err: fmt.Errorf("negative coin amount: %v", -1),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
