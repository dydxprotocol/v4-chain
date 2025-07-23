package keeper_test

import (
	"math"
	"math/big"
	"testing"

	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	auth_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/auth"
	bank_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	sample_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/sample"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	revsharetypes "github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

func TestWithdrawFundsFromSubaccountToAccount_DepositFundsFromAccountToSubaccount_Success(t *testing.T) {
	tests := map[string]struct {
		testTransferFundToAccount bool
		asset                     asstypes.Asset

		// Subaccount state.
		assetPositions     []*types.AssetPosition
		perpetualPositions []*types.PerpetualPosition

		// Module account state.
		subaccountModuleAccBalance *big.Int
		accAddressBalance          *big.Int

		collateralPoolAddr sdk.AccAddress

		// Transfer details.
		quantums *big.Int

		// Expectations.
		expectedAssetPositions              []*types.AssetPosition
		expectedQuoteBalance                *big.Int
		expectedSubaccountsModuleAccBalance *big.Int
		expectedAccAddressBalance           *big.Int
	}{
		"WithdrawFundsFromSubaccountToAccount: send from subaccount to an account address": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			collateralPoolAddr:                  types.ModuleAddress,
			expectedQuoteBalance:                big.NewInt(0),    // 500 - 500
			expectedSubaccountsModuleAccBalance: big.NewInt(100),  // 600 - 100
			expectedAccAddressBalance:           big.NewInt(3000), // 500 + 2500
		},
		"WithdrawFundsFromSubaccountToAccount: send from isolated subaccount to an account address": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
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
				DenomExponent:    int32(-6), // $1 = 1_000_000 coin unit.
				HasMarket:        false,
				MarketId:         uint32(0),
				AtomicResolution: int32(-7), // $1 = 10_000_000 quantums
			},
			accAddressBalance:          big.NewInt(2_500_000),  // $2.5
			subaccountModuleAccBalance: big.NewInt(10_000_000), // $10
			quantums:                   big.NewInt(20_000_001), // $2.0000001, only $2 transferred.
			assetPositions: testutil.CreateUsdcAssetPositions(
				big.NewInt(30_000_001),
			), // $3.0001
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			collateralPoolAddr:                  types.ModuleAddress,
			expectedQuoteBalance:                big.NewInt(10_000_001), // $1.0001, untransfered $0.0001 remains.
			expectedSubaccountsModuleAccBalance: big.NewInt(8_000_000),  // $8
			expectedAccAddressBalance:           big.NewInt(4_500_000),  // $2.5 + $2
		},
		"DepositFundsFromAccountToSubaccount: send from account to subaccount": {
			testTransferFundToAccount:  false,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(200),
			accAddressBalance:          big.NewInt(2000),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(150)),
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			collateralPoolAddr:                  types.ModuleAddress,
			expectedQuoteBalance:                big.NewInt(650),  // 150 + 500
			expectedSubaccountsModuleAccBalance: big.NewInt(700),  // 200 + 500
			expectedAccAddressBalance:           big.NewInt(1500), // 2000 - 500
		},
		"DepositFundsFromAccountToSubaccount: send from account to isolated subaccount": {
			testTransferFundToAccount:  false,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(200),
			accAddressBalance:          big.NewInt(2000),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(150)),
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
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
				DenomExponent:    int32(-6), // $1 = 1000_000 coin unit.
				HasMarket:        false,
				MarketId:         uint32(0),
				AtomicResolution: int32(-5), // $1 = 100_000 quantums
			},
			subaccountModuleAccBalance: big.NewInt(2_000_000),                                  // $2
			accAddressBalance:          big.NewInt(9_000_000),                                  // $9
			quantums:                   big.NewInt(502_100),                                    // $5.021
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(105_000)), // $1.05
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			collateralPoolAddr:                  types.ModuleAddress,
			expectedQuoteBalance:                big.NewInt(607_100),   // $1.05 + $5.021
			expectedSubaccountsModuleAccBalance: big.NewInt(7_021_000), // $2 + $5.021
			expectedAccAddressBalance:           big.NewInt(3_979_000), // $9 - $5.021
		},
		"DepositFundsFromAccountToSubaccount: new balance reaches max int64": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: new(big.Int).SetUint64(math.MaxUint64 - 100),
			quantums:                   big.NewInt(500),
			assetPositions: testutil.CreateUsdcAssetPositions(
				new(big.Int).SetUint64(math.MaxUint64 - 100),
			),
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			collateralPoolAddr: types.ModuleAddress,
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
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, _, _, _, _ :=
				keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})

			// Set up test account address.
			addressStr := sample_testutil.AccAddress()
			testAccAddress, err := sdk.AccAddressFromBech32(addressStr)
			require.NoError(t, err)

			testAcc := authtypes.NewBaseAccount(testAccAddress, nil, accountKeeper.NextAccountNumber(ctx), 0)
			accountKeeper.SetAccount(ctx, testAcc)

			if tc.accAddressBalance.Sign() > 0 {
				// Mint asset in the receipt/sender account address for transfer.
				err := bank_testutil.FundAccount(
					ctx,
					testAccAddress,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.accAddressBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			_, err = assetsKeeper.CreateAsset(
				ctx,
				tc.asset.Id,
				tc.asset.Symbol,
				tc.asset.Denom,
				tc.asset.DenomExponent,
				tc.asset.HasMarket,
				tc.asset.MarketId,
				tc.asset.AtomicResolution,
			)
			require.NoError(t, err)

			subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.AssetPositions = tc.assetPositions
			subaccount.PerpetualPositions = tc.perpetualPositions

			keeper.SetSubaccount(ctx, subaccount)

			if tc.subaccountModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					tc.collateralPoolAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

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
			subaccountsModuleAccBalance := bankKeeper.GetBalance(ctx, tc.collateralPoolAddr, tc.asset.Denom)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedSubaccountsModuleAccBalance)),
				subaccountsModuleAccBalance,
			)

			// Check the test account balance has been updated as expected.
			testAccountBalance := bankKeeper.GetBalance(
				ctx, testAccAddress,
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedAccAddressBalance)),
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
		assetPositions     []*types.AssetPosition
		perpetualPositions []*types.PerpetualPosition

		// Module account state.
		subaccountModuleAccBalance *big.Int
		accAddressBalance          *big.Int

		collateralPoolAddr sdk.AccAddress

		// Transfer details
		quantums *big.Int

		// Optional. Defaults to an arbitrary test account if nil
		optionalRecipient sdk.AccAddress

		// Expectations.
		expectedErr error
	}{
		"WithdrawFundsFromSubaccountToAccount: recipient is blocked address": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(0),
			subaccountModuleAccBalance: big.NewInt(0),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			collateralPoolAddr: types.ModuleAddress,
			optionalRecipient:  authtypes.NewModuleAddress(distrtypes.ModuleName),
			expectedErr:        sdkerrors.ErrUnauthorized,
		},
		"WithdrawFundsFromSubaccountToAccount: subaccount does not have enough balance to transfer": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(1000),
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(100)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                types.ErrFailedToUpdateSubaccounts,
		},
		"WithdrawFundsFromSubaccountToAccount: subaccounts module account does not have enough balance": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(400),
			accAddressBalance:          big.NewInt(5000),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                sdkerrors.ErrInsufficientFunds,
		},
		"WithdrawFundsFromSubaccountToAccount: isolated market subaccounts module account does not have enough balance": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(400),
			accAddressBalance:          big.NewInt(5000),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			perpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			expectedErr: sdkerrors.ErrInsufficientFunds,
		},
		"WithdrawFundsFromSubaccountToAccount: transfer quantums is zero": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(0),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                types.ErrAssetTransferQuantumsNotPositive,
		},
		"WithdrawFundsFromSubaccountToAccount: transfer quantums is negative": {
			testTransferFundToAccount:  true,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(-100),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                types.ErrAssetTransferQuantumsNotPositive,
		},
		"WithdrawFundsFromSubaccountToAccount: do not support assets other than USDC": {
			testTransferFundToAccount:  true,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.BtcUsd,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                types.ErrAssetTransferThroughBankNotImplemented,
		},
		"WithdrawFundsFromSubaccountToAccount: asset ID doesn't exist": {
			testTransferFundToAccount:  true,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.Usdc,
			skipSetUpUsdc:              true,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                asstypes.ErrAssetDoesNotExist,
		},
		"DepositFundsFromAccountToSubaccount: fee-collector does not have enough balance to transfer": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(100),
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(2000),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                sdkerrors.ErrInsufficientFunds,
		},
		"DepositFundsFromAccountToSubaccount: transfer quantums is zero": {
			testTransferFundToAccount:  false,
			asset:                      *constants.Usdc,
			accAddressBalance:          big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			quantums:                   big.NewInt(0),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                types.ErrAssetTransferQuantumsNotPositive,
		},
		"DepositFundsFromAccountToSubaccount: do not support assets other than USDC": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(500),
			asset:                      *constants.BtcUsd,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                types.ErrAssetTransferThroughBankNotImplemented,
		},
		"DepositFundsFromAccountToSubaccount: failure, asset ID doesn't exist": {
			testTransferFundToAccount:  false,
			accAddressBalance:          big.NewInt(500),
			skipSetUpUsdc:              true,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(500),
			quantums:                   big.NewInt(500),
			assetPositions:             testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			collateralPoolAddr:         types.ModuleAddress,
			expectedErr:                asstypes.ErrAssetDoesNotExist,
		},
		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, _, _, _, _ :=
				keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})

			// Set up test account address.
			var err error
			if tc.optionalRecipient == nil {
				addressStr := sample_testutil.AccAddress()
				tc.optionalRecipient, err = sdk.AccAddressFromBech32(addressStr)
				require.NoError(t, err)
			}

			testAcc := authtypes.NewBaseAccount(tc.optionalRecipient, nil, accountKeeper.NextAccountNumber(ctx), 0)
			accountKeeper.SetAccount(ctx, testAcc)

			if tc.accAddressBalance.Sign() > 0 {
				// Mint asset in the receipt/sender account address for transfer.
				err := bank_testutil.FundAccount(
					ctx,
					tc.optionalRecipient,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.accAddressBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if !tc.skipSetUpUsdc {
				// Always create USDC as the first asset unless specificed to skip.
				err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
				require.NoError(t, err)
			}

			if tc.asset.Denom != constants.Usdc.Denom {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					tc.asset.Id,
					tc.asset.Symbol,
					tc.asset.Denom,
					tc.asset.DenomExponent,
					tc.asset.HasMarket,
					tc.asset.MarketId,
					tc.asset.AtomicResolution,
				)
				require.NoError(t, err)
			}

			subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.AssetPositions = tc.assetPositions
			subaccount.PerpetualPositions = tc.perpetualPositions

			keeper.SetSubaccount(ctx, subaccount)

			if tc.subaccountModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					tc.collateralPoolAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			// Test either WithdrawFundsFromSubaccountToAccount or DepositFundsFromAccountToSubaccount.
			if tc.testTransferFundToAccount {
				err = keeper.WithdrawFundsFromSubaccountToAccount(
					ctx,
					*subaccount.Id,
					tc.optionalRecipient,
					tc.asset.Id,
					tc.quantums,
				)
			} else {
				err = keeper.DepositFundsFromAccountToSubaccount(
					ctx,
					tc.optionalRecipient,
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
			subaccountsModuleAccBalance := bankKeeper.GetBalance(ctx, tc.collateralPoolAddr, tc.asset.Denom)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
				subaccountsModuleAccBalance,
			)

			// Check the test account balance stays the same.
			testAccountBalance := bankKeeper.GetBalance(
				ctx, tc.optionalRecipient,
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.accAddressBalance)),
				testAccountBalance,
			)
		})
	}
}

