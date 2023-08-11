package keeper

import (
	indexerevents "github.com/dydxprotocol/v4/indexer/events"
	"github.com/dydxprotocol/v4/x/prices/types"
)

// GenerateMarketPriceUpdateEvents takes in a slice of markets and returns a slice of price updates.
func GenerateMarketPriceUpdateEvents(markets []types.Market) []*indexerevents.MarketEvent {
	events := make([]*indexerevents.MarketEvent, 0, len(markets))
	for _, market := range markets {
		events = append(
			events,
			indexerevents.NewMarketPriceUpdateEvent(
				market.Id,
				market.Price,
			),
		)
	}
	return events
}
