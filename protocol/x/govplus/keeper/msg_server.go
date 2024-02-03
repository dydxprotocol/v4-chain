package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/x/govplus/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SlashValidator(
	goCtx context.Context,
	msg *types.MsgSlashValidator,
) (*types.MsgSlashValidatorResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	consAddr, err := sdk.ConsAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrValidatorAddress
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	_, err = k.stakingKeeper.Slash(
		ctx,
		consAddr,
		int64(msg.InfractionHeight), // Casting from uint32
		sdk.TokensToConsensusPower(
			sdkmath.NewIntFromBigInt(msg.TokensAtInfractionHeight.BigInt()), sdk.DefaultPowerReduction),
		msg.SlashFactor,
	)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"error occurred when slashing validator",
			err,
		)
		return nil, err
	}
	return &types.MsgSlashValidatorResponse{}, nil
}
