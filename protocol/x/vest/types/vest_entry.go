package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (entry VestEntry) Validate() error {
	if entry.VesterAccount == "" {
		return errorsmod.Wrap(ErrInvalidVesterAccount, "vester account cannot be empty")
	}

	if entry.TreasuryAccount == "" {
		return errorsmod.Wrap(ErrInvalidTreasuryAccount, "treasury account cannot be empty")
	}

	if err := sdk.ValidateDenom(entry.Denom); err != nil {
		return errorsmod.Wrap(ErrInvalidDenom, err.Error())
	}

	if !entry.StartTime.Before(entry.EndTime) {
		return errorsmod.Wrapf(ErrInvalidStartAndEndTimes, "start_time = %v, end_time = %v", entry.StartTime, entry.EndTime)
	}

	if entry.StartTime.Location().String() != "UTC" {
		return errorsmod.Wrap(ErrInvalidTimeZone, "start_time must be in UTC")
	}

	if entry.EndTime.Location().String() != "UTC" {
		return errorsmod.Wrap(ErrInvalidTimeZone, "end_time must be in UTC")
	}
	return nil
}