func TestTransferFundsFromSubaccountToSubaccount_Success(t *testing.T) {
	tests := map[string]struct {
		// transfer details
		asset    asstypes.Asset
		quantums *big.Int

		// Subaccount state.
		senderAssetPositions        []*types.AssetPosition
		senderPerpetualPositions    []*types.PerpetualPosition
		recipientAssetPositions     []*types.AssetPosition
		recipientPerpetualPositions []*types.PerpetualPosition

		// Module account state.
		senderCollateralPoolBalance    *big.Int
		recipientCollateralPoolBalance *big.Int

		senderCollateralPoolAddr    sdk.AccAddress
		recipientCollateralPoolAddr sdk.AccAddress

		// Expectations.
		expectedSenderAssetPositions           []*types.AssetPosition
		expectedRecipientAssetPositions        []*types.AssetPosition
		expectedSenderQuoteBalance             *big.Int
		expectedRecipientQuoteBalance          *big.Int
		expectedSenderCollateralPoolBalance    *big.Int
		expectedRecipientCollateralPoolBalance *big.Int
	}{
		"Send USDC from non-isolated subaccount to non-isolated subaccount": {
			asset:                *constants.Usdc,
			quantums:             big.NewInt(500),
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			senderCollateralPoolBalance:            big.NewInt(1100), // 500 + 600
			recipientCollateralPoolBalance:         big.NewInt(1100), // same collateral pool, same balance
			senderCollateralPoolAddr:               types.ModuleAddress,
			recipientCollateralPoolAddr:            types.ModuleAddress,
			expectedRecipientAssetPositions:        testutil.CreateUsdcAssetPositions(big.NewInt(1100)),
			expectedSenderQuoteBalance:             big.NewInt(0),    // 500 - 500
			expectedRecipientQuoteBalance:          big.NewInt(1100), // 500 + 600
			expectedSenderCollateralPoolBalance:    big.NewInt(1100), // no changes to collateral pools
			expectedRecipientCollateralPoolBalance: big.NewInt(1100),
		},
		"Send USDC from isolated subaccount to non-isolated subaccount": {
			asset:                *constants.Usdc,
			quantums:             big.NewInt(500),
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			senderCollateralPoolBalance:    big.NewInt(600),
			recipientCollateralPoolBalance: big.NewInt(700),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			recipientCollateralPoolAddr:            types.ModuleAddress,
			expectedRecipientAssetPositions:        testutil.CreateUsdcAssetPositions(big.NewInt(1100)),
			expectedSenderQuoteBalance:             big.NewInt(0),    // 500 - 500
			expectedRecipientQuoteBalance:          big.NewInt(1100), // 500 + 600
			expectedSenderCollateralPoolBalance:    big.NewInt(100),  // 600 - 500
			expectedRecipientCollateralPoolBalance: big.NewInt(1200), // 700 + 500
		},
		"Send USDC from non-isolated subaccount to isolated subaccount": {
			asset:                *constants.Usdc,
			quantums:             big.NewInt(500),
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			senderCollateralPoolBalance:    big.NewInt(600),
			recipientCollateralPoolBalance: big.NewInt(700),
			senderCollateralPoolAddr:       types.ModuleAddress,
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			expectedRecipientAssetPositions:        testutil.CreateUsdcAssetPositions(big.NewInt(1100)),
			expectedSenderQuoteBalance:             big.NewInt(0),    // 500 - 500
			expectedRecipientQuoteBalance:          big.NewInt(1100), // 500 + 600
			expectedSenderCollateralPoolBalance:    big.NewInt(100),  // 600 - 500
			expectedRecipientCollateralPoolBalance: big.NewInt(1200), // 700 + 500
		},
		"Send USDC from isolated subaccount to isolated subaccount (same perp)": {
			asset:                *constants.Usdc,
			quantums:             big.NewInt(500),
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			senderCollateralPoolBalance:    big.NewInt(1100), // 500 + 600
			recipientCollateralPoolBalance: big.NewInt(1100), // same collateral pool, same balance
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			expectedRecipientAssetPositions:        testutil.CreateUsdcAssetPositions(big.NewInt(1100)),
			expectedSenderQuoteBalance:             big.NewInt(0),    // 500 - 500
			expectedRecipientQuoteBalance:          big.NewInt(1100), // 500 + 600
			expectedSenderCollateralPoolBalance:    big.NewInt(1100), // no changes to collateral pools
			expectedRecipientCollateralPoolBalance: big.NewInt(1100),
		},
		"Send USDC from isolated subaccount to isolated subaccount (different perp)": {
			asset:                *constants.Usdc,
			quantums:             big.NewInt(500),
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISO2Long,
			},
			senderCollateralPoolBalance:    big.NewInt(600),
			recipientCollateralPoolBalance: big.NewInt(700),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISO2Long.PerpetualId),
			),
			expectedRecipientAssetPositions:        testutil.CreateUsdcAssetPositions(big.NewInt(1100)),
			expectedSenderQuoteBalance:             big.NewInt(0),    // 500 - 500
			expectedRecipientQuoteBalance:          big.NewInt(1100), // 500 + 600
			expectedSenderCollateralPoolBalance:    big.NewInt(100),  // 600 - 500
			expectedRecipientCollateralPoolBalance: big.NewInt(1200), // 700 + 500
		},
		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
		// TODO(CORE-169): Add tests for when the input quantums is rounded down to
		// a integer denom amount.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, _, _, _, _ :=
				keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})

			// Set up test account address.
			addressStr := sample_testutil.AccAddress()
			testAccAddress, err := sdk.AccAddressFromBech32(addressStr)
			require.NoError(t, err)

			testAcc := authtypes.NewBaseAccount(testAccAddress, nil, accountKeeper.NextAccountNumber(ctx), 0)
			accountKeeper.SetAccount(ctx, testAcc)

			_, err = assetsKeeper.CreateAsset(
				ctx,
				tc.asset.Id,
				tc.asset.Symbol,
				tc.asset.Denom,
				tc.asset.DenomExponent,
				tc.asset.HasMarket,
				tc.asset.MarketId,
				tc.asset.AtomicResolution,
			)
			require.NoError(t, err)

			subaccounts := createNSubaccount(keeper, ctx, 2, big.NewInt(1_000))
			senderSubaccount := subaccounts[0]
			recipientSubaccount := subaccounts[1]
			senderSubaccount.AssetPositions = tc.senderAssetPositions
			senderSubaccount.PerpetualPositions = tc.senderPerpetualPositions
			recipientSubaccount.AssetPositions = tc.recipientAssetPositions
			recipientSubaccount.PerpetualPositions = tc.recipientPerpetualPositions

			keeper.SetSubaccount(ctx, senderSubaccount)
			keeper.SetSubaccount(ctx, recipientSubaccount)

			if tc.senderCollateralPoolBalance.Sign() > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					tc.senderCollateralPoolAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.senderCollateralPoolBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}
			if tc.recipientCollateralPoolBalance.Sign() > 0 &&
				!tc.recipientCollateralPoolAddr.Equals(tc.senderCollateralPoolAddr) {
				err := bank_testutil.FundAccount(
					ctx,
					tc.recipientCollateralPoolAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.recipientCollateralPoolBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			err = keeper.TransferFundsFromSubaccountToSubaccount(
				ctx,
				*senderSubaccount.Id,
				*recipientSubaccount.Id,
				tc.asset.Id,
				tc.quantums,
			)

			require.NoError(t, err)

			// Check the subaccount has been updated as expected.
			updatedSenderSubaccount := keeper.GetSubaccount(ctx, *senderSubaccount.Id)
			if tc.expectedSenderAssetPositions != nil {
				require.Equal(t,
					tc.expectedSenderAssetPositions,
					updatedSenderSubaccount.AssetPositions,
				)
			}
			require.Equal(t,
				tc.expectedSenderQuoteBalance,
				updatedSenderSubaccount.GetUsdcPosition(),
			)

			updatedRecipientSubaccount := keeper.GetSubaccount(ctx, *recipientSubaccount.Id)
			if tc.expectedRecipientAssetPositions != nil {
				require.Equal(t,
					tc.expectedRecipientAssetPositions,
					updatedRecipientSubaccount.AssetPositions,
				)
			}
			require.Equal(t,
				tc.expectedRecipientQuoteBalance,
				updatedRecipientSubaccount.GetUsdcPosition(),
			)

			// Check the subaccount module balance.
			senderCollateralPoolAddrBalance := bankKeeper.GetBalance(ctx, tc.senderCollateralPoolAddr, tc.asset.Denom)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedSenderCollateralPoolBalance)),
				senderCollateralPoolAddrBalance,
			)
			recipientCollateralPoolAddrBalance := bankKeeper.GetBalance(ctx, tc.recipientCollateralPoolAddr, tc.asset.Denom)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedRecipientCollateralPoolBalance)),
				recipientCollateralPoolAddrBalance,
			)
		})
	}
}

