package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	prices "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

type PricesKeeper interface {
	GetMarketPrice(
		ctx sdk.Context,
		id uint32,
	) (market prices.MarketPrice, err error)
	// Methods imported from prices should be defined here
}
