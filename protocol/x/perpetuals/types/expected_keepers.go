package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	epochstypes "github.com/dydxprotocol/v4/x/epochs/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
)

type PricesKeeper interface {
	GetMarket(
		ctx sdk.Context,
		id uint32,
	) (market pricestypes.Market, err error)
}

// PricePremiumGetter defines an interface that returns price premium for a perpetual.
type PricePremiumGetter interface {
	GetPricePremiumForPerpetual(
		ctx sdk.Context,
		perpetualId uint32,
		params GetPricePremiumParams,
	) (
		premiumPpm int32,
		err error,
	)
}

// AccountKeeper defines the expected account keeper used for simulations.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
}

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// EpochsKeeper defines the expected epochs keeper to get epoch info.
type EpochsKeeper interface {
	NumBlocksSinceEpochStart(
		ctx sdk.Context,
		id epochstypes.EpochInfoName,
	) (uint32, error)
	MustGetFundingTickEpochInfo(
		ctx sdk.Context,
	) epochstypes.EpochInfo
	MustGetFundingSampleEpochInfo(
		ctx sdk.Context,
	) epochstypes.EpochInfo
}
