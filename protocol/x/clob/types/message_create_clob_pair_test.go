package types_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateClobPair_ValidateBasic(t *testing.T) {
	tests := []struct {
		desc        string
		msg         types.MsgCreateClobPair
		expectedErr string
	}{
		{
			desc: "Invalid Metadata (SpotClobMetadata)",
			msg: types.MsgCreateClobPair{
				Authority: lib.GovModuleAddress.String(),
				ClobPair: types.ClobPair{
					Metadata:         &types.ClobPair_SpotClobMetadata{},
					StepBaseQuantums: 1,
					SubticksPerTick:  1,
					Status:           types.ClobPair_STATUS_ACTIVE,
				},
			},
			expectedErr: "is not a perpetual CLOB",
		},
		{
			desc: "Empty authority",
			msg: types.MsgCreateClobPair{
				Authority: "",
				ClobPair: types.ClobPair{
					Metadata:         &types.ClobPair_PerpetualClobMetadata{},
					StepBaseQuantums: 1,
					SubticksPerTick:  1,
					Status:           types.ClobPair_STATUS_ACTIVE,
				},
			},
			expectedErr: "authority cannot be empty: Authority is invalid",
		},
		{
			desc: "Unsupported Status",
			msg: types.MsgCreateClobPair{
				Authority: lib.GovModuleAddress.String(),
				ClobPair: types.ClobPair{
					Metadata:         &types.ClobPair_PerpetualClobMetadata{},
					StepBaseQuantums: 1,
					SubticksPerTick:  1,
					Status:           types.ClobPair_STATUS_PAUSED,
				},
			},
			expectedErr: "has unsupported status",
		},
		{
			desc: "StepBaseQuantums <= 0",
			msg: types.MsgCreateClobPair{
				Authority: lib.GovModuleAddress.String(),
				ClobPair: types.ClobPair{
					Metadata:         &types.ClobPair_PerpetualClobMetadata{},
					StepBaseQuantums: 0,
					SubticksPerTick:  1,
					Status:           types.ClobPair_STATUS_ACTIVE,
				},
			},
			expectedErr: "StepBaseQuantums must be > 0.",
		},
		{
			desc: "SubticksPerTick <= 0",
			msg: types.MsgCreateClobPair{
				Authority: lib.GovModuleAddress.String(),
				ClobPair: types.ClobPair{
					Metadata:         &types.ClobPair_PerpetualClobMetadata{},
					StepBaseQuantums: 1,
					SubticksPerTick:  0,
					Status:           types.ClobPair_STATUS_ACTIVE,
				},
			},
			expectedErr: "SubticksPerTick must be > 0",
		},
		{
			desc: "Valid ClobPair",
			msg: types.MsgCreateClobPair{
				Authority: lib.GovModuleAddress.String(),
				ClobPair: types.ClobPair{
					Metadata:         &types.ClobPair_PerpetualClobMetadata{},
					StepBaseQuantums: 1,
					SubticksPerTick:  1,
					Status:           types.ClobPair_STATUS_ACTIVE,
				},
			},
			expectedErr: "",
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}
