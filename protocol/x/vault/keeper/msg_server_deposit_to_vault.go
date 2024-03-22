package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// DepositToVault deposits from a subaccount to a vault.
func (k msgServer) DepositToVault(
	goCtx context.Context,
	msg *types.MsgDepositToVault,
) (*types.MsgDepositToVaultResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	err := k.Keeper.HandleMsgDepositToVault(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositToVaultResponse{}, nil
}

// HandleMsgDepositToVault handles a MsgDepositToVault.
func (k Keeper) HandleMsgDepositToVault(
	ctx sdk.Context,
	msg *types.MsgDepositToVault,
) error {
	// Mint shares for the vault.
	err := k.MintShares(
		ctx,
		*msg.VaultId,
		msg.SubaccountId.Owner,
		msg.QuoteQuantums.BigInt(),
	)
	if err != nil {
		return err
	}

	// Transfer from sender subaccount to vault.
	// Note: Transfer should take place after minting shares for
	// shares calculation to be correct.
	err = k.subaccountsKeeper.TransferFundsFromSubaccountToSubaccount(
		ctx,
		*msg.SubaccountId,
		*msg.VaultId.ToSubaccountId(),
		0, // assetId
		msg.QuoteQuantums.BigInt(),
	)
	if err != nil {
		return err
	}

	// Emit metric on vault equity.
	equity, err := k.GetVaultEquity(ctx, *msg.VaultId)
	if err != nil {
		log.ErrorLogWithError(ctx, "Failed to get vault equity", err, "vaultId", *msg.VaultId)
	} else {
		msg.VaultId.SetGaugeWithLabels(
			metrics.VaultEquity,
			float32(equity.Int64()),
		)
	}

	return nil
}
