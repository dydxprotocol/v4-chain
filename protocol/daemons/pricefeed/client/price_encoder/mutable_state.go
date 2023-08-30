package price_encoder

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"sync"
)

// mutableState stores all mutable state for the price encoder. Unlike the price fetcher, the
// price encoder needs access to multiple parts of the mutable market state, so we store it all.
type mutableState struct {
	sync.Mutex
	mutableExchangeConfig *types.MutableExchangeMarketConfig
	marketToMutableConfig map[types.MarketId]*types.MutableMarketConfig
}

// Update updates all fields in the mutableState atomically. Updates are required
// to be atomic in order to keep the mutableState consistent. This method expects
// validation to occur in the PriceEncoder. This method is synchronized.
func (ms *mutableState) Update(
	mutableExchangeConfig *types.MutableExchangeMarketConfig,
	marketToMutableConfig map[types.MarketId]*types.MutableMarketConfig,
) {
	ms.Lock()
	defer ms.Unlock()

	ms.mutableExchangeConfig = mutableExchangeConfig
	ms.marketToMutableConfig = marketToMutableConfig
}

// GetPriceConversionDetailsForMarket returns the price conversion details for the given market.
func (ms *mutableState) GetPriceConversionDetailsForMarket(id types.MarketId) (priceConversionDetails, error) {
	ms.Lock()
	defer ms.Unlock()

	marketConfig, ok := ms.mutableExchangeConfig.MarketToMarketConfig[id]
	if !ok {
		return priceConversionDetails{}, fmt.Errorf(
			"market config for market %v not found on exchange '%v'",
			id,
			ms.mutableExchangeConfig.Id,
		)
	}

	var adjustDetails *adjustByMarketDetails
	if marketConfig.AdjustByMarket != nil {
		adjustByMarketConfig, ok := ms.marketToMutableConfig[*marketConfig.AdjustByMarket]
		if !ok {
			return priceConversionDetails{}, fmt.Errorf(
				"mutable market config for adjust-by market %v not found on exchange '%v'",
				*marketConfig.AdjustByMarket,
				ms.mutableExchangeConfig.Id,
			)
		}
		adjustDetails = &adjustByMarketDetails{
			MarketId:     *marketConfig.AdjustByMarket,
			Exponent:     adjustByMarketConfig.Exponent,
			MinExchanges: adjustByMarketConfig.MinExchanges,
		}
	}

	mutableMarketConfig, ok := ms.marketToMutableConfig[id]
	if !ok {
		return priceConversionDetails{}, fmt.Errorf(
			"mutable market config for market %v not found on exchange '%v'",
			id,
			ms.mutableExchangeConfig.Id,
		)
	}
	return priceConversionDetails{
		Invert:                marketConfig.Invert,
		Exponent:              mutableMarketConfig.Exponent,
		AdjustByMarketDetails: adjustDetails,
	}, nil
}
