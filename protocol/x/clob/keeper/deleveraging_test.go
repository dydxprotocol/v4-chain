package keeper_test

import (
	"errors"
	"math"
	"math/big"
	"testing"
	"time"

	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/indexer_manager"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"

	sdkmath "cosmossdk.io/math"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	clobtest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/clob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	perptest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/perpetuals"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/memclob"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetInsuranceFundBalanceInQuoteQuantums(t *testing.T) {
	tests := map[string]struct {
		// Setup
		assets               []assettypes.Asset
		insuranceFundBalance *big.Int
		perpetualId          uint32
		perpetual            *perptypes.Perpetual

		// Expectations.
		expectedInsuranceFundBalance *big.Int
		expectedError                error
	}{
		"can get zero balance": {
			assets: []assettypes.Asset{
				*constants.TDai,
			},
			perpetualId:                  0,
			insuranceFundBalance:         new(big.Int),
			expectedInsuranceFundBalance: big.NewInt(0),
		},
		"can get positive balance": {
			assets: []assettypes.Asset{
				*constants.TDai,
			},
			perpetualId:                  0,
			insuranceFundBalance:         big.NewInt(100),
			expectedInsuranceFundBalance: big.NewInt(100),
		},
		"can get greater than MaxUint64 balance": {
			assets: []assettypes.Asset{
				*constants.TDai,
			},
			perpetualId: 0,
			insuranceFundBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				new(big.Int).SetUint64(math.MaxUint64),
			),
			expectedInsuranceFundBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				new(big.Int).SetUint64(math.MaxUint64),
			),
		},
		"can get zero balance - isolated market": {
			assets: []assettypes.Asset{
				*constants.TDai,
			},
			perpetualId:                  3, // Isolated market.
			insuranceFundBalance:         new(big.Int),
			expectedInsuranceFundBalance: big.NewInt(0),
		},
		"can get positive balance - isolated market": {
			assets: []assettypes.Asset{
				*constants.TDai,
			},
			perpetualId:                  3, // Isolated market.
			insuranceFundBalance:         big.NewInt(100),
			expectedInsuranceFundBalance: big.NewInt(100),
		},
		"panics when asset not found in state": {
			assets:        []assettypes.Asset{},
			perpetualId:   0,
			expectedError: errors.New("GetInsuranceFundBalanceInQuoteQuantums: TDai asset not found in state"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

			ctx := ks.Ctx.WithIsCheckTx(true)
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			keepertest.CreateTestPerpetuals(t, ctx, ks.PerpetualsKeeper)

			for _, a := range tc.assets {
				_, err := ks.AssetsKeeper.CreateAsset(
					ks.Ctx,
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

			insuranceFundAddr, err := ks.PerpetualsKeeper.GetInsuranceFundModuleAddress(ks.Ctx, tc.perpetualId)
			require.NoError(t, err)

			if tc.insuranceFundBalance != nil {
				bankMock.On(
					"GetBalance",
					mock.Anything,
					insuranceFundAddr,
					constants.TDai.Denom,
				).Return(
					sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(tc.insuranceFundBalance)),
				)
			}

			if tc.expectedError != nil {
				require.PanicsWithValue(
					t,
					tc.expectedError.Error(),
					func() {
						ks.ClobKeeper.GetInsuranceFundBalanceInQuoteQuantums(ks.Ctx, tc.perpetualId)
					},
				)
			} else {
				require.Equal(
					t,
					tc.expectedInsuranceFundBalance,
					ks.ClobKeeper.GetInsuranceFundBalanceInQuoteQuantums(ks.Ctx, tc.perpetualId),
				)
			}
		})
	}
}

