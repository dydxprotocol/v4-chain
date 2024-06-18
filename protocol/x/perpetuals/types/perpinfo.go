package types

import (
	errorsmod "cosmossdk.io/errors"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// PerpInfo contains all information needed to calculate margin requirements for a perpetual.
type PerpInfo struct {
	Perpetual     Perpetual
	Price         pricestypes.MarketPrice
	LiquidityTier LiquidityTier
}

// PerpInfos is a map of PerpInfo objects, keyed by perpetualId.
type PerpInfos map[uint32]PerpInfo

// MustGet returns the PerpInfo for the given perpetualId, or panics if it does not exist.
func (pi PerpInfos) MustGet(perpetualId uint32) PerpInfo {
	p, ok := pi[perpetualId]

	if !ok {
		panic(errorsmod.Wrapf(
			ErrPerpetualInfoDoesNotExist,
			"perpetualId: %d",
			perpetualId,
		))
	}

	return p
}
