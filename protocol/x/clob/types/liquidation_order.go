package types

import (
	"bytes"
	"crypto/sha256"
	"math/big"
	"sort"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var _ MatchableOrder = &LiquidationOrder{}

// LiquidationOrder is used to represent an IOC liquidation order.
type LiquidationOrder struct {
	// Information about this liquidation order.
	perpetualLiquidationInfo PerpetualLiquidationInfo
	// CLOB pair ID of the CLOB pair the liquidation order will be matched against.
	clobPairId ClobPairId
	// True if this is a buy order liquidating a short position, false if vice versa.
	isBuy bool
	// The number of base quantums for this liquidation order.
	quantums satypes.BaseQuantums
	// The subticks this order will be submitted at.
	subticks Subticks
}

// NewLiquidationOrder creates and returns a new liquidation order.
// This function will panic if the caller attempts to create a liquidation order with a non-perpetual
// CLOB pair.
func NewLiquidationOrder(
	subaccountId satypes.SubaccountId,
	clobPair ClobPair,
	isBuy bool,
	quantums satypes.BaseQuantums,
	subticks Subticks,
) *LiquidationOrder {
	// If this is not a perpetual CLOB, panic.
	perpetualClobMetadata := clobPair.GetPerpetualClobMetadata()
	if perpetualClobMetadata == nil {
		panic("NewLiquidationOrder: Attempting to create liquidation order with a non-perpetual CLOB pair")
	}

	return &LiquidationOrder{
		perpetualLiquidationInfo: PerpetualLiquidationInfo{
			SubaccountId: subaccountId,
			PerpetualId:  perpetualClobMetadata.PerpetualId,
		},
		clobPairId: clobPair.GetClobPairId(),
		isBuy:      isBuy,
		quantums:   quantums,
		subticks:   subticks,
	}
}

// SortedLiquidationOrders is type alias for `*LiquidationOrder` which supports deterministic
// sorting. Orders are first ordered by `ClobPairId` of the order,
// followed by the side of the order, followed by the fillable price, followed by the size of the order,
// and finally by order hashes.
type SortedLiquidationOrders []LiquidationOrder

// The below methods are required to implement `sort.Interface` for sorting using the sort package.
var _ sort.Interface = SortedLiquidationOrders{}

func (s SortedLiquidationOrders) Len() int {
	return len(s)
}

func (s SortedLiquidationOrders) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortedLiquidationOrders) Less(i, j int) bool {
	x := s[i]
	y := s[j]
	if x.GetClobPairId() != y.GetClobPairId() {
		return x.GetClobPairId() < y.GetClobPairId()
	}

	// Buy orders before sell orders.
	if x.IsBuy() != y.IsBuy() {
		return x.IsBuy()
	}

	if x.GetOrderSubticks() != y.GetOrderSubticks() {
		if x.IsBuy() {
			// Buy orders are sorted in descending order.
			return x.GetOrderSubticks() > y.GetOrderSubticks()
		} else {
			// Sell orders are sorted in ascending order.
			return x.GetOrderSubticks() < y.GetOrderSubticks()
		}
	}

	// Sort by order size.
	if x.GetBaseQuantums() != y.GetBaseQuantums() {
		return x.GetBaseQuantums() > y.GetBaseQuantums()
	}

	xHash := x.GetOrderHash()
	yHash := y.GetOrderHash()
	return bytes.Compare(xHash[:], yHash[:]) == -1
}

// IsBuy returns true if this is a buy order, false if not.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) IsBuy() bool {
	return lo.isBuy
}

// GetOrderHash returns the SHA256 hash of the `PerpetualLiquidationInfo` field.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) GetOrderHash() OrderHash {
	perpetualLiquidationInfoBytes, err := lo.perpetualLiquidationInfo.Marshal()
	if err != nil {
		panic(err)
	}
	return sha256.Sum256(perpetualLiquidationInfoBytes)
}

// GetBaseQuantums returns the quantums of this liquidation order.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) GetBaseQuantums() satypes.BaseQuantums {
	return satypes.BaseQuantums(lo.quantums)
}

// GetClobPairId returns the CLOB pair ID of this liquidation order.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) GetClobPairId() ClobPairId {
	return lo.clobPairId
}

// GetOrderSubticks returns the subticks of this liquidation order.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) GetOrderSubticks() Subticks {
	return Subticks(lo.subticks)
}

// GetSubaccountId returns the subaccount ID that is being liquidated.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) GetSubaccountId() satypes.SubaccountId {
	return lo.perpetualLiquidationInfo.SubaccountId
}

// IsLiquidation always returns true since this order is a liquidation order.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) IsLiquidation() bool {
	return true
}

// MustGetOrder always panics since there is no underlying `Order` type for a liquidation.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) MustGetOrder() Order {
	panic("MustGetOrder: No underlying order on a LiquidationOrder type.")
}

// MustGetLiquidatedPerpetualId returns the perpetual ID that this perpetual order is liquidating.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) MustGetLiquidatedPerpetualId() uint32 {
	return lo.perpetualLiquidationInfo.PerpetualId
}

// IsReduceOnly returns whether this is a reduce-only order. This always returns false
// for liquidation orders.
func (o *LiquidationOrder) IsReduceOnly() bool {
	return false
}

// GetDeltaQuantums returns the delta quantums of this liquidation order.
func (o *LiquidationOrder) GetDeltaQuantums() *big.Int {
	deltaQuantums := o.GetBaseQuantums().ToBigInt()
	if !o.IsBuy() {
		deltaQuantums.Neg(deltaQuantums)
	}
	return deltaQuantums
}
