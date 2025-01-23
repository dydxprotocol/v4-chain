package events

import (
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

// NewPerpetualMarketCreateEvent creates a PerpetualMarketCreateEvent
// representing creation of a perpetual market.
func NewPerpetualMarketCreateEvent(
	id uint32,
	clobPairId uint32,
	ticker string,
	marketId uint32,
	status clobtypes.ClobPair_Status,
	quantumConversionExponent int32,
	atomicResolution int32,
	subticksPerTick uint32,
	stepBaseQuantums uint64,
	liquidityTier uint32,
	marketType perptypes.PerpetualMarketType,
	defaultFundingPpm int32,
) *PerpetualMarketCreateEventV3 {
	return &PerpetualMarketCreateEventV3{
		Id:                        id,
		ClobPairId:                clobPairId,
		Ticker:                    ticker,
		MarketId:                  marketId,
		Status:                    v1.ConvertToClobPairStatus(status),
		QuantumConversionExponent: quantumConversionExponent,
		AtomicResolution:          atomicResolution,
		SubticksPerTick:           subticksPerTick,
		StepBaseQuantums:          stepBaseQuantums,
		LiquidityTier:             liquidityTier,
		MarketType:                v1.ConvertToPerpetualMarketType(marketType),
		DefaultFunding8HrPpm:      defaultFundingPpm,
	}
}
