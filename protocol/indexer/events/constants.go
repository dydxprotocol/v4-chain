package events

const (
	// Cosmos event attribute values for the subtype attribute for different indexer events.
	SubtypeOrderFill        = "order_fill"
	SubtypeSubaccountUpdate = "subaccount_update"
	SubtypeTransfer         = "transfer"
	SubtypeMarket           = "market"
	SubtypeFundingValues    = "funding_values"
	SubtypeStatefulOrder    = "stateful_order"
)

var OnChainEventSubtypes = []string{
	SubtypeOrderFill,
	SubtypeSubaccountUpdate,
	SubtypeTransfer,
	SubtypeMarket,
	SubtypeFundingValues,
	SubtypeStatefulOrder,
}
