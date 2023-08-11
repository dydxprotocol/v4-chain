package types

import (
	"fmt"
	"sort"
)

// SortedStatefulOrderPlacements is type alias for `[]OrderId` which supports deterministic
// sorting. Stateful order placements are first ordered by block height, followed by
// transaction index.
// This list assumes that all order placements are unique and will panic if
// any two order placements have the same block height and transaction index.
type SortedStatefulOrderPlacements []StatefulOrderPlacement

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedStatefulOrderPlacements{}

func (s SortedStatefulOrderPlacements) Len() int {
	return len(s)
}

func (s SortedStatefulOrderPlacements) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedStatefulOrderPlacements) Less(i, j int) bool {
	si := s[i]
	sj := s[j]

	if si.BlockHeight != sj.BlockHeight {
		return si.BlockHeight < sj.BlockHeight
	}

	if si.TransactionIndex != sj.TransactionIndex {
		return si.TransactionIndex < sj.TransactionIndex
	}

	panic(
		fmt.Errorf(
			"Less: stateful order placements (%v) and (%v) have the same block height and transaction index",
			si,
			sj,
		),
	)
}
