package types

import (
	"fmt"
	"sort"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

const (
	OrderIdFlags_ShortTerm    = uint32(0)
	OrderIdFlags_Conditional  = uint32(32)
	OrderIdFlags_LongTerm     = uint32(64)
	OrderIdFlags_Twap         = uint32(128)
	OrderIdFlags_TwapSuborder = uint32(256)
)

// IsShortTermOrder returns true if this order ID is for a short-term order, false if
// not (which implies the order ID is for a long-term or conditional order).
// Note that all short-term orders will have the `OrderFlags` field set to 0.
func (o *OrderId) IsShortTermOrder() bool {
	return o.OrderFlags == OrderIdFlags_ShortTerm
}

// IsConditionalOrder returns true if this order ID is for a conditional order, false if
// not (which implies the order ID is for a short-term or long-term order).
func (o *OrderId) IsConditionalOrder() bool {
	// If the third bit in the first byte is set and no other bits are set,
	// this is a conditional order.
	// Note that 32 in decimal == 0x20 in hex == 0b00100000 in binary.
	return o.OrderFlags == OrderIdFlags_Conditional
}

// IsTwapOrder returns true if this order ID is for a TWAP order, false if
// not (which implies the order ID is for a short-term or long-term order).
func (o *OrderId) IsTwapOrder() bool {
	return o.OrderFlags == OrderIdFlags_Twap
}

// IsLongTermOrder returns true if this order ID is for a long-term order, false if
// not (which implies the order ID is for a short-term or conditional order).
func (o *OrderId) IsLongTermOrder() bool {
	// If the second bit in the first byte is set and no other bits are set,
	// this is a long-term order.
	// Note that 64 in decimal == 0x40 in hex == 0b01000000 in binary.
	return o.OrderFlags == OrderIdFlags_LongTerm
}

func (o *OrderId) IsTwapSuborder() bool {
	return o.OrderFlags == OrderIdFlags_TwapSuborder
}

// IsStatefulOrder returns whether this order is a stateful order, which is true for Long-Term
// and conditional orders and false for Short-Term orders.
func (o *OrderId) IsStatefulOrder() bool {
	return o.IsLongTermOrder() || o.IsConditionalOrder() || o.IsTwapOrder() || o.IsTwapSuborder()
}

// MustBeStatefulOrder panics if the orderId is not a stateful order, else it does nothing.
func (o *OrderId) MustBeStatefulOrder() {
	if !o.IsStatefulOrder() {
		panic(
			fmt.Sprintf(
				"MustBeStatefulOrder: called with non-stateful order ID (%+v)",
				*o,
			),
		)
	}
}

// MustBeConditionalOrder panics if the orderId is not a conditional order, else it does nothing.
func (o *OrderId) MustBeConditionalOrder() {
	if !o.IsConditionalOrder() {
		panic(
			fmt.Sprintf(
				"MustBeConditionalOrder: called with non-conditional order ID (%+v)",
				o,
			),
		)
	}
}

// MustBeShortTermOrder panics if the orderId is not a short term order, else it does nothing.
func (o *OrderId) MustBeShortTermOrder() {
	if o.IsStatefulOrder() {
		panic(
			fmt.Sprintf(
				"MustBeShortTermOrder: called with stateful order ID (%+v)",
				*o,
			),
		)
	}
}

// Validate performs checks on the OrderId. It performs the following checks:
// - Validates subaccount id
// - checks OrderFlags for validity
func (o *OrderId) Validate() error {
	subaccountId := o.GetSubaccountId()
	if err := subaccountId.Validate(); err != nil {
		return err
	}
	if !o.IsShortTermOrder() && !o.IsStatefulOrder() {
		return errorsmod.Wrapf(ErrInvalidOrderFlag, "orderId: %v", o)
	}
	return nil
}

// ToStateKey returns a bytes representation of a OrderId for use as a state key.
// The key uses the proto marshaling of the object such that it can be unmarshalled in
// the same way if it needs to be.
func (o *OrderId) ToStateKey() []byte {
	b, err := o.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

// SortedOrders is type alias for `*OrderId` which supports deterministic
// sorting. Orders are first ordered by string comparison
// of their `Subaccount` owner, followed by integer comparison of their
// `Subaccount` number, followed by `ClientId` of the order, followed by `OrderFlags`,
// and finally by `ClobPairId` of the order.
// If two `*OrderIds` have equal Owners, Numbers, ClientIds, OrderFlags, and ClobPairId, they
// are assumed to be equal, and their sorted order is not deterministic.
type SortedOrders []OrderId

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedOrders{}

func (s SortedOrders) Len() int {
	return len(s)
}

func (s SortedOrders) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedOrders) Less(i, j int) bool {
	si := s[i]
	sj := s[j]

	if si.SubaccountId.Owner != sj.SubaccountId.Owner {
		return si.SubaccountId.Owner < sj.SubaccountId.Owner
	}

	if si.SubaccountId.Number != sj.SubaccountId.Number {
		return si.SubaccountId.Number < sj.SubaccountId.Number
	}

	if si.ClientId != sj.ClientId {
		return si.ClientId < sj.ClientId
	}

	if si.OrderFlags != sj.OrderFlags {
		return si.OrderFlags < sj.OrderFlags
	}

	if si.ClobPairId != sj.ClobPairId {
		return si.ClobPairId < sj.ClobPairId
	}

	return false
}

// MustSortAndHaveNoDuplicates is a wrapper around SortedOrders, which is for deterministic sorting.
// Mutates input slice.
// This function checks for duplicate OrderIds first, and panics if a duplicate exists.
// Orders are first ordered by string comparison
// of their `Subaccount` owner, followed by integer comparison of their
// `Subaccount` number, followed by `ClientId` of the order, followed by `OrderFlags`,
// and finally by `ClobPairId` of the order.
func MustSortAndHaveNoDuplicates(orderIds []OrderId) {
	orderIdSet := make(map[OrderId]struct{}, len(orderIds))
	for _, orderId := range orderIds {
		if _, exists := orderIdSet[orderId]; exists {
			panic(fmt.Errorf("cannot sort orders with duplicate order id %+v", orderId))
		}
		orderIdSet[orderId] = struct{}{}
	}
	sort.Sort(SortedOrders(orderIds))
}

// GetOrderIdLabels returns the telemetry labels of this order ID.
func (o *OrderId) GetOrderIdLabels() []metrics.Label {
	return []metrics.Label{
		metrics.GetLabelForIntValue(metrics.OrderFlag, int(o.GetOrderFlags())),
		metrics.GetLabelForIntValue(metrics.ClobPairId, int(o.GetClobPairId())),
	}
}
