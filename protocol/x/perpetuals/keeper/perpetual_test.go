package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/common"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	perptest "github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestModifyPerpetual_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	numLiquidityTiers := 4
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 100)
	numMarkets := pricesKeeper.GetNumMarkets(ctx)
	for i, item := range perps {
		// Modify each field arbitrarily and
		// verify the fields were modified in state.
		ticker := fmt.Sprintf("foo_%v", i)
		marketId := uint32(i*2) % numMarkets
		defaultFundingPpm := int32(i * 2)
		liquidityTier := uint32((i + 1) % numLiquidityTiers)
		retItem, err := keeper.ModifyPerpetual(
			ctx,
			item.Params.Id,
			ticker,
			marketId,
			defaultFundingPpm,
			liquidityTier,
		)
		require.NoError(t, err)
		newItem, err := keeper.GetPerpetual(ctx, item.Params.Id)
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
}

func TestCreatePerpetual_Failure(t *testing.T) {
	tests := map[string]struct {
		id                uint32
		ticker            string
		marketId          uint32
		atomicResolution  int32
		defaultFundingPpm int32
		liquidityTier     uint32
		expectedError     error
	}{
		"Price doesn't exist": {
			id:                0,
			ticker:            "ticker",
			marketId:          999,
			atomicResolution:  -10,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			expectedError:     errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, fmt.Sprint(999)),
		},
		"Positive default funding magnitude exceeds maximum": {
			id:                0,
			ticker:            "ticker",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: int32(lib.OneMillion + 1),
			liquidityTier:     0,
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
			expectedError:     errorsmod.Wrap(types.ErrDefaultFundingPpmMagnitudeExceedsMax, fmt.Sprint(math.MinInt32)),
		},
		"Ticker is an empty string": {
			id:                0,
			ticker:            "",
			marketId:          0,
			atomicResolution:  -10,
			defaultFundingPpm: 0,
			liquidityTier:     0,
			expectedError:     types.ErrTickerEmptyString,
		},
	}

	// Test setup.
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	keepertest.CreateNMarkets(t, ctx, pricesKeeper, 1)
	// Create Liquidity Tiers
	keepertest.CreateTestLiquidityTiers(t, ctx, keeper)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := keeper.CreatePerpetual(
				ctx,
				tc.id,
				tc.ticker,
				tc.marketId,
				tc.atomicResolution,
				tc.defaultFundingPpm,
				tc.liquidityTier,
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
	}

	// Test setup.
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 1)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := keeper.ModifyPerpetual(
				ctx,
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

func TestGetPerpetual_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 10)

	for _, perp := range perps {
		rst, err := keeper.GetPerpetual(ctx,
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
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals
	keepertest.CreateTestLiquidityTiers(t, ctx, keeper)

	perps := []types.Perpetual{
		*perptest.GeneratePerpetual(perptest.WithId(0)),
		*perptest.GeneratePerpetual(perptest.WithId(5)),
		*perptest.GeneratePerpetual(perptest.WithId(20)),
		*perptest.GeneratePerpetual(perptest.WithId(999)),
	}

	_, err := pricesKeeper.CreateMarket(
		ctx,
		// `ExchangeConfigJson` is left unset as it is not used by the server.
		pricestypes.MarketParam{
			Id:                0,
			Pair:              "marketName",
			Exponent:          -10,
			MinExchanges:      uint32(1),
			MinPriceChangePpm: uint32(50),
		},
		pricestypes.MarketPrice{
			Id:       0,
			Exponent: -10,
			Price:    1_000, // leave this as a placeholder b/c we cannot set the price to 0
		},
	)
	require.NoError(t, err)

	for perp := range perps {
		_, err := keeper.CreatePerpetual(
			ctx,
			perps[perp].Params.Id,
			perps[perp].Params.Ticker,
			perps[perp].Params.MarketId,
			perps[perp].Params.AtomicResolution,
			perps[perp].Params.DefaultFundingPpm,
			perps[perp].Params.LiquidityTier,
		)
		require.NoError(t, err)
	}

	for _, perp := range perps {
		// Test if HasPerpetual correctly identifies an existing perpetual
		found := keeper.HasPerpetual(ctx, perp.Params.Id)
		require.True(t, found, "Expected to find perpetual with id %d, but it was not found", perp.Params.Id)
	}

	found := keeper.HasPerpetual(ctx, 9999)
	require.False(t, found, "Expected not to find perpetual with id 9999, but it was found")
}

func TestGetPerpetual_NotFound(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := keeper.GetPerpetual(
		ctx,
		nonExistentPerpetualId,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetPerpetuals_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 10)

	require.ElementsMatch(t,
		nullify.Fill(perps),                        //nolint:staticcheck
		nullify.Fill(keeper.GetAllPerpetuals(ctx)), //nolint:staticcheck
	)
}

func TestGetAllPerpetuals_Sorted(t *testing.T) {
	// Setup context and keepers
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals
	keepertest.CreateTestLiquidityTiers(t, ctx, keeper)

	perps := []types.Perpetual{
		*perptest.GeneratePerpetual(perptest.WithId(999)),
		*perptest.GeneratePerpetual(perptest.WithId(5)),
		*perptest.GeneratePerpetual(perptest.WithId(0)),
		*perptest.GeneratePerpetual(perptest.WithId(20)),
		*perptest.GeneratePerpetual(perptest.WithId(1)),
	}

	_, err := pricesKeeper.CreateMarket(
		ctx,
		// `ExchangeConfigJson` is left unset as it is not used by the server.
		pricestypes.MarketParam{
			Id:                0,
			Pair:              "marketName",
			Exponent:          -10,
			MinExchanges:      uint32(1),
			MinPriceChangePpm: uint32(50),
		},
		pricestypes.MarketPrice{
			Id:       0,
			Exponent: -10,
			Price:    1_000, // leave this as a placeholder b/c we cannot set the price to 0
		},
	)
	require.NoError(t, err)

	for perp := range perps {
		_, err := keeper.CreatePerpetual(
			ctx,
			perps[perp].Params.Id,
			perps[perp].Params.Ticker,
			perps[perp].Params.MarketId,
			perps[perp].Params.AtomicResolution,
			perps[perp].Params.DefaultFundingPpm,
			perps[perp].Params.LiquidityTier,
		)
		require.NoError(t, err)
	}

	got := keeper.GetAllPerpetuals(ctx)
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

func TestGetMarginRequirements_Success(t *testing.T) {
	oneBip := math.Pow10(2)
	oneTrillion := 1_000_000_000_000
	tests := map[string]struct {
		price                           uint64
		exponent                        int32
		baseCurrencyAtomicResolution    int32
		bigBaseQuantums                 *big.Int
		initialMarginPpm                uint32
		maintenanceFractionPpm          uint32
		basePositionNotional            uint64
		bigExpectedInitialMarginPpm     *big.Int
		bigExpectedMaintenanceMarginPpm *big.Int
	}{
		"InitialMargin 2 BIPs, MaintenanceMargin 1 BIP, positive exponent, atomic resolution 8": {
			price:                           5_555,
			exponent:                        2,
			baseCurrencyAtomicResolution:    -8,
			bigBaseQuantums:                 big.NewInt(7_000),
			initialMarginPpm:                uint32(oneBip * 2),
			maintenanceFractionPpm:          uint32(500_000), // 50% of IM
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(7_777),
			bigExpectedMaintenanceMarginPpm: big.NewInt(3_889),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 50 BIPs, atomic resolution 4": {
			price:                           5_555,
			exponent:                        0,
			baseCurrencyAtomicResolution:    -4,
			bigBaseQuantums:                 big.NewInt(7_000),
			initialMarginPpm:                uint32(oneBip * 100),
			maintenanceFractionPpm:          uint32(500_000), // 50% of IM
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(38_885_000),
			bigExpectedMaintenanceMarginPpm: big.NewInt(19_442_500),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 50 BIPs, positive exponent, atomic resolution 0": {
			price:                           42,
			exponent:                        5,
			baseCurrencyAtomicResolution:    -0,
			bigBaseQuantums:                 big.NewInt(88),
			initialMarginPpm:                uint32(oneBip * 100),
			maintenanceFractionPpm:          uint32(500_000),             // 50% of IM
			basePositionNotional:            uint64(369_600_000_000_000), // same as quote quantums
			bigExpectedInitialMarginPpm:     big.NewInt(3_696_000_000_000),
			bigExpectedMaintenanceMarginPpm: big.NewInt(1_848_000_000_000),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 50 BIPs, negative exponent, atomic resolution 6": {
			price:                           42_000_000,
			exponent:                        -2,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(-5_000),
			initialMarginPpm:                uint32(oneBip * 100),
			maintenanceFractionPpm:          uint32(500_000), // 50% of IM
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(21_000_000),
			bigExpectedMaintenanceMarginPpm: big.NewInt(10_500_000),
		},
		"InitialMargin 10_000 BIPs (max), MaintenanceMargin 10_000 BIPs (max), atomic resolution 6": {
			price:                           5_555,
			exponent:                        0,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(7_000),
			initialMarginPpm:                uint32(oneBip * 10_000),
			maintenanceFractionPpm:          uint32(1_000_000), // 100% of IM
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(38_885_000),
			bigExpectedMaintenanceMarginPpm: big.NewInt(38_885_000),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 100 BIPs, atomic resolution 6": {
			price:                           5_555,
			exponent:                        0,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(7_000),
			initialMarginPpm:                uint32(oneBip * 100),
			maintenanceFractionPpm:          uint32(1_000_000), // 100% of IM
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(388_850),
			bigExpectedMaintenanceMarginPpm: big.NewInt(388_850),
		},
		"InitialMargin 0.02 BIPs, MaintenanceMargin 0.01 BIPs, positive exponent, atomic resolution 6": {
			price:                           5_555,
			exponent:                        3,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(-7_000),
			initialMarginPpm:                uint32(oneBip * 0.02),
			maintenanceFractionPpm:          uint32(500_000), // 50% of IM
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(77_770),
			bigExpectedMaintenanceMarginPpm: big.NewInt(38_885),
		},
		"InitialMargin 0 BIPs (min), MaintenanceMargin 0 BIPs (min), atomic resolution 6": {
			price:                           5_555,
			exponent:                        0,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(7_000),
			initialMarginPpm:                uint32(oneBip * 0),
			maintenanceFractionPpm:          uint32(1_000_000), // 100% of IM,
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(0),
			bigExpectedMaintenanceMarginPpm: big.NewInt(0),
		},
		"Price is zero, atomic resolution 6": {
			price:                           0,
			exponent:                        1,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(-7_000),
			initialMarginPpm:                uint32(oneBip * 1),
			maintenanceFractionPpm:          uint32(1_000_000), // 100% of IM,
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(0),
			bigExpectedMaintenanceMarginPpm: big.NewInt(0),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              new(big.Int).SetUint64(math.MaxUint64),
			initialMarginPpm:             uint32(oneBip * 1),
			maintenanceFractionPpm:       uint32(1_000_000), // 100% of IM,
			basePositionNotional:         math.MaxUint64,
			// As both `price` and `bigBaseQuantums` are `MaxUint64`, `bigQuoteQuantums` (`= price * bigBaseQuantums`)
			// has a value much higher than `MaxUint64` (3402823669209384634264811192843491082250).
			// Now that `bigQuoteQuantums` has a much higher value than `basePositionNotional`, which is
			// only the max value of a `uint64`, `marginAdjustmentPpm` is a very big value (13_581_879_131_294_591).
			// Thus, adjusted initial margin (initial margin * margin adjustment) is capped at 100%,
			// so adjusted initial margin in quote quantums = `100% * bigQuoteQuantums` = `bigQuoteQuantums`.
			bigExpectedInitialMarginPpm: big_testutil.MustFirst(
				new(big.Int).SetString("3402823669209384634264811192843491082250", 10),
			),
			bigExpectedMaintenanceMarginPpm: big_testutil.MustFirst(
				new(big.Int).SetString("3402823669209384634264811192843491082250", 10),
			),
		},
		"InitialMargin 100 BIPs, MaintenanceMargin 50 BIPs, atomic resolution 6, margin adjusted": {
			price:                        5_555,
			exponent:                     0,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(7_000),
			initialMarginPpm:             uint32(oneBip * 100),
			maintenanceFractionPpm:       uint32(500_000), // 50% of IM
			basePositionNotional:         uint64(1_000_000),
			// marginAdjustmentPpm = sqrt(quoteQuantums * (OneMillion * OneMillion) / basePositionNotional)
			// = sqrt(38_885_000 * 1_000_000) ~= 6235783
			// initialMarginPpmQuoteQuantums = initialMarginPpm * quoteQuantums * marginAdjustmentPpm / 1_000_000 / 1_000_000
			// = 10_000 * 38_885_000 * 6235783 / 1_000_000 / 1_000_000 ~= 2_424_784
			bigExpectedInitialMarginPpm:     big.NewInt(2_424_784),
			bigExpectedMaintenanceMarginPpm: big.NewInt(1_212_392),
		},
		"InitialMargin 20%, MaintenanceMargin 10%, atomic resolution 6, margin adjusted": {
			price:                        36_750,
			exponent:                     0,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(12_000),
			initialMarginPpm:             uint32(200_000),
			maintenanceFractionPpm:       uint32(500_000), // 50% of IM
			basePositionNotional:         uint64(100_000_000),
			// quoteQuantums = 36_750 * 12_000 = 441_000_000
			// marginAdjustmentPpm = sqrt(quoteQuantums * (OneMillion * OneMillion) / basePositionNotional)
			// = sqrt(441_000_000 * (OneMillion * OneMillion) / 100_000_000) ~= 2_100_000
			// initialMarginPpmQuoteQuantums = initialMarginPpm * quoteQuantums * marginAdjustmentPpm / 1_000_000 / 1_000_000
			// = 200_000 * 441_000_000 * 2_100_000 / 1_000_000 / 1_000_000 ~= 185_220_000
			bigExpectedInitialMarginPpm:     big.NewInt(185_220_000),
			bigExpectedMaintenanceMarginPpm: big.NewInt(92_610_000),
		},
		"InitialMargin 5%, MaintenanceMargin 3%, atomic resolution 6, margin adjusted": {
			price:                        123_456,
			exponent:                     0,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(74_523),
			initialMarginPpm:             uint32(50_000),
			maintenanceFractionPpm:       uint32(600_000), // 60% of IM
			basePositionNotional:         uint64(100_000_000),
			// quoteQuantums = 123_456 * 74_523 = 9_200_311_488
			// marginAdjustmentPpm = sqrt(quoteQuantums * (OneMillion * OneMillion) / basePositionNotional)
			// = sqrt(9_200_311_488 * (OneMillion * OneMillion) / 100_000_000) ~= 9_591_825
			// initialMarginPpmQuoteQuantums = initialMarginPpm * quoteQuantums * marginAdjustmentPpm / 1_000_000 / 1_000_000
			// = 50_000 * 9_200_311_488 * 9_591_825 / 1_000_000 / 1_000_000 ~= 4_412_388_886
			bigExpectedInitialMarginPpm:     big.NewInt(4_412_388_886),
			bigExpectedMaintenanceMarginPpm: big.NewInt(2_647_433_332),
		},
		"InitialMargin 25%, MaintenanceMargin 15%, atomic resolution 6, margin adjusted and IM capped at 100% of notional": {
			price:                        123_456,
			exponent:                     0,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              big.NewInt(74_523),
			initialMarginPpm:             uint32(250_000),
			maintenanceFractionPpm:       uint32(600_000), // 60% of IM
			basePositionNotional:         uint64(100_000_000),
			// quoteQuantums = 123_456 * 74_523 = 9_200_311_488
			// marginAdjustmentPpm = sqrt(quoteQuantums * (OneMillion * OneMillion) / basePositionNotional)
			// = sqrt(9_200_311_488 * (OneMillion * OneMillion) / 100_000_000) ~= 9_591_825
			// After adjustment, initial margin is capped at 100% of notional (quote quantums).
			bigExpectedInitialMarginPpm:     big.NewInt(9_200_311_488),
			bigExpectedMaintenanceMarginPpm: big.NewInt(5_520_186_893),
		},
		"InitialMargin 10_000 BIPs (max), MaintenanceMargin 10_000 BIPs (max), atomic resolution 6, margin adjusted": {
			price:                           5_555,
			exponent:                        0,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(7_000),
			initialMarginPpm:                uint32(oneBip * 10_000),
			maintenanceFractionPpm:          uint32(1_000_000), // 100% of IM
			basePositionNotional:            uint64(oneTrillion),
			bigExpectedInitialMarginPpm:     big.NewInt(38_885_000),
			bigExpectedMaintenanceMarginPpm: big.NewInt(38_885_000),
		},
		"InitialMargin 0 BIPs (min), MaintenanceMargin 0 BIPs (min), atomic resolution 6, margin adjusted": {
			price:                           5_555,
			exponent:                        0,
			baseCurrencyAtomicResolution:    -6,
			bigBaseQuantums:                 big.NewInt(7_000),
			initialMarginPpm:                uint32(oneBip * 0),
			maintenanceFractionPpm:          uint32(1_000_000), // 100% of IM,
			basePositionNotional:            uint64(1_000_000),
			bigExpectedInitialMarginPpm:     big.NewInt(0),
			bigExpectedMaintenanceMarginPpm: big.NewInt(0),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Individual test setup.
			ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
			// Create a new market param and price.
			marketId := pricesKeeper.GetNumMarkets(ctx)
			_, err := pricesKeeper.CreateMarket(
				ctx,
				// `ExchangeConfigJson` is left unset as it is not used by the server.
				pricestypes.MarketParam{
					Id:                marketId,
					Pair:              "marketName",
					Exponent:          tc.exponent,
					MinExchanges:      uint32(1),
					MinPriceChangePpm: uint32(50),
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    1_000, // leave this as a placeholder b/c we cannot set the price to 0
				},
			)
			require.NoError(t, err)

			// Update `Market.price`. By updating prices this way, we can simulate conditions where the oracle
			// price may become 0.
			err = pricesKeeper.UpdateMarketPrices(
				ctx,
				[]*pricestypes.MsgUpdateMarketPrices_MarketPrice{pricestypes.NewMarketPriceUpdate(
					marketId,
					tc.price,
				)},
			)
			require.NoError(t, err)

			// Create `LiquidityTier` struct.
			_, err = keeper.CreateLiquidityTier(
				ctx,
				"name",
				tc.initialMarginPpm,
				tc.maintenanceFractionPpm,
				tc.basePositionNotional,
				1, // dummy impact notional value
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := keeper.CreatePerpetual(
				ctx,
				0,                               // PerpetualId
				"getMarginRequirementsTicker",   // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
			)
			require.NoError(t, err)

			// Verify initial and maintenance margin requirements are calculated correctly.
			bigInitialMargin, bigMaintenanceMargin, err := keeper.GetMarginRequirements(
				ctx,
				perpetual.Params.Id,
				tc.bigBaseQuantums,
			)
			require.NoError(t, err)

			if tc.bigExpectedInitialMarginPpm.Cmp(bigInitialMargin) != 0 {
				t.Fatalf(
					"%s: expectedInitialMargin: %s, initialMargin: %s",
					name,
					tc.bigExpectedInitialMarginPpm.String(),
					bigInitialMargin.String())
			}

			if tc.bigExpectedMaintenanceMarginPpm.Cmp(bigMaintenanceMargin) != 0 {
				t.Fatalf(
					"%s: expectedMaintenanceMargin: %s, maintenanceMargin: %s",
					name,
					tc.bigExpectedMaintenanceMarginPpm.String(),
					bigMaintenanceMargin.String())
			}
		})
	}
}

func TestGetMarginRequirements_PerpetualNotFound(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, _, err := keeper.GetMarginRequirements(
		ctx,
		nonExistentPerpetualId,
		big.NewInt(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetMarginRequirements_MarketNotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, storeKey := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.PerpetualKeyPrefix))
	perpetualStore.Set(types.PerpetualKey(
		perpetual.Params.Id,
	), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, _, err := keeper.GetMarginRequirements(
		ctx,
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

func TestGetMarginRequirements_LiquidityTierNotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, storeKey := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad LiquidityTier.
	nonExistentLiquidityTier := uint32(999)
	perpetual.Params.LiquidityTier = nonExistentLiquidityTier
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.PerpetualKeyPrefix))
	perpetualStore.Set(types.PerpetualKey(
		perpetual.Params.Id,
	), b)

	// Getting margin requirements for perpetual with bad LiquidityTier should return an error.
	_, _, err := keeper.GetMarginRequirements(
		ctx,
		perpetual.Params.Id,
		big.NewInt(-1),
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
		"Positive exponent, atomic resolution 6, long position": {
			price:                               5_555,
			exponent:                            2,
			baseCurrencyAtomicResolution:        -6,
			bigBaseQuantums:                     big.NewInt(7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(3_888_500_000),
		},
		"Positive exponent, atomic resolution 6, short position": {
			price:                               5_555,
			exponent:                            2,
			baseCurrencyAtomicResolution:        -6,
			bigBaseQuantums:                     big.NewInt(-7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(-3_888_500_000),
		},
		"Negative exponent, atomic resolution 6, short position": {
			price:                               5_555,
			exponent:                            -2,
			baseCurrencyAtomicResolution:        -6,
			bigBaseQuantums:                     big.NewInt(-7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(-388_850),
		},
		"Zero exponent, atomic resolution 6, short position": {
			price:                               5_555,
			exponent:                            0,
			baseCurrencyAtomicResolution:        -6,
			bigBaseQuantums:                     big.NewInt(-7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(-38_885_000),
		},
		"Positive exponent, atomic resolution 4, long position": {
			price:                               5_555,
			exponent:                            4,
			baseCurrencyAtomicResolution:        -4,
			bigBaseQuantums:                     big.NewInt(7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(38_885_000_000_000),
		},
		"Positive exponent, atomic resolution 0, long position": {
			price:                               5_555,
			exponent:                            4,
			baseCurrencyAtomicResolution:        -0,
			bigBaseQuantums:                     big.NewInt(7_000),
			bigExpectedNetNotionalQuoteQuantums: big.NewInt(388_850_000_000_000_000),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     1,
			baseCurrencyAtomicResolution: -6,
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
			ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, keeper)

			// Create a new market param and price.
			marketId := uint32(0)
			_, err := pricesKeeper.CreateMarket(
				ctx,
				pricestypes.MarketParam{
					Id:                marketId,
					Pair:              "marketName",
					Exponent:          tc.exponent,
					MinExchanges:      uint32(1),
					MinPriceChangePpm: uint32(50),
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    tc.price,
				},
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := keeper.CreatePerpetual(
				ctx,
				0,                               // PerpetualId
				"GetNetNotionalTicker",          // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
			)
			require.NoError(t, err)

			// Verify collateral requirements are calculated correctly.
			bigNotionalQuoteQuantums, err := keeper.GetNetNotional(
				ctx,
				perpetual.Params.Id,
				tc.bigBaseQuantums,
			)
			require.NoError(t, err)

			if tc.bigExpectedNetNotionalQuoteQuantums.Cmp(bigNotionalQuoteQuantums) != 0 {
				t.Fatalf(
					"%s: expectedNetNotionalQuoteQuantums: %s, collateralQuoteQuantums: %s",
					name,
					tc.bigExpectedNetNotionalQuoteQuantums.String(),
					bigNotionalQuoteQuantums.String(),
				)
			}
		})
	}
}

func TestGetNetNotional_PerpetualNotFound(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := keeper.GetNetNotional(
		ctx,
		nonExistentPerpetualId,
		big.NewInt(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetNetNotional_MarketNotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, storeKey := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.PerpetualKeyPrefix))
	perpetualStore.Set(types.PerpetualKey(
		perpetual.Params.Id,
	), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, err := keeper.GetNetNotional(
		ctx,
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
		"Positive exponent, atomic resolution 6, long position": {
			price:                              5_555,
			exponent:                           2,
			baseCurrencyAtomicResolution:       -6,
			bigQuoteQuantums:                   big.NewInt(3_888_500_000),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(7_000),
		},
		"Positive exponent, atomic resolution 6, short position": {
			price:                              5_555,
			exponent:                           2,
			baseCurrencyAtomicResolution:       -6,
			bigQuoteQuantums:                   big.NewInt(-3_888_500_000),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(-7_000),
		},
		"Negative exponent, atomic resolution 6, short position": {
			price:                              5_555,
			exponent:                           -2,
			baseCurrencyAtomicResolution:       -6,
			bigQuoteQuantums:                   big.NewInt(-388_850),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(-7_000),
		},
		"Zero exponent, atomic resolution 6, short position": {
			price:                              5_555,
			exponent:                           0,
			baseCurrencyAtomicResolution:       -6,
			bigQuoteQuantums:                   big.NewInt(-38_885_000),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(-7_000),
		},
		"Positive exponent, atomic resolution 4, long position": {
			price:                              5_555,
			exponent:                           4,
			baseCurrencyAtomicResolution:       -4,
			bigQuoteQuantums:                   big.NewInt(38_885_000_000_000),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(7_000),
		},
		"Positive exponent, atomic resolution 0, long position": {
			price:                              5_555,
			exponent:                           4,
			baseCurrencyAtomicResolution:       -0,
			bigQuoteQuantums:                   big.NewInt(388_850_000_000_000_000),
			bigExpectedNetNotionalBaseQuantums: big.NewInt(7_000),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     1,
			baseCurrencyAtomicResolution: -6,
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
			ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, keeper)
			// Create a new market param and price.
			marketId := pricesKeeper.GetNumMarkets(ctx)
			_, err := pricesKeeper.CreateMarket(
				ctx,
				pricestypes.MarketParam{
					Id:                marketId,
					Pair:              "marketName",
					Exponent:          tc.exponent,
					MinExchanges:      uint32(1),
					MinPriceChangePpm: uint32(50),
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    tc.price,
				},
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := keeper.CreatePerpetual(
				ctx,
				0,                               // PerpetualId
				"GetNetNotionalTicker",          // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
			)
			require.NoError(t, err)

			// Verify collateral requirements are calculated correctly.
			bigNotionalBaseQuantums, err := keeper.GetNotionalInBaseQuantums(
				ctx,
				perpetual.Params.Id,
				tc.bigQuoteQuantums,
			)
			require.NoError(t, err)

			if tc.bigExpectedNetNotionalBaseQuantums.Cmp(bigNotionalBaseQuantums) != 0 {
				t.Fatalf(
					"%s: expectedNetNotionalBaseQuantums: %s, collateralBaseQuantums: %s",
					name,
					tc.bigExpectedNetNotionalBaseQuantums.String(),
					bigNotionalBaseQuantums.String(),
				)
			}
		})
	}
}

func TestGetNotionalInBaseQuantums_PerpetualNotFound(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := keeper.GetNotionalInBaseQuantums(
		ctx,
		nonExistentPerpetualId,
		big.NewInt(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetNotionalInBaseQuantums_MarketNotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, storeKey := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.PerpetualKeyPrefix))
	perpetualStore.Set(types.PerpetualKey(
		perpetual.Params.Id,
	), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, err := keeper.GetNotionalInBaseQuantums(
		ctx,
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
		"Positive exponent, atomic resolution 6, long position": {
			price:                                 5_555,
			exponent:                              2,
			baseCurrencyAtomicResolution:          -6,
			bigBaseQuantums:                       big.NewInt(7_000),
			bigExpectedNetCollateralQuoteQuantums: big.NewInt(3_888_500_000),
		},
		"Positive exponent, atomic resolution 6, short position": {
			price:                                 5_555,
			exponent:                              2,
			baseCurrencyAtomicResolution:          -6,
			bigBaseQuantums:                       big.NewInt(-7_000),
			bigExpectedNetCollateralQuoteQuantums: big.NewInt(-3_888_500_000),
		},
		"Negative exponent, atomic resolution 6, short position": {
			price:                                 5_555,
			exponent:                              -2,
			baseCurrencyAtomicResolution:          -6,
			bigBaseQuantums:                       big.NewInt(-7_000),
			bigExpectedNetCollateralQuoteQuantums: big.NewInt(-388_850),
		},
		"Zero exponent, atomic resolution 6, short position": {
			price:                                 5_555,
			exponent:                              0,
			baseCurrencyAtomicResolution:          -6,
			bigBaseQuantums:                       big.NewInt(-7_000),
			bigExpectedNetCollateralQuoteQuantums: big.NewInt(-38_885_000),
		},
		"Positive exponent, atomic resolution 4, long position": {
			price:                                 5_555,
			exponent:                              4,
			baseCurrencyAtomicResolution:          -4,
			bigBaseQuantums:                       big.NewInt(7_000),
			bigExpectedNetCollateralQuoteQuantums: big.NewInt(38_885_000_000_000),
		},
		"Positive exponent, atomic resolution 0, long position": {
			price:                                 5_555,
			exponent:                              4,
			baseCurrencyAtomicResolution:          -0,
			bigBaseQuantums:                       big.NewInt(7_000),
			bigExpectedNetCollateralQuoteQuantums: big.NewInt(388_850_000_000_000_000),
		},
		"Price and quantums are max uints": {
			price:                        math.MaxUint64,
			exponent:                     1,
			baseCurrencyAtomicResolution: -6,
			bigBaseQuantums:              new(big.Int).SetUint64(math.MaxUint64),
			bigExpectedNetCollateralQuoteQuantums: big_testutil.MustFirst(
				new(big.Int).SetString("3402823669209384634264811192843491082250", 10),
			),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Test suite setup.
			ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, keeper)
			// Test setup.
			// Create a new market.
			marketId := pricesKeeper.GetNumMarkets(ctx)
			_, err := pricesKeeper.CreateMarket(
				ctx,
				pricestypes.MarketParam{
					Id:                marketId,
					Pair:              "marketName",
					Exponent:          tc.exponent,
					MinExchanges:      uint32(1),
					MinPriceChangePpm: uint32(50),
				},
				pricestypes.MarketPrice{
					Id:       marketId,
					Exponent: tc.exponent,
					Price:    tc.price,
				},
			)
			require.NoError(t, err)

			// Create `Perpetual` struct with baseAssetAtomicResolution and marketId.
			perpetual, err := keeper.CreatePerpetual(
				ctx,
				0,                               // PerpetualId
				"GetNetCollateralTicker",        // Ticker
				marketId,                        // MarketId
				tc.baseCurrencyAtomicResolution, // AtomicResolution
				int32(0),                        // DefaultFundingPpm
				0,                               // LiquidityTier
			)
			require.NoError(t, err)

			// Verify collateral requirements are calculated correctly.
			bigCollateralQuoteQuantums, err := keeper.GetNetCollateral(
				ctx,
				perpetual.Params.Id,
				tc.bigBaseQuantums,
			)
			require.NoError(t, err)

			if tc.bigExpectedNetCollateralQuoteQuantums.Cmp(bigCollateralQuoteQuantums) != 0 {
				t.Fatalf(
					"%s: expectedNetCollateralQuoteQuantums: %s, collateralQuoteQuantums: %s",
					name,
					tc.bigExpectedNetCollateralQuoteQuantums.String(),
					bigCollateralQuoteQuantums.String(),
				)
			}
		})
	}
}

func TestGetNetCollateral_PerpetualNotFound(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, err := keeper.GetNetCollateral(
		ctx,
		nonExistentPerpetualId,
		big.NewInt(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestGetNetCollateral_MarketNotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, storeKey := keepertest.PerpetualsKeepers(t)

	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 1)
	perpetual := perps[0]

	// Store the perpetual with a bad MarketId.
	nonExistentMarketId := uint32(999)
	perpetual.Params.MarketId = nonExistentMarketId
	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	b := cdc.MustMarshal(&perpetual)
	perpetualStore := prefix.NewStore(ctx.KVStore(storeKey), types.KeyPrefix(types.PerpetualKeyPrefix))
	perpetualStore.Set(types.PerpetualKey(
		perpetual.Params.Id,
	), b)

	// Getting margin requirements for perpetual with bad MarketId should return an error.
	_, err := keeper.GetNetCollateral(
		ctx,
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

func TestGetSettlement_Success(t *testing.T) {
	tests := map[string]struct {
		quantums              *big.Int
		perpetualFundingIndex *big.Int
		prevIndex             *big.Int
		expectedSettlement    *big.Int
	}{
		"is long, index went from negative to positive": {
			quantums:              big.NewInt(30_000),
			prevIndex:             big.NewInt(-100),
			perpetualFundingIndex: big.NewInt(100),
			expectedSettlement:    big.NewInt(-6),
		},
		"is long, index unchanged": {
			quantums:              big.NewInt(1_000_000),
			prevIndex:             big.NewInt(100),
			perpetualFundingIndex: big.NewInt(100),
			expectedSettlement:    big.NewInt(0),
		},
		"is long, index went from positive to zero": {
			quantums:              big.NewInt(10_000_000),
			prevIndex:             big.NewInt(100),
			perpetualFundingIndex: big.NewInt(0),
			expectedSettlement:    big.NewInt(1_000),
		},
		"is long, index went from positive to negative": {
			quantums:              big.NewInt(10_000_000),
			prevIndex:             big.NewInt(100),
			perpetualFundingIndex: big.NewInt(-200),
			expectedSettlement:    big.NewInt(3_000),
		},
		"is short, index went from negative to positive": {
			quantums:              big.NewInt(-30_000),
			prevIndex:             big.NewInt(-100),
			perpetualFundingIndex: big.NewInt(100),
			expectedSettlement:    big.NewInt(6),
		},
		"is short, index unchanged": {
			quantums:              big.NewInt(-1_000),
			prevIndex:             big.NewInt(100),
			perpetualFundingIndex: big.NewInt(100),
			expectedSettlement:    big.NewInt(0),
		},
		"is short, index went from positive to zero": {
			quantums:              big.NewInt(-5_000_000),
			prevIndex:             big.NewInt(100),
			perpetualFundingIndex: big.NewInt(0),
			expectedSettlement:    big.NewInt(-500),
		},
		"is short, index went from positive to negative": {
			quantums:              big.NewInt(-5_000_000),
			prevIndex:             big.NewInt(100),
			perpetualFundingIndex: big.NewInt(-50),
			expectedSettlement:    big.NewInt(-750),
		},
		"rounding - negative settlement should round toward negative infinity (is short)": {
			quantums:              big.NewInt(-5_500_000),
			prevIndex:             big.NewInt(1),
			perpetualFundingIndex: big.NewInt(0),
			expectedSettlement:    big.NewInt(-6),
		},
		"rounding - negative settlement should round toward negative infinity (is long)": {
			quantums:              big.NewInt(5_500_000),
			prevIndex:             big.NewInt(0),
			perpetualFundingIndex: big.NewInt(1),
			expectedSettlement:    big.NewInt(-6),
		},
		"rounding - positive settlement should round toward zero (is short)": {
			quantums:              big.NewInt(-5_500_000),
			prevIndex:             big.NewInt(0),
			perpetualFundingIndex: big.NewInt(1),
			expectedSettlement:    big.NewInt(5),
		},
		"rounding - positive settlement should round toward zero (is long)": {
			quantums:              big.NewInt(5_500_000),
			prevIndex:             big.NewInt(1),
			perpetualFundingIndex: big.NewInt(0),
			expectedSettlement:    big.NewInt(5),
		},
	}

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Test suite setup.
			ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, keeper)

			perps, err := keepertest.CreateNPerpetuals(t, ctx, keeper, pricesKeeper, 1)
			require.NoError(t, err)

			perpetualId := perps[0].Params.Id

			// Since FundingIndex starts at zero, tc.perpetualFundingIndex will be
			// the current FundingIndex.
			err = keeper.ModifyFundingIndex(ctx, perpetualId, tc.perpetualFundingIndex)
			require.NoError(t, err)

			bigNetSettlement, newFundingIndex, err := keeper.GetSettlement(
				ctx,
				perpetualId,
				tc.quantums,
				tc.prevIndex,
			)
			require.NoError(t, err)

			require.Equal(t,
				tc.perpetualFundingIndex,
				newFundingIndex,
			)

			require.Equal(t,
				0,
				tc.expectedSettlement.Cmp(bigNetSettlement),
			)
		})
	}
}

func TestGetSettlement_PerpetualNotFound(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)
	_, _, err := keeper.GetSettlement(
		ctx,
		nonExistentPerpetualId, // perpetualId
		big.NewInt(-100),       // quantum
		big.NewInt(0),          // index
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrPerpetualDoesNotExist, fmt.Sprint(nonExistentPerpetualId)).Error())
	require.ErrorIs(t, err, types.ErrPerpetualDoesNotExist)
}

func TestModifyFundingIndex_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 100)

	for _, perp := range perps {
		testFundingIndexDelta := big.NewInt(2*(int64(perp.Params.Id)%2) - 1)

		err := keeper.ModifyFundingIndex(
			ctx,
			perp.Params.Id,
			testFundingIndexDelta,
		)
		require.NoError(t, err)

		newPerp, err := keeper.GetPerpetual(ctx, perp.Params.Id)
		require.NoError(t, err)

		require.Equal(
			t,
			testFundingIndexDelta,
			newPerp.FundingIndex.BigInt(),
		)
	}
}

