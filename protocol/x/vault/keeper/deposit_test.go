package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestMintShares(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Existing megavault equity.
		equity *big.Int
		// Existing megavault TotalShares.
		totalShares *big.Int
		// Owner that deposits.
		owner string
		// Existing owner shares.
		ownerShares *big.Int
		// Quote quantums to deposit.
		quantumsToDeposit *big.Int

		/* --- Expectations --- */
		// Expected TotalShares after minting.
		expectedTotalShares *big.Int
		// Expected OwnerShares after minting.
		expectedOwnerShares *big.Int
		// Expected error.
		expectedErr error
	}{
		"Equity 0, TotalShares 0, OwnerShares 0, Deposit 1000": {
			equity:            big.NewInt(0),
			totalShares:       big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			ownerShares:       big.NewInt(0),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `1_000` shares.
			expectedTotalShares: big.NewInt(1_000),
			expectedOwnerShares: big.NewInt(1_000),
		},
		"Equity 0, TotalShares 0, OwnerShares non-existent, Deposit 12345654321": {
			equity:            big.NewInt(0),
			totalShares:       big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(12_345_654_321),
			// Should mint `12_345_654_321` shares.
			expectedTotalShares: big.NewInt(12_345_654_321),
			expectedOwnerShares: big.NewInt(12_345_654_321),
		},
		"Equity 1000, TotalShares 0, OwnerShares non-existent, Deposit 500": {
			equity:            big.NewInt(1_000),
			totalShares:       big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(500),
			// Should mint `500` shares.
			expectedTotalShares: big.NewInt(500),
			expectedOwnerShares: big.NewInt(500),
		},
		"Equity 4000, TotalShares 5000, OwnerShares 2500, Deposit 1000": {
			equity:            big.NewInt(4_000),
			totalShares:       big.NewInt(5_000),
			owner:             constants.AliceAccAddress.String(),
			ownerShares:       big.NewInt(2_500),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `1_250` shares.
			expectedTotalShares: big.NewInt(6_250),
			expectedOwnerShares: big.NewInt(3_750),
		},
		"Equity 1_000_000, TotalShares 2_000, OwnerShares 1, Deposit 1_000": {
			equity:            big.NewInt(1_000_000),
			totalShares:       big.NewInt(2_000),
			owner:             constants.BobAccAddress.String(),
			ownerShares:       big.NewInt(1),
			quantumsToDeposit: big.NewInt(1_000),
			// Should mint `2` shares.
			expectedTotalShares: big.NewInt(2_002),
			expectedOwnerShares: big.NewInt(3),
		},
		"Equity 8000, TotalShares 4000, OwnerShares 101, Deposit 455": {
			equity:            big.NewInt(8_000),
			totalShares:       big.NewInt(4_000),
			owner:             constants.CarlAccAddress.String(),
			ownerShares:       big.NewInt(101),
			quantumsToDeposit: big.NewInt(455),
			// Should mint `227.5` shares, round down to 227.
			expectedTotalShares: big.NewInt(4_227),
			expectedOwnerShares: big.NewInt(328),
		},
		"Equity 123456, TotalShares 654321, OwnerShares 0, Deposit 123456789": {
			equity:            big.NewInt(123_456),
			totalShares:       big.NewInt(654_321),
			owner:             constants.DaveAccAddress.String(),
			quantumsToDeposit: big.NewInt(123_456_789),
			// Should mint `654_325_181.727` shares, round down to 654_325_181.
			expectedTotalShares: big.NewInt(654_979_502),
			expectedOwnerShares: big.NewInt(654_325_181),
		},
		"Equity 1000000, TotalShares 1000, OwnerShares 0, Deposit 9_900": {
			equity:            big.NewInt(1_000_000),
			totalShares:       big.NewInt(1_000),
			owner:             constants.DaveAccAddress.String(),
			quantumsToDeposit: big.NewInt(9_900),
			// Should mint `9_900 * 1_000 / 1_000_000` shares, round down to 9.
			expectedTotalShares: big.NewInt(1_009),
			expectedOwnerShares: big.NewInt(9),
		},
		"Equity -1, TotalShares 10, Deposit 1": {
			equity:            big.NewInt(-1),
			totalShares:       big.NewInt(10),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(1),
			expectedErr:       vaulttypes.ErrNonPositiveEquity,
		},
		"Equity 1, TotalShares 1, Deposit 0": {
			equity:            big.NewInt(1),
			totalShares:       big.NewInt(1),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(0),
			expectedErr:       vaulttypes.ErrInvalidQuoteQuantums,
		},
		"Equity 0, TotalShares 0, Deposit -1": {
			equity:            big.NewInt(0),
			totalShares:       big.NewInt(0),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(-1),
			expectedErr:       vaulttypes.ErrInvalidQuoteQuantums,
		},
		"Equity 1000, TotalShares 1, Deposit 100": {
			equity:            big.NewInt(1_000),
			totalShares:       big.NewInt(1),
			owner:             constants.AliceAccAddress.String(),
			quantumsToDeposit: big.NewInt(100),
			expectedErr:       vaulttypes.ErrZeroSharesToMint,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize vault with its existing equity.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										0,
										tc.equity,
									),
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Set existing total shares.
			err := tApp.App.VaultKeeper.SetTotalShares(
				ctx,
				vaulttypes.BigIntToNumShares(tc.totalShares),
			)
			require.NoError(t, err)
			// Set existing owner shares if specified.
			if tc.ownerShares != nil {
				err := tApp.App.VaultKeeper.SetOwnerShares(
					ctx,
					tc.owner,
					vaulttypes.BigIntToNumShares(tc.ownerShares),
				)
				require.NoError(t, err)
			}

			// Mint shares.
			mintedShares, err := tApp.App.VaultKeeper.MintShares(
				ctx,
				tc.owner,
				tc.quantumsToDeposit,
			)
			if tc.expectedErr != nil {
				// Check that error is as expected.
				require.ErrorContains(t, err, tc.expectedErr.Error())
				require.Nil(t, mintedShares)
				// Check that TotalShares is unchanged.
				totalShares := tApp.App.VaultKeeper.GetTotalShares(ctx)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(tc.totalShares),
					totalShares,
				)
				// Check that OwnerShares is unchanged.
				ownerShares, _ := tApp.App.VaultKeeper.GetOwnerShares(ctx, tc.owner)
				require.Equal(t, vaulttypes.BigIntToNumShares(tc.ownerShares), ownerShares)
			} else {
				require.NoError(t, err)
				// Check that minted shares is as expected.
				require.Equal(
					t,
					new(big.Int).Sub(tc.expectedTotalShares, tc.totalShares),
					mintedShares,
				)
				// Check that TotalShares is as expected.
				totalShares := tApp.App.VaultKeeper.GetTotalShares(ctx)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(tc.expectedTotalShares),
					totalShares,
				)
				// Check that OwnerShares is as expected.
				ownerShares, exists := tApp.App.VaultKeeper.GetOwnerShares(ctx, tc.owner)
				require.True(t, exists)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(tc.expectedOwnerShares),
					ownerShares,
				)
			}
		})
	}
}
