package lib

import (
	"golang.org/x/exp/constraints"
	"sort"
)

// Sortable[K] attaches the methods of sort.Interface to []K, sorting in increasing order.
type Sortable[K constraints.Ordered] []K

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
// int is used as an example.
var _ sort.Interface = Sortable[int]{}

func (s Sortable[K]) Len() int {
	return len(s)
}

func (s Sortable[K]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sortable[K]) Less(i, j int) bool {
	return s[i] < s[j]
}
