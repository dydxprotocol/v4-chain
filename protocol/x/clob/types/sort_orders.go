package types

import (
	"fmt"
	"sort"
)

// SortedLongTermOrderPlacements is type alias for `[]LongTermOrderPlacement` which supports deterministic
// sorting. Long Term order placements are first ordered by block height, followed by
// transaction index.
// This list assumes that all order placements are unique and will panic if
// any two order placements have the same block height and transaction index.
type SortedLongTermOrderPlacements []LongTermOrderPlacement

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedLongTermOrderPlacements{}

func (s SortedLongTermOrderPlacements) Len() int {
	return len(s)
}

func (s SortedLongTermOrderPlacements) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedLongTermOrderPlacements) Less(i, j int) bool {
	si := s[i]
	sj := s[j]

	siPlacementIndex := si.GetPlacementIndex()
	sjPlacementIndex := sj.GetPlacementIndex()

	if siPlacementIndex.BlockHeight != sjPlacementIndex.BlockHeight {
		return siPlacementIndex.BlockHeight < sjPlacementIndex.BlockHeight
	}

	if siPlacementIndex.TransactionIndex != sjPlacementIndex.TransactionIndex {
		return siPlacementIndex.TransactionIndex < sjPlacementIndex.TransactionIndex
	}

	panic(
		fmt.Errorf(
			"Less: long term order placements (%v) and (%v) have the same block height and transaction index",
			si,
			sj,
		),
	)
}

// StatefulOrderPlacement represents any type of order placement that is stored in state.
// An stateful order placement exposes a `TransactionOrdering` transaction index that is used when
// sorting the stateful order placements. Ordering is first by block height, then transaction index.
// The following types are currently supported:
// - LongTermOrderPlacement
// - ConditionalOrderPlacement
type StatefulOrderPlacement interface {
	GetTransactionIndex() TransactionOrdering
}

// GetTransactionIndex returns a LongTermOrderPlacement's placement index.
func (orderPlacement *LongTermOrderPlacement) GetTransactionIndex() TransactionOrdering {
	return orderPlacement.GetPlacementIndex()
}

// GetTransactionIndex returns a ConditionalOrderPlacement's TransactionOrdering.
// If the conditional order is triggered, the trigger index is used. If the conditional order
// is placed but not triggered, the placement index is used.
func (orderPlacement *ConditionalOrderPlacement) GetTransactionIndex() TransactionOrdering {
	if orderPlacement.GetTriggerIndex() != nil {
		return *orderPlacement.GetTriggerIndex()
	}
	return orderPlacement.GetPlacementIndex()
}

// SortedStatefulOrderPlacement is type alias for `[]StatefulOrderPlacement` which supports deterministic
// sorting. Stateful orders must expose one TransactionOrdering to be sorted on. Ordering is first by
// block height, then transaction index.
// The following types are currently supported:
// - LongTermOrderPlacement
// - ConditionalOrderPlacement
// This list assumes that all order placements are unique and will panic if
// any two stateful order placements have the same TransactionOrdering.
type SortedStatefulOrderPlacement []StatefulOrderPlacement

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedStatefulOrderPlacement{}

func (s SortedStatefulOrderPlacement) Len() int {
	return len(s)
}

func (s SortedStatefulOrderPlacement) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedStatefulOrderPlacement) Less(i, j int) bool {
	si := s[i]
	sj := s[j]

	siPlacementIndex := si.GetTransactionIndex()
	sjPlacementIndex := sj.GetTransactionIndex()

	if siPlacementIndex.BlockHeight != sjPlacementIndex.BlockHeight {
		return siPlacementIndex.BlockHeight < sjPlacementIndex.BlockHeight
	}

	if siPlacementIndex.TransactionIndex != sjPlacementIndex.TransactionIndex {
		return siPlacementIndex.TransactionIndex < sjPlacementIndex.TransactionIndex
	}

	panic(
		fmt.Errorf(
			"Less: stateful order placements (%v) and (%v) have the same TransactionOrdering",
			si,
			sj,
		),
	)
}

// SortedClobPairId is type alias for `[]ClobPairId` which supports deterministic
// sorting of clob pair Ids.
// This list assumes that all clob pair ids are unique and will panic if
// any two order placements have the same block height and transaction index.
type SortedClobPairId []ClobPairId

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedClobPairId{}

func (s SortedClobPairId) Len() int {
	return len(s)
}

func (s SortedClobPairId) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedClobPairId) Less(i, j int) bool {
	si := s[i]
	sj := s[j]

	return uint32(si) < uint32(sj)
}
