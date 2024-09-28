package keeper_test

import (
	"math/big"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	bank_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/bank"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	asstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	subaccountskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

func TestDepositYieldToSubaccount(t *testing.T) {
	// default subaccount id, the first subaccount id generated when calling createNSubaccount
	defaultSubaccountId := types.SubaccountId{
		Owner:  "0",
		Number: 0,
	}

	tests := map[string]struct {
		// state
		perpetuals                 []perptypes.Perpetual
		fundsInTDaiPool            *big.Int
		collateralPoolTDaiBalances map[string]*big.Int

		// subaccount state (include perp positions for isolated markets)
		perpetualPositions []*types.PerpetualPosition

		// test data
		subaccountId     types.SubaccountId
		amountToTransfer *big.Int

		// expectations
		expectedCollateralPoolTDaiBalances map[string]*big.Int
		expectedTDaiYieldPoolBalance       *big.Int
		expectedErr                        error
	}{
		"Success: deposits basic yield to account": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: big.NewInt(10_000_000_000),

			expectedTDaiYieldPoolBalance: big.NewInt(90_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(110_000_000_000),
			},
		},
		"Success: deposits all yield present tDai pool": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(1),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: big.NewInt(100_000_000_000),

			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_001),
			},
		},
		"Success: deposits tDai yield into isolated collateral pool": {
			fundsInTDaiPool: big.NewInt(2),
			collateralPoolTDaiBalances: map[string]*big.Int{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): big.NewInt(100_000_000_000),
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},

			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: big.NewInt(1),

			expectedTDaiYieldPoolBalance: big.NewInt(1),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): big.NewInt(100_000_000_001),
				types.ModuleAddress.String(): big.NewInt(0),
			},
		},
		"Success: deposits tDai amount greater than max uint64": {
			fundsInTDaiPool: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("100000000000000000000000000"),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(1),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("100000000000000000000000000"),

			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("100000000000000000000000001"),
			},
		},
		"Success: depositing nil transfer amount results in no change": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: nil,

			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},
		},
		"Success: depositing zero transfer amount results in no change": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: big.NewInt(0),

			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},
		},
		"Success: deposits when subaccount is not previously set": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},

			subaccountId: types.SubaccountId{
				Owner:  "1",
				Number: 1,
			},
			amountToTransfer: big.NewInt(10_000_000_000),

			expectedTDaiYieldPoolBalance: big.NewInt(90_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(110_000_000_000),
			},
		},
		"Failure: not enough yield in pool to deposit": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: big.NewInt(110_000_000_000),

			expectedErr: errorsmod.Wrapf(
				sdkerrors.ErrInsufficientFunds,
				"spendable balance 100000000000utdai is smaller than 110000000000utdai",
			),

			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},
		},
		"Failure: attempts to deposit negative yield": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: big.NewInt(-100),

			expectedErr: types.ErrTryingToDepositNegativeYield,

			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},
		},
		"Failure: collateral pool not found": {
			fundsInTDaiPool: big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},

			subaccountId:     defaultSubaccountId,
			amountToTransfer: big.NewInt(100),

			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},

			expectedErr: errorsmod.Wrap(perptypes.ErrPerpetualDoesNotExist, lib.UintToString(uint64(3))),

			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]*big.Int{
				types.ModuleAddress.String(): big.NewInt(100_000_000_000),
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(
				t,
				true,
			)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)
			rateLimitKeeper.SetSDAIPrice(ctx, rate)
			rateLimitKeeper.SetAssetYieldIndex(ctx, big.NewRat(1, 1))

			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))

			for _, p := range tc.perpetuals {
				perpetualsKeeper.SetPerpetualForTest(
					ctx,
					p,
				)
			}

			for collateralPoolAddr, TDaiBal := range tc.collateralPoolTDaiBalances {
				err := bank_testutil.FundAccount(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(TDaiBal)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.fundsInTDaiPool != nil {
				err := bank_testutil.FundModuleAccount(
					ctx,
					ratelimittypes.TDaiPoolAccount,
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.fundsInTDaiPool)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.PerpetualPositions = tc.perpetualPositions
			keeper.SetSubaccount(ctx, subaccount)

			err := keeper.DepositYieldToSubaccount(ctx, tc.subaccountId, tc.amountToTransfer)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}

			for collateralPoolAddr, expectedTDaiBal := range tc.expectedCollateralPoolTDaiBalances {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(expectedTDaiBal)),
					TDaiBal,
				)
			}

			if tc.expectedTDaiYieldPoolBalance != nil {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(ratelimittypes.TDaiPoolAccount),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.expectedTDaiYieldPoolBalance)),
					TDaiBal,
				)
			}
		})
	}
}

