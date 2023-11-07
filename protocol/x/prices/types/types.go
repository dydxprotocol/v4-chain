package types

import (
	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PricesKeeper interface {
	// Market related.
	CreateMarket(
		ctx sdk.Context,
		param MarketParam,
		price MarketPrice,
	) (createdMarketParam MarketParam, err error)

	ModifyMarketParam(
		ctx sdk.Context,
		param MarketParam,
	) (updatedMarketParam MarketParam, err error)

	UpdateMarketPrices(
		ctx sdk.Context,
		updates []*MsgUpdateMarketPrices_MarketPrice,
	) (err error)

	GetAllMarketParamPrices(ctx sdk.Context) (marketPramPrices []MarketParamPrice, err error)
	GetMarketParam(ctx sdk.Context, id uint32) (marketParam MarketParam, exists bool)
	GetMarketIdToValidIndexPrice(ctx sdk.Context) (marketIdToIndexPrice map[uint32]MarketPrice)
	GetAllMarketParams(ctx sdk.Context) (marketParams []MarketParam)
	GetMarketPrice(ctx sdk.Context, id uint32) (marketPrice MarketPrice, err error)
	GetAllMarketPrices(ctx sdk.Context) (marketPrices []MarketPrice)
	HasAuthority(authority string) bool

	// Validation related.
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *MsgUpdateMarketPrices,
		performNonDeterministicValidation bool,
	) error

	// Proposal related.
	UpdateSmoothedPrices(
		ctx sdk.Context,
		linearInterpolateFunc func(v0 uint64, v1 uint64, ppm uint32) (uint64, error),
	) error

	// Misc.
	Logger(ctx sdk.Context) log.Logger
}
