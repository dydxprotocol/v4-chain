package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/lib/metrics"
)

// ExchangeToMarketPrices maintains price info for multiple exchanges. Each exchange can support
// prices from multiple market sources. Methods are goroutine safe in the underlying MarketToPrice
// objects.
type ExchangeToMarketPrices struct {
	// {k: exchangeFeedId, v: market prices, read-write lock}
	ExchangeMarketPrices map[ExchangeFeedId]*MarketToPrice
}

// NewExchangeToMarketPrices creates a new ExchangeToMarketPrices. Since `ExchangeToMarketPrices` is
// not goroutine safe to write to, all key-value store creation is done on initialization.
// Validation is also done to verify `exchangeFeedIds` is a valid input.
func NewExchangeToMarketPrices(exchangeFeedIds []ExchangeFeedId) (*ExchangeToMarketPrices, error) {
	// Verify `ExchangeToMarketPrices` will not be initialized without `exchangeFeedIds`.
	if len(exchangeFeedIds) == 0 {
		return nil, errors.New("exchangeFeedIds must not be empty")
	}

	exchangeToMarketPrices := &ExchangeToMarketPrices{
		ExchangeMarketPrices: make(map[MarketId]*MarketToPrice, len(exchangeFeedIds)),
	}

	for _, exchangeFeedId := range exchangeFeedIds {
		// Verify no `exchangeFeedIds` are duplicates.
		if _, ok := exchangeToMarketPrices.ExchangeMarketPrices[exchangeFeedId]; ok {
			return nil, fmt.Errorf("exchangeFeedId: %d appears twice in request", exchangeFeedId)
		}

		exchangeToMarketPrices.ExchangeMarketPrices[exchangeFeedId] = NewMarketToPrice()
	}

	return exchangeToMarketPrices, nil
}

// UpdatePrice updates a price for a market for an exchange. Prices are only updated if the
// timestamp on the updates are greater than the timestamp on existing prices. NOTE:
// `UpdatePrice` will only ever read from `ExchangeMarketPrices` and calls a
// goroutine-safe method on the fetched `MarketToPrice`.
// Note: if an invalid `exchangeFeedId` is being written to the `UpdatePrice` it is possible the
// underlying map was corrupted or the priceDaemon logic is invalid. Therefore, `UpdatePrice`
// will panic.
func (exchangeToMarketPrices *ExchangeToMarketPrices) UpdatePrice(
	exchangeFeedId ExchangeFeedId,
	marketPriceTimestamp *MarketPriceTimestamp,
) {
	// Measure latency to update price in in-memory map.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.PriceEncoderUpdatePrice,
		metrics.Latency,
	)

	exchangeToMarketPrices.ExchangeMarketPrices[exchangeFeedId].UpdatePrice(marketPriceTimestamp)
}

// GetAllPrices returns a map of exchangeFeedIds to a list of all `MarketPriceTimestamps` for the exchange.
func (exchangeToMarketPrices *ExchangeToMarketPrices) GetAllPrices() map[MarketId][]MarketPriceTimestamp {
	// Measure latency to get all prices from in-memory map.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.GetAllPrices_MarketIdToPrice,
		metrics.Latency,
	)

	exchangeFeedIdToPrices := make(
		map[MarketId][]MarketPriceTimestamp,
		len(exchangeToMarketPrices.ExchangeMarketPrices),
	)

	for exchangeFeedId, mtp := range exchangeToMarketPrices.ExchangeMarketPrices {
		marketPrices := mtp.GetAllPrices()
		exchangeFeedIdToPrices[exchangeFeedId] = marketPrices
	}

	return exchangeFeedIdToPrices
}
