package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/nullify"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	priceskeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	// firstValidAssetId is the first valid asset ID after the reserved `assetId=0` for TDAI.
	firstValidAssetId = uint32(1)
)

// createNAssets creates n test assets with id 1 to n (0 is reserved for TDai)
func createNAssets(
	t *testing.T,
	ctx sdk.Context,
	keeper *keeper.Keeper,
	pricesKeeper *priceskeeper.Keeper,
	n int,
) ([]types.Asset, error) {
	items := make([]types.Asset, n)

	keepertest.CreateNMarkets(t, ctx, pricesKeeper, n)

	for i := range items {
		hasMarket := i%2 == 0
		var marketId uint32
		if hasMarket {
			marketId = uint32(i)
		}
		asset, err := keeper.CreateAsset(
			ctx,
			uint32(i+1),                 // AssetId
			fmt.Sprintf("symbol-%v", i), // Symbol
			fmt.Sprintf("denom-%v", i),  // Denom
			int32(i),                    // DenomExponent
			hasMarket,                   // HasMarket
			marketId,                    // MarketId
			int32(i),                    // AtomicResolution
			"1/1",                       // AssetYieldIndex
		)
		if err != nil {
			return items, err
		}

		items[i] = asset
	}

	return items, nil
}

func TestCreateAsset_MarketNotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	// Throws error when creating asset for invalid marketId.
	_, err := keeper.CreateAsset(
		ctx,
		firstValidAssetId,
		"foo-symbol", // symbol
		"foo-denom",  // denom
		-6,           // denomExponent
		true,
		uint32(999),
		int32(-1),
		"1/1", // AssetYieldIndex
	)
	require.EqualError(t, err, errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, "999").Error())

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)
}

func TestCreateAsset_InvalidTDaiAsset(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	// Throws error when creating an asset with id 0 that's not TDAI.
	_, err := keeper.CreateAsset(
		ctx,
		0,
		"foo-symbol", // symbol
		"foo-denom",  // denom
		-6,           // denomExponent
		true,
		uint32(999),
		int32(-1),
		"1/1", // AssetYieldIndex
	)
	require.ErrorIs(t, err, types.ErrTDaiMustBeAssetZero)

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)

	// Throws error when creating asset TDAI with id other than 0.
	_, err = keeper.CreateAsset(
		ctx,
		1,
		constants.TDai.Symbol,        // symbol
		constants.TDai.Denom,         // denom
		constants.TDai.DenomExponent, // denomExponent
		true,
		uint32(999),
		int32(-1),
		"1/1", // AssetYieldIndex
	)
	require.ErrorIs(t, err, types.ErrTDaiMustBeAssetZero)

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)

	// Throws error when creating asset TDAI with unexpected denom exponent.
	_, err = keeper.CreateAsset(
		ctx,
		0,
		constants.TDai.Symbol, // symbol
		constants.TDai.Denom,  // denom
		-9,                    // denomExponent
		true,
		uint32(999),
		int32(-1),
		"1/1", // AssetYieldIndex
	)
	require.ErrorIs(t, err, types.ErrUnexpectedTDaiDenomExponent)

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)
}

func TestCreateAsset_MarketIdInvalid(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	// Throws error when creating asset for invalid marketId.
	_, err := keeper.CreateAsset(
		ctx,
		firstValidAssetId,
		"foo-symbol", // symbol
		"foo-denom",  // denom
		-6,           // denomExponent
		false,
		uint32(1),
		int32(-1),
		"1/1", // AssetYieldIndex
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrInvalidMarketId, "Market ID: 1").Error())

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)
}

func TestCreateAsset_AssetAlreadyExists(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)

	keepertest.CreateNMarkets(t, ctx, pricesKeeper, 1)

	_, err := keeper.CreateAsset(
		ctx,
		firstValidAssetId,
		"BTC",       // symbol
		"btc-denom", // denom
		-6,          // denomExponent
		false,       // hasMarket
		0,           // marketId
		10,          // atomicResolution
		"1/1",       // AssetYieldIndex
	)
	require.NoError(t, err)

	// Create a new asset with identical denom
	_, err = keeper.CreateAsset(
		ctx,
		2,
		"BTC",       // symbol
		"btc-denom", // denom
		-6,          // denomExponent
		false,       // hasMarket
		0,           // marketId
		10,          // atomicResolution
		"1/1",       // AssetYieldIndex
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrAssetDenomAlreadyExists, "btc-denom").Error())

	// Create a new asset with the same ID
	_, err = keeper.CreateAsset(
		ctx,
		firstValidAssetId,
		"BTC-COPY",       // symbol
		"btc-denom-copy", // denom
		-6,               // denomExponent
		false,            // hasMarket
		0,                // marketId
		10,               // atomicResolution
		"1/1",            // AssetYieldIndex
	)
	require.ErrorIs(t, err, types.ErrAssetIdAlreadyExists)
}

