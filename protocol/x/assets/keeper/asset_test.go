package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

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
			fmt.Sprintf("symbol-%v", i), // Symbol
			fmt.Sprintf("denom-%v", i),  // Denom
			int32(i),                    // DenomExponent
			hasMarket,                   // HasMarket
			marketId,                    // MarketId
			int32(i),                    // AtomicResolution
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
		"foo-symbol", // symbol
		"foo-denom",  // denom
		-6,           // denomExponent
		true,
		uint32(999),
		int32(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, "999").Error())

	// Does not create an asset.
	numAssets := keeper.GetNumAssets(ctx)
	require.Equal(t, uint32(0), numAssets)
}

func TestCreateAsset_MarketIdInvalid(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	// Throws error when creating asset for invalid marketId.
	_, err := keeper.CreateAsset(
		ctx,
		"foo-symbol", // symbol
		"foo-denom",  // denom
		-6,           // denomExponent
		false,
		uint32(1),
		int32(-1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrInvalidMarketId, "Market ID: 1").Error())

	// Does not create an asset.
	numAssets := keeper.GetNumAssets(ctx)
	require.Equal(t, uint32(0), numAssets)
}

func TestCreateAsset_AssetAlreadyExists(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)

	keepertest.CreateNMarkets(t, ctx, pricesKeeper, 1)

	_, err := keeper.CreateAsset(
		ctx,
		"BTC",       // symbol
		"btc-denom", // denom
		-6,          // denomExponent
		false,       // hasMarket
		0,           // marketId
		10,          // atomicResolution
	)
	require.NoError(t, err)

	// Create a new asset with identical denom
	_, err = keeper.CreateAsset(
		ctx,
		"BTC",       // symbol
		"btc-denom", // denom
		-6,          // denomExponent
		false,       // hasMarket
		0,           // marketId
		10,          // atomicResolution
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrAssetDenomAlreadyExists, "btc-denom").Error())
}

func TestModifyAsset_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	items, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	numMarkets := pricesKeeper.GetNumMarkets(ctx)
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
		newItem, err := keeper.GetAsset(ctx, item.Id)
		require.NoError(t, err)
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
		uint32(0),
		true,
		uint32(1),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrAssetDoesNotExist, "0").Error())
	require.ErrorIs(t, err, types.ErrAssetDoesNotExist)

	// Actually create the asset
	_, err = createNAssets(t, ctx, keeper, pricesKeeper, 1)
	require.NoError(t, err)

	// Expect no issue with modifying the asset now
	_, err = keeper.ModifyAsset(
		ctx,
		uint32(0),
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
		uint32(0),
		true,
		uint32(999),
	)
	require.EqualError(t, err, errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, "999").Error())
}

func TestGetDenomById_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	items, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	for _, item := range items {
		denom, err := keeper.GetDenomById(
			ctx,
			item.Id,
		)
		require.NoError(t, err)
		require.Equal(t,
			item.Denom,
			denom,
		)
	}
}

func TestGetDenomById_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	_, err := keeper.GetDenomById(
		ctx,
		0,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrAssetDoesNotExist, "0").Error())
}

func TestGetIdByDenom_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	items, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	for _, item := range items {
		id, err := keeper.GetIdByDenom(ctx,
			item.Denom,
		)
		require.NoError(t, err)
		require.Equal(t,
			item.Id,
			id,
		)
	}
}

func TestGetIdByDenom_NotFound(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	nonExistingDenom := "non-existent-denom"

	_, err = keeper.GetIdByDenom(ctx,
		nonExistingDenom,
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrNoAssetWithDenom, nonExistingDenom).Error())
}

func TestGetAsset_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	items, err := createNAssets(t, ctx, keeper, pricesKeeper, 10)
	require.NoError(t, err)

	for _, item := range items {
		rst, err := keeper.GetAsset(ctx,
			item.Id,
		)
		require.NoError(t, err)
		require.Equal(t,
			nullify.Fill(&item), //nolint:staticcheck
			nullify.Fill(&rst),  //nolint:staticcheck
		)
	}
}

func TestGetAsset_NotFound(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := keeper.GetAsset(ctx,
		uint32(0),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrAssetDoesNotExist, "0").Error())
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

func TestGetAllAssets_MissingAsset(t *testing.T) {
	ctx, keeper, _, _, _, storeKey := keepertest.AssetsKeepers(t, true)

	// Write some bad data to the store
	store := ctx.KVStore(storeKey)
	store.Set(types.KeyPrefix(types.NumAssetsKey), lib.Uint32ToBytes(20))

	// Expect a panic
	require.Panics(t, func() { keeper.GetAllAssets(ctx) })
}

func TestModifyLongInterest_Success(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := createNAssets(t, ctx, keeper, pricesKeeper, 1)
	require.NoError(t, err)
	assetId := uint32(0)

	// Set long interest to positive number.
	asset, err := keeper.ModifyLongInterest(
		ctx,
		assetId,
		true,
		uint64(10),
	)
	require.NoError(t, err)
	getAsset, err := keeper.GetAsset(ctx, assetId)
	require.NoError(t, err)
	require.Equal(t, asset, getAsset)
	require.Equal(t, uint64(10), asset.LongInterest)

	// Decrease long interest.
	asset, err = keeper.ModifyLongInterest(
		ctx,
		assetId,
		false,
		uint64(7),
	)
	require.NoError(t, err)
	getAsset, err = keeper.GetAsset(ctx, assetId)
	require.NoError(t, err)
	require.Equal(t, asset, getAsset)
	require.Equal(t, uint64(3), asset.LongInterest)

	// Set long interest to zero.
	asset, err = keeper.ModifyLongInterest(
		ctx,
		assetId,
		false,
		uint64(3),
	)
	require.NoError(t, err)
	getAsset, err = keeper.GetAsset(ctx, assetId)
	require.NoError(t, err)
	require.Equal(t, asset, getAsset)
	require.Equal(t, uint64(0), asset.LongInterest)
}

