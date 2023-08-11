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
