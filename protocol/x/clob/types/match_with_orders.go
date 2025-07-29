package types

import (
	"errors"

	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// MatchWithOrders represents a match which occurred between two orders and the amount that was matched.
type MatchWithOrders struct {
	MakerOrder          MatchableOrder
	TakerOrder          MatchableOrder
	FillAmount          satypes.BaseQuantums
	MakerFee            int64
	TakerFee            int64
	MakerBuilderFee     uint64
	TakerBuilderFee     uint64
	MakerOrderRouterFee uint64
	TakerOrderRouterFee uint64
}

// Validate performs stateless validation on an order match. This validation does
// not perform any state reads, or memclob reads.
//
// This validation ensures:
//   - Order match does not constitute a self-trade.
//   - Order match contains a `fillAmount` greater than 0.
//   - Orders in match are for the same `ClobPairId`.
//   - Orders in match are for opposing sides.
//   - Orders are crossing.
//   - The minimum of the takerOrder and makerOrder initial quantums does not exceed the FillAmount.
//   - The maker order referenced in the match is not a liquidation order.
//   - The maker order referenced in the match is not an IOC order.
func (match MatchWithOrders) Validate() error {
	makerOrder := match.MakerOrder
	takerOrder := match.TakerOrder
	fillAmount := match.FillAmount
	// Make sure the maker and taker order are not for the same Subaccount.
	if makerOrder.GetSubaccountId() == takerOrder.GetSubaccountId() {
		return errors.New("Match constitutes a self-trade")
	}

	// Make sure the fill amount is greater than zero.
	if fillAmount == 0 {
		return errors.New("fillAmount must be greater than 0")
	}

	// Make sure the maker and taker order are for the same `ClobPair`.
	if makerOrder.GetClobPairId() != takerOrder.GetClobPairId() {
		return errors.New("ClobPairIds do not match in match")
	}

	// Make sure the maker and taker order are for opposing sides of the book.
	if makerOrder.IsBuy() == takerOrder.IsBuy() {
		return errors.New("Orders are not on opposing sides of the book in match")
	}

	// Make sure the maker and taker order cross.
	if makerOrder.IsBuy() {
		if makerOrder.GetOrderSubticks() < takerOrder.GetOrderSubticks() {
			return errors.New("Orders do not cross in match")
		}
	} else {
		if takerOrder.GetOrderSubticks() < makerOrder.GetOrderSubticks() {
			return errors.New("Orders do not cross in match")
		}
	}

	// Verify that the minimum of the `makerOrder` and `takerOrder` initial quantums does not exceed the `fillAmount`.
	if fillAmount > makerOrder.GetBaseQuantums() || fillAmount > takerOrder.GetBaseQuantums() {
		return errors.New("Minimum initial order quantums exceeds fill amount")
	}

	// Make sure the maker order is not a liquidation order.
	if makerOrder.IsLiquidation() {
		return errors.New("Liquidation order cannot be matched as a maker order")
	}

	// Make sure the maker order is not an IOC order.
	unwrappedMakerOrder := makerOrder.MustGetOrder()
	if unwrappedMakerOrder.RequiresImmediateExecution() {
		return errors.New("IOC order cannot be matched as a maker order")
	}
	return nil
}
