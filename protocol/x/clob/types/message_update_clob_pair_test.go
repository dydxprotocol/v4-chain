package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateClobPair_GetSigners(t *testing.T) {
	msg := types.MsgUpdateClobPair{
		Authority: constants.AliceAccAddress.String(),
	}
	require.Equal(t, []sdk.AccAddress{constants.AliceAccAddress}, msg.GetSigners())
}

func TestMsgUpdateClobPair_ValidateBasic(t *testing.T) {
	tests := []struct {
		desc        string
		clobPair    types.ClobPair
		expectedErr string
	}{
		{
			desc: "Invalid Metadata (SpotClobMetadata)",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_SpotClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "is not a perpetual CLOB",
		},
		{
			desc: "UNSPECIFIED Status",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_PAUSED,
			},
			expectedErr: "has unsupported status",
		},
		{
			desc: "invalid negative status integer",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           -1,
			},
			expectedErr: "has unsupported status",
		},
		{
			desc: "invalid positive status integer",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  1,
				Status:           100,
			},
			expectedErr: "has unsupported status",
		},
		{
			desc: "StepBaseQuantums <= 0",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 0,
				SubticksPerTick:  1,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "StepBaseQuantums must be > 0.",
		},
		{
			desc: "SubticksPerTick <= 0",
			clobPair: types.ClobPair{
				Metadata:         &types.ClobPair_PerpetualClobMetadata{},
				StepBaseQuantums: 1,
				SubticksPerTick:  0,
				Status:           types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: "SubticksPerTick must be > 0",
		},
		{
			desc: "Valid ClobPair",
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
				ClobPair: tc.clobPair,
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
