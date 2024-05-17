package types

import (
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/api"
	pricefeedmetrics "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	gometrics "github.com/hashicorp/go-metrics"
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

		// Measure invalid price updates inserted into the in-memory map.
		if exists && !isUpdated {
			telemetry.IncrCounterWithLabels(
				[]string{metrics.PricefeedServer, metrics.UpdatePrice, metrics.Invalid, metrics.Count},
				1,
				[]gometrics.Label{
					pricefeedmetrics.GetLabelForMarketId(etp.marketId),
					pricefeedmetrics.GetLabelForExchangeId(exchangeId),
				},
			)
		}
	}
}

// GetValidPrices returns a list of "valid" prices. Prices are considered
// "valid" iff the last update time is greater than or equal to the given cutoff time.
func (etp *ExchangeToPrice) GetValidPrices(
	logger log.Logger,
	cutoffTime time.Time,
) []uint64 {
	validExchangePricesForMarket := make([]uint64, 0, len(etp.exchangeToPriceTimestamp))
	for exchangeId, priceTimestamp := range etp.exchangeToPriceTimestamp {
		// PriceTimestamp returns price if the last update time is valid.
		if price, ok := priceTimestamp.GetValidPrice(cutoffTime); ok {
			validExchangePricesForMarket = append(validExchangePricesForMarket, price)
		} else {
			// Price is invalid.
			logger.Warn(
				"GetValidPrice returned invalid price. This likely means stale prices.",
				metrics.ExchangeId,
				exchangeId,
				metrics.MarketId,
				etp.marketId,
			)
		}
	}
	return validExchangePricesForMarket
}