func TestModifyAsset_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	items, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	numMarkets := keepertest.GetNumMarkets(t, ctx, pricesKeeper)
	for i, item := range items {
		// Modify each field arbitrarily and
		// verify the fields were modified in state
		hasMarket := (i%2 == 0)
		marketId := uint32(i*2) % numMarkets
		retItem, err := keeper.ModifyAsset(
			ctx,
			item.Id,
			hasMarket,
			marketId,
		)
		require.NoError(t, err)
		newItem, exists := keeper.GetAsset(ctx, item.Id)
		require.True(t, exists)
		require.Equal(t,
			retItem,
			newItem,
		)
		require.Equal(t,
			fmt.Sprintf("denom-%v", i),
			newItem.Denom,
		)
		require.Equal(t,
			hasMarket,
			newItem.HasMarket,
		)
		require.Equal(t,
			marketId,
			newItem.MarketId,
		)
		require.Equal(t,
			int32(i),
			newItem.AtomicResolution,
		)
	}
}

func TestModifyAsset_NotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)

	// Expect error when modifying non-existent asset
	_, err := keeper.ModifyAsset(
		ctx,
		firstValidAssetId,
		true,
		uint32(1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrAssetDoesNotExist, "1").Error())
	require.ErrorIs(t, err, types.ErrAssetDoesNotExist)

	// Actually create the asset
	_, err = createNAssets(t, ctx, keeper, pricesKeeper, 1)
	require.NoError(t, err)

	// Expect no issue with modifying the asset now
	_, err = keeper.ModifyAsset(
		ctx,
		firstValidAssetId,
		true,
		uint32(0),
	)
	require.NoError(t, err)
}

func TestModifyAsset_MarketNotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := createNAssets(t, ctx, keeper, pricesKeeper, 1)
	require.NoError(t, err)

	_, err = keeper.ModifyAsset(
		ctx,
		firstValidAssetId,
		true,
		uint32(999),
	)
	require.EqualError(t, err, errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, "999").Error())
}

func TestGetAsset_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	items, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	for _, item := range items {
		rst, exists := keeper.GetAsset(ctx,
			item.Id,
		)
		require.True(t, exists)
		require.Equal(t,
			nullify.Fill(&item), //nolint:staticcheck
			nullify.Fill(&rst),  //nolint:staticcheck
		)
	}
}

func TestGetAsset_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, exists := keeper.GetAsset(ctx,
		uint32(0),
	)
	require.False(t, exists)
}

func TestGetAllAssets_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	items, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	require.ElementsMatch(t,
		nullify.Fill(items),                    //nolint:staticcheck
		nullify.Fill(keeper.GetAllAssets(ctx)), //nolint:staticcheck
	)
}

func TestGetNetCollateral(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := createNAssets(t, ctx, keeper, pricesKeeper, 2)
	require.NoError(t, err)

	netCollateral, err := keeper.GetNetCollateral(
		ctx,
		types.AssetTDai.Id,
		new(big.Int).SetInt64(100),
	)
	require.NoError(t, err)
	require.Equal(t, new(big.Int).SetInt64(100), netCollateral)

	_, err = keeper.GetNetCollateral(
		ctx,
		uint32(1),
		new(big.Int).SetInt64(100),
	)
	require.EqualError(t, types.ErrNotImplementedMulticollateral, err.Error())

	_, err = keeper.GetNetCollateral(
		ctx,
		uint32(1),
		new(big.Int).SetInt64(-100),
	)
	require.EqualError(t, types.ErrNotImplementedMargin, err.Error())
}

