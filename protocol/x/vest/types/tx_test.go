package types_test

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	validAuthority = authtypes.NewModuleAddress("authority")
	validVestEntry = types.VestEntry{
		VesterAccount:   "vester_account",
		TreasuryAccount: "treasury_account",
		Denom:           "denom",
		StartTime:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:         time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
	}
)

func TestMsgSetVestEntry_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg    *types.MsgSetVestEntry
		expErr error
	}{
		"valid": {
			msg: types.NewMsgSetVestEntry(validAuthority.String(), validVestEntry),
		},
		"invalid authority": {
			msg: &types.MsgSetVestEntry{
				Authority: "invalid",
			},
			expErr: types.ErrInvalidAuthority,
		},
		"invalid entry": {
			msg:    types.NewMsgSetVestEntry(validAuthority.String(), types.VestEntry{}),
			expErr: types.ErrInvalidVesterAccount,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgDeleteVestEntry_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg    *types.MsgDeleteVestEntry
		expErr error
	}{
		"valid": {
			msg: types.NewMsgDeleteVestEntry(validAuthority.String(), "vester_account"),
		},
		"invalid authority": {
			msg:    types.NewMsgDeleteVestEntry("invalid", "vester_account"),
			expErr: types.ErrInvalidAuthority,
		},
		"invalid vester account": {
			msg:    types.NewMsgDeleteVestEntry(validAuthority.String(), ""),
			expErr: types.ErrInvalidVesterAccount,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
