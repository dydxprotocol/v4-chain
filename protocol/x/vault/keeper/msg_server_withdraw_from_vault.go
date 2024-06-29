package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// DepositToVault deposits from a subaccount to a vault.
func (k msgServer) WithdrawFromVault(
	goCtx context.Context,
	msg *types.MsgWithdrawFromVault,
) (*types.MsgWithdrawFromVaultResponse, error) {
	return &types.MsgWithdrawFromVaultResponse{}, nil
}
