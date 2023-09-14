package types

import (
	"bytes"
	"sort"
)

// OrderHash is used to represent the SHA256 hash of an order.
type OrderHash [32]byte

func (oh OrderHash) ToBytes() []byte {
	return oh[:]
}

type SortedOrderHashes []OrderHash

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedOrderHashes{}

func (s SortedOrderHashes) Len() int {
	return len(s)
}

func (s SortedOrderHashes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedOrderHashes) Less(i, j int) bool {
	return bytes.Compare(s[i].ToBytes(), s[j].ToBytes()) < 0
}