func TestGetMarginRequirements(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := createNAssets(t, ctx, keeper, pricesKeeper, 2)
	require.NoError(t, err)

	initial, maintenance, err := keeper.GetMarginRequirements(
		ctx,
		types.AssetTDai.Id,
		new(big.Int).SetInt64(100),
	)
	require.NoError(t, err)
	require.Equal(t, new(big.Int), initial)
	require.Equal(t, new(big.Int), maintenance)

	initial, maintenance, err = keeper.GetMarginRequirements(
		ctx,
		uint32(1),
		new(big.Int).SetInt64(100),
	)
	require.NoError(t, err)
	require.Equal(t, new(big.Int), initial)
	require.Equal(t, new(big.Int), maintenance)

	_, _, err = keeper.GetMarginRequirements(
		ctx,
		uint32(1),
		new(big.Int).SetInt64(-100),
	)
	require.EqualError(t, types.ErrNotImplementedMargin, err.Error())
}

func TestConvertAssetToCoin_Success(t *testing.T) {
	testSymbol := "TEST_SYMBOL"
	testDenom := "test_denom"

	tests := map[string]struct {
		denomExponent             int32
		atomicResolution          int32
		quantumsToConvert         *big.Int
		expectedCoin              sdk.Coin
		expectedConvertedQuantums *big.Int
	}{
		"atomicResolution < denomExponent, divisble, DenomExponent=-6, AtomicResolution=-8": {
			denomExponent:             -6,
			atomicResolution:          -8,
			quantumsToConvert:         big.NewInt(1100),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(11)),
			expectedConvertedQuantums: big.NewInt(1100),
		},
		"atomicResolution < denomExponent, divisble, DenomExponent=-6, AtomicResolution=-7": {
			denomExponent:             -6,
			atomicResolution:          -7,
			quantumsToConvert:         big.NewInt(1120),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(112)),
			expectedConvertedQuantums: big.NewInt(1120),
		},
		"atomicResolution < denomExponent, not divisble, DenomExponent=-6,AtomicResolution=-8": {
			denomExponent:             -6,
			atomicResolution:          -8,
			quantumsToConvert:         big.NewInt(1125),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(11)), // 11.25 rounded down
			expectedConvertedQuantums: big.NewInt(1100),                           // 11 * 100
		},
		"atomicResolution < denomExponent, not, divisble, DenomExponent=-6, AtomicResolution=-7": {
			denomExponent:             -6,
			atomicResolution:          -7,
			quantumsToConvert:         big.NewInt(1125),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(112)), // 112.5 rounded down
			expectedConvertedQuantums: big.NewInt(1120),                            // 112 * 10
		},
		"atomicResolution < denomExponent, not, divisble, DenomExponent=1, AtomicResolution=-3": {
			denomExponent:             1,
			atomicResolution:          -3,
			quantumsToConvert:         big.NewInt(123456),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(12)), // 12.3456 rounded down
			expectedConvertedQuantums: big.NewInt(120000),                         // 12*10000
		},
		"atomicResolution = denomExponent, DenomExponent=-6, AtomicResolution=-6": {
			denomExponent:             -6,
			atomicResolution:          -6,
			quantumsToConvert:         big.NewInt(1500),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(1500)),
			expectedConvertedQuantums: big.NewInt(1500),
		},
		"atomicResolution = denomExponent, DenomExponent=-6, AtomicResolution=-6, large input": {
			denomExponent:             -6,
			atomicResolution:          -6,
			quantumsToConvert:         big.NewInt(12345678),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(12345678)),
			expectedConvertedQuantums: big.NewInt(12345678),
		},
		"atomicResolution > denomExponent": {
			denomExponent:             -6,
			atomicResolution:          -4,
			quantumsToConvert:         big.NewInt(275),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(27500)),
			expectedConvertedQuantums: big.NewInt(275),
		},
		"atomicResolution > denomExponent, positive AtomicResoluton": {
			denomExponent:             -2,
			atomicResolution:          1,
			quantumsToConvert:         big.NewInt(275),
			expectedCoin:              sdk.NewCoin(testDenom, sdkmath.NewInt(275000)),
			expectedConvertedQuantums: big.NewInt(275),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

			// Create test asset with the given DenomExponent and AtomicResolution values
			asset, err := keeper.CreateAsset(
				ctx,
				firstValidAssetId,
				testSymbol,
				testDenom,
				tc.denomExponent,
				false,
				0,
				tc.atomicResolution,
				"1/1", // AssetYieldIndex
			)
			require.NoError(t, err)

			// Call ConvertAssetToCoin
			convertedQuantums, coin, err := keeper.ConvertAssetToCoin(ctx, asset.Id, tc.quantumsToConvert)

			// Check for successful conversion
			require.NoError(t, err)
			require.NotNil(t, convertedQuantums)

			// Check if the converted quantums and denom amount are as expected
			require.Equal(t, tc.expectedConvertedQuantums, convertedQuantums)
			require.Equal(t, tc.expectedCoin, coin)

			assetEvents := keepertest.GetAssetCreateEventsFromIndexerBlock(ctx, keeper)
			require.Len(t, assetEvents, 1)

			expectedEvent := indexerevents.NewAssetCreateEvent(
				asset.Id,
				testSymbol,
				false,
				0,
				tc.atomicResolution,
			)
			require.Contains(t, assetEvents, expectedEvent)
		})
	}
}

