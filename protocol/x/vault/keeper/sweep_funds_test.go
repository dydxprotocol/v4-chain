package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/gogoproto/proto"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

func TestSweepMainVaultBankBalance(t *testing.T) {
	tests := map[string]struct {
		// Bank balance of main vault
		bankBalance int64
		// Subaccount balance of main vault
		subaccountBalance *big.Int
		// Expected bank balance of main vault
		expectedBankBalance int64
		// Expected subaccount balance of main vault
		expectedSubaccountBalance *big.Int
	}{
		"Zero bank balance, zero subaccount balance": {
			bankBalance:               0,
			subaccountBalance:         big.NewInt(0),
			expectedBankBalance:       0,
			expectedSubaccountBalance: big.NewInt(0),
		},
		"100_000_000 quantums bank balance, zero subaccount balance": {
			bankBalance:               100_000_000,
			subaccountBalance:         big.NewInt(0),
			expectedBankBalance:       0,
			expectedSubaccountBalance: big.NewInt(100_000_000),
		},
		"100_000_000 quantums bank balance, 50_000_000 subaccount balance": {
			bankBalance:               100_000_000,
			subaccountBalance:         big.NewInt(50_000_000),
			expectedBankBalance:       0,
			expectedSubaccountBalance: big.NewInt(150_000_000),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Initialize vaults with their equities.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										assettypes.AssetUsdc.Id,
										tc.subaccountBalance,
									),
								},
							},
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
							Address: vaulttypes.MegavaultMainAddress.String(),
							Coins: sdktypes.Coins{
								sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(tc.bankBalance)),
							},
						})
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			k.SweepMainVaultBankBalance(ctx)

			mainVaultSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, vaulttypes.MegavaultMainSubaccount)
			require.Equal(t, tc.expectedSubaccountBalance, mainVaultSubaccount.AssetPositions[0].Quantums.BigInt())
			mainVaultBankBalance := tApp.App.BankKeeper.GetBalance(
				ctx,
				vaulttypes.MegavaultMainAddress,
				constants.Usdc.Denom,
			).Amount
			require.Equal(t, sdkmath.NewIntFromBigInt(big.NewInt(tc.expectedBankBalance)), mainVaultBankBalance)
		})
	}
}

func TestSweepMainVaultBankBalance_EndBlock(t *testing.T) {
	tests := map[string]struct {
		// Bank balance of main vault
		bankBalance int64
		// Subaccount balance of main vault
		subaccountBalance uint64
		// Expected bank balance of main vault
		expectedBankBalance int64
		// Expected subaccount balance of main vault
		expectedSubaccountBalance *big.Int
	}{
		"Zero bank balance, zero subaccount balance": {
			bankBalance:               0,
			subaccountBalance:         0,
			expectedBankBalance:       0,
			expectedSubaccountBalance: big.NewInt(0),
		},
		"100_000_000 quantums bank balance, zero subaccount balance": {
			bankBalance:               100_000_000,
			subaccountBalance:         0,
			expectedBankBalance:       0,
			expectedSubaccountBalance: big.NewInt(100_000_000),
		},
		"100_000_000 quantums bank balance, 50_000_000 subaccount balance": {
			bankBalance:               100_000_000,
			subaccountBalance:         50_000_000,
			expectedBankBalance:       0,
			expectedSubaccountBalance: big.NewInt(150_000_000),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(
										assettypes.AssetUsdc.Id,
										big.NewInt(0),
									),
								},
							},
						}
					},
				)
				return genesis
			}).Build()

			// Fund the subaccount and bank balance of megavault
			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})
			if tc.subaccountBalance > 0 {
				var msg proto.Message
				depositMsg := sendingtypes.MsgDepositToSubaccount{
					Sender:    constants.AliceAccAddress.String(),
					Recipient: vaulttypes.MegavaultMainSubaccount,
					AssetId:   constants.Usdc.Id,
					Quantums:  tc.subaccountBalance,
				}
				msg = &depositMsg
				for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: constants.AliceAccAddress.String(),
						Gas:                  1000000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					msg,
				) {
					resp := tApp.CheckTx(checkTx)
					require.Condition(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			if tc.bankBalance > 0 {
				var msg proto.Message
				bankSendMsg := banktypes.MsgSend{
					FromAddress: constants.AliceAccAddress.String(),
					ToAddress:   vaulttypes.MegavaultMainAddress.String(),
					Amount: sdktypes.Coins{
						sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(tc.bankBalance)),
					},
				}
				msg = &bankSendMsg
				for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
					ctx,
					tApp.App,
					testapp.MustMakeCheckTxOptions{
						AccAddressForSigning: constants.AliceAccAddress.String(),
						Gas:                  1000000,
						FeeAmt:               constants.TestFeeCoins_5Cents,
					},
					msg,
				) {
					resp := tApp.CheckTx(checkTx)
					require.Condition(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
				}
			}
			// Advance block to execute EndBlocker
			ctx = tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{})

			mainVaultSubaccount := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, vaulttypes.MegavaultMainSubaccount)
			require.Equal(t, tc.expectedSubaccountBalance, mainVaultSubaccount.AssetPositions[0].Quantums.BigInt())
			mainVaultBankBalance := tApp.App.BankKeeper.GetBalance(
				ctx,
				vaulttypes.MegavaultMainAddress,
				constants.Usdc.Denom,
			).Amount
			require.Equal(t, sdkmath.NewIntFromBigInt(big.NewInt(tc.expectedBankBalance)), mainVaultBankBalance)
		})
	}
}