func TestIsValidInsuranceFundDelta(t *testing.T) {
	tests := map[string]struct {
		// Setup
		insuranceFundBalance *big.Int
		insuranceFundDelta   *big.Int

		// Expectations.
		expectedIsValidInsuranceFundDelta bool
	}{
		"valid: zero insurance fund delta": {
			insuranceFundBalance: big.NewInt(9_998_000_000), // $9,998
			insuranceFundDelta:   big.NewInt(0),

			expectedIsValidInsuranceFundDelta: true,
		},
		"valid: zero insurance fund delta and zero balance": {
			insuranceFundBalance: big.NewInt(0), // $0
			insuranceFundDelta:   big.NewInt(0),

			expectedIsValidInsuranceFundDelta: true,
		},
		"valid: positive insurance fund delta": {
			insuranceFundBalance: big.NewInt(9_998_000_000), // $9,998
			insuranceFundDelta:   big.NewInt(1_000_000),

			expectedIsValidInsuranceFundDelta: true,
		},
		"valid: positive insurance fund delta and zero balance": {
			insuranceFundBalance: big.NewInt(0), // $0
			insuranceFundDelta:   big.NewInt(1_000_000),

			expectedIsValidInsuranceFundDelta: true,
		},
		"valid: negative insurance fund delta - insurance fund is still positive after delta": {
			insuranceFundBalance: big.NewInt(9_998_000_000), // $10,000
			insuranceFundDelta:   big.NewInt(-1_000_000),

			expectedIsValidInsuranceFundDelta: true,
		},
		"valid: negative insurance fund delta - insurance fund has zero balance after delta": {
			insuranceFundBalance: big.NewInt(10_000_000_000),
			insuranceFundDelta:   big.NewInt(-10_000_000_000),

			expectedIsValidInsuranceFundDelta: true,
		},
		"invalid: negative insurance fund delta - insurance fund is negative after delta": {
			insuranceFundBalance: big.NewInt(10_000_000_000),
			insuranceFundDelta:   big.NewInt(-10_000_000_001),

			expectedIsValidInsuranceFundDelta: false,
		},
		"invalid: negative insurance fund delta - insurance fund was empty and is negative after delta": {
			insuranceFundBalance: big.NewInt(0),
			insuranceFundDelta:   big.NewInt(-1),

			expectedIsValidInsuranceFundDelta: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			ctx := ks.Ctx.WithIsCheckTx(true)
			keepertest.CreateTestMarkets(t, ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, ks.PerpetualsKeeper)

			keepertest.CreateTestPerpetuals(t, ctx, ks.PerpetualsKeeper)

			bankMock.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(tc.insuranceFundBalance)),
			)
			require.Equal(
				t,
				tc.expectedIsValidInsuranceFundDelta,
				ks.ClobKeeper.IsValidInsuranceFundDelta(
					ks.Ctx,
					tc.insuranceFundDelta,
					0,
				),
			)
		})
	}
}

func TestCanDeleverageSubaccount(t *testing.T) {
	tests := map[string]struct {
		// Setup
		liquidationConfig             types.LiquidationsConfig
		insuranceFundBalance          *big.Int
		subaccount                    satypes.Subaccount
		marketIdToOraclePriceOverride map[uint32]uint64
		clobPairs                     []types.ClobPair

		// Expectations.
		expectedShouldDeleverageAtBankruptcyPrice bool
		expectedShouldDeleverageAtOraclePrice     bool
	}{
		`Cannot deleverage when subaccount has positive TNC`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_001), // $10,000.000001
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_000_000_000, // $50,000 / BTC
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedShouldDeleverageAtBankruptcyPrice: false,
			expectedShouldDeleverageAtOraclePrice:     false,
		},
		`Cannot deleverage when subaccount has zero TNC`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_001), // $10,000.000001
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_499_000_000, // $54,999 / BTC
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedShouldDeleverageAtBankruptcyPrice: false,
			expectedShouldDeleverageAtOraclePrice:     false,
		},
		`Can deleverage when subaccount has negative TNC`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_000), // $10,000
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_500_000_000, // $55,000 / BTC
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc,
			},

			expectedShouldDeleverageAtBankruptcyPrice: true,
			expectedShouldDeleverageAtOraclePrice:     false,
		},
		`Can deleverage when subaccount has negative TNC and clob pair has status FINAL_SETTLEMENT`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_000), // $10,000
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_500_000_000, // $55,000 / BTC
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},

			expectedShouldDeleverageAtBankruptcyPrice: true,
			expectedShouldDeleverageAtOraclePrice:     false,
		},
		`Can final settle deleverage when subaccount has positive TNC and clob pair has status FINAL_SETTLEMENT`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_001), // $10,000.000001
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_000_000_000, // $50,000 / BTC
			},
			clobPairs: []types.ClobPair{
				constants.ClobPair_Btc_Final_Settlement,
			},

			expectedShouldDeleverageAtBankruptcyPrice: false,
			expectedShouldDeleverageAtOraclePrice:     true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, mockIndexerEventManager)

			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(0, 1))

			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Initialize the liquidations config.
			err = ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, tc.liquidationConfig)
			require.NoError(t, err)

			bankMock.On(
				"GetBalance",
				mock.Anything,
				perptypes.InsuranceFundModuleAddress,
				constants.TDai.Denom,
			).Return(
				sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(tc.insuranceFundBalance)),
			)
			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create test markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Update the prices on the test markets.
			for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
				err := ks.PricesKeeper.UpdateSpotAndPnlMarketPrices(
					ks.Ctx,
					&pricestypes.MarketPriceUpdate{
						MarketId:  marketId,
						SpotPrice: oraclePrice,
						PnlPrice:  oraclePrice,
					},
				)
				require.NoError(t, err)
			}

			perpetuals := []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
			}
			for _, perpetual := range perpetuals {
				_, err = ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					perpetual.Params.Id,
					perpetual.Params.Ticker,
					perpetual.Params.MarketId,
					perpetual.Params.AtomicResolution,
					perpetual.Params.DefaultFundingPpm,
					perpetual.Params.LiquidityTier,
					perpetual.Params.MarketType,
					perpetual.Params.DangerIndexPpm,
					perpetual.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					perpetual.YieldIndex,
				)
				require.NoError(t, err)
			}

			for i, clobPair := range tc.clobPairs {
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexerevents.PerpetualMarketEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewPerpetualMarketCreateEvent(
							clobPair.MustGetPerpetualId(),
							clobPair.Id,
							perpetuals[i].Params.Ticker,
							perpetuals[i].Params.MarketId,
							clobPair.Status,
							clobPair.QuantumConversionExponent,
							perpetuals[i].Params.AtomicResolution,
							clobPair.SubticksPerTick,
							clobPair.StepBaseQuantums,
							perpetuals[i].Params.LiquidityTier,
							perpetuals[i].Params.MarketType,
							perpetuals[i].Params.DangerIndexPpm,
							perpetuals[i].Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
						),
					),
				).Once().Return()

				_, err = ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					clobPair.Id,
					clobPair.MustGetPerpetualId(),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.subaccount)

			shouldDeleverageAtBankruptcyPrice, shouldDeleverageAtOraclePrice, err := ks.ClobKeeper.CanDeleverageSubaccount(
				ks.Ctx,
				*tc.subaccount.Id,
				0,
			)
			require.NoError(t, err)
			require.Equal(
				t,
				tc.expectedShouldDeleverageAtBankruptcyPrice,
				shouldDeleverageAtBankruptcyPrice,
			)
			require.Equal(
				t,
				tc.expectedShouldDeleverageAtOraclePrice,
				shouldDeleverageAtOraclePrice,
			)
		})
	}
}

