package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/revshare/types"
)

func (k msgServer) SetOrderRouterRevShares(
	goCtx context.Context,
	msg *types.MsgSetOrderRouterRevShares,
) (*types.MsgSetOrderRouterRevSharesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if sender is authorized to set revenue share
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	for _, orderRouterRevShare := range msg.OrderRouterRevShares {
		// Maximum fee share is 500_000 ppm
		if orderRouterRevShare.SharePpm > lib.OneHundredThousand*5 {
			return nil, errorsmod.Wrapf(
				types.ErrInvalidRevenueSharePpm,
				"rev share safety violation: rev shares greater than or equal to allowed amount",
			)
		}
		if err := k.Keeper.SetOrderRouterRevShares(
			ctx, orderRouterRevShare.Address, orderRouterRevShare.SharePpm,
		); err != nil {
			return nil, err
		}
	}

	return &types.MsgSetOrderRouterRevSharesResponse{}, nil
}
