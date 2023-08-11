package prepare

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ ContextHelper = &ContextHelperImpl{}

// TODO(DEC-1248): delete once height == 0 issue is resolved.
type ContextHelper interface {
	Height(ctx sdk.Context) int64
}

// ContextHelperImpl implements ContextHelper interface.
type ContextHelperImpl struct{}

// Height returns current block height.
func (t *ContextHelperImpl) Height(ctx sdk.Context) int64 {
	return ctx.BlockHeight()
}
