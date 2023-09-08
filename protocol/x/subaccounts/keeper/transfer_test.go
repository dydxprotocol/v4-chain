package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	"math"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	auth_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/auth"
	bank_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	sample_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestWithdrawFundsFromSubaccountToAccount_DepositFundsFromAccountToSubaccount_Success(t *testing.T) {
	tests := map[string]struct {
		testTransferFundToAccount bool
		asset                     asstypes.Asset

		// Subaccount state.
		assetPositions []*types.AssetPosition

		// Module account state.
		subaccountModuleAccBalance *big.Int
		accAddressBalance          *big.Int

		// Transfer details.
		quantums *big.Int

		// Expectations.
		expectedAssetPositions              []*types.AssetPosition
		expectedQuoteBalance                *big.Int
		expectedSubaccountsModuleAccBalance *big.Int
		expectedAccAddressBalance           *big.Int
	}{
		"WithdrawFundsFromSubaccountToAccount: send from subaccount to an account address": {
			testTransferFundToAccount:           true,
			asset:                               *constants.Usdc,
			accAddressBalance:                   big.NewInt(2500),
			subaccountModuleAccBalance:          big.NewInt(600),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedQuoteBalance:                big.NewInt(0),    // 500 - 500
			expectedSubaccountsModuleAccBalance: big.NewInt(100),  // 600 - 100
			expectedAccAddressBalance:           big.NewInt(3000), // 500 + 2500
		},
		"WithdrawFundsFromSubaccountToAccount: DenomeExponent > AtomicResolution, quantums rounded down": {
			testTransferFundToAccount: true,
			asset: asstypes.Asset{
				Id:               0,
				Symbol:           "USDC",
				Denom:            asstypes.AssetUsdc.Denom,
				DenomExponent:    int32(-3), // $1 = 1_000 coin unit.
				HasMarket:        false,
				MarketId:         uint32(0),
				AtomicResolution: int32(-4), // $1 = 10_000 quantums
			},
			accAddressBalance:          big.NewInt(2_500),  // $2.5
			subaccountModuleAccBalance: big.NewInt(10_000), // $10
			quantums:                   big.NewInt(20_001), // $2.0001, only $2 transfered.
			assetPositions: keepertest.CreateUsdcAssetPosition(
				big.NewInt(30_001),
			), // $3.0001
			expectedQuoteBalance:                big.NewInt(10_001), // $1.0001, untransfered $0.0001 remains.
			expectedSubaccountsModuleAccBalance: big.NewInt(8_000),  // $8
			expectedAccAddressBalance:           big.NewInt(4_500),  // $2.5 + $2
		},
		"DepositFundsFromAccountToSubaccount: send from account to subaccount": {
			testTransferFundToAccount:           false,
			asset:                               *constants.Usdc,
			subaccountModuleAccBalance:          big.NewInt(200),
			accAddressBalance:                   big.NewInt(2000),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(150)),
			expectedQuoteBalance:                big.NewInt(650),  // 150 + 500
			expectedSubaccountsModuleAccBalance: big.NewInt(700),  // 200 + 500
			expectedAccAddressBalance:           big.NewInt(1500), // 2000 - 500
		},
		"DepositFundsFromAccountToSubaccount: send from account to subaccount, DenomExponent < AtomicResolution": {
			testTransferFundToAccount: false,
			asset: asstypes.Asset{
				Id:               0,
				Symbol:           "USDC",
				Denom:            asstypes.AssetUsdc.Denom,
				DenomExponent:    int32(-4), // $1 = 10_000 coin unit.
				HasMarket:        false,
				MarketId:         uint32(0),
				AtomicResolution: int32(-3), // $1 = 1_000 quantums
			},
			subaccountModuleAccBalance:          big.NewInt(20_000),                                    // $2
			accAddressBalance:                   big.NewInt(90_000),                                    // $9
			quantums:                            big.NewInt(5_021),                                     // $5.021
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(1_050)), // $1.05
			expectedQuoteBalance:                big.NewInt(6_071),                                     // $1.05 + $5.021
			expectedSubaccountsModuleAccBalance: big.NewInt(70_210),                                    // $2 + $5.021
			expectedAccAddressBalance:           big.NewInt(39_790),                                    // $9 - $5.021
		},
		"DepositFundsFromAccountToSubaccount: new balance reaches max int64": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: new(big.Int).SetUint64(math.MaxUint64 - 100),
			quantums:                   big.NewInt(500),
			assetPositions: keepertest.CreateUsdcAssetPosition(
				new(big.Int).SetUint64(math.MaxUint64 - 100),
			),
			expectedQuoteBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				big.NewInt(400),
			),
			expectedSubaccountsModuleAccBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				big.NewInt(400),
			),
			expectedAccAddressBalance: big.NewInt(0),
		},

		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
		// TODO(CORE-169): Add tests for when the input quantums is rounded down to
		// a integer denom amount.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, _, accountKeeper, bankKeeper, assetsKeeper, _ := keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})

			// Set up test account address.
			addressStr := sample_testutil.AccAddress()
			testAccAddress, err := sdk.AccAddressFromBech32(addressStr)
			require.NoError(t, err)

			testAcc := authtypes.NewBaseAccount(testAccAddress, nil, 0, 0)
			accountKeeper.SetAccount(ctx, testAcc)

			if tc.accAddressBalance.Sign() > 0 {
				// Mint asset in the receipt/sender account address for transfer.
				err := bank_testutil.FundAccount(
					ctx,
					testAccAddress,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.accAddressBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.subaccountModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					types.ModuleName,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			_, err = assetsKeeper.CreateAsset(
				ctx,
				tc.asset.Symbol,
				tc.asset.Denom,
				tc.asset.DenomExponent,
				tc.asset.HasMarket,
				tc.asset.MarketId,
				tc.asset.AtomicResolution,
			)
			require.NoError(t, err)

			subaccount := createNSubaccount(keeper, ctx, 1)[0]
			subaccount.AssetPositions = tc.assetPositions

			keeper.SetSubaccount(ctx, subaccount)

			// Test either WithdrawFundsFromSubaccountToAccount or DepositFundsFromAccountToSubaccount.
			if tc.testTransferFundToAccount {
				err = keeper.WithdrawFundsFromSubaccountToAccount(
					ctx,
					*subaccount.Id,
					testAccAddress,
					tc.asset.Id,
					tc.quantums,
				)
			} else {
				err = keeper.DepositFundsFromAccountToSubaccount(
					ctx,
					testAccAddress,
					*subaccount.Id,
					tc.asset.Id,
					tc.quantums,
				)
			}

			require.NoError(t, err)

			// Check the subaccount has been updated as expected.
			updatedSubaccount := keeper.GetSubaccount(ctx, *subaccount.Id)
			if tc.expectedAssetPositions != nil {
				require.Equal(t,
					tc.expectedAssetPositions,
					updatedSubaccount.AssetPositions,
				)
			}
			require.Equal(t,
				tc.expectedQuoteBalance,
				updatedSubaccount.GetUsdcPosition(),
			)

			// Check the subaccount module balance.
			subaccountsModuleAccBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(types.ModuleName),
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.expectedSubaccountsModuleAccBalance)),
				subaccountsModuleAccBalance,
			)

			// Check the test account balance has been updated as expected.
			testAccountBalance := bankKeeper.GetBalance(
				ctx, testAccAddress,
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.expectedAccAddressBalance)),
				testAccountBalance,
			)
		})
	}
}

