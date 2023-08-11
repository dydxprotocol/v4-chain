package events

// NewMarketPriceUpdateEvent creates a MarketEvent representing an update in the priceWithExponent of a market.
func NewMarketPriceUpdateEvent(
	marketId uint32,
	priceWithExponent uint64,
) *MarketEvent {
	priceUpdateEventProto := MarketPriceUpdateEvent{
		PriceWithExponent: priceWithExponent,
	}
	return &MarketEvent{
		MarketId: marketId,
		Event: &MarketEvent_PriceUpdate{
			PriceUpdate: &priceUpdateEventProto,
		},
	}
}

// NewMarketModifyEvent creates a MarketEvent representing an update to a market.
func NewMarketModifyEvent(
	marketId uint32,
	pair string,
	minPriceChangePpm uint32,
) *MarketEvent {
	marketModifyEventProto := MarketModifyEvent{
		Base: &MarketBaseEvent{
			Pair:              pair,
			MinPriceChangePpm: minPriceChangePpm,
		},
	}
	marketEventProto := MarketEvent{
		MarketId: marketId,
		Event: &MarketEvent_MarketModify{
			MarketModify: &marketModifyEventProto,
		},
	}
	return &marketEventProto
}

// NewMarketCreateEvent creates a MarketEvent representing a new market.
func NewMarketCreateEvent(
	marketId uint32,
	pair string,
	minPriceChangePpm uint32,
	exponent int32,
) *MarketEvent {
	marketCreateEventProto := MarketCreateEvent{
		Base: &MarketBaseEvent{
			Pair:              pair,
			MinPriceChangePpm: minPriceChangePpm,
		},
		Exponent: exponent,
	}
	return &MarketEvent{
		MarketId: marketId,
		Event: &MarketEvent_MarketCreate{
			MarketCreate: &marketCreateEventProto,
		},
	}
}
