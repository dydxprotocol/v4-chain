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
	OrdersFilledDuringLatestBlockKey = "OrdersFilledDuringLatestBlockKey/value"
	// StatefulOrderPlacementKeyPrefix is the key to retrieve a stateful order and information about
	// when it was placed.
	StatefulOrderPlacementKeyPrefix = "StatefulOrderPlacement/value/"
	// StatefulOrdersTimeSlicePrefix is the key to retrieve a unique list of the stateful orders that
	// expire at a given timestamp, sorted by order ID.
	StatefulOrdersTimeSlicePrefix = "StatefulOrdersTimeSlice/value/"
	// LastCommittedBlockTimeKey defines the key that stores the block time of the previously committed block.
	LastCommittedBlockTimeKey = "LastCommittedBlockTime/value"
	// NextStatefulOrderBlockTransactionIndexKey is the transient store key that stores the next
	// transaction index to use for the next newly-placed stateful order.
	NextStatefulOrderBlockTransactionIndexKey = "NextStatefulOrderBlockTransactionIndex/value"
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