func TestModifyFundingIndex_PerpetualDoesNotExist(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	nonExistentPerpetualId := uint32(0)

	err := keeper.ModifyFundingIndex(
		ctx,
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
			ctx, perpsKeeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
			// Create liquidity tiers and perpetuals,
			_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, perpsKeeper, pricesKeeper, 1)

			// Set up intitial funding index, should succeed.
			err := perpsKeeper.ModifyFundingIndex(
				ctx,
				tc.perpetualId,
				tc.initialFundingIndex,
			)
			require.NoError(t, err)

			err = perpsKeeper.ModifyFundingIndex(
				ctx,
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
			ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)

			sampleTailsRemovalFunc := keeper.GetRemoveSampleTailsFunc(ctx, tc.removalRatePpm)
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
			ctx, perpsKeeper, pricesKeeper, epochsKeeper, _ := keepertest.PerpetualsKeepers(t)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Create liquidity tiers.
			keepertest.CreateTestLiquidityTiers(t, ctx, perpsKeeper)

			// Create test perpetuals.
			// 1BTC = $50,000.
			oldPerps := make([]types.Perpetual, len(tc.testPerpetuals))
			for i, p := range tc.testPerpetuals {
				perp, err := perpsKeeper.CreatePerpetual(
					ctx,
					p.Params.Id,
					p.Params.Ticker,
					p.Params.MarketId,
					p.Params.AtomicResolution,
					p.Params.DefaultFundingPpm,
					p.Params.LiquidityTier,
				)
				require.NoError(t, err)
				oldPerps[i] = perp
			}

			// Create funding-tick epoch.
			err := epochsKeeper.CreateEpochInfo(
				ctx,
				epochstypes.EpochInfo{
					Name:                   string(epochstypes.FundingTickEpochInfoName),
					Duration:               tc.testFundingTickDuration,
					CurrentEpochStartBlock: testCurrentFundingTickEpochStartBlock,
					CurrentEpoch:           testCurrentEpoch,
				},
			)
			require.NoError(t, err)
			// Create funding-sample epoch.
			err = epochsKeeper.CreateEpochInfo(
				ctx,
				epochstypes.EpochInfo{
					Name:     string(epochstypes.FundingSampleEpochInfoName),
					Duration: tc.testFundingSampleDuration,
				},
			)
			require.NoError(t, err)

			// Insert test funding sample.
			keepertest.PopulateTestPremiumStore(
				t,
				ctx,
				perpsKeeper,
				oldPerps,
				tc.testFundingSamples,
				false, // isVote
			)

			perpsKeeper.MaybeProcessNewFundingTickEpoch(
				// Current block is the start of a new epoch for funding-tick.
				ctx.WithBlockHeight(int64(testCurrentFundingTickEpochStartBlock)))

			for i, p := range oldPerps {
				newPerp, err := perpsKeeper.GetPerpetual(ctx, p.Params.Id)
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
				allSamples := perpsKeeper.GetPremiumSamples(ctx)
				require.NoError(t, err)
				for _, marketPremiums := range allSamples.AllMarketPremiums {
					require.Equal(t, 0, len(marketPremiums.Premiums))
				}
			}

			fundingEvents := getFundingBlockEventsFromIndexerBlock(ctx, perpsKeeper)
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
	perpsKeeper *keeper.Keeper,
) []*indexerevents.FundingEventV1 {
	block := perpsKeeper.GetIndexerEventManager().ProduceBlock(ctx)
	var fundingEvents []*indexerevents.FundingEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeFundingValues {
			continue
		}
		if _, ok := event.OrderingWithinBlock.(*indexer_manager.IndexerTendermintEvent_BlockEvent_); ok {
			bytes := indexer_manager.GetBytesFromEventData(event.Data)
			unmarshaler := common.UnmarshalerImpl{}
			var fundingEvent indexerevents.FundingEventV1
			err := unmarshaler.Unmarshal(bytes, &fundingEvent)
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
			ctx, perpsKeeper, pricesKeeper, epochsKeeper, _ := keepertest.PerpetualsKeepers(t)

			// Create the default markets.
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			// Insert test funding sample.
			keepertest.PopulateTestPremiumStore(
				t,
				ctx,
				perpsKeeper,
				tc.testPerpetuals,
				tc.testPremiumSamples,
				false, // isVote
			)

			// Create test epochs.
			for _, epochInfo := range tc.testEpochs {
				err := epochsKeeper.CreateEpochInfo(
					ctx,
					epochInfo,
				)
				require.NoError(t, err)
			}

			initialEvents := ctx.EventManager().ABCIEvents()

			require.PanicsWithError(
				t,
				tc.expectedError.Error(),
				func() {
					perpsKeeper.MaybeProcessNewFundingTickEpoch(
						ctx.WithBlockHeight(int64(testCurrentFundingTickEpochStartBlock)))
				},
			)

			// Verify that no new events were emitted.
			laterEvents := ctx.EventManager().ABCIEvents()
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

	ctx, perpsKeeper, pricesKeeper, epochsKeeper, _ := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers and perpetuals,
	perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, perpsKeeper, pricesKeeper, 100)

	err := epochsKeeper.CreateEpochInfo(
		ctx,
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

	perpsKeeper.MaybeProcessNewFundingTickEpoch(
		// Current block is not start of a new epoch for funding-tick.
		ctx.WithBlockHeight(int64(testCurrentFundingTickEpochStartBlock + 1)))

	for _, perp := range perps {
		newPerp, err := perpsKeeper.GetPerpetual(ctx, perp.Params.Id)
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

	ctx, perpsKeeper, _, epochsKeeper, _ := keepertest.PerpetualsKeepers(t)
	err := epochsKeeper.CreateEpochInfo(
		ctx,
		epochstypes.EpochInfo{
			Name:                   string(epochstypes.FundingSampleEpochInfoName),
			Duration:               testDuration,
			CurrentEpochStartBlock: testCurrentFundingSampleEpochStartBlock,
			CurrentEpoch:           testCurrentEpoch,
		},
	)
	require.NoError(t, err)

	msgAddPremiumVotes := perpsKeeper.GetAddPremiumVotes(
		ctx.WithBlockHeight(int64(testCurrentFundingSampleEpochStartBlock)),
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
		expectedNumSamples                  int
	}{
		"Positive premium": {
			currentFundingSampleEpochStartBlock: 23,
			blockHeight:                         23,
			samplePremiumPpm:                    100,
			numPerpetuals:                       10,
			expectedNumSamples:                  10,
		},
		"Negative premium": {
			currentFundingSampleEpochStartBlock: 24,
			blockHeight:                         24,
			samplePremiumPpm:                    -150,
			numPerpetuals:                       10,
			expectedNumSamples:                  10,
		},
		"Not start of new funding-sample epoch, still produce samples": {
			currentFundingSampleEpochStartBlock: 24,
			blockHeight:                         25,
			samplePremiumPpm:                    100,
			numPerpetuals:                       10,
			expectedNumSamples:                  10,
		},
		"Zero premiums": {
			currentFundingSampleEpochStartBlock: 24,
			blockHeight:                         24,
			samplePremiumPpm:                    0,
			numPerpetuals:                       10,
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

			ctx,
				perpsKeeper,
				pricesKeeper,
				epochsKeeper,
				_ := keepertest.PerpetualsKeepersWithClobHelpers(
				t,
				&mockPricePremiumGetter,
			)

			// Create liquidity tiers and perpetuals,
			_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, perpsKeeper, pricesKeeper, tc.numPerpetuals)

			err := epochsKeeper.CreateEpochInfo(
				ctx,
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

			msgAddPremiumVotes := perpsKeeper.GetAddPremiumVotes(
				ctx.WithBlockHeight(int64(tc.blockHeight)),
			)

			mockPricePremiumGetter.AssertNumberOfCalls(
				t,
				"GetPricePremiumForPerpetual",
				tc.numPerpetuals,
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
		ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)

		premiumSamples := tc.getPremiumFunc(keeper, ctx)
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
		ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)

		// Create liquidity tiers and perpetuals,
		numPerpetuals := 10
		perps := keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, numPerpetuals)

		// Insert one round of premiums for all perps.
		firstPremiums := make([]types.FundingPremium, numPerpetuals)
		for i, perp := range perps {
			firstPremiums[i] = types.FundingPremium{
				PerpetualId: perp.Params.Id,
				// -1000 for even Ids, +1000 for odd Ids.
				PremiumPpm: 1_000 * (2*(int32(perp.Params.Id)%2) - 1),
			}
		}

		err := tc.addPremiumFunc(keeper, ctx, firstPremiums)
		require.NoError(t, err)

		// Check each perp has expected number of premiums stored after first around of addPremiumFunc().
		firstStoredPremiums := tc.getPremiumFunc(keeper, ctx)

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
			keeper,
			ctx,
			secondPremiums,
		)
		require.NoError(t, err)

		// Check each perp has expected number of premiums stored after second round of addPremiumFunc().
		secondStoredPremiums := tc.getPremiumFunc(keeper, ctx)

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
		ctx, keeper, pricesKeeper, _, _ := keepertest.PerpetualsKeepers(t)
		nonExistentPerpetualId := uint32(1000)

		newPremiums := []types.FundingPremium{
			{
				PerpetualId: nonExistentPerpetualId,
				PremiumPpm:  -1_000,
			},
		}

		// Create liquidity tiers and perpetuals,
		_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 3)

		err := tc.addPremiumFunc(keeper, ctx, newPremiums)
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

func TestModifyOpenInterest_NotImplemented(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	_, err := keeper.ModifyOpenInterest(
		ctx,
		0,
		true,
		0,
	)
	require.ErrorIs(t, err, types.ErrNotImplementedOpenInterest)
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
			ctx, keeper, pricesKeeper, epochsKeeper, _ := keepertest.PerpetualsKeepers(t)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)

			// Create funding-sample epoch.
			err := epochsKeeper.CreateEpochInfo(
				ctx,
				epochstypes.EpochInfo{
					Name:                   string(epochstypes.FundingSampleEpochInfoName),
					Duration:               testDuration,
					CurrentEpochStartBlock: tc.currentEpochStartBlock,
					CurrentEpoch:           testCurrentEpoch,
				},
			)
			require.NoError(t, err)

			// Create liquidity tiers and perpetuals,
			_ = keepertest.CreateLiquidityTiersAndNPerpetuals(t, ctx, keeper, pricesKeeper, 4)
			require.NoError(t, err)

			err = keeper.SetMinNumVotesPerSample(ctx, tc.minNumVotesPerSample)
			require.NoError(t, err)
			keeper.SetPremiumVotes(ctx, tc.premiumVotes)
			keeper.SetPremiumSamples(ctx, tc.prevPremiumSamples)

			initialEvents := ctx.EventManager().ABCIEvents()

			if tc.panicErr != nil {
				require.PanicsWithError(
					t,
					tc.panicErr.Error(),
					func() {
						keeper.MaybeProcessNewFundingSampleEpoch(ctx.WithBlockHeight(tc.currentBlockHeight))
					},
				)

				laterEvents := ctx.EventManager().ABCIEvents()
				require.ElementsMatch(t,
					initialEvents,
					laterEvents,
				)
				return
			}

			keeper.MaybeProcessNewFundingSampleEpoch(ctx.WithBlockHeight(tc.currentBlockHeight))

			require.Equal(t,
				tc.expectedPremiumVotes,
				keeper.GetPremiumVotes(ctx),
			)

			require.Equal(t,
				tc.expectedPremiumSamples,
				keeper.GetPremiumSamples(ctx),
			)
		})
	}
}

