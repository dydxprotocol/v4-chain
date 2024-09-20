package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AssetsKeeper interface {
	ConvertAssetToCoin(ctx sdk.Context, assetId uint32, quantums *big.Int) (*big.Int, sdk.Coin, error)
	CreateAsset(
		ctx sdk.Context,
		assetId uint32,
		symbol string,
		denom string,
		denomExponent int32,
		hasMarket bool,
		marketId uint32,
		atomicResolution int32,
	) (
		Asset,
		error,
	)

	GetAsset(ctx sdk.Context, id uint32) (Asset, bool)

	GetAllAssets(ctx sdk.Context) []Asset

	IsPositionUpdatable(ctx sdk.Context, id uint32) (bool, error)

	ModifyAsset(ctx sdk.Context, id uint32, hasMarket bool, marketId uint32) (Asset, error)
}
