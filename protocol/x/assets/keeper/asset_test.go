package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/nullify"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/keeper"
	assetslib "github.com/dydxprotocol/v4-chain/protocol/x/assets/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	priceskeeper "github.com/dydxprotocol/v4-chain/protocol/x/prices/keeper"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

const (
	// firstValidAssetId is the first valid asset ID after the reserved `assetId=0` for USDC.
	firstValidAssetId = uint32(1)
)

// createNAssets creates n test assets with id 1 to n (0 is reserved for USDC)
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
	)
	require.EqualError(t, err, errorsmod.Wrap(pricestypes.ErrMarketPriceDoesNotExist, "999").Error())

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)
}

func TestCreateAsset_InvalidUsdcAsset(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)

	// Throws error when creating an asset with id 0 that's not USDC.
	_, err := keeper.CreateAsset(
		ctx,
		0,
		"foo-symbol", // symbol
		"foo-denom",  // denom
		-6,           // denomExponent
		true,
		uint32(999),
		int32(-1),
	)
	require.ErrorIs(t, err, types.ErrUsdcMustBeAssetZero)

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)

	// Throws error when creating asset USDC with id other than 0.
	_, err = keeper.CreateAsset(
		ctx,
		1,
		constants.Usdc.Symbol,        // symbol
		constants.Usdc.Denom,         // denom
		constants.Usdc.DenomExponent, // denomExponent
		true,
		uint32(999),
		int32(-1),
	)
	require.ErrorIs(t, err, types.ErrUsdcMustBeAssetZero)

	// Does not create an asset.
	require.Len(t, keeper.GetAllAssets(ctx), 0)

	// Throws error when creating asset USDC with unexpected denom exponent.
	_, err = keeper.CreateAsset(
		ctx,
		0,
		constants.Usdc.Symbol, // symbol
		constants.Usdc.Denom,  // denom
		-9,                    // denomExponent
		true,
		uint32(999),
		int32(-1),
	)
	require.ErrorIs(t, err, types.ErrUnexpectedUsdcDenomExponent)

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

func TestGetNetCollateralAndMarginRequirements(t *testing.T) {
	tests := map[string]struct {
		assetId     uint32
		bigQuantums *big.Int
		expectedNC  *big.Int
		expectedIMR *big.Int
		expectedMMR *big.Int
		expectedErr error
	}{
		"USDC asset. Positive Balance": {
			assetId:     types.AssetUsdc.Id,
			bigQuantums: big.NewInt(100),
			expectedNC:  big.NewInt(100),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
		"USDC asset. Negative Balance": {
			assetId:     types.AssetUsdc.Id,
			bigQuantums: big.NewInt(-100),
			expectedNC:  big.NewInt(-100),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
		"USDC asset. Zero Balance": {
			assetId:     types.AssetUsdc.Id,
			bigQuantums: big.NewInt(0),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
		"Non USDC asset. Positive Balance": {
			assetId:     uint32(1),
			bigQuantums: big.NewInt(100),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: types.ErrNotImplementedMulticollateral,
		},
		"Non USDC asset. Negative Balance": {
			assetId:     uint32(1),
			bigQuantums: big.NewInt(-100),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: types.ErrNotImplementedMargin,
		},
		"Non USDC asset. Zero Balance": {
			assetId:     uint32(1),
			bigQuantums: big.NewInt(0),
			expectedNC:  big.NewInt(0),
			expectedIMR: big.NewInt(0),
			expectedMMR: big.NewInt(0),
			expectedErr: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, keeper, pricesKeeper, _, _, _ := keepertest.AssetsKeepers(t, true)
			_, err := createNAssets(t, ctx, keeper, pricesKeeper, 2)
			require.NoError(t, err)

			risk, err := assetslib.GetNetCollateralAndMarginRequirements(
				tc.assetId,
				tc.bigQuantums,
			)

			require.Equal(t, tc.expectedNC, risk.NC)
			require.Equal(t, tc.expectedIMR, risk.IMR)
			require.Equal(t, tc.expectedMMR, risk.MMR)

			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
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
		-50, /* invalid asset atomic resolution */
	)
	require.NoError(t, err)
	_, _, err = keeper.ConvertAssetToCoin(ctx, 2, big.NewInt(100))
	require.ErrorIs(
		t,
		err,
		types.ErrInvalidAssetAtomicResolution,
	)
}

func TestIsPositionUpdatable(t *testing.T) {
	ctx, keeper, _, _, _, _ := keepertest.AssetsKeepers(t, true)
	require.NoError(t, keepertest.CreateUsdcAsset(ctx, keeper))

	// Check Usdc asset is updatable.
	updatable, err := keeper.IsPositionUpdatable(ctx, types.AssetUsdc.Id)
	require.NoError(t, err)
	require.True(t, updatable)

	// Return error for non-existent asset
	_, err = keeper.IsPositionUpdatable(ctx, 100)
	require.ErrorContains(t, err, "Asset does not exist")
}
