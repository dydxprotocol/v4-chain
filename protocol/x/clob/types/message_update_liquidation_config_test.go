package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateLiquidationsConfig_GetSigners(t *testing.T) {
	msg := types.MsgUpdateLiquidationsConfig{
		Authority: constants.AliceAccAddress.String(),
	}
	require.Equal(t, []sdk.AccAddress{constants.AliceAccAddress}, msg.GetSigners())
}

func TestMsgUpdateLiquidationsConfig_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg           types.MsgUpdateLiquidationsConfig
		expectedError string
	}{
		"valid": {
			msg: types.MsgUpdateLiquidationsConfig{
				Authority:          constants.AliceAccAddress.String(),
				LiquidationsConfig: constants.LiquidationsConfig_No_Limit,
			},
		},
		"invalid liquidations config": {
			msg: types.MsgUpdateLiquidationsConfig{
				Authority: constants.AliceAccAddress.String(),
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 5_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm: 0,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
				},
			},
			expectedError: "0 is not a valid BankruptcyAdjustmentPpm",
		},
		"invalid authority": {
			msg: types.MsgUpdateLiquidationsConfig{
				LiquidationsConfig: constants.LiquidationsConfig_No_Limit,
			},
			expectedError: "authority cannot be empty",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()

			if tc.expectedError != "" {
				require.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
