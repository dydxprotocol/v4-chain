package keeper_test

import (
	"fmt"
	"math"
	"math/big"
	"sort"
	"testing"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	lttest "github.com/dydxprotocol/v4-chain/protocol/testutil/liquidity_tier"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	pricefeed_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/pricefeed"
	pricestest "github.com/dydxprotocol/v4-chain/protocol/testutil/prices"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestModifyPerpetual_Success(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	numLiquidityTiers := 4
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 100)
	numMarkets := keepertest.GetNumMarkets(t, pc.Ctx, pc.PricesKeeper)
	expectedIndexerEvents := make([]*indexerevents.UpdatePerpetualEventV3, len(perps))
	for i, item := range perps {
		// Modify each field arbitrarily and
		// verify the fields were modified in state.
		ticker := fmt.Sprintf("foo_%v", i)
		marketId := uint32(i*2) % numMarkets
		defaultFundingPpm := int32(i * 2)
		liquidityTier := uint32((i + 1) % numLiquidityTiers)
		retItem, err := pc.PerpetualsKeeper.ModifyPerpetual(
			pc.Ctx,
			item.Params.Id,
			ticker,
			marketId,
			defaultFundingPpm,
			liquidityTier,
		)
		require.NoError(t, err)

		// Record the indexer event expected to emit from above `ModifyPerpetual`.
		expectedIndexerEvents[i] = &indexerevents.UpdatePerpetualEventV3{
			Id:                   item.Params.Id,
			Ticker:               ticker,
			MarketId:             marketId,
			AtomicResolution:     item.Params.AtomicResolution,
			LiquidityTier:        liquidityTier,
			MarketType:           v1.ConvertToPerpetualMarketType(item.Params.MarketType),
			DefaultFunding8HrPpm: defaultFundingPpm,
		}

		// Verify updatedp perpetual in store.
		newItem, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, item.Params.Id)
		require.NoError(t, err)
		require.Equal(
			t,
			retItem,
			newItem,
		)
		require.Equal(
			t,
			ticker,
			newItem.Params.Ticker,
		)
		require.Equal(
			t,
			marketId,
			newItem.Params.MarketId,
		)
		require.Equal(
			t,
			int32(i),
			newItem.Params.AtomicResolution,
		)
		require.Equal(
			t,
			defaultFundingPpm,
			newItem.Params.DefaultFundingPpm,
		)
		require.Equal(
			t,
			liquidityTier,
			newItem.Params.LiquidityTier,
		)
	}

	// Verify that expected indexer events were emitted.
	emittedIndexerEvents := getUpdatePerpetualEventsFromIndexerBlock(pc.Ctx, pc.PerpetualsKeeper)
	require.Equal(t, emittedIndexerEvents, expectedIndexerEvents)
}

// getUpdatePerpetualEventsFromIndexerBlock returns all UpdatePerpetual events from the indexer block.
func getUpdatePerpetualEventsFromIndexerBlock(
	ctx sdk.Context,
	perpetualsKeeper *keeper.Keeper,
) []*indexerevents.UpdatePerpetualEventV3 {
	block := perpetualsKeeper.GetIndexerEventManager().ProduceBlock(ctx)
	var updatePerpetualEvents []*indexerevents.UpdatePerpetualEventV3
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeUpdatePerpetual {
			continue
		}
		if _, ok := event.OrderingWithinBlock.(*indexer_manager.IndexerTendermintEvent_TransactionIndex); ok {
			var updatePerpetualEvent indexerevents.UpdatePerpetualEventV3
			err := proto.Unmarshal(event.DataBytes, &updatePerpetualEvent)
			if err != nil {
				panic(err)
			}
			updatePerpetualEvents = append(updatePerpetualEvents, &updatePerpetualEvent)
		}
	}
	return updatePerpetualEvents
}

func TestCreatePerpetual_Failure(t *testing.T) {
	tests := map[string]struct {
		id                uint32
		ticker            string
		marketId          uint32
		atomicResolution  int32
		defaultFundingPpm int32
		liquidityTier     uint32
		marketType        types.PerpetualMarketType
		expectedError     error
	}{
		"Price doesn't exist": {
			id:                0,
			ticker:            "ticker",
			marketId:          999,
			atomicResolution:  -10,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			marketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			expectedError:     errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, fmt.Sprint(999)),
		},
		"Positive default funding magnitude exceeds maximum": {
			id:                0,
			ticker:            "ticker",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: int32(lib.OneMillion + 1),
			liquidityTier:     0,
			marketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			expectedError: errorsmod.Wrap(
				types.ErrDefaultFundingPpmMagnitudeExceedsMax,
				fmt.Sprint(int32(lib.OneMillion+1)),
			),
		},
		"Negative default funding magnitude exceeds maximum": {
			id:                0,
			ticker:            "ticker",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: 0 - int32(lib.OneMillion) - 1,
			liquidityTier:     0,
			marketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			expectedError: errorsmod.Wrap(
				types.ErrDefaultFundingPpmMagnitudeExceedsMax,
				fmt.Sprint(0-int32(lib.OneMillion)-1),
			),
		},
		"Negative default funding magnitude exceeds maximum due to overflow": {
			id:                0,
			ticker:            "ticker",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: math.MinInt32,
			liquidityTier:     0,
			marketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			expectedError:     errorsmod.Wrap(types.ErrDefaultFundingPpmMagnitudeExceedsMax, fmt.Sprint(math.MinInt32)),
		},
		"Ticker is an empty string": {
			id:                0,
			ticker:            "",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			marketType:        types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			expectedError:     types.ErrTickerEmptyString,
		},
		"Unspecified market type": {
			id:                0,
			ticker:            "",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			expectedError: errorsmod.Wrap(
				types.ErrInvalidMarketType,
				fmt.Sprintf("market type %v", types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED),
			),
		},
		"Invalid market type": {
			id:                0,
			ticker:            "",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			marketType:        3,
			expectedError: errorsmod.Wrap(
				types.ErrInvalidMarketType,
				fmt.Sprintf("market type %v", 3),
			),
		},
	}

	// Test setup.
	pc := keepertest.PerpetualsKeepers(t)
	keepertest.CreateNMarkets(t, pc.Ctx, pc.PricesKeeper, 1)
	// Create Liquidity Tiers
	keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper)
	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := pc.PerpetualsKeeper.CreatePerpetual(
				pc.Ctx,
				tc.id,
				tc.ticker,
				tc.marketId,
				tc.atomicResolution,
				tc.defaultFundingPpm,
				tc.liquidityTier,
				tc.marketType,
			)

			require.Error(t, err)
			require.EqualError(t, err, tc.expectedError.Error())
		})
	}
}

func TestModifyPerpetual_Failure(t *testing.T) {
	tests := map[string]struct {
		id                uint32
		ticker            string
		marketId          uint32
		defaultFundingPpm int32
		liquidityTier     uint32
		expectedError     error
	}{
		"Perpetual doesn't exist": {
			id:                999,
			ticker:            "ticker",
			marketId:          0,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			expectedError:     errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(999)),
		},
		"Price doesn't exist": {
			id:                0,
			ticker:            "ticker",
			marketId:          999,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			expectedError:     errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, fmt.Sprint(999)),
		},
		"Ticker is an empty string": {
			id:                0,
			ticker:            "",
			marketId:          0,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			expectedError:     types.ErrTickerEmptyString,
		},
		"Modified to empty liquidity tier": {
			id:                0,
			ticker:            "ticker",
			marketId:          0,
			defaultFundingPpm: 0,
			liquidityTier:     999,
			expectedError:     errorsmod.Wrap(types.ErrLiquidityTierDoesNotExist, fmt.Sprint(999)),
		},
	}

	// Test setup.
	pc := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := pc.PerpetualsKeeper.ModifyPerpetual(
				pc.Ctx,
				tc.id,
				tc.ticker,
				tc.marketId,
				tc.defaultFundingPpm,
				tc.liquidityTier,
			)

			require.Error(t, err)
			require.EqualError(t, err, tc.expectedError.Error())
		})
	}
}

func TestSetPerpetualMarketType(t *testing.T) {
	tests := map[string]struct {
		currType      types.PerpetualMarketType
		newType       types.PerpetualMarketType
		errorExpected bool
		expectedError error
	}{
		"success - set unspecified to cross": {
			currType:      types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED,
			newType:       types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			errorExpected: false,
		},
		"success - set isolated to cross": {
			currType:      types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
			newType:       types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			errorExpected: false,
		},
		"failure - setting to unspecified": {
			currType:      types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			newType:       types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED,
			errorExpected: true,
			expectedError: errorsmod.Wrap(
				types.ErrInvalidMarketType,
				fmt.Sprintf(
					"invalid market type %v for perpetual %d",
					types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED, 0,
				),
			),
		},
		"failure - market type already set to cross": {
			currType:      types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
			newType:       types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
			errorExpected: true,
			expectedError: errorsmod.Wrap(
				types.ErrInvalidMarketType,
				fmt.Sprintf(
					"perpetual %d already has market type %v and cannot be changed",
					0,
					types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
				),
			),
		},
	}

	// Test setup.
	for name, tc := range tests {
		t.Run(
			name, func(t *testing.T) {
				pc := keepertest.PerpetualsKeepers(t)
				// Create liquidity tiers and perpetuals,
				perp := keepertest.CreateLiquidityTiersAndNPerpetuals(
					t,
					pc.Ctx,
					pc.PerpetualsKeeper,
					pc.PricesKeeper,
					1,
				)[0]
				perp.Params.MarketType = tc.currType
				pc.PerpetualsKeeper.SetPerpetualForTest(pc.Ctx, perp)

				_, err := pc.PerpetualsKeeper.SetPerpetualMarketType(
					pc.Ctx,
					perp.Params.Id,
					tc.newType,
				)

				if tc.errorExpected {
					require.EqualError(t, err, tc.expectedError.Error())
				} else {
					require.NoError(t, err)

					rst, err := pc.PerpetualsKeeper.GetPerpetual(
						pc.Ctx,
						perp.Params.Id,
					)
					require.NoError(t, err)
					require.Equal(t, tc.newType, rst.Params.MarketType)
				}
			},
		)
	}
}

func TestGetPerpetual_Success(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 10)

	for _, perp := range perps {
		rst, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx,
			perp.Params.Id,
		)
		require.NoError(t, err)
		require.Equal(t,
			nullify.Fill(&perp), //nolint:staticcheck
			nullify.Fill(&rst),  //nolint:staticcheck
		)
	}
}

