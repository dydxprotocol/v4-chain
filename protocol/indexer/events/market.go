package events

// NewMarketPriceUpdateEvent creates a MarketEvent representing an update in the priceWithExponent of a market.
func NewMarketPriceUpdateEvent(
	marketId uint32,
	priceWithExponent uint64,
) *MarketEventV1 {
	priceUpdateEventProto := MarketPriceUpdateEventV1{
		PriceWithExponent: priceWithExponent,
	}
	return &MarketEventV1{
		MarketId: marketId,
		Event: &MarketEventV1_PriceUpdate{
			PriceUpdate: &priceUpdateEventProto,
		},
	}
}

// NewMarketModifyEvent creates a MarketEvent representing an update to a market.
func NewMarketModifyEvent(
	marketId uint32,
	pair string,
	minPriceChangePpm uint32,
) *MarketEventV1 {
	marketModifyEventProto := MarketModifyEventV1{
		Base: &MarketBaseEventV1{
			Pair:              pair,
			MinPriceChangePpm: minPriceChangePpm,
		},
	}
	marketEventProto := MarketEventV1{
		MarketId: marketId,
		Event: &MarketEventV1_MarketModify{
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
) *MarketEventV1 {
	marketCreateEventProto := MarketCreateEventV1{
		Base: &MarketBaseEventV1{
			Pair:              pair,
			MinPriceChangePpm: minPriceChangePpm,
		},
		Exponent: exponent,
	}
	return &MarketEventV1{
		MarketId: marketId,
		Event: &MarketEventV1_MarketCreate{
			MarketCreate: &marketCreateEventProto,
		},
	}
}
