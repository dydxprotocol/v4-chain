package keeper

import (
	indexerevents "github.com/StreamFinance-Protocol/stream-chain/protocol/indexer/events"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
)

// GenerateMarketPriceUpdateIndexerEvents takes in a slice of market prices
// and returns a slice of price updates.
func GenerateMarketPriceUpdateIndexerEvent(
	market types.MarketPrice,
) *indexerevents.MarketEventV1 {
	return indexerevents.NewMarketPriceUpdateEvent(
		market.Id,
		market.Price,
	)
}
