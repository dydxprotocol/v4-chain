package events

// NewUpdatePerpetualEventV1 creates a UpdatePerpetualEventV1 representing
// update of a perpetual.
func NewUpdatePerpetualEventV1(
	id uint32,
	ticker string,
	marketId uint32,
	atomicResolution int32,
	liquidityTier uint32,
	dangerIndexPpm uint32,
	isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock string,
) *UpdatePerpetualEventV1 {
	return &UpdatePerpetualEventV1{
		Id:               id,
		Ticker:           ticker,
		MarketId:         marketId,
		AtomicResolution: atomicResolution,
		LiquidityTier:    liquidityTier,
		DangerIndexPpm:   dangerIndexPpm,
		IsolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
	}
}
