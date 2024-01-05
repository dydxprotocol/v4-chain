package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

var (
	// validAuthority is a valid bech32 address.
	validAuthority = constants.AliceAccAddress.String()
)

func TestMsgUpdateClobPair_ValidateBasic(t *testing.T) {
	tests := []struct {
		desc        string
		authority   string
		clobPair    types.ClobPair
		expectedErr string
	}{
		{
			desc:      "Invalid Metadata (SpotClobMetadata)",
			authority: validAuthority,
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_SpotClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "is not a perpetual CLOB",
		},
		{
			desc:      "UNSPECIFIED Status",
			authority: validAuthority,
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_PAUSED,
			},
			expectedErr: "has unsupported status",
		},
		{
			desc:      "invalid negative status integer",
			authority: validAuthority,
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           -1,
			},
			expectedErr: "has unsupported status",
		},
		{
			desc:      "invalid positive status integer",
			authority: validAuthority,
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           100,
			},
			expectedErr: "has unsupported status",
		},
		{
			desc:      "StepBaseQuantums <= 0",
			authority: validAuthority,
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 0,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "StepBaseQuantums must be > 0.",
		},
		{
			desc:      "SubticksPerTick <= 0",
			authority: validAuthority,
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  0,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "SubticksPerTick must be > 0",
		},
		{
			desc: "Invalid authority",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "Authority is invalid",
		},
		{
			desc:      "Valid ClobPair",
			authority: validAuthority,
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			msg := types.MsgUpdateClobPair{
				Authority: tc.authority,
				ClobPair:  tc.clobPair,
			}
			err := msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
