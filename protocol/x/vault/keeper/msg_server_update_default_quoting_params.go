package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// UpdateDefaultQuotingParams updates the default quoting parameters of the vault module.
func (k msgServer) UpdateDefaultQuotingParams(
	goCtx context.Context,
	msg *types.MsgUpdateDefaultQuotingParams,
) (*types.MsgUpdateDefaultQuotingParamsResponse, error) {
	if !k.HasAuthority(msg.Authority) {
		return nil, errorsmod.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority %s",
			msg.Authority,
		)
	}

	ctx := lib.UnwrapSDKContext(goCtx, types.ModuleName)
	if err := k.SetDefaultQuotingParams(ctx, &msg.DefaultQuotingParams); err != nil {
		return nil, err
	}

	return &types.MsgUpdateDefaultQuotingParamsResponse{}, nil
}
