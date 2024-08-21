package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	libante "github.com/dydxprotocol/v4-chain/protocol/lib/ante"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var DefaultWasmExecRateLimitConfig = clobtypes.MaxPerNBlocksRateLimit{
	Limit:     15,
	NumBlocks: 1,
}

func (k Keeper) RatelimitWasmExute(
	ctx sdk.Context,
	sender string,
) error {
	if !libante.ShouldRateLimit(ctx) {
		return nil
	}
	return k.wasmExecRatelimiter.RateLimit(ctx, sender)
}

func (k Keeper) PruneRateLimits(ctx sdk.Context) {
	k.wasmExecRatelimiter.PruneRateLimits(ctx)
}

func (k *Keeper) SetRateLimiter(rateLimiter rate_limit.RateLimiter[string]) {
	k.wasmExecRatelimiter = rateLimiter
}
