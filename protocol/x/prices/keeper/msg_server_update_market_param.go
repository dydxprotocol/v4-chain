package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func (k msgServer) UpdateMarketParam(
	goCtx context.Context,
	msg *types.MsgUpdateMarketParam,
) (*types.MsgUpdateMarketParamResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, err := k.Keeper.ModifyMarketParam(ctx, msg.MarketParam); err != nil {
		return nil, err
	}

	return &types.MsgUpdateMarketParamResponse{}, nil
}
