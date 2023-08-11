package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		// Corresponds to module account address: dydx16wrau2x4tsg033xfrrdpae6kxfn9kyuerr5jjp
		TreasuryAccount: "rewards_treasury",
		// The exact denom to be used for rewards is TBD, so using eth as a placeholder.
		// Note that `eth_rewards_denom` is not an actual denom for a coin, since any
		// form of Ethereum token will likely exist as an IBC token.
		Denom:            "eth_rewards_denom",
		DenomExponent:    -6,
		MarketId:         1,
		FeeMultiplierPpm: 990_000, // 0.99
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.TreasuryAccount == "" {
		return fmt.Errorf("treasury account cannot have empty name")
	}
	if err := sdk.ValidateDenom(p.Denom); err != nil {
		return err
	}
	return nil
}
