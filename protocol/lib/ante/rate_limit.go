package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// The rate limiting is only performed during `CheckTx`.
// Rate limiting during `ReCheckTx` might result in over counting.
func ShouldRateLimit(ctx sdk.Context) bool {
	return ctx.IsCheckTx() && !ctx.IsReCheckTx()
}
