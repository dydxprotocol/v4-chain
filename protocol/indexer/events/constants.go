package events

const (
	// Cosmos event attribute values for the subtype attribute for different indexer events.
	// Keep these constants in sync with:
	// https://github.com/dydxprotocol/indexer/blob/master/services/ender/src/lib/types.ts.
	// Ender uses these to maintain a mapping between event type and event proto.
	SubtypeOrderFill        = "order_fill"
	SubtypeSubaccountUpdate = "subaccount_update"
	SubtypeTransfer         = "transfer"
	SubtypeMarket           = "market"
	SubtypeFundingValues    = "funding_values"
	SubtypeStatefulOrder    = "stateful_order"
	SubtypeAsset            = "asset"
	SubtypePerpetualMarket  = "perpetual_market"
)

var OnChainEventSubtypes = []string{
	SubtypeOrderFill,
	SubtypeSubaccountUpdate,
	SubtypeTransfer,
	SubtypeMarket,
	SubtypeFundingValues,
	SubtypeStatefulOrder,
	SubtypeAsset,
	SubtypePerpetualMarket,
}