func TestHasPerpetual(t *testing.T) {
	// Setup context and keepers
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals
	keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper)
	perps := []types.Perpetual{
		*perptest.GeneratePerpetual(perptest.WithId(0)),
		*perptest.GeneratePerpetual(perptest.WithId(5)),
		*perptest.GeneratePerpetual(perptest.WithId(20)),
		*perptest.GeneratePerpetual(perptest.WithId(999)),
	}

	_, err := keepertest.CreateTestMarket(
		t,
		pc.Ctx,
		pc.PricesKeeper,
		pricestypes.MarketParam{
			Id:                 0,
			Pair:               "base-quote",
			Exponent:           -10,
			MinExchanges:       uint32(1),
			MinPriceChangePpm:  uint32(50),
			ExchangeConfigJson: "{}",
		},
		pricestypes.MarketPrice{
			Id:       0,
			Exponent: -10,
			Price:    1_000, // leave this as a placeholder b/c we cannot set the price to 0
		},
	)
	require.NoError(t, err)

	for perp := range perps {
		_, err := pc.PerpetualsKeeper.CreatePerpetual(
			pc.Ctx,
			perps[perp].Params.Id,
			perps[perp].Params.Ticker,
			perps[perp].Params.MarketId,
			perps[perp].Params.AtomicResolution,
			perps[perp].Params.DefaultFundingPpm,
			perps[perp].Params.LiquidityTier,
			perps[perp].Params.MarketType,
		)
		require.NoError(t, err)
	}

	for _, perp := range perps {
		// Test if HasPerpetual correctly identifies an existing perpetual
		found := pc.PerpetualsKeeper.HasPerpetual(pc.Ctx, perp.Params.Id)
		require.True(t, found, "Expected to find perpetual with id %d, but it was not found", perp.Params.Id)
	}

	found := pc.PerpetualsKeeper.HasPerpetual(pc.Ctx, 9999)
	require.False(t, found, "Expected not to find perpetual with id 9999, but it was found")
}

func TestGetPerpetual_NotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := pc.PerpetualsKeeper.GetPerpetual(
		pc.Ctx,
		nonExistentPerpetualId,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetPerpetuals_Success(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 10)

	require.ElementsMatch(t,
		nullify.Fill(perps), //nolint:staticcheck
		nullify.Fill(pc.PerpetualsKeeper.GetAllPerpetuals(pc.Ctx)), //nolint:staticcheck
	)
}

func TestGetAllPerpetuals_Sorted(t *testing.T) {
	// Setup context and keepers
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals
	keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper)
	perps := []types.Perpetual{
		*perptest.GeneratePerpetual(perptest.WithId(999)),
		*perptest.GeneratePerpetual(perptest.WithId(5)),
		*perptest.GeneratePerpetual(perptest.WithId(0)),
		*perptest.GeneratePerpetual(perptest.WithId(20)),
		*perptest.GeneratePerpetual(perptest.WithId(1)),
	}

	_, err := keepertest.CreateTestMarket(
		t,
		pc.Ctx,
		pc.PricesKeeper,
		pricestypes.MarketParam{
			Id:                 0,
			Pair:               "base-quote",
			Exponent:           -10,
			MinExchanges:       uint32(1),
			MinPriceChangePpm:  uint32(50),
			ExchangeConfigJson: "{}",
		},
		pricestypes.MarketPrice{
			Id:       0,
			Exponent: -10,
			Price:    1_000, // leave this as a placeholder b/c we cannot set the price to 0
		},
	)
	require.NoError(t, err)

	for perp := range perps {
		_, err := pc.PerpetualsKeeper.CreatePerpetual(
			pc.Ctx,
			perps[perp].Params.Id,
			perps[perp].Params.Ticker,
			perps[perp].Params.MarketId,
			perps[perp].Params.AtomicResolution,
			perps[perp].Params.DefaultFundingPpm,
			perps[perp].Params.LiquidityTier,
			perps[perp].Params.MarketType,
		)
		require.NoError(t, err)
	}

	got := pc.PerpetualsKeeper.GetAllPerpetuals(pc.Ctx)
	require.Equal(
		t,
		[]types.Perpetual{
			*perptest.GeneratePerpetual(perptest.WithId(0)),
			*perptest.GeneratePerpetual(perptest.WithId(1)),
			*perptest.GeneratePerpetual(perptest.WithId(5)),
			*perptest.GeneratePerpetual(perptest.WithId(20)),
			*perptest.GeneratePerpetual(perptest.WithId(999)),
		},
		got,
	)
}

func TestModifyOpenInterest_Failure(t *testing.T) {
	testCases := map[string]struct {
		id                uint32
		initOpenInterest  *big.Int
		openInterestDelta *big.Int
		err               error
	}{
		"Would become negative": {
			id:                0,
			initOpenInterest:  big.NewInt(1_000),
			openInterestDelta: big.NewInt(-1_001),
			err:               types.ErrOpenInterestWouldBecomeNegative,
		},
		"Non-existent perp Id": {
			id:                1111,
			initOpenInterest:  big.NewInt(1_000),
			openInterestDelta: big.NewInt(100),
			err:               types.ErrPerpetualDoesNotExist,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			pc := keepertest.PerpetualsKeepers(t)
			perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)

			// Set up initial open interest
			require.NoError(t, pc.PerpetualsKeeper.ModifyOpenInterest(
				pc.Ctx,
				perps[0].Params.Id,
				tc.initOpenInterest,
			))

			err := pc.PerpetualsKeeper.ModifyOpenInterest(
				pc.Ctx,
				tc.id,
				tc.openInterestDelta,
			)
			require.ErrorContains(t, err, tc.err.Error())
		})
	}
}

func TestModifyOpenInterest_Mixed(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 100)

	for _, perp := range perps {
		openInterestDeltaBaseQuantums := big.NewInt(2_000_000_000*(int64(perp.Params.Id)%2) - 1)

		// Add `openInterestDeltaBaseQuantums` to open interest which is initially 0.
		err := pc.PerpetualsKeeper.ModifyOpenInterest(
			pc.Ctx,
			perp.Params.Id,
			openInterestDeltaBaseQuantums,
		)

		// If Id is even, the modification is negagive and should fail.
		if perp.Params.Id%2 == 0 {
			require.ErrorContains(t,
				err,
				types.ErrOpenInterestWouldBecomeNegative.Error(),
			)

			newPerp, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, perp.Params.Id)
			require.NoError(t, err)

			require.Equal(
				t,
				big.NewInt(0), // open interest should remain 0
				newPerp.OpenInterest.BigInt(),
			)
		} else {
			require.NoError(t, err)

			newPerp, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, perp.Params.Id)
			require.NoError(t, err)

			require.Equal(
				t,
				openInterestDeltaBaseQuantums,
				newPerp.OpenInterest.BigInt(),
			)
		}

		// Add `openInterestDeltaBaseQuantums` again
		err = pc.PerpetualsKeeper.ModifyOpenInterest(
			pc.Ctx,
			perp.Params.Id,
			openInterestDeltaBaseQuantums,
		)
		// If Id is even, the modification is negagive and should fail.
		if perp.Params.Id%2 == 0 {
			require.ErrorContains(t,
				err,
				types.ErrOpenInterestWouldBecomeNegative.Error(),
			)

			newPerp, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, perp.Params.Id)
			require.NoError(t, err)

			require.Equal(
				t,
				big.NewInt(0), // open interest should remain 0
				newPerp.OpenInterest.BigInt(),
			)
		} else {
			require.NoError(t, err)

			newPerp, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, perp.Params.Id)
			require.NoError(t, err)

			require.Equal(
				t,
				// open interest should be 2 * delta now
				openInterestDeltaBaseQuantums.Mul(
					openInterestDeltaBaseQuantums,
					big.NewInt(2),
				),
				newPerp.OpenInterest.BigInt(),
			)
		}
	}
}

func TestGetPerpetualAndMarketPriceAndLiquidityTier_PerpetualNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, _, _, err := pc.PerpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(
		pc.Ctx,
		nonExistentPerpetualId,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetPerpetualAndMarketPriceAndLiquidityTier_MarketNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(pc.Ctx.KVStore(pc.StoreKey), []byte(types.PerpetualKeyPrefix))
	perpetualStore.Set(lib.Uint32ToKey(perpetual.Params.Id), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, _, _, err := pc.PerpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(
		pc.Ctx,
		perpetual.Params.Id,
	)

	require.EqualError(
		t,
		err,
		errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, fmt.Sprint(nonExistentMarketId)).Error(),
	)
	require.ErrorIs(t, err, pricestypes.ErrMarketPriceDoesNotExist)
}

func TestGetPerpetualAndMarketPriceAndLiquidityTier_LiquidityTierNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad LiquidityTier.
	nonExistentLiquidityTier := uint32(999)
	perpetual.Params.LiquidityTier = nonExistentLiquidityTier
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(pc.Ctx.KVStore(pc.StoreKey), []byte(types.PerpetualKeyPrefix))
	perpetualStore.Set(lib.Uint32ToKey(perpetual.Params.Id), b)

	// Getting margin requirements for perpetual with bad LiquidityTier should return an error.
	_, _, _, err := pc.PerpetualsKeeper.GetPerpetualAndMarketPriceAndLiquidityTier(
		pc.Ctx,
		perpetual.Params.Id,
	)

	require.EqualError(
		t,
		err,
		errorsmod.Wrap(types.ErrLiquidityTierDoesNotExist, fmt.Sprint(nonExistentLiquidityTier)).Error(),
	)
	require.ErrorIs(t, err, types.ErrLiquidityTierDoesNotExist)
}

