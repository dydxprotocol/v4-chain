package keeper_test

import (
	"errors"
	"math"
	"math/big"
	"testing"
	"time"

	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetInsuranceFundBalance(t *testing.T) {
	tests := map[string]struct {
		// Setup
		assets               []assettypes.Asset
		insuranceFundBalance *big.Int

		// Expectations.
		expectedInsuranceFundBalance *big.Int
		expectedError                error
	}{
		"can get zero balance": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance:         new(big.Int),
			expectedInsuranceFundBalance: big.NewInt(0),
		},
		"can get positive balance": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance:         big.NewInt(100),
			expectedInsuranceFundBalance: big.NewInt(100),
		},
		"can get greater than MaxUint64 balance": {
			assets: []assettypes.Asset{
				*constants.Usdc,
			},
			insuranceFundBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				new(big.Int).SetUint64(math.MaxUint64),
			),
			expectedInsuranceFundBalance: new(big.Int).Add(
				new(big.Int).SetUint64(math.MaxUint64),
				new(big.Int).SetUint64(math.MaxUint64),
			),
		},
		"panics when asset not found in state": {
			assets:        []assettypes.Asset{},
			expectedError: errors.New("GetInsuranceFundBalance: Usdc asset not found in state"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

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
				)
				require.NoError(t, err)
			}

			if tc.insuranceFundBalance != nil {
				bankMock.On(
					"GetBalance",
					mock.Anything,
					types.InsuranceFundModuleAddress,
					constants.Usdc.Denom,
				).Return(
					sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(tc.insuranceFundBalance)),
				)
			}

			if tc.expectedError != nil {
				require.PanicsWithValue(
					t,
					tc.expectedError.Error(),
					func() {
						ks.ClobKeeper.GetInsuranceFundBalance(ks.Ctx)
					},
				)
			} else {
				require.Equal(
					t,
					tc.expectedInsuranceFundBalance,
					ks.ClobKeeper.GetInsuranceFundBalance(ks.Ctx),
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

			err := keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			bankMock.On(
				"GetBalance",
				mock.Anything,
				types.InsuranceFundModuleAddress,
				constants.Usdc.Denom,
			).Return(
				sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(tc.insuranceFundBalance)),
			)
			require.Equal(
				t,
				tc.expectedIsValidInsuranceFundDelta,
				ks.ClobKeeper.IsValidInsuranceFundDelta(
					ks.Ctx,
					tc.insuranceFundDelta,
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

		// Expectations.
		expectedCanDeleverageSubaccount bool
	}{
		`Cannot deleverage when subaccount has positive TNC`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_001), // $10,000.000001
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_000_000_000, // $50,000 / BTC
			},

			expectedCanDeleverageSubaccount: false,
		},
		`Cannot deleverage when subaccount has zero TNC`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_001), // $10,000.000001
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_499_000_000, // $54,999 / BTC
			},

			expectedCanDeleverageSubaccount: false,
		},
		`Can deleverage when subaccount has negative TNC`: {
			liquidationConfig:    constants.LiquidationsConfig_No_Limit,
			insuranceFundBalance: big.NewInt(10_000_000_000), // $10,000
			subaccount:           constants.Carl_Num0_1BTC_Short_54999USD,
			marketIdToOraclePriceOverride: map[uint32]uint64{
				constants.BtcUsd.MarketId: 5_500_000_000, // $55,000 / BTC
			},

			expectedCanDeleverageSubaccount: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup keeper state.
			memClob := memclob.NewMemClobPriceTimePriority(false)
			bankMock := &mocks.BankKeeper{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, &mocks.IndexerEventManager{})

			err := keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			// Initialize the liquidations config.
			err = ks.ClobKeeper.InitializeLiquidationsConfig(ks.Ctx, tc.liquidationConfig)
			require.NoError(t, err)

			bankMock.On(
				"GetBalance",
				mock.Anything,
				types.InsuranceFundModuleAddress,
				constants.Usdc.Denom,
			).Return(
				sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewIntFromBigInt(tc.insuranceFundBalance)),
			)

			// Create test markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			// Update the prices on the test markets.
			for marketId, oraclePrice := range tc.marketIdToOraclePriceOverride {
				err := ks.PricesKeeper.UpdateMarketPrices(
					ks.Ctx,
					[]*pricestypes.MsgUpdateMarketPrices_MarketPrice{
						{
							MarketId: marketId,
							Price:    oraclePrice,
						},
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
				)
				require.NoError(t, err)
			}

			ks.SubaccountsKeeper.SetSubaccount(ks.Ctx, tc.subaccount)

			canDeleverageSubaccount, err := ks.ClobKeeper.CanDeleverageSubaccount(
				ks.Ctx,
				*tc.subaccount.Id,
			)
			require.NoError(t, err)
			require.Equal(
				t,
				tc.expectedCanDeleverageSubaccount,
				canDeleverageSubaccount,
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
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 1 BTC short is $54,999 and we close both positions at this price.
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(50_000_000_000 + 54_999_000_000),
					),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
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
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(100_000_000_000 - 54_999_000_000),
					),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
		},
		"Can get multiple offsetting subaccounts": {
			subaccounts: []satypes.Subaccount{
				constants.Carl_Num0_1BTC_Short_54999USD,
				{
					Id: &constants.Dave_Num0,
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
			expectedSubaccounts: []satypes.Subaccount{
				{
					Id: &constants.Carl_Num0,
				},
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 0.5 BTC short is $27,499.5 and we close both positions at this price.
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(50_000_000_000 + 27_499_500_000),
					),
				},
				{
					Id: &constants.Dave_Num1,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 0.5 BTC short is $27,499.5 and we close both positions at this price.
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(50_000_000_000 + 27_499_500_000),
					),
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
					Id: &constants.Carl_Num0,
				},
				constants.Carl_Num1_1BTC_Short,
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 1 BTC short is $54,999 and we close both positions at this price.
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(50_000_000_000 + 54_999_000_000),
					),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
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
					Id: &constants.Carl_Num0,
				},
				constants.Dave_Num1_1ETH_Long_50000USD,
				{
					Id: &constants.Dave_Num0,
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 1 BTC short is $54,999 and we close both positions at this price.
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(50_000_000_000 + 54_999_000_000),
					),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
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
					Id: &constants.Carl_Num0,
				},
				constants.Dave_Num0_1BTC_Long_50001USD_Short,
				{
					Id: &constants.Dave_Num1,
					// TNC of liquidated subaccount is $0, which means the bankruptcy price
					// to close 1 BTC short is $50,000 and we close both positions at this price.
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(50_000_000_000 + 50_000_000_000),
					),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num1,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: new(big.Int),
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
						},
					},
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(-3_000_000_000),
					),
				},
				{
					Id: &constants.Dave_Num0,
					AssetPositions: keepertest.CreateUsdcAssetPosition(
						big.NewInt(50_000_000_000 + 50_000_000_000),
					),
				},
			},
			expectedFills: []types.MatchPerpetualDeleveraging_Fill{
				{
					OffsettingSubaccountId: constants.Dave_Num0,
					FillAmount:             100_000_000,
				},
			},
			expectedQuantumsRemaining: big.NewInt(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
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
				)
				require.NoError(t, err)
			}

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
						),
					),
				).Return()
			}

			fills, deltaQuantumsRemaining := ks.ClobKeeper.OffsetSubaccountPerpetualPosition(
				ks.Ctx,
				tc.liquidatedSubaccountId,
				tc.perpetualId,
				tc.deltaQuantums,
			)
			require.Equal(t, tc.expectedFills, fills)
			require.True(t, tc.expectedQuantumsRemaining.Cmp(deltaQuantumsRemaining) == 0)

			for _, subaccount := range tc.expectedSubaccounts {
				require.Equal(t, subaccount, ks.SubaccountsKeeper.GetSubaccount(ks.Ctx, *subaccount.Id))
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
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(50_000_000_000 + 54_999_000_000),
				),
			},
		},
		"Liquidated: under-collateralized, TNC > 0, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(-45_001_000_000 + 54_999_000_000),
				),
			},
		},
		"Liquidated: under-collateralized, TNC > 0, offsetting: under-collateralized, TNC == 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(-50_000_000_000 + 54_999_000_000),
				),
			},
		},
		"Liquidated: under-collateralized, TNC > 0, offsetting: under-collateralized, TNC < 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
				// to close 1 BTC short is $54,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(-50_001_000_000 + 54_999_000_000)),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: well-collateralized": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $0, which means the bankruptcy price
				// to close 1 BTC short is $50,000 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $0, which means the bankruptcy price
				// to close 1 BTC short is $50,000 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(-45_001_000_000 + 50_000_000_000),
				),
			},
		},
		"Liquidated: under-collateralized, TNC == 0, offsetting: under-collateralized, TNC == 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_50000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $0, which means the bankruptcy price
				// to close 1 BTC short is $50,000 and we close both positions at this price.
				// USDC of this suabccount is -$50,000 + $50,000 = $0.
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
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $-1, which means the bankruptcy price
				// to close 1 BTC short is $49,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(50_000_000_000 + 49_999_000_000),
				),
			},
		},
		"Liquidated: under-collateralized, TNC < 0, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_49999USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_45001USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				// TNC of liquidated subaccount is $-1, which means the bankruptcy price
				// to close 1 BTC short is $49,999 and we close both positions at this price.
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(-45_001_000_000 + 49_999_000_000),
				),
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
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(54_999_000_000 - 5_499_900_000),
				),
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(-90_000_000), // -0.9 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					// TNC of liquidated subaccount is $4,999, which means the bankruptcy price
					// to close 0.1 BTC short is $5,499.9 and we close both positions at this price.
					big.NewInt(50_000_000_000 + 5_499_900_000),
				),
				PerpetualPositions: []*satypes.PerpetualPosition{
					{
						PerpetualId:  0,
						Quantums:     dtypes.NewInt(90_000_000), // 0.9 BTC
						FundingIndex: dtypes.ZeroInt(),
					},
				},
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
			},
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
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
					},
				},
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(50_000_000_000 + 50_500_000_000),
				),
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
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			for _, p := range []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			} {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
				)
				require.NoError(t, err)
			}

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
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(100_000_000_000 - 50_000_000_000),
				),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
			},
		},
		"Liquidated: well-collateralized, offsetting: under-collateralized, TNC > 0": {
			liquidatedSubaccount: constants.Dave_Num0_1BTC_Long_50000USD,
			offsettingSubaccount: constants.Carl_Num0_1BTC_Short_54999USD,
			deltaQuantums:        big.NewInt(-100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(54_999_000_000 - 50_000_000_000),
				),
			},
		},
		"Liquidated: well-collateralized, offsetting: under-collateralized, TNC == 0": {
			liquidatedSubaccount: constants.Carl_Num0_1BTC_Short_100000USD,
			offsettingSubaccount: constants.Dave_Num0_1BTC_Long_50000USD_Short,
			deltaQuantums:        big.NewInt(100_000_000), // 1 BTC

			expectedLiquidatedSubaccount: satypes.Subaccount{
				Id: &constants.Carl_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(100_000_000_000 - 50_000_000_000),
				),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
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
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(54_999_000_000 - 50_000_000_000),
				),
			},
			expectedOffsettingSubaccount: satypes.Subaccount{
				Id: &constants.Dave_Num0,
				AssetPositions: keepertest.CreateUsdcAssetPosition(
					big.NewInt(50_000_000_000 + 50_000_000_000),
				),
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
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			for _, p := range []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			} {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
				)
				require.NoError(t, err)
			}

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
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)
			require.NoError(
				t,
				ks.PricesKeeper.UpdateMarketPrices(ks.Ctx, []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
					{
						MarketId: uint32(0),
						Price:    4_999_999_937, // Set the price to some large prime number.
					},
				}),
			)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)

			err := keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper)
			require.NoError(t, err)

			for _, p := range []perptypes.Perpetual{
				constants.BtcUsd_20PercentInitial_10PercentMaintenance,
				constants.EthUsd_20PercentInitial_10PercentMaintenance,
			} {
				_, err := ks.PerpetualsKeeper.CreatePerpetual(
					ks.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
				)
				require.NoError(t, err)
			}

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
