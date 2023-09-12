package clob

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// MustPerpetualId is a wrapper around ClobPair.GetPerpetualId() which panics if an error is returned.
func MustPerpetualId(clobPair clobtypes.ClobPair) uint32 {
	perpetualId, err := clobPair.GetPerpetualId()
	if err != nil {
		panic(err)
	}
	return perpetualId
}

type ClobModifierOption func(cp *clobtypes.ClobPair)

func WithId(id uint32) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.Id = id
	}
}

func WithStepBaseQuantums(bq satypes.BaseQuantums) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.StepBaseQuantums = bq.ToUint64()
	}
}

func WithStatus(status clobtypes.ClobPair_Status) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.Status = status
	}
}

func WithSubticksPerTick(subticks uint32) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.SubticksPerTick = subticks
	}
}

func WithQuantumConversionExponent(exponent int32) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.QuantumConversionExponent = exponent
	}
}

func WithPerpetualMetadata(metadata *clobtypes.ClobPair_PerpetualClobMetadata) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.Metadata = metadata
	}
}

func WithPerpetualId(perpetualId uint32) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.Metadata = &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: perpetualId,
			},
		}
	}
}

func WithSpotMetadata(metadata *clobtypes.ClobPair_SpotClobMetadata) ClobModifierOption {
	return func(cp *clobtypes.ClobPair) {
		cp.Metadata = metadata
	}
}

// GenerateClobPair returns a `ClobPair` object set to default values.
// Passing in `ClobModifierOption` methods alters the value of the `ClobPair` returned.
// It will start with the default, valid `ClobPair` value defined within the method
// and make the requested modifications before returning the object.
//
// Example usage:
// `GenerateClobPair(WithQuantumConversionExponent(25))`
// This will start with the default `ClobPair` object defined within the method and
// return the newly-created object after overriding the values of
// `QuantumConversionExponent` to 25.
func GenerateClobPair(optionalModifications ...ClobModifierOption) *clobtypes.ClobPair {
	clobPair := &clobtypes.ClobPair{
		Metadata: &clobtypes.ClobPair_PerpetualClobMetadata{
			PerpetualClobMetadata: &clobtypes.PerpetualClobMetadata{
				PerpetualId: 0,
			},
		},
		Id:                        0,
		StepBaseQuantums:          5,
		SubticksPerTick:           10,
		QuantumConversionExponent: -8,
		Status:                    clobtypes.ClobPair_STATUS_ACTIVE,
	}

	for _, opt := range optionalModifications {
		opt(clobPair)
	}

	return clobPair
}
