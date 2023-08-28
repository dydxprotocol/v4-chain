package types

import (
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// SupportedClobPairStatusTransitions has keys corresponding to currently-supported
// ClobPairStatus types with values equal to the set of ClobPairStatus types that
// may be transitioned to from this state. Note the keys of this map may be
// a subset of the types defined in the proto for ClobPairStatus.
var SupportedClobPairStatusTransitions = map[ClobPairStatus]map[ClobPairStatus]struct{}{
	ClobPairStatus_ACTIVE: {},
	ClobPairStatus_INITIALIZING: {
		ClobPairStatus_ACTIVE: struct{}{},
	},
}

// IsSupportedClobPairStatus returns true if the provided ClobPairStatus is in the list
// of currently supported ClobPairStatus types. Else, returns false.
func IsSupportedClobPairStatus(clobPairStatus ClobPairStatus) bool {
	_, exists := SupportedClobPairStatusTransitions[clobPairStatus]
	return exists
}

// IsSupportedClobPairStatusTransition returns true if it is considered valid to transition from
// the first provided ClobPairStatus to the second provided ClobPairStatus. Else, returns false.
func IsSupportedClobPairStatusTransition(from ClobPairStatus, to ClobPairStatus) bool {
	_, exists := SupportedClobPairStatusTransitions[from][to]
	return exists
}

func (c *ClobPair) GetClobPairSubticksPerTick() SubticksPerTick {
	return SubticksPerTick(c.SubticksPerTick)
}

func (c *ClobPair) GetClobPairMinOrderBaseQuantums() satypes.BaseQuantums {
	return satypes.BaseQuantums(c.StepBaseQuantums)
}

// GetPerpetualId returns the `PerpetualId` for the provided `clobPair`.
func (c *ClobPair) GetPerpetualId() (uint32, error) {
	perpetualClobMetadata := c.GetPerpetualClobMetadata()
	if perpetualClobMetadata == nil {
		return 0, ErrAssetOrdersNotImplemented
	}

	return perpetualClobMetadata.PerpetualId, nil
}

// MustGetPerpetualId returns the `PerpetualId` for the provided `clobPair`.
// Will panic if `GetPerpetualId` returns an error.
func (c *ClobPair) MustGetPerpetualId() uint32 {
	id, err := c.GetPerpetualId()
	if err != nil {
		panic(err)
	}
	return id
}

// GetId returns the `ClobPairId` for the provided `clobPair`.
func (c *ClobPair) GetClobPairId() ClobPairId {
	return ClobPairId(c.Id)
}
