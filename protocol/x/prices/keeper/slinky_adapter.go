package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"
	oracletypes "github.com/dydxprotocol/slinky/x/oracle/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

/*
 * This package implements the OracleKeeper interface from Slinky's currencypair package
 *
 * It is required in order to convert between x/prices types and the data which slinky stores on chain
 * via Vote Extensions. Using this compatibility layer, we can now use the x/prices keeper as the backing
 * store for converting to and from Slinky's on chain price data.
 */

// getCurrencyPairIDStore returns a prefix store for market IDs corresponding to currency pairs.
func (k Keeper) getCurrencyPairIDStore(ctx sdk.Context) prefix.Store {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.CurrencyPairIDPrefix))
}

func (k Keeper) GetCurrencyPairFromID(ctx sdk.Context, id uint64) (cp slinkytypes.CurrencyPair, found bool) {
	mp, found := k.GetMarketParam(ctx, uint32(id))
	if !found {
		return cp, false
	}

	cp, err := slinky.MarketPairToCurrencyPair(mp.Pair)
	if err != nil {
		k.Logger(ctx).Error("CurrencyPairFromString", "error", err)
		return cp, false
	}

	return cp, true
}

func (k Keeper) GetIDForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (uint64, bool) {
	// Try to get corresponding market ID of the currency pair from the store
	marketId, found := k.GetCurrencyPairIDFromStore(ctx, cp)
	if found {
		return uint64(marketId), true
	}

	return 0, false
}

func (k Keeper) GetCurrencyPairIDFromStore(ctx sdk.Context, cp slinkytypes.CurrencyPair) (marketId uint64, found bool) {
	currencyPairString := cp.String()
	currencyPairIDStore := k.getCurrencyPairIDStore(ctx)
	var result gogotypes.UInt64Value
	b := currencyPairIDStore.Get([]byte(currencyPairString))
	if b == nil {
		return 0, false
	} else {
		k.cdc.MustUnmarshal(b, &result)
		return result.Value, true
	}
}

func (k Keeper) AddCurrencyPairIDToStore(ctx sdk.Context, marketId uint32, cp slinkytypes.CurrencyPair) {
	currencyPairString := cp.String()
	currencyPairIDStore := k.getCurrencyPairIDStore(ctx)
	value := gogotypes.UInt64Value{Value: uint64(marketId)}
	b := k.cdc.MustMarshal(&value)
	currencyPairIDStore.Set([]byte(currencyPairString), b)
}

func (k Keeper) RemoveCurrencyPairFromStore(ctx sdk.Context, cp slinkytypes.CurrencyPair) {
	currencyPairString := cp.String()
	currencyPairIDStore := k.getCurrencyPairIDStore(ctx)
	currencyPairIDStore.Delete([]byte(currencyPairString))
}

func (k Keeper) GetPriceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (oracletypes.QuotePrice, error) {
	id, found := k.GetIDForCurrencyPair(ctx, cp)
	if !found {
		return oracletypes.QuotePrice{}, fmt.Errorf("currency pair %s not found", cp.String())
	}
	mp, err := k.GetMarketPrice(ctx, uint32(id))
	if err != nil {
		return oracletypes.QuotePrice{}, fmt.Errorf("currency pair %s not found", cp.String())
	}
	return oracletypes.QuotePrice{
		Price: math.NewIntFromUint64(mp.Price),
	}, nil
}

func (k Keeper) GetNumCurrencyPairs(ctx sdk.Context) (uint64, error) {
	marketPriceStore := k.getMarketPriceStore(ctx)

	var numMarketPrices uint64

	iterator := marketPriceStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		marketPrice := types.MarketPrice{}
		k.cdc.MustUnmarshal(iterator.Value(), &marketPrice)
		numMarketPrices++
	}

	return numMarketPrices, nil
}

// GetNumRemovedCurrencyPairs is currently a no-op since we don't support removing Markets right now.
func (k Keeper) GetNumRemovedCurrencyPairs(_ sdk.Context) (uint64, error) {
	return 0, nil
}

// GetAllCurrencyPairs is not used with the DefaultCurrencyPair strategy.
// See https://github.com/dydxprotocol/slinky/blob/main/abci/strategies/currencypair/default.go
func (k Keeper) GetAllCurrencyPairs(_ sdk.Context) []slinkytypes.CurrencyPair {
	return nil
}
