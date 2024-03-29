package keeper

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func (k Keeper) Vault(
	goCtx context.Context,
	req *types.QueryVaultRequest,
) (*types.QueryVaultResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	vaultId := types.VaultId{
		Type:   req.Type,
		Number: req.Number,
	}

	// Get total shares.
	totalShares, exists := k.GetTotalShares(ctx, vaultId)
	if !exists {
		return nil, status.Error(codes.NotFound, "vault not found")
	}

	// Get vault equity.
	equity, err := k.GetVaultEquity(ctx, vaultId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get vault equity")
	}

	// Get vault inventory.
	clobPair, exists := k.clobKeeper.GetClobPair(ctx, clobtypes.ClobPairId(vaultId.Number))
	if !exists {
		return nil, status.Error(codes.Internal, fmt.Sprintf("clob pair %d doesn't exist", vaultId.Number))
	}
	perpId := clobPair.Metadata.(*clobtypes.ClobPair_PerpetualClobMetadata).PerpetualClobMetadata.PerpetualId
	inventory := k.GetVaultInventoryInPerpetual(ctx, vaultId, perpId)

	return &types.QueryVaultResponse{
		VaultId:      vaultId,
		SubaccountId: *vaultId.ToSubaccountId(),
		Equity:       equity.Uint64(),
		Inventory:    inventory.Uint64(),
		TotalShares:  totalShares.NumShares.BigInt().Uint64(),
	}, nil
}
