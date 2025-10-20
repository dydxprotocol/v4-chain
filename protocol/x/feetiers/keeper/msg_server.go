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

// SetMarketFeeDiscountParams sets or updates fee discount parameters for specific CLOB pairs
func (k msgServer) SetMarketFeeDiscountParams(
	goCtx context.Context,
	msg *types.MsgSetMarketFeeDiscountParams,
) (*types.MsgSetMarketFeeDiscountParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Process each market fee discount in the message
	for _, marketDiscount := range msg.Params {
		// Validate the fee discount parameters with the current block time
		if err := marketDiscount.Validate(ctx.BlockTime()); err != nil {
			return nil, errorsmod.Wrapf(
				err,
				"invalid market fee discount parameters for CLOB pair ID %d",
				marketDiscount.ClobPairId,
			)
		}

		// Set the market fee discount parameters
		if err := k.Keeper.SetPerMarketFeeDiscountParams(ctx, marketDiscount); err != nil {
			return nil, errorsmod.Wrapf(
				err,
				"failed to set market fee discount for CLOB pair ID %d",
				marketDiscount.ClobPairId,
			)
		}
	}

	return &types.MsgSetMarketFeeDiscountParamsResponse{}, nil
}

func (k msgServer) SetStakingTiers(
	goCtx context.Context,
	msg *types.MsgSetStakingTiers,
) (*types.MsgSetStakingTiersResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	if err := k.Keeper.SetStakingTiers(ctx, msg.StakingTiers); err != nil {
		return nil, err
	}

	return &types.MsgSetStakingTiersResponse{}, nil
}
