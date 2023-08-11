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
	// {k: exchangeId, v: market prices, read-write lock}
	ExchangeMarketPrices map[ExchangeId]*MarketToPrice
}

// NewExchangeToMarketPrices creates a new ExchangeToMarketPrices. Since `ExchangeToMarketPrices` is
// not goroutine safe to write to, all key-value store creation is done on initialization.
// Validation is also done to verify `exchangeIds` is a valid input.
func NewExchangeToMarketPrices(exchangeIds []ExchangeId) (*ExchangeToMarketPrices, error) {
	// Verify `ExchangeToMarketPrices` will not be initialized without `exchangeIds`.
	if len(exchangeIds) == 0 {
		return nil, errors.New("exchangeIds must not be empty")
	}

	exchangeToMarketPrices := &ExchangeToMarketPrices{
		ExchangeMarketPrices: make(map[ExchangeId]*MarketToPrice, len(exchangeIds)),
	}

	for _, exchangeId := range exchangeIds {
		// Verify no `exchangeIds` are duplicates.
		if _, ok := exchangeToMarketPrices.ExchangeMarketPrices[exchangeId]; ok {
			return nil, fmt.Errorf("exchangeId: '%v' appears twice in request", exchangeId)
		}

		exchangeToMarketPrices.ExchangeMarketPrices[exchangeId] = NewMarketToPrice()
	}

	return exchangeToMarketPrices, nil
}

// UpdatePrice updates a price for a market for an exchange. Prices are only updated if the
// timestamp on the updates are greater than the timestamp on existing prices. NOTE:
// `UpdatePrice` will only ever read from `ExchangeMarketPrices` and calls a
// goroutine-safe method on the fetched `MarketToPrice`.
// Note: if an invalid `exchangeId` is being written to the `UpdatePrice` it is possible the
// underlying map was corrupted or the priceDaemon logic is invalid. Therefore, `UpdatePrice`
// will panic.
func (exchangeToMarketPrices *ExchangeToMarketPrices) UpdatePrice(
	exchangeId ExchangeId,
	marketPriceTimestamp *MarketPriceTimestamp,
) {
	// Measure latency to update price in in-memory map.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.PriceEncoderUpdatePrice,
		metrics.Latency,
	)

	exchangeToMarketPrices.ExchangeMarketPrices[exchangeId].UpdatePrice(marketPriceTimestamp)
}

// GetAllPrices returns a map of exchangeIds to a list of all `MarketPriceTimestamps` for the exchange.
func (exchangeToMarketPrices *ExchangeToMarketPrices) GetAllPrices() map[ExchangeId][]MarketPriceTimestamp {
	// Measure latency to get all prices from in-memory map.
	defer telemetry.ModuleMeasureSince(
		metrics.PricefeedDaemon,
		time.Now(),
		metrics.GetAllPrices_MarketIdToPrice,
		metrics.Latency,
	)

	exchangeIdToPrices := make(
		map[ExchangeId][]MarketPriceTimestamp,
		len(exchangeToMarketPrices.ExchangeMarketPrices),
	)

	for exchangeId, mtp := range exchangeToMarketPrices.ExchangeMarketPrices {
		marketPrices := mtp.GetAllPrices()
		exchangeIdToPrices[exchangeId] = marketPrices
	}

	return exchangeIdToPrices
}
