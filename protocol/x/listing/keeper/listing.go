package keeper

import (
	"math"

	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/skip-mev/slinky/x/marketmap/types/tickermetadata"
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
	// TODO: change to use util from marketmap when available
	marketMapPair, err := slinky.MarketPairToCurrencyPair(ticker)
	if err != nil {
		return 0, err
	}
	marketMapDetails, err := k.MarketMapKeeper.GetMarket(ctx, marketMapPair.String())
	if err != nil {
		return 0, types.ErrMarketNotFound
	}

	// Create a new market
	market, err := k.PricesKeeper.CreateMarket(
		ctx,
		pricestypes.MarketParam{
			Id:   marketId,
			Pair: ticker,
			// Set the price exponent to the negative of the number of decimals
			Exponent:           int32(marketMapDetails.Ticker.Decimals) * -1,
			MinExchanges:       uint32(marketMapDetails.Ticker.MinProviderCount),
			MinPriceChangePpm:  types.MinPriceChangePpm_LongTail,
			ExchangeConfigJson: "{}", // Placeholder. TODO (TRA-513): Deprecate this field
		},
		pricestypes.MarketPrice{
			Id:       marketId,
			Exponent: int32(marketMapDetails.Ticker.Decimals) * -1,
			Price:    0,
		},
	)
	if err != nil {
		return 0, err
	}

	return market.Id, nil
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

// Function to wrap the creation of a new perpetual
// Note: This will only list long-tail/isolated markets
func (k Keeper) CreatePerpetual(
	ctx sdk.Context,
	marketId uint32,
	ticker string,
) (perpetualId uint32, err error) {
	perpetualId = k.PerpetualsKeeper.AcquireNextPerpetualID(ctx)

	// Get reference price from market map
	// TODO: change to use util from marketmap when available
	marketMapPair, err := slinky.MarketPairToCurrencyPair(ticker)
	if err != nil {
		return 0, err
	}
	marketMapDetails, err := k.MarketMapKeeper.GetMarket(ctx, marketMapPair.String())
	if err != nil {
		return 0, err
	}
	metadata, err := tickermetadata.DyDxFromJSONString(marketMapDetails.Ticker.Metadata_JSON)
	if err != nil {
		return 0, err
	}
	if metadata.ReferencePrice == 0 {
		return 0, types.ErrReferencePriceZero
	}

	// calculate atomic resolution from reference price
	atomicResolution := types.ResolutionOffset - int32(math.Floor(math.Log10(float64(metadata.ReferencePrice))))

	// Create a new perpetual
	perpetual, err := k.PerpetualsKeeper.CreatePerpetual(
		ctx,
		perpetualId,
		ticker,
		marketId,
		atomicResolution,
		types.DefaultFundingPpm,
		types.LiquidityTier_Isolated,
		perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED,
	)
	if err != nil {
		return 0, err
	}

	return perpetual.GetId(), nil
}

// Function to set listing vault deposit params in module store
func (k Keeper) SetListingVaultDepositParams(
	ctx sdk.Context,
	params types.ListingVaultDepositParams,
) error {
	// Validate the params
	if err := params.Validate(); err != nil {
		return err
	}

	// Store the params in the module store
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.ListingVaultDepositParamsKey)
	store.Set(key, k.cdc.MustMarshal(&params))
	return nil
}

// Function to get listing vault deposit params from module store
func (k Keeper) GetListingVaultDepositParams(
	ctx sdk.Context,
) (vaultDepositParams types.ListingVaultDepositParams) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.ListingVaultDepositParamsKey))
	k.cdc.MustUnmarshal(b, &vaultDepositParams)
	return vaultDepositParams
}