func TestCreateLiquidityTier_Success(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	for i, lt := range constants.LiquidityTiers {
		// Create LiquidityTier without error.
		_, err := keeper.CreateLiquidityTier(
			ctx,
			lt.Name,
			lt.InitialMarginPpm,
			lt.MaintenanceFractionPpm,
			lt.BasePositionNotional,
			lt.ImpactNotional,
		)
		require.NoError(t, err)

		// Validate number of LiquidityTiers in store.
		numLiquidityTiers := keeper.GetNumLiquidityTiers(ctx)
		require.Equal(t, uint32(i+1), numLiquidityTiers)

		// Validate fields of LiquidityTier object in store.
		liquidityTier, err := keeper.GetLiquidityTier(ctx, uint32(i))
		require.NoError(t, err)
		require.Equal(t, lt.Id, liquidityTier.Id)
		require.Equal(t, lt.Name, liquidityTier.Name)
		require.Equal(t, lt.InitialMarginPpm, liquidityTier.InitialMarginPpm)
		require.Equal(t, lt.MaintenanceFractionPpm, liquidityTier.MaintenanceFractionPpm)
		require.Equal(t, lt.BasePositionNotional, liquidityTier.BasePositionNotional)
		require.Equal(t, lt.ImpactNotional, liquidityTier.ImpactNotional)
	}
}

