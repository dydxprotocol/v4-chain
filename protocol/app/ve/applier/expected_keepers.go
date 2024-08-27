package price_writer

import (
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PriceApplierPricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *types.MarketPriceUpdate,
	) (isSpotValid bool, isPnlValid bool)

	GetAllMarketParams(ctx sdk.Context) []types.MarketParam

	UpdateSpotAndPnlMarketPrices(
		ctx sdk.Context,
		update *types.MarketPriceUpdate,
	) error

	UpdateSpotPrice(
		ctx sdk.Context,
		update *types.MarketSpotPriceUpdate,
	) error

	UpdatePnlPrice(
		ctx sdk.Context,
		update *types.MarketPnlPriceUpdate,
	) error

	GetMarketParam(
		ctx sdk.Context,
		id uint32,
	) (
		market types.MarketParam,
		exists bool,
	)
}
