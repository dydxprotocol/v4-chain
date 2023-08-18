package types

import (
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

const MaxFeePpm = 100000 // 10%

// SupportedClobPairStatusTransitions has keys corresponding to currently-supported
// ClobPair_Status types with values equal to the set of ClobPair_Status types that
// may be transitioned to from this state. Note the keys of this map may be
// a subset of the types defined in the proto for ClobPair_Status.
var SupportedClobPairStatusTransitions = map[ClobPair_Status]map[ClobPair_Status]struct{}{
	ClobPair_STATUS_ACTIVE: {},
	ClobPair_STATUS_POST_ONLY: {
		ClobPair_STATUS_ACTIVE: struct{}{},
	},
}

// IsSupportedClobPairStatus returns true if the provided ClobPair_Status is in the list
// of currently supported ClobPair_Status types. Else, returns false.
func IsSupportedClobPairStatus(clobPairStatus ClobPair_Status) bool {
	_, exists := SupportedClobPairStatusTransitions[clobPairStatus]
	return exists
}

// IsSupportedClobPairStatusTransition returns true if it is considered valid to transition from
// the first provided ClobPair_Status to the second provided ClobPair_Status. Else, returns false.
func IsSupportedClobPairStatusTransition(from ClobPair_Status, to ClobPair_Status) bool {
	_, exists := SupportedClobPairStatusTransitions[from][to]
	return exists
}

func (c *ClobPair) GetClobPairSubticksPerTick() SubticksPerTick {
	return SubticksPerTick(c.SubticksPerTick)
}

func (c *ClobPair) GetClobPairMinOrderBaseQuantums() satypes.BaseQuantums {
	return satypes.BaseQuantums(c.MinOrderBaseQuantums)
}

// Get fee rate in ppm. Returns the taker fee for taker orders, otherwise returns the maker fee.
func (c *ClobPair) GetFeePpm(isTaker bool) uint32 {
	if isTaker {
		return c.TakerFeePpm
	}
	return c.MakerFeePpm
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
