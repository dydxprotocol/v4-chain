package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/protectedaccounts"
)

func (entry VestEntry) Validate() error {
	if entry.VesterAccount == "" {
		return errorsmod.Wrap(ErrInvalidVesterAccount, "vester account cannot be empty")
	}

	if entry.TreasuryAccount == "" {
		return errorsmod.Wrap(ErrInvalidTreasuryAccount, "treasury account cannot be empty")
	}

	if protectedaccounts.IsProtectedModuleName(entry.VesterAccount) {
		return errorsmod.Wrapf(
			ErrInvalidVesterAccount,
			"vester account '%s' is a protected module account",
			entry.VesterAccount,
		)
	}

	if protectedaccounts.IsProtectedModuleName(entry.TreasuryAccount) {
		return errorsmod.Wrapf(
			ErrInvalidTreasuryAccount,
			"treasury account '%s' is a protected module account",
			entry.TreasuryAccount,
		)
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
