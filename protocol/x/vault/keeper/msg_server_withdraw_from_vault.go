package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// WithdrawFromVault attempts to withdraw an asset from a vault to a subaccount.
func (k msgServer) WithdrawFromVault(
	goCtx context.Context,
	msg *types.MsgWithdrawFromVault,
) (*types.MsgWithdrawFromVaultResponse, error) {
	// TODO(TRA-461): Validate.
	// TODO(TRA-462): Calculate effective amount to withdraw + shares to redeem with slippage and user equity.
	// TODO(TRA-461): Redeem shares for the vault.
	// TODO(TRA-461): Transfer asset from vault to recipient subaccount.
	// should transfer happen after redeeming shares? why?
	// TODO(TRA-461): emit metric on vault equity.
	// TODO(TRA-461): Get info on shares after the withdrawal for the response.
	return &types.MsgWithdrawFromVaultResponse{}, nil
}
