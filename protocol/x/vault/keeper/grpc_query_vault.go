package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/store/prefix"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
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
		Equity:       dtypes.NewIntFromBigInt(equity),
		Inventory:    dtypes.NewIntFromBigInt(inventory),
		TotalShares:  totalShares.NumShares.BigInt().Uint64(),
	}, nil
}

func (k Keeper) AllVaults(
	c context.Context,
	req *types.QueryAllVaultsRequest,
) (*types.QueryAllVaultsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := lib.UnwrapSDKContext(c, types.ModuleName)

	var vaults []*types.QueryVaultResponse

	totalSharesStore := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.TotalSharesKeyPrefix))

	pageRes, err := query.Paginate(totalSharesStore, req.Pagination, func(key []byte, value []byte) error {
		vaultId, err := types.GetVaultIdFromStateKey(key)
		if err != nil {
			return err
		}

		vault, err := k.Vault(c, &types.QueryVaultRequest{
			Type:   vaultId.Type,
			Number: vaultId.Number,
		})
		if err != nil {
			return err
		}

		vaults = append(vaults, vault)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllVaultsResponse{
		Vaults:     vaults,
		Pagination: pageRes,
	}, nil
}
