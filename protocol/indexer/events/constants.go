package events

const (
	// Cosmos event attribute values for the subtype attribute for different indexer events.
	// Keep these constants in sync with:
	// https://github.com/dydxprotocol/indexer/blob/master/services/ender/src/lib/types.ts.
	// Ender uses these to maintain a mapping between event type and event proto.
	SubtypeOrderFill          = "order_fill"
	SubtypeSubaccountUpdate   = "subaccount_update"
	SubtypeTransfer           = "transfer"
	SubtypeMarket             = "market"
	SubtypeFundingValues      = "funding_values"
	SubtypeStatefulOrder      = "stateful_order"
	SubtypeAsset              = "asset"
	SubtypePerpetualMarket    = "perpetual_market"
	SubtypeLiquidityTier      = "liquidity_tier"
	SubtypeUpdatePerpetual    = "update_perpetual"
	SubtypeUpdateClobPair     = "update_clob_pair"
	SubtypeDeleveraging       = "deleveraging"
	SubtypeOpenInterestUpdate = "open_interest_update"
	SubtypeYieldParams        = "yield_params"
)

const (
	// Indexer event versions.
	OrderFillEventVersion         uint32 = 1
	SubaccountUpdateEventVersion  uint32 = 1
	TransferEventVersion          uint32 = 1
	MarketEventVersion            uint32 = 1
	FundingValuesEventVersion     uint32 = 1
	StatefulOrderEventVersion     uint32 = 1
	AssetEventVersion             uint32 = 1
	PerpetualMarketEventVersion   uint32 = 2
	LiquidityTierEventVersion     uint32 = 2
	UpdatePerpetualEventVersion   uint32 = 1
	UpdateClobPairEventVersion    uint32 = 1
	DeleveragingEventVersion      uint32 = 1
	OpenInterestUpdateVersion     uint32 = 1
	UpdateYieldParamsEventVersion uint32 = 1
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
	SubtypeLiquidityTier,
	SubtypeUpdatePerpetual,
	SubtypeUpdateClobPair,
	SubtypeDeleveraging,
	SubtypeYieldParams,
}