func TestAddYieldToSubaccount(t *testing.T) {
	tests := map[string]struct {
		// state
		perpetuals []perptypes.Perpetual
		assets     []*asstypes.Asset
		// Only set when specified. Defaults to 0/1.
		globalAssetYieldIndex *big.Rat
		availableYield        *big.Int

		// subaccount state
		perpetualPositions        []*types.PerpetualPosition
		assetPositions            []*types.AssetPosition
		subaccountAssetYieldIndex string

		// additional test state
		PerpIdToPerp map[uint32]perptypes.Perpetual

		// expectations
		expectedTotalYield         *big.Int
		expectedPerpetualPositions []*types.PerpetualPosition
		expectedAssetPositions     []*types.AssetPosition
		expectedErr                error

		// Only contains the updated perpetual positions, to assert against the events included.
		expectedAssetYieldIndex string
	}{
		"Successfully adds yield when only tDai asset claims yield and no perpetual position exists": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			expectedErr:               nil,
			expectedTotalYield:        big.NewInt(100_000_000_000),
			expectedAssetYieldIndex:   big.NewRat(2, 1).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
		},
		"Successfully adds yield when only tDai asset claims yield, but perpetual position exists": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: constants.BtcUsd_NoMarginRequirement,
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(100_000_000_000),
			expectedAssetYieldIndex: big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
		},
		"Successfully adds yield when only perp position claims yield and no asset position exists": {
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(100_000_000),
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000),
				},
			},
		},
		"Successfully adds yield when only perp position claims yield, but asset position exists": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(100_000_000),
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_100_000_000),
				},
			},
		},
		"Successfully adds yield when asset position and perp position claim yield": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(2, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(100_100_000_000),
			expectedAssetYieldIndex: big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_100_000_000),
				},
			},
		},
		"Successfully claims yield when multiple perp position claim yield and no asset position exists": {
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(12, 11),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				1: {
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(94444443),
			expectedAssetYieldIndex: big.NewRat(12, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(94444443),
				},
			},
		},
		"Successfully claims yield when multiple perp position claim yield and asset position exists but doesn't claim yield": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(12, 11).String(),
			globalAssetYieldIndex:     big.NewRat(12, 11),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				1: {
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(94444443),
			expectedAssetYieldIndex: big.NewRat(12, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_094_444_443),
				},
			},
		},
		"Successfully claims yield when multiple perp position claim yield and asset position exists and claims yield": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(12, 11).String(),
			globalAssetYieldIndex:     big.NewRat(12, 10),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				1: {
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(10_094_444_443),
			expectedAssetYieldIndex: big.NewRat(12, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(110_094_444_443),
				},
			},
		},
		"Successfully claims yield when tDai position and perp postions claim yield and one perp position is non-standard": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(12, 11).String(),
			globalAssetYieldIndex:     big.NewRat(12, 10),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
				{
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
				{
					PerpetualId:  uint32(constants.PerpetualPosition_OneISOLong.PerpetualId),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				1: {
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
				constants.PerpetualPosition_OneISOLong.PerpetualId: {
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(10_014_444_443),
			expectedAssetYieldIndex: big.NewRat(12, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
				{
					PerpetualId:  uint32(constants.PerpetualPosition_OneISOLong.PerpetualId),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(110_014_444_443),
				},
			},
		},
		"Successfully claims yield with negative and positive position": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(12, 11).String(),
			globalAssetYieldIndex:     big.NewRat(12, 10),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(constants.PerpetualPosition_OneISOLong.PerpetualId),
					Quantums:     dtypes.NewInt(-1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				constants.PerpetualPosition_OneISOLong.PerpetualId: {
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(9_000_000_000),
			expectedAssetYieldIndex: big.NewRat(12, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(constants.PerpetualPosition_OneISOLong.PerpetualId),
					Quantums:     dtypes.NewInt(-1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(109_000_000_000),
				},
			},
		},
		"Successfully adds 0 yield when yield negative: only tDai position open": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(-100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			expectedErr:               nil,
			expectedTotalYield:        big.NewInt(0),
			expectedAssetYieldIndex:   big.NewRat(1, 1).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(-100_000_000_000),
				},
			},
		},
		"Successfully adds 0 yield when negative yield: asset position and perp position claim yield and negative position with more yield": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(2, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(2, 1).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(0),
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(2, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
		},
		"Successfully adds 0 yield: asset position and perp position cancel each other out": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(0),
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
		},
		"Successfully adds 0 yield: general yield index is 0": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},

			expectedErr:             nil,
			expectedTotalYield:      big.NewInt(0),
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
		},
		"Successfully handled insufficient yield available when only tDai asset claims yield and no perpetual position exists": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			availableYield:            big.NewInt(1_000_000),
			expectedErr:               nil,
			expectedTotalYield:        big.NewInt(1_000_000),
			expectedAssetYieldIndex:   big.NewRat(2, 1).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_001_000_000),
				},
			},
		},
		"Failure: subaccount's asset yield index is empty": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: "",
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},

			expectedErr: types.ErrYieldIndexUninitialized,
		},
		"Failure: global asset yield index is negative": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(-1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			expectedErr:               types.ErrGlobalYieldIndexNegative,
		},
		"Failure: global asset yield less than asset yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(2, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			expectedErr:               types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
		},
		"Failure: other asset than tDai": {
			assetPositions: []*types.AssetPosition{
				{
					AssetId:  1,
					Quantums: dtypes.NewInt(100_000_000_000),
					Index:    0,
				},
			},
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			expectedErr:               assettypes.ErrNotImplementedMulticollateral,
		},
		"Failure: Tries to add yield when subaccount has perp position that is not in perpIdToPerp map": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				1: constants.BtcUsd_NoMarginRequirement,
			},

			expectedErr: errorsmod.Wrap(
				perptypes.ErrPerpetualDoesNotExist, "0",
			),
		},
		"Failure: Tries to add yield when perp general yield index is nil": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   "",
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   "",
				},
			},

			expectedErr: types.ErrYieldIndexUninitialized,
		},
		"Failure: Tries to add yield when perp general yield index is negative": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(-1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(-1, 1).String(),
				},
			},

			expectedErr: types.ErrGlobalYieldIndexNegative,
		},
		"Failure: Tries to add yield when perp general yield index is less than perp position yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			availableYield:            big.NewInt(1_000_000_000_000_000_000),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(2, 1).String(),
				},
			},
			PerpIdToPerp: map[uint32]perptypes.Perpetual{
				0: {
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},

			expectedErr: types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, _, _, assetsKeeper, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(
				t,
				true,
			)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			rateLimitKeeper.SetSDAIPrice(ctx, rate)
			globalAssetYieldIndex := big.NewRat(1, 1)
			if tc.globalAssetYieldIndex != nil {
				globalAssetYieldIndex = tc.globalAssetYieldIndex
			}
			rateLimitKeeper.SetAssetYieldIndex(ctx, globalAssetYieldIndex)

			availableYield := big.NewInt(0)
			if tc.availableYield != nil {
				availableYield = tc.availableYield
			}

			// Always creates TDai asset first
			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Id,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
					a.AssetYieldIndex,
				)
				require.NoError(t, err)
			}

			for _, p := range tc.perpetuals {
				perpetualsKeeper.SetPerpetualForTest(
					ctx,
					p,
				)
			}

			subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.PerpetualPositions = tc.perpetualPositions
			subaccount.AssetPositions = tc.assetPositions
			subaccount.AssetYieldIndex = tc.subaccountAssetYieldIndex
			keeper.SetSubaccount(ctx, subaccount)

			newSubaccount, totalNewYield, err := subaccountskeeper.AddYieldToSubaccount(subaccount, tc.PerpIdToPerp, globalAssetYieldIndex, availableYield)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Nil(t, totalNewYield)
				require.Empty(t, newSubaccount)
			} else {
				require.NoError(t, err)
				require.Equal(t, 0, tc.expectedTotalYield.Cmp(totalNewYield),
					"Expected Yield: %v. Got: %v", tc.expectedTotalYield, totalNewYield)
				require.Equal(t, len(newSubaccount.PerpetualPositions), len(tc.expectedPerpetualPositions))
				for i, ep := range tc.expectedPerpetualPositions {
					require.Equal(t, *ep, *newSubaccount.PerpetualPositions[i])
				}
				require.Equal(t, len(newSubaccount.AssetPositions), len(tc.expectedAssetPositions))
				for i, ep := range tc.expectedAssetPositions {
					require.Equal(t, *ep, *newSubaccount.AssetPositions[i])
				}
				require.Equal(t, 0, globalAssetYieldIndex.Cmp(ratelimitkeeper.ConvertStringToBigRatWithPanicOnErr(newSubaccount.AssetYieldIndex)),
					"Expected AssetYieldIndex %v. Got %v.", globalAssetYieldIndex, newSubaccount.AssetYieldIndex,
				)
			}
		})
	}
}

