package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// BatchCancel performs batch order cancellation functionality for short term orders.
// For now, MsgBatchCancel only handles short term orders, so this code path should never
// be reached since the message is excluded from the mempool.
func (k msgServer) BatchCancel(
	goCtx context.Context,
	msg *types.MsgBatchCancel,
) (resp *types.MsgBatchCancelResponse, err error) {
	return &types.MsgBatchCancelResponse{}, nil
}
