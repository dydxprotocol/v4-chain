package lib

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MustParseCoinsNormalized(coinStr string) sdk.Coins {
	coins, err := sdk.ParseCoinsNormalized(coinStr)
	if err != nil {
		panic(err)
	}
	return coins
}
