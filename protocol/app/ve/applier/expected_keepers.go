package price_writer

import (
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PriceApplierPricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *pricestypes.MarketPriceUpdates,
	) error

	GetAllMarketParams(ctx sdk.Context) []pricestypes.MarketParam

	UpdateMarketPrice(
		ctx sdk.Context,
		update *pricestypes.MarketPriceUpdates_MarketPriceUpdate,
	) error
	GetMarketParam(
		ctx sdk.Context,
		id uint32,
	) (
		market pricestypes.MarketParam,
		exists bool,
	)
}
