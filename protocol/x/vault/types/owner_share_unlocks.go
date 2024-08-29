package types

import (
	"math/big"
)

// Validate validates `OwnerShareUnlocks`.
func (o OwnerShareUnlocks) Validate() error {
	if o.OwnerAddress == "" {
		return ErrEmptyOwnerAddress
	}

	return nil
}

func (o OwnerShareUnlocks) GetTotalLockedShares() *big.Int {
	totalLockedShares := big.NewInt(0)
	for _, unlock := range o.ShareUnlocks {
		totalLockedShares.Add(totalLockedShares, unlock.Shares.NumShares.BigInt())
	}
	return totalLockedShares
}
