package types

// Module name and store keys
const (
	// ModuleName defines the module name
	ModuleName = "clob"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_" + ModuleName

	// TransientStoreKey defines the primary module transient store key
	TransientStoreKey = "tmp_" + ModuleName
)

// Below key prefixes are not explicitly used to read/write to state, but rather used to iterate over
// certain groups of items stored in state.
const (
	// StatefulOrderKeyPrefix is the prefix key for all long term orders and all conditional orders,
	// both triggered and untriggered.
	StatefulOrderKeyPrefix = "SO/"

	// PlacedStatefulOrderKeyPrefix is the prefix key for placed long term orders and triggered
	// conditional orders. It represents all stateful orders that should be placed upon the memclob
	// during app start up.
	PlacedStatefulOrderKeyPrefix = StatefulOrderKeyPrefix + "P/"

	// PrunableOrdersKeyPrefix is the prefix key for orders prunable at a certain height.
	PrunableOrdersKeyPrefix = "PO/"
)

// State
const (
	// LiquidationsConfigKey is the key to retrieve the liquidations config.
	LiquidationsConfigKey = "LiqCfg"

	// EquityTierLimitConfigKey is the key to retrieve the equity tier limit configuration.
	EquityTierLimitConfigKey = "EqTierCfg"

	// BlockRateLimitConfigKey is the key to retrieve the block rate limit configuration.
	BlockRateLimitConfigKey = "RateLimCfg"

	// ClobPairKeyPrefix is the prefix to retrieve all ClobPair
	ClobPairKeyPrefix = "Clob:"

	// OrderAmountFilledKeyPrefix is the prefix to retrieve the fill amount for an order.
	OrderAmountFilledKeyPrefix = "Fill:"

	// Deprecated: LegacyBlockHeightToPotentiallyPrunableOrdersPrefix is the prefix to retrieve a list of
	// potentially prunable short term orders by block height. Should not be used after migrating to
	// key-per-order format.
	LegacyBlockHeightToPotentiallyPrunableOrdersPrefix = "ExpHt:"

	// Deprecated: LegacyStatefulOrdersTimeSlicePrefix is the key to retrieve a unique list of the stateful
	// orders that expire at a given timestamp, sorted by order ID. Do not use.
	LegacyStatefulOrdersTimeSlicePrefix = "ExpTm:"

	// StatefulOrdersTimeSliceKeyPrefix is used to store orders that expire at a certain time.
	// The specifier should be replaced with the time.
	StatefulOrdersExpirationsKeyPrefix = "Exp/%s:"

	// TriggeredConditionalOrderKeyPrefix is the key to retrieve an triggered conditional order and
	// information about when it was triggered.
	TriggeredConditionalOrderKeyPrefix = PlacedStatefulOrderKeyPrefix + "T:"

	// TWAPOrderKeyPrefix is the key to retrieve a TWAP order and information about when it was placed.
	TWAPOrderKeyPrefix = "TWAP:"

	// TWAPTriggerOrderKeyPrefix is the key to retrieve TWAP suborder information.
	TWAPTriggerOrderKeyPrefix = "TWAP/T:"

	// LongTermOrderPlacementKeyPrefix is the key to retrieve a long term order and information about
	// when it was placed.
	LongTermOrderPlacementKeyPrefix = PlacedStatefulOrderKeyPrefix + "L:"

	// UntriggeredConditionalOrderKeyPrefix is the key to retrieve an untriggered conditional order and
	// information about when it was placed.
	UntriggeredConditionalOrderKeyPrefix = StatefulOrderKeyPrefix + "U:"

	// ConditionalOrderTriggerPriceIndexKeyPrefix is the prefix for the secondary trigger-price index.
	// IMPORTANT: this prefix must NOT start with StatefulOrderKeyPrefix ("SO/") to avoid being
	// picked up by getAllOrdersIterator / GetAllStatefulOrders, which would try to unmarshal
	// the empty-value index entries as LongTermOrderPlacement.
	// Keys are structured as:
	//   <prefix> <clobPairId:4 big-endian> <directionByte:1> <triggerSubticks:8 big-endian>
	//     <placementSequence:8 big-endian> <orderId state key>
	// directionByte = 0x00 for LTE-direction (trigger when oracle ≤ triggerPrice)
	//               = 0x01 for GTE-direction (trigger when oracle ≥ triggerPrice)
	// triggerSubticks is big-endian so byte ordering is monotonically increasing in price.
	ConditionalOrderTriggerPriceIndexKeyPrefix = "TPIdx:"

	// NextClobPairIDKey is the key to retrieve the next ClobPair ID to be used.
	NextClobPairIDKey = "NextClobPairID"

	// LeverageKeyPrefix is the prefix for leverage storage
	LeverageKeyPrefix = "Leverage:"
)

