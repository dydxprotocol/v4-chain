package events

import "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

// NewPerpetualMarketCreateEvent creates a PerpetualMarketCreateEvent
// representing creation of a perpetual market.
func NewPerpetualMarketCreateEvent(
	id uint32,
	clobPairId uint32,
	ticker string,
	marketId uint32,
	status types.ClobPair_Status,
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
		Status:                    convertToClobPairStatus(status),
		QuantumConversionExponent: quantumConversionExponent,
		AtomicResolution:          atomicResolution,
		SubticksPerTick:           subticksPerTick,
		MinOrderBaseQuantums:      minOrderBaseQuantums,
		StepBaseQuantums:          stepBaseQuantums,
		LiquidityTier:             liquidityTier,
	}
}

func convertToClobPairStatus(status types.ClobPair_Status) ClobPairStatus {
	switch status {
	case types.ClobPair_STATUS_UNSPECIFIED:
		return ClobPairStatus_CLOB_PAIR_STATUS_UNSPECIFIED
	case types.ClobPair_STATUS_ACTIVE:
		return ClobPairStatus_CLOB_PAIR_STATUS_ACTIVE
	case types.ClobPair_STATUS_PAUSED:
		return ClobPairStatus_CLOB_PAIR_STATUS_PAUSED
	case types.ClobPair_STATUS_CANCEL_ONLY:
		return ClobPairStatus_CLOB_PAIR_STATUS_CANCEL_ONLY
	case types.ClobPair_STATUS_POST_ONLY:
		return ClobPairStatus_CLOB_PAIR_STATUS_POST_ONLY
	default:
		panic("invalid clob pair status")
	}
}
