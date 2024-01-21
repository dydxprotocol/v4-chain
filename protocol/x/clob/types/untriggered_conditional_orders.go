package types

import "bytes"

type MinConditionalOrderHeap []Order

func (h MinConditionalOrderHeap) Len() int {
	return len(h)
}

func (h MinConditionalOrderHeap) Less(i, j int) bool {
	x, y := h[i], h[j]
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

type MaxConditionalOrderHeap []Order

func (h MaxConditionalOrderHeap) Len() int {
	return len(h)
}

func (h MaxConditionalOrderHeap) Less(i, j int) bool {
	x, y := h[i], h[j]
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