func TestConvertAssetToCoin_Failure(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	// Test convert asset with invalid asset ID.
	_, _, err := keeper.ConvertAssetToCoin(
		ctx,
		firstValidAssetId, /* invalid asset ID */
		big.NewInt(100),
	)

	require.ErrorIs(
		t,
		err,
		types.ErrAssetDoesNotExist,
	)

	// Test convert asset with invalid denom exponent.
	_, err = keeper.CreateAsset(
		ctx,
		firstValidAssetId,
		"TEST-SYMBOL-1",
		"test-denom-1",
		-50, /* invalid denom exponent */
		false,
		0,
		-6,
		"1/1", // AssetYieldIndex
	)
	require.NoError(t, err)

	_, _, err = keeper.ConvertAssetToCoin(ctx, 1, big.NewInt(100))
	require.ErrorIs(
		t,
		err,
		types.ErrInvalidDenomExponent,
	)

	// Test convert asset with invalid denom exponent.
	_, err = keeper.CreateAsset(
		ctx,
		2,
		"TEST-SYMBOL-2",
		"test-denom-2",
		-6,
		false,
		0,
		-50,   /* invalid asset atomic resolution */
		"1/1", // AssetYieldIndex
	)
	require.NoError(t, err)
	_, _, err = keeper.ConvertAssetToCoin(ctx, 2, big.NewInt(100))
	require.ErrorIs(
		t,
		err,
		types.ErrInvalidAssetAtomicResolution,
	)
}

func TestConvertCoinToAsset_Success(t *testing.T) {
	testSymbol := "TEST_SYMBOL"
	testDenom := "test_denom"

	tests := map[string]struct {
		denomExponent               int32
		atomicResolution            int32
		coinToConvert               sdk.Coin
		expectedQuantums            *big.Int
		expectedConvertedCoinAmount *big.Int
	}{
		"atomicResolution < denomExponent, DenomExponent=-6, AtomicResolution=-8": {
			denomExponent:               -6,
			atomicResolution:            -8,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(11)),
			expectedQuantums:            big.NewInt(1100),
			expectedConvertedCoinAmount: big.NewInt(11),
		},
		"atomicResolution < denomExponent, DenomExponent=-6, AtomicResolution=-7": {
			denomExponent:               -6,
			atomicResolution:            -7,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(112)),
			expectedQuantums:            big.NewInt(1120),
			expectedConvertedCoinAmount: big.NewInt(112),
		},
		"atomicResolution < denomExponent, DenomExponent=1, AtomicResolution=-3": {
			denomExponent:               1,
			atomicResolution:            -3,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(12)),
			expectedQuantums:            big.NewInt(120000),
			expectedConvertedCoinAmount: big.NewInt(12),
		},
		"atomicResolution = denomExponent, DenomExponent=-6, AtomicResolution=-6": {
			denomExponent:               -6,
			atomicResolution:            -6,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(1500)),
			expectedQuantums:            big.NewInt(1500),
			expectedConvertedCoinAmount: big.NewInt(1500),
		},
		"atomicResolution = denomExponent, DenomExponent=-6, AtomicResolution=-6, large input": {
			denomExponent:               -6,
			atomicResolution:            -6,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(12345678)),
			expectedQuantums:            big.NewInt(12345678),
			expectedConvertedCoinAmount: big.NewInt(12345678),
		},
		"atomicResolution > denomExponent": {
			denomExponent:               -6,
			atomicResolution:            -4,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(27500)),
			expectedQuantums:            big.NewInt(275),
			expectedConvertedCoinAmount: big.NewInt(27500),
		},
		"atomicResolution > denomExponent, positive AtomicResolution": {
			denomExponent:               -2,
			atomicResolution:            1,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(275000)),
			expectedQuantums:            big.NewInt(275),
			expectedConvertedCoinAmount: big.NewInt(275000),
		},
		"atomicResolution > denomExponent, positive AtomicResolution. Uneven coin number.": {
			denomExponent:               -2,
			atomicResolution:            1,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(123456)),
			expectedQuantums:            big.NewInt(123),
			expectedConvertedCoinAmount: big.NewInt(123000),
		},
		"atomicResolution > denomExponent, positive AtomicResolution. Zero conversion.": {
			denomExponent:               -2,
			atomicResolution:            1,
			coinToConvert:               sdk.NewCoin(testDenom, sdkmath.NewInt(999)),
			expectedQuantums:            big.NewInt(0),
			expectedConvertedCoinAmount: big.NewInt(0),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

			// Create test asset with the given DenomExponent and AtomicResolution values
			asset, err := keeper.CreateAsset(
				ctx,
				firstValidAssetId,
				testSymbol,
				testDenom,
				tc.denomExponent,
				false,
				0,
				tc.atomicResolution,
				"1/1", // AssetYieldIndex
			)
			require.NoError(t, err)

			// Call ConvertCoinToAsset
			quantums, convertedCoinAmount, err := keeper.ConvertCoinToAsset(ctx, asset.Id, tc.coinToConvert)

			// Check for successful conversion
			require.NoError(t, err)

			require.NotNil(t, quantums)
			require.Equal(t, tc.expectedQuantums, quantums)

			require.NotNil(t, convertedCoinAmount)
			require.Equal(t, tc.expectedConvertedCoinAmount, convertedCoinAmount)

			assetEvents := keepertest.GetAssetCreateEventsFromIndexerBlock(ctx, keeper)
			require.Len(t, assetEvents, 1)

			expectedEvent := indexerevents.NewAssetCreateEvent(
				asset.Id,
				testSymbol,
				false,
				0,
				tc.atomicResolution,
			)
			require.Contains(t, assetEvents, expectedEvent)
		})
	}
}

