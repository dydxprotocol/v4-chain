package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestQueryMegavaultWithdrawalInfo(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Quote quantums that main vault has.
		mainVaultBalance *big.Int
		// Total shares.
		totalShares uint64
		// Vaults.
		vaults []VaultSetup
		// Query request.
		req *vaulttypes.QueryMegavaultWithdrawalInfoRequest

		/* --- Expectations --- */
		res         *vaulttypes.QueryMegavaultWithdrawalInfoResponse
		expectedErr string
	}{
		"Success: Withdraw 10%, 100 quantums in main vault, one sub-vault with 0 leverage and 50 equity": {
			mainVaultBalance: big.NewInt(100),
			totalShares:      50,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob0,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
					assetQuoteQuantums:   big.NewInt(50),
					positionBaseQuantums: big.NewInt(0),
					clobPair:             constants.ClobPair_Btc,
					perpetual:            constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[0],
					marketPrice:          constants.TestMarketPrices[0],
				},
			},
			req: &vaulttypes.QueryMegavaultWithdrawalInfoRequest{
				SharesToWithdraw: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(5),
				},
			},
			res: &vaulttypes.QueryMegavaultWithdrawalInfoResponse{
				SharesToWithdraw: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(5),
				},
				ExpectedQuoteQuantums: dtypes.NewInt(15),
				MegavaultEquity:       dtypes.NewInt(150),
				TotalShares: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(50),
				},
			},
		},
		"Success: Withdraw ~0.65%, 5_471_283_193_197 quantums in main vault," +
			"one sub-vault with 1 leverage and 1_500_000_000 equity": {
			mainVaultBalance: big.NewInt(5_471_283_193_197),
			totalShares:      128_412_843_128,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob1,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
					// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
					// equity = -1_500_000_000 + 3_000_000_000 = 1_500_000_000
					assetQuoteQuantums:   big.NewInt(-1_500_000_000),
					positionBaseQuantums: big.NewInt(1_000_000_000),
					clobPair:             constants.ClobPair_Eth,
					perpetual:            constants.EthUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[1],
					marketPrice:          constants.TestMarketPrices[1],
				},
			},
			req: &vaulttypes.QueryMegavaultWithdrawalInfoRequest{
				SharesToWithdraw: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(831_571_304),
				},
			},
			res: &vaulttypes.QueryMegavaultWithdrawalInfoResponse{
				SharesToWithdraw: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(831_571_304),
				},
				// expected_quote_quantums
				// = shares_to_withdraw / total_shares * main_vault_balance +
				// shares_to_withdraw / total_shares * sub_vault_equity * (1 - slippage)
				// = 831_571_304 / 128_412_843_128 * 5_471_283_193_197 +
				// 831_571_304 / 128_412_843_128 * 1_500_000_000 * (1 - 0.4)
				// = 35_430_740_327 + 5_828_187 = 35_436_568_514
				ExpectedQuoteQuantums: dtypes.NewInt(35_436_568_514),
				// megavault_equity
				// = main_vault_balance + sub_vault_equity
				// = 5_471_283_193_197 + 1_500_000_000
				// = 5_472_783_193_197
				MegavaultEquity: dtypes.NewInt(5_472_783_193_197),
				TotalShares: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(128_412_843_128),
				},
			},
		},
		"Error: Withdraw more than total shares": {
			mainVaultBalance: big.NewInt(100),
			totalShares:      50,
			req: &vaulttypes.QueryMegavaultWithdrawalInfoRequest{
				SharesToWithdraw: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(51),
				},
			},
			expectedErr: vaulttypes.ErrInvalidSharesToWithdraw.Error(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						subaccounts := []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assetstypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromBigInt(tc.mainVaultBalance),
									},
								},
							},
						}
						for _, vault := range tc.vaults {
							subaccounts = append(subaccounts, satypes.Subaccount{
								Id: vault.id.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assetstypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromBigInt(vault.assetQuoteQuantums),
									},
								},
								PerpetualPositions: []*satypes.PerpetualPosition{
									{
										PerpetualId: vault.perpetual.Params.Id,
										Quantums:    dtypes.NewIntFromBigInt(vault.positionBaseQuantums),
									},
								},
							})
						}
						genesisState.Subaccounts = subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.TotalShares = vaulttypes.NumShares{
							NumShares: dtypes.NewIntFromUint64(tc.totalShares),
						}
						vaults := make([]vaulttypes.Vault, len(tc.vaults))
						for i, vault := range tc.vaults {
							vaults[i] = vaulttypes.Vault{
								VaultId:     vault.id,
								VaultParams: vault.params,
							}
						}
						genesisState.Vaults = vaults
					},
				)
				// Initialize prices.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *pricestypes.GenesisState) {
						marketParams := make([]pricestypes.MarketParam, len(tc.vaults))
						marketPrices := make([]pricestypes.MarketPrice, len(tc.vaults))
						for i, vault := range tc.vaults {
							vault.marketParam.Id = vault.id.Number
							marketParams[i] = vault.marketParam
							vault.marketPrice.Id = vault.id.Number
							marketPrices[i] = vault.marketPrice
						}
						genesisState.MarketParams = marketParams
						genesisState.MarketPrices = marketPrices
					},
				)
				// Initialize perpetuals.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.LiquidityTiers = constants.LiquidityTiers
						perpetuals := make([]perptypes.Perpetual, len(tc.vaults))
						for i, vault := range tc.vaults {
							vault.perpetual.Params.Id = vault.id.Number
							perpetuals[i] = vault.perpetual
						}
						genesisState.Perpetuals = perpetuals
					},
				)
				// Initialize clob pairs.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *clobtypes.GenesisState) {
						clobPairs := make([]clobtypes.ClobPair, len(tc.vaults))
						for i, vault := range tc.vaults {
							vault.clobPair.Id = vault.id.Number
							clobPairs[i] = vault.clobPair
						}
						genesisState.ClobPairs = clobPairs
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			response, err := tApp.App.VaultKeeper.MegavaultWithdrawalInfo(ctx, tc.req)
			if tc.expectedErr != "" {
				require.ErrorContains(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.res, response)
			}
		})
	}
}
