package keeper_test

import (
	"bytes"
	"math"
	"math/big"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	assetstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func TestMsgAllocateToVault(t *testing.T) {
	tests := map[string]struct {
		// Operator.
		operator string
		// Number of quote quantums main vault has.
		mainVaultQuantums uint64
		// Number of quote quantums sub vault has.
		subVaultQuantums uint64
		// Existing vault params, if any.
		vaultParams *vaulttypes.VaultParams
		// Msg.
		msg *vaulttypes.MsgAllocateToVault
		// Signer of above msg.
		signer string
		// A string that CheckTx response should contain, if any.
		checkTxResponseContains string
		// Whether CheckTx fails.
		checkTxFails bool
		// Whether DeliverTx fails.
		deliverTxFails bool
	}{
		"Success - Allocate 50 to Vault Clob 0, Existing vault params": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  0,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(50),
			},
			signer: constants.AliceAccAddress.String(),
		},
		"Success - Allocate 77 to Vault Clob 1, Non-existent Vault Params": {
			operator:          constants.BobAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.BobAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(77),
			},
			signer: constants.BobAccAddress.String(),
		},
		"Failure - Operator Authority, allocating more than max uint64 quantums": {
			operator:          constants.CarlAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority: constants.CarlAccAddress.String(),
				VaultId:   constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewIntFromBigInt(
					new(big.Int).Add(
						new(big.Int).SetUint64(math.MaxUint64),
						new(big.Int).SetUint64(1),
					),
				),
			},
			checkTxResponseContains: "QuoteQuantums must be positive and less than 2^64",
			checkTxFails:            true,
			signer:                  constants.CarlAccAddress.String(),
		},
		"Failure - Operator Authority, allocating zero quantums": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(0),
			},
			checkTxResponseContains: "QuoteQuantums must be positive",
			checkTxFails:            true,
			signer:                  constants.AliceAccAddress.String(),
		},
		"Failure - Operator Authority, allocating negative quantums": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(-1),
			},
			checkTxResponseContains: "QuoteQuantums must be positive",
			checkTxFails:            true,
			signer:                  constants.AliceAccAddress.String(),
		},
		"Failure - Operator Authority, Insufficient quantums to allocate to Vault Clob 0, Existing vault params": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			vaultParams: &vaulttypes.VaultParams{
				Status: vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
			},
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob0,
				QuoteQuantums: dtypes.NewInt(101),
			},
			signer:         constants.AliceAccAddress.String(),
			checkTxFails:   false,
			deliverTxFails: true,
		},
		"Failure - Operator Authority, No corresponding clob pair": {
			operator:          constants.AliceAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  0,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority: constants.AliceAccAddress.String(),
				VaultId: vaulttypes.VaultId{
					Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
					Number: 727,
				},
				QuoteQuantums: dtypes.NewInt(1),
			},
			signer:         constants.AliceAccAddress.String(),
			checkTxFails:   false,
			deliverTxFails: true,
		},
		"Failure - Invalid Authority, Non-existent Vault Params": {
			operator:          constants.BobAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     constants.AliceAccAddress.String(),
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(77),
			},
			signer:         constants.AliceAccAddress.String(),
			checkTxFails:   false,
			deliverTxFails: true,
		},
		"Failure - Empty Authority, Existing vault params": {
			operator:          constants.BobAccAddress.String(),
			mainVaultQuantums: 100,
			subVaultQuantums:  15,
			msg: &vaulttypes.MsgAllocateToVault{
				Authority:     "",
				VaultId:       constants.Vault_Clob1,
				QuoteQuantums: dtypes.NewInt(77),
			},
			signer:                  constants.BobAccAddress.String(),
			checkTxResponseContains: vaulttypes.ErrInvalidAuthority.Error(),
			checkTxFails:            true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Set megavault operator.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.OperatorParams = vaulttypes.OperatorParams{
							Operator: tc.operator,
						}
						if tc.vaultParams != nil {
							genesisState.Vaults = []vaulttypes.Vault{
								{
									VaultId:     tc.msg.VaultId,
									VaultParams: *tc.vaultParams,
								},
							}
						}
					},
				)
				// Set balances of main vault and sub vault.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assetstypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromUint64(tc.mainVaultQuantums),
									},
								},
							},
							{
								Id: tc.msg.VaultId.ToSubaccountId(),
								AssetPositions: []*satypes.AssetPosition{
									{
										AssetId:  assetstypes.AssetUsdc.Id,
										Quantums: dtypes.NewIntFromUint64(tc.subVaultQuantums),
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

			// Invoke CheckTx.
			CheckTx_MsgAllocateToVault := testapp.MustMakeCheckTx(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: tc.signer,
					Gas:                  constants.TestGasLimit,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				tc.msg,
			)
			checkTxResp := tApp.CheckTx(CheckTx_MsgAllocateToVault)

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
							if bytes.Equal(tx, CheckTx_MsgAllocateToVault.Tx) {
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

			// Check expectations.
			mainVault := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, vaulttypes.MegavaultMainSubaccount)
			subVault := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *tc.msg.VaultId.ToSubaccountId())
			require.Len(t, subVault.AssetPositions, 1)
			if tc.deliverTxFails {
				// Verify that main vault and sub vault balances are unchanged.
				require.Len(t, mainVault.AssetPositions, 1)
				require.Equal(
					t,
					tc.mainVaultQuantums,
					mainVault.AssetPositions[0].Quantums.BigInt().Uint64(),
				)
				require.Equal(
					t,
					tc.subVaultQuantums,
					subVault.AssetPositions[0].Quantums.BigInt().Uint64(),
				)

				// Verify that vault params is unchanged.
				vaultParams, exists := k.GetVaultParams(ctx, tc.msg.VaultId)
				if tc.vaultParams != nil {
					require.True(t, exists)
					require.Equal(t, *tc.vaultParams, vaultParams)
				} else {
					require.False(t, exists)
				}
			} else {
				// Verify that main vault and sub vault balances are updated.
				expectedMainVaultQuantums := tc.mainVaultQuantums - tc.msg.QuoteQuantums.BigInt().Uint64()
				if expectedMainVaultQuantums == 0 {
					require.Len(t, mainVault.AssetPositions, 0)
				} else {
					require.Len(t, mainVault.AssetPositions, 1)
					require.Equal(
						t,
						expectedMainVaultQuantums,
						mainVault.AssetPositions[0].Quantums.BigInt().Uint64(),
					)
				}
				require.Equal(
					t,
					tc.subVaultQuantums+tc.msg.QuoteQuantums.BigInt().Uint64(),
					subVault.AssetPositions[0].Quantums.BigInt().Uint64(),
				)

				// Verify that vault params is initialized if didn't exist before.
				vaultParams, exists := k.GetVaultParams(ctx, tc.msg.VaultId)
				require.True(t, exists)
				if tc.vaultParams != nil {
					require.Equal(t, *tc.vaultParams, vaultParams)
				} else {
					require.Equal(
						t,
						vaulttypes.VaultParams{
							Status: vaulttypes.VaultStatus_VAULT_STATUS_STAND_BY,
						},
						vaultParams,
					)
				}
			}
		})
	}
}
