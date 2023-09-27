package lib

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MustParseCoinsNormalized parses a string of coins and panics on error.
func MustParseCoinsNormalized(coinStr string) sdk.Coins {
	coins, err := sdk.ParseCoinsNormalized(coinStr)
	if err != nil {
		panic(err)
	}
	return coins
}
