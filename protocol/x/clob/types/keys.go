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
	TransientStoreKey = "transient_" + ModuleName
)

const (
	// ProcessProposerMatchesEventsKey is the key to retrieve information about how to update
	// memclob state based on the latest block.
	ProcessProposerMatchesEventsKey = "process_proposer_matches_events"

	// SubaccountLiquidationInfoKeyPrefix is the prefix to retrieve the liquidation information
	// for a subaccount within the last block.
	SubaccountLiquidationInfoKeyPrefix = "subaccount_liquidation_info/"

	// LiquidationsConfigKey is the key to retrieve the liquidations config.
	LiquidationsConfigKey = "liquidations_config"

	// EquityTierLimitConfigKey is the key to retrieve the equity tier limit configuration.
	EquityTierLimitConfigKey = "equity_tier_limit_config"

	// BlockRateLimitConfigKey is the key to retrieve the block rate limit configuration.
	BlockRateLimitConfigKey = "block_rate_limit_config"

	// ClobPairKeyPrefix is the prefix to retrieve all ClobPair
	ClobPairKeyPrefix = "clob_pair/"
)

const (
	// OrderAmountFilledKeyPrefix is the prefix to retrieve the fill amount for an order.
	OrderAmountFilledKeyPrefix = "order_amount_filled/"

	// BlockHeightToPotentiallyPrunableOrdersPrefix is the prefix to retrieve a list of potentially prunable orders
	// by block height.
	BlockHeightToPotentiallyPrunableOrdersPrefix = "block_height_to_potentially_prunable_orders/"

	// OrdersFilledDuringLatestBlockKey is the key to retrieve the list of orders filled during the latest block.
	OrdersFilledDuringLatestBlockKey = "orders_filled_during_latest_block"
)

// Below key prefixes are not explicitly used to read/write to state, but rather used to iterate over
// certain groups of items stored in state.
const (
	// StatefulOrderKeyPrefix is the prefix key for all long term orders and all conditional orders,
	// both triggered and untriggered.
	StatefulOrderKeyPrefix = "stateful_order_placement/"

	// PlacedStatefulOrderKeyPrefix is the prefix key for placed long term orders and triggered
	// conditional orders. It represents all stateful orders that should be placed upon the memclob
	// during app start up.
	PlacedStatefulOrderKeyPrefix = StatefulOrderKeyPrefix + "placed/"
)

// Store / Memstore
const (
	// TriggeredConditionalOrderKeyPrefix is the key to retrieve an triggered conditional order and
	// information about when it was triggered.
	TriggeredConditionalOrderKeyPrefix = PlacedStatefulOrderKeyPrefix + "conditional/"

	// LongTermOrderPlacementKeyPrefix is the key to retrieve a long term order and information about
	// when it was placed.
	LongTermOrderPlacementKeyPrefix = PlacedStatefulOrderKeyPrefix + "long_term/"

	// UntriggeredConditionalOrderKeyPrefix is the key to retrieve an untriggered conditional order and
	// information about when it was placed.
	UntriggeredConditionalOrderKeyPrefix = StatefulOrderKeyPrefix + "untriggered/conditional/"

	// StatefulOrdersTimeSlicePrefix is the key to retrieve a unique list of the stateful orders that
	// expire at a given timestamp, sorted by order ID.
	StatefulOrdersTimeSlicePrefix = "stateful_orders_time_slice/"
)

// Transient Store
const (
	// NextStatefulOrderBlockTransactionIndexKey is the transient store key that stores the next
	// transaction index to use for the next newly-placed stateful order.
	NextStatefulOrderBlockTransactionIndexKey = "next_stateful_order_block_transaction_index"

	// UncommittedStatefulOrderPlacementKeyPrefix is the key to retrieve an uncommitted stateful order and information
	// about when it was placed. uncommitted orders are orders that this validator is aware of that have yet to be
	// committed to a block and are stored in a transient store.
	UncommittedStatefulOrderPlacementKeyPrefix = StatefulOrderKeyPrefix + "uncommitted/long_term/"

	// UncommittedStatefulOrderCancellationKeyPrefix is the key to retrieve an uncommitted stateful order cancellation.
	// uncommitted cancelleations are cancellations that this validator is aware of that have yet to be
	// committed to a block and are stored in a transient store.
	UncommittedStatefulOrderCancellationKeyPrefix = "stateful_order_cancellation/uncommitted/long_term/"

	// UncommittedStatefulOrderCountPrefix is the key to retrieve an uncommitted stateful order count.
	// uncommitted orders are orders that this validator is aware of that have yet to be committed to a block and
	// are stored in a transient store. This count represents the number of uncommitted stateful
	// `placements - cancellations`.
	UncommittedStatefulOrderCountPrefix = "stateful_order_count/uncommitted/long_term/"

	// ToBeCommittedStatefulOrderCountPrefix is the key to retrieve the to be committed stateful order count.
	// To be committed orders are orders that this validator is aware of during `DeliverTx` that are in the process
	// or being committed to a block and are stored in a transient store. This count represents the number of to
	// be committed stateful `placements - removals`.
	ToBeCommittedStatefulOrderCountPrefix = "stateful_order_count/to_be_committed/long_term/"
)

// Module Accounts
const (
	// InsuranceFundName defines the root string for the insurance fund account address
	InsuranceFundName = "insurance_fund"
)
