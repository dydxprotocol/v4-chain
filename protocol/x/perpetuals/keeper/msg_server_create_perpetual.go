package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
)

func (k msgServer) CreatePerpetual(
	goCtx context.Context,
	msg *types.MsgCreatePerpetual,
) (*types.MsgCreatePerpetualResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := k.Keeper.CreatePerpetual(
		ctx,
		msg.Params.Id,
		msg.Params.Ticker,
		msg.Params.MarketId,
		msg.Params.AtomicResolution,
		msg.Params.DefaultFundingPpm,
		msg.Params.LiquidityTier,
	)
	if err != nil {
		return &types.MsgCreatePerpetualResponse{}, err
	}

	return &types.MsgCreatePerpetualResponse{}, nil
}
