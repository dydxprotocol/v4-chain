package events

import (
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// NewPerpetualMarketCreateEvent creates a PerpetualMarketCreateEvent
// representing creation of a perpetual market.
func NewPerpetualMarketCreateEvent(
	id uint32,
	clobPairId uint32,
	ticker string,
	marketId uint32,
	status types.ClobPairStatus,
	quantumConversionExponent int32,
	atomicResolution int32,
	subticksPerTick uint32,
	minOrderBaseQuantums uint64,
	stepBaseQuantums uint64,
	liquidityTier uint32,
) *PerpetualMarketCreateEventV1 {
	return &PerpetualMarketCreateEventV1{
		Id:                        id,
		ClobPairId:                clobPairId,
		Ticker:                    ticker,
		MarketId:                  marketId,
		Status:                    v1.ConvertToClobPairStatus(status),
		QuantumConversionExponent: quantumConversionExponent,
		AtomicResolution:          atomicResolution,
		SubticksPerTick:           subticksPerTick,
		MinOrderBaseQuantums:      minOrderBaseQuantums,
		StepBaseQuantums:          stepBaseQuantums,
		LiquidityTier:             liquidityTier,
	}
}