func TestTransferFundsFromSubaccountToSubaccount_Failure(t *testing.T) {
	tests := map[string]struct {
		skipSetUpUsdc bool
		asset         asstypes.Asset

		// Subaccount state.
		senderAssetPositions        []*types.AssetPosition
		senderPerpetualPositions    []*types.PerpetualPosition
		recipientAssetPositions     []*types.AssetPosition
		recipientPerpetualPositions []*types.PerpetualPosition

		// Module account state.
		senderCollateralPoolBalance    *big.Int
		recipientCollateralPoolBalance *big.Int

		senderCollateralPoolAddr    sdk.AccAddress
		recipientCollateralPoolAddr sdk.AccAddress

		// Transfer details
		quantums *big.Int

		// Expectations.
		expectedErr error
	}{
		"Send from non-isolated subaccount to non-isolated subaccount, sender does not have enough balance": {
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(100)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCShort,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCShort,
			},
			senderCollateralPoolBalance:    big.NewInt(1100),
			recipientCollateralPoolBalance: big.NewInt(1100),
			senderCollateralPoolAddr:       types.ModuleAddress,
			recipientCollateralPoolAddr:    types.ModuleAddress,
			quantums:                       big.NewInt(500),
			expectedErr:                    types.ErrFailedToUpdateSubaccounts,
		},
		"Send between isolated subaccounts (same perp), sender does not have enough balance": {
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(100)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOShort,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOShort,
			},
			senderCollateralPoolBalance:    big.NewInt(1100),
			recipientCollateralPoolBalance: big.NewInt(1100),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOShort.PerpetualId),
			),
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOShort.PerpetualId),
			),
			quantums:    big.NewInt(500),
			expectedErr: types.ErrFailedToUpdateSubaccounts,
		},
		"Send between isolated subaccounts (different perp), sender does not have enough balance": {
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(100)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOShort,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISO2Short,
			},
			senderCollateralPoolBalance:    big.NewInt(500),
			recipientCollateralPoolBalance: big.NewInt(600),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOShort.PerpetualId),
			),
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISO2Short.PerpetualId),
			),
			quantums:    big.NewInt(500),
			expectedErr: types.ErrFailedToUpdateSubaccounts,
		},
		"Send from isolated subaccount to non-isolated subaccount, sender does not have enough balance": {
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(100)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOShort,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCShort,
			},
			senderCollateralPoolBalance:    big.NewInt(500),
			recipientCollateralPoolBalance: big.NewInt(600),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOShort.PerpetualId),
			),
			recipientCollateralPoolAddr: types.ModuleAddress,
			quantums:                    big.NewInt(500),
			expectedErr:                 types.ErrFailedToUpdateSubaccounts,
		},
		"Send from non-isolated subaccount to isolated subaccount, collateral pool does not have enough balance": {
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			senderCollateralPoolBalance:    big.NewInt(100),
			recipientCollateralPoolBalance: big.NewInt(600),
			senderCollateralPoolAddr:       types.ModuleAddress,
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			quantums:    big.NewInt(500),
			expectedErr: sdkerrors.ErrInsufficientFunds,
		},
		"Send from isolated subaccount to non-isolated subaccount, collateral pool does not have enough balance": {
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneBTCLong,
			},
			senderCollateralPoolBalance:    big.NewInt(100),
			recipientCollateralPoolBalance: big.NewInt(600),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			recipientCollateralPoolAddr: types.ModuleAddress,
			quantums:                    big.NewInt(500),
			expectedErr:                 sdkerrors.ErrInsufficientFunds,
		},
		"Send between isolated subaccounts (different perp), collateral pool does not have enough balance": {
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISO2Long,
			},
			senderCollateralPoolBalance:    big.NewInt(100),
			recipientCollateralPoolBalance: big.NewInt(600),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISO2Long.PerpetualId),
			),
			quantums:    big.NewInt(500),
			expectedErr: sdkerrors.ErrInsufficientFunds,
		},
		"Do not support assets other than USDC": {
			asset:                *constants.BtcUsd,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISO2Long,
			},
			senderCollateralPoolBalance:    big.NewInt(100),
			recipientCollateralPoolBalance: big.NewInt(600),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISO2Long.PerpetualId),
			),
			quantums:    big.NewInt(500),
			expectedErr: types.ErrAssetTransferThroughBankNotImplemented,
		},
		"Asset ID doesn't exist": {
			skipSetUpUsdc:        true,
			asset:                *constants.Usdc,
			senderAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(500)),
			senderPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISOLong,
			},
			recipientAssetPositions: testutil.CreateUsdcAssetPositions(big.NewInt(600)),
			recipientPerpetualPositions: []*types.PerpetualPosition{
				&constants.PerpetualPosition_OneISO2Long,
			},
			senderCollateralPoolBalance:    big.NewInt(100),
			recipientCollateralPoolBalance: big.NewInt(600),
			senderCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
			),
			recipientCollateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISO2Long.PerpetualId),
			),
			quantums:    big.NewInt(500),
			expectedErr: asstypes.ErrAssetDoesNotExist,
		},
		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, _, _, _, _ :=
				keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})

			// Set up test account address.
			addressStr := sample_testutil.AccAddress()
			testAccAddress, err := sdk.AccAddressFromBech32(addressStr)
			require.NoError(t, err)

			testAcc := authtypes.NewBaseAccount(testAccAddress, nil, accountKeeper.NextAccountNumber(ctx), 0)
			accountKeeper.SetAccount(ctx, testAcc)

			if !tc.skipSetUpUsdc {
				// Always create USDC as the first asset unless specificed to skip.
				err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
				require.NoError(t, err)
			}

			if tc.asset.Denom != constants.Usdc.Denom {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					tc.asset.Id,
					tc.asset.Symbol,
					tc.asset.Denom,
					tc.asset.DenomExponent,
					tc.asset.HasMarket,
					tc.asset.MarketId,
					tc.asset.AtomicResolution,
				)
				require.NoError(t, err)
			}

			subaccounts := createNSubaccount(keeper, ctx, 2, big.NewInt(1_000))
			senderSubaccount := subaccounts[0]
			recipientSubaccount := subaccounts[1]
			senderSubaccount.AssetPositions = tc.senderAssetPositions
			senderSubaccount.PerpetualPositions = tc.senderPerpetualPositions
			recipientSubaccount.AssetPositions = tc.recipientAssetPositions
			recipientSubaccount.PerpetualPositions = tc.recipientPerpetualPositions

			keeper.SetSubaccount(ctx, senderSubaccount)
			keeper.SetSubaccount(ctx, recipientSubaccount)

			if tc.senderCollateralPoolBalance.Sign() > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					tc.senderCollateralPoolAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.senderCollateralPoolBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}
			if tc.recipientCollateralPoolBalance.Sign() > 0 &&
				!tc.recipientCollateralPoolAddr.Equals(tc.senderCollateralPoolAddr) {
				err := bank_testutil.FundAccount(
					ctx,
					tc.recipientCollateralPoolAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.recipientCollateralPoolBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			err = keeper.TransferFundsFromSubaccountToSubaccount(
				ctx,
				*senderSubaccount.Id,
				*recipientSubaccount.Id,
				tc.asset.Id,
				tc.quantums,
			)

			require.ErrorIs(t,
				err,
				tc.expectedErr,
			)

			// Check the subaccount has been updated as expected.
			updatedSenderSubaccount := keeper.GetSubaccount(ctx, *senderSubaccount.Id)
			require.Equal(t,
				tc.senderAssetPositions,
				updatedSenderSubaccount.AssetPositions,
			)

			updatedRecipientSubaccount := keeper.GetSubaccount(ctx, *recipientSubaccount.Id)
			require.Equal(t,
				tc.recipientAssetPositions,
				updatedRecipientSubaccount.AssetPositions,
			)

			// Check the subaccount module balance.
			senderCollateralPoolAddrBalance := bankKeeper.GetBalance(ctx, tc.senderCollateralPoolAddr, tc.asset.Denom)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.senderCollateralPoolBalance)),
				senderCollateralPoolAddrBalance,
			)
			recipientCollateralPoolAddrBalance := bankKeeper.GetBalance(ctx, tc.recipientCollateralPoolAddr, tc.asset.Denom)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.recipientCollateralPoolBalance)),
				recipientCollateralPoolAddrBalance,
			)
		})
	}
}