func TestGetNetNotional_Success(t *testing.T) {
	tests := map[string]struct {
		price                               uint64
		exponent                            int32
		baseCurrencyAtomicResolution        int32
		bigBaseQuantums                     *big.Int
		bigExpectedNetNotionalQuoteQuantums *big.Int
	}{
		// TODO: Add back tests for positive and zero exponent once x/marketmap supports them
		"Negative exponent, atomic resolution 6, long position": {
			price:                               5_555,
			exponent:                            -2,
			baseCurrencyAtomicResolution:        -6,
			bigBaseQuantums:                     big.NewInt(7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(388_850),
		},
		"Negative exponent, atomic resolution 6, short position": {
			price:                               5_555,
			exponent:                            -2,
			baseCurrencyAtomicResolution:        -6,
			bigBaseQuantums:                     big.NewInt(-7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(-388_850),
		},
		"Negative exponent, atomic resolution 4, long position": {
			price:                               5_555,
			exponent:                            -2,
			baseCurrencyAtomicResolution:        -4,
			bigBaseQuantums:                     big.NewInt(7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(38_885_000),
		},
		"Negative exponent, atomic resolution 0, long position": {
			price:                               5_555,
			exponent:                            -2,
			baseCurrencyAtomicResolution:        -0,
			bigBaseQuantums:                     big.NewInt(7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(38_8850_000_000),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -4,
			bigBaseQuantums:              new(big.Int).SetUint64(math.MaxUint64),
			bigExpectedNetNotionalQuoteQuantums: big_testutil.MustFirst(
				new(big.Int).SetString("3402823669209384634264811192843491082250", 10),
			),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Test suite setup.
			pc := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper)
			// Create a new market param and price.
			marketId := uint32(0)
			_, err := keepertest.CreateTestMarket(
				t,
				pc.Ctx,
				pc.PricesKeeper,
				pricestypes.MarketParam{
					Id:                 marketId,
					Pair:               "base-quote",
					Exponent:           tc.exponent,
					MinExchanges:       uint32(1),
					MinPriceChangePpm:  uint32(50),
					ExchangeConfigJson: "{}",
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    tc.price,
				},
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := pc.PerpetualsKeeper.CreatePerpetual(
				pc.Ctx,
				0,                               // PerpetualId
				"GetNetNotionalTicker",          // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
				types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS, // MarketType
			)
			require.NoError(t, err)

			// Verify collateral requirements are calculated correctly.
			bigNotionalQuoteQuantums, err := pc.PerpetualsKeeper.GetNetNotional(
				pc.Ctx,
				perpetual.Params.Id,
				tc.bigBaseQuantums,
			)
			require.NoError(t, err)

			require.Equal(t, tc.bigExpectedNetNotionalQuoteQuantums, bigNotionalQuoteQuantums, "Net notional mismatch")
		})
	}
}

func TestGetNetNotional_PerpetualNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := pc.PerpetualsKeeper.GetNetNotional(
		pc.Ctx,
		nonExistentPerpetualId,
		big.NewInt(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetNetNotional_MarketNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(pc.Ctx.KVStore(pc.StoreKey), []byte(types.PerpetualKeyPrefix))
	perpetualStore.Set(lib.Uint32ToKey(perpetual.Params.Id), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, err := pc.PerpetualsKeeper.GetNetNotional(
		pc.Ctx,
		perpetual.Params.Id,
		big.NewInt(-1),
	)
	expectedErrorStr := fmt.Sprintf(
		"Market ID %d does not exist on perpetual ID %d",
		perpetual.Params.MarketId,
		perpetual.Params.Id,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrMarketDoesNotExist, expectedErrorStr).Error())
	require.ErrorIs(t, err, types.ErrMarketDoesNotExist)
}

func TestGetNotionalInBaseQuantums_Success(t *testing.T) {
	tests := map[string]struct {
		price                              uint64
		exponent                           int32
		baseCurrencyAtomicResolution       int32
		bigQuoteQuantums                   *big.Int
		bigExpectedNetNotionalBaseQuantums *big.Int
	}{
		// TODO: Add back tests for positive and zero exponent once x/marketmap supports them
		"Negative exponent, atomic resolution 6, long position": {
			price:                              5_555,
			exponent:                           -2,
			baseCurrencyAtomicResolution:       -6,
			bigQuoteQuantums:                   big.NewInt(388_850),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(7_000),
		},
		"Negative exponent, atomic resolution 6, short position": {
			price:                              5_555,
			exponent:                           -2,
			baseCurrencyAtomicResolution:       -6,
			bigQuoteQuantums:                   big.NewInt(-388_850),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(-7_000),
		},
		"Negative exponent, atomic resolution 4, long position": {
			price:                              5_555,
			exponent:                           -2,
			baseCurrencyAtomicResolution:       -4,
			bigQuoteQuantums:                   big.NewInt(38_885_000),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(7_000),
		},
		"Negative exponent, atomic resolution 0, long position": {
			price:                              5_555,
			exponent:                           -2,
			baseCurrencyAtomicResolution:       -0,
			bigQuoteQuantums:                   big.NewInt(388_850_000_000),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(7_000),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -4,
			bigQuoteQuantums: big_testutil.MustFirst(
				new(big.Int).SetString("3402823669209384634264811192843491082250", 10),
			),
			bigExpectedNetNotionalBaseQuantums: new(big.Int).SetUint64(math.MaxUint64),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Test suite setup.
			pc := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper) // Create a new market param and price.
			marketId := keepertest.GetNumMarkets(t, pc.Ctx, pc.PricesKeeper)
			_, err := keepertest.CreateTestMarket(
				t,
				pc.Ctx,
				pc.PricesKeeper,
				pricestypes.MarketParam{
					Id:                 marketId,
					Pair:               "base-quote",
					Exponent:           tc.exponent,
					MinExchanges:       uint32(1),
					MinPriceChangePpm:  uint32(50),
					ExchangeConfigJson: "{}",
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    tc.price,
				},
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := pc.PerpetualsKeeper.CreatePerpetual(
				pc.Ctx,
				0,                               // PerpetualId
				"GetNetNotionalTicker",          // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
				types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS, // MarketType
			)
			require.NoError(t, err)

			// Verify collateral requirements are calculated correctly.
			bigNotionalBaseQuantums, err := pc.PerpetualsKeeper.GetNotionalInBaseQuantums(
				pc.Ctx,
				perpetual.Params.Id,
				tc.bigQuoteQuantums,
			)
			require.NoError(t, err)

			require.Equal(t, tc.bigExpectedNetNotionalBaseQuantums, bigNotionalBaseQuantums, "Net notional mismatch")
		})
	}
}

func TestGetNotionalInBaseQuantums_PerpetualNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := pc.PerpetualsKeeper.GetNotionalInBaseQuantums(
		pc.Ctx,
		nonExistentPerpetualId,
		big.NewInt(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetNotionalInBaseQuantums_MarketNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(pc.Ctx.KVStore(pc.StoreKey), []byte(types.PerpetualKeyPrefix))
	perpetualStore.Set(lib.Uint32ToKey(perpetual.Params.Id), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, err := pc.PerpetualsKeeper.GetNotionalInBaseQuantums(
		pc.Ctx,
		perpetual.Params.Id,
		big.NewInt(-1),
	)
	expectedErrorStr := fmt.Sprintf(
		"Market ID %d does not exist on perpetual ID %d",
		perpetual.Params.MarketId,
		perpetual.Params.Id,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrMarketDoesNotExist, expectedErrorStr).Error())
	require.ErrorIs(t, err, types.ErrMarketDoesNotExist)
}

func TestGetNetCollateral_Success(t *testing.T) {
	tests := map[string]struct {
		price                                 uint64
		exponent                              int32
		baseCurrencyAtomicResolution          int32
		bigBaseQuantums                       *big.Int
		bigExpectedNetCollateralQuoteQuantums *big.Int
	}{
		// TODO: Add back tests for positive and zero exponent once x/marketmap supports them
		"Negative exponent, atomic resolution 6, short position": {
			price:                                 5_555,
			exponent:                              -2,
			baseCurrencyAtomicResolution:          -6,
			bigBaseQuantums:                       big.NewInt(-7_000),
			bigExpectedNetCollateralQuoteQuantums: big.NewInt(-388_850),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     -1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              new(big.Int).SetUint64(math.MaxUint64),
			bigExpectedNetCollateralQuoteQuantums: big_testutil.MustFirst(
				new(big.Int).SetString("34028236692093846342648111928434910822", 10),
			),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Test suite setup.
			pc := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper) // Test setup.
			// Create a new market.
			marketId := keepertest.GetNumMarkets(t, pc.Ctx, pc.PricesKeeper)
			_, err := keepertest.CreateTestMarket(
				t,
				pc.Ctx,
				pc.PricesKeeper,
				pricestypes.MarketParam{
					Id:                 marketId,
					Pair:               "base-quote",
					Exponent:           tc.exponent,
					MinExchanges:       uint32(1),
					MinPriceChangePpm:  uint32(50),
					ExchangeConfigJson: "{}",
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    tc.price,
				},
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := pc.PerpetualsKeeper.CreatePerpetual(
				pc.Ctx,
				0,                               // PerpetualId
				"GetNetCollateralTicker",        // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
				types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS, // MarketType
			)
			require.NoError(t, err)

			// Verify collateral requirements are calculated correctly.
			bigCollateralQuoteQuantums, err := pc.PerpetualsKeeper.GetNetCollateral(
				pc.Ctx,
				perpetual.Params.Id,
				tc.bigBaseQuantums,
			)
			require.NoError(t, err)

			require.Equal(t, tc.bigExpectedNetCollateralQuoteQuantums, bigCollateralQuoteQuantums, "Net collateral mismatch")
		})
	}
}

func TestGetNetCollateral_PerpetualNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := pc.PerpetualsKeeper.GetNetCollateral(
		pc.Ctx,
		nonExistentPerpetualId,
		big.NewInt(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetNetCollateral_MarketNotFound(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	cdc := codec.NewProtoCodec(module.InterfaceRegistry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(pc.Ctx.KVStore(pc.StoreKey), []byte(types.PerpetualKeyPrefix))
	perpetualStore.Set(lib.Uint32ToKey(perpetual.Params.Id), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, err := pc.PerpetualsKeeper.GetNetCollateral(
		pc.Ctx,
		perpetual.Params.Id,
		big.NewInt(-1),
	)
	expectedErrorStr := fmt.Sprintf(
		"Market ID %d does not exist on perpetual ID %d",
		perpetual.Params.MarketId,
		perpetual.Params.Id,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrMarketDoesNotExist, expectedErrorStr).Error())
	require.ErrorIs(t, err, types.ErrMarketDoesNotExist)
}

func TestModifyFundingIndex_Success(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 100)

	for _, perp := range perps {
		testFundingIndexDelta := big.NewInt(2*(int64(perp.Params.Id)%2) - 1)

		err := pc.PerpetualsKeeper.ModifyFundingIndex(
			pc.Ctx,
			perp.Params.Id,
			testFundingIndexDelta,
		)
		require.NoError(t, err)

		newPerp, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, perp.Params.Id)
		require.NoError(t, err)

		require.Equal(
			t,
			testFundingIndexDelta,
			newPerp.FundingIndex.BigInt(),
		)
	}
}

func TestModifyFundingIndex_PerpetualDoesNotExist(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)

	err := pc.PerpetualsKeeper.ModifyFundingIndex(
		pc.Ctx,
		nonExistentPerpetualId,
		big.NewInt(1),
	)

	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestModifyFundingIndex_IntegerOverflowUnderflow(t *testing.T) {
	tests := map[string]struct {
		perpetualId         uint32
		fundingIndexDelta   *big.Int
		initialFundingIndex *big.Int
	}{
		"funding index overflow": {
			perpetualId:         0,
			fundingIndexDelta:   big.NewInt(math.MaxInt64),
			initialFundingIndex: big.NewInt(1),
		},
		"funding index underflow": {
			perpetualId:         0,
			fundingIndexDelta:   big.NewInt(math.MinInt64),
			initialFundingIndex: big.NewInt(-1),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pc := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers and perpetuals,
			_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 1)

			// Set up intitial funding index, should succeed.
			err := pc.PerpetualsKeeper.ModifyFundingIndex(
				pc.Ctx,
				tc.perpetualId,
				tc.initialFundingIndex,
			)
			require.NoError(t, err)

			err = pc.PerpetualsKeeper.ModifyFundingIndex(
				pc.Ctx,
				tc.perpetualId,
				tc.fundingIndexDelta,
			)
			require.NoError(t, err)
		})
	}
}

func TestGetRemoveSampleTailsFunc(t *testing.T) {
	tests := map[string]struct {
		removalRatePpm uint32
		input          []int32
		expectedOutput []int32
	}{
		"25%, input length = 0": {
			removalRatePpm: 250_000,
			input:          []int32{},
			expectedOutput: []int32{},
		},
		"25%, input length = 1": {
			removalRatePpm: 250_000,
			input:          []int32{0},
			expectedOutput: []int32{0},
		},
		"25%, input length = 3": {
			removalRatePpm: 250_000,
			input:          []int32{0, -1, -3},
			expectedOutput: []int32{-3, -1}, // bottomRemoval = 0, topRemoval = 1
		},
		"25%, input length = 4": {
			removalRatePpm: 250_000,
			input:          []int32{0, -1, -3, 5},
			expectedOutput: []int32{-1, 0}, // bottomRemoval = 1, topRemoval = 1
		},
		"25%, input length = 5": {
			removalRatePpm: 250_000,
			input:          []int32{0, -1, -3, 5, 7},
			expectedOutput: []int32{-1, 0, 5}, // bottomRemoval = 1, topRemoval = 1
		},
		"25%, input length = 6": {
			removalRatePpm: 250_000,
			input:          []int32{0, -1, -3, -5, 5, 7},
			expectedOutput: []int32{-3, -1, 0}, // bottomRemoval = 1, topRemoval = 2
		},
		"10%, input length = 5": {
			removalRatePpm: 100_000,
			input:          []int32{0, -1, -3, -5, 5},
			expectedOutput: []int32{-5, -3, -1, 0}, // bottomRemoval = 0, topRemoval = 1
		},
		"80%, invalid removal ratio, skips removing samples": {
			removalRatePpm: 800_000,
			input:          []int32{0, -1, -3, -5, 5},
			expectedOutput: []int32{0, -1, -3, -5, 5},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pc := keepertest.PerpetualsKeepers(t)

			sampleTailsRemovalFunc := pc.PerpetualsKeeper.GetRemoveSampleTailsFunc(pc.Ctx, tc.removalRatePpm)
			output := sampleTailsRemovalFunc(tc.input)

			require.Equal(t,
				tc.expectedOutput,
				output,
			)
		})
	}
}

func TestMaybeProcessNewFundingTickEpoch_ProcessNewEpoch(t *testing.T) {
	tests := map[string]struct {
		testFundingSampleDuration        uint32
		testFundingTickDuration          uint32
		testPerpetuals                   []types.Perpetual
		testFundingSamples               []int32
		expectedFundingIndexDeltas       []*big.Int
		expectedFundingIndexDeltaStrings []string
		fundingRatesAndIndices           []indexerevents.FundingUpdateV1
	}{
		"Success: 60 equivalent samples of 0.001 percent, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// Premium sample = 0.001%, length = 60.
			testFundingSamples:               constants.GenerateConstantFundingPremiums(1000, 60),
			expectedFundingIndexDeltaStrings: []string{"625"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 1_000,
					FundingIndex:    dtypes.NewInt(625),
				},
			},
		},
		"Success: 60 equivalent samples of -0.001 percent, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// Premium sample = -0.001%, length = 60.
			testFundingSamples:               constants.GenerateConstantFundingPremiums(-1000, 60),
			expectedFundingIndexDeltaStrings: []string{"-625"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: -1_000,
					FundingIndex:    dtypes.NewInt(-625),
				},
			},
		},
		"Success: 60 equivalent samples of int32 max, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// Premium sample = MaxInt32, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(math.MaxInt32, 60),
			// 8-hr funding rate capped at 1_500_000, prorated funding rate thus is 187_500
			// funding index delta = 187_500 * 5_000_000_000 * 10^(-5 + -10) * 10^6 = 937500
			expectedFundingIndexDeltaStrings: []string{"937500"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 1_500_000,
					FundingIndex:    dtypes.NewInt(937500),
				},
			},
		},
		"Success: 60 equivalent samples of int32 min, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// Premium sample = MinInt32, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(math.MinInt32, 60),
			// 8-hr funding rate capped at -1_500_000, prorated funding rate thus is -187_500
			expectedFundingIndexDeltaStrings: []string{"-937500"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: -1_500_000,
					FundingIndex:    dtypes.NewInt(-937500),
				},
			},
		},
		"Success: 60 equivalent samples of 0.2 percent, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// Premium sample = 0.2%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(2000, 60),
			// 8-hr funding rate is 2_000, prorated funding rate thus is 250
			expectedFundingIndexDeltaStrings: []string{"1250"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 2_000,
					FundingIndex:    dtypes.NewInt(1250),
				},
			},
		},
		"Success: 60 equivalent samples of 15 percent, 60 samples expected, 20% IM, 18% MM": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM,
			},
			// Premium sample = 15%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(150_000, 60),
			// 8-hr funding rate capped at 120_000, prorated funding rate thus is 15_000
			expectedFundingIndexDeltaStrings: []string{"75000"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM.GetId(),
					FundingValuePpm: 120_000,
					FundingIndex:    dtypes.NewInt(75000),
				},
			},
		},
		"Success: 60 equivalent samples of -15 percent, 60 samples expected, 20% IM, 18% MM": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM,
			},
			// Premium sample = -15%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(-150_000, 60),
			// 8-hr funding rate capped at -120_000, prorated funding rate thus is -15_000
			expectedFundingIndexDeltaStrings: []string{"-75000"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM.GetId(),
					FundingValuePpm: -120_000,
					FundingIndex:    dtypes.NewInt(-75000),
				},
			},
		},
		"Success: 60 equivalent samples of 12 percent, 60 samples expected, 20% IM, 18% MM": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM,
			},
			// Premium sample = 12%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(120_000, 60),
			// 8-hr funding rate capped at 120_000 (same value as sample)
			// prorated funding rate for one hour thus is 15_000
			expectedFundingIndexDeltaStrings: []string{"75000"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM.GetId(),
					FundingValuePpm: 120_000,
					FundingIndex:    dtypes.NewInt(75000),
				},
			},
		},
		"Success: 60 equivalent samples of -12 percent, 60 samples expected, 20% IM, 18% MM": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM,
			},
			// Premium sample = -12%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(-120_000, 60),
			// 8-hr funding rate is -120_000, prorated funding rate thus is -15_000
			expectedFundingIndexDeltaStrings: []string{"-75000"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution_20IM_18MM.GetId(),
					FundingValuePpm: -120_000,
					FundingIndex:    dtypes.NewInt(-75000),
				},
			},
		},
		"Success: 60 equivalent samples of 12 percent, 60 samples expected, IM equal to MM": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_100PercentMarginRequirement,
			},
			// Premium sample = 12%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(120_000, 60),
			// funding rate is clamped to 0
			expectedFundingIndexDeltaStrings: []string{"0"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_100PercentMarginRequirement.GetId(),
					FundingValuePpm: 0,
					FundingIndex:    dtypes.NewInt(0),
				},
			},
		},
		"Success: 60 equivalent samples of 0.001 percent, 60 samples expected, two perpetuals": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
				constants.EthUsd_0DefaultFunding_9AtomicResolution,
			},
			// Premium sample = 0.001%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(1000, 60),
			// 8-hr funding rate is 1000, prorated funding rate thus is 125
			expectedFundingIndexDeltaStrings: []string{
				// Btc: 0.001% at $50000 (price_value= 5_000_000_000, base_atomic = -10, price_exponenet = -5)
				"625",
				// Eth: 0.001% at $3000, (price_value= 3_000_000_000, base_atomic = -9, price_exponenet = -6)
				"375",
			},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 1_000,
					FundingIndex:    dtypes.NewInt(625),
				},
				{
					PerpetualId:     constants.EthUsd_0DefaultFunding_9AtomicResolution.GetId(),
					FundingValuePpm: 1_000,
					FundingIndex:    dtypes.NewInt(375),
				},
			},
		},
		"Success: 60 equivalent samples of 0.001 percent, 120 samples expected": {
			testFundingSampleDuration: 30,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// Premium sample = 0.001%, length = 60.
			testFundingSamples: constants.GenerateConstantFundingPremiums(1000, 60),
			// Average sampled rate = 0.0005% due to padding of zeros.
			expectedFundingIndexDeltaStrings: []string{"312"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 500,
					FundingIndex:    dtypes.NewInt(312),
				},
			},
		},
		"Success: 30 equivalent samples of 0.001 percent, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// Premium sample = 0.001%, length = 30. Average sampled rate = 0.0005% due to padding of zeros.
			testFundingSamples:               constants.GenerateConstantFundingPremiums(1000, 30),
			expectedFundingIndexDeltaStrings: []string{"312"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 500,
					FundingIndex:    dtypes.NewInt(312),
				},
			},
		},
		"Success: 60 equivalent samples of 0.001 percent, default funding = 0.001 percent, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0_001Percent_DefaultFunding_10AtomicResolution,
			},
			// Premium sample = 0.001%, length = 60.
			testFundingSamples:               constants.GenerateConstantFundingPremiums(1000, 60),
			expectedFundingIndexDeltaStrings: []string{"1250"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0_001Percent_DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 2_000, // 0.001% (premium) + 0.001% (default funding)
					FundingIndex:    dtypes.NewInt(1250),
				},
			},
		},
		"Success: more than expected funding samples recorded, 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_10AtomicResolution,
			},
			// 25 samples of 0.005% and 75 samples of 0.001%.
			testFundingSamples: constants.GenerateFundingSamplesWithValues(
				[]int32{1000, 5000},
				[]uint32{75, 25},
			),
			expectedFundingIndexDeltaStrings: []string{"1250"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: 2_000,
					FundingIndex:    dtypes.NewInt(1250),
				},
			},
		},
		"Success: funding index greater than MaxInt64 (doesn't overflow), 60 samples expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_0AtomicResolution,
			},
			testFundingSamples: constants.GenerateConstantFundingPremiums(math.MaxInt32, 60),
			// 8-hr funding rate capped at 6_000_000, prorated funding rate thus is 750_000
			expectedFundingIndexDeltaStrings: []string{"37500000000000000"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_0AtomicResolution.GetId(),
					FundingValuePpm: 6_000_000,
					FundingIndex:    dtypes.NewInt(37500000000000000),
				},
			},
		},
		"Success: no samples": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_0DefaultFunding_0AtomicResolution,
			},
			testFundingSamples:               []int32{},
			expectedFundingIndexDeltaStrings: []string{"0"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_0DefaultFunding_0AtomicResolution.GetId(),
					FundingValuePpm: 0,
					FundingIndex:    dtypes.NewInt(0),
				},
			},
		},
		"Success: no samples, -0.001 percent default funding, negative funding rate expected": {
			testFundingSampleDuration: 60,
			testFundingTickDuration:   3600,
			testPerpetuals: []types.Perpetual{
				constants.BtcUsd_NegativeDefaultFunding_10AtomicResolution,
			},
			testFundingSamples:               []int32{},
			expectedFundingIndexDeltaStrings: []string{"-625"},
			fundingRatesAndIndices: []indexerevents.FundingUpdateV1{
				{
					PerpetualId:     constants.BtcUsd_NegativeDefaultFunding_10AtomicResolution.GetId(),
					FundingValuePpm: -1000,
					FundingIndex:    dtypes.NewInt(-625),
				},
			},
		},
	}

	testCurrentFundingTickEpochStartBlock := uint32(23)
	testCurrentEpoch := uint32(1)

	for name, tc := range tests {
		t.Run(name, func(*testing.T) {
			pc := keepertest.PerpetualsKeepers(t)
			ctx := pc.Ctx.WithTxBytes(constants.TestTxBytes)
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, pc.PricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, pc.PerpetualsKeeper)

			// Create test perpetuals.
			// 1BTC = $50,000.
			oldPerps := make([]types.Perpetual, len(tc.testPerpetuals))
			for i, p := range tc.testPerpetuals {
				perp, err := pc.PerpetualsKeeper.CreatePerpetual(
					pc.Ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
					p.Params.MarketType,
				)
				require.NoError(t, err)
				oldPerps[i] = perp
			}

			// Create funding-tick epoch.
			err := pc.EpochsKeeper.CreateEpochInfo(
				pc.Ctx,
				epochstypes.EpochInfo{
					Name:                   string(epochstypes.FundingTickEpochInfoName),
					Duration:               tc.testFundingTickDuration,
					CurrentEpochStartBlock: testCurrentFundingTickEpochStartBlock,
					CurrentEpoch:           testCurrentEpoch,
				},
			)
			require.NoError(t, err)
			// Create funding-sample epoch.
			err = pc.EpochsKeeper.CreateEpochInfo(
				pc.Ctx,
				epochstypes.EpochInfo{
					Name:     string(epochstypes.FundingSampleEpochInfoName),
					Duration: tc.testFundingSampleDuration,
				},
			)
			require.NoError(t, err)

			// Insert test funding sample.
			keepertest.PopulateTestPremiumStore(
				t,
				pc.Ctx,
				pc.PerpetualsKeeper,
				oldPerps,
				tc.testFundingSamples,
				false, // isVote
			)

			pc.PerpetualsKeeper.MaybeProcessNewFundingTickEpoch(
				// Current block is the start of a new epoch for funding-tick.
				pc.Ctx.WithBlockHeight(int64(testCurrentFundingTickEpochStartBlock)))

			for i, p := range oldPerps {
				newPerp, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, p.Params.Id)
				require.NoError(t, err)

				// Set `expectedFundingIndexDelta` either to the provided big.Int or to a big.Int of the
				// provided string value.
				expectedFundingIndexDelta, ok := new(big.Int).SetString(
					tc.expectedFundingIndexDeltaStrings[i],
					10,
				)
				if !ok {
					panic("invalid expectedFundingIndexDeltaString")
				}

				actualFundingIndexDelta := new(big.Int).Sub(
					newPerp.FundingIndex.BigInt(),
					oldPerps[i].FundingIndex.BigInt(),
				)
				require.Equal(
					t,
					expectedFundingIndexDelta,
					actualFundingIndexDelta,
				)

				// Check that all recorded funding samples from the previous epoch were deleted.
				allSamples := pc.PerpetualsKeeper.GetPremiumSamples(pc.Ctx)
				require.NoError(t, err)
				for _, marketPremiums := range allSamples.AllMarketPremiums {
					require.Equal(t, 0, len(marketPremiums.Premiums))
				}
			}

			fundingEvents := getFundingBlockEventsFromIndexerBlock(pc.Ctx, pc.PerpetualsKeeper)
			expectedFundingEvent := indexerevents.NewFundingRatesAndIndicesEvent(
				tc.fundingRatesAndIndices,
			)
			require.Contains(t, fundingEvents, expectedFundingEvent)
		})
	}
}

