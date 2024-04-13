package types

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// ProductKeeper represents a generic interface for a keeper
// of a product.
type ProductKeeper interface {
	GetNetCollateral(
		ctx sdk.Context,
		id uint32,
		quantums *int256.Int,
	) (
		netCollateralQuoteQuantums *int256.Int,
		err error,
	)
	GetMarginRequirements(
		ctx sdk.Context,
		id uint32,
		Quantumsuantums *int256.Int,
	) (
		initialMarginQuoteQuantums *int256.Int,
		maintenanceMarginQuoteQuantums *int256.Int,
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
		quantums *int256.Int,
	) (
		convertedQuantums *int256.Int,
		coin sdk.Coin,
		err error,
	)
}

type PerpetualsKeeper interface {
	ProductKeeper
	GetSettlementPpm(
		ctx sdk.Context,
		perpetualId uint32,
		quantums *int256.Int,
		index *int256.Int,
	) (
		netSettlement *int256.Int,
		newFundingIndex *int256.Int,
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
	IsIsolatedPerpetual(ctx sdk.Context, perpetualId uint32) (bool, error)
	ModifyOpenInterest(ctx sdk.Context, perpetualId uint32, quantums *int256.Int) error
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
}

type BlocktimeKeeper interface {
	GetDowntimeInfoFor(ctx sdk.Context, duration time.Duration) blocktimetypes.AllDowntimeInfo_DowntimeInfo
}