func TestOffsetSubaccountPerpetualPosition(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		subaccounts []satypes.Subaccount

		// Parameters.
		liquidatedSubaccountId satypes.SubaccountId
		perpetualId            uint32
		deltaQuantums          *big.Int

		// Expectations.
		expectedSubaccounts       []satypes.Subaccount
		expectedFills             []types.MatchPerpetualDeleveraging_Fill
		expectedQuantumsRemaining *big.Int
		// Expected remaining OI after test.
		// The test initializes each perp with default open interest of 1 full coin.
		expectedOpenInterest *big.Int
	}{
		"Can get one offsetting subaccount for deleveraged short": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 1 BTC short is $54,999 and we close both positions at this price.
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(50_000_000_000 + 54_999_000_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
			expectedOpenInterest:      new(big.Int), // fully deleveraged
		},
		"Can get one offsetting subaccount for deleveraged long": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Long_54999USD,
				constants.Dave_Num0_1BTC_Short_100000USD,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(-100_000_000),
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(100_000_000_000 - 54_999_000_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
			expectedOpenInterest:      new(big.Int), // fully deleveraged
		},
		"Can get multiple offsetting subaccounts": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				{
					Id: &constants.Dave_Num0,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_50_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(50_000_000), // 0.5 BTC
						},
					},
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num1,
					AssetPositions: []*satypes.AssetPosition{
						&constants.TDai_Asset_50_000,
					},
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId: 0,
							Quantums:    dtypes.NewInt(50_000_000), // 0.5 BTC
						},
					},
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 0.5 BTC short is $27,499.5 and we close both positions at this price.
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(50_000_000_000 + 27_499_500_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num1,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 0.5 BTC short is $27,499.5 and we close both positions at this price.
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(50_000_000_000 + 27_499_500_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             50_000_000,
				},
				{
					OffsettingSubaccountId: constants.Dave_Num1,
					FillAmount:             50_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
			expectedOpenInterest:      new(big.Int), // fully deleveraged
		},
		"Skips subaccounts with positions on the same side": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Carl_Num1_1BTC_Short,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id:                 constants.Carl_Num1_1BTC_Short.Id,
					AssetPositions:     constants.Carl_Num1_1BTC_Short.AssetPositions,
					PerpetualPositions: constants.Carl_Num1_1BTC_Short.PerpetualPositions,
					MarginEnabled:      constants.Carl_Num1_1BTC_Short.MarginEnabled,
					AssetYieldIndex:    big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 1 BTC short is $54,999 and we close both positions at this price.
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(50_000_000_000 + 54_999_000_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
			expectedOpenInterest:      new(big.Int), // fully deleveraged
		},
		"Skips subaccounts with no open position for the given perpetual": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				constants.Dave_Num1_1ETH_Long_50000USD, // ETH
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id:                 constants.Dave_Num1_1ETH_Long_50000USD.Id,
					AssetPositions:     constants.Dave_Num1_1ETH_Long_50000USD.AssetPositions,
					PerpetualPositions: constants.Dave_Num1_1ETH_Long_50000USD.PerpetualPositions,
					MarginEnabled:      constants.Dave_Num1_1ETH_Long_50000USD.MarginEnabled,
					AssetYieldIndex:    big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 1 BTC short is $54,999 and we close both positions at this price.
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(50_000_000_000 + 54_999_000_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
			expectedOpenInterest:      new(big.Int), // fully deleveraged
		},
		"Skips subaccounts with non-overlapping bankruptcy prices": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50000USD,
				constants.Dave_Num0_1BTC_Long_50001USD_Short,
				constants.Dave_Num1_1BTC_Long_50000USD,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id:              &constants.Carl_Num0,
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id:                 constants.Dave_Num0_1BTC_Long_50001USD_Short.Id,
					AssetPositions:     constants.Dave_Num0_1BTC_Long_50001USD_Short.AssetPositions,
					PerpetualPositions: constants.Dave_Num0_1BTC_Long_50001USD_Short.PerpetualPositions,
					MarginEnabled:      constants.Dave_Num0_1BTC_Long_50001USD_Short.MarginEnabled,
					AssetYieldIndex:    big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num1,
					// TNC of liquidated subaccount is $0, which means the bankruptcy price
					// to close 1 BTC short is $50,000 and we close both positions at this price.
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(50_000_000_000 + 50_000_000_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num1,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
			expectedOpenInterest:      new(big.Int), // fully deleveraged
		},
		"Returns an error if not enough subaccounts to fully deleverage liquidated subaccount's position": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_50000USD,
				constants.Dave_Num0_1BTC_Long_50001USD_Short,
			},
			liquidatedSubaccountId:    constants.Carl_Num0,
			perpetualId:               0,
			deltaQuantums:             big.NewInt(100_000_000),
			expectedSubaccounts:       nil,
			expectedFills:             []types.MatchPerpetualDeleveraging_Fill{},
			expectedQuantumsRemaining: big.NewInt(100_000_000),
			expectedOpenInterest:      big.NewInt(100_000_000),
		},
		"Can offset subaccount with multiple positions, first position is offset leaving TNC constant": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_1ETH_Long_47000USD,
				constants.Dave_Num0_1BTC_Long_50000USD,
			},
			liquidatedSubaccountId: constants.Carl_Num0,
			perpetualId:            0,
			deltaQuantums:          big.NewInt(100_000_000),
			expectedSubaccounts: []satypes.Subaccount{
				// Carl's BTC short position is offset by Dave's BTC long position at $50,000 leaving
				// his ETH long position untouched and dropping his asset position to -$3000.
				{
					Id: &constants.Carl_Num0,
					PerpetualPositions: []*satypes.PerpetualPosition{
						{
							PerpetualId:  1,
							Quantums:     dtypes.NewInt(1_000_000_000), // 1 ETH
							FundingIndex: dtypes.NewInt(0),
							YieldIndex:   big.NewRat(0, 1).String(),
						},
					},
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(-3_000_000_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: keepertest.CreateTDaiAssetPosition(
						big.NewInt(50_000_000_000 + 50_000_000_000),
					),
					AssetYieldIndex: big.NewRat(0, 1).String(),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: big.NewInt(0),
			expectedOpenInterest:      new(big.Int), // fully deleveraged
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, mockIndexerEventManager)
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(0, 1))

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			perps := []perptypes.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
				constants.EthUsd_100PercentMarginRequirement,
			}
			for _, p := range perps {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				perps,
			)

			clobPairs := []types.ClobPair{
				constants.ClobPair_Btc,
				constants.ClobPair_Eth,
			}
			for i, clobPair := range clobPairs {
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypePerpetualMarket,
					indexerevents.PerpetualMarketEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewPerpetualMarketCreateEvent(
							clobPair.MustGetPerpetualId(),
							clobPair.Id,
							perps[i].Params.Ticker,
							perps[i].Params.MarketId,
							clobPair.Status,
							clobPair.QuantumConversionExponent,
							perps[i].Params.AtomicResolution,
							clobPair.SubticksPerTick,
							clobPair.StepBaseQuantums,
							perps[i].Params.LiquidityTier,
							perps[i].Params.MarketType,
							perps[i].Params.DangerIndexPpm,
							perps[i].Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
						),
					),
				).Once().Return()

				_, err = ks.ClobKeeper.CreatePerpetualClobPair(
					ks.Ctx,
					clobPair.Id,
					clobPair.MustGetPerpetualId(),
					satypes.BaseQuantums(clobPair.StepBaseQuantums),
					clobPair.QuantumConversionExponent,
					clobPair.SubticksPerTick,
					clobPair.Status,
				)
				require.NoError(t, err)
			}

			for _, subaccount := range tc.subaccounts {
				ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, subaccount)
			}

			ks.BlockTimeKeeper.SetPreviousBlockInfo(ks.Ctx, &blocktimetypes.BlockInfo{
				Timestamp: time.Unix(5, 0),
			})
			// check that an event is emitted per fill
			for _, fill := range tc.expectedFills {
				fillAmount := new(big.Int).SetUint64(fill.FillAmount)
				if tc.deltaQuantums.Sign() < 0 {
					fillAmount = new(big.Int).Neg(fillAmount)
				}
				bankruptcyPriceQuoteQuantums, err := ks.ClobKeeper.GetBankruptcyPriceInQuoteQuantums(
					ks.Ctx,
					tc.liquidatedSubaccountId,
					tc.perpetualId,
					fillAmount,
				)
				require.NoError(t, err)
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypeDeleveraging,
					indexerevents.DeleveragingEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewDeleveragingEvent(
							tc.liquidatedSubaccountId,
							fill.OffsettingSubaccountId,
							tc.perpetualId,
							satypes.BaseQuantums(fill.FillAmount),
							satypes.BaseQuantums(bankruptcyPriceQuoteQuantums.Uint64()),
							tc.deltaQuantums.Sign() > 0,
							false,
						),
					),
				).Return()
			}

			positions := clobtest.GetOpenPositionsFromSubaccounts(tc.subaccounts)
			ks.ClobKeeper.DaemonDeleveragingInfo.UpdateSubaccountsWithPositions(positions)
			fills, deltaQuantumsRemaining := ks.ClobKeeper.OffsetSubaccountPerpetualPosition(
				ks.Ctx,
				tc.liquidatedSubaccountId,
				tc.perpetualId,
				tc.deltaQuantums,
				false, // TODO, add tests where final settlement is true
			)
			require.Equal(t, tc.expectedFills, fills)
			require.True(t, tc.expectedQuantumsRemaining.Cmp(deltaQuantumsRemaining) == 0)

			for _, subaccount := range tc.expectedSubaccounts {
				require.Equal(t, subaccount, ks.SubaccountsKeeper.GetSubaccount(ks.Ctx, *subaccount.Id))
			}

			if tc.expectedOpenInterest != nil {
				gotPerp, err := ks.PerpetualsKeeper.GetPerpetual(ks.Ctx, tc.perpetualId)
				require.NoError(t, err)
				require.Zero(t,
					tc.expectedOpenInterest.Cmp(gotPerp.OpenInterest.BigInt()),
					"expected open interest %s, got %s",
					tc.expectedOpenInterest.String(),
					gotPerp.OpenInterest.String(),
				)
			}
		})
	}
}

