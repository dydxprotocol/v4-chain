package types

import "bytes"

// MinConditionalOrderHeap is type alias for `[]Order` which implements "container/heap"
// interface.
//
// This is a _MIN_ heap. Orders are compared by their `ConditionalOrderTriggerSubticks` field, and then by
// their `OrderHash` if the trigger subticks are equal.
type MinConditionalOrderHeap []Order

func (h MinConditionalOrderHeap) Len() int {
	return len(h)
}

func (h MinConditionalOrderHeap) Less(i, j int) bool {
	x, y := h[i], h[j]

	// If the trigger subticks are the same, sort by order hash.
	// This is required for determinism in the case of multiple orders with the same trigger subticks.
	if x.ConditionalOrderTriggerSubticks == y.ConditionalOrderTriggerSubticks {
		xHash := x.GetOrderHash()
		yHash := y.GetOrderHash()
		return bytes.Compare(xHash[:], yHash[:]) == -1
	}
	return x.ConditionalOrderTriggerSubticks < y.ConditionalOrderTriggerSubticks
}

func (h MinConditionalOrderHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinConditionalOrderHeap) Push(x interface{}) {
	*h = append(*h, x.(Order))
}

func (h *MinConditionalOrderHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// MaxConditionalOrderHeap is type alias for `[]Order` which implements "container/heap"
// interface.
//
// This is a _MAX_ heap. Orders are compared by their `ConditionalOrderTriggerSubticks` field, and then by
// their `OrderHash` if the trigger subticks are equal.
type MaxConditionalOrderHeap []Order

func (h MaxConditionalOrderHeap) Len() int {
	return len(h)
}

func (h MaxConditionalOrderHeap) Less(i, j int) bool {
	x, y := h[i], h[j]

	// If the trigger subticks are the same, sort by order hash.
	// This is required for determinism in the case of multiple orders with the same trigger subticks.
	if x.ConditionalOrderTriggerSubticks == y.ConditionalOrderTriggerSubticks {
		xHash := x.GetOrderHash()
		yHash := y.GetOrderHash()
		return bytes.Compare(xHash[:], yHash[:]) == -1
	}
	return x.ConditionalOrderTriggerSubticks > y.ConditionalOrderTriggerSubticks
}

func (h MaxConditionalOrderHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxConditionalOrderHeap) Push(x interface{}) {
	*h = append(*h, x.(Order))
}

func (h *MaxConditionalOrderHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
