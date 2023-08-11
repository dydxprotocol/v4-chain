package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PricesKeeper interface {
	// Market related.
	CreateMarket(
		ctx sdk.Context,
		pair string,
		exponent int32,
		exchanges []uint32,
		minExchanges uint32,
		minPriceChangePpm uint32,
	) (createdMarket Market, err error)

	ModifyMarket(
		ctx sdk.Context,
		id uint32,
		pair string,
		exchanges []uint32,
		minExchanges uint32,
		minPriceChangePpm uint32,
	) (updatedMarket Market, err error)

	UpdateMarketPrices(
		ctx sdk.Context,
		updates []*MsgUpdateMarketPrices_MarketPrice,
		sendIndexerPriceUpdates bool,
	) (err error)

	GetMarket(ctx sdk.Context, id uint32) (market Market, err error)
	GetAllMarkets(ctx sdk.Context) (markets []Market)
	GetNumMarkets(ctx sdk.Context) (numMarkets uint32)

	// Exchange related.
	CreateExchangeFeed(
		ctx sdk.Context,
		name string,
		memo string,
	) (createdExchange ExchangeFeed, err error)

	ModifyExchangeFeed(
		ctx sdk.Context,
		id uint32,
		memo string,
	) (updatedExchange ExchangeFeed, err error)

	GetExchangeFeed(ctx sdk.Context, id uint32) (exchange ExchangeFeed, err error)
	GetAllExchangeFeeds(ctx sdk.Context) (exchanges []ExchangeFeed)
	GetNumExchangeFeeds(ctx sdk.Context) (numExchanges uint32)

	// Validation related.
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *MsgUpdateMarketPrices,
		performNonDeterministicValidation bool,
	) error

	// Proposal related.
	UpdateSmoothedPrices(
		ctx sdk.Context,
	) error
}
