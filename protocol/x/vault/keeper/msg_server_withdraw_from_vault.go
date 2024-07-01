package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// TODO: implement
func (k msgServer) WithdrawFromVault(
	goCtx context.Context,
	msg *types.MsgWithdrawFromVault,
) (*types.MsgWithdrawFromVaultResponse, error) {
	return &types.MsgWithdrawFromVaultResponse{}, nil
}
