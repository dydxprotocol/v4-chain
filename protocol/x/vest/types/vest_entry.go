package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (entry VestEntry) Validate() error {
	if entry.VesterAccount == "" {
		return sdkerrors.Wrapf(ErrInvalidVesterAccount, "vester account cannot be empty")
	}

	if entry.TreasuryAccount == "" {
		return sdkerrors.Wrapf(ErrInvalidTreasuryAccount, "treasury account cannot be empty")
	}

	if err := sdk.ValidateDenom(entry.Denom); err != nil {
		return sdkerrors.Wrapf(ErrInvalidDenom, err.Error())
	}

	if !entry.StartTime.Before(entry.EndTime) {
		return sdkerrors.Wrapf(ErrInvalidStartAndEndTimes, "start_time = %v, end_time = %v", entry.StartTime, entry.EndTime)
	}

	if entry.StartTime.Location().String() != "UTC" {
		return sdkerrors.Wrapf(ErrInvalidTimeZone, "start_time must be in UTC")
	}

	if entry.EndTime.Location().String() != "UTC" {
		return sdkerrors.Wrapf(ErrInvalidTimeZone, "start_time must be in UTC")
	}
	return nil
}