func TestProcessDeleveraging(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		liquidatedSubaccount satypes.Subaccount
		offsettingSubaccount satypes.Subaccount
		deltaQuantums        *big.Int

		// Expectations.
		expectedLiquidatedSubaccount satypes.Subaccount
		expectedOffsettingSubaccount satypes.Subaccount
		expectedErr                  error
	}{
		// Categorizing subaccounts into four groups:
		// 1. Well-collateralized
		// 2. Liquidatable, but TNC > 0
		// 3. Liquidatable, TNC == 0
		// 4. Liquidatable, TNC < 0
		//
		// Here, we construct table tests for only 3x4 permutations of the above groups
		// since liquidatedSubaccount shouldn't be well-collateralized.
		"Liquidated: under-collateralized, TNC > 0, offsetting: well-collateralized": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(50_000_000_000 + 54_999_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC > 0, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(-45_001_000_000 + 54_999_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC > 0, offsetting: under-collateralized, TNC == 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(-50_000_000_000 + 54_999_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC > 0, offsetting: under-collateralized, TNC < 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(-50_001_000_000 + 54_999_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: well-collateralized": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $0, which means the bankruptcy price
				// to close 1 BTC short is $50,000 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $0, which means the bankruptcy price
				// to close 1 BTC short is $50,000 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(-45_001_000_000 + 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: under-collateralized, TNC == 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $0, which means the bankruptcy price
				// to close 1 BTC short is $50,000 and we close both positions at this price.
				// TDai of this suabccount is -$50,000 + $50,000 = $0.
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: under-collateralized, TNC < 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			// TNC of liquidated subaccount is $0, which means the bankruptcy price
			// to close 1 BTC short is $50,000.
			// TNC of offsetting subaccount is $-1, which means the bankruptcy price
			// to close 1 BTC long is $50,001.
			// Since the bankruptcy prices do not overlap,
			// i.e. bankruptcy price of long > bankruptcy price of short,
			// state transitions aren't valid.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
		"Liquidated: under-collateralized, TNC < 0, offsetting: well-collateralized": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $-1, which means the bankruptcy price
				// to close 1 BTC short is $49,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(50_000_000_000 + 49_999_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC < 0, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id:              &constants.Carl_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $-1, which means the bankruptcy price
				// to close 1 BTC short is $49,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(-45_001_000_000 + 49_999_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC < 0, offsetting: under-collateralized, TNC == 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			// TNC of liquidated subaccount is $-1, which means the bankruptcy price
			// to close 1 BTC short is $49,999.
			// TNC of offsetting subaccount is $0, which means the bankruptcy price
			// to close 1 BTC long is $50,000.
			// Since the bankruptcy prices do not overlap,
			// i.e. bankruptcy price of long > bankruptcy price of short,
			// state transitions aren't valid.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
		"Liquidated: under-collateralized, TNC < 0, offsetting: under-collateralized, TNC < 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			// TNC of liquidated subaccount is $-1, which means the bankruptcy price
			// to close 1 BTC short is $49,999.
			// TNC of offsetting subaccount is $-1, which means the bankruptcy price
			// to close 1 BTC long is $50,001.
			// Since the bankruptcy prices do not overlap,
			// i.e. bankruptcy price of long > bankruptcy price of short, state transitions aren't valid.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
		`Liquidated: under-collateralized, TNC > 0, offsetting: well-collateralized -
		can deleverage a partial position`: {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(10_000_000), // 0.1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(54_999_000_000 - 5_499_900_000),
				),
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-90_000_000), // -0.9 BTC
						FundingIndex: dtypes.ZeroInt(),
						YieldIndex:   big.NewRat(0, 1).String(),
					},
				},
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 0.1 BTC short is $5,499.9 and we close both positions at this price.
					big.NewInt(50_000_000_000 + 5_499_900_000),
				),
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(90_000_000), // 0.9 BTC
						FundingIndex: dtypes.ZeroInt(),
						YieldIndex:   big.NewRat(0, 1).String(),
					},
				},
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		`Liquidated: under-collateralized, TNC < 0, offsetting: under-collateralized, TNC < 0 -
		can not deleverage paritial positions`: {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50001USD_Short,
			deltaQuantums:        big.NewInt(10_000_000), // 0.1 BTC

			// TNC of liquidated subaccount is $-1, which means the bankruptcy price
			// to close 1 BTC short is $49,999.
			// TNC of offsetting subaccount is $-1, which means the bankruptcy price
			// to close 1 BTC long is $50,001.
			// Since the bankruptcy prices do not overlap,
			// i.e. bankruptcy price of long > bankruptcy price of short,
			// state transitions aren't valid.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
		`Liquidated: under-collateralized, TNC > 0, offsetting: well-collatearlized -
		can deleverage when there are multiple positions`: {
			liquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: []*satypes.AssetPosition{
					{
						AssetId:  0,
						Quantums: dtypes.NewInt(80_800_000_000), // $80,800
					},
				},
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId: 0,
						Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
					},
					{
						PerpetualId: 1,
						Quantums:    dtypes.NewInt(-10_000_000_000), // -10 ETH
					},
				},
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					// TNC of liquidated subaccount is $800, MMR(BTC) = $5,000, MMR(ETH) = $3,000,
					// which means the bankruptcy price to close 1 BTC short is $50,500
					// and we close both positions at this price.
					big.NewInt(80_800_000_000 - 50_500_000_000),
				),
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId:  1,
						Quantums:     dtypes.NewInt(-10_000_000_000), // -10 ETH
						FundingIndex: dtypes.ZeroInt(),
						YieldIndex:   big.NewRat(0, 1).String(),
					},
				},
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(50_000_000_000 + 50_500_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Fails when deltaQuantums is invalid with respect to liquidated subaccounts's position side": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(-100_000_000), // -1 BTC

			expectedErr: types.ErrInvalidPerpetualPositionSizeDelta,
		},
		"Fails when deltaQuantums is invalid with respect to liquidated subaccounts's position size": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(500_000_000), // 5 BTC

			expectedErr: types.ErrInvalidPerpetualPositionSizeDelta,
		},
		"Fails when deltaQuantums is invalid with respect to offsetting subaccounts's position side": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Carl_Num1_1BTC_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedErr: types.ErrInvalidPerpetualPositionSizeDelta,
		},
		"Fails when deltaQuantums is invalid with respect to offsetting subaccounts's position size": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_01BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedErr: types.ErrInvalidPerpetualPositionSizeDelta,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, mockIndexerEventManager)
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(0, 1))

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			testPerps := []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			}
			for _, p := range testPerps {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				testPerps,
			)

			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.liquidatedSubaccount)
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.offsettingSubaccount)

			bankruptcyPriceQuoteQuantums := new(big.Int)
			if tc.expectedErr == nil {
				bankruptcyPriceQuoteQuantums, err = ks.ClobKeeper.GetBankruptcyPriceInQuoteQuantums(
					ks.Ctx,
					*tc.liquidatedSubaccount.GetId(),
					uint32(0),
					tc.deltaQuantums,
				)
				require.NoError(t, err)

				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypeDeleveraging,
					indexerevents.DeleveragingEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewDeleveragingEvent(
							*tc.liquidatedSubaccount.GetId(),
							*tc.offsettingSubaccount.GetId(),
							uint32(0),
							satypes.BaseQuantums(new(big.Int).Abs(tc.deltaQuantums).Uint64()),
							satypes.BaseQuantums(bankruptcyPriceQuoteQuantums.Uint64()),
							tc.deltaQuantums.Sign() > 0,
							false,
						),
					),
				).Return()
			}
			err = ks.ClobKeeper.ProcessDeleveraging(
				ks.Ctx,
				*tc.liquidatedSubaccount.GetId(),
				*tc.offsettingSubaccount.GetId(),
				uint32(0),
				tc.deltaQuantums,
				bankruptcyPriceQuoteQuantums,
			)
			if tc.expectedErr == nil {
				require.NoError(t, err)

				actualLiquidated := ks.SubaccountsKeeper.GetSubaccount(ks.Ctx, *tc.liquidatedSubaccount.GetId())
				require.Equal(
					t,
					tc.expectedLiquidatedSubaccount,
					actualLiquidated,
				)

				actualOffsetting := ks.SubaccountsKeeper.GetSubaccount(ks.Ctx, *tc.offsettingSubaccount.GetId())
				require.Equal(
					t,
					tc.expectedOffsettingSubaccount,
					actualOffsetting,
				)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

// Note that final settlement matches piggyback off of the deleveraging operation. Because of this
// the pair of subaccounts offsetting each other are still referred to as "liquidated subaccount" and
// "offsetting subaccount" in the test cases below.
func TestProcessDeleveragingAtOraclePrice(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		liquidatedSubaccount satypes.Subaccount
		offsettingSubaccount satypes.Subaccount
		deltaQuantums        *big.Int

		// Expectations.
		expectedLiquidatedSubaccount satypes.Subaccount
		expectedOffsettingSubaccount satypes.Subaccount
		expectedErr                  error
	}{
		"Liquidated: well-collateralized, offsetting: well-collateralized": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_100000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(100_000_000_000 - 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: well-collateralized, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			offsettingSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			deltaQuantums:        big.NewInt(-100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(54_999_000_000 - 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: well-collateralized, offsetting: under-collateralized, TNC == 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_100000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(100_000_000_000 - 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id:              &constants.Dave_Num0,
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: well-collateralized, offsetting: under-collateralized, TNC < 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_100000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			// Negative TNC account closing at oracle price is an invalid state transition.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
		"Liquidated: under-collateralized, TNC > 0, offsetting: well-collateralized": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(54_999_000_000 - 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateTDaiAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
				AssetYieldIndex: big.NewRat(0, 1).String(),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: under-collateralized, TNC < 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			// Negative TNC account closing at oracle price is an invalid state transition.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
		"Liquidated: under-collateralized, TNC < 0, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			// Negative TNC account closing at oracle price is an invalid state transition.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
		"Liquidated: under-collateralized, TNC < 0, offsetting: well-collateralized": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			// Negative TNC account closing at oracle price is an invalid state transition.
			expectedErr: satypes.ErrFailedToUpdateSubaccounts,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, mockIndexerEventManager)
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(0, 1))

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			testPerps := []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			}
			for _, p := range testPerps {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				testPerps,
			)

			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.liquidatedSubaccount)
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.offsettingSubaccount)

			fillPriceQuoteQuantums, err := ks.PerpetualsKeeper.GetNetNotional(
				ks.Ctx,
				uint32(0),
				tc.deltaQuantums,
			)
			fillPriceQuoteQuantums.Neg(fillPriceQuoteQuantums)
			require.NoError(t, err)

			if tc.expectedErr == nil {
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypeDeleveraging,
					indexerevents.DeleveragingEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewDeleveragingEvent(
							*tc.liquidatedSubaccount.GetId(),
							*tc.offsettingSubaccount.GetId(),
							uint32(0),
							satypes.BaseQuantums(new(big.Int).Abs(tc.deltaQuantums).Uint64()),
							satypes.BaseQuantums(fillPriceQuoteQuantums.Uint64()),
							tc.deltaQuantums.Sign() > 0,
							false,
						),
					),
				).Return()
			}
			err = ks.ClobKeeper.ProcessDeleveraging(
				ks.Ctx,
				*tc.liquidatedSubaccount.GetId(),
				*tc.offsettingSubaccount.GetId(),
				uint32(0),
				tc.deltaQuantums,
				fillPriceQuoteQuantums,
			)
			if tc.expectedErr == nil {
				require.NoError(t, err)

				actualLiquidated := ks.SubaccountsKeeper.GetSubaccount(ks.Ctx, *tc.liquidatedSubaccount.GetId())
				require.Equal(
					t,
					tc.expectedLiquidatedSubaccount,
					actualLiquidated,
				)

				actualOffsetting := ks.SubaccountsKeeper.GetSubaccount(ks.Ctx, *tc.offsettingSubaccount.GetId())
				require.Equal(
					t,
					tc.expectedOffsettingSubaccount,
					actualOffsetting,
				)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}

