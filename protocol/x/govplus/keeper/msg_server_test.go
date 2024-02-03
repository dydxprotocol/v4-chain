package keeper_test

import (
	"context"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/govplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/govplus/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, context.Context) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.GovPlusKeeper

	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, k)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestSlashValidator(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	testCases := []struct {
		name      string
		input     *types.MsgSlashValidator
		expErr    bool
		expErrMsg string
	}{
		{
			name: "invalid authority",
			input: &types.MsgSlashValidator{
				Authority: "invalid",
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
		{
			name: "bad address",
			input: &types.MsgSlashValidator{
				Authority:        lib.GovModuleAddress.String(),
				ValidatorAddress: "asdfasdfasdfaf",
			},
			expErr:    true,
			expErrMsg: "Could not convert validator consensus address",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.SlashValidator(ctx, tc.input)
			if tc.expErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expErrMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
