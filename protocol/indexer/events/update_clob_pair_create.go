package events

import (
	v1 "github.com/dydxprotocol/v4-chain/protocol/indexer/protocol/v1"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// NewUpdateClobPairEvent creates a UpdateClobPairEventV1
// representing update of a clob pair.
func NewUpdateClobPairEvent(
	clobPairId uint32,
	status types.ClobPair_Status,
	quantumConversionExponent int32,
	subticksPerTick uint32,
	stepBaseQuantums uint64,
) *UpdateClobPairEventV1 {
	return &UpdateClobPairEventV1{
		ClobPairId:                clobPairId,
		Status:                    v1.ConvertToClobPairStatus(status),
		QuantumConversionExponent: quantumConversionExponent,
		SubticksPerTick:           subticksPerTick,
		StepBaseQuantums:          stepBaseQuantums,
	}
}
