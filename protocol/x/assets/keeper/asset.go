package keeper

import (
	"math/big"
	"sort"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

func (k Keeper) CreateAsset(
	ctx sdk.Context,
	assetId uint32,
	symbol string,
	denom string,
	denomExponent int32,
	hasMarket bool,
	marketId uint32,
	atomicResolution int32,
) (types.Asset, error) {
	if prevAsset, exists := k.GetAsset(ctx, assetId); exists {
		return types.Asset{}, errorsmod.Wrapf(
			types.ErrAssetIdAlreadyExists,
			"previous asset = %v",
			prevAsset,
		)
	}

	if assetId == types.AssetUsdc.Id {
		// Ensure assetId zero is always USDC. This is a protocol-wide invariant.
		if denom != types.AssetUsdc.Denom {
			return types.Asset{}, types.ErrUsdcMustBeAssetZero
		}

		// Confirm that USDC asset has the expected denom exponent (-6).
		// This is an important invariant before coin-to-quote-quantum conversion
		// is correctly implemented. See CLOB-871 for details.
		if denomExponent != types.AssetUsdc.DenomExponent {
			return types.Asset{}, errorsmod.Wrapf(
				types.ErrUnexpectedUsdcDenomExponent,
				"expected = %v, actual = %v",
				types.AssetUsdc.DenomExponent,
				denomExponent,
			)
		}
	}

	// Ensure USDC is not created with a non-zero assetId. This is a protocol-wide invariant.
	if assetId != types.AssetUsdc.Id && denom == types.AssetUsdc.Denom {
		return types.Asset{}, types.ErrUsdcMustBeAssetZero
	}

	// Ensure the denom is unique versus existing assets.
	allAssets := k.GetAllAssets(ctx)
	for _, asset := range allAssets {
		if asset.Denom == denom {
			return types.Asset{}, errorsmod.Wrap(types.ErrAssetDenomAlreadyExists, denom)
		}
	}

	// Create the asset
	asset := types.Asset{
		Id:               assetId,
		Symbol:           symbol,
		Denom:            denom,
		DenomExponent:    denomExponent,
		HasMarket:        hasMarket,
		MarketId:         marketId,
		AtomicResolution: atomicResolution,
	}

	// Validate market
	if hasMarket {
		if _, err := k.pricesKeeper.GetMarketPrice(ctx, marketId); err != nil {
			return asset, err
		}
	} else if marketId > 0 {
		return asset, errorsmod.Wrapf(
			types.ErrInvalidMarketId,
			"Market ID: %v",
			marketId,
		)
	}

	// Store the new asset
	k.setAsset(ctx, asset)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeAsset,
		indexerevents.AssetEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewAssetCreateEvent(
				assetId,
				asset.Symbol,
				asset.HasMarket,
				asset.MarketId,
				asset.AtomicResolution,
			),
		),
	)

	return asset, nil
}

func (k Keeper) ModifyAsset(
	ctx sdk.Context,
	id uint32,
	hasMarket bool,
	marketId uint32,
) (types.Asset, error) {
	// Get asset
	asset, exists := k.GetAsset(ctx, id)
	if !exists {
		return asset, errorsmod.Wrap(types.ErrAssetDoesNotExist, lib.UintToString(id))
	}

	// Validate market
	if _, err := k.pricesKeeper.GetMarketPrice(ctx, marketId); err != nil {
		return asset, err
	}

	// Modify asset
	asset.HasMarket = hasMarket
	asset.MarketId = marketId

	// Store the modified asset
	k.setAsset(ctx, asset)

	return asset, nil
}

func (k Keeper) setAsset(
	ctx sdk.Context,
	asset types.Asset,
) {
	b := k.cdc.MustMarshal(&asset)
	assetStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AssetKeyPrefix))
	assetStore.Set(lib.Uint32ToKey(asset.Id), b)
}

func (k Keeper) GetAsset(
	ctx sdk.Context,
	id uint32,
) (val types.Asset, exists bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AssetKeyPrefix))

	b := store.Get(lib.Uint32ToKey(id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllAssets(
	ctx sdk.Context,
) (list []types.Asset) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AssetKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Asset
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Id < list[j].Id
	})

	return list
}

// ConvertAssetToCoin converts the given `assetId` and `quantums` used in `x/asset`,
// to an `sdk.Coin` in correspoding `denom` and `amount` used in `x/bank`.
// Also outputs `convertedQuantums` which has the equal value as converted `sdk.Coin`.
// The conversion is done with the formula:
//
//	denom_amount = quantums * 10^(atomic_resolution - denom_exponent)
//
// If the resulting `denom_amount` is not an integer, it is rounded down,
// and `convertedQuantums` of the equal value is returned. The upstream
// transfer function should adjust asset balance with `convertedQuantums`
// to ensure that that no fund is ever lost in the conversion due to rounding error.
//
// Example:
// Assume `denom_exponent` = -7, `atomic_resolution` = -8,
// ConvertAssetToCoin(`101 quantums`) should output:
// - `convertedQuantums` = 100 quantums
// -  converted coin amount = 10 coin
func (k Keeper) ConvertAssetToCoin(
	ctx sdk.Context,
	assetId uint32,
	quantums *big.Int,
) (
	convertedQuantums *big.Int,
	coin sdk.Coin,
	err error,
) {
	asset, exists := k.GetAsset(ctx, assetId)
	if !exists {
		return nil, sdk.Coin{}, errorsmod.Wrap(
			types.ErrAssetDoesNotExist, lib.UintToString(assetId))
	}

	if lib.AbsInt32(asset.AtomicResolution) > types.MaxAssetUnitExponentAbs {
		return nil, sdk.Coin{}, errorsmod.Wrapf(
			types.ErrInvalidAssetAtomicResolution,
			"asset: %+v",
			asset,
		)
	}

	if lib.AbsInt32(asset.DenomExponent) > types.MaxAssetUnitExponentAbs {
		return nil, sdk.Coin{}, errorsmod.Wrapf(
			types.ErrInvalidDenomExponent,
			"asset: %+v",
			asset,
		)
	}

	exponent := asset.AtomicResolution - asset.DenomExponent
	p10, inverse := lib.BigPow10(exponent)
	var resultDenom *big.Int
	var resultQuantums *big.Int
	if inverse {
		resultDenom = new(big.Int).Div(quantums, p10)
		resultQuantums = new(big.Int).Mul(resultDenom, p10)
	} else {
		resultDenom = new(big.Int).Mul(quantums, p10)
		resultQuantums = new(big.Int).Div(resultDenom, p10)
	}

	return resultQuantums, sdk.NewCoin(asset.Denom, sdkmath.NewIntFromBigInt(resultDenom)), nil
}

// IsPositionUpdatable returns whether position of an asset is updatable.
func (k Keeper) IsPositionUpdatable(
	ctx sdk.Context,
	id uint32,
) (
	updatable bool,
	err error,
) {
	_, exists := k.GetAsset(ctx, id)
	if !exists {
		return false, errorsmod.Wrap(types.ErrAssetDoesNotExist, lib.UintToString(id))
	}
	// All existing assets are by default updatable.
	return true, nil
}
