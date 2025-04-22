package vote_extensions

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/dydxprotocol/slinky/pkg/types"

	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PricesKeeper is the expected interface for the x/price keeper used by the vote extension handlers
type PricesKeeper interface {
	GetCurrencyPairFromID(ctx sdk.Context, id uint64) (cp oracletypes.CurrencyPair, found bool)
	GetValidMarketPriceUpdates(ctx sdk.Context) *pricestypes.MsgUpdateMarketPrices
	UpdateMarketPrices(
		ctx sdk.Context,
		updates []*pricestypes.MsgUpdateMarketPrices_MarketPrice,
	) (err error)
}
