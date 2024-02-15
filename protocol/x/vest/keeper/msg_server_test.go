package keeper_test

import (
	"context"
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"

	"github.com/dydxprotocol/v4-chain/protocol/x/vest/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

var (
	GovAuthority = authtypes.NewModuleAddress(govtypes.ModuleName).String()
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, context.Context) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, k)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestMsgSetVestEntry(t *testing.T) {
	_, ms, ctx := setupMsgServer(t)

	testCases := []struct {
		name        string
		input       *types.MsgSetVestEntry
		expectedErr string
	}{
		{
			name: "valid params",
			input: &types.MsgSetVestEntry{
				Authority: GovAuthority,
				Entry:     TestValidEntry,
			},
			expectedErr: "",
		},
		{
			name: "invalid authority",
			input: &types.MsgSetVestEntry{
				Authority: "invalid",
				Entry:     TestValidEntry,
			},
			expectedErr: "invalid authority",
		},
		{
			name: "invalid params: invalid denom",
			input: &types.MsgSetVestEntry{
				Authority: GovAuthority,
				Entry: types.VestEntry{
					VesterAccount:   TestVesterAccount,
					TreasuryAccount: TestTreasuryAccount,
					Denom:           "invaliddenom!",
				},
			},
			expectedErr: "invalid denom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ms.SetVestEntry(ctx, tc.input)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgDeleteVestEntry(t *testing.T) {
	testCases := []struct {
		name        string
		input       *types.MsgDeleteVestEntry
		expectedErr string
	}{
		{
			name: "valid params",
			input: &types.MsgDeleteVestEntry{
				Authority:     GovAuthority,
				VesterAccount: TestVesterAccount,
			},
			expectedErr: "",
		},
		{
			name: "invalid authority",
			input: &types.MsgDeleteVestEntry{
				Authority:     "invalid",
				VesterAccount: TestVesterAccount,
			},
			expectedErr: "invalid authority",
		},
		{
			name: "delete non-existent entry",
			input: &types.MsgDeleteVestEntry{
				Authority:     GovAuthority,
				VesterAccount: "non_existent_vester",
			},
			expectedErr: "account is not associated with a vest entry",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, ms, goCtx := setupMsgServer(t)
			ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

			// Set up valid entry
			err := k.SetVestEntry(ctx, TestValidEntry)
			require.NoError(t, err)

			_, err = ms.DeleteVestEntry(goCtx, tc.input)
			if tc.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
