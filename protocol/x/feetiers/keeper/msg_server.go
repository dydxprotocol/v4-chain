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

// SetFeeHolidayParams sets or updates fee holidays for specific CLOB pairs
func (k msgServer) SetFeeHolidayParams(
	goCtx context.Context,
	msg *types.MsgSetFeeHolidayParams,
) (*types.MsgSetFeeHolidayParamsResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	// Process each fee holiday in the message
	for _, feeHoliday := range msg.Params {
		// Validate the fee holiday parameters
		if err := feeHoliday.Validate(ctx.BlockTime()); err != nil {
			return nil, errorsmod.Wrapf(
				err,
				"invalid fee holiday parameters for CLOB pair ID %d",
				feeHoliday.ClobPairId,
			)
		}

		// Set the fee holiday parameters
		if err := k.Keeper.SetFeeHolidayParams(ctx, feeHoliday); err != nil {
			return nil, errorsmod.Wrapf(
				err,
				"failed to set fee holiday for CLOB pair ID %d",
				feeHoliday.ClobPairId,
			)
		}
	}

	return &types.MsgSetFeeHolidayParamsResponse{}, nil
}
