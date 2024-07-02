package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// DefaultParams returns a default set of `x/revshare` market mapper parameters.
func DefaultParams() MarketMapperRevenueShareParams {
	return MarketMapperRevenueShareParams{
		Address:         authtypes.NewModuleAddress(authtypes.FeeCollectorName).String(),
		RevenueSharePpm: 0,
		ValidDays:       0,
	}
}

// Validate validates `x/revshare` parameters.
func (p MarketMapperRevenueShareParams) Validate() error {
	// Address must be a valid address
	_, err := types.AccAddressFromBech32(p.Address)
	if err != nil {
		return ErrInvalidAddress
	}

	// Revenue share ppm must be less than 1000000 (100%)
	if p.RevenueSharePpm >= 1000000 {
		return ErrInvalidRevenueSharePpm
	}

	return nil
}
