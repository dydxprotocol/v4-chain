package keeper

import (
	"context"
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

type msgServer struct {
	Keeper
}

// RegisterAffiliate implements types.MsgServer.
// This is only valid if a referee signs the message
// since the referee field is annotated with cosmos.msg.v1.signer
// in protos. This ensures that only referee is returned
// as a signer when GetSigners is called for authentication.
// For example, if Alice is the referee and Bob is the affiliate,
// then only Alice can register Bob as an affiliate. Any
// other signer that sends this message will be rejected.
func (k msgServer) RegisterAffiliate(ctx context.Context,
	msg *types.MsgRegisterAffiliate) (*types.MsgRegisterAffiliateResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := k.Keeper.RegisterAffiliate(sdkCtx, msg.Referee, msg.Affiliate)
	if err != nil {
		return nil, err
	}
	return &types.MsgRegisterAffiliateResponse{}, nil
}

func (k msgServer) UpdateAffiliateTiers(ctx context.Context,
	msg *types.MsgUpdateAffiliateTiers) (*types.MsgUpdateAffiliateTiersResponse, error) {
	if !k.Keeper.HasAuthority(msg.Authority) {
		return nil, errors.New("invalid authority")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	unconditionalRevShareConfig, err := k.revShareKeeper.GetUnconditionalRevShareConfigParams(sdkCtx)
	if err != nil {
		return nil, err
	}
	marketMapperRevShareParams := k.revShareKeeper.GetMarketMapperRevenueShareParams(sdkCtx)

	if !k.revShareKeeper.ValidateRevShareSafety(msg.Tiers, unconditionalRevShareConfig, marketMapperRevShareParams) {
		return nil, errorsmod.Wrapf(
			types.ErrRevShareSafetyViolation,
			"rev share safety violation",
		)
	}

	err = k.Keeper.UpdateAffiliateTiers(sdkCtx, msg.Tiers)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdateAffiliateTiersResponse{}, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
