package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	prices "github.com/dydxprotocol/v4/x/prices/types"
)

type PricesKeeper interface {
	GetMarket(
		ctx sdk.Context,
		id uint32,
	) (market prices.Market, err error)
	// Methods imported from prices should be defined here
}
