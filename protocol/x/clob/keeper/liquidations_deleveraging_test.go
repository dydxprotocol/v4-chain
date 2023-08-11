package keeper_test

import (
	"errors"
	"math"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4/dtypes"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	assettypes "github.com/dydxprotocol/v4/x/assets/types"
	"github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/dydxprotocol/v4/x/clob/types"
	perptypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetInsuranceFundBalance(t *testing.T) {
	tests := map[string]struct {
		// Setup
		assets               []assettypes.Asset
		insuranceFundBalance *big.Int

		// Expectations.
		expectedInsuranceFundBalance uint64
		expectedError                error
	}{
		"can get zero balance": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance:         new(big.Int),
			expectedInsuranceFundBalance: 0,
		},
		"can get positive balance": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance:         new(big.Int).SetInt64(100),
			expectedInsuranceFundBalance: 100,
		},
		"can get max uint64 balance": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance:         new(big.Int).SetUint64(math.MaxUint64),
			expectedInsuranceFundBalance: math.MaxUint64,
		},
		"panics when asset not found in state": {
			assets:        []assettypes.Asset{},
			expectedError: errors.New("GetInsuranceFundBalance: Usdc asset not found in state"),
		},
		"panics when amount is greater than uint64": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				big.NewInt(1),
			),
			expectedError: errors.New("Uint64() out of bounds"),
		},
		"panics when amount is negative": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance: big.NewInt(-1),
			expectedError:        errors.New("Uint64() out of bounds"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ctx,
				clobKeeper,
				_,
				assetsKeeper,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memClob, bankMock, &mocks.IndexerEventManager{})

			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
				)
				require.NoError(t, err)
			}

			if tc.insuranceFundBalance != nil {
				if tc.insuranceFundBalance.IsUint64() {
					bankMock.On(
						"GetBalance",
						mock.Anything,
						authtypes.NewModuleAddress(types.InsuranceFundName),
						constants.Usdc.Denom,
					).Return(
						sdk.NewCoin(constants.Usdc.Denom, sdk.NewIntFromBigInt(tc.insuranceFundBalance)),
					)
				} else {
					bankMock.On(
						"GetBalance",
						mock.Anything,
						authtypes.NewModuleAddress(types.InsuranceFundName),
						constants.Usdc.Denom,
					).Panic("Uint64() out of bounds")
				}
			}

			if tc.expectedError != nil {
				require.PanicsWithValue(
					t,
					tc.expectedError.Error(),
					func() {
						clobKeeper.GetInsuranceFundBalance(ctx)
					},
				)
			} else {
				require.Equal(
					t,
					tc.expectedInsuranceFundBalance,
					clobKeeper.GetInsuranceFundBalance(ctx),
				)
			}
		})
	}
}

func TestShouldPerformDeleveraging(t *testing.T) {
	tests := map[string]struct {
		// Setup
		liquidationConfig    types.LiquidationsConfig
		insuranceFundBalance *big.Int
		insuranceFundDelta   *big.Int

		// Expectations.
		expectedShouldPerformDeleveraging bool
	}{
		"zero insurance fund delta": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(9_998_000_000), // $9,998
			insuranceFundDelta:   big.NewInt(0),

			expectedShouldPerformDeleveraging: false,
		},
		"zero insurance fund delta - insurance fund balance is greater than deleveraging threshold": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(20_000_000_000), // $20,000
			insuranceFundDelta:   big.NewInt(0),

			expectedShouldPerformDeleveraging: false,
		},
		"positive insurance fund delta": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(9_998_000_000), // $9,998
			insuranceFundDelta:   big.NewInt(1_000_000),

			expectedShouldPerformDeleveraging: false,
		},
		"positive insurance fund delta - insurance fund after applying delta is greater than deleveraging threshold": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(20_000_000_000), // $20,000
			insuranceFundDelta:   big.NewInt(1_000_000),

			expectedShouldPerformDeleveraging: false,
		},
		"negative insurance fund delta - initial balance is less than deleveraging threshold": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(9_998_000_000), // $10,000
			insuranceFundDelta:   big.NewInt(-1_000_000),

			expectedShouldPerformDeleveraging: true,
		},
		"negative insurance fund delta - initial balance is greater than deleveraging threshold": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(20_000_000_000), // $20,000
			insuranceFundDelta:   big.NewInt(-1_000_000),

			expectedShouldPerformDeleveraging: false,
		},
		"negative insurance fund delta - insurance fund balance can go from above threshold to below threshold": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(10_000_000_000), // $10,000
			insuranceFundDelta:   big.NewInt(-1_000_000),

			expectedShouldPerformDeleveraging: false,
		},
		"negative insurance fund delta - abs delta is greater than max insurance fund quantums for deleverging ": {
			liquidationConfig:    constants.LiquidationsConfig_10bMaxInsuranceFundQuantumsForDeleveraging,
			insuranceFundBalance: big.NewInt(10_000_000_000),
			insuranceFundDelta:   big.NewInt(-10_000_000_001),

			expectedShouldPerformDeleveraging: true,
		},
		"negative insurance fund delta - max insurance fund quantums for deleveraging is zero": {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_000),
			insuranceFundDelta:   big.NewInt(-10_000_000_001),

			expectedShouldPerformDeleveraging: true,
		},
		"negative insurance fund delta - max insurance fund quantums for deleveraging is max uint64": {
			liquidationConfig: types.LiquidationsConfig{
				MaxLiquidationFeePpm:                    5_000,
				MaxInsuranceFundQuantumsForDeleveraging: math.MaxUint64,
				FillablePriceConfig:                     constants.FillablePriceConfig_Default,
				PositionBlockLimits:                     constants.PositionBlockLimits_No_Limit,
				SubaccountBlockLimits:                   constants.SubaccountBlockLimits_No_Limit,
			},
			insuranceFundBalance: new(big.Int).SetUint64(math.MaxUint64 - 1),
			insuranceFundDelta:   big.NewInt(-10_000_000_000),

			expectedShouldPerformDeleveraging: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ctx,
				clobKeeper,
				_,
				assetsKeeper,
				_,
				_,
				_,
				_ := keepertest.ClobKeepers(t, memClob, bankMock, &mocks.IndexerEventManager{})

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

			// Initialize the liquidations config.
			err = clobKeeper.InitializeLiquidationsConfig(ctx, tc.liquidationConfig)
			require.NoError(t, err)

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(types.InsuranceFundName),
				constants.Usdc.Denom,
			).Return(
				sdk.NewCoin(constants.Usdc.Denom, sdk.NewIntFromBigInt(tc.insuranceFundBalance)),
			)
			require.Equal(
				t,
				tc.expectedShouldPerformDeleveraging,
				clobKeeper.ShouldPerformDeleveraging(
					ctx,
					tc.insuranceFundDelta,
				),
			)
		})
	}
}