func TestWithdrawFundsFromSubaccountToAccount_DepositFundsFromAccountToSubaccount_Failure(t *testing.T) {
	tests := map[string]struct {
		skipSetUpUsdc             bool
		testTransferFundToAccount bool
		asset                     asstypes.Asset

		// Subaccount state.
		assetPositions []*types.AssetPosition

		// Module account state.
		subaccountModuleAccBalance *big.Int
		accAddressBalance          *big.Int

		// Transfer details
		quantums *big.Int

		// Expectations.
		expectedErr error
	}{
		"WithdrawFundsFromSubaccountToAccount: subaccount does not have enough balance to transfer": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(1000),
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(100)),
			expectedErr:                types.ErrFailedToUpdateSubaccounts,
		},
		"WithdrawFundsFromSubaccountToAccount: subaccounts module account does not have enough balance": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(400),
			accAddressBalance:          big.NewInt(5000),
			quantums:                   big.NewInt(500),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                sdkerrors.ErrInsufficientFunds,
		},
		"WithdrawFundsFromSubaccountToAccount: transfer quantums is zero": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(0),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                types.ErrAssetTransferQuantumsNotPositive,
		},
		"WithdrawFundsFromSubaccountToAccount: transfer quantums is negative": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(-100),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                types.ErrAssetTransferQuantumsNotPositive,
		},
		"WithdrawFundsFromSubaccountToAccount: do not support assets other than USDC": {
			testTransferFundToAccount:  true,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.BtcUsd,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                types.ErrAssetTransferThroughBankNotImplemented,
		},
		"WithdrawFundsFromSubaccountToAccount: asset ID doesn't exist": {
			testTransferFundToAccount:  true,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.Usdc,
			skipSetUpUsdc:              true,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                asstypes.ErrAssetDoesNotExist,
		},
		"DepositFundsFromAccountToSubaccount: fee-collector does not have enough balance to transfer": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(100),
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(2000),
			quantums:                   big.NewInt(500),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                sdkerrors.ErrInsufficientFunds,
		},
		"DepositFundsFromAccountToSubaccount: transfer quantums is zero": {
			testTransferFundToAccount:  false,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(0),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                types.ErrAssetTransferQuantumsNotPositive,
		},
		"DepositFundsFromAccountToSubaccount: do not support assets other than USDC": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.BtcUsd,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                types.ErrAssetTransferThroughBankNotImplemented,
		},
		"DepositFundsFromAccountToSubaccount: failure, asset ID doesn't exist": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(500),
			skipSetUpUsdc:              true,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                asstypes.ErrAssetDoesNotExist,
		},
		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, _, accountKeeper, bankKeeper, assetsKeeper, _ := keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})

			// Set up test account address.
			addressStr := sample_testutil.AccAddress()
			testAccAddress, err := sdk.AccAddressFromBech32(addressStr)
			require.NoError(t, err)

			testAcc := authtypes.NewBaseAccount(testAccAddress, nil, 0, 0)
			accountKeeper.SetAccount(ctx, testAcc)

			if tc.accAddressBalance.Sign() > 0 {
				// Mint asset in the receipt/sender account address for transfer.
				err := bank_testutil.FundAccount(
					ctx,
					testAccAddress,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.accAddressBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.subaccountModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					types.ModuleName,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if !tc.skipSetUpUsdc {
				// Always create USDC as the first asset unless specificed to skip.
				_, err = assetsKeeper.CreateAsset(
					ctx,
					constants.Usdc.Symbol,
					constants.Usdc.Denom,
					constants.Usdc.DenomExponent,
					constants.Usdc.HasMarket,
					constants.Usdc.MarketId,
					constants.Usdc.AtomicResolution,
				)
				require.NoError(t, err)
			}

			if tc.asset.Denom != constants.Usdc.Denom {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					tc.asset.Symbol,
					tc.asset.Denom,
					tc.asset.DenomExponent,
					tc.asset.HasMarket,
					tc.asset.MarketId,
					tc.asset.AtomicResolution,
				)
				require.NoError(t, err)
			}

			subaccount := createNSubaccount(keeper, ctx, 1)[0]
			subaccount.AssetPositions = tc.assetPositions

			keeper.SetSubaccount(ctx, subaccount)

			// Test either WithdrawFundsFromSubaccountToAccount or DepositFundsFromAccountToSubaccount.
			if tc.testTransferFundToAccount {
				err = keeper.WithdrawFundsFromSubaccountToAccount(
					ctx,
					*subaccount.Id,
					testAccAddress,
					tc.asset.Id,
					tc.quantums,
				)
			} else {
				err = keeper.DepositFundsFromAccountToSubaccount(
					ctx,
					testAccAddress,
					*subaccount.Id,
					tc.asset.Id,
					tc.quantums,
				)
			}

			require.ErrorIs(t,
				err,
				tc.expectedErr,
			)

			// Check the subaccount balance stays the same.
			updatedSubaccount := keeper.GetSubaccount(ctx, *subaccount.Id)

			require.Equal(t,
				tc.assetPositions[0].GetBigQuantums(),
				updatedSubaccount.GetUsdcPosition(),
			)

			// Check the subaccount module balance stays the same.
			subaccountsModuleAccBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(types.ModuleName),
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
				subaccountsModuleAccBalance,
			)

			// Check the test account balance stays the same.
			testAccountBalance := bankKeeper.GetBalance(
				ctx, testAccAddress,
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.accAddressBalance)),
				testAccountBalance,
			)
		})
	}
}

