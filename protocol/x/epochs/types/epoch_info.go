package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Validate does basic validation for epoch info.
func (epoch EpochInfo) Validate() error {
	if epoch.Name == "" {
		return sdkerrors.Wrap(ErrEmptyEpochInfoName, "EpochInfo Name is empty")
	}
	if epoch.Duration == 0 {
		return sdkerrors.Wrap(ErrDurationIsZero, "Duration is zero")
	}
	// `CurrentEpoch` should be zero if and only if `CurrentEpochStartBlock` is zero.
	if (epoch.CurrentEpoch == 0) != (epoch.CurrentEpochStartBlock == 0) {
		return sdkerrors.Wrap(
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
