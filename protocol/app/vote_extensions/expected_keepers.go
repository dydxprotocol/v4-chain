package vote_extensions

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type PricesKeeper interface {
	UpdateSmoothedPrices(
		ctx sdk.Context,
		linearInterpolateFunc func(v0 uint64, v1 uint64, ppm uint32) (uint64, error),
	) error
	UpdateMarketPrices(
		ctx sdk.Context,
		updates []*pricestypes.MsgUpdateMarketPrices_MarketPrice,
	) (err error)
}
