package aggregator

import (
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PriceApplierPricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *pricestypes.MarketPriceUpdates,
		performNonDeterministicValidation bool,
	) error

	GetValidMarketPriceUpdates(ctx sdk.Context) *pricestypes.MarketPriceUpdates
	GetAllMarketParams(ctx sdk.Context) []pricestypes.MarketParam
	GetMarketPriceUpdateFromBytes(id uint32, bz []byte) (*pricestypes.MarketPriceUpdates_MarketPriceUpdate, error)

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
