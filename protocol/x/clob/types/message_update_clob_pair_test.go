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
	tests := map[string]struct {
		status        types.ClobPair_Status
		expectedError string
	}{
		"valid status": {
			status: types.ClobPair_STATUS_ACTIVE,
		},
		"invalid unsupported status": {
			status:        types.ClobPair_STATUS_UNSPECIFIED,
			expectedError: "has unsupported status",
		},
		"invalid negative out of bounds status": {
			status:        -1,
			expectedError: "has unsupported status",
		},
		"invalid positive out of bounds status": {
			status:        100,
			expectedError: "has unsupported status",
		},
		// TODO add more
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			clobPair := constants.ClobPair_Btc
			clobPair.Status = tc.status
			msg := types.MsgUpdateClobPair{
				ClobPair: &clobPair,
			}
			err := msg.ValidateBasic()

			if tc.expectedError != "" {
				require.ErrorContains(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