func TestCreateLiquidityTier_Failure(t *testing.T) {
	tests := map[string]struct {
		id                     uint32
		name                   string
		initialMarginPpm       uint32
		maintenanceFractionPpm uint32
		basePositionNotional   uint64
		impactNotional         uint64
		expectedError          error
	}{
		"Initial Margin Ppm exceeds maximum": {
			id:                     0,
			name:                   "Large-Cap",
			initialMarginPpm:       lib.OneMillion + 1,
			maintenanceFractionPpm: 500_000,
			basePositionNotional:   uint64(lib.OneMillion),
			impactNotional:         uint64(lib.OneMillion),
			expectedError:          errorsmod.Wrap(types.ErrInitialMarginPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Maintenance Fraction Ppm exceeds maximum": {
			id:                     1,
			name:                   "Medium-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion + 1,
			basePositionNotional:   uint64(lib.OneMillion),
			impactNotional:         uint64(lib.OneMillion),
			expectedError:          errorsmod.Wrap(types.ErrMaintenanceFractionPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Base Position Notional is zero": {
			id:                     1,
			name:                   "Small-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion,
			basePositionNotional:   uint64(0),
			impactNotional:         uint64(lib.OneMillion),
			expectedError:          errorsmod.Wrap(types.ErrBasePositionNotionalIsZero, fmt.Sprint(0)),
		},
		"Impact Notional is zero": {
			id:                     1,
			name:                   "Small-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion,
			basePositionNotional:   uint64(lib.OneMillion),
			impactNotional:         uint64(0),
			expectedError:          errorsmod.Wrap(types.ErrImpactNotionalIsZero, fmt.Sprint(0)),
		},
	}

	// Test setup.
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := keeper.CreateLiquidityTier(
				ctx,
				tc.name,
				tc.initialMarginPpm,
				tc.maintenanceFractionPpm,
				tc.basePositionNotional,
				tc.impactNotional,
			)

			require.Error(t, err)
			require.EqualError(t, err, tc.expectedError.Error())
		})
	}
}

