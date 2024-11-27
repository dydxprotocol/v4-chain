package events

import (
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// NewUpdatePerpetualEvent creates a UpdatePerpetualEventV2 representing
// update of a perpetual.
func NewUpdatePerpetualEvent(
	id uint32,
	ticker string,
	marketId uint32,
	atomicResolution int32,
	liquidityTier uint32,
	marketType perptypes.PerpetualMarketType,
) *UpdatePerpetualEventV2 {
	return &UpdatePerpetualEventV2{
		Id:               id,
		Ticker:           ticker,
		MarketId:         marketId,
		AtomicResolution: atomicResolution,
		LiquidityTier:    liquidityTier,
		MarketType:       v1.ConvertToPerpetualMarketType(marketType),
	}
}
