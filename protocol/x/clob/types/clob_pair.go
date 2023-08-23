package types

import (
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

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
