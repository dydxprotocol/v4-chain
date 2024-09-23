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
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

// DepositInstance represents an instance of a deposit to test.
type DepositInstance struct {
	// Depositor.
	depositor satypes.SubaccountId
	// Amount to deposit (in quote quantums).
	depositAmount *big.Int
	// Signer of the message.
	msgSigner string

	// A string that CheckTx response should contain, if any.
	checkTxResponseContains string
	// Whether CheckTx fails.
	checkTxFails bool
	// Whether DeliverTx fails.
	deliverTxFails bool
	// Expected owner shares for depositor above.
	expectedOwnerShares *big.Int
}

// DepositorSetup represents the setup of a depositor.
type DepositorSetup struct {
	// Depositor.
	depositor satypes.SubaccountId
	// Initial balance of the depositor (in quote quantums).
	depositorBalance *big.Int
}

func TestMsgDepositToMegavault(t *testing.T) {
	tests := map[string]struct {
		/* --- Setup --- */
		// Depositor setups.
		depositorSetups []DepositorSetup
		// Instances of deposits.
		depositInstances []DepositInstance

		/* --- Expectations --- */
		// Total shares after each of the above deposit instances.
		totalSharesHistory []*big.Int
		// Megavault equity after each of the above deposit instances.
		equityHistory []*big.Int
	}{
		"Two successful deposits, Same depositor": {
			depositorSetups: []DepositorSetup{
				{
					depositor:        constants.Alice_Num0,
					depositorBalance: big.NewInt(1000),
				},
			},
			depositInstances: []DepositInstance{
				{
					depositor:           constants.Alice_Num0,
					depositAmount:       big.NewInt(123),
					msgSigner:           constants.Alice_Num0.Owner,
					expectedOwnerShares: big.NewInt(123),
				},
				{
					depositor:           constants.Alice_Num0,
					depositAmount:       big.NewInt(321),
					msgSigner:           constants.Alice_Num0.Owner,
					expectedOwnerShares: big.NewInt(444),
				},
			},
			totalSharesHistory: []*big.Int{
				big.NewInt(123),
				big.NewInt(444),
			},
			equityHistory: []*big.Int{
				big.NewInt(123),
				big.NewInt(444),
			},
		},
		"Two successful deposits, Different depositors": {
			depositorSetups: []DepositorSetup{
				{
					depositor:        constants.Alice_Num0,
					depositorBalance: big.NewInt(1_000),
				},
				{
					depositor:        constants.Bob_Num1,
					depositorBalance: big.NewInt(500),
				},
			},
			depositInstances: []DepositInstance{
				{
					depositor:           constants.Alice_Num0,
					depositAmount:       big.NewInt(1_000),
					msgSigner:           constants.Alice_Num0.Owner,
					expectedOwnerShares: big.NewInt(1_000),
				},
				{
					depositor:           constants.Bob_Num1,
					depositAmount:       big.NewInt(500),
					msgSigner:           constants.Bob_Num1.Owner,
					expectedOwnerShares: big.NewInt(500),
				},
			},
			totalSharesHistory: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(1_500),
			},
			equityHistory: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(1_500),
			},
		},
		"One successful deposit, One failed deposit due to insufficient balance": {
			depositorSetups: []DepositorSetup{
				{
					depositor:        constants.Alice_Num0,
					depositorBalance: big.NewInt(1_000),
				},
				{
					depositor:        constants.Bob_Num1,
					depositorBalance: big.NewInt(500),
				},
			},
			depositInstances: []DepositInstance{
				{
					depositor:           constants.Alice_Num0,
					depositAmount:       big.NewInt(1_000),
					msgSigner:           constants.Alice_Num0.Owner,
					expectedOwnerShares: big.NewInt(1_000),
				},
				{
					depositor:           constants.Bob_Num1,
					depositAmount:       big.NewInt(501), // Greater than balance.
					msgSigner:           constants.Bob_Num1.Owner,
					deliverTxFails:      true,
					expectedOwnerShares: nil,
				},
			},
			totalSharesHistory: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(1_000),
			},
			equityHistory: []*big.Int{
				big.NewInt(1_000),
				big.NewInt(1_000),
			},
		},
		"One failed deposit due to incorrect signer, One successful deposit": {
			depositorSetups: []DepositorSetup{
				{
					depositor:        constants.Alice_Num0,
					depositorBalance: big.NewInt(1_000),
				},
				{
					depositor:        constants.Bob_Num1,
					depositorBalance: big.NewInt(500),
				},
			},
			depositInstances: []DepositInstance{
				{
					depositor:               constants.Bob_Num1,
					depositAmount:           big.NewInt(500),
					msgSigner:               constants.Alice_Num0.Owner, // Incorrect signer.
					checkTxFails:            true,
					checkTxResponseContains: "does not match signer address",
					expectedOwnerShares:     nil,
				},
				{
					depositor:           constants.Alice_Num0,
					depositAmount:       big.NewInt(1_000),
					msgSigner:           constants.Alice_Num0.Owner,
					expectedOwnerShares: big.NewInt(1_000),
				},
			},
			totalSharesHistory: []*big.Int{
				big.NewInt(0),
				big.NewInt(1_000),
			},
			equityHistory: []*big.Int{
				big.NewInt(0),
				big.NewInt(1_000),
			},
		},
		"Three failed deposits due to invalid deposit amount": {
			depositorSetups: []DepositorSetup{
				{
					depositor:        constants.Alice_Num0,
					depositorBalance: big.NewInt(1_000),
				},
				{
					depositor:        constants.Bob_Num0,
					depositorBalance: big.NewInt(1_000),
				},
			},
			depositInstances: []DepositInstance{
				{
					depositor:               constants.Alice_Num0,
					depositAmount:           big.NewInt(0),
					msgSigner:               constants.Alice_Num0.Owner,
					checkTxFails:            true,
					checkTxResponseContains: vaulttypes.ErrInvalidQuoteQuantums.Error(),
					expectedOwnerShares:     nil,
				},
				{
					depositor:               constants.Bob_Num0,
					depositAmount:           big.NewInt(-1),
					msgSigner:               constants.Bob_Num0.Owner,
					checkTxFails:            true,
					checkTxResponseContains: vaulttypes.ErrInvalidQuoteQuantums.Error(),
					expectedOwnerShares:     nil,
				},
				{
					depositor: constants.Bob_Num0,
					depositAmount: new(big.Int).Add(
						new(big.Int).SetUint64(math.MaxUint64),
						big.NewInt(1),
					),
					msgSigner:               constants.Bob_Num0.Owner,
					checkTxFails:            true,
					checkTxResponseContains: vaulttypes.ErrInvalidQuoteQuantums.Error(),
					expectedOwnerShares:     nil,
				},
			},
			totalSharesHistory: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
				big.NewInt(0),
			},
			equityHistory: []*big.Int{
				big.NewInt(0),
				big.NewInt(0),
				big.NewInt(0),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize tApp and ctx.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Initialize balances of depositors.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						subaccounts := make([]satypes.Subaccount, len(tc.depositorSetups))
						for i, setup := range tc.depositorSetups {
							subaccounts[i] = satypes.Subaccount{
								Id: &(setup.depositor),
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										0,
										setup.depositorBalance,
									),
								},
							}
						}
						genesisState.Subaccounts = subaccounts
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Simulate each deposit instance.
			for i, depositInstance := range tc.depositInstances {
				// Construct message.
				msgDepositToMegavault := vaulttypes.MsgDepositToMegavault{
					SubaccountId:  &(depositInstance.depositor),
					QuoteQuantums: dtypes.NewIntFromBigInt(depositInstance.depositAmount),
				}

				// Invoke CheckTx.
				CheckTx_MsgDepositToMegavault := testapp.MustMakeCheckTx(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: depositInstance.msgSigner,
						Gas:                  constants.TestGasLimit,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					&msgDepositToMegavault,
				)
				checkTxResp := tApp.CheckTx(CheckTx_MsgDepositToMegavault)

				// Check that CheckTx response log contains expected string, if any.
				if depositInstance.checkTxResponseContains != "" {
					require.Contains(t, checkTxResp.Log, depositInstance.checkTxResponseContains)
				}
				// Check that CheckTx succeeds or errors out as expected.
				if depositInstance.checkTxFails {
					require.Conditionf(t, checkTxResp.IsErr, "Expected CheckTx to error. Response: %+v", checkTxResp)
					return
				}
				require.Conditionf(t, checkTxResp.IsOK, "Expected CheckTx to succeed. Response: %+v", checkTxResp)

				// Advance to next block (and check that DeliverTx is as expected).
				nextBlock := uint32(ctx.BlockHeight()) + 1
				if depositInstance.deliverTxFails {
					// Check that DeliverTx fails on `msgDepositToMegavault`.
					ctx = tApp.AdvanceToBlock(nextBlock, testapp.AdvanceToBlockOptions{
						ValidateFinalizeBlock: func(
							context sdktypes.Context,
							request abcitypes.RequestFinalizeBlock,
							response abcitypes.ResponseFinalizeBlock,
						) (haltChain bool) {
							for i, tx := range request.Txs {
								if bytes.Equal(tx, CheckTx_MsgDepositToMegavault.Tx) {
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

				// Check that total shares is as expected.
				totalShares := tApp.App.VaultKeeper.GetTotalShares(ctx)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(tc.totalSharesHistory[i]),
					totalShares,
				)
				// Check that owner shares of the depositor is as expected.
				ownerShares, _ := tApp.App.VaultKeeper.GetOwnerShares(
					ctx,
					depositInstance.depositor.Owner,
				)
				require.Equal(
					t,
					vaulttypes.BigIntToNumShares(depositInstance.expectedOwnerShares),
					ownerShares,
				)
				// Check that equity of megavault is as expected.
				vaultEquity, err := tApp.App.VaultKeeper.GetMegavaultEquity(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.equityHistory[i], vaultEquity)
			}
		})
	}
}
