package types

import (
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/daemons/pricefeed/types"
	"github.com/dydxprotocol/v4/lib/metrics"
)

// ExchangeToPrice maintains multiple prices from different exchanges for
// the same market, along with the last time the each exchange price was updated.
type ExchangeToPrice struct {
	exchangeToPriceTimestamp map[uint32]*types.PriceTimestamp
}

// NewExchangeToPrice creates a new ExchangeToPrice
func NewExchangeToPrice() *ExchangeToPrice {
	return &ExchangeToPrice{
		exchangeToPriceTimestamp: make(map[uint32]*types.PriceTimestamp),
	}
}

// UpdatePrices updates prices given a list of prices from different exchanges.
// Prices are only updated if the timestamp on the updates are greater than
// the timestamp on existing prices.
func (etp *ExchangeToPrice) UpdatePrices(updates []*api.ExchangePrice) {
	for _, exchangePrice := range updates {
		exchangeFeedId := exchangePrice.ExchangeFeedId
		priceTimestamp, ok := etp.exchangeToPriceTimestamp[exchangeFeedId]
		if !ok {
			priceTimestamp = types.NewPriceTimestamp()
			etp.exchangeToPriceTimestamp[exchangeFeedId] = priceTimestamp
		}

		isUpdated := priceTimestamp.UpdatePrice(exchangePrice.Price, exchangePrice.LastUpdateTime)

		validity := metrics.Valid
		if !isUpdated {
			validity = metrics.Invalid
		}

		// Measure count of valid and invalid prices inserted into the in-memory map.
		telemetry.IncrCounter(1, metrics.PricefeedServer, metrics.UpdatePrice, validity, metrics.Count)
	}
}

// GetValidPrices returns a list of "valid" prices. Prices are considered
// "valid" iff
// 1) the price's source exchange is within a given set of valid exchange IDs.
// 2) the last update time is greater than or equal to the given cutoff time.
func (etp *ExchangeToPrice) GetValidPrices(
	validExchangeFeedIds map[uint32]bool,
	cutoffTime time.Time,
) []uint64 {
	validExchangePricesForMarket := make([]uint64, 0, len(etp.exchangeToPriceTimestamp))
	for exchangeFeedId, priceTimestamp := range etp.exchangeToPriceTimestamp {
		validity := metrics.Valid

		if _, valid := validExchangeFeedIds[exchangeFeedId]; valid {
			// PriceTimestamp returns price if the last update time is valid.
			if price, ok := priceTimestamp.GetValidPrice(cutoffTime); ok {
				validExchangePricesForMarket =
					append(validExchangePricesForMarket, price)
			} else {
				// Price is invalid.
				validity = metrics.PriceIsInvalid
			}
		} else {
			// ExchangeFeedId is invalid.
			validity = metrics.ExchangeFeedIsInvalid
		}

		// Measure count of valid and invalid prices fetched from the in-memory map.
		telemetry.IncrCounterWithLabels(
			[]string{
				metrics.PricefeedServer,
				metrics.GetValidPrices,
				validity,
				metrics.Count,
			},
			1,
			[]gometrics.Label{
				pricefeedmetrics.GetLabelForExchangeFeedId(exchangeFeedId),
			},
		)
	}
	return validExchangePricesForMarket
}