func TestTransferFundsFromSubaccountToModule_TransferFundsFromModuleToSubaccount(t *testing.T) {
	tests := map[string]struct {
		skipSetUpUsdc            bool
		testTransferFundToModule bool
		asset                    asstypes.Asset

		// Subaccount state.
		assetPositions []*types.AssetPosition

		// Module account state.
		subaccountModuleAccBalance *big.Int
		otherModuleAccBalance      *big.Int

		// Transfer details
		quantums        *big.Int
		otherModuleName string

		// Expectations.
		expectedAssetPositions              []*types.AssetPosition
		expectedErr                         error
		expectedQuoteBalance                *big.Int
		expectedSubaccountsModuleAccBalance *big.Int
		expectedOtherModuleAccBalance       *big.Int
	}{
		"TransferFundsFromSubaccountToModule: successfully send from subaccount to fee-collector module account": {
			testTransferFundToModule:            true,
			otherModuleName:                     authtypes.FeeCollectorName,
			asset:                               *constants.Usdc,
			otherModuleAccBalance:               big.NewInt(2500),
			subaccountModuleAccBalance:          big.NewInt(600),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedQuoteBalance:                big.NewInt(0),    // 500 - 5
			expectedSubaccountsModuleAccBalance: big.NewInt(100),  // 600 - 100
			expectedOtherModuleAccBalance:       big.NewInt(3000), // 500 + 2500
		},
		"TransferFundsFromSubaccountToModule: failure, subaccount does not have enough balance to transfer": {
			testTransferFundToModule:            true,
			otherModuleName:                     authtypes.FeeCollectorName,
			otherModuleAccBalance:               big.NewInt(0),
			asset:                               *constants.Usdc,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(100)),
			expectedErr:                         types.ErrFailedToUpdateSubaccounts,
			expectedQuoteBalance:                big.NewInt(100),
			expectedSubaccountsModuleAccBalance: big.NewInt(500),
			expectedOtherModuleAccBalance:       big.NewInt(0),
		},
		"TransferFundsFromSubaccountToModule: failure, subaccounts module account does not have enough balance": {
			testTransferFundToModule:            true,
			otherModuleName:                     authtypes.FeeCollectorName,
			asset:                               *constants.Usdc,
			otherModuleAccBalance:               big.NewInt(0),
			subaccountModuleAccBalance:          big.NewInt(400),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(400),
			expectedOtherModuleAccBalance:       big.NewInt(0),
		},
		"TransferFundsFromSubaccountToModule: failure, transfer quantums is zero": {
			testTransferFundToModule:            true,
			otherModuleName:                     authtypes.FeeCollectorName,
			asset:                               *constants.Usdc,
			otherModuleAccBalance:               big.NewInt(2500),
			subaccountModuleAccBalance:          big.NewInt(600),
			quantums:                            big.NewInt(0),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(600),
			expectedOtherModuleAccBalance:       big.NewInt(2500),
			expectedErr:                         types.ErrAssetTransferQuantumsNotPositive,
		},
		"TransferFundsFromSubaccountToModule: failure, do not support assets other than USDC": {
			testTransferFundToModule:            true,
			otherModuleName:                     authtypes.FeeCollectorName,
			otherModuleAccBalance:               big.NewInt(500),
			asset:                               *constants.BtcUsd,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                         types.ErrAssetTransferThroughBankNotImplemented,
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(500),
			expectedOtherModuleAccBalance:       big.NewInt(500),
		},
		"TransferFundsFromSubaccountToModule: failure, asset ID doesn't exist": {
			testTransferFundToModule:            true,
			otherModuleName:                     authtypes.FeeCollectorName,
			otherModuleAccBalance:               big.NewInt(500),
			asset:                               *constants.Usdc,
			skipSetUpUsdc:                       true,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                         asstypes.ErrAssetDoesNotExist,
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(500),
			expectedOtherModuleAccBalance:       big.NewInt(500),
		},
		"TransferFundsFromModuleToSubaccount: successfully send from fee-collector module account to subaccount": {
			testTransferFundToModule:            false,
			otherModuleName:                     authtypes.FeeCollectorName,
			asset:                               *constants.Usdc,
			subaccountModuleAccBalance:          big.NewInt(200),
			otherModuleAccBalance:               big.NewInt(2000),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(150)),
			expectedQuoteBalance:                big.NewInt(650),  // 150 + 500
			expectedSubaccountsModuleAccBalance: big.NewInt(700),  // 200 + 500
			expectedOtherModuleAccBalance:       big.NewInt(1500), // 2000 - 500
		},
		"TransferFundsFromModuleToSubaccount: failure, fee-collector does not have enough balance to transfer": {
			testTransferFundToModule:            false,
			otherModuleName:                     authtypes.FeeCollectorName,
			otherModuleAccBalance:               big.NewInt(100),
			asset:                               *constants.Usdc,
			subaccountModuleAccBalance:          big.NewInt(2000),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(2000),
			expectedOtherModuleAccBalance:       big.NewInt(100),
		},
		"TransferFundsFromModuleToSubaccount: failure, transfer quantums is zero": {
			testTransferFundToModule:            false,
			otherModuleName:                     authtypes.FeeCollectorName,
			asset:                               *constants.Usdc,
			otherModuleAccBalance:               big.NewInt(2500),
			subaccountModuleAccBalance:          big.NewInt(600),
			quantums:                            big.NewInt(0),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(600),
			expectedOtherModuleAccBalance:       big.NewInt(2500),
			expectedErr:                         types.ErrAssetTransferQuantumsNotPositive,
		},
		"TransferFundsFromModuleToSubaccount: failure, do not support assets other than USDC": {
			testTransferFundToModule:            false,
			otherModuleName:                     authtypes.FeeCollectorName,
			otherModuleAccBalance:               big.NewInt(500),
			asset:                               *constants.BtcUsd,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                         types.ErrAssetTransferThroughBankNotImplemented,
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(500),
			expectedOtherModuleAccBalance:       big.NewInt(500),
		},
		"TransferFundsFromModuleToSubaccount: failure, asset ID doesn't exist": {
			testTransferFundToModule:            false,
			otherModuleName:                     authtypes.FeeCollectorName,
			otherModuleAccBalance:               big.NewInt(500),
			asset:                               *constants.Usdc,
			skipSetUpUsdc:                       true,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(500),
			assetPositions:                      keepertest.CreateUsdcAssetPosition(big.NewInt(500)),
			expectedErr:                         asstypes.ErrAssetDoesNotExist,
			expectedQuoteBalance:                big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(500),
			expectedOtherModuleAccBalance:       big.NewInt(500),
		},
		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, _, accountKeeper, bankKeeper, assetsKeeper, _ := keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Create a mint module for funding test accounts and module accounts.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, minttypes.ModuleName, []string{authtypes.Minter})
			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})
			// Set up receiver module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, tc.otherModuleName, []string{authtypes.Minter})

			// This currently assumes the 1 base denom = 1 base quantum.
			// TODO(DEC-714): Implement conversion of assets between `assets` and
			// `bank` modules
			bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
				Base:    tc.asset.Denom,
				Display: tc.asset.Denom,
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    tc.asset.Denom,
						Exponent: 0,
					},
				},
			})

			// Mint asset in the receipt/sender module account for transfer.
			if tc.otherModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					tc.otherModuleName,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.otherModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.subaccountModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					types.ModuleName,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			// Always create USDC as the first asset.
			if !tc.skipSetUpUsdc {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					constants.Usdc.Symbol,
					constants.Usdc.Denom,
					constants.Usdc.DenomExponent,
					constants.Usdc.HasMarket,
					constants.Usdc.MarketId,
					constants.Usdc.AtomicResolution,
				)
				require.NoError(t, err)
			}

			if tc.asset.Denom != constants.Usdc.Denom {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					tc.asset.Symbol,
					tc.asset.Denom,
					tc.asset.DenomExponent,
					tc.asset.HasMarket,
					tc.asset.MarketId,
					tc.asset.AtomicResolution,
				)
				require.NoError(t, err)
			}

			subaccount := createNSubaccount(keeper, ctx, 1)[0]
			subaccount.AssetPositions = tc.assetPositions

			keeper.SetSubaccount(ctx, subaccount)

			// Test either TransferFundsFromSubaccountToModule or TransferFundsFromModuleToSubaccount.
			var err error
			if tc.testTransferFundToModule {
				err = keeper.TransferFundsFromSubaccountToModule(
					ctx,
					*subaccount.Id,
					authtypes.FeeCollectorName,
					tc.asset.Id,
					tc.quantums,
				)
			} else {
				err = keeper.TransferFundsFromModuleToSubaccount(
					ctx,
					authtypes.FeeCollectorName,
					*subaccount.Id,
					tc.asset.Id,
					tc.quantums,
				)
			}

			if tc.expectedErr != nil {
				require.ErrorIs(t,
					err,
					tc.expectedErr,
				)
			} else {
				require.NoError(t, err)
			}

			// Check the subaccount has been updated as expected.
			updatedSubaccount := keeper.GetSubaccount(ctx, *subaccount.Id)
			if tc.expectedAssetPositions != nil {
				require.Equal(t,
					tc.expectedAssetPositions,
					updatedSubaccount.AssetPositions,
				)
			}
			require.Equal(t,
				tc.expectedQuoteBalance,
				updatedSubaccount.GetUsdcPosition(),
			)

			// Check the subaccount module balance.
			subaccountsModuleAccBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(types.ModuleName),
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.expectedSubaccountsModuleAccBalance)),
				subaccountsModuleAccBalance,
			)

			// Check the other module account balance has been updated as expected.
			toModuleBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(tc.otherModuleName),
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.expectedOtherModuleAccBalance)),
				toModuleBalance,
			)
		})
	}
}

