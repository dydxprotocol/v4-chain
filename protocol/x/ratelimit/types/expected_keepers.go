package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines the expected bank keeper used for simulations.
type BankKeeper interface {
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
}