func TestClaimYieldForSubaccountFromIdAndSetNewState(t *testing.T) {
	// default subaccount id, the first subaccount id generated when calling createNSubaccount
	defaultSubaccountId := types.SubaccountId{
		Owner:  "0",
		Number: 0,
	}

	tests := map[string]struct {
		// state
		perpetuals []perptypes.Perpetual
		assets     []*asstypes.Asset
		// Only set when specified. Defaults to 0/1.
		// Set perpYieldIndex in the perpetuals state.
		globalAssetYieldIndex *big.Rat
		fundsInTDaiPool       *big.Int

		// subaccount state
		perpetualPositions        []*types.PerpetualPosition
		assetPositions            []*types.AssetPosition
		subaccountAssetYieldIndex string

		// collateral pool state
		collateralPoolTDaiBalances map[string]int64

		// extra test state
		subaccountId *types.SubaccountId

		// expectations
		expectedCollateralPoolTDaiBalances map[string]int64
		expectedPerpetualPositions         []*types.PerpetualPosition
		expectedAssetPositions             []*types.AssetPosition
		expectedTDaiYieldPoolBalance       *big.Int
		expectedErr                        error
		expectedAssetYieldIndex            string
	}{
		"Successfully claims no yield for tDai asset position and no other position exists when global asset yield index is 1": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Successfully claims yield for tDai asset position and no other position exists when subaccount's asset yield index 0": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
		},
		"Successfully claims yield for tDai asset position and no other position exists when subaccount's asset yield index 1": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
		},
		"Successfully claims yield for tDai asset position when perp with no yield exists": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
		},
		"Successfully claims yield for one perp with no asset positions existing before yield claim": {
			globalAssetYieldIndex: big.NewRat(1, 1),
			fundsInTDaiPool:       big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(199_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 101_000_000_000,
			},
		},
		"Successfully claims yield for one perp with asset position existing but not claiming yield": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(101_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(199_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 101_000_000_000,
			},
		},
		"Successfully claims yield for tDai asset and one perp": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_100_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(99_900_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_100_000_000,
			},
		},
		"Successfully claims yield when multiple perp positions are open and tDai position open": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(50_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(13, 11).String(),
			globalAssetYieldIndex:     big.NewRat(26, 11),
			fundsInTDaiPool:           big.NewInt(222_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 50_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(26, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_094_444_443), // Total Yield: 50_094_444_443
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(171_905_555_557),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_094_444_443,
			},
		},
		"Successfully claims all yield in tDaiPool": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(3, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(3, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(300_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 300_000_000_000,
			},
		},
		"Successfully claims yield for isolated market": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(5, 4),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(5, 4).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(124_920_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(175_080_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 124_920_000_000,
				types.ModuleAddress.String(): 0,
			},
		},
		"Successfully does not claim yield when asset yield index is already updated": {
			globalAssetYieldIndex: big.NewRat(5, 4),
			fundsInTDaiPool:       big.NewInt(1_200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
			subaccountAssetYieldIndex: big.NewRat(5, 4).String(),
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(5, 4).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000), // Yield Collected: 0 tDAI
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(1_200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
		},
		"Succesfully does not claim yield when negative positions cancel out positive position yield claims": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1000, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1000, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Negative general asset yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(-1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedErr:             types.ErrGlobalYieldIndexNegative,
			expectedAssetYieldIndex: big.NewRat(-1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Asset yield index in account higher than in general ": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 2),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedErr:             types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
			expectedAssetYieldIndex: big.NewRat(1, 2).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Negative general perp yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(-1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedErr:             types.ErrGlobalYieldIndexNegative,
			expectedAssetYieldIndex: big.NewRat(-1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Perp yield index in subaccount higher than in general": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedErr:             types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Perp yield index in subaccount badly initialized": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   "",
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedErr:             types.ErrYieldIndexUninitialized,
			expectedAssetYieldIndex: big.NewRat(1, 2).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   "",
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Successful yield claim: not enough yield in tdai pool so we take what's available": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(1),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_001),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_001,
			},
		},
		"Fails yield claim: no open positions": {
			globalAssetYieldIndex: big.NewRat(1, 1),
			fundsInTDaiPool:       big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			subaccountId:                 &defaultSubaccountId,
			expectedErr:                  types.ErrNoYieldToClaim,
			expectedAssetYieldIndex:      big.NewRat(1, 1).String(),
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: subaccountId is nil": {
			globalAssetYieldIndex: big.NewRat(1, 1),
			fundsInTDaiPool:       big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			subaccountId:                 nil,
			expectedErr:                  types.ErrSubaccountIdIsNil,
			expectedAssetYieldIndex:      big.NewRat(1, 1).String(),
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Successfully claims 0 yield when subaccount's yield is negative": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(10_000, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			subaccountId:            &defaultSubaccountId,
			expectedAssetYieldIndex: big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(10_000, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(
				t,
				true,
			)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			rateLimitKeeper.SetSDAIPrice(ctx, rate)
			globalAssetYieldIndex := big.NewRat(1, 1)
			if tc.globalAssetYieldIndex != nil {
				globalAssetYieldIndex = tc.globalAssetYieldIndex
			}
			rateLimitKeeper.SetAssetYieldIndex(ctx, globalAssetYieldIndex)

			// Always creates TDai asset first
			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Id,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
					a.AssetYieldIndex,
				)
				require.NoError(t, err)
			}

			for _, p := range tc.perpetuals {
				perpetualsKeeper.SetPerpetualForTest(
					ctx,
					p,
				)
			}

			for collateralPoolAddr, TDaiBal := range tc.collateralPoolTDaiBalances {
				err := bank_testutil.FundAccount(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewInt(TDaiBal)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.fundsInTDaiPool != nil {
				err := bank_testutil.FundModuleAccount(
					ctx,
					ratelimittypes.TDaiPoolAccount,
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.fundsInTDaiPool)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			subaccount := createNSubaccount(keeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.PerpetualPositions = tc.perpetualPositions
			subaccount.AssetPositions = tc.assetPositions
			subaccountYieldIndex := constants.AssetYieldIndex_Zero
			if tc.subaccountAssetYieldIndex != "" {
				subaccountYieldIndex = tc.subaccountAssetYieldIndex
			}
			subaccount.AssetYieldIndex = subaccountYieldIndex
			keeper.SetSubaccount(ctx, subaccount)
			subaccountId := *subaccount.Id

			err := keeper.ClaimYieldForSubaccountFromIdAndSetNewState(ctx, tc.subaccountId)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			newSubaccount := keeper.GetSubaccount(ctx, subaccountId)
			require.Equal(t, len(newSubaccount.PerpetualPositions), len(tc.expectedPerpetualPositions))
			for i, ep := range tc.expectedPerpetualPositions {
				require.Equal(t, *ep, *newSubaccount.PerpetualPositions[i])
			}
			require.Equal(t, len(newSubaccount.AssetPositions), len(tc.expectedAssetPositions))
			for i, ep := range tc.expectedAssetPositions {
				require.Equal(t, *ep, *newSubaccount.AssetPositions[i])
			}
			if tc.expectedErr == nil {
				require.Equal(t, 0, globalAssetYieldIndex.Cmp(ratelimitkeeper.ConvertStringToBigRatWithPanicOnErr(newSubaccount.AssetYieldIndex)),
					"Expected AssetYieldIndex %v. Got %v.", globalAssetYieldIndex, newSubaccount.AssetYieldIndex,
				)
			}

			for collateralPoolAddr, expectedTDaiBal := range tc.expectedCollateralPoolTDaiBalances {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewInt(expectedTDaiBal)),
					TDaiBal,
				)
			}

			if tc.expectedTDaiYieldPoolBalance != nil {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(ratelimittypes.TDaiPoolAccount),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.expectedTDaiYieldPoolBalance)),
					TDaiBal,
				)

			}
		})
	}
}
