package types_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSetLimitParams_ValidateBasic(t *testing.T) {
	validAuthority := "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu"
	validDenom := "uatom"
	validPeriod := time.Hour
	validBaselineMinimum := dtypes.NewIntFromBigInt(big.NewInt(1000000))
	validBaselineTvlPpm := uint32(500000)

	validLimitParams := types.LimitParams{
		Denom: validDenom,
		Limiters: []types.Limiter{
			{
				Period:          validPeriod,
				BaselineMinimum: validBaselineMinimum,
				BaselineTvlPpm:  validBaselineTvlPpm,
			},
		},
	}

	tests := []struct {
		name string
		msg  types.MsgSetLimitParams
		err  string
	}{
		{
			name: "valid message",
			msg: types.MsgSetLimitParams{
				Authority:   validAuthority,
				LimitParams: validLimitParams,
			},
			err: "",
		},
		{
			name: "invalid authority",
			msg: types.MsgSetLimitParams{
				Authority:   "invalid_address",
				LimitParams: validLimitParams,
			},
			err: types.ErrInvalidAuthority.Error(),
		},
		{
			name: "empty limit params",
			msg: types.MsgSetLimitParams{
				Authority:   validAuthority,
				LimitParams: types.LimitParams{},
			},
			err: "invalid denom",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.err != "" {
				require.ErrorContains(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgSetLimitParams_GetSigners(t *testing.T) {
	tests := []struct {
		name      string
		authority string
		expected  []sdk.AccAddress
	}{
		{
			name:      "valid authority",
			authority: "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu",
			expected:  []sdk.AccAddress{sdk.MustAccAddressFromBech32("cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu")},
		},
		{
			name:      "invalid authority",
			authority: "invalid_address",
			expected:  []sdk.AccAddress{nil},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			msg := types.MsgSetLimitParams{
				Authority: tc.authority,
			}
			signers := msg.GetSigners()
			require.Equal(t, tc.expected, signers)
		})
	}
}
