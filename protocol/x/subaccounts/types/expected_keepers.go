package types

import (
	"context"
	"math/big"
	"time"

	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ProductKeeper represents a generic interface for a keeper
// of a product.
type ProductKeeper interface {
	GetNetCollateral(
		ctx sdk.Context,
		id uint32,
		bigQuantums *big.Int,
	) (
		bigNetCollateralQuoteQuantums *big.Int,
		err error,
	)
	GetMarginRequirements(
		ctx sdk.Context,
		id uint32,
		bigQuantums *big.Int,
	) (
		bigInitialMarginQuoteQuantums *big.Int,
		bigMaintenanceMarginQuoteQuantums *big.Int,
		err error,
	)
	IsPositionUpdatable(
		ctx sdk.Context,
		id uint32,
	) (
		updatable bool,
		err error,
	)
}

type AssetsKeeper interface {
	ProductKeeper
	ConvertAssetToCoin(
		ctx sdk.Context,
		assetId uint32,
		quantums *big.Int,
	) (
		convertedQuantums *big.Int,
		coin sdk.Coin,
		err error,
	)
	ConvertCoinToAsset(ctx sdk.Context, assetId uint32, coin sdk.Coin) (quantums *big.Int, convertedDenom *big.Int, err error)
	ConvertAssetToFullCoin(
		ctx sdk.Context,
		assetId uint32,
		quantums *big.Int,
	) (
		convertedQuantums *big.Int,
		fullCoinAmount *big.Int,
		err error,
	)
}

type PerpetualsKeeper interface {
	ProductKeeper
	GetSettlementPpm(
		ctx sdk.Context,
		perpetualId uint32,
		quantums *big.Int,
		index *big.Int,
	) (
		bigNetSettlement *big.Int,
		newFundingIndex *big.Int,
		err error,
	)
	GetPerpetual(
		ctx sdk.Context,
		perpetualId uint32,
	) (
		perpetual perptypes.Perpetual,
		err error,
	)
	GetAllPerpetuals(ctx sdk.Context) []perptypes.Perpetual
	GetInsuranceFundName(ctx sdk.Context, perpetualId uint32) (string, error)
	GetInsuranceFundModuleAddress(ctx sdk.Context, perpetualId uint32) (sdk.AccAddress, error)
	ModifyOpenInterest(ctx sdk.Context, perpetualId uint32, bigQuantums *big.Int) error
	IsIsolatedPerpetual(ctx sdk.Context, perpetualId uint32) (bool, error)
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
}

type BlocktimeKeeper interface {
	GetDowntimeInfoFor(ctx sdk.Context, duration time.Duration) blocktimetypes.AllDowntimeInfo_DowntimeInfo
}

type RatelimitKeeper interface {
	GetAssetYieldIndex(ctx sdk.Context) (yieldIndex *big.Rat, found bool)
}
