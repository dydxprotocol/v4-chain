package memclob

import (
	"fmt"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// memclobCancels is a utility struct used for storing order cancelations that each expire at a given block height.
type memclobCancels struct {
	// A map of all known canceled order IDs mapped to their expiry block.
	orderIdToExpiry map[types.OrderId]uint32
	// A map from a block height to a set of all canceled order IDs that expire at the block.
	expiryToOrderIds map[uint32]map[types.OrderId]bool
}

// newMemclobCancels returns a new `memclobCancels`.
func newMemclobCancels() *memclobCancels {
	return &memclobCancels{
		orderIdToExpiry:  make(map[types.OrderId]uint32),
		expiryToOrderIds: make(map[uint32]map[types.OrderId]bool),
	}
}

// get returns the `tilBlock` expiry of an order cancelation and a bool indicating whether the expiry exists.
func (c *memclobCancels) get(
	orderId types.OrderId,
) (
	tilBlock uint32,
	exists bool,
) {
	tilBlock, exists = c.orderIdToExpiry[orderId]
	return tilBlock, exists
}

// addShortTermCancel adds an `orderId` to all data structures, expiring at block `tilBlock`.
// Panics if the `orderId` already exists in the data structures.
func (c *memclobCancels) addShortTermCancel(
	orderId types.OrderId,
	tilBlock uint32,
) {
	orderId.MustBeShortTermOrder()
	// Add the `orderId` to the `orderIdToExpiry` map, panicing if it already exists.
	if _, exists := c.orderIdToExpiry[orderId]; exists {
		panic(fmt.Sprintf(
			"mustAddCancel: orderId %+v already exists in orderIdToExpiry",
			orderId,
		))
	}
	c.orderIdToExpiry[orderId] = tilBlock

	// Fetch a reference to the `expiryToOrderIds` for this `tilBlock`, creating it if it does not already exist.
	orderIdsInBlock, exists := c.expiryToOrderIds[tilBlock]
	if !exists {
		orderIdsInBlock = make(map[types.OrderId]bool)
		c.expiryToOrderIds[tilBlock] = orderIdsInBlock
	} else if _, exists = orderIdsInBlock[orderId]; exists {
		panic(fmt.Sprintf(
			"memclobCancels#add: orderId %+v already exists in expiryToOrderIds[%d]",
			orderId,
			tilBlock,
		))
	}

	// Set the `OrderId` in the `expiryToOrderIds` data structure for the new cancel's `tilBlock`.
	orderIdsInBlock[orderId] = true
}

// remove removes an `orderId` from all data structs. Panics if the `orderId` is not found.
func (c *memclobCancels) remove(
	orderId types.OrderId,
) {
	// Panic if the `orderId` does not exist in the `orderIdToExpiry` map.
	goodTilBlock, exists := c.orderIdToExpiry[orderId]
	if !exists {
		panic(fmt.Sprintf(
			"memclobCancels#remove: orderId %+v does not exist in orderIdToExpiry",
			orderId,
		))
	}

	// Panic if the `orderId` does not exist in the appropriate submap of `expiryToOrderIds`.
	expiryToOrderIdsForBlock, exists := c.expiryToOrderIds[goodTilBlock]
	if !exists {
		panic(fmt.Sprintf(
			"memclobCancels#remove: %d does not exist in expiryToOrderIds",
			goodTilBlock,
		))
	}
	if _, exists = expiryToOrderIdsForBlock[orderId]; !exists {
		panic(fmt.Sprintf(
			"memclobCancels#remove: orderId %+v does not exist in expiryToOrderIds[%d]",
			orderId,
			goodTilBlock,
		))
	}

	// Delete the `orderId` from the `orderIdToExpiry` map.
	delete(c.orderIdToExpiry, orderId)

	// Delete the `orderId` from the `expiryToOrderIds` submap.
	// If this is the last order in the submap, delete the submap.
	if len(expiryToOrderIdsForBlock) == 1 {
		delete(c.expiryToOrderIds, goodTilBlock)
	} else {
		delete(expiryToOrderIdsForBlock, orderId)
	}
}

// removeAllAtBlock iterates through (and removes from all data structures) the IDs that expire at a certain `block`.
func (c *memclobCancels) removeAllAtBlock(
	block uint32,
) {
	orderIds, exists := c.expiryToOrderIds[block]

	// If map entry does not exist, return early.
	if !exists {
		return
	}

	// Remove all ids. This also removes `block` from the `expiryToOrderIds` map.
	for id := range orderIds {
		c.remove(id)
	}
}
