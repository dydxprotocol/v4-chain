package types

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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

// GetTotalFilledQuantums gets the total filled quantums from a MatchPerpetualDeleveraging match.
func (m *MatchPerpetualDeleveraging) GetTotalFilledQuantums() *big.Int {
	totalQuantums := big.NewInt(0)
	for _, fill := range m.GetFills() {
		totalQuantums.Add(totalQuantums, new(big.Int).SetUint64(fill.GetFillAmount()))
	}
	return totalQuantums
}