func TestModifyLiquidityTier_Success(t *testing.T) {
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	for _, lt := range constants.LiquidityTiers {
		_, err := keeper.CreateLiquidityTier(
			ctx,
			lt.Name,
			lt.InitialMarginPpm,
			lt.MaintenanceFractionPpm,
			lt.BasePositionNotional,
			lt.ImpactNotional,
		)
		require.NoError(t, err)
	}

	for i, lt := range constants.LiquidityTiers {
		// Modify each field arbitrarily and
		// verify the fields are modified in state.
		name := fmt.Sprintf("foo_%v", i)
		initialMarginPpm := uint32(i * 2)
		maintenanceFractionPpm := uint32(i * 2)
		basePositionNotional := uint64((i + 1) * 1_000_000)
		impactNotional := uint64((i + 1) * 500_000_000)
		modifiedLt, err := keeper.ModifyLiquidityTier(
			ctx,
			lt.Id,
			name,
			initialMarginPpm,
			maintenanceFractionPpm,
			basePositionNotional,
			impactNotional,
		)
		require.NoError(t, err)
		obtainedLt, err := keeper.GetLiquidityTier(ctx, lt.Id)
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
			basePositionNotional,
			obtainedLt.BasePositionNotional,
		)
		require.Equal(
			t,
			impactNotional,
			obtainedLt.ImpactNotional,
		)
	}
	liquidityTierUpsertEvents := keepertest.GetLiquidityTierUpsertEventsFromIndexerBlock(ctx, keeper)
	require.Len(t, liquidityTierUpsertEvents, len(constants.LiquidityTiers)*2)
}

