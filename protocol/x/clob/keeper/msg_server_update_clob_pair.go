package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

func (k msgServer) UpdateClobPair(
	goCtx context.Context,
	msg *types.MsgUpdateClobPair,
) (resp *types.MsgUpdateClobPairResponse, err error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	if err := k.Keeper.UpdateClobPair(
		ctx,
		msg.GetClobPair(),
	); err != nil {
		return nil, err
	}

	return &types.MsgUpdateClobPairResponse{}, nil
}
