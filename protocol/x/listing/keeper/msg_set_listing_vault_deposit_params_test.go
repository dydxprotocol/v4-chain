package keeper_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	"github.com/stretchr/testify/require"
)

func TestSetListingVaultDepositParams(t *testing.T) {
	tests := map[string]struct {
		// Msg.
		msg *types.MsgSetListingVaultDepositParams
		// Expected error
		expectedErr string
	}{
		"Success - Set vault deposit params": {
			msg: &types.MsgSetListingVaultDepositParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.ListingVaultDepositParams{
					NewVaultDepositAmount:  dtypes.NewIntFromUint64(10_000),
					MainVaultDepositAmount: dtypes.NewIntFromUint64(0),
					NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
				},
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetListingVaultDepositParams{
				Authority: constants.AliceAccAddress.String(),
				Params: types.ListingVaultDepositParams{
					NewVaultDepositAmount:  dtypes.NewIntFromUint64(10_000),
					MainVaultDepositAmount: dtypes.NewIntFromUint64(0),
					NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty authority": {
			msg: &types.MsgSetListingVaultDepositParams{
				Params: types.ListingVaultDepositParams{
					NewVaultDepositAmount:  dtypes.NewIntFromUint64(10_000),
					MainVaultDepositAmount: dtypes.NewIntFromUint64(0),
					NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Invalid deposit amount": {
			msg: &types.MsgSetListingVaultDepositParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.ListingVaultDepositParams{
					NewVaultDepositAmount:  dtypes.NewIntFromUint64(0),
					MainVaultDepositAmount: dtypes.NewIntFromUint64(0),
					NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
				},
			},
			expectedErr: "invalid vault deposit amount",
		},
		"Failure - Invalid num blocks to lock shares": {
			msg: &types.MsgSetListingVaultDepositParams{
				Authority: lib.GovModuleAddress.String(),
				Params: types.ListingVaultDepositParams{
					NewVaultDepositAmount:  dtypes.NewIntFromUint64(10_000),
					MainVaultDepositAmount: dtypes.NewIntFromUint64(0),
					NumBlocksToLockShares:  0,
				},
			},
			expectedErr: "invalid number of blocks to lock shares",
		},
	}

	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				tApp := testapp.NewTestAppBuilder(t).Build()
				ctx := tApp.InitChain()
				k := tApp.App.ListingKeeper
				ms := keeper.NewMsgServerImpl(k)
				_, err := ms.SetListingVaultDepositParams(ctx, tc.msg)
				if tc.expectedErr != "" {
					require.Error(t, err)
					require.Contains(t, err.Error(), tc.expectedErr)
				} else {
					params := k.GetListingVaultDepositParams(ctx)
					require.Equal(t, tc.msg.Params, params)
				}
			},
		)
	}
}
