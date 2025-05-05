package types

import (
	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"
	oracletypes "github.com/dydxprotocol/slinky/x/oracle/types"
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
	GetExponent(ctx sdk.Context, ticker string) (exponent int32, err error)
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

	// Currency Pair ID cache
	AddCurrencyPairIDToStore(ctx sdk.Context, id uint32, cp slinkytypes.CurrencyPair)

	// Slinky compat
	GetCurrencyPairFromID(ctx sdk.Context, id uint64) (cp slinkytypes.CurrencyPair, found bool)
	GetIDForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (uint64, bool)
	GetPriceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (oracletypes.QuotePrice, error)

	SetNextMarketID(ctx sdk.Context, nextID uint32)
}