// getFundingBlockEventsFromIndexerBlock returns all funding events from the indexer block.
func getFundingBlockEventsFromIndexerBlock(
	ctx sdk.Context,
	perpetualsKeeper *keeper.Keeper,
) []*indexerevents.FundingEventV1 {
	block := perpetualsKeeper.GetIndexerEventManager().ProduceBlock(ctx)
	var fundingEvents []*indexerevents.FundingEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeFundingValues {
			continue
		}
		if _, ok := event.OrderingWithinBlock.(*indexer_manager.IndexerTendermintEvent_BlockEvent_); ok {
			var fundingEvent indexerevents.FundingEventV1
			err := proto.Unmarshal(event.DataBytes, &fundingEvent)
			if err != nil {
				panic(err)
			}
			fundingEvents = append(fundingEvents, &fundingEvent)
		}
	}
	return fundingEvents
}

func TestMaybeProcessNewFundingTickEpoch_Failure(t *testing.T) {
	tests := map[string]struct {
		testEpochs         []epochstypes.EpochInfo
		testPerpetuals     []types.Perpetual
		testPremiumSamples []int32
		expectedError      error
	}{
		"No `funding-sample` epoch info found": {
			testEpochs: []epochstypes.EpochInfo{
				{
					Name:                   string(epochstypes.FundingTickEpochInfoName),
					CurrentEpochStartBlock: 23,
					CurrentEpoch:           1,
					Duration:               3600,
				},
			},
			expectedError: errorsmod.Wrapf(
				epochstypes.ErrEpochInfoNotFound,
				"name: %s",
				epochstypes.FundingSampleEpochInfoName,
			),
		},
	}

	testCurrentFundingTickEpochStartBlock := uint32(23)

	for name, tc := range tests {
		t.Run(name, func(*testing.T) {
			pc := keepertest.PerpetualsKeepers(t)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, pc.Ctx, pc.PricesKeeper)

			// Insert test funding sample.
			keepertest.PopulateTestPremiumStore(
				t,
				pc.Ctx,
				pc.PerpetualsKeeper,
				tc.testPerpetuals,
				tc.testPremiumSamples,
				false, // isVote
			)

			// Create test epochs.
			for _, epochInfo := range tc.testEpochs {
				err := pc.EpochsKeeper.CreateEpochInfo(
					pc.Ctx,
					epochInfo,
				)
				require.NoError(t, err)
			}

			initialEvents := pc.Ctx.EventManager().ABCIEvents()

			require.PanicsWithError(
				t,
				tc.expectedError.Error(),
				func() {
					pc.PerpetualsKeeper.MaybeProcessNewFundingTickEpoch(
						pc.Ctx.WithBlockHeight(int64(testCurrentFundingTickEpochStartBlock)))
				},
			)

			// Verify that no new events were emitted.
			laterEvents := pc.Ctx.EventManager().ABCIEvents()
			require.ElementsMatch(t,
				initialEvents,
				laterEvents,
			)
		})
	}
}