func TestProcessDeleveraging_Rounding(t *testing.T) {
	tests := map[string]struct {
		// Setup.
		liquidatedSubaccount satypes.Subaccount
		offsettingSubaccount satypes.Subaccount
		deltaQuantums        *big.Int

		// Expectations.
		expectedErr error
	}{
		// Rounding tests.
		"Can deleverage short positions correctly after rounding": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(49_999_991),
		},
		"Can deleverage long position correctly after rounding": {
			liquidatedSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			offsettingSubaccount: constants.Carl_Num0_1BTC_Short_100000USD,
			deltaQuantums:        big.NewInt(-49_999_991),
		},
		"Can deleverage short positions correctly after rounding - negative TNC": {
			liquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: []*satypes.AssetPosition{
					{
						AssetId:  0,
						Quantums: dtypes.NewInt(45_001_000_000), // $45,001, TNC = -$4,999
					},
				},
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId: 0,
						Quantums:    dtypes.NewInt(-100_000_000), // -1 BTC
					},
				},
			},
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(49_999_991),
		},
		"Can deleverage long positions correctly after rounding - negative TNC": {
			liquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: []*satypes.AssetPosition{
					{
						AssetId:  0,
						Quantums: dtypes.NewInt(-50_000_000_000 - 4_999_000_000),
					},
				},
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId: 0,
						Quantums:    dtypes.NewInt(100_000_000),
					},
				},
			},
			offsettingSubaccount: constants.Carl_Num0_1BTC_Short_100000USD,
			deltaQuantums:        big.NewInt(-49_999_991),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, mockIndexerEventManager)
			ks.RatelimitKeeper.SetAssetYieldIndex(ks.Ctx, big.NewRat(0, 1))

			bankMock.On(
				"GetBalance",
				mock.Anything,
				authtypes.NewModuleAddress(ratelimittypes.TDaiPoolAccount),
				constants.TDai.Denom,
			).Return(sdk.NewCoin(constants.TDai.Denom, sdkmath.NewIntFromBigInt(new(big.Int).SetUint64(1_000_000_000_000))))

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)
			require.NoError(
				t,
				ks.PricesKeeper.UpdateSpotAndPnlMarketPrices(ks.Ctx, &pricestypes.MarketPriceUpdate{
					MarketId:  uint32(0),
					SpotPrice: 4_999_999_937, // Set the price to some large prime number.
					PnlPrice:  4_999_999_937, // Set the price to some large prime number.
				}),
			)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateTDaiAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			testPerps := []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			}
			for _, p := range testPerps {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
					p.Params.DangerIndexPpm,
					p.Params.IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
					p.YieldIndex,
				)
				require.NoError(t, err)
			}

			perptest.SetUpDefaultPerpOIsForTest(
				t,
				ks.Ctx,
				ks.PerpetualsKeeper,
				testPerps,
			)

			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.liquidatedSubaccount)
			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.offsettingSubaccount)
			bankruptcyPriceQuoteQuantums, err := ks.ClobKeeper.GetBankruptcyPriceInQuoteQuantums(
				ks.Ctx,
				*tc.liquidatedSubaccount.GetId(),
				uint32(0),
				tc.deltaQuantums,
			)
			require.NoError(t, err)

			if tc.expectedErr == nil {
				mockIndexerEventManager.On("AddTxnEvent",
					ks.Ctx,
					indexerevents.SubtypeDeleveraging,
					indexerevents.DeleveragingEventVersion,
					indexer_manager.GetBytes(
						indexerevents.NewDeleveragingEvent(
							*tc.liquidatedSubaccount.GetId(),
							*tc.offsettingSubaccount.GetId(),
							uint32(0),
							satypes.BaseQuantums(new(big.Int).Abs(tc.deltaQuantums).Uint64()),
							satypes.BaseQuantums(bankruptcyPriceQuoteQuantums.Uint64()),
							tc.deltaQuantums.Sign() > 0,
							false,
						),
					),
				).Return()
			}
			err = ks.ClobKeeper.ProcessDeleveraging(
				ks.Ctx,
				*tc.liquidatedSubaccount.GetId(),
				*tc.offsettingSubaccount.GetId(),
				uint32(0),
				tc.deltaQuantums,
				bankruptcyPriceQuoteQuantums,
			)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedErr.Error())
			}
		})
	}
}
