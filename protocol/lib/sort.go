package lib

import (
	"sort"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// SortSubticks sorts a slice of Subticks in ascending order.
func SortSubticks(subticks []types.Subticks) {
	sort.Slice(subticks, func(i, j int) bool {
		return subticks[i] < subticks[j]
	})
}

// SortSubticksDesc sorts a slice of Subticks in descending order.
func SortSubticksDesc(subticks []types.Subticks) {
	sort.Slice(subticks, func(i, j int) bool {
		return subticks[i] > subticks[j]
	})
}
