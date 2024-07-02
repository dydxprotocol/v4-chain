package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// DepositToVault deposits from a subaccount to a vault.
func (k msgServer) DepositToVault(
	goCtx context.Context,
	msg *types.MsgDepositToVault,
) (*types.MsgDepositToVaultResponse, error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	quoteQuantums := msg.QuoteQuantums.BigInt()

	// Mint vault shares for the depositor.
	err := k.MintShares(
		ctx,
		*msg.VaultId,
		msg.SubaccountId.Owner,
		quoteQuantums,
	)
	if err != nil {
		return nil, err
	}

	// Transfer from depositor's subaccount to vault.
	// IMPORTANT: Transfer should take place after minting shares for
	// shares calculation to be correct. This is because minting shares
	// depends on the vault's current equity. Therefore, if you transfer
	// before minting shares, then minting shares will be based on the
	// vault's equity after the deposit.
	err = k.sendingKeeper.ProcessTransfer(
		ctx,
		&sendingtypes.Transfer{
			Sender:    *msg.SubaccountId,
			Recipient: *msg.VaultId.ToSubaccountId(),
			AssetId:   assettypes.AssetUsdc.Id,
			Amount:    quoteQuantums.Uint64(),
		},
	)
	if err != nil {
		return nil, err
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

	return &types.MsgDepositToVaultResponse{}, nil
}