func TestDistributeFees(t *testing.T) {
	tests := map[string]struct {
		skipSetUpUsdc bool

		// Module account state.
		subaccountModuleAccBalance *big.Int
		feeModuleAccBalance        *big.Int
		marketMapperAccBalance     *big.Int

		// Transfer details.
		asset asstypes.Asset
		fill  clobtypes.FillForProcess

		collateralPoolAddr                  sdk.AccAddress
		affiliateRevShareAcctAddr           string
		marketMapperRevShareAcctAddr        string
		buySideOrderRouterRevShareAcctAddr  string
		sellSideOrderRouterRevShareAcctAddr string
		unconditionalRevShareAcctAddr       string
		revShare                            revsharetypes.RevSharesForFill

		// Expectations.
		expectedErr                             error
		expectedSubaccountsModuleAccBalance     *big.Int
		expectedFeeModuleAccBalance             *big.Int
		expectedMarketMapperAccBalance          *big.Int
		expectedBuySideOrderRouterRevShare      *big.Int
		expectedSellSideOrderRouterRevShare     *big.Int
		expectedAffiliateAccBalance             *big.Int
		expectedUnconditionalRevShareAccBalance *big.Int
	}{
		"success - send to fee-collector module account": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			marketMapperAccBalance:     big.NewInt(0),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(0),
				MarketId:                          uint32(0),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			collateralPoolAddr:                      types.ModuleAddress,
			affiliateRevShareAcctAddr:               "",
			marketMapperRevShareAcctAddr:            constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr:           "",
			expectedSubaccountsModuleAccBalance:     big.NewInt(100),  // 600 - 500
			expectedFeeModuleAccBalance:             big.NewInt(3000), // 500 + 2500
			revShare:                                revsharetypes.RevSharesForFill{},
			expectedMarketMapperAccBalance:          big.NewInt(0),
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
		},
		"success - send to fee-collector module account from isolated market account": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(3),
				MarketId:                          uint32(3),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.IntToString(3),
			),
			affiliateRevShareAcctAddr:               "",
			marketMapperRevShareAcctAddr:            constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr:           "",
			expectedSubaccountsModuleAccBalance:     big.NewInt(100),  // 600 - 500
			expectedFeeModuleAccBalance:             big.NewInt(3000), // 500 + 2500
			marketMapperAccBalance:                  big.NewInt(0),
			revShare:                                revsharetypes.RevSharesForFill{},
			expectedMarketMapperAccBalance:          big.NewInt(0),
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
		},
		"failure - subaccounts module does not have sufficient funds": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(300),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(3),
				MarketId:                          uint32(3),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			collateralPoolAddr:                      types.ModuleAddress,
			affiliateRevShareAcctAddr:               "",
			marketMapperRevShareAcctAddr:            constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr:           "",
			expectedSubaccountsModuleAccBalance:     big.NewInt(300),
			expectedFeeModuleAccBalance:             big.NewInt(2500),
			expectedErr:                             sdkerrors.ErrInsufficientFunds,
			marketMapperAccBalance:                  big.NewInt(0),
			revShare:                                revsharetypes.RevSharesForFill{},
			expectedMarketMapperAccBalance:          big.NewInt(0),
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
		},
		"failure - isolated markets subaccounts module does not have sufficient funds": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(300),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(3),
				MarketId:                          uint32(3),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.IntToString(3),
			),
			affiliateRevShareAcctAddr:               "",
			marketMapperRevShareAcctAddr:            constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr:           "",
			expectedSubaccountsModuleAccBalance:     big.NewInt(300),
			expectedFeeModuleAccBalance:             big.NewInt(2500),
			expectedErr:                             sdkerrors.ErrInsufficientFunds,
			marketMapperAccBalance:                  big.NewInt(0),
			revShare:                                revsharetypes.RevSharesForFill{},
			expectedMarketMapperAccBalance:          big.NewInt(0),
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
		},
		"failure - asset ID doesn't exist": {
			feeModuleAccBalance:        big.NewInt(1500),
			skipSetUpUsdc:              true,
			asset:                      *constants.Usdc,
			subaccountModuleAccBalance: big.NewInt(500),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(3),
				MarketId:                          uint32(3),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			collateralPoolAddr:                      types.ModuleAddress,
			affiliateRevShareAcctAddr:               "",
			marketMapperRevShareAcctAddr:            constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr:           "",
			expectedErr:                             asstypes.ErrAssetDoesNotExist,
			expectedSubaccountsModuleAccBalance:     big.NewInt(500),
			expectedFeeModuleAccBalance:             big.NewInt(1500),
			marketMapperAccBalance:                  big.NewInt(0),
			revShare:                                revsharetypes.RevSharesForFill{},
			expectedMarketMapperAccBalance:          big.NewInt(0),
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
		},
		"failure - asset other than USDC not supported": {
			feeModuleAccBalance:        big.NewInt(1500),
			asset:                      *constants.BtcUsd,
			subaccountModuleAccBalance: big.NewInt(500),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(3),
				MarketId:                          uint32(3),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			collateralPoolAddr:                      types.ModuleAddress,
			affiliateRevShareAcctAddr:               "",
			marketMapperRevShareAcctAddr:            constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr:           "",
			expectedErr:                             types.ErrAssetTransferThroughBankNotImplemented,
			expectedSubaccountsModuleAccBalance:     big.NewInt(500),
			expectedFeeModuleAccBalance:             big.NewInt(1500),
			marketMapperAccBalance:                  big.NewInt(0),
			revShare:                                revsharetypes.RevSharesForFill{},
			expectedMarketMapperAccBalance:          big.NewInt(0),
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
		},
		"success - distribute fees to market mapper and fee collector": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			marketMapperAccBalance:     big.NewInt(0),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(4),
				MarketId:                          uint32(4),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			expectedSubaccountsModuleAccBalance:     big.NewInt(100),  // 600 - 500
			expectedFeeModuleAccBalance:             big.NewInt(2950), // 2500 + 500 - 50
			expectedMarketMapperAccBalance:          big.NewInt(50),   // 0 + 50
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.IntToString(4),
			),
			affiliateRevShareAcctAddr:     "",
			marketMapperRevShareAcctAddr:  constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr: "",
			revShare: revsharetypes.RevSharesForFill{
				AffiliateRevShare: nil,
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(0),
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(50),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            0,       // 0%
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
				},
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(50),
						RevSharePpm:       100_000, // 10%
					},
				},
			},
		},
		"success - distribute fees to market mapper, fee collector, and order router rev share": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			marketMapperAccBalance:     big.NewInt(0),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(4),
				MarketId:                          uint32(4),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
				MakerOrderRouterAddr:              constants.BobAccAddress.String(),
				TakerOrderRouterAddr:              constants.AliceAccAddress.String(),
			},
			expectedSubaccountsModuleAccBalance:     big.NewInt(100),  // 600 - 500
			expectedFeeModuleAccBalance:             big.NewInt(2850), // 2500 + 500 - 150
			expectedMarketMapperAccBalance:          big.NewInt(50),   // 0 + 50
			expectedBuySideOrderRouterRevShare:      big.NewInt(50),   // 0 + 50
			expectedSellSideOrderRouterRevShare:     big.NewInt(50),   // 0 + 50
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.IntToString(4),
			),
			affiliateRevShareAcctAddr:           "",
			buySideOrderRouterRevShareAcctAddr:  constants.BobAccAddress.String(),
			sellSideOrderRouterRevShareAcctAddr: constants.CarlAccAddress.String(),
			marketMapperRevShareAcctAddr:        constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr:       "",
			revShare: revsharetypes.RevSharesForFill{
				AffiliateRevShare: nil,
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(50),
					revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE:            big.NewInt(50),
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(50),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            100_000, // 10%
					revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE:            100_000, // 10%
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 100_000, // 10%
				},
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(50),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER,
						QuoteQuantums:     big.NewInt(50),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         constants.CarlAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER,
						QuoteQuantums:     big.NewInt(50),
						RevSharePpm:       100_000, // 10%
					},
				},
			},
		},
		"success - market mapper rev share rounded down to 0": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(100),
			subaccountModuleAccBalance: big.NewInt(200),
			marketMapperAccBalance:     big.NewInt(0),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(5),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(4),
				FillQuoteQuantums:                 big.NewInt(9),
				ProductId:                         uint32(4),
				MarketId:                          uint32(4),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			expectedSubaccountsModuleAccBalance:     big.NewInt(191), // 200 - 9
			expectedFeeModuleAccBalance:             big.NewInt(109), // 100 + 9
			expectedMarketMapperAccBalance:          big.NewInt(0),
			expectedAffiliateAccBalance:             big.NewInt(0),
			expectedUnconditionalRevShareAccBalance: big.NewInt(0),
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.IntToString(4),
			),
			affiliateRevShareAcctAddr:     "",
			marketMapperRevShareAcctAddr:  constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr: "",
			revShare: revsharetypes.RevSharesForFill{
				AffiliateRevShare: nil,
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(0),
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(0),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            0, // 0%
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 0, // 0%
				},
				AllRevShares: []revsharetypes.RevShare{},
			},
		},
		"success - distribute fees to market mapper, unconditional rev share, affiliate and fee collector": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			marketMapperAccBalance:     big.NewInt(0),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(4),
				MarketId:                          uint32(4),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
			},
			expectedSubaccountsModuleAccBalance:     big.NewInt(100),  // 600 - 500
			expectedFeeModuleAccBalance:             big.NewInt(2892), // 2500 + 500 - 108
			expectedMarketMapperAccBalance:          big.NewInt(48),   // 10% of 488
			expectedAffiliateAccBalance:             big.NewInt(12),   // 5%  of 250
			expectedUnconditionalRevShareAccBalance: big.NewInt(48),   // 10%  of 488
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.IntToString(4),
			),
			affiliateRevShareAcctAddr:     constants.BobAccAddress.String(),
			marketMapperRevShareAcctAddr:  constants.AliceAccAddress.String(),
			unconditionalRevShareAcctAddr: constants.CarlAccAddress.String(),
			revShare: revsharetypes.RevSharesForFill{
				AffiliateRevShare: &revsharetypes.RevShare{
					Recipient:         constants.BobAccAddress.String(),
					RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
					RevShareType:      revsharetypes.REV_SHARE_TYPE_AFFILIATE,
					QuoteQuantums:     big.NewInt(12),
					RevSharePpm:       50_000, // 5%
				},
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(12),
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(96),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            50_000,  // 5%
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 200_000, // 20%
				},
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_AFFILIATE,
						QuoteQuantums:     big.NewInt(12),
						RevSharePpm:       50_000, // 5%
					},
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(48),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         constants.CarlAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(48),
						RevSharePpm:       100_000, // 10%
					},
				},
			},
		},
		"success - distribute fees to market mapper, unconditional rev share, order router, and fee collector": {
			asset:                      *constants.Usdc,
			feeModuleAccBalance:        big.NewInt(2500),
			subaccountModuleAccBalance: big.NewInt(600),
			marketMapperAccBalance:     big.NewInt(0),
			fill: clobtypes.FillForProcess{
				TakerAddr:                         constants.AliceAccAddress.String(),
				TakerFeeQuoteQuantums:             big.NewInt(250),
				MakerAddr:                         constants.BobAccAddress.String(),
				MakerFeeQuoteQuantums:             big.NewInt(250),
				FillQuoteQuantums:                 big.NewInt(500),
				ProductId:                         uint32(4),
				MarketId:                          uint32(4),
				MonthlyRollingTakerVolumeQuantums: 1_000_000,
				TakerOrderRouterAddr:              constants.BobAccAddress.String(),
				MakerOrderRouterAddr:              constants.DaveAccAddress.String(),
			},
			expectedSubaccountsModuleAccBalance:     big.NewInt(100),  // 600 - 500
			expectedFeeModuleAccBalance:             big.NewInt(2873), // 2500 + 500 - 127
			expectedMarketMapperAccBalance:          big.NewInt(46),   // 10% of 465
			expectedBuySideOrderRouterRevShare:      big.NewInt(20),   // 8% of 250
			expectedSellSideOrderRouterRevShare:     big.NewInt(15),   // 6% of 250
			expectedUnconditionalRevShareAccBalance: big.NewInt(46),   // 10% of 465
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.IntToString(4),
			),
			marketMapperRevShareAcctAddr:        constants.AliceAccAddress.String(),
			buySideOrderRouterRevShareAcctAddr:  constants.BobAccAddress.String(),
			unconditionalRevShareAcctAddr:       constants.CarlAccAddress.String(),
			sellSideOrderRouterRevShareAcctAddr: constants.DaveAccAddress.String(),
			revShare: revsharetypes.RevSharesForFill{
				FeeSourceToQuoteQuantums: map[revsharetypes.RevShareFeeSource]*big.Int{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            big.NewInt(20),
					revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE:            big.NewInt(15),
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: big.NewInt(92),
				},
				FeeSourceToRevSharePpm: map[revsharetypes.RevShareFeeSource]uint32{
					revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE:            80_000,  // 8%
					revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE:            60_000,  // 6%
					revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE: 200_000, // 20%
				},
				AllRevShares: []revsharetypes.RevShare{
					{
						Recipient:         constants.AliceAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_MARKET_MAPPER,
						QuoteQuantums:     big.NewInt(46),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         constants.CarlAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_NET_PROTOCOL_REVENUE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_UNCONDITIONAL,
						QuoteQuantums:     big.NewInt(46),
						RevSharePpm:       100_000, // 10%
					},
					{
						Recipient:         constants.BobAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_TAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER,
						QuoteQuantums:     big.NewInt(20),
						RevSharePpm:       80_000, // 8%
					},
					{
						Recipient:         constants.DaveAccAddress.String(),
						RevShareFeeSource: revsharetypes.REV_SHARE_FEE_SOURCE_MAKER_FEE,
						RevShareType:      revsharetypes.REV_SHARE_TYPE_ORDER_ROUTER,
						QuoteQuantums:     big.NewInt(15),
						RevSharePpm:       60_000, // 6%
					},
				},
			},
		},
		// TODO(DEC-715): Add more test for non-USDC assets, after asset update
		// is implemented.
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper,
				bankKeeper, assetsKeeper, _, _, _, _ :=
				keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			keepertest.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)
			keepertest.CreateTestPerpetuals(t, ctx, perpetualsKeeper)
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
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.feeModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.subaccountModuleAccBalance.Sign() > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					tc.collateralPoolAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.subaccountModuleAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			marketMapperAddr, err := sdk.AccAddressFromBech32(tc.marketMapperRevShareAcctAddr)
			require.NoError(t, err)

			if tc.marketMapperAccBalance.Sign() > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					marketMapperAddr,
					sdk.Coins{
						sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.marketMapperAccBalance)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			// Always create USDC as the first asset.
			if !tc.skipSetUpUsdc {
				err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
				require.NoError(t, err)
			}

			if tc.asset.Denom != constants.Usdc.Denom {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					tc.asset.Id,
					tc.asset.Symbol,
					tc.asset.Denom,
					tc.asset.DenomExponent,
					tc.asset.HasMarket,
					tc.asset.MarketId,
					tc.asset.AtomicResolution,
				)
				require.NoError(t, err)
			}

			err = keeper.DistributeFees(ctx, tc.asset.Id, tc.revShare, tc.fill)

			if tc.expectedErr != nil {
				require.ErrorIs(t,
					err,
					tc.expectedErr,
				)
			} else {
				require.NoError(t, err)
			}

			// Check the subaccount module balance.
			subaccountsModuleAccBalance := bankKeeper.GetBalance(ctx, tc.collateralPoolAddr, tc.asset.Denom)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedSubaccountsModuleAccBalance)),
				subaccountsModuleAccBalance,
			)

			// Check the fee module account balance has been updated as expected.
			toModuleBalance := bankKeeper.GetBalance(
				ctx, authtypes.NewModuleAddress(authtypes.FeeCollectorName),
				tc.asset.Denom,
			)
			require.Equal(t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedFeeModuleAccBalance)),
				toModuleBalance,
			)

			// Check the market mapper account balance has been updated as expected.
			marketMapperBalance := bankKeeper.GetBalance(
				ctx, marketMapperAddr,
				tc.asset.Denom,
			)
			require.Equal(
				t,
				sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedMarketMapperAccBalance)),
				marketMapperBalance,
			)

			if tc.expectedBuySideOrderRouterRevShare != nil && tc.expectedBuySideOrderRouterRevShare.Sign() > 0 {
				// Check the buy side order router rev share account balance has been updated as expected.
				buySideOrderRouterRevShareAddr, err := sdk.AccAddressFromBech32(tc.buySideOrderRouterRevShareAcctAddr)
				require.NoError(t, err)
				buySideOrderRouterRevShareBalance := bankKeeper.GetBalance(
					ctx, buySideOrderRouterRevShareAddr,
					tc.asset.Denom,
				)
				require.Equal(
					t,
					sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedBuySideOrderRouterRevShare)),
					buySideOrderRouterRevShareBalance,
				)
			}

			if tc.expectedSellSideOrderRouterRevShare != nil && tc.expectedSellSideOrderRouterRevShare.Sign() > 0 {
				// Check the sell side order router rev share account balance has been updated as expected.
				sellSideOrderRouterRevShareAddr, err := sdk.AccAddressFromBech32(tc.sellSideOrderRouterRevShareAcctAddr)
				require.NoError(t, err)
				sellSideOrderRouterRevShareBalance := bankKeeper.GetBalance(
					ctx, sellSideOrderRouterRevShareAddr,
					tc.asset.Denom,
				)
				require.Equal(
					t,
					sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedSellSideOrderRouterRevShare)),
					sellSideOrderRouterRevShareBalance,
				)
			}

			// Check the unconditional rev share account balance has been updated as expected.
			if tc.expectedUnconditionalRevShareAccBalance.Sign() > 0 {
				unconditionalRevShareAddr, err := sdk.AccAddressFromBech32(tc.unconditionalRevShareAcctAddr)
				require.NoError(t, err)
				unconditionalRevShareBalance := bankKeeper.GetBalance(
					ctx, unconditionalRevShareAddr,
					tc.asset.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedUnconditionalRevShareAccBalance)),
					unconditionalRevShareBalance,
				)
			}

			// Check the affiliate account balance has been updated as expected.
			if tc.expectedAffiliateAccBalance != nil && tc.expectedAffiliateAccBalance.Sign() > 0 {
				affiliateAddr, err := sdk.AccAddressFromBech32(tc.affiliateRevShareAcctAddr)
				require.NoError(t, err)
				affiliateBalance := bankKeeper.GetBalance(
					ctx, affiliateAddr,
					tc.asset.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(tc.asset.Denom, sdkmath.NewIntFromBigInt(tc.expectedAffiliateAccBalance)),
					affiliateBalance,
				)
			}
		})
	}
}