func TestTransferFeesToFeeCollectorModule(t *testing.T) {
	tests := map[string]struct {
		skipSetUpUsdc bool

		// Module account state.
		subaccountModuleAccBalance *big.Int
		feeModuleAccBalance        *big.Int

		// Transfer details.
		asset    asstypes.Asset
		quantums *big.Int

		// Expectations.
		expectedErr                         error
		expectedSubaccountsModuleAccBalance *big.Int
		expectedFeeModuleAccBalance         *big.Int
	}{
		"success - send to fee-collector module account": {
			asset:                               *constants.Usdc,
			feeModuleAccBalance:                 big.NewInt(2500),
			subaccountModuleAccBalance:          big.NewInt(600),
			quantums:                            big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(100),  // 600 - 500
			expectedFeeModuleAccBalance:         big.NewInt(3000), // 500 + 2500
		},
		"success - quantums is zero": {
			asset:                               *constants.Usdc,
			feeModuleAccBalance:                 big.NewInt(2500),
			subaccountModuleAccBalance:          big.NewInt(600),
			quantums:                            big.NewInt(0),
			expectedSubaccountsModuleAccBalance: big.NewInt(600),  // 600
			expectedFeeModuleAccBalance:         big.NewInt(2500), // 2500
		},
		"failure - subaccounts module does not have sufficient funds": {
			asset:                               *constants.Usdc,
			feeModuleAccBalance:                 big.NewInt(2500),
			subaccountModuleAccBalance:          big.NewInt(300),
			quantums:                            big.NewInt(500),
			expectedSubaccountsModuleAccBalance: big.NewInt(300),
			expectedFeeModuleAccBalance:         big.NewInt(2500),
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
		},
		"failure - asset ID doesn't exist": {
			feeModuleAccBalance:                 big.NewInt(1500),
			skipSetUpUsdc:                       true,
			asset:                               *constants.Usdc,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(500),
			expectedErr:                         asstypes.ErrAssetDoesNotExist,
			expectedSubaccountsModuleAccBalance: big.NewInt(500),
			expectedFeeModuleAccBalance:         big.NewInt(1500),
		},
		"failure - asset other than USDC not supported": {
			feeModuleAccBalance:                 big.NewInt(1500),
			asset:                               *constants.BtcUsd,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(500),
			expectedErr:                         types.ErrAssetTransferThroughBankNotImplemented,
			expectedSubaccountsModuleAccBalance: big.NewInt(500),
			expectedFeeModuleAccBalance:         big.NewInt(1500),
		},
		"success - transfer quantums is negative": {
			feeModuleAccBalance:                 big.NewInt(1500),
			asset:                               *constants.Usdc,
			subaccountModuleAccBalance:          big.NewInt(500),
			quantums:                            big.NewInt(-500),
			expectedSubaccountsModuleAccBalance: big.NewInt(1000),
			expectedFeeModuleAccBalance:         big.NewInt(1000),
		},
		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, _, accountKeeper, bankKeeper, assetsKeeper, _ := keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})
			// Set up receiver module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, authtypes.FeeCollectorName, []string{authtypes.Minter})

			// This currently assumes the 1 base denom = 1 base quantum.
			// TODO(DEC-714): Implement conversion of assets between `assets` and
			// `bank` modules
			bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
				Base:    tc.asset.Denom,
				Display: tc.asset.Denom,
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    tc.asset.Denom,
						Exponent: 0,
					},
				},
			})

			// Mint asset in the receipt/sender module account for transfer.
			if tc.feeModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					authtypes.FeeCollectorName,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.feeModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.subaccountModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					types.ModuleName,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			// Always create USDC as the first asset.
			if !tc.skipSetUpUsdc {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					constants.Usdc.Symbol,
					constants.Usdc.Denom,
					constants.Usdc.DenomExponent,
					constants.Usdc.HasMarket,
					constants.Usdc.MarketId,
					constants.Usdc.AtomicResolution,
				)
				require.NoError(t, err)
			}

			if tc.asset.Denom != constants.Usdc.Denom {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					tc.asset.Symbol,
					tc.asset.Denom,
					tc.asset.DenomExponent,
					tc.asset.HasMarket,
					tc.asset.MarketId,
					tc.asset.AtomicResolution,
				)
				require.NoError(t, err)
			}

			err := keeper.TransferFeesToFeeCollectorModule(
				ctx,
				tc.asset.Id,
				tc.quantums,
			)

			if tc.expectedErr != nil {
				require.ErrorIs(t,
					err,
					tc.expectedErr,
				)
			} else {
				require.NoError(t, err)
			}

			// Check the subaccount module balance.
			subaccountsModuleAccBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(types.ModuleName),
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.expectedSubaccountsModuleAccBalance)),
				subaccountsModuleAccBalance,
			)

			// Check the fee module account balance has been updated as expected.
			toModuleBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdk.NewIntFromBigInt(tc.expectedFeeModuleAccBalance)),
				toModuleBalance,
			)
		})
	}
}

