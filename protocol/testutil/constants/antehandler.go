package constants

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	EmptyAnteHandler = func(c sdk.Context, t sdk.Tx, b bool) (sdk.Context, error) { return c, nil }
)
