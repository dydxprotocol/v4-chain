package price_writer

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PriceApplierPricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *pricestypes.MarketPriceUpdate,
	) (isSpotValid bool, isPnlValid bool)

	GetAllMarketParams(ctx sdk.Context) []pricestypes.MarketParam

	UpdateSpotAndPnlMarketPrices(
		ctx sdk.Context,
		update *pricestypes.MarketPriceUpdate,
	) error

	UpdateSpotPrice(
		ctx sdk.Context,
		update *pricestypes.MarketSpotPriceUpdate,
	) error

	UpdatePnlPrice(
		ctx sdk.Context,
		update *types.MarketPnlPriceUpdate,
	) error

	GetMarketParam(
		ctx sdk.Context,
		id uint32,
	) (
		market pricestypes.MarketParam,
		exists bool,
	)
}
