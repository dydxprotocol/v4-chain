package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"strings"
)

/*
 * This package implements the OracleKeeper interface from Slinky's currencypair package
 *
 * It is required in order to convert between x/prices types and the data which slinky stores on chain
 * via Vote Extensions. Using this compatibility layer, we can now use the x/prices keeper as the backing
 * store for converting to and from Slinky's on chain price data.
 */

func (k Keeper) GetAllCurrencyPairs(ctx sdk.Context) []oracletypes.CurrencyPair {
	mps := k.GetAllMarketParams(ctx)
	cps := make([]oracletypes.CurrencyPair, len(mps))
	for i, mp := range mps {
		cp, err := oracletypes.CurrencyPairFromString(mp.Pair)
		if err != nil {
			cps[i] = cp
		}
	}
	return cps
}

func (k Keeper) GetCurrencyPairFromID(ctx sdk.Context, id uint64) (cp oracletypes.CurrencyPair, found bool) {
	mp, found := k.GetMarketParam(ctx, uint32(id))
	if !found {
		return cp, false
	}
	cp, err := oracletypes.CurrencyPairFromString(mp.Pair)
	if err != nil {
		return cp, false
	}
	return cp, true
}

func (k Keeper) GetIDForCurrencyPair(ctx sdk.Context, cp oracletypes.CurrencyPair) (uint64, bool) {
	mps := k.GetAllMarketParams(ctx)
	for _, mp := range mps {
		if strings.EqualFold(mp.Pair, cp.String()) {
			return uint64(mp.Id), true
		}
	}
	return 0, false
}

func (k Keeper) GetPriceForCurrencyPair(ctx sdk.Context, cp oracletypes.CurrencyPair) (oracletypes.QuotePrice, error) {
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
