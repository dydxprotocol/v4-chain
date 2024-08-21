package keeper_test

import (
	"math/big"
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
					NewVaultDepositAmount:  dtypes.NewIntFromBigInt(big.NewInt(100_000_000)),
					MainVaultDepositAmount: dtypes.NewIntFromBigInt(big.NewInt(0)),
					NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
				},
			},
		},
		"Failure - Invalid Authority": {
			msg: &types.MsgSetListingVaultDepositParams{
				Authority: constants.AliceAccAddress.String(),
				Params: types.ListingVaultDepositParams{
					NewVaultDepositAmount:  dtypes.NewIntFromBigInt(big.NewInt(100_000_000)),
					MainVaultDepositAmount: dtypes.NewIntFromBigInt(big.NewInt(0)),
					NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
				},
			},
			expectedErr: "invalid authority",
		},
		"Failure - Empty authority": {
			msg: &types.MsgSetListingVaultDepositParams{
				Params: types.ListingVaultDepositParams{
					NewVaultDepositAmount:  dtypes.NewIntFromBigInt(big.NewInt(100_000_000)),
					MainVaultDepositAmount: dtypes.NewIntFromBigInt(big.NewInt(0)),
					NumBlocksToLockShares:  30 * 24 * 3600, // 30 days
				},
			},
			expectedErr: "invalid authority",
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
