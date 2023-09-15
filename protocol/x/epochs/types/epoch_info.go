package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
)

// Validate does basic validation for epoch info.
func (epoch EpochInfo) Validate() error {
	if epoch.Name == "" {
		return errorsmod.Wrap(ErrEmptyEpochInfoName, "EpochInfo Name is empty")
	}
	if epoch.Duration == 0 {
		return errorsmod.Wrap(ErrDurationIsZero, "Duration is zero")
	}
	// `CurrentEpoch` should be zero if and only if `CurrentEpochStartBlock` is zero.
	if (epoch.CurrentEpoch == 0) != (epoch.CurrentEpochStartBlock == 0) {
		return errorsmod.Wrap(
			ErrInvalidCurrentEpochAndCurrentEpochStartBlockTuple,
			fmt.Sprintf(
				"CurrentEpoch: %d, CurrentEpochStartBlock: %v",
				epoch.CurrentEpoch,
				epoch.CurrentEpochStartBlock,
			),
		)
	}
	return nil
}

// GetEpochInfoName returns Id from epoch info.
func (epoch EpochInfo) GetEpochInfoName() EpochInfoName {
	return EpochInfoName(epoch.Name)
}
