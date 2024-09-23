package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

func (k Keeper) VaultParams(
	goCtx context.Context,
	req *types.QueryVaultParamsRequest,
) (*types.QueryVaultParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	vaultId := types.VaultId{
		Type:   req.Type,
		Number: req.Number,
	}

	// Get vault params.
	vaultParams, exists := k.GetVaultParams(ctx, vaultId)
	if !exists {
		return nil, status.Error(codes.NotFound, "vault not found")
	}

	return &types.QueryVaultParamsResponse{
		VaultId:     vaultId,
		VaultParams: vaultParams,
	}, nil
}
