package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/rate_limit"
)

type RatelimitKeeper interface {
	RatelimitWasmExute(
		ctx sdk.Context,
		sender string,
	) error
	SetRateLimiter(rateLimiter rate_limit.RateLimiter[string])
}
