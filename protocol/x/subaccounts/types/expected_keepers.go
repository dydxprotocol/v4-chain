package types

import (
	"context"
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type AssetsKeeper interface {
	IsPositionUpdatable(
		ctx sdk.Context,
		id uint32,
	) (
		updatable bool,
		err error,
	)
	ConvertAssetToCoin(
		ctx sdk.Context,
		assetId uint32,
		quantums *big.Int,
	) (
		convertedQuantums *big.Int,
		coin sdk.Coin,
		err error,
	)
	GetAsset(
		ctx sdk.Context,
		id uint32,
	) (
		val assettypes.Asset,
		exists bool,
	)
}

type PerpetualsKeeper interface {
	IsPositionUpdatable(
		ctx sdk.Context,
		id uint32,
	) (
		updatable bool,
		err error,
	)
	GetPerpetual(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		perpetual perptypes.Perpetual,
		err error,
	)
	GetPerpetualAndMarketPrice(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		perptypes.Perpetual,
		pricestypes.MarketPrice,
		error,
	)
	GetPerpetualAndMarketPriceAndLiquidityTier(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		perptypes.Perpetual,
		pricestypes.MarketPrice,
		perptypes.LiquidityTier,
		error,
	)
	GetLiquidityTier(
		ctx sdk.Context,
		id uint32,
	) (
		perptypes.LiquidityTier,
		error,
	)
	GetAllPerpetuals(ctx sdk.Context) []perptypes.Perpetual
	GetInsuranceFundName(ctx sdk.Context, perpetualId uint32) (string, error)
	GetInsuranceFundModuleAddress(ctx sdk.Context, perpetualId uint32) (sdk.AccAddress, error)
	IsIsolatedPerpetual(ctx sdk.Context, perpetualId uint32) (bool, error)
	ModifyOpenInterest(ctx sdk.Context, perpetualId uint32, bigQuantums *big.Int) error
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromAccountToModule(
		ctx context.Context,
		senderAddr sdk.AccAddress,
		recipientModule string,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToAccount(ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx context.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	BlockedAddr(addr sdk.AccAddress) bool
}

type BlocktimeKeeper interface {
	GetDowntimeInfoFor(ctx sdk.Context, duration time.Duration) blocktimetypes.AllDowntimeInfo_DowntimeInfo
}
