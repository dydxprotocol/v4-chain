package types

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
}