func TestMaybeProcessNewFundingTickEpoch_NoNewEpoch(t *testing.T) {
	testCurrentFundingTickEpochStartBlock := uint32(23)
	testCurrentEpoch := uint32(1)

	pc := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 100)

	err := pc.EpochsKeeper.CreateEpochInfo(
		pc.Ctx,
		epochstypes.EpochInfo{
			Name:                   string(epochstypes.FundingTickEpochInfoName),
			Duration:               3600,
			CurrentEpochStartBlock: testCurrentFundingTickEpochStartBlock,
			CurrentEpoch:           testCurrentEpoch,
		},
	)
	if err != nil {
		require.NoError(t, err)
	}

	pc.PerpetualsKeeper.MaybeProcessNewFundingTickEpoch(
		// Current block is not start of a new epoch for funding-tick.
		pc.Ctx.WithBlockHeight(int64(testCurrentFundingTickEpochStartBlock + 1)))

	for _, perp := range perps {
		newPerp, err := pc.PerpetualsKeeper.GetPerpetual(pc.Ctx, perp.Params.Id)
		require.NoError(t, err)

		require.Equal(t,
			perp.FundingIndex,
			newPerp.FundingIndex,
		)
	}
}

func TestGetAddPremiumVotes_NoPremiumVotes(t *testing.T) {
	testCurrentEpoch := uint32(5)
	testDuration := uint32(60)
	testCurrentFundingSampleEpochStartBlock := uint32(3)

	pc := keepertest.PerpetualsKeepers(t)
	err := pc.EpochsKeeper.CreateEpochInfo(
		pc.Ctx,
		epochstypes.EpochInfo{
			Name:                   string(epochstypes.FundingSampleEpochInfoName),
			Duration:               testDuration,
			CurrentEpochStartBlock: testCurrentFundingSampleEpochStartBlock,
			CurrentEpoch:           testCurrentEpoch,
		},
	)
	require.NoError(t, err)

	msgAddPremiumVotes := pc.PerpetualsKeeper.GetAddPremiumVotes(
		pc.Ctx.WithBlockHeight(int64(testCurrentFundingSampleEpochStartBlock)),
	)
	// We don't panic but only log an error if there are no new premium votes.
	require.Equal(t, 0, len(msgAddPremiumVotes.Votes))
}

