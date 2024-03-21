package keeper

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
)

/*
 * This package implements the OracleKeeper interface from Slinky's currencypair package
 *
 * It is required in order to convert between x/prices types and the data which slinky stores on chain
 * via Vote Extensions. Using this compatibility layer, we can now use the x/prices keeper as the backing
 * store for converting to and from Slinky's on chain price data.
 */

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
	mps := k.GetAllMarketParams(ctx)
	for _, mp := range mps {
		mpCp, err := slinky.MarketPairToCurrencyPair(mp.Pair)
		if err != nil {
			k.Logger(ctx).Error("market param pair invalid format", "pair", mp.Pair)
			continue
		}
		if strings.EqualFold(mpCp.String(), cp.String()) {
			return uint64(mp.Id), true
		}
	}
	return 0, false
}

func (k Keeper) GetPriceForCurrencyPair(ctx sdk.Context, cp slinkytypes.CurrencyPair) (oracletypes.QuotePrice, error) {
	id, found := k.GetIDForCurrencyPair(ctx, cp)
	if !found {
		return oracletypes.QuotePrice{}, fmt.Errorf("id for currency pair %s not found", cp.String())
	}
	mp, err := k.GetMarketPrice(ctx, uint32(id))
	if err != nil {
		return oracletypes.QuotePrice{}, fmt.Errorf("market price not found for currency pair %s", cp.String())
	}
	return oracletypes.QuotePrice{
		Price: math.NewIntFromUint64(mp.Price),
	}, nil
}
