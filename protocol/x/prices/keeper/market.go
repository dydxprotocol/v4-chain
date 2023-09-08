package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// CreateMarket creates a new market param in the store along with a new market price
// for that market param.
// This is the only path to creating new MarketPrices, so if we have a param
// defined for a market, we should expect to see a price defined, and vice versa.
func (k Keeper) CreateMarket(
	ctx sdk.Context,
	marketParam types.MarketParam,
	marketPrice types.MarketPrice,
) (types.MarketParam, error) {
	// Validate input.
	if err := marketParam.Validate(); err != nil {
		return types.MarketParam{}, err
	}
	if err := marketPrice.ValidateFromParam(marketParam); err != nil {
		return types.MarketParam{}, err
	}

	paramBytes := k.cdc.MustMarshal(&marketParam)
	priceBytes := k.cdc.MustMarshal(&marketPrice)

	marketParamStore := k.newMarketParamStore(ctx)
	marketParamStore.Set(types.MarketKey(marketParam.Id), paramBytes)

	marketPriceStore := k.newMarketPriceStore(ctx)
	marketPriceStore.Set(types.MarketKey(marketPrice.Id), priceBytes)

	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeMarket,
		indexer_manager.GetB64EncodedEventMessage(
			indexerevents.NewMarketCreateEvent(
				marketParam.Id,
				marketParam.Pair,
				marketParam.MinPriceChangePpm,
				marketParam.Exponent,
			),
		),
	)
	return marketParam, nil
}

// GetAllMarketParamPrices returns a slice of MarketParam, MarketPrice tuples for all markets.
func (k Keeper) GetAllMarketParamPrices(ctx sdk.Context) ([]types.MarketParamPrice, error) {
	marketParams := k.GetAllMarketParams(ctx)
	marketPrices := k.GetAllMarketPrices(ctx)

	if len(marketParams) != len(marketPrices) {
		return nil, errorsmod.Wrap(types.ErrMarketPricesAndParamsDontMatch, "market param and price lengths do not match")
	}

	marketParamPrices := make([]types.MarketParamPrice, len(marketParams))
	for i, param := range marketParams {
		marketParamPrices[i].Param = param
		price := marketPrices[i]
		if param.Id != price.Id {
			return nil, errorsmod.Wrap(types.ErrMarketPricesAndParamsDontMatch,
				fmt.Sprintf("market param and price ids do not match: %d != %d", param.Id, price.Id))
		}
		marketParamPrices[i].Price = price
	}
	return marketParamPrices, nil
}

// GetNumMarkets returns the total number of markets.
// Panics if the length of market params and prices does not match.
func (k Keeper) GetNumMarkets(
	ctx sdk.Context,
) uint32 {
	marketParams := k.GetAllMarketParams(ctx)
	marketPrices := k.GetAllMarketPrices(ctx)

	if len(marketParams) != len(marketPrices) {
		panic(errorsmod.Wrap(types.ErrMarketPricesAndParamsDontMatch, "market param and price lengths do not match"))
	}

	return lib.MustConvertIntegerToUint32(len(marketParams))
}