func TestGetAddPremiumVotes_Success(t *testing.T) {
	testCurrentEpoch := uint32(1)
	testDuration := uint32(60)

	tests := map[string]struct {
		currentFundingSampleEpochStartBlock uint32
		blockHeight                         int64
		samplePremiumPpm                    int32
		numPerpetuals                       int
		// Should be <= `numPerpetuals`.
		numPerpetualsWithValidIndexPrice int
		expectedNumSamples               int
	}{
		"Positive premium": {
			currentFundingSampleEpochStartBlock: 23,
			blockHeight:                         23,
			samplePremiumPpm:                    100,
			numPerpetuals:                       10,
			numPerpetualsWithValidIndexPrice:    10,
			expectedNumSamples:                  10,
		},
		"Positive premium, only 1 perpetual has valid index price": {
			currentFundingSampleEpochStartBlock: 23,
			blockHeight:                         23,
			samplePremiumPpm:                    100,
			numPerpetuals:                       10,
			numPerpetualsWithValidIndexPrice:    1,
			expectedNumSamples:                  1,
		},
		"Negative premium": {
			currentFundingSampleEpochStartBlock: 24,
			blockHeight:                         24,
			samplePremiumPpm:                    -150,
			numPerpetuals:                       10,
			numPerpetualsWithValidIndexPrice:    10,
			expectedNumSamples:                  10,
		},
		"Not start of new funding-sample epoch, still produce samples": {
			currentFundingSampleEpochStartBlock: 24,
			blockHeight:                         25,
			samplePremiumPpm:                    100,
			numPerpetuals:                       10,
			numPerpetualsWithValidIndexPrice:    10,
			expectedNumSamples:                  10,
		},
		"Zero premiums": {
			currentFundingSampleEpochStartBlock: 24,
			blockHeight:                         24,
			samplePremiumPpm:                    0,
			numPerpetuals:                       10,
			numPerpetualsWithValidIndexPrice:    10,
			expectedNumSamples:                  0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockPricePremiumGetter := mocks.PerpetualsClobKeeper{}
			mockPricePremiumGetter.On(
				"GetPricePremiumForPerpetual",
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
				mock.Anything,
			).Return(tc.samplePremiumPpm, nil)

			pc := keepertest.PerpetualsKeepersWithClobHelpers(t, &mockPricePremiumGetter)

			// MockTimeProvider needed for to use `constants.TimeT` as cutoff time of index price cache query.
			pc.MockTimeProvider.On("Now").Return(constants.TimeT)

			pc.IndexPriceCache.UpdatePrices(
				pricefeed_testutil.GetTestMarketPriceUpdates(
					tc.numPerpetualsWithValidIndexPrice,
				),
			)

			// Create liquidity tiers and perpetuals,
			_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, tc.numPerpetuals)

			err := pc.EpochsKeeper.CreateEpochInfo(
				pc.Ctx,
				epochstypes.EpochInfo{
					Name:                   string(epochstypes.FundingSampleEpochInfoName),
					Duration:               testDuration,
					CurrentEpochStartBlock: tc.currentFundingSampleEpochStartBlock,
					CurrentEpoch:           testCurrentEpoch,
				},
			)
			if err != nil {
				require.NoError(t, err)
			}

			msgAddPremiumVotes := pc.PerpetualsKeeper.GetAddPremiumVotes(
				pc.Ctx.WithBlockHeight(int64(tc.blockHeight)),
			)

			mockPricePremiumGetter.AssertNumberOfCalls(
				t,
				"GetPricePremiumForPerpetual",
				tc.numPerpetualsWithValidIndexPrice,
			)

			// Check that new premium votes are returned.
			require.NotNil(t, msgAddPremiumVotes)
			require.Equal(t, tc.expectedNumSamples, len(msgAddPremiumVotes.Votes))

			require.True(t,
				sort.SliceIsSorted(msgAddPremiumVotes.Votes, func(p, q int) bool {
					return msgAddPremiumVotes.Votes[p].PerpetualId < msgAddPremiumVotes.Votes[q].PerpetualId
				}))

			for i, sample := range msgAddPremiumVotes.Votes {
				if i > 0 {
					// Check samples are unique and are sorted by perpetual.Params.Id
					require.True(
						t,
						msgAddPremiumVotes.Votes[i-1].PerpetualId < sample.PerpetualId,
					)
				}
				require.Equal(t, tc.samplePremiumPpm, sample.PremiumPpm)
			}
		})
	}
}

func TestGetPremiumStore_DefaultValue(t *testing.T) {
	testCases := map[string]struct {
		getPremiumFunc func(
			*keeper.Keeper,
			sdk.Context,
		) types.PremiumStore
	}{
		"GetPremiumSamples": {
			getPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
			) types.PremiumStore {
				return keeper.GetPremiumSamples(ctx)
			},
		},
		"GetPremiumVotes": {
			getPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
			) types.PremiumStore {
				return keeper.GetPremiumVotes(ctx)
			},
		},
	}

	for _, tc := range testCases {
		pc := keepertest.PerpetualsKeepers(t)

		premiumSamples := tc.getPremiumFunc(pc.PerpetualsKeeper, pc.Ctx)
		require.Equal(t, 0, len(premiumSamples.AllMarketPremiums))
	}
}

func TestAddPremiums_Success(t *testing.T) {
	testCases := map[string]struct {
		addPremiumFunc func(
			*keeper.Keeper,
			sdk.Context,
			[]types.FundingPremium,
		) error
		getPremiumFunc func(
			*keeper.Keeper,
			sdk.Context,
		) types.PremiumStore
	}{
		"AddPremiumSamples": {
			addPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
				samples []types.FundingPremium,
			) error {
				return keeper.AddPremiumSamples(ctx, samples)
			},
			getPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
			) types.PremiumStore {
				return keeper.GetPremiumSamples(ctx)
			},
		},
		"AddPremiumVotes": {
			addPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
				votes []types.FundingPremium,
			) error {
				return keeper.AddPremiumVotes(ctx, votes)
			},
			getPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
			) types.PremiumStore {
				return keeper.GetPremiumVotes(ctx)
			},
		},
	}

	for _, tc := range testCases {
		pc := keepertest.PerpetualsKeepers(t)

		// Create liquidity tiers and perpetuals,
		numPerpetuals := 10
		perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, numPerpetuals)

		// Insert one round of premiums for all perps.
		firstPremiums := make([]types.FundingPremium, numPerpetuals)
		for i, perp := range perps {
			firstPremiums[i] = types.FundingPremium{
				PerpetualId: perp.Params.Id,
				// -1000 for even Ids, +1000 for odd Ids.
				PremiumPpm: 1_000 * (2*(int32(perp.Params.Id)%2) - 1),
			}
		}

		err := tc.addPremiumFunc(pc.PerpetualsKeeper, pc.Ctx, firstPremiums)
		require.NoError(t, err)

		// Check each perp has expected number of premiums stored after first around of addPremiumFunc().
		firstStoredPremiums := tc.getPremiumFunc(pc.PerpetualsKeeper, pc.Ctx)

		require.Equal(t,
			uint32(1),
			firstStoredPremiums.NumPremiums,
		)

		marketSamplesMap := firstStoredPremiums.GetMarketPremiumsMap()

		for _, perp := range perps {
			entries := marketSamplesMap[perp.Params.Id].Premiums

			require.Equal(t,
				1,
				len(entries),
			)
			require.Equal(t,
				1000*(2*(int32(perp.Params.Id)%2)-1),
				entries[0],
			)
		}

		// Insert another round of samples for only perps with even Ids.
		secondPremiums := make([]types.FundingPremium, numPerpetuals/2)
		for i := range secondPremiums {
			secondPremiums[i] = types.FundingPremium{
				PerpetualId: uint32(2 * i),
				PremiumPpm:  -1_000,
			}
		}
		err = tc.addPremiumFunc(
			pc.PerpetualsKeeper,
			pc.Ctx,
			secondPremiums,
		)
		require.NoError(t, err)

		// Check each perp has expected number of premiums stored after second round of addPremiumFunc().
		secondStoredPremiums := tc.getPremiumFunc(pc.PerpetualsKeeper, pc.Ctx)

		require.Equal(t,
			uint32(2),
			secondStoredPremiums.NumPremiums,
		)

		marketSamplesMap = secondStoredPremiums.GetMarketPremiumsMap()

		for _, perp := range perps {
			entries := marketSamplesMap[perp.Params.Id].Premiums
			if perp.Params.Id%2 == 0 {
				// Even perpetuals should have two samples of -1000.
				require.Equal(t,
					2,
					len(entries),
				)
				require.Equal(t,
					[]int32{-1000, -1000},
					entries,
				)
			} else {
				// Odd perpetuals shold have one sample of 1000.
				require.Equal(t,
					1,
					len(entries),
				)
				require.Equal(t,
					[]int32{1000},
					entries,
				)
			}
		}
	}
}

func TestAddPremiums_NonExistingPerpetuals(t *testing.T) {
	testCases := map[string]struct {
		addPremiumFunc func(
			*keeper.Keeper,
			sdk.Context,
			[]types.FundingPremium,
		) error
	}{
		"AddPremiumSamples": {
			addPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
				samples []types.FundingPremium,
			) error {
				return keeper.AddPremiumSamples(ctx, samples)
			},
		},
		"AddPremiumVotes": {
			addPremiumFunc: func(
				keeper *keeper.Keeper,
				ctx sdk.Context,
				votes []types.FundingPremium,
			) error {
				return keeper.AddPremiumVotes(ctx, votes)
			},
		},
	}

	for _, tc := range testCases {
		pc := keepertest.PerpetualsKeepers(t)
		nonExistentPerpetualId := uint32(1000)

		newPremiums := []types.FundingPremium{
			{
				PerpetualId: nonExistentPerpetualId,
				PremiumPpm:  -1_000,
			},
		}

		// Create liquidity tiers and perpetuals,
		_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 3)

		err := tc.addPremiumFunc(pc.PerpetualsKeeper, pc.Ctx, newPremiums)
		require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
		require.Error(t,
			err,
			errorsmod.Wrapf(
				types.ErrPerpetualDoesNotExist,
				"perpetual ID = %d",
				1000,
			).Error(),
		)
	}
}

