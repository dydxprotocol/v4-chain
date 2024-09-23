package keeper_test

import (
	"bytes"
	"math"
	"math/big"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
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

type VaultSetup struct {
	id                   vaulttypes.VaultId
	params               vaulttypes.VaultParams
	assetQuoteQuantums   *big.Int
	positionBaseQuantums *big.Int
	clobPair             clobtypes.ClobPair
	perpetual            perptypes.Perpetual
	marketParam          pricestypes.MarketParam
	marketPrice          pricestypes.MarketPrice
	postWithdrawalEquity *big.Int
}

func TestMsgWithdrawFromMegavault(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Quote quantums that main vault has.
		mainVaultBalance *big.Int
		// Total shares before withdrawal.
		totalShares uint64
		// Owner address.
		owner string
		// Owner total shares.
		ownerTotalShares uint64
		// Owner locked shares.
		ownerLockedShares uint64
		// Vaults.
		vaults []VaultSetup
		// Shares to withdraw.
		sharesToWithdraw int64
		// Minimum quote quantums to redeem.
		minQuoteQuantums int64

		/* --- Expectations --- */
		// A string that CheckTx response contains, if any.
		checkTxResponseContains string
		// Whether CheckTx should fail.
		checkTxFails bool
		// Whether DeliverTx should fail.
		deliverTxFails bool
		// Quote quantums that should be redeemed.
		redeemedQuoteQuantums uint64
		// Total shares after withdrawal.
		expectedTotalShares uint64
		// Owner shares after withdrawal.
		expectedOwnerShares uint64
	}{
		"Success: Withdraw some unlocked shares (5% of total), No sub-vaults, Redeemed quantums = Min quantums": {
			mainVaultBalance:      big.NewInt(100),
			totalShares:           200,
			owner:                 constants.AliceAccAddress.String(),
			ownerTotalShares:      50,
			ownerLockedShares:     25,
			sharesToWithdraw:      10,
			minQuoteQuantums:      5,
			deliverTxFails:        false,
			redeemedQuoteQuantums: 5,   // 5% of 100
			expectedTotalShares:   190, // 200 - 10
			expectedOwnerShares:   40,  // 50 - 10
		},
		"Success: Withdraw all unlocked shares (8% of total), No sub-vaults, Redeemed quantums > Min quantums": {
			mainVaultBalance:      big.NewInt(1_234),
			totalShares:           500,
			owner:                 constants.BobAccAddress.String(),
			ownerTotalShares:      47,
			ownerLockedShares:     7,
			sharesToWithdraw:      40,
			minQuoteQuantums:      95,
			deliverTxFails:        false,
			redeemedQuoteQuantums: 98,  // 1234 * 0.08 = 98.72 ~= 98 (rounded down)
			expectedTotalShares:   460, // 500 - 40
			expectedOwnerShares:   7,   // 47 - 40
		},
		"Success: Withdraw all shares (100% of total), No sub-vaults, Redeemed quantums = Min quantums": {
			mainVaultBalance:      big.NewInt(654_321),
			totalShares:           787_565,
			owner:                 constants.CarlAccAddress.String(),
			ownerTotalShares:      787_565,
			ownerLockedShares:     0,
			sharesToWithdraw:      787_565,
			minQuoteQuantums:      654_321,
			deliverTxFails:        false,
			redeemedQuoteQuantums: 654_321, // all main vault balance
			expectedTotalShares:   0,
			expectedOwnerShares:   0,
		},
		"Failure: Withdraw some unlocked shares (1% of total), No sub-vaults, Redeemed quantums rounds down to 0": {
			mainVaultBalance:      big.NewInt(99),
			totalShares:           200,
			owner:                 constants.AliceAccAddress.String(),
			ownerTotalShares:      10,
			ownerLockedShares:     5,
			sharesToWithdraw:      2,
			minQuoteQuantums:      0,
			deliverTxFails:        true,
			redeemedQuoteQuantums: 0,   // 99 * 2 / 200 = 0.99 ~= 0 (rounded down)
			expectedTotalShares:   200, // unchanged
			expectedOwnerShares:   10,  // unchanged
		},
		"Failure: Withdraw more than locked shares": {
			mainVaultBalance:    big.NewInt(100),
			totalShares:         500,
			owner:               constants.AliceAccAddress.String(),
			ownerTotalShares:    100,
			ownerLockedShares:   20,
			sharesToWithdraw:    81,
			minQuoteQuantums:    1,
			deliverTxFails:      true,
			expectedTotalShares: 500, // unchanged
			expectedOwnerShares: 100, // unchanged
		},
		"Failure: Withdraw zero shares": {
			mainVaultBalance:        big.NewInt(100),
			totalShares:             500,
			owner:                   constants.AliceAccAddress.String(),
			ownerTotalShares:        100,
			ownerLockedShares:       20,
			sharesToWithdraw:        0,
			minQuoteQuantums:        1,
			checkTxResponseContains: vaulttypes.ErrNonPositiveShares.Error(),
			checkTxFails:            true,
			expectedTotalShares:     500, // unchanged
			expectedOwnerShares:     100, // unchanged
		},
		"Failure: Withdraw negative shares": {
			mainVaultBalance:        big.NewInt(100),
			totalShares:             500,
			owner:                   constants.AliceAccAddress.String(),
			ownerTotalShares:        100,
			ownerLockedShares:       20,
			sharesToWithdraw:        -1,
			minQuoteQuantums:        1,
			checkTxResponseContains: vaulttypes.ErrNonPositiveShares.Error(),
			checkTxFails:            true,
			expectedTotalShares:     500, // unchanged
			expectedOwnerShares:     100, // unchanged
		},
		"Failure: Withdraw some unlocked shares (8% of total), one sub-vault, Redeemed quantums < Min quantums": {
			mainVaultBalance:  big.NewInt(1_234),
			totalShares:       500,
			owner:             constants.BobAccAddress.String(),
			ownerTotalShares:  47,
			ownerLockedShares: 7,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob0,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
					},
					assetQuoteQuantums:   big.NewInt(400),
					positionBaseQuantums: big.NewInt(0),
					clobPair:             constants.ClobPair_Btc,
					perpetual:            constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[0],
					marketPrice:          constants.TestMarketPrices[0],
					postWithdrawalEquity: big.NewInt(400), // unchanged
				},
			},
			sharesToWithdraw:      40,
			minQuoteQuantums:      131, // greater than redeemed quote quantums
			deliverTxFails:        true,
			redeemedQuoteQuantums: 130, // 1234 * 0.08 + 400 * 0.08 = 130.72 ~= 130 (rounded down)
			expectedTotalShares:   500, // unchanged
			expectedOwnerShares:   47,  // unchanged
		},
		"Success: Withdraw some unlocked shares (8% of total), one deactivated sub-vault is excluded": {
			mainVaultBalance:  big.NewInt(1_234),
			totalShares:       500,
			owner:             constants.DaveAccAddress.String(),
			ownerTotalShares:  47,
			ownerLockedShares: 7,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob0,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_DEACTIVATED,
					},
					assetQuoteQuantums:   big.NewInt(-400),
					positionBaseQuantums: big.NewInt(0),
					clobPair:             constants.ClobPair_Btc,
					perpetual:            constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[0],
					marketPrice:          constants.TestMarketPrices[0],
					postWithdrawalEquity: big.NewInt(-400), // unchanged
				},
			},
			sharesToWithdraw:      40,
			minQuoteQuantums:      50,
			deliverTxFails:        false,
			redeemedQuoteQuantums: 98,  // 1234 * 0.08 = 98.72 ~= 98 (rounded down)
			expectedTotalShares:   460, // 500 - 40
			expectedOwnerShares:   7,   // 47 - 40
		},
		"Success: Withdraw some unlocked shares (0.4444% of total), 888_888 quantums in main vault, " +
			"one quoting sub-vault with negative equity": {
			mainVaultBalance:  big.NewInt(888_888),
			totalShares:       1_000_000,
			owner:             constants.AliceAccAddress.String(),
			ownerTotalShares:  9999,
			ownerLockedShares: 134,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob0,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
					assetQuoteQuantums:   big.NewInt(-345),
					positionBaseQuantums: big.NewInt(0),
					clobPair:             constants.ClobPair_Btc,
					perpetual:            constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[0],
					marketPrice:          constants.TestMarketPrices[0],
					postWithdrawalEquity: big.NewInt(-345),
				},
			},
			sharesToWithdraw:      4444,
			minQuoteQuantums:      123,
			deliverTxFails:        false,
			redeemedQuoteQuantums: 3_950,   // 888_888 * 4444 / 1_000_000 ~= 3950 (sub-vault is skipped)
			expectedTotalShares:   995_556, // 1_000_000 - 4444
			expectedOwnerShares:   5_555,   // 9999 - 4444
		},
		"Success: Withdraw some unlocked shares (~0.67% of total), 0 quantums in main vault, " +
			"one quoting sub-vault with 0 leverage": {
			mainVaultBalance:  big.NewInt(0),
			totalShares:       987_654,
			owner:             constants.AliceAccAddress.String(),
			ownerTotalShares:  9999,
			ownerLockedShares: 134,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob0,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
					},
					assetQuoteQuantums:   big.NewInt(345),
					positionBaseQuantums: big.NewInt(0),
					clobPair:             constants.ClobPair_Btc,
					perpetual:            constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[0],
					marketPrice:          constants.TestMarketPrices[0],
					postWithdrawalEquity: big.NewInt(343), // 345 - 2
				},
			},
			sharesToWithdraw:      6666,
			minQuoteQuantums:      2,
			deliverTxFails:        false,
			redeemedQuoteQuantums: 2,       // 345 * 6666 / 987654 ~= 2.32 ~= 2 (rounded down)
			expectedTotalShares:   980_988, // 987654 - 6666
			expectedOwnerShares:   3333,    // 9999 - 6666
		},
		"Success: Withdraw some unlocked shares (10% of total), 500 quantums in main vault, " +
			"one stand-by sub-vault with 0 leverage, one close-only sub-vault with 1.5 leverage": {
			mainVaultBalance:  big.NewInt(500),
			totalShares:       1_000,
			owner:             constants.AliceAccAddress.String(),
			ownerTotalShares:  120,
			ownerLockedShares: 15,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob0,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
					},
					assetQuoteQuantums:   big.NewInt(345),
					positionBaseQuantums: big.NewInt(0),
					clobPair:             constants.ClobPair_Btc,
					perpetual:            constants.BtcUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[0],
					marketPrice:          constants.TestMarketPrices[0],
					postWithdrawalEquity: big.NewInt(311), // 345 - 34
				},
				{
					id: constants.Vault_Clob1,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
						QuotingParams: &vaulttypes.QuotingParams{
							Layers:                           3,
							SpreadMinPpm:                     3_000,
							SpreadBufferPpm:                  1_500,
							SkewFactorPpm:                    3_000_000,
							OrderSizePctPpm:                  100_000,
							OrderExpirationSeconds:           60,
							ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000),
						},
					},
					// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
					// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
					assetQuoteQuantums:   big.NewInt(-1_000_000_000),
					positionBaseQuantums: big.NewInt(1_000_000_000),
					clobPair:             constants.ClobPair_Eth,
					perpetual:            constants.EthUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[1],
					marketPrice:          constants.TestMarketPrices[1],
					postWithdrawalEquity: big.NewInt(1_829_775_000), // 2_000_000_000 - 170_225_000
				},
			},
			sharesToWithdraw: 100,
			minQuoteQuantums: 170_000_000,
			deliverTxFails:   false,
			// Main vault withdrawal + sub-vault 0 withdrawal + sub-vault 1 withdrawal
			// = 500 * 100 / 1_000 + 345 * 100 / 1_000 + 2_000_000_000 * 100 / 1_000 * (1 - 1191/8000)
			// ~= 50 + 34 + 170225000
			// = 170225084
			redeemedQuoteQuantums: 170_225_084,
			expectedTotalShares:   900, // 1_000 - 100
			expectedOwnerShares:   20,  // 120 - 100
		},
		"Success: Withdraw all shares (100% of total), 500 quantums in main vault, " +
			"one close-only sub-vault with 1.5 leverage": {
			mainVaultBalance:  big.NewInt(500),
			totalShares:       1_000,
			owner:             constants.AliceAccAddress.String(),
			ownerTotalShares:  1_000,
			ownerLockedShares: 0,
			vaults: []VaultSetup{
				{
					id: constants.Vault_Clob1,
					params: vaulttypes.VaultParams{
						Status: vaulttypes.VaultStatus_VAULT_STATUS_CLOSE_ONLY,
					},
					// open_notional = 1_000_000_000 * 10^-9 * 3_000 * 10^6 = 3_000_000_000
					// leverage = 3_000_000_000 / (-1_000_000_000 + 3_000_000_000) = 1.5
					assetQuoteQuantums:   big.NewInt(-1_000_000_000),
					positionBaseQuantums: big.NewInt(1_000_000_000),
					clobPair:             constants.ClobPair_Eth,
					perpetual:            constants.EthUsd_20PercentInitial_10PercentMaintenance,
					marketParam:          constants.TestMarketParams[1],
					marketPrice:          constants.TestMarketPrices[1],
					postWithdrawalEquity: big.NewInt(600_000_000), // 2_000_000_000 - 1_400_000_000
				},
			},
			sharesToWithdraw: 1_000,
			minQuoteQuantums: 1_400_000_500,
			deliverTxFails:   false,
			// Main vault withdrawal + sub-vault 0 withdrawal
			// = 500 * 1_000 / 1_000 + 2_000_000_000 * 1_000 / 1_000 * (1 - leverage * imf)
			// = 500 + 2_000_000_000 * (1 - 1.5 * 0.2)
			redeemedQuoteQuantums: 1_400_000_500,
			expectedTotalShares:   0,
			expectedOwnerShares:   0,
		},
		"Failure: Withdraw more than max uint64": {
			mainVaultBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				new(big.Int).SetUint64(1),
			),
			totalShares:       155,
			owner:             constants.CarlAccAddress.String(),
			ownerTotalShares:  155,
			ownerLockedShares: 0,
			sharesToWithdraw:  155,
			minQuoteQuantums:  1,
			// fails as owner redeems more than max uint64 quote quantums.
			deliverTxFails:      true,
			expectedTotalShares: 155, // unchanged
			expectedOwnerShares: 155, // unchanged
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
						genesisState.OwnerShares = []vaulttypes.OwnerShare{
							{
								Owner: tc.owner,
								Shares: vaulttypes.NumShares{
									NumShares: dtypes.NewIntFromUint64(tc.ownerTotalShares),
								},
							},
						}
						genesisState.AllOwnerShareUnlocks = []vaulttypes.OwnerShareUnlocks{
							{
								OwnerAddress: tc.owner,
								ShareUnlocks: []vaulttypes.ShareUnlock{
									{
										Shares: vaulttypes.NumShares{
											NumShares: dtypes.NewIntFromUint64(tc.ownerLockedShares),
										},
										UnlockBlockHeight: 7, // dummy height
									},
								},
							},
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

			// Construct message.
			msgWithdrawFromMegavault := vaulttypes.MsgWithdrawFromMegavault{
				SubaccountId: satypes.SubaccountId{
					Owner:  tc.owner,
					Number: 0,
				},
				Shares: vaulttypes.NumShares{
					NumShares: dtypes.NewInt(tc.sharesToWithdraw),
				},
				MinQuoteQuantums: dtypes.NewInt(tc.minQuoteQuantums),
			}

			preMegavaultEquity, err := tApp.App.VaultKeeper.GetMegavaultEquity(ctx)
			require.NoError(t, err)
			preOwnerEquity, err := tApp.App.VaultKeeper.GetSubaccountEquity(ctx, msgWithdrawFromMegavault.SubaccountId)
			require.NoError(t, err)

			// Invoke CheckTx.
			CheckTx_MsgWithdrawFromMegavault := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: tc.owner,
					Gas:                  constants.TestGasLimit,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				&msgWithdrawFromMegavault,
			)
			checkTxResp := tApp.CheckTx(CheckTx_MsgWithdrawFromMegavault)
			// Check that CheckTx response log contains expected string, if any.
			if tc.checkTxResponseContains != "" {
				require.Contains(t, checkTxResp.Log, tc.checkTxResponseContains)
			}
			// Check that CheckTx succeeds or errors out as expected.
			if tc.checkTxFails {
				require.Conditionf(t, checkTxResp.IsErr, "Expected CheckTx to error. Response: %+v", checkTxResp)
				return
			}
			require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

			// Advance to next block (and check that DeliverTx is as expected).
			nextBlock := uint32(ctx.BlockHeight()) + 1
			if tc.deliverTxFails {
				// Check that DeliverTx fails on `msgDepositToMegavault`.
				ctx = tApp.AdvanceToBlock(nextBlock, testapp.AdvanceToBlockOptions{
					ValidateFinalizeBlock: func(
						context sdktypes.Context,
						request abcitypes.RequestFinalizeBlock,
						response abcitypes.ResponseFinalizeBlock,
					) (haltChain bool) {
						for i, tx := range request.Txs {
							if bytes.Equal(tx, CheckTx_MsgWithdrawFromMegavault.Tx) {
								require.True(t, response.TxResults[i].IsErr())
							} else {
								require.True(t, response.TxResults[i].IsOK())
							}
						}
						return false
					},
				})
			} else {
				ctx = tApp.AdvanceToBlock(nextBlock, testapp.AdvanceToBlockOptions{})
			}

			// Check total shares.
			totalShares := tApp.App.VaultKeeper.GetTotalShares(ctx)
			require.Equal(
				t,
				new(big.Int).SetUint64(tc.expectedTotalShares),
				totalShares.NumShares.BigInt(),
			)
			// Check owner shares.
			ownerShares, exists := tApp.App.VaultKeeper.GetOwnerShares(
				ctx,
				tc.owner,
			)
			if tc.expectedOwnerShares == 0 {
				require.False(t, exists)
			} else {
				require.True(t, exists)
				require.Equal(
					t,
					new(big.Int).SetUint64(tc.expectedOwnerShares),
					ownerShares.NumShares.BigInt(),
				)
			}
			// Check equity of owner, megavault, and each sub-vault.
			postOwnerEquity, err := tApp.App.VaultKeeper.GetSubaccountEquity(ctx, msgWithdrawFromMegavault.SubaccountId)
			require.NoError(t, err)
			postMegavaultEquity, err := tApp.App.VaultKeeper.GetMegavaultEquity(ctx)
			require.NoError(t, err)
			if tc.deliverTxFails {
				require.Equal(t, preOwnerEquity, postOwnerEquity)
				require.Equal(t, preMegavaultEquity, postMegavaultEquity)
			} else {
				require.Equal(
					t,
					preOwnerEquity.Uint64()+tc.redeemedQuoteQuantums,
					postOwnerEquity.Uint64(),
				)
				require.Equal(
					t,
					preMegavaultEquity.Uint64()-tc.redeemedQuoteQuantums,
					postMegavaultEquity.Uint64(),
				)
			}
			for _, vault := range tc.vaults {
				subVaultEquity, err := tApp.App.VaultKeeper.GetSubaccountEquity(ctx, *vault.id.ToSubaccountId())
				require.NoError(t, err)
				require.Equal(
					t,
					vault.postWithdrawalEquity,
					subVaultEquity,
				)
			}
		})
	}
}
