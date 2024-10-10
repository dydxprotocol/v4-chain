package types_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

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
					InsuranceFundFeePpm: 5_000,
					ValidatorFeePpm:     200_000,
					LiquidityFeePpm:     800_000,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm: 0,
					},
					MaxCumulativeInsuranceFundDelta: uint64(1_000_000_000_000),
				},
			},
			expectedError: "0 is not a valid BankruptcyAdjustmentPpm",
		},
		"invalid authority": {
			msg: types.MsgUpdateLiquidationsConfig{
				LiquidationsConfig: constants.LiquidationsConfig_No_Limit,
			},
			expectedError: "Authority is invalid",
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
