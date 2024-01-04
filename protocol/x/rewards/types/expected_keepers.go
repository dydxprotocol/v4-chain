package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assets "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoinsFromModuleToAccount(
		ctx context.Context,
		senderModule string,
		recipientAddr sdk.AccAddress,
		amt sdk.Coins,
	) error
}

type FeeTiersKeeper interface {
	GetLowestMakerFee(ctx sdk.Context) int32
}

type PricesKeeper interface {
	GetMarketPrice(
		ctx sdk.Context,
		id uint32,
	) (market prices.MarketPrice, err error)
}

type AssetsKeeper interface {
	GetAsset(
		ctx sdk.Context,
		id uint32,
	) (val assets.Asset, exists bool)
}
