package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CreateClobPair handles `MsgCreateClobPair`.
func (k msgServer) CreateClobPair(
	goCtx context.Context,
	msg *types.MsgCreateClobPair,
) (*types.MsgCreateClobPairResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO(DEC-1535): update this when additional clob pair types are supported.
	if _, err := k.Keeper.CreatePerpetualClobPair(
		ctx,
		msg.ClobPair.Id,
		// `MsgCreateClobPair.ValidateBasic` ensures that `msg.ClobPair.Metadata` is `PerpetualClobMetadata`.
		msg.ClobPair.MustGetPerpetualId(),
		satypes.BaseQuantums(msg.ClobPair.StepBaseQuantums),
		msg.ClobPair.QuantumConversionExponent,
		msg.ClobPair.SubticksPerTick,
		msg.ClobPair.Status,
	); err != nil {
		return nil, err
	}
	return &types.MsgCreateClobPairResponse{}, nil
}
