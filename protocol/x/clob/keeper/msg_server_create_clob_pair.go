package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// CreateClobPair handles `MsgCreateClobPair`.
func (k msgServer) CreateClobPair(
	goCtx context.Context,
	msg *types.MsgCreateClobPair,
) (resp *types.MsgCreateClobPairResponse, err error) {
	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)

	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	perpetualId, err := msg.ClobPair.GetPerpetualId()
	if err != nil {
		return nil, err
	}

	// TODO(DEC-1535): update this when additional clob pair types are supported.
	if _, err := k.Keeper.CreatePerpetualClobPair(
		ctx,
		msg.ClobPair.Id,
		// `MsgCreateClobPair.ValidateBasic` ensures that `msg.ClobPair.Metadata` is `PerpetualClobMetadata`.
		perpetualId,
		satypes.BaseQuantums(msg.ClobPair.StepBaseQuantums),
		msg.ClobPair.QuantumConversionExponent,
		msg.ClobPair.SubticksPerTick,
		msg.ClobPair.Status,
	); err != nil {
		return nil, err
	}
	return &types.MsgCreateClobPairResponse{}, nil
}
