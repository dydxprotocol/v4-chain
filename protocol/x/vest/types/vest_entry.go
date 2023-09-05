package types

import (
	moderrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (entry VestEntry) Validate() error {
	if entry.VesterAccount == "" {
		return moderrors.Wrapf(ErrInvalidVesterAccount, "vester account cannot be empty")
	}

	if entry.TreasuryAccount == "" {
		return moderrors.Wrapf(ErrInvalidTreasuryAccount, "treasury account cannot be empty")
	}

	if err := sdk.ValidateDenom(entry.Denom); err != nil {
		return moderrors.Wrapf(ErrInvalidDenom, err.Error())
	}

	if !entry.StartTime.Before(entry.EndTime) {
		return moderrors.Wrapf(ErrInvalidStartAndEndTimes, "start_time = %v, end_time = %v", entry.StartTime, entry.EndTime)
	}

	if entry.StartTime.Location().String() != "UTC" {
		return moderrors.Wrapf(ErrInvalidTimeZone, "start_time must be in UTC")
	}

	if entry.EndTime.Location().String() != "UTC" {
		return moderrors.Wrapf(ErrInvalidTimeZone, "start_time must be in UTC")
	}
	return nil
}