func TestModifyLongInterest_CannotNegative(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := createNAssets(t, ctx, keeper, pricesKeeper, 1)
	require.NoError(t, err)
	assetId := uint32(0)

	// Set long interest to positive number.
	asset, err := keeper.ModifyLongInterest(
		ctx,
		assetId,
		true,
		uint64(10),
	)
	require.NoError(t, err)
	getAsset, err := keeper.GetAsset(ctx, assetId)
	require.NoError(t, err)
	require.Equal(t, asset, getAsset)

	// Fails if long interest would be negative.
	asset, err = keeper.ModifyLongInterest(
		ctx,
		assetId,
		false,
		uint64(12),
	)
	require.EqualError(t, err, errorsmod.Wrap(types.ErrNegativeLongInterest, "0").Error())
	getAsset, err = keeper.GetAsset(ctx, assetId)
	require.NoError(t, err)
	require.Equal(t, asset, getAsset)
}

func TestGetNetCollateral(t *testing.T) {
	ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
	_, err := createNAssets(t, ctx, keeper, pricesKeeper, 2)
	require.NoError(t, err)

	netCollateral, err := keeper.GetNetCollateral(
		ctx,
		lib.UsdcAssetId,
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
		lib.UsdcAssetId,
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
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(11)),
			expectedConvertedQuantums: big.NewInt(1100),
		},
		"atomicResolution < denomExponent, divisble, DenomExponent=-6, AtomicResolution=-7": {
			denomExponent:             -6,
			atomicResolution:          -7,
			quantumsToConvert:         big.NewInt(1120),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(112)),
			expectedConvertedQuantums: big.NewInt(1120),
		},
		"atomicResolution < denomExponent, not divisble, DenomExponent=-6,AtomicResolution=-8": {
			denomExponent:             -6,
			atomicResolution:          -8,
			quantumsToConvert:         big.NewInt(1125),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(11)), // 11.25 rounded down
			expectedConvertedQuantums: big.NewInt(1100),                       // 11 * 100
		},
		"atomicResolution < denomExponent, not, divisble, DenomExponent=-6, AtomicResolution=-7": {
			denomExponent:             -6,
			atomicResolution:          -7,
			quantumsToConvert:         big.NewInt(1125),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(112)), // 112.5 rounded down
			expectedConvertedQuantums: big.NewInt(1120),                        // 112 * 10
		},
		"atomicResolution < denomExponent, not, divisble, DenomExponent=1, AtomicResolution=-3": {
			denomExponent:             1,
			atomicResolution:          -3,
			quantumsToConvert:         big.NewInt(123456),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(12)), // 12.3456 rounded down
			expectedConvertedQuantums: big.NewInt(120000),                     // 12*10000
		},
		"atomicResolution = denomExponent, DenomExponent=-6, AtomicResolution=-6": {
			denomExponent:             -6,
			atomicResolution:          -6,
			quantumsToConvert:         big.NewInt(1500),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(1500)),
			expectedConvertedQuantums: big.NewInt(1500),
		},
		"atomicResolution = denomExponent, DenomExponent=-6, AtomicResolution=-6, large input": {
			denomExponent:             -6,
			atomicResolution:          -6,
			quantumsToConvert:         big.NewInt(12345678),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(12345678)),
			expectedConvertedQuantums: big.NewInt(12345678),
		},
		"atomicResolution > denomExponent": {
			denomExponent:             -6,
			atomicResolution:          -4,
			quantumsToConvert:         big.NewInt(275),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(27500)),
			expectedConvertedQuantums: big.NewInt(275),
		},
		"atomicResolution > denomExponent, positive AtomicResoluton": {
			denomExponent:             -2,
			atomicResolution:          1,
			quantumsToConvert:         big.NewInt(275),
			expectedCoin:              sdk.NewCoin(testDenom, sdk.NewInt(275000)),
			expectedConvertedQuantums: big.NewInt(275),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

			// Create test asset with the given DenomExponent and AtomicResolution values
			asset, err := keeper.CreateAsset(
				ctx,
				testSymbol,
				testDenom,
				tc.denomExponent,
				false,
				0,
				tc.atomicResolution,
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
		1, /* invalid asset ID */
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
		"TEST-SYMBOL-1",
		"test-denom-1",
		-50, /* invalid denom exponent */
		false,
		0,
		-6,
	)
	require.NoError(t, err)
	_, _, err = keeper.ConvertAssetToCoin(ctx, 0, big.NewInt(100))
	require.ErrorIs(
		t,
		err,
		types.ErrInvalidDenomExponent,
	)

	// Test convert asset with invalid denom exponent.
	_, err = keeper.CreateAsset(
		ctx,
		"TEST-SYMBOL-2",
		"test-denom-2",
		-6,
		false,
		0,
		-50, /* invalid asset atomic resolution */
	)
	require.NoError(t, err)
	_, _, err = keeper.ConvertAssetToCoin(ctx, 1, big.NewInt(100))
	require.ErrorIs(
		t,
		err,
		types.ErrInvalidAssetAtomicResolution,
	)
}