func TestTransferInsuranceFundPayments(t *testing.T) {
	tests := map[string]struct {
		skipSetUpUsdc bool

		// Module account state.
		subaccountModuleAccBalance int64
		insuranceFundBalance       int64
		perpetual                  perptypes.Perpetual
		collateralPoolAddr         sdk.AccAddress

		// Transfer details.
		quantums *big.Int

		// Expectations.
		panics                              bool
		expectedErr                         error
		expectedSubaccountsModuleAccBalance int64
		expectedInsuranceFundBalance        int64
	}{
		"success - send to insurance fund module account": {
			perpetual:                           constants.BtcUsd_SmallMarginRequirement,
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          600,
			quantums:                            big.NewInt(500),
			collateralPoolAddr:                  types.ModuleAddress,
			expectedSubaccountsModuleAccBalance: 100,  // 600 - 500
			expectedInsuranceFundBalance:        3000, // 2500 + 500
		},
		"success - send from insurance fund module account": {
			perpetual:                           constants.BtcUsd_SmallMarginRequirement,
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          600,
			quantums:                            big.NewInt(-500),
			collateralPoolAddr:                  types.ModuleAddress,
			expectedSubaccountsModuleAccBalance: 1100, // 600 + 500
			expectedInsuranceFundBalance:        2000, // 2500 - 500
		},
		"success - can send zero payment": {
			perpetual:                           constants.BtcUsd_SmallMarginRequirement,
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          600,
			quantums:                            big.NewInt(0),
			collateralPoolAddr:                  types.ModuleAddress,
			expectedSubaccountsModuleAccBalance: 600,
			expectedInsuranceFundBalance:        2500,
		},
		"success - send to isolated insurance fund account": {
			perpetual:                  constants.IsoUsd_IsolatedMarket,
			insuranceFundBalance:       2500,
			subaccountModuleAccBalance: 600,
			quantums:                   big.NewInt(500),
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.GetId()),
			),
			expectedSubaccountsModuleAccBalance: 100,  // 600 - 500
			expectedInsuranceFundBalance:        3000, // 2500 + 500
		},
		"success - send from isolated insurance fund account": {
			perpetual:                  constants.IsoUsd_IsolatedMarket,
			insuranceFundBalance:       2500,
			subaccountModuleAccBalance: 600,
			quantums:                   big.NewInt(-500),
			collateralPoolAddr: authtypes.NewModuleAddress(
				types.ModuleName + ":" + lib.UintToString(constants.IsoUsd_IsolatedMarket.GetId()),
			),
			expectedSubaccountsModuleAccBalance: 1100, // 600 + 500
			expectedInsuranceFundBalance:        2000, // 2500 - 500
		},
		"failure - subaccounts module does not have sufficient funds": {
			perpetual:                           constants.BtcUsd_SmallMarginRequirement,
			insuranceFundBalance:                2500,
			subaccountModuleAccBalance:          300,
			quantums:                            big.NewInt(500),
			collateralPoolAddr:                  types.ModuleAddress,
			expectedSubaccountsModuleAccBalance: 300,
			expectedInsuranceFundBalance:        2500,
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
		},
		"failure - insurance fund does not have sufficient funds": {
			perpetual:                           constants.BtcUsd_SmallMarginRequirement,
			insuranceFundBalance:                300,
			subaccountModuleAccBalance:          2500,
			quantums:                            big.NewInt(-500),
			collateralPoolAddr:                  types.ModuleAddress,
			expectedSubaccountsModuleAccBalance: 2500,
			expectedInsuranceFundBalance:        300,
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
		},
		"failure - isolated market insurance fund does not have sufficient funds": {
			perpetual:                           constants.IsoUsd_IsolatedMarket,
			insuranceFundBalance:                300,
			subaccountModuleAccBalance:          2500,
			quantums:                            big.NewInt(-500),
			expectedSubaccountsModuleAccBalance: 2500,
			expectedInsuranceFundBalance:        300,
			expectedErr:                         sdkerrors.ErrInsufficientFunds,
		},
		"panics - asset doesn't exist": {
			perpetual:                           constants.BtcUsd_SmallMarginRequirement,
			insuranceFundBalance:                1500,
			skipSetUpUsdc:                       true,
			subaccountModuleAccBalance:          500,
			quantums:                            big.NewInt(500),
			collateralPoolAddr:                  types.ModuleAddress,
			expectedErr:                         errorsmod.Wrap(asstypes.ErrAssetDoesNotExist, lib.UintToString(uint32(0))),
			expectedSubaccountsModuleAccBalance: 500,
			expectedInsuranceFundBalance:        1500,
			panics:                              true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpsKeeper, accountKeeper, bankKeeper, assetsKeeper, _, _, _, _ :=
				keepertest.SubaccountsKeepers(t, true)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpsKeeper)

			// Set up Subaccounts module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, types.ModuleName, []string{})
			// Set up insurance fund module account.
			auth_testutil.CreateTestModuleAccount(ctx, accountKeeper, perptypes.InsuranceFundName, []string{})

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

			_, err := perpsKeeper.CreatePerpetual(
				ctx,
				tc.perpetual.GetId(),
				tc.perpetual.Params.Ticker,
				tc.perpetual.Params.MarketId,
				tc.perpetual.Params.AtomicResolution,
				tc.perpetual.Params.DefaultFundingPpm,
				tc.perpetual.Params.LiquidityTier,
				tc.perpetual.Params.MarketType,
			)
			require.NoError(t, err)

			insuranceFundAddr, err := perpsKeeper.GetInsuranceFundModuleAddress(ctx, tc.perpetual.GetId())
			require.NoError(t, err)

			// Mint asset in the receipt/sender module account for transfer.
			if tc.insuranceFundBalance > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					insuranceFundAddr,
					sdk.Coins{
						sdk.NewInt64Coin(constants.Usdc.Denom, tc.insuranceFundBalance),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.subaccountModuleAccBalance > 0 {
				err := bank_testutil.FundAccount(
					ctx,
					tc.collateralPoolAddr,
					sdk.Coins{
						sdk.NewInt64Coin(constants.Usdc.Denom, tc.subaccountModuleAccBalance),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if !tc.skipSetUpUsdc {
				err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
				require.NoError(t, err)
			}

			if tc.expectedErr != nil {
				if tc.panics {
					require.PanicsWithError(
						t,
						tc.expectedErr.Error(),
						func() {
							//nolint:errcheck
							keeper.TransferInsuranceFundPayments(ctx, tc.quantums, tc.perpetual.GetId())
						},
					)
				} else {
					require.ErrorIs(
						t,
						keeper.TransferInsuranceFundPayments(ctx, tc.quantums, tc.perpetual.GetId()),
						tc.expectedErr,
					)
				}
			} else {
				require.NoError(t, keeper.TransferInsuranceFundPayments(ctx, tc.quantums, tc.perpetual.GetId()))
			}

			// Check the subaccount module balance.
			subaccountsModuleAccBalance := bankKeeper.GetBalance(ctx, tc.collateralPoolAddr, constants.Usdc.Denom)
			require.Equal(
				t,
				sdk.NewInt64Coin(constants.Usdc.Denom, tc.expectedSubaccountsModuleAccBalance),
				subaccountsModuleAccBalance,
			)

			// Check the fee module account balance has been updated as expected.
			toModuleBalance := bankKeeper.GetBalance(
				ctx, insuranceFundAddr,
				constants.Usdc.Denom,
			)
			require.Equal(t,
				sdk.NewInt64Coin(constants.Usdc.Denom, tc.expectedInsuranceFundBalance),
				toModuleBalance,
			)
		})
	}
}
