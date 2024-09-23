package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMegavaultAllOwnerShares(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Request.
		req *vaulttypes.QueryMegavaultAllOwnerSharesRequest
		// Owner shares.
		ownerShares map[string]*big.Int

		/* --- Expectations --- */
		expectedOwnerShares []*vaulttypes.OwnerShare
		expectedErr         string
	}{
		"Success": {
			req: &vaulttypes.QueryMegavaultAllOwnerSharesRequest{},
			ownerShares: map[string]*big.Int{
				constants.Alice_Num0.Owner: big.NewInt(100),
				constants.Bob_Num0.Owner:   big.NewInt(200),
			},
			expectedOwnerShares: []*vaulttypes.OwnerShare{
				{
					Owner: constants.Alice_Num0.Owner,
					Shares: vaulttypes.NumShares{
						NumShares: dtypes.NewInt(100),
					},
				},
				{
					Owner: constants.Bob_Num0.Owner,
					Shares: vaulttypes.NumShares{
						NumShares: dtypes.NewInt(200),
					},
				},
			},
		},
		"Success: no owners": {
			req:         &vaulttypes.QueryMegavaultAllOwnerSharesRequest{},
			ownerShares: map[string]*big.Int{},
		},
		"Error: nil request": {
			req:         nil,
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
				err := k.SetOwnerShares(ctx, owner, vaulttypes.BigIntToNumShares(shares))
				require.NoError(t, err)
				totalShares.Add(totalShares, shares)
			}
			err := k.SetTotalShares(ctx, vaulttypes.BigIntToNumShares(totalShares))
			require.NoError(t, err)

			// Check OwnerShares query response is as expected.
			response, err := k.MegavaultAllOwnerShares(ctx, tc.req)
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

func TestMegavaultOwnerShares(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Request.
		req *vaulttypes.QueryMegavaultOwnerSharesRequest
		// Owner address.
		ownerAddress string
		// Owner shares.
		ownerShares *big.Int
		// Share unlocks.
		shareUnlocks []vaulttypes.ShareUnlock
		// Total shares.
		totalShares *big.Int
		// Megavault equity.
		megavaultEquity *big.Int

		/* --- Expectations --- */
		expectedOwnerEquity             *big.Int
		expectedOwnerWithdrawableEquity *big.Int
		expectedErr                     string
	}{
		"Success with zero unlock": {
			req: &vaulttypes.QueryMegavaultOwnerSharesRequest{
				Address: constants.AliceAccAddress.String(),
			},
			ownerAddress:                    constants.AliceAccAddress.String(),
			ownerShares:                     big.NewInt(100),
			totalShares:                     big.NewInt(5_000),
			megavaultEquity:                 big.NewInt(2_000_000),
			expectedOwnerEquity:             big.NewInt(40_000),
			expectedOwnerWithdrawableEquity: big.NewInt(40_000),
		},
		"Success with one unlock": {
			req: &vaulttypes.QueryMegavaultOwnerSharesRequest{
				Address: constants.AliceAccAddress.String(),
			},
			ownerAddress: constants.AliceAccAddress.String(),
			ownerShares:  big.NewInt(100),
			shareUnlocks: []vaulttypes.ShareUnlock{
				{
					Shares: vaulttypes.NumShares{
						NumShares: dtypes.NewInt(17),
					},
					UnlockBlockHeight: 123,
				},
			},
			totalShares:         big.NewInt(5_000),
			megavaultEquity:     big.NewInt(2_000_000),
			expectedOwnerEquity: big.NewInt(40_000),
			// 40_000 * (100 - 17) / 100 = 33_200
			expectedOwnerWithdrawableEquity: big.NewInt(33_200),
		},
		"Success with two unlocks": {
			req: &vaulttypes.QueryMegavaultOwnerSharesRequest{
				Address: constants.BobAccAddress.String(),
			},
			ownerAddress: constants.BobAccAddress.String(),
			ownerShares:  big.NewInt(47_123),
			shareUnlocks: []vaulttypes.ShareUnlock{
				{
					Shares: vaulttypes.NumShares{
						NumShares: dtypes.NewInt(1_234),
					},
					UnlockBlockHeight: 905,
				},
				{
					Shares: vaulttypes.NumShares{
						NumShares: dtypes.NewInt(4_444),
					},
					UnlockBlockHeight: 1023,
				},
			},
			totalShares:     big.NewInt(358_791_341),
			megavaultEquity: big.NewInt(753_582_314_912),
			// 753_582_314_912 * 47_123 / 358_791_341 ~= 98_974_126
			expectedOwnerEquity: big.NewInt(98_974_126),
			// 98_974_126 * (47_123 - 5_678) / 47_123 ~= 87_048_419
			expectedOwnerWithdrawableEquity: big.NewInt(87_048_419),
		},
		"Success: owner has 0 shares": {
			req: &vaulttypes.QueryMegavaultOwnerSharesRequest{
				Address: constants.AliceAccAddress.String(),
			},
			ownerAddress:    constants.AliceAccAddress.String(),
			ownerShares:     big.NewInt(0),
			totalShares:     big.NewInt(5_000),
			megavaultEquity: big.NewInt(2_000_000),
		},
		"Error: owner not found": {
			req: &vaulttypes.QueryMegavaultOwnerSharesRequest{
				Address: constants.BobAccAddress.String(),
			},
			ownerAddress:    constants.AliceAccAddress.String(),
			ownerShares:     big.NewInt(100),
			shareUnlocks:    []vaulttypes.ShareUnlock{},
			totalShares:     big.NewInt(5_000),
			megavaultEquity: big.NewInt(2_000_000),
			expectedErr:     "owner not found",
		},
		"Error: nil request": {
			req:             nil,
			ownerAddress:    constants.AliceAccAddress.String(),
			ownerShares:     big.NewInt(100),
			shareUnlocks:    []vaulttypes.ShareUnlock{},
			totalShares:     big.NewInt(5_000),
			megavaultEquity: big.NewInt(2_000_000),
			expectedErr:     "invalid request",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize megavault main vault with its equity.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  constants.Usdc.Id,
										Quantums: dtypes.NewIntFromBigInt(tc.megavaultEquity),
									},
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			err := k.SetTotalShares(ctx, vaulttypes.BigIntToNumShares(tc.totalShares))
			require.NoError(t, err)
			err = k.SetOwnerShares(
				ctx,
				tc.ownerAddress,
				vaulttypes.BigIntToNumShares(tc.ownerShares),
			)
			require.NoError(t, err)
			err = k.SetOwnerShareUnlocks(
				ctx,
				tc.ownerAddress,
				vaulttypes.OwnerShareUnlocks{
					OwnerAddress: tc.ownerAddress,
					ShareUnlocks: tc.shareUnlocks,
				},
			)
			require.NoError(t, err)

			// Check OwnerShares query response is as expected.
			res, err := k.MegavaultOwnerShares(ctx, tc.req)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(
					t,
					vaulttypes.QueryMegavaultOwnerSharesResponse{
						Address:            tc.req.Address,
						Shares:             vaulttypes.BigIntToNumShares(tc.ownerShares),
						ShareUnlocks:       tc.shareUnlocks,
						Equity:             dtypes.NewIntFromBigInt(tc.expectedOwnerEquity),
						WithdrawableEquity: dtypes.NewIntFromBigInt(tc.expectedOwnerWithdrawableEquity),
					},
					*res,
				)
			}
		})
	}
}
