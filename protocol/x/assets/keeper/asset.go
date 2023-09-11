package keeper

import (
	errorsmod "cosmossdk.io/errors"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"math/big"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

func (k Keeper) CreateAsset(
	ctx sdk.Context,
	symbol string,
	denom string,
	denomExponent int32,
	hasMarket bool,
	marketId uint32,
	atomicResolution int32,
) (types.Asset, error) {
	// Get the `nextId`
	nextId := k.GetNumAssets(ctx)

	_, found := k.internalGetIdByDenom(ctx, denom)
	if found {
		return types.Asset{}, errorsmod.Wrap(types.ErrAssetDenomAlreadyExists, denom)
	}

	// Create the asset
	asset := types.Asset{
		Id:               nextId,
		Symbol:           symbol,
		Denom:            denom,
		DenomExponent:    denomExponent,
		HasMarket:        hasMarket,
		MarketId:         marketId,
		AtomicResolution: atomicResolution,
		LongInterest:     0,
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

	// Store the new `numAssets`
	k.setNumAssets(ctx, nextId+1)

	// Store the denom-to-asset-id mapping
	k.setDenomToId(ctx, asset.Denom, asset.Id)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeAsset,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewAssetCreateEvent(
				nextId,
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
	asset, err := k.GetAsset(ctx, id)
	if err != nil {
		return asset, err
	}

	// Validate market
	if _, err = k.pricesKeeper.GetMarketPrice(ctx, marketId); err != nil {
		return asset, err
	}

	// Modify asset
	asset.HasMarket = hasMarket
	asset.MarketId = marketId

	// Store the modified asset
	k.setAsset(ctx, asset)

	return asset, nil
}

func (k Keeper) ModifyLongInterest(
	ctx sdk.Context,
	id uint32,
	isIncrease bool,
	delta uint64,
) (types.Asset, error) {
	// Get asset
	asset, err := k.GetAsset(ctx, id)
	if err != nil {
		return asset, err
	}

	// Validate delta
	if !isIncrease && delta > asset.LongInterest {
		return asset, errorsmod.Wrap(types.ErrNegativeLongInterest, lib.Uint32ToString(id))
	}

	// Modify asset
	if isIncrease {
		asset.LongInterest += delta
	} else {
		asset.LongInterest -= delta
	}

	// Store the modified asset
	k.setAsset(ctx, asset)
	return asset, nil
}

func (k Keeper) GetNumAssets(
	ctx sdk.Context,
) uint32 {
	store := ctx.KVStore(k.storeKey)
	var rawBytes []byte = store.Get(types.KeyPrefix(types.NumAssetsKey))
	return lib.BytesToUint32(rawBytes)
}

func (k Keeper) setNumAssets(
	ctx sdk.Context,
	numAssets uint32,
) {
	// Get necessary stores
	store := ctx.KVStore(k.storeKey)

	// Set `numAssets`
	store.Set(types.KeyPrefix(types.NumAssetsKey), lib.Uint32ToBytes(numAssets))
}

func (k Keeper) setAsset(
	ctx sdk.Context,
	asset types.Asset,
) {
	b := k.cdc.MustMarshal(&asset)
	assetStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetKeyPrefix))
	assetStore.Set(types.AssetKey(asset.Id), b)
}

func (k Keeper) setDenomToId(
	ctx sdk.Context,
	denom string,
	id uint32,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomToIdKeyPrefix))
	store.Set(types.KeyPrefix(denom), lib.Uint32ToBytes(id))
}

func (k Keeper) internalGetIdByDenom(
	ctx sdk.Context,
	denom string,
) (
	id uint32,
	found bool,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomToIdKeyPrefix))

	idBytes := store.Get(types.KeyPrefix(denom))
	if idBytes == nil {
		return 0, false
	}

	return lib.BytesToUint32(idBytes), true
}

// GetIdByDenom returns the `id` of the asset with a given `denom`.
// Returns an error if the `denom` does not exist.
func (k Keeper) GetIdByDenom(
	ctx sdk.Context,
	denom string,
) (
	id uint32,
	err error,
) {
	id, found := k.internalGetIdByDenom(ctx, denom)

	if !found {
		return 0, errorsmod.Wrap(types.ErrNoAssetWithDenom, denom)
	}

	return id, nil
}

