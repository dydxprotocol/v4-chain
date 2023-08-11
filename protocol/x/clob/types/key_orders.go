package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
)

const (
	// OrderAmountFilledKeyPrefix is the prefix to retrieve the fill amount for an order.
	OrderAmountFilledKeyPrefix = "OrderAmount/value/"
	// BlockHeightToPotentiallyPrunableOrdersPrefix is the prefix to retrieve a list of potentially prunable orders
	// by block height.
	BlockHeightToPotentiallyPrunableOrdersPrefix = "BlockHeightToPotentiallyPrunableOrders/value/"
	// OrdersFilledDuringLatestBlockKey is the key to retrieve the list of orders filled during the latest block.
	OrdersFilledDuringLatestBlockKey = "OrdersFilledDuringLatestBlockKey/value/"
	// LongTermOrderPlacementKeyPrefix is the key to retrieve a long term order and information about
	// when it was placed.
	LongTermOrderPlacementKeyPrefix = "StatefulOrderPlacement/Placed/LongTerm/value/"
	// UncommittedStatefulOrderPlacementKeyPrefix is the key to retrieve an uncommitted stateful order and information
	// about when it was placed. Uncommitted orders are orders that this validator is aware of that have yet to be
	// committed to a block and are stored in a transient store.
	UncommittedStatefulOrderPlacementKeyPrefix = "StatefulOrderPlacement/Uncommitted/LongTerm/value/"
	// UncommittedStatefulOrderCancellationKeyPrefix is the key to retrieve an uncommitted stateful order cancellation.
	// Uncommitted cancelleations are cancellations that this validator is aware of that have yet to be
	// committed to a block and are stored in a transient store.
	UncommittedStatefulOrderCancellationKeyPrefix = "StatefulOrderCancellation/Uncommitted/LongTerm/value/"
	// UntriggeredConditionalOrderKeyPrefix is the key to retrieve an untriggered conditional order and
	// information about when it was placed.
	UntriggeredConditionalOrderKeyPrefix = "StatefulOrderPlacement/Untriggered/Conditional/value/"
	// TriggeredConditionalOrderKeyPrefix is the key to retrieve an triggered conditional order and
	// information about when it was triggered.
	TriggeredConditionalOrderKeyPrefix = "StatefulOrderPlacement/Placed/Conditional/value/"
	// StatefulOrdersTimeSlicePrefix is the key to retrieve a unique list of the stateful orders that
	// expire at a given timestamp, sorted by order ID.
	StatefulOrdersTimeSlicePrefix = "StatefulOrdersTimeSlice/value/"
	// LastCommittedBlockTimeKey defines the key that stores the block time of the previously committed block.
	LastCommittedBlockTimeKey = "LastCommittedBlockTime/value/"
	// NextStatefulOrderBlockTransactionIndexKey is the transient store key that stores the next
	// transaction index to use for the next newly-placed stateful order.
	NextStatefulOrderBlockTransactionIndexKey = "NextStatefulOrderBlockTransactionIndex/value/"
)

// Below key prefixes are not explicitly used to read/write to state, but rather used to iterate over
// certain groups of items stored in state.
const (
	// PlacedStatefulOrderKeyPrefix is the prefix key for placed long term orders and triggered
	// conditional orders. It represents all stateful orders that should be placed upon the memclob
	// during app start up.
	PlacedStatefulOrderKeyPrefix = "StatefulOrderPlacement/Placed/"
	// StatefulOrderKeyPrefix is the prefix key for all long term orders and all conditional orders,
	// both triggered and untriggered.
	StatefulOrderKeyPrefix = "StatefulOrderPlacement/"
)

// OrderIdKey returns the order ID marshaled to bytes. The bytes are meant to be
// used as a subkey for getting or setting an order in state.
func OrderIdKey(
	id OrderId,
) []byte {
	idBytes, err := id.Marshal()
	if err != nil {
		panic(err)
	}

	return idBytes
}

// BlockHeightToPotentiallyPrunableOrdersKey returns the store key to retrieve a list of potentially prunable orders
// by `blockHeight`.
func BlockHeightToPotentiallyPrunableOrdersKey(
	id uint32,
) []byte {
	return lib.Uint32ToBytes(id)
}

// GetTimeSliceKey returns a key prefix for indexing an
// slice of stateful orders placements or cancellations expiring at a timestamp.
func GetTimeSliceKey(timestamp time.Time) []byte {
	return sdk.FormatTimeBytes(timestamp)
}
