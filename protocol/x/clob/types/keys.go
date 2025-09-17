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
