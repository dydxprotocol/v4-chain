package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
)

type msgServer struct {
	Keeper types.FeeTiersKeeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper types.FeeTiersKeeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) UpdatePerpetualFeeParams(
	goCtx context.Context,
	msg *types.MsgUpdatePerpetualFeeParams,
) (*types.MsgUpdatePerpetualFeeParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	if err := k.Keeper.SetPerpetualFeeParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePerpetualFeeParamsResponse{}, nil
}

// SetFeeDiscountCampaignParams sets or updates fee discount campaigns for specific CLOB pairs
func (k msgServer) SetFeeDiscountCampaignParams(
	goCtx context.Context,
	msg *types.MsgSetFeeDiscountCampaignParams,
) (*types.MsgSetFeeDiscountCampaignParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Process each fee discount campaign in the message
	for _, campaign := range msg.Params {
		// Validate the fee discount campaign parameters with the current block time
		if err := campaign.Validate(ctx.BlockTime()); err != nil {
			return nil, errorsmod.Wrapf(
				err,
				"invalid fee discount campaign parameters for CLOB pair ID %d",
				campaign.ClobPairId,
			)
		}

		// Set the fee discount campaign parameters
		if err := k.Keeper.SetFeeDiscountCampaignParams(ctx, campaign); err != nil {
			return nil, errorsmod.Wrapf(
				err,
				"failed to set fee discount campaign for CLOB pair ID %d",
				campaign.ClobPairId,
			)
		}
	}

	return &types.MsgSetFeeDiscountCampaignParamsResponse{}, nil
}