// Memstore
const (
	// KeyMemstoreInitialized is the key to check if the memstore has been initialized.
	KeyMemstoreInitialized = "MemstoreInit"

	// ProcessProposerMatchesEventsKey is the key to retrieve information about how to update
	// memclob state based on the latest block.
	ProcessProposerMatchesEventsKey = "ProposerEvents"

	// The following Delivered keys used to be a part of ProcessProposerMatchesEvents but were taken out to unnecessary
	// serde of a big monolithic value.

	// OrderedDeliveredLongTermOrderIndexKey stores the next index to be used for OrderedDeliveredLongTermOrderKeyPrefix
	OrderedDeliveredLongTermOrderIndexKey = "DLTOIdx"
	// OrderedDeliveredLongTermOrderKeyPrefix is used to store placed orders for memclob placement in PrepareCheckState.
	OrderedDeliveredLongTermOrderKeyPrefix = "DLTO:"

	// OrderedDeliveredConditionalOrdexIndexKey stores the next index to be used for
	// OrderedDeliveredConditionalOrderKeyPrefix
	OrderedDeliveredConditionalOrderIndexKey = "DCOIdx"
	// OrderedDeliveredConditionalOrderKeyPrefix is used to store placed orders for memclob placement in PrepareCheckState.
	OrderedDeliveredConditionalOrderKeyPrefix = "DCIdx:"

	// DeliveredCancelKeyPrefix is used to store placed orders for memclob placement in PrepareCheckState.
	DeliveredCancelKeyPrefix = "DCancel:"

	// StatefulOrderCountPrefix is the key to retrieve the stateful order count. The stateful order count
	// represents the number of stateful orders stored in state.
	StatefulOrderCountPrefix = "NumSO:"

	// UntriggeredConditionalOrderCountGlobalKey is the single memstore key that holds the global count
	// of resting untriggered conditional orders across all subaccounts and clob pairs.
	// Maintained at the same four lifecycle hooks as the trigger-price index (Packet 1).
	UntriggeredConditionalOrderCountGlobalKey = "NumUcond"

	// UntriggeredConditionalOrderCountPerSubaccountPrefix is the memstore key prefix for the
	// per-subaccount count of resting untriggered conditional orders.
	// Full key = prefix + subaccountId.ToStateKey().
	UntriggeredConditionalOrderCountPerSubaccountPrefix = "NumUcondSa:"

	// ConditionalOrderTriggerConfigKey is the single persistent-store key holding the
	// consensus-level configuration that gates the bounded conditional-order trigger path.
	// When absent (the default), MaybeTriggerConditionalOrders runs the legacy full-scan
	// behavior that is byte-for-byte identical to the pre-fix implementation, so the fix can be
	// shipped on a rolling basis and activated later at a governed height without a version split.
	// When present and enabled, incremental activation starts; the bounded crossing-priority path
	// runs only after ConditionalOrderTriggerIndexReadyKey is set. These values are consensus state,
	// so every node agrees deterministically on which path executes.
	ConditionalOrderTriggerConfigKey = "CondTrigCfg"

	// ConditionalOrderTriggerNextClobPairKey stores the clob-pair id at which the next bounded
	// conditional-order scheduling pass should begin. Rotating this cursor prevents fixed ascending
	// pair order from starving later markets when the trigger budget is smaller than the number of
	// active markets.
	ConditionalOrderTriggerNextClobPairKey = "CondTrigNextPair"

	// ConditionalOrderTriggerIndexActivation* keys persist the incremental trigger-price-index
	// activation state. The legacy trigger path remains authoritative until ReadyKey is set.
	ConditionalOrderTriggerIndexActivationPhaseKey   = "CondTrigIdxPhase"
	ConditionalOrderTriggerIndexActivationCursorKey  = "CondTrigIdxCursor"
	ConditionalOrderTriggerIndexActivationClearedKey = "CondTrigIdxCleared"
	ConditionalOrderTriggerIndexActivationIndexedKey = "CondTrigIdxIndexed"
	ConditionalOrderTriggerIndexReadyKey             = "CondTrigIdxReady"
)

// Transient Store
const (
	// SubaccountLiquidationInfoKeyPrefix is the prefix to retrieve the liquidation information
	// for a subaccount within the last block.
	SubaccountLiquidationInfoKeyPrefix = "SaLiqInfo:"

	// NextStatefulOrderBlockTransactionIndexKey is the transient store key that stores the next
	// transaction index to use for the next newly-placed stateful order.
	NextStatefulOrderBlockTransactionIndexKey = "NextTxIdx"

	// UncommittedStatefulOrderPlacementKeyPrefix is the key to retrieve an uncommitted stateful order and information
	// about when it was placed. Uncommitted orders are orders that this validator is aware of that have yet to be
	// committed to a block and are stored in a transient store.
	UncommittedStatefulOrderPlacementKeyPrefix = "UncmtSO:"

	// UncommittedStatefulOrderCancellationKeyPrefix is the key to retrieve an uncommitted stateful order cancellation.
	// Uncommitted cancelleations are cancellations that this validator is aware of that have yet to be
	// committed to a block and are stored in a transient store.
	UncommittedStatefulOrderCancellationKeyPrefix = "UncmtSOCxl:"

	// UncommittedStatefulOrderCountPrefix is the key to retrieve an uncommitted stateful order count.
	// Uncommitted orders are orders that this validator is aware of that have yet to be committed to a block and
	// are stored in a transient store. This count represents the number of uncommitted stateful
	// `placements - cancellations`.
	UncommittedStatefulOrderCountPrefix = "NumUncmtSO:"

	// MinTradePricePrefix is the key prefix to retrieve the min trade price for a perpetual.
	// This is meant to be used for improved conditional order triggering.
	MinTradePricePrefix = "MinTrade:"

	// MaxTradePricePrefix is the key prefix to retrieve the max trade price for a perpetual.
	// This is meant to be used for improved conditional order triggering.
	MaxTradePricePrefix = "MaxTrade:"
)

// FinalizeBlock event staging
const (
	StagedEventsCountKey  = "StgEvtCnt"
	StagedEventsKeyPrefix = "StgEvt:"
)
