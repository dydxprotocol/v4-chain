package types_test

import (
	fmt "fmt"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSendFromAccountToAccount_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg types.MsgSendFromAccountToAccount
		err error
	}{
		"Valid": {
			msg: types.MsgSendFromAccountToAccount{
				Authority: validAuthority,
				Sender:    constants.BobAccAddress.String(),
				Recipient: constants.AliceAccAddress.String(),
				Coin:      sdk.NewCoin("adv4tnt", sdkmath.NewInt(1)),
			},
		},
		"Valid - larger amount": {
			msg: types.MsgSendFromAccountToAccount{
				Authority: validAuthority,
				Sender:    constants.CarlAccAddress.String(),
				Recipient: constants.DaveAccAddress.String(),
				Coin:      sdk.NewCoin("adv4tnt", sdkmath.NewInt(100000)),
			},
		},
		"Invalid authority": {
			msg: types.MsgSendFromAccountToAccount{
				Authority: "",
			},
			err: types.ErrInvalidAuthority,
		},
		"Invalid sender address": {
			msg: types.MsgSendFromAccountToAccount{
				Authority: validAuthority,
				Sender:    "invalid_address",
				Recipient: constants.AliceAccAddress.String(),
				Coin:      sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
			},
			err: types.ErrInvalidAccountAddress,
		},
		"Invalid recipient address": {
			msg: types.MsgSendFromAccountToAccount{
				Authority: validAuthority,
				Sender:    constants.BobAccAddress.String(),
				Recipient: "invalid_address",
				Coin:      sdk.NewCoin("adv4tnt", sdkmath.NewInt(100)),
			},
			err: types.ErrInvalidAccountAddress,
		},
		"Invalid coin denom": {
			msg: types.MsgSendFromAccountToAccount{
				Authority: validAuthority,
				Sender:    constants.BobAccAddress.String(),
				Recipient: constants.CarlAccAddress.String(),
				Coin: sdk.Coin{
					Denom:  "7coin",
					Amount: sdkmath.NewInt(100),
				},
			},
			err: fmt.Errorf("invalid denom: %s", "7coin"),
		},
		"Invalid coin amount": {
			msg: types.MsgSendFromAccountToAccount{
				Authority: validAuthority,
				Sender:    constants.BobAccAddress.String(),
				Recipient: constants.CarlAccAddress.String(),
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
