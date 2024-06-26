package sending_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/gogoproto/proto"

	sdkmath "cosmossdk.io/math"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	feetiertypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestTransfer_Isolated_Non_Isolated_Subaccounts(t *testing.T) {
	tests := map[string]struct {
		// State.
		subaccounts            []satypes.Subaccount
		collateralPoolBalances map[string]int64

		// Parameters.
		senderSubaccountId   satypes.SubaccountId
		receiverSubaccountId satypes.SubaccountId
		quantums             uint64

		// Configuration.
		liquidityTiers []perptypes.LiquidityTier
		perpetuals     []perptypes.Perpetual
		clobPairs      []clobtypes.ClobPair

		// Expectations.
		expectedSubaccounts            []satypes.Subaccount
		expectedCollateralPoolBalances map[string]int64
		expectedErr                    string
	}{
		`Can transfer from isolated subaccount to non-isolated subaccount, and coins are sent from
		isolated subaccount collateral pool to cross collateral pool`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // $10,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
			},
			senderSubaccountId:   constants.Alice_Num0,
			receiverSubaccountId: constants.Bob_Num0,
			quantums:             100_000_000, // $100
			liquidityTiers:       constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			expectedSubaccounts: []satypes.Subaccount{
				testutil.ChangeUsdcBalance(constants.Alice_Num0_1ISO_LONG_10_000USD, -100_000_000),
				testutil.ChangeUsdcBalance(constants.Bob_Num0_10_000USD, 100_000_000),
			},
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_100_000_000, // $10,100 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 9_900_000_000, // $9,900 USDC
			},
			expectedErr: "",
		},
		`Can transfer from non-isolated subaccount to isolated subaccount, and coins are sent from
		cross collateral pool to isolated subaccount collateral pool`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // $10,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
			},
			senderSubaccountId:   constants.Bob_Num0,
			receiverSubaccountId: constants.Alice_Num0,
			quantums:             100_000_000, // $100
			liquidityTiers:       constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			expectedSubaccounts: []satypes.Subaccount{
				testutil.ChangeUsdcBalance(constants.Alice_Num0_1ISO_LONG_10_000USD, 100_000_000),
				testutil.ChangeUsdcBalance(constants.Bob_Num0_10_000USD, -100_000_000),
			},
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 9_900_000_000, // $9,900 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_100_000_000, // $10,100 USDC
			},
			expectedErr: "",
		},
		`Can transfer from isolated subaccount to isolated subaccount in different isolated markets, and
		coins are sent from one isolated collateral pool to the other isolated collateral pool`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_1ISO2_LONG_10_000USD,
			},
			collateralPoolBalances: map[string]int64{
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.Iso2Usd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
			},
			senderSubaccountId:   constants.Alice_Num0,
			receiverSubaccountId: constants.Bob_Num0,
			quantums:             100_000_000, // $100
			liquidityTiers:       constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			expectedSubaccounts: []satypes.Subaccount{
				testutil.ChangeUsdcBalance(constants.Alice_Num0_1ISO_LONG_10_000USD, -100_000_000),
				testutil.ChangeUsdcBalance(constants.Bob_Num0_1ISO2_LONG_10_000USD, 100_000_000),
			},
			expectedCollateralPoolBalances: map[string]int64{
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 9_900_000_000, // $9,900 USDC
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.Iso2Usd_IsolatedMarket.Params.Id),
				).String(): 10_100_000_000, // $10,100 USDC
			},
			expectedErr: "",
		},
		`Can't transfer from isolated subaccount to non-isolated subaccount if collateral pool has 
		insufficient funds`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_10_000USD,
			},
			collateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // $10,000 USDC
				// Isolated perpetual collateral pool has no entry and is therefore empty ($0).
			},
			senderSubaccountId:   constants.Alice_Num0,
			receiverSubaccountId: constants.Bob_Num0,
			quantums:             100_000_000, // $100
			liquidityTiers:       constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD, // No changes as transfer should fail.
				constants.Bob_Num0_10_000USD,
			},
			expectedCollateralPoolBalances: map[string]int64{
				satypes.ModuleAddress.String(): 10_000_000_000, // No changes to collateral pools as transfer fails.
			},
			expectedErr: "insufficient funds",
		},
		`Can't transfer from isolated subaccount to isolated subaccount with different isolated 
		perpetuals if collateral pool has insufficient funds`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_1ISO2_LONG_10_000USD,
			},
			collateralPoolBalances: map[string]int64{
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.Iso2Usd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
				// Isolated perpetual collateral pool has no entry and is therefore empty ($0).
			},
			senderSubaccountId:   constants.Alice_Num0,
			receiverSubaccountId: constants.Bob_Num0,
			quantums:             100_000_000, // $100
			liquidityTiers:       constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
				constants.Iso2Usd_IsolatedMarket,
			},
			expectedSubaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD, // No changes as transfer should fail.
				constants.Bob_Num0_1ISO2_LONG_10_000USD,
			},
			expectedCollateralPoolBalances: map[string]int64{
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.Iso2Usd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // No changes to collateral pools as transfer fails.
			},
			expectedErr: "insufficient funds",
		},
		`Can transfer from isolated subaccount to isolated subaccount in the same isolated markets, no
		coins are sent`: {
			subaccounts: []satypes.Subaccount{
				constants.Alice_Num0_1ISO_LONG_10_000USD,
				constants.Bob_Num0_1ISO_LONG_10_000USD,
			},
			collateralPoolBalances: map[string]int64{
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // $10,000 USDC
			},
			senderSubaccountId:   constants.Alice_Num0,
			receiverSubaccountId: constants.Bob_Num0,
			quantums:             100_000_000, // $100
			liquidityTiers:       constants.LiquidityTiers,
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
				constants.IsoUsd_IsolatedMarket,
			},
			expectedSubaccounts: []satypes.Subaccount{
				testutil.ChangeUsdcBalance(constants.Alice_Num0_1ISO_LONG_10_000USD, -100_000_000),
				testutil.ChangeUsdcBalance(constants.Bob_Num0_1ISO_LONG_10_000USD, 100_000_000),
			},
			expectedCollateralPoolBalances: map[string]int64{
				authtypes.NewModuleAddress(
					satypes.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.Params.Id),
				).String(): 10_000_000_000, // No change
			},
			expectedErr: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Configure the test application.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *assettypes.GenesisState) {
						genesisState.Assets = []assettypes.Asset{
							*constants.Usdc,
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						// If the collateral pool address is already in bank genesis state, update it.
						foundPools := make(map[string]struct{})
						for i, bal := range genesisState.Balances {
							usdcBal, exists := tc.collateralPoolBalances[bal.Address]
							if exists {
								foundPools[bal.Address] = struct{}{}
								genesisState.Balances[i] = banktypes.Balance{
									Address: bal.Address,
									Coins: sdktypes.Coins{
										sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(usdcBal)),
									},
								}
							}
						}
						// If the collateral pool address is not in the bank genesis state, add it.
						for addr, bal := range tc.collateralPoolBalances {
							_, exists := foundPools[addr]
							if exists {
								continue
							}
							genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
								Address: addr,
								Coins: sdktypes.Coins{
									sdktypes.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(bal)),
								},
							})
						}
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *perptypes.GenesisState) {
						genesisState.Params = constants.PerpetualsGenesisParams
						genesisState.LiquidityTiers = tc.liquidityTiers
						genesisState.Perpetuals = tc.perpetuals
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = tc.subaccounts
					},
				)
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *feetiertypes.GenesisState) {
						genesisState.Params = constants.PerpetualFeeParamsNoFee
					},
				)
				return genesis
			}).Build()

			ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

			// Send transfer message.
			var msg proto.Message
			transferMsg := sendingtypes.MsgCreateTransfer{
				Transfer: &sendingtypes.Transfer{
					Sender:    tc.senderSubaccountId,
					Recipient: tc.receiverSubaccountId,
					AssetId:   constants.Usdc.Id,
					Amount:    tc.quantums,
				},
			}
			msg = &transferMsg
			for _, checkTx := range testapp.MustMakeCheckTxsWithSdkMsg(
				ctx,
				tApp.App,
				testapp.MustMakeCheckTxOptions{
					AccAddressForSigning: tc.senderSubaccountId.Owner,
					Gas:                  1000000,
					FeeAmt:               constants.TestFeeCoins_5Cents,
				},
				msg,
			) {
				resp := tApp.CheckTx(checkTx)
				require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)
			}

			// Verify test expectations, subaccount and collateral pools should be updated.
			ctx = tApp.AdvanceToBlock(
				3,
				testapp.AdvanceToBlockOptions{
					ValidateFinalizeBlock: func(
						ctx sdktypes.Context,
						request abcitypes.RequestFinalizeBlock,
						response abcitypes.ResponseFinalizeBlock,
					) (haltchain bool) {
						execResult := response.TxResults[1]
						if tc.expectedErr != "" {
							// Note the first TX is MsgProposedOperations, the second is all other TXs.
							execResult := response.TxResults[1]
							require.True(t, execResult.IsErr())
							require.Equal(t, sdkerrors.ErrInsufficientFunds.ABCICode(), execResult.Code)
							require.Contains(t, execResult.Log, tc.expectedErr)
						} else {
							require.False(t, execResult.IsErr())
						}
						return false
					},
				},
			)
			for _, expectedSubaccount := range tc.expectedSubaccounts {
				require.Equal(
					t,
					expectedSubaccount,
					tApp.App.SubaccountsKeeper.GetSubaccount(ctx, *expectedSubaccount.Id),
				)
			}
			for addr, expectedBalance := range tc.expectedCollateralPoolBalances {
				require.Equal(
					t,
					sdkmath.NewIntFromBigInt(big.NewInt(expectedBalance)),
					tApp.App.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(addr), constants.Usdc.Denom).Amount,
				)
			}
		})
	}
}
