package types

import (
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

// ExchangeToPrice maintains multiple prices from different exchanges for
// the same market, along with the last time the each exchange price was updated.
type ExchangeToPrice struct {
	marketId                 uint32
	exchangeToPriceTimestamp map[string]*types.PriceTimestamp
}

// NewExchangeToPrice creates a new ExchangeToPrice. It takes a market ID, which is used in logging and metrics to
// identify the market these exchange prices are for. The market ID does not otherwise affect the behavior
// of the ExchangeToPrice.
func NewExchangeToPrice(marketId uint32) *ExchangeToPrice {
	return &ExchangeToPrice{
		marketId:                 marketId,
		exchangeToPriceTimestamp: make(map[string]*types.PriceTimestamp),
	}
}

// UpdatePrices updates prices given a list of prices from different exchanges.
// Prices are only updated if the timestamp on the updates are greater than
// the timestamp on existing prices.
func (etp *ExchangeToPrice) UpdatePrices(updates []*api.ExchangePrice) {
	for _, exchangePrice := range updates {
		exchangeId := exchangePrice.ExchangeId
		priceTimestamp, exists := etp.exchangeToPriceTimestamp[exchangeId]
		if !exists {
			priceTimestamp = types.NewPriceTimestamp()
			etp.exchangeToPriceTimestamp[exchangeId] = priceTimestamp
		}

		isUpdated := priceTimestamp.UpdatePrice(exchangePrice.Price, exchangePrice.LastUpdateTime)

		validity := metrics.Valid
		if exists && !isUpdated {
			validity = metrics.Invalid
		}

		// Measure count of valid and invalid prices inserted into the in-memory map.
		telemetry.IncrCounterWithLabels(
			[]string{metrics.PricefeedServer, metrics.UpdatePrice, validity, metrics.Count},
			1,
			[]gometrics.Label{
				pricefeedmetrics.GetLabelForMarketId(etp.marketId),
				pricefeedmetrics.GetLabelForExchangeId(exchangeId),
			},
		)
	}
}

// GetValidPrices returns a list of "valid" prices. Prices are considered
// "valid" iff the last update time is greater than or equal to the given cutoff time.
func (etp *ExchangeToPrice) GetValidPrices(
	cutoffTime time.Time,
) []uint64 {
	validExchangePricesForMarket := make([]uint64, 0, len(etp.exchangeToPriceTimestamp))
	for exchangeId, priceTimestamp := range etp.exchangeToPriceTimestamp {
		validity := metrics.Valid

		// PriceTimestamp returns price if the last update time is valid.
		if price, ok := priceTimestamp.GetValidPrice(cutoffTime); ok {
			validExchangePricesForMarket = append(validExchangePricesForMarket, price)
		} else {
			// Price is invalid.
			validity = metrics.PriceIsInvalid
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
				pricefeedmetrics.GetLabelForExchangeId(exchangeId),
				pricefeedmetrics.GetLabelForMarketId(etp.marketId),
			},
		)
	}
	return validExchangePricesForMarket
}
