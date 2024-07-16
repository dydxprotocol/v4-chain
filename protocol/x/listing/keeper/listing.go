package keeper

import (
	gogotypes "github.com/cosmos/gogoproto/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

// Function to set hard cap on listed markets in module store
func (k Keeper) SetMarketsHardCap(ctx sdk.Context, hardCap uint32) error {
	store := ctx.KVStore(k.storeKey)
	value := gogotypes.UInt32Value{Value: hardCap}
	store.Set([]byte(types.HardCapForMarketsKey), k.cdc.MustMarshal(&value))
	return nil
}

// Function to get hard cap on listed markets from module store
func (k Keeper) GetMarketsHardCap(ctx sdk.Context) (hardCap uint32) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.HardCapForMarketsKey))
	var result gogotypes.UInt32Value
	k.cdc.MustUnmarshal(b, &result)
	return result.Value
}

// Function to wrap the creation of a new market
// Note: This will only list long-tail/isolated markets
func (k Keeper) CreateMarket(
	ctx sdk.Context,
	ticker string,
) (marketId uint32, err error) {
	marketId = k.PricesKeeper.AcquireNextMarketID(ctx)

	// Get market details from marketmap
	market, err := k.MarketMapKeeper.GetMarket(ctx, ticker)
	if err != nil {
		return 0, err
	}

	// Create a new market
	_, err = k.PricesKeeper.CreateMarket(
		ctx,
		pricestypes.MarketParam{
			Id:   marketId,
			Pair: ticker,
			// Set the price exponent to the negative of the number of decimals
			Exponent: int32(market.Ticker.Decimals * -1),
			MinExchanges: uint32(market.Ticker.MinProviderCount),
			MinPriceChangePpm:
		},
		pricestypes.MarketPrice{
			Id: marketId,
		},


}

// Function to wrap the creation of a new clob pair
// Note: This will only list long-tail/isolated markets
func (k Keeper) CreateClobPair(
	ctx sdk.Context,
	perpetualId uint32,
) (clobPairId uint32, err error) {
	clobPairId = k.ClobKeeper.AcquireNextClobPairID(ctx)

	// Create a new clob pair
	clobPair, err := k.ClobKeeper.CreatePerpetualClobPair(
		ctx,
		clobPairId,
		perpetualId,
		satypes.BaseQuantums(types.DefaultStepBaseQuantums),
		types.DefaultQuantumConversionExponent,
		types.SubticksPerTick_LongTail,
		clobtypes.ClobPair_STATUS_ACTIVE,
	)
	if err != nil {
		return 0, err
	}

	return clobPair.Id, nil
}