func TestMaybeProcessNewFundingSampleEpoch(t *testing.T) {
	testDuration := uint32(60)
	testCurrentEpoch := uint32(5)

	tests := map[string]struct {
		currentEpochStartBlock uint32
		currentBlockHeight     int64
		minNumVotesPerSample   uint32
		premiumVotes           types.PremiumStore
		prevPremiumSamples     types.PremiumStore
		expectedPremiumSamples types.PremiumStore
		expectedPremiumVotes   types.PremiumStore
		panicErr               error
		expectedEvents         []sdk.Event
	}{
		"Not new epoch": {
			currentEpochStartBlock: 23,
			currentBlockHeight:     25,
			minNumVotesPerSample:   1,
			premiumVotes: types.PremiumStore{
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{51, 51, -100, -100},
					},
				},
			},
			prevPremiumSamples:     types.PremiumStore{},
			expectedPremiumSamples: types.PremiumStore{},
			expectedPremiumVotes: types.PremiumStore{
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{51, 51, -100, -100},
					},
				},
			},
		},
		"New epoch, empty premium samples storage": {
			currentEpochStartBlock: 23,
			currentBlockHeight:     23,
			minNumVotesPerSample:   1,
			premiumVotes: types.PremiumStore{
				NumPremiums: 4,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{51, 51, -100, -100},
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{51, 51, 100, 100},
					},
				},
			},
			prevPremiumSamples: types.PremiumStore{},
			expectedPremiumSamples: types.PremiumStore{
				NumPremiums: 1,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{-25},
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{76},
					},
				},
			},
			expectedPremiumVotes: types.PremiumStore{}, // reset to empty
		},
		"New epoch, add new sample to existing samples, skip zero samples": {
			currentEpochStartBlock: 23,
			currentBlockHeight:     23,
			minNumVotesPerSample:   1,
			premiumVotes: types.PremiumStore{
				NumPremiums: 6,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{-1000, 1000, 51, -50, 100, -100}, // median = 1
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{0, 0, 1, 2, 3, 4}, // median = 2
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{-1000, -500, -5, 5, 500, 1000}, // median = 0
					},
				},
			},
			prevPremiumSamples: types.PremiumStore{
				NumPremiums: 2,

				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{100, 101},
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{1000},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{-1000},
					},
				},
			},
			expectedPremiumSamples: types.PremiumStore{
				NumPremiums: 3,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{100, 101, 1},
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{1000, 2},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{-1000}, // unchanged
					},
				},
			},
			expectedPremiumVotes: types.PremiumStore{}, // reset to empty
		},
		"New epoch, add zero paddings, NumPremiums > MinNumVotesPerSample": {
			currentEpochStartBlock: 23,
			currentBlockHeight:     23,
			minNumVotesPerSample:   2,
			premiumVotes: types.PremiumStore{
				NumPremiums: 4,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1000}, // median([0, 0, 0, 1000]) = 0
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{-5, -10}, // median([-10, -5, 0, 0]) = -3
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{200, -100, 100}, // median([-100, 0, 100, 200]) = 50
					},
				},
			},
			prevPremiumSamples: types.PremiumStore{},
			expectedPremiumSamples: types.PremiumStore{
				NumPremiums: 1,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 1,
						Premiums:    []int32{-3},
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{50},
					},
				},
			},
			expectedPremiumVotes: types.PremiumStore{}, // reset to empty
		},
		"New epoch, add zero paddings, NumPremiums < MinNumVotesPerSample": {
			currentEpochStartBlock: 23,
			currentBlockHeight:     23,
			minNumVotesPerSample:   6,
			premiumVotes: types.PremiumStore{
				NumPremiums: 5,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1000}, // median([0, 0, 0, 0, 0, 1000]) = 0
					},
					{
						PerpetualId: 1,
						Premiums:    []int32{20, -20}, // median([-20, 0, 0, 0, 0, 20]) = 0
					},
					{
						PerpetualId: 2,
						Premiums:    []int32{-5, -10, 2, -1}, // median([-10, -5, -1, 0, 0, 2]) = -1
					},
					{
						PerpetualId: 3,
						Premiums:    []int32{200, -100, 100, 30, 40}, // median([-100, 0, 30, 40, 100, 200]) = 35
					},
				},
			},
			prevPremiumSamples: types.PremiumStore{},
			expectedPremiumSamples: types.PremiumStore{
				NumPremiums: 1,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 2,
						Premiums:    []int32{-1},
					},
					{
						PerpetualId: 3,
						Premiums:    []int32{35},
					},
				},
			},
			expectedPremiumVotes: types.PremiumStore{}, // reset to empty
		},
		"Panic: `NumPremiums` < premium entries length": {
			currentEpochStartBlock: 23,
			currentBlockHeight:     23,
			minNumVotesPerSample:   1,
			premiumVotes: types.PremiumStore{
				NumPremiums: 1,
				AllMarketPremiums: []types.MarketPremiums{
					{
						PerpetualId: 0,
						Premiums:    []int32{1, 2, 3},
					},
				},
			},
			panicErr: fmt.Errorf(
				"marketPremiums (%+v) has more non-zero premiums than total number of premiums (%d)",
				types.MarketPremiums{
					PerpetualId: 0,
					Premiums:    []int32{1, 2, 3},
				},
				1,
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pc := keepertest.PerpetualsKeepers(t)
			pc.Ctx = pc.Ctx.WithTxBytes(constants.TestTxBytes)

			// Create funding-sample epoch.
			err := pc.EpochsKeeper.CreateEpochInfo(
				pc.Ctx,
				epochstypes.EpochInfo{
					Name:                   string(epochstypes.FundingSampleEpochInfoName),
					Duration:               testDuration,
					CurrentEpochStartBlock: tc.currentEpochStartBlock,
					CurrentEpoch:           testCurrentEpoch,
				},
			)
			require.NoError(t, err)

			// Create liquidity tiers and perpetuals,
			_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 4)
			require.NoError(t, err)

			params := pc.PerpetualsKeeper.GetParams(pc.Ctx)
			err = pc.PerpetualsKeeper.SetParams(
				pc.Ctx,
				types.Params{
					FundingRateClampFactorPpm: params.FundingRateClampFactorPpm,
					PremiumVoteClampFactorPpm: params.PremiumVoteClampFactorPpm,
					MinNumVotesPerSample:      tc.minNumVotesPerSample,
				},
			)
			require.NoError(t, err)
			pc.PerpetualsKeeper.SetPremiumVotes(pc.Ctx, tc.premiumVotes)
			pc.PerpetualsKeeper.SetPremiumSamples(pc.Ctx, tc.prevPremiumSamples)

			initialEvents := pc.Ctx.EventManager().ABCIEvents()

			if tc.panicErr != nil {
				require.PanicsWithError(
					t,
					tc.panicErr.Error(),
					func() {
						pc.PerpetualsKeeper.MaybeProcessNewFundingSampleEpoch(pc.Ctx.WithBlockHeight(tc.currentBlockHeight))
					},
				)

				laterEvents := pc.Ctx.EventManager().ABCIEvents()
				require.ElementsMatch(t,
					initialEvents,
					laterEvents,
				)
				return
			}

			pc.PerpetualsKeeper.MaybeProcessNewFundingSampleEpoch(pc.Ctx.WithBlockHeight(tc.currentBlockHeight))

			require.Equal(t,
				tc.expectedPremiumVotes,
				pc.PerpetualsKeeper.GetPremiumVotes(pc.Ctx),
			)

			require.Equal(t,
				tc.expectedPremiumSamples,
				pc.PerpetualsKeeper.GetPremiumSamples(pc.Ctx),
			)
		})
	}
}

func TestGetAllLiquidityTiers_Sorted(t *testing.T) {
	// Setup context and keepers
	pc := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals
	// keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper)
	lts := []types.LiquidityTier{
		*lttest.GenerateLiquidityTier(lttest.WithId(0)),
		*lttest.GenerateLiquidityTier(lttest.WithId(100)),
		*lttest.GenerateLiquidityTier(lttest.WithId(5)),
		*lttest.GenerateLiquidityTier(lttest.WithId(72)),
		*lttest.GenerateLiquidityTier(lttest.WithId(16)),
	}

	for _, lt := range lts {
		_, err := pc.PerpetualsKeeper.SetLiquidityTier(
			pc.Ctx,
			lt.Id,
			lt.Name,
			lt.InitialMarginPpm,
			lt.MaintenanceFractionPpm,
			lt.ImpactNotional,
			lt.OpenInterestLowerCap,
			lt.OpenInterestUpperCap,
		)
		require.NoError(t, err)
	}

	got := pc.PerpetualsKeeper.GetAllLiquidityTiers(pc.Ctx)
	require.Equal(
		t,
		[]types.LiquidityTier{
			*lttest.GenerateLiquidityTier(lttest.WithId(0)),
			*lttest.GenerateLiquidityTier(lttest.WithId(5)),
			*lttest.GenerateLiquidityTier(lttest.WithId(16)),
			*lttest.GenerateLiquidityTier(lttest.WithId(72)),
			*lttest.GenerateLiquidityTier(lttest.WithId(100)),
		},
		got,
	)
}

func TestHasLiquidityTier(t *testing.T) {
	// Setup context and keepers
	pc := keepertest.PerpetualsKeepers(t)

	lts := []types.LiquidityTier{
		*lttest.GenerateLiquidityTier(lttest.WithId(0)),
		*lttest.GenerateLiquidityTier(lttest.WithId(5)),
		*lttest.GenerateLiquidityTier(lttest.WithId(16)),
		*lttest.GenerateLiquidityTier(lttest.WithId(72)),
		*lttest.GenerateLiquidityTier(lttest.WithId(100)),
	}

	for _, lt := range lts {
		_, err := pc.PerpetualsKeeper.SetLiquidityTier(
			pc.Ctx,
			lt.Id,
			lt.Name,
			lt.InitialMarginPpm,
			lt.MaintenanceFractionPpm,
			lt.ImpactNotional,
			lt.OpenInterestLowerCap,
			lt.OpenInterestUpperCap,
		)
		require.NoError(t, err)
	}

	for _, lt := range lts {
		// Test if HasLiquidityTier correctly identifies an existing liquidity tier.
		require.True(t, pc.PerpetualsKeeper.HasLiquidityTier(pc.Ctx, lt.Id))
	}

	found := pc.PerpetualsKeeper.HasLiquidityTier(pc.Ctx, 9999)
	require.False(t, found, "Expected not to find liquidity tier with id 9999, but it was found")
}

func TestCreateLiquidityTier_Success(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	for _, lt := range constants.LiquidityTiers {
		// Create LiquidityTier without error.
		_, err := pc.PerpetualsKeeper.SetLiquidityTier(
			pc.Ctx,
			lt.Id,
			lt.Name,
			lt.InitialMarginPpm,
			lt.MaintenanceFractionPpm,
			lt.ImpactNotional,
			lt.OpenInterestLowerCap,
			lt.OpenInterestUpperCap,
		)
		require.NoError(t, err)

		// Validate liquidity tier exists in store.
		require.True(t, pc.PerpetualsKeeper.HasLiquidityTier(pc.Ctx, lt.Id))

		// Validate fields of LiquidityTier object in store.
		liquidityTier, err := pc.PerpetualsKeeper.GetLiquidityTier(pc.Ctx, lt.Id)
		require.NoError(t, err)
		require.Equal(t, lt.Id, liquidityTier.Id)
		require.Equal(t, lt.Name, liquidityTier.Name)
		require.Equal(t, lt.InitialMarginPpm, liquidityTier.InitialMarginPpm)
		require.Equal(t, lt.MaintenanceFractionPpm, liquidityTier.MaintenanceFractionPpm)
		require.Equal(t, lt.ImpactNotional, liquidityTier.ImpactNotional)
	}
}

func TestSetLiquidityTier_New_Failure(t *testing.T) {
	tests := map[string]struct {
		id                     uint32
		name                   string
		initialMarginPpm       uint32
		maintenanceFractionPpm uint32
		impactNotional         uint64
		openInterestLowerCap   uint64
		openInterestUpperCap   uint64
		expectedError          error
	}{
		"Initial Margin Ppm exceeds maximum": {
			id:                     0,
			name:                   "Large-Cap",
			initialMarginPpm:       lib.OneMillion + 1,
			maintenanceFractionPpm: 500_000,
			impactNotional:         uint64(lib.OneMillion),
			openInterestLowerCap:   0,
			openInterestUpperCap:   0,
			expectedError:          errorsmod.Wrap(types.ErrInitialMarginPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Maintenance Fraction Ppm exceeds maximum": {
			id:                     1,
			name:                   "Medium-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion + 1,
			impactNotional:         uint64(lib.OneMillion),
			openInterestLowerCap:   0,
			openInterestUpperCap:   0,
			expectedError:          errorsmod.Wrap(types.ErrMaintenanceFractionPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Impact Notional is zero": {
			id:                     1,
			name:                   "Small-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion,
			impactNotional:         uint64(0),
			openInterestLowerCap:   0,
			openInterestUpperCap:   0,
			expectedError:          types.ErrImpactNotionalIsZero,
		},
	}

	// Test setup.
	pc := keepertest.PerpetualsKeepers(t)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := pc.PerpetualsKeeper.SetLiquidityTier(
				pc.Ctx,
				tc.id,
				tc.name,
				tc.initialMarginPpm,
				tc.maintenanceFractionPpm,
				tc.impactNotional,
				tc.openInterestLowerCap,
				tc.openInterestUpperCap,
			)

			require.Error(t, err)
			require.EqualError(t, err, tc.expectedError.Error())
		})
	}
}