func TestTransferInsuranceFundPayments(t *testing.T) {
	tests := map[string]struct {
		skipSetUpUsdc bool

		// Module account state.
		subaccountModuleAccBalance int64
		insuranceFundBalance       int64

		// Transfer details.
		quantums *big.Int

		// Expectations.
		panics                              bool
		expectedErr                         error
		expectedSubaccountsModuleAccBalance int64
		expectedInsuranceFundBalance        int64
	}{
		"success - send to insurance fund module account": {
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          600,
			quantums:                            big.NewInt(500),
			expectedSubaccountsModuleAccBalance: 100,  // 600 - 500
			expectedInsuranceFundBalance:        3000, // 2500 + 500
		},
		"success - send from insurance fund module account": {
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          600,
			quantums:                            big.NewInt(-500),
			expectedSubaccountsModuleAccBalance: 1100, // 600 + 500
			expectedInsuranceFundBalance:        2000, // 2500 - 500
		},
		"success - can send zero payment": {
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          600,
			quantums:                            big.NewInt(0),
			expectedSubaccountsModuleAccBalance: 600,
			expectedInsuranceFundBalance:        2500,
		},
		"failure - subaccounts module does not have sufficient funds": {
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          300,
			quantums:                            big.NewInt(500),
			expectedSubaccountsModuleAccBalance: 300,
			expectedInsuranceFundBalance:        2500,
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
		},
		"failure - insurance fund does not have sufficient funds": {
			insuranceFundBalance:                300,
			subaccountModuleAccBalance:          2500,
			quantums:                            big.NewInt(-500),
			expectedSubaccountsModuleAccBalance: 2500,
			expectedInsuranceFundBalance:        300,
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
		},
		"panics - asset doesn't exist": {
			insuranceFundBalance:                1500,
			skipSetUpUsdc:                       true,
			subaccountModuleAccBalance:          500,
			quantums:                            big.NewInt(500),
			expectedErr:                         errorsmod.Wrap(asstypes.ErrAssetDoesNotExist, lib.Uint32ToString(0)),
			expectedSubaccountsModuleAccBalance: 500,
			expectedInsuranceFundBalance:        1500,
			panics:                              true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, _, accountKeeper, bankKeeper, assetsKeeper, _ := keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})
			// Set up insurance fund module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, clobtypes.InsuranceFundName, []string{})

			bankKeeper.SetDenomMetaData(ctx, banktypes.Metadata{
				Base:    constants.Usdc.Denom,
				Display: constants.Usdc.Denom,
				DenomUnits: []*banktypes.DenomUnit{
					{
						Denom:    constants.Usdc.Denom,
						Exponent: 0,
					},
				},
			})

			// Mint asset in the receipt/sender module account for transfer.
			if tc.insuranceFundBalance > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					clobtypes.InsuranceFundName,
					sdk.Coins{
						sdk.NewInt64Coin(constants.Usdc.Denom, tc.insuranceFundBalance),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.subaccountModuleAccBalance > 0 {
				err := bank_testutil.FundModuleAccount(
					ctx,
					types.ModuleName,
					sdk.Coins{
						sdk.NewInt64Coin(constants.Usdc.Denom, tc.subaccountModuleAccBalance),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if !tc.skipSetUpUsdc {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					constants.Usdc.Symbol,
					constants.Usdc.Denom,
					constants.Usdc.DenomExponent,
					constants.Usdc.HasMarket,
					constants.Usdc.MarketId,
					constants.Usdc.AtomicResolution,
				)
				require.NoError(t, err)
			}

			if tc.expectedErr != nil {
				if tc.panics {
					require.PanicsWithError(
						t,
						tc.expectedErr.Error(),
						func() {
							//nolint:errcheck
							keeper.TransferInsuranceFundPayments(ctx, tc.quantums)
						},
					)
				} else {
					require.ErrorIs(
						t,
						keeper.TransferInsuranceFundPayments(ctx, tc.quantums),
						tc.expectedErr,
					)
				}
			} else {
				require.NoError(t, keeper.TransferInsuranceFundPayments(ctx, tc.quantums))
			}

			// Check the subaccount module balance.
			subaccountsModuleAccBalance := bankKeeper.GetBalance(
				ctx,
				authtypes.NewModuleAddress(types.ModuleName),
				constants.Usdc.Denom,
			)
			require.Equal(
				t,
				sdk.NewInt64Coin(constants.Usdc.Denom, tc.expectedSubaccountsModuleAccBalance),
				subaccountsModuleAccBalance,
			)

			// Check the fee module account balance has been updated as expected.
			toModuleBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(clobtypes.InsuranceFundName),
				constants.Usdc.Denom,
			)
			require.Equal(t,
				sdk.NewInt64Coin(constants.Usdc.Denom, tc.expectedInsuranceFundBalance),
				toModuleBalance,
			)
		})
	}
}
