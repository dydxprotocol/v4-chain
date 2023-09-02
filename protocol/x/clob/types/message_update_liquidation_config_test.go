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
		liquidationsConfig types.LiquidationsConfig
		expectedError      error
	}{
		"valid": {
			liquidationsConfig: constants.LiquidationsConfig_No_Limit,
		},
		"invalid": {
			liquidationsConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm: 5_000,
				FillablePriceConfig: types.FillablePriceConfig{
					BankruptcyAdjustmentPpm: 0,
				},
				PositionBlockLimits:   constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits: constants.SubaccountBlockLimits_No_Limit,
			},
			expectedError: types.ErrInvalidLiquidationsConfig,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := types.MsgUpdateLiquidationsConfig{
				LiquidationsConfig: tc.liquidationsConfig,
			}
			err := msg.ValidateBasic()

			if tc.expectedError != nil {
				require.ErrorContains(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
