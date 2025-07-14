package types

import (
	"crypto/sha256"
	"math/big"

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

// ToStreamLiquidationOrder converts the LiquidationOrder to a StreamLiquidationOrder
// to be emitted by full node streaming.
func (lo *LiquidationOrder) ToStreamLiquidationOrder() *StreamLiquidationOrder {
	return &StreamLiquidationOrder{
		LiquidationInfo: &lo.perpetualLiquidationInfo,
		ClobPairId:      uint32(lo.clobPairId),
		IsBuy:           lo.isBuy,
		Quantums:        lo.quantums.ToUint64(),
		Subticks:        lo.subticks.ToUint64(),
	}
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

// MustGetLiquidationOrder returns the underlying `LiquidationOrder` type.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (lo *LiquidationOrder) MustGetLiquidationOrder() LiquidationOrder {
	return *lo
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

// GetOrderRouterAddress returns the order router address of this liquidation order.
// This function is necessary for the `LiquidationOrder` type to implement the `MatchableOrder` interface.
func (o *LiquidationOrder) GetOrderRouterAddress() string {
	return ""
}

// GetDeltaQuantums returns the delta quantums of this liquidation order.
func (o *LiquidationOrder) GetDeltaQuantums() *big.Int {
	deltaQuantums := o.GetBaseQuantums().ToBigInt()
	if !o.IsBuy() {
		deltaQuantums.Neg(deltaQuantums)
	}
	return deltaQuantums
}
