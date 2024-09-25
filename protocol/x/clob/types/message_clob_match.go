package types

// NewClobMatchFromMatchOrders creates a `ClobMatch` from the provided `MatchOrders`.
func NewClobMatchFromMatchOrders(
	msgMatchOrders *MatchOrders,
) *ClobMatch {
	return &ClobMatch{
		Match: &ClobMatch_MatchOrders{
			MatchOrders: msgMatchOrders,
		},
	}
}

// NewClobMatchFromMatchPerpetualLiquidation creates a `ClobMatch` from the provided
// `MatchPerpetualLiquidation`.
func NewClobMatchFromMatchPerpetualLiquidation(
	msgMatchPerpetualLiquidation *MatchPerpetualLiquidation,
) *ClobMatch {
	return &ClobMatch{
		Match: &ClobMatch_MatchPerpetualLiquidation{
			MatchPerpetualLiquidation: msgMatchPerpetualLiquidation,
		},
	}
}

// GetAllOrderIds returns a set of orderIds involved in a ClobMatch.
// It assumes the ClobMatch is valid (no duplicate orderIds in fills)
func (clobMatch *ClobMatch) GetAllOrderIds() (orderIds map[OrderId]struct{}) {
	orderIds = make(map[OrderId]struct{})
	if matchOrders := clobMatch.GetMatchOrders(); matchOrders != nil {
		orderIds[matchOrders.GetTakerOrderId()] = struct{}{}
		for _, makerFill := range matchOrders.GetFills() {
			orderIds[makerFill.GetMakerOrderId()] = struct{}{}
		}
	}
	if matchOrders := clobMatch.GetMatchPerpetualLiquidation(); matchOrders != nil {
		for _, makerFill := range matchOrders.GetFills() {
			orderIds[makerFill.GetMakerOrderId()] = struct{}{}
		}
	}
	return orderIds
}