// GetDenomById returns the `denom` of the asset with a given `id`.
// Returns an error if the `id` does not exist.
func (k Keeper) GetDenomById(
	ctx sdk.Context,
	id uint32,
) (
	denom string,
	err error,
) {
	asset, err := k.GetAsset(ctx, id)
	if err != nil {
		return "", err
	}

	return asset.Denom, nil
}

func (k Keeper) GetAsset(
	ctx sdk.Context,
	id uint32,
) (val types.Asset, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetKeyPrefix))

	b := store.Get(types.AssetKey(id))
	if b == nil {
		return val, errorsmod.Wrap(types.ErrAssetDoesNotExist, lib.Uint32ToString(id))
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, nil
}

func (k Keeper) GetAllAssets(
	ctx sdk.Context,
) []types.Asset {
	num := k.GetNumAssets(ctx)
	assets := make([]types.Asset, num)

	for i := uint32(0); i < num; i++ {
		asset, err := k.GetAsset(ctx, i)
		if err != nil {
			panic(err)
		}

		assets[i] = asset
	}

	return assets
}

// GetNetCollateral returns the net collateral that a given position (quantums)
// for a given assetId contributes to an account.
func (k Keeper) GetNetCollateral(
	ctx sdk.Context,
	id uint32,
	bigQuantums *big.Int,
) (
	bigNetCollateralQuoteQuantums *big.Int,
	err error,
) {
	if id == lib.UsdcAssetId {
		return new(big.Int).Set(bigQuantums), nil
	}

	// Get asset
	_, err = k.GetAsset(ctx, id)
	if err != nil {
		return big.NewInt(0), err
	}

	// Balance is zero.
	if bigQuantums.BitLen() == 0 {
		return big.NewInt(0), nil
	}

	// Balance is positive.
	// TODO(DEC-581): add multi-collateral support.
	if bigQuantums.Sign() == 1 {
		return big.NewInt(0), types.ErrNotImplementedMulticollateral
	}

	// Balance is negative.
	// TODO(DEC-582): add margin-trading support.
	return big.NewInt(0), types.ErrNotImplementedMargin
}

// GetMarginRequirements returns the initial and maintenance margin-
// requirements for a given position size for a given assetId.
func (k Keeper) GetMarginRequirements(
	ctx sdk.Context,
	id uint32,
	bigQuantums *big.Int,
) (
	bigInitialMarginQuoteQuantums *big.Int,
	bigMaintenanceMarginQuoteQuantums *big.Int,
	err error,
) {
	// QuoteBalance does not contribute to any margin requirements.
	if id == lib.UsdcAssetId {
		return big.NewInt(0), big.NewInt(0), nil
	}

	// Get asset
	_, err = k.GetAsset(ctx, id)
	if err != nil {
		return big.NewInt(0), big.NewInt(0), err
	}

	// Balance is zero or positive.
	if bigQuantums.Sign() >= 0 {
		return big.NewInt(0), big.NewInt(0), nil
	}

	// Balance is negative.
	// TODO(DEC-582): margin-trading
	return big.NewInt(0), big.NewInt(0), types.ErrNotImplementedMargin
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
	asset, err := k.GetAsset(ctx, assetId)
	if err != nil {
		return nil, sdk.Coin{}, err
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

	bigRatDenomAmount := lib.BigMulPow10(
		quantums,
		asset.AtomicResolution-asset.DenomExponent,
	)

	// round down to get denom amount that was converted.
	bigConvertedDenomAmount := lib.BigRatRound(bigRatDenomAmount, false)

	bigRatConvertedQuantums := lib.BigMulPow10(
		bigConvertedDenomAmount,
		asset.DenomExponent-asset.AtomicResolution,
	)

	bigConvertedQuantums := bigRatConvertedQuantums.Num()

	return bigConvertedQuantums, sdk.NewCoin(asset.Denom, sdk.NewIntFromBigInt(bigConvertedDenomAmount)), nil
}
