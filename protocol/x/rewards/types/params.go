package types

import (
	sdkerrors "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		// Corresponds to module account address: dydx16wrau2x4tsg033xfrrdpae6kxfn9kyuerr5jjp
		TreasuryAccount:  "rewards_treasury",
		Denom:            "dv4tnt",
		DenomExponent:    -6,
		MarketId:         1,
		FeeMultiplierPpm: 990_000, // 0.99
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.TreasuryAccount == "" {
		return sdkerrors.Wrap(ErrInvalidTreasuryAccount, "treasury account cannot have empty name")
	}

	if p.FeeMultiplierPpm > lib.OneMillion {
		return sdkerrors.Wrap(ErrInvalidFeeMultiplierPpm, "FeeMultiplierPpm cannot be greater than 1_000_000 (100%)")
	}

	if err := sdk.ValidateDenom(p.Denom); err != nil {
		return err
	}
	return nil
}
