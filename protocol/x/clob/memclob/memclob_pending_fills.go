package memclob

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/x/clob/types"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// memclobPendingFills is a utility struct used for storing pending fills.
type memclobPendingFills struct {
	// Ordered queue of pending fills, which will be proposed in a block if
	// the validator is the next block proposer.
	// TODO(DEC-1924): Deprecate this after CLOB refactor is completed.
	fills []types.PendingFill
	// A map from the order hashes to the matched orders themselves for each order that
	// was matched.
	// TODO(DEC-1924): Deprecate this after CLOB refactor is completed.
	orderHashToMatchableOrder map[types.OrderHash]types.MatchableOrder
	// A map from order ID to the orders themselves for each order that
	// was matched. Note: there may be multiple distinct orders with the same
	// ID that are matched. In that case, only the "greatest" of any such orders
	// is maintained in this map.
	// TODO(DEC-1924): Deprecate this after CLOB refactor is completed.
	orderIdToOrder map[types.OrderId]types.Order
	// A struct encapsulating data required for determining the operations to propose in the next
	// block if this validator is the block proposer.
	operationsToPropose types.OperationsToPropose
}

// newMemclobPendingFills returns a new `memclobPendingFills`.
func newMemclobPendingFills() *memclobPendingFills {
	return &memclobPendingFills{
		fills:                     make([]types.PendingFill, 0),
		orderHashToMatchableOrder: make(map[types.OrderHash]types.MatchableOrder),
		orderIdToOrder:            make(map[types.OrderId]types.Order),
		operationsToPropose:       *types.NewOperationsToPropose(),
	}
}

// mustGetMatchedOrderByHash returns a matched order by hash. Panics if the order cannot be found.
func (m *memclobPendingFills) mustGetMatchedOrderByHash(
	ctx sdk.Context,
	hash types.OrderHash,
) types.MatchableOrder {
	order, exists := m.orderHashToMatchableOrder[hash]
	if !exists {
		panic(fmt.Sprintf("mustGetMatchedOrderByHash: Order does not exist in matched order map %v", hash))
	}
	return order
}

// resizeReduceOnlyMatchIfNecessary resizes a reduce-only match if it would change or increase
// the position side of the subaccount, and returns the resized match.
func (m *memclobPendingFills) resizeReduceOnlyMatchIfNecessary(
	ctx sdk.Context,
	subaccountId satypes.SubaccountId,
	clobPairId types.ClobPairId,
	currentPositionSize *big.Int,
	newlyMatchedAmount satypes.BaseQuantums,
	isBuy bool,
) satypes.BaseQuantums {
	// Get the signed size of the new match.
	newMatchSize := newlyMatchedAmount.ToBigInt()
	if !isBuy {
		newMatchSize.Neg(newMatchSize)
	}

	// If the match is not on the opposite side of the position, then the match is invalid.
	// Note that this can occur for reduce-only maker orders if the maker subaccount's position side
	// changes during the matching loop, and this should never happen for taker orders.
	if currentPositionSize.Sign()*newMatchSize.Sign() != -1 {
		return satypes.BaseQuantums(0)
	}

	// The match is on the opposite side of the position. Return the minimum of the match size and
	// position size to ensure the new match does not change the subaccount's position side.
	absPositionSize := new(big.Int).Abs(currentPositionSize)
	absNewMatchSize := new(big.Int).Abs(newMatchSize)
	maxMatchSize := lib.BigMin(absPositionSize, absNewMatchSize)
	return satypes.BaseQuantums(maxMatchSize.Uint64())
}

// mustAddOrderToOperationsToPropose adds an order to the operations to propose.
func (m *memclobPendingFills) mustAddOrderToOperationsToPropose(
	ctx sdk.Context,
	order types.Order,
	isPreexistingStatefulOrder bool,
) {
	if isPreexistingStatefulOrder {
		m.operationsToPropose.AddPreexistingStatefulOrderPlacementToOperationsQueue(
			order,
		)
	} else {
		m.operationsToPropose.AddOrderPlacementToOperationsQueue(
			order,
		)
	}
}

// isMakerOrderInOperationsToPropose returns true if the order is in the operations to propose,
// false if not. Note this function panics if called with orders that aren't assigned nonces,
// which is why it shouldn't be used for taker orders.
func (m *memclobPendingFills) isMakerOrderInOperationsToPropose(
	ctx sdk.Context,
	order types.Order,
	isPreexistingStatefulOrder bool,
) bool {
	if isPreexistingStatefulOrder {
		return m.operationsToPropose.IsPreexistingStatefulOrderInOperationsQueue(
			order,
		)
	} else {
		return m.operationsToPropose.IsOrderPlacementInOperationsQueue(
			order,
		)
	}
}
