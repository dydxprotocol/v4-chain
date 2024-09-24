package heap

import (
	"container/heap"
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
)

// LiquidationPriority represents an item in our priority queue
type LiquidationPriority struct {
	SubaccountId types.SubaccountId
	Priority     *big.Float
	Index        int // The index of the item in the heap
}

// LiquidationPriorityHeap is a min-heap of LiquidationPriority items
type LiquidationPriorityHeap []*LiquidationPriority

// Len returns the number of elements in the heap
func (h LiquidationPriorityHeap) Len() int { return len(h) }

// Less defines the ordering of items in the heap
func (h LiquidationPriorityHeap) Less(i, j int) bool {
	return h[i].Priority.Cmp(h[j].Priority) < 0
}

// Swap swaps the elements with indexes i and j
func (h LiquidationPriorityHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].Index = i
	h[j].Index = j
}

// Push adds an element to the heap
func (h *LiquidationPriorityHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*LiquidationPriority)
	item.Index = n
	*h = append(*h, item)
}

// Pop removes and returns the minimum element (according to Less) from the heap
func (h *LiquidationPriorityHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	item.Index = -1 // for safety
	old[n-1] = nil  // avoid memory leak
	return item
}

// NewLiquidationPriorityHeap creates and initializes a new LiquidationPriorityHeap
func NewLiquidationPriorityHeap() *LiquidationPriorityHeap {
	h := &LiquidationPriorityHeap{}
	heap.Init(h)
	return h
}

// AddSubaccount adds a new subaccount to the heap
func (h *LiquidationPriorityHeap) AddSubaccount(subaccountId types.SubaccountId, priority *big.Float) {
	heap.Push(h, &LiquidationPriority{
		SubaccountId: subaccountId,
		Priority:     priority,
	})
}

// PopLowestPriority removes and returns the subaccount with the lowest priority
func (h *LiquidationPriorityHeap) PopLowestPriority() *LiquidationPriority {
	if h.Len() == 0 {
		return nil
	}
	return heap.Pop(h).(*LiquidationPriority)
}

// UpdatePriority updates the priority of a subaccount in the heap
func (h *LiquidationPriorityHeap) UpdatePriority(item *LiquidationPriority, newPriority *big.Float) bool {
	if item.Index < 0 || item.Index >= len(*h) {
		return false
	}
	if (*h)[item.Index] != item {
		return false
	}
	item.Priority = newPriority
	heap.Fix(h, item.Index)
	return true
}
