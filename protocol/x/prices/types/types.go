package types

import (
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
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

	GetValidMarketPriceUpdates(
		ctx sdk.Context,
	) *MsgUpdateMarketPrices

	// Misc.
	Logger(ctx sdk.Context) log.Logger

	// Slinky compat
	GetCurrencyPairFromID(ctx sdk.Context, id uint64) (cp slinkytypes.CurrencyPair, found bool)
	GetIDForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (uint64, bool)
	GetPriceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (oracletypes.QuotePrice, error)
	GetPrevBlockCPCounter(ctx sdk.Context) (uint64, error)
}
