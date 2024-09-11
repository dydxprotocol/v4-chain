package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// UnlockShares unlocks an owner's shares.
func (k msgServer) UnlockShares(
	goCtx context.Context,
	msg *types.MsgUnlockShares,
) (*types.MsgUnlockSharesResponse, error) {
	// Check if authority is valid.
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if msg.OwnerAddress == "" {
		return nil, fmt.Errorf("owner address cannot be empty")
	}

	unlockedShares, err := k.Keeper.UnlockShares(ctx, msg.OwnerAddress)
	if err != nil {
		return nil, err
	}

	return &types.MsgUnlockSharesResponse{
		UnlockedShares: unlockedShares,
	}, nil
}