func TestModifyLiquidityTier_Success(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	for _, lt := range constants.LiquidityTiers {
		_, err := pc.PerpetualsKeeper.SetLiquidityTier(
			pc.Ctx,
			lt.Id,
			lt.Name,
			lt.InitialMarginPpm,
			lt.MaintenanceFractionPpm,
			lt.ImpactNotional,
			lt.OpenInterestLowerCap,
			lt.OpenInterestUpperCap,
		)
		require.NoError(t, err)
	}

	for i, lt := range constants.LiquidityTiers {
		// Modify each field arbitrarily and
		// verify the fields are modified in state.
		name := fmt.Sprintf("foo_%v", i)
		initialMarginPpm := uint32(i * 2)
		maintenanceFractionPpm := uint32(i * 2)
		impactNotional := uint64((i + 1) * 500_000_000)
		openInterestLowerCap := uint64(0)
		openInterestUpperCap := uint64(0)
		modifiedLt, err := pc.PerpetualsKeeper.SetLiquidityTier(
			pc.Ctx,
			lt.Id,
			name,
			initialMarginPpm,
			maintenanceFractionPpm,
			impactNotional,
			openInterestLowerCap,
			openInterestUpperCap,
		)
		require.NoError(t, err)
		obtainedLt, err := pc.PerpetualsKeeper.GetLiquidityTier(pc.Ctx, lt.Id)
		require.NoError(t, err)
		require.Equal(
			t,
			modifiedLt,
			obtainedLt,
		)
		require.Equal(
			t,
			name,
			obtainedLt.Name,
		)
		require.Equal(
			t,
			initialMarginPpm,
			obtainedLt.InitialMarginPpm,
		)
		require.Equal(
			t,
			maintenanceFractionPpm,
			obtainedLt.MaintenanceFractionPpm,
		)
		require.Equal(
			t,
			impactNotional,
			obtainedLt.ImpactNotional,
		)
	}
	liquidityTierUpsertEvents := keepertest.GetLiquidityTierUpsertEventsFromIndexerBlock(pc.Ctx, pc.PerpetualsKeeper)
	require.Len(t, liquidityTierUpsertEvents, len(constants.LiquidityTiers)*2)
}

func TestSetLiquidityTier_Existing_Failure(t *testing.T) {
	tests := map[string]struct {
		id                     uint32
		name                   string
		initialMarginPpm       uint32
		maintenanceFractionPpm uint32
		impactNotional         uint64
		openInterestLowerCap   uint64
		openInterestUpperCap   uint64
		expectedError          error
	}{
		"Initial Margin Ppm exceeds maximum": {
			id:                     0,
			name:                   "Large-Cap",
			initialMarginPpm:       lib.OneMillion + 1,
			maintenanceFractionPpm: 500_000,
			impactNotional:         uint64(lib.OneMillion),
			openInterestLowerCap:   0,
			openInterestUpperCap:   0,
			expectedError:          errorsmod.Wrap(types.ErrInitialMarginPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Maintenance Fraction Ppm exceeds maximum": {
			id:                     1,
			name:                   "Medium-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion + 1,
			impactNotional:         uint64(lib.OneMillion),
			openInterestLowerCap:   0,
			openInterestUpperCap:   0,
			expectedError:          errorsmod.Wrap(types.ErrMaintenanceFractionPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Impact Notional is zero": {
			id:                     1,
			name:                   "Small-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion,
			impactNotional:         uint64(0),
			openInterestLowerCap:   0,
			openInterestUpperCap:   0,
			expectedError:          types.ErrImpactNotionalIsZero,
		},
		"Invalid open interest caps": {
			id:                     1,
			name:                   "Small-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion,
			impactNotional:         uint64(lib.OneMillion),
			openInterestLowerCap:   50_000_000_000_000,
			openInterestUpperCap:   25_000_000_000_000,
			expectedError: errorsmod.Wrapf(
				types.ErrOpenInterestLowerCapLargerThanUpperCap,
				"open_interest_lower_cap: %d, open_interest_upper_cap: %d",
				50_000_000_000_000,
				25_000_000_000_000,
			),
		},
	}

	// Test setup.
	pc := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers.
	keepertest.CreateTestLiquidityTiers(t, pc.Ctx, pc.PerpetualsKeeper)
	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := pc.PerpetualsKeeper.SetLiquidityTier(
				pc.Ctx,
				tc.id,
				tc.name,
				tc.initialMarginPpm,
				tc.maintenanceFractionPpm,
				tc.impactNotional,
				tc.openInterestLowerCap,
				tc.openInterestUpperCap,
			)

			require.Error(t, err)
			require.EqualError(t, err, tc.expectedError.Error())
		})
	}
}

func TestSetParams(t *testing.T) {
	tests := map[string]struct {
		params      types.Params
		expectedErr string
	}{
		"Success": {
			params: types.Params{
				FundingRateClampFactorPpm: 6_000_000,
				PremiumVoteClampFactorPpm: 60_000_000,
				MinNumVotesPerSample:      15,
			},
		},
		"Failure: Funding Rate Clamp is 0": {
			params: types.Params{
				FundingRateClampFactorPpm: 0,
				PremiumVoteClampFactorPpm: 60_000_000,
				MinNumVotesPerSample:      15,
			},
			expectedErr: types.ErrFundingRateClampFactorPpmIsZero.Error(),
		},
		"Failure: Premium Vote Clamp is 0": {
			params: types.Params{
				FundingRateClampFactorPpm: 6_000_000,
				PremiumVoteClampFactorPpm: 0,
				MinNumVotesPerSample:      15,
			},
			expectedErr: types.ErrPremiumVoteClampFactorPpmIsZero.Error(),
		},
	}

	// Test setup.
	pc := keepertest.PerpetualsKeepers(t)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			initialParams := pc.PerpetualsKeeper.GetParams(pc.Ctx)

			// Set Params.
			err := pc.PerpetualsKeeper.SetParams(pc.Ctx, tc.params)

			if tc.expectedErr == "" {
				require.NoError(t, err)
				// Check that params in store are updated.
				got := pc.PerpetualsKeeper.GetParams(pc.Ctx)
				require.Equal(t, tc.params, got)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
				// Check that params in store are unchanged.
				got := pc.PerpetualsKeeper.GetParams(pc.Ctx)
				require.Equal(t, initialParams, got)
			}
		})
	}
}

func TestIsPositionUpdatable(t *testing.T) {
	testCases := map[string]struct {
		perp              types.Perpetual
		marketParamPrice  pricestypes.MarketParamPrice
		queryPerpId       uint32
		expectedUpdatable bool
		expectedErr       string
	}{
		"Updatable": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithId(1),
				perptest.WithMarketId(1),
			),
			queryPerpId: 1,
			marketParamPrice: *pricestest.GenerateMarketParamPrice(
				pricestest.WithId(1),
				pricestest.WithPriceValue(1000), // non-zero
			),
			expectedUpdatable: true,
		},
		"Not updatable due to zero oracle price": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithId(1),
				perptest.WithMarketId(1),
			),
			queryPerpId: 1,
			marketParamPrice: *pricestest.GenerateMarketParamPrice(
				pricestest.WithId(1),
				pricestest.WithPriceValue(0),
			),
			expectedUpdatable: false,
		},
		"Error: Perp Id not found": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithId(1),
				perptest.WithMarketId(1),
			),
			queryPerpId: 100, // doesn't exist
			marketParamPrice: *pricestest.GenerateMarketParamPrice(
				pricestest.WithId(1),
			),
			expectedErr: "Perpetual does not exist",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			pc := keepertest.PerpetualsKeepers(t)
			keepertest.CreateTestPricesAndPerpetualMarkets(
				t,
				pc.Ctx,
				pc.PerpetualsKeeper,
				pc.PricesKeeper,
				[]types.Perpetual{tc.perp},
				[]pricestypes.MarketParamPrice{tc.marketParamPrice},
			)

			updatable, err := pc.PerpetualsKeeper.IsPositionUpdatable(pc.Ctx, tc.queryPerpId)
			if tc.expectedErr == "" {
				require.NoError(t, err)
				require.Equal(t, tc.expectedUpdatable, updatable)
			} else {
				require.ErrorContains(t, err, tc.expectedErr)
			}
		})
	}
}

func TestIsIsolatedPerpetual(t *testing.T) {
	testCases := map[string]struct {
		perp     types.Perpetual
		expected bool
	}{
		"Isolated Perpetual": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithMarketType(types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED),
			),
			expected: true,
		},
		"Cross Perpetual": {
			perp: *perptest.GeneratePerpetual(
				perptest.WithMarketType(types.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS),
			),
			expected: false,
		},
	}

	for name, tc := range testCases {
		t.Run(
			name, func(t *testing.T) {
				pc := keepertest.PerpetualsKeepers(t)
				pc.PerpetualsKeeper.SetPerpetualForTest(pc.Ctx, tc.perp)
				isIsolated, err := pc.PerpetualsKeeper.IsIsolatedPerpetual(pc.Ctx, tc.perp.Params.Id)
				require.NoError(t, err)
				require.Equal(t, tc.expected, isIsolated)
			},
		)
	}
}

func TestAcquireNextPerpetualID(t *testing.T) {
	pc := keepertest.PerpetualsKeepers(t)
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, pc.Ctx, pc.PerpetualsKeeper, pc.PricesKeeper, 10)

	// Get the highest perpetual id from the existing markets
	highestPerpetualId := uint32(0)
	for _, perp := range perps {
		if perp.Params.Id > highestPerpetualId {
			highestPerpetualId = perp.Params.Id
		}
	}

	// Acquire the next perpetual id
	nextPerpetualId := pc.PerpetualsKeeper.AcquireNextPerpetualID(pc.Ctx)
	require.Equal(t, highestPerpetualId+1, nextPerpetualId)

	// Verify the next perpetual id is stored in the store
	nextPerpetualIdFromStore := pc.PerpetualsKeeper.GetNextPerpetualID(pc.Ctx)
	require.Equal(t, nextPerpetualId+1, nextPerpetualIdFromStore)

	// Create a perpetual with the next perpetual ID outside of acquire flow
	perp := perptest.GeneratePerpetual(
		perptest.WithId(nextPerpetualIdFromStore),
		perptest.WithMarketId(nextPerpetualIdFromStore),
	)
	pc.PerpetualsKeeper.SetPerpetualForTest(pc.Ctx, *perp)

	// Verify the next perpetual id is incremented
	nextPerpetualId = pc.PerpetualsKeeper.AcquireNextPerpetualID(pc.Ctx)
	require.Equal(t, nextPerpetualIdFromStore+1, nextPerpetualId)
}
