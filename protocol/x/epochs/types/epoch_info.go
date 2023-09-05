package types

import (
	moderrors "cosmossdk.io/errors"
	"fmt"
)

// Validate does basic validation for epoch info.
func (epoch EpochInfo) Validate() error {
	if epoch.Name == "" {
		return moderrors.Wrap(ErrEmptyEpochInfoName, "EpochInfo Name is empty")
	}
	if epoch.Duration == 0 {
		return moderrors.Wrap(ErrDurationIsZero, "Duration is zero")
	}
	// `CurrentEpoch` should be zero if and only if `CurrentEpochStartBlock` is zero.
	if (epoch.CurrentEpoch == 0) != (epoch.CurrentEpochStartBlock == 0) {
		return moderrors.Wrap(
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
