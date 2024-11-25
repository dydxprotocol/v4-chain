package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
)

func (k msgServer) UpgradeIsolatedPerpetualToCross(
	goCtx context.Context,
	msg *types.MsgUpgradeIsolatedPerpetualToCross,
) (*types.MsgUpgradeIsolatedPerpetualToCrossResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	err := k.Keeper.UpgradeIsolatedPerpetualToCross(
		ctx,
		msg.PerpetualId,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpgradeIsolatedPerpetualToCrossResponse{}, nil
}
