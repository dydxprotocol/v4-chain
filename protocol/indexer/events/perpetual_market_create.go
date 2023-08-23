package events

// NewPerpetualMarketCreateEvent creates a PerpetualMarketCreateEvent
// representing creation of a perpetual market.
func NewPerpetualMarketCreateEvent(
	id uint32,
	symbol string,
	hasMarket bool,
	marketId uint32,
	atomicResolution int32,
) *PerpetualMarketCreateEventV1 {
	return &PerpetualMarketCreateEventV1{
		Id:               id,
		Symbol:           symbol,
		HasMarket:        hasMarket,
		MarketId:         marketId,
		AtomicResolution: atomicResolution,
	}
}