func TestMustGetOffsettingSubaccountsForDeleveraging(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		subaccounts []satypes.Subaccount

		// Parameters.
		liquidatedSubaccountId satypes.SubaccountId
		perpetualId            uint32
		deltaQuantums          *big.Int

		// Expectations.
		expectedSubaccounts []satypes.SubaccountId
		panics              bool
	}{
		"Can get one offsetting subaccount": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts:    []satypes.SubaccountId{constants.Dave_Num0},
		},
		"Can get multiple offsetting subaccounts": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				{
					Id:             &constants.Dave_Num0,
					AssetPositions: keepertest.CreateUsdcAssetPosition(big.NewInt(50_000_000_000)),
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(50_000_000), // 0.5 BTC
						},
					},
				},
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.Usdc_Asset_50_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(50_000_000), // 0.5 BTC
						},
					},
				},
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts:    []satypes.SubaccountId{constants.Dave_Num1, constants.Dave_Num0},
		},
		"Skips subaccounts with positions on the same side": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Carl_Num1_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts:    []satypes.SubaccountId{constants.Dave_Num0},
		},
		"Skips subaccounts with different perpetual position": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				{
					Id:             &constants.Dave_Num1,
					AssetPositions: keepertest.CreateUsdcAssetPosition(big.NewInt(50_000_000_000)),
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 1,                            // Different perpetual.
							Quantums:    dtypes.NewInt(1_000_000_000), // -1 BTC
						},
					},
				},
				constants.Dave_Num0_1BTC_Long,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts:    []satypes.SubaccountId{constants.Dave_Num0},
		},
		"Returns empty slice if subaccount to be deleveraged has negative net collateral": {
			subaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
					// Carl_Num0 has negative net collateral.
					AssetPositions: keepertest.CreateUsdcAssetPosition(big.NewInt(49_999_999_999)),
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
						},
					},
				},
				constants.Dave_Num0_1BTC_Long,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts:    nil,
		},
		"Skips subaccounts with negative net collateral": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				{
					Id: &constants.Dave_Num1,
					// Dave_Num1 has negative net collateral.
					AssetPositions: keepertest.CreateUsdcAssetPosition(big.NewInt(-50_000_000_001)),
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(100_000_000), // -1 BTC
						},
					},
				},
				constants.Dave_Num0_1BTC_Long,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts:    []satypes.SubaccountId{constants.Dave_Num0},
		},
		"Panics when deltaQuantums is zero": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short,
				constants.Dave_Num0_1BTC_Long,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(0),
			panics:                 true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			ctx,
				clobKeeper,
				pricesKeeper,
				assetsKeeper,
				perpKeeper,
				subaccountsKeeper,
				_,
				_ := keepertest.ClobKeepers(t, memClob, &mocks.BankKeeper{}, &mocks.IndexerEventManager{})

			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpKeeper)

			err := keepertest.CreateUsdcAsset(ctx, assetsKeeper)
			require.NoError(t, err)

			for _, p := range []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
			} {
				_, err := perpKeeper.CreatePerpetual(
					ctx,
					p.Ticker,
					p.MarketId,
					p.AtomicResolution,
					p.DefaultFundingPpm,
					p.LiquidityTier,
				)
				require.NoError(t, err)
			}

			for _, subaccount := range tc.subaccounts {
				subaccountsKeeper.SetSubaccount(ctx, subaccount)
			}

			if tc.panics {
				require.Panics(t, func() {
					clobKeeper.MustGetOffsettingSubaccountsForDeleveraging(
						ctx,
						tc.liquidatedSubaccountId,
						tc.perpetualId,
						tc.deltaQuantums,
					)
				})
			} else {
				offsettingSubaccounts := clobKeeper.MustGetOffsettingSubaccountsForDeleveraging(
					ctx,
					tc.liquidatedSubaccountId,
					tc.perpetualId,
					tc.deltaQuantums,
				)
				require.Equal(t, tc.expectedSubaccounts, offsettingSubaccounts)
			}
		})
	}
}
