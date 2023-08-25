package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestMsgServerSetClobPairStatus(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.ClobKeeper
	msgServer := keeper.NewMsgServerImpl(k)
	wrappedCtx := sdk.WrapSDKContext(ctx)

	tests := map[string]struct {
		testMsg      types.MsgSetClobPairStatus
		expectedResp *types.MsgSetClobPairStatusResponse
		expectedErr  string
	}{
		"Success": {
			testMsg: types.MsgSetClobPairStatus{
				Authority:      k.GetGovAuthority(),
				ClobPairId:     0,
				ClobPairStatus: types.ClobPair_STATUS_ACTIVE,
			},
			expectedResp: &types.MsgSetClobPairStatusResponse{},
		},
		"Failure: invalid authority": {
			testMsg: types.MsgSetClobPairStatus{
				Authority:      "12345",
				ClobPairId:     0,
				ClobPairStatus: types.ClobPair_STATUS_ACTIVE,
			},
			expectedErr: fmt.Sprintf(
				"invalid authority: expected %s, got %s",
				k.GetGovAuthority(),
				"12345",
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := msgServer.SetClobPairStatus(wrappedCtx, &tc.testMsg)

			// Assert msg server response.
			require.Equal(t, tc.expectedResp, resp)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
