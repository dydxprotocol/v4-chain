package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

func (k msgServer) UpdatePerpetualParams(
	goCtx context.Context,
	msg *types.MsgUpdatePerpetualParams,
) (*types.MsgUpdatePerpetualParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := k.Keeper.ModifyPerpetual(
		ctx,
		msg.PerpetualParams.Id,
		msg.PerpetualParams.Ticker,
		msg.PerpetualParams.MarketId,
		msg.PerpetualParams.DefaultFundingPpm,
		msg.PerpetualParams.LiquidityTier,
	)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdatePerpetualParamsResponse{}, nil
}
