package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k msgServer) SetClobPairStatus(
	goCtx context.Context,
	msg *types.MsgSetClobPairStatus,
) (*types.MsgSetClobPairStatusResponse, error) {
	if !k.Keeper.HasAuthority(msg.GetAuthority()) {
		return nil, sdkerrors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.Keeper.SetClobPairStatus(
		ctx,
		types.ClobPairId(msg.GetClobPairId()),
		types.ClobPair_Status(msg.GetClobPairStatus()),
	); err != nil {
		return nil, err
	}

	return &types.MsgSetClobPairStatusResponse{}, nil
}