func TestModifyLiquidityTier_Failure(t *testing.T) {
	tests := map[string]struct {
		id                     uint32
		name                   string
		initialMarginPpm       uint32
		maintenanceFractionPpm uint32
		basePositionNotional   uint64
		impactNotional         uint64
		expectedError          error
	}{
		"Initial Margin Ppm exceeds maximum": {
			id:                     0,
			name:                   "Large-Cap",
			initialMarginPpm:       lib.OneMillion + 1,
			maintenanceFractionPpm: 500_000,
			basePositionNotional:   uint64(lib.OneMillion),
			impactNotional:         uint64(lib.OneMillion),
			expectedError:          errorsmod.Wrap(types.ErrInitialMarginPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Maintenance Fraction Ppm exceeds maximum": {
			id:                     1,
			name:                   "Medium-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion + 1,
			basePositionNotional:   uint64(lib.OneMillion),
			impactNotional:         uint64(lib.OneMillion),
			expectedError:          errorsmod.Wrap(types.ErrMaintenanceFractionPpmExceedsMax, fmt.Sprint(lib.OneMillion+1)),
		},
		"Base Position Notional is zero": {
			id:                     1,
			name:                   "Small-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion,
			basePositionNotional:   uint64(0),
			impactNotional:         uint64(lib.OneMillion),
			expectedError:          errorsmod.Wrap(types.ErrBasePositionNotionalIsZero, fmt.Sprint(0)),
		},
		"Impact Notional is zero": {
			id:                     1,
			name:                   "Small-Cap",
			initialMarginPpm:       500_000,
			maintenanceFractionPpm: lib.OneMillion,
			basePositionNotional:   uint64(lib.OneMillion),
			impactNotional:         uint64(0),
			expectedError:          errorsmod.Wrap(types.ErrImpactNotionalIsZero, fmt.Sprint(0)),
		},
	}

	// Test setup.
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)
	// Create liquidity tiers.
	keepertest.CreateTestLiquidityTiers(t, ctx, keeper)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := keeper.ModifyLiquidityTier(
				ctx,
				tc.id,
				tc.name,
				tc.initialMarginPpm,
				tc.maintenanceFractionPpm,
				tc.basePositionNotional,
				tc.impactNotional,
			)

			require.Error(t, err)
			require.EqualError(t, err, tc.expectedError.Error())
		})
	}
}

