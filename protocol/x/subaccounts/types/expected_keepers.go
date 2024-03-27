package types

import (
	"context"
	"math/big"
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
		bigQuantums *big.Int,
	) (
		bigNetCollateralQuoteQuantums *big.Int,
		err error,
	)
	GetNetCollateralInt256(
		ctx sdk.Context,
		id uint32,
		bigQuantums *int256.Int,
	) (
		bigNetCollateralQuoteQuantums *int256.Int,
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
	GetMarginRequirementsInt256(
		ctx sdk.Context,
		id uint32,
		bigQuantums *int256.Int,
	) (
		bigInitialMarginQuoteQuantums *int256.Int,
		bigMaintenanceMarginQuoteQuantums *int256.Int,
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
	IsIsolatedPerpetual(ctx sdk.Context, perpetualId uint32) (bool, error)
	ModifyOpenInterest(ctx sdk.Context, perpetualId uint32, bigQuantums *big.Int) error
	ModifyOpenInterestInt256(ctx sdk.Context, perpetualId uint32, quantums *int256.Int) error
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