func TestConvertCoinToAsset_Failure(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	testDenom := "test_denom"

	// Test convert coin with invalid asset ID.
	quantums, convertedDenom, err := keeper.ConvertCoinToAsset(
		ctx,
		firstValidAssetId, /* invalid asset ID */
		sdk.NewCoin(testDenom, sdkmath.NewInt(100)),
	)

	require.ErrorIs(
		t,
		err,
		types.ErrAssetDoesNotExist,
	)

	require.Nil(t, quantums)
	require.Nil(t, convertedDenom)

	// Test convert coin with invalid denom exponent.
	_, err = keeper.CreateAsset(
		ctx,
		firstValidAssetId,
		"TEST-SYMBOL-1",
		testDenom,
		-50, /* invalid denom exponent */
		false,
		0,
		-6,
		"1/1", // AssetYieldIndex
	)
	require.NoError(t, err)

	quantums, convertedDenom, err = keeper.ConvertCoinToAsset(ctx, firstValidAssetId, sdk.NewCoin(testDenom, sdkmath.NewInt(100)))
	require.ErrorIs(
		t,
		err,
		types.ErrInvalidDenomExponent,
	)

	require.Nil(t, quantums)
	require.Nil(t, convertedDenom)

	// Test convert coin with invalid atomic resolution.
	_, err = keeper.CreateAsset(
		ctx,
		2,
		"TEST-SYMBOL-2",
		"test-denom-2",
		-6,
		false,
		0,
		-50,   /* invalid asset atomic resolution */
		"1/1", // AssetYieldIndex
	)
	require.NoError(t, err)
	quantums, convertedDenom, err = keeper.ConvertCoinToAsset(ctx, 2, sdk.NewCoin("test-denom-2", sdkmath.NewInt(100)))
	require.ErrorIs(
		t,
		err,
		types.ErrInvalidAssetAtomicResolution,
	)

	require.Nil(t, quantums)
	require.Nil(t, convertedDenom)
}

func TestIsPositionUpdatable(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)
	require.NoError(t, keepertest.CreateTDaiAsset(ctx, keeper))

	// Check TDai asset is updatable.
	updatable, err := keeper.IsPositionUpdatable(ctx, types.AssetTDai.Id)
	require.NoError(t, err)
	require.True(t, updatable)

	// Return error for non-existent asset
	_, err = keeper.IsPositionUpdatable(ctx, 100)
	require.ErrorContains(t, err, "Asset does not exist")
}
