package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestOwnerShares(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Request.
		req *vaulttypes.QueryOwnerSharesRequest
		// Vault ID.
		vaultId vaulttypes.VaultId
		// Owner shares.
		ownerShares map[string]*big.Int

		/* --- Expectations --- */
		expectedOwnerShares []*vaulttypes.OwnerShare
		expectedErr         string
	}{
		"Success": {
			req: &vaulttypes.QueryOwnerSharesRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 0,
			},
			vaultId: constants.Vault_Clob0,
			ownerShares: map[string]*big.Int{
				constants.Alice_Num0.Owner: big.NewInt(100),
				constants.Bob_Num0.Owner:   big.NewInt(200),
			},
			expectedOwnerShares: []*vaulttypes.OwnerShare{
				{
					Owner: constants.Alice_Num0.Owner,
					Shares: &vaulttypes.NumShares{
						NumShares: dtypes.NewInt(100),
					},
				},
				{
					Owner: constants.Bob_Num0.Owner,
					Shares: &vaulttypes.NumShares{
						NumShares: dtypes.NewInt(200),
					},
				},
			},
		},
		"Error: vault not found": {
			req: &vaulttypes.QueryOwnerSharesRequest{
				Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
				Number: 1,
			},
			vaultId: constants.Vault_Clob0,
			ownerShares: map[string]*big.Int{
				constants.Alice_Num0.Owner: big.NewInt(100),
				constants.Bob_Num0.Owner:   big.NewInt(200),
			},
			expectedErr: "vault not found",
		},
		"Error: nil request": {
			req:     nil,
			vaultId: constants.Vault_Clob0,
			ownerShares: map[string]*big.Int{
				constants.Alice_Num0.Owner: big.NewInt(100),
				constants.Bob_Num0.Owner:   big.NewInt(200),
			},
			expectedErr: "invalid request",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			// Set total and owner shares.
			totalShares := big.NewInt(0)
			for owner, shares := range tc.ownerShares {
				err := k.SetOwnerShares(ctx, tc.vaultId, owner, vaulttypes.BigIntToNumShares(shares))
				require.NoError(t, err)
				totalShares.Add(totalShares, shares)
			}

			// Set total shares.
			err := k.SetTotalShares(ctx, tc.vaultId, vaulttypes.BigIntToNumShares(totalShares))
			require.NoError(t, err)

			// Check Vault query response is as expected.
			response, err := k.OwnerShares(ctx, tc.req)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.ElementsMatch(
					t,
					tc.expectedOwnerShares,
					response.OwnerShares,
				)
				require.Equal(
					t,
					&query.PageResponse{
						NextKey: nil,
						Total:   uint64(len(tc.expectedOwnerShares)),
					},
					response.Pagination,
				)
			}
		})
	}
}
