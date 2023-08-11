package types

import (
	"github.com/dydxprotocol/v4/lib"
	satypes "github.com/dydxprotocol/v4/x/subaccounts/types"
)

// MakerFillWithOrder represents the filled amount of a matched maker order,
// along with the `Order` representing the matched maker order.
type MakerFillWithOrder struct {
	MakerFill

	// The `Order` representing the matched maker order.
	Order Order
}

// MakerFillsWithOrderToMakerFills converts a slice of `MakerFillWithOrder` to a
// slice of `MakerFill`.
func MakerFillsWithOrderToMakerFills(mfos []MakerFillWithOrder) []MakerFill {
	mfs := make([]MakerFill, len(mfos))
	for i, mf := range mfos {
		mfs[i] = mf.MakerFill
	}
	return mfs
}

// GetMakerSubaccountIds gets a list of SubaccountIds belonging to subaccounts which placed the maker orders
// in the Fills property of this MatchOrdersNew object.
func (m *MatchOrders) GetMakerSubaccountIds() []satypes.SubaccountId {
	return lib.MapSlice(
		m.Fills,
		func(fill MakerFill) satypes.SubaccountId {
			return fill.MakerOrderId.SubaccountId
		},
	)
}

// GetMakerSubaccountIds gets a list of SubaccountIds belonging to subaccounts which placed the maker orders
// in the Fills property of this MatchPerpetualLiquidationNew object.
func (m *MatchPerpetualLiquidation) GetMakerSubaccountIds() []satypes.SubaccountId {
	return lib.MapSlice(
		m.Fills,
		func(fill MakerFill) satypes.SubaccountId {
			return fill.MakerOrderId.SubaccountId
		},
	)
}
