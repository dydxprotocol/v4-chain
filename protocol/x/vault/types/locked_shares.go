package types

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
)

// Validate validates `LockedShares`.
func (l LockedShares) Validate() error {
	if l.OwnerAddress == "" {
		return errorsmod.Wrapf(ErrInvalidLockedShares, "empty owner address")
	}

	totalSharesToLock := big.NewInt(0)
	for _, unlockDetail := range l.UnlockDetails {
		totalSharesToLock.Add(totalSharesToLock, unlockDetail.Shares.NumShares.BigInt())
	}
	totalSharesLocked := l.TotalLockedShares.NumShares.BigInt()
	if totalSharesToLock.Cmp(totalSharesLocked) != 0 {
		return errorsmod.Wrapf(
			ErrInvalidLockedShares,
			"total shares locked (%s) not equal to total shares to unlock (%s)",
			totalSharesLocked,
			totalSharesToLock,
		)
	}

	return nil
}