func TestSetFundingRateClampFactorPpm(t *testing.T) {
	tests := map[string]struct {
		fundingRateClampFactorPpm uint32
		expectedError             error
	}{
		"Sets successfully": {
			fundingRateClampFactorPpm: 6_000_000,
			expectedError:             nil,
		},
		"Sets successfully: max funding rate": {
			fundingRateClampFactorPpm: math.MaxUint32,
			expectedError:             nil,
		},
		"Failure: funding rate clamp factor ppm is zero": {
			fundingRateClampFactorPpm: 0,
			expectedError:             types.ErrFundingRateClampFactorPpmIsZero,
		},
	}

	// Test setup.
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := keeper.SetFundingRateClampFactorPpm(ctx, tc.fundingRateClampFactorPpm)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				// Check that value in store is as expected.
				got := keeper.GetFundingRateClampFactorPpm(ctx)
				require.Equal(t, tc.fundingRateClampFactorPpm, got)
			}
		})
	}
}

func TestSetPremiumVoteClampFactorPpm(t *testing.T) {
	tests := map[string]struct {
		premiumVoteClampFactorPpm uint32
		expectedError             error
	}{
		"Sets successfully": {
			premiumVoteClampFactorPpm: 60_000_000,
			expectedError:             nil,
		},
		"Sets successfully: max uint32": {
			premiumVoteClampFactorPpm: math.MaxUint32,
			expectedError:             nil,
		},
		"Failure: premium vote clamp factor ppm is zero": {
			premiumVoteClampFactorPpm: 0,
			expectedError:             types.ErrPremiumVoteClampFactorPpmIsZero,
		},
	}

	// Test setup.
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := keeper.SetPremiumVoteClampFactorPpm(ctx, tc.premiumVoteClampFactorPpm)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)
				// Check that value in store is as expected.
				got := keeper.GetPremiumVoteClampFactorPpm(ctx)
				require.Equal(t, tc.premiumVoteClampFactorPpm, got)
			}
		})
	}
}

func TestSetMinNumVotesPerSample(t *testing.T) {
	tests := map[string]struct {
		minNumVotesPerSample uint32
	}{
		"Sets successfully: zero": {
			minNumVotesPerSample: 0,
		},
		"Sets successfully: default genesis value": {
			minNumVotesPerSample: 15,
		},
		"Sets successfully: max uint32": {
			minNumVotesPerSample: math.MaxUint32,
		},
	}

	// Test setup.
	ctx, keeper, _, _, _ := keepertest.PerpetualsKeepers(t)

	// Run tests.
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := keeper.SetMinNumVotesPerSample(ctx, tc.minNumVotesPerSample)

			require.NoError(t, err)
			// Check that value in store is as expected.
			got := keeper.GetMinNumVotesPerSample(ctx)
			require.Equal(t, tc.minNumVotesPerSample, got)
		})
	}
}
