package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	types "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePerpetual_GetSigners(t *testing.T) {
	msg := types.MsgCreatePerpetual{
		Authority: constants.AliceAccAddress.String(),
	}
	require.Equal(t, []sdk.AccAddress{constants.AliceAccAddress}, msg.GetSigners())
}

func TestMsgCreatePerpetual_ValidateBasic(t *testing.T) {
	tests := []struct {
		desc        string
		msg         types.MsgCreatePerpetual
		expectedErr string
	}{
		{
			desc:        "Empty authority",
			msg:         types.MsgCreatePerpetual{},
			expectedErr: "authority cannot be empty",
		},
		{
			desc: "Empty ticker",
			msg: types.MsgCreatePerpetual{
				Authority: "test",
				Params: types.PerpetualParams{
					Ticker: "",
				},
			},
			expectedErr: "Ticker must be non-empty string",
		},
		{
			desc: "DefaultFundingPpm >= MaxDefaultFundingPpmAbs",
			msg: types.MsgCreatePerpetual{
				Authority: "test",
				Params: types.PerpetualParams{
					Ticker:            "test",
					DefaultFundingPpm: 100_000_000,
				},
			},
			expectedErr: "DefaultFundingPpm magnitude exceeds maximum value",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			require.ErrorContains(t, err, tc.expectedErr)
		})
	}
}
