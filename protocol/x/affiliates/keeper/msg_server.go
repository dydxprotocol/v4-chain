package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

type msgServer struct {
	Keeper
}

// RegisterAffiliate implements types.MsgServer.
func (k msgServer) RegisterAffiliate(ctx context.Context,
	msg *types.MsgRegisterAffiliate) (*types.MsgRegisterAffiliateResponse, error) {
	return nil, nil
}

func (k msgServer) UpdateAffiliateTiers(ctx context.Context,
	msg *types.MsgUpdateAffiliateTiers) (*types.MsgUpdateAffiliateTiersResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	unconditionalRevShareConfig, err := k.revShareKeeper.GetUnconditionalRevShareConfigParams(sdkCtx)
	if err != nil {
		return nil, err
	}
	marketMapperRevShareParams := k.revShareKeeper.GetMarketMapperRevenueShareParams(sdkCtx)

	if !k.revShareKeeper.ValidateRevShareSafety(*msg.Tiers, unconditionalRevShareConfig, marketMapperRevShareParams) {
		return nil, errorsmod.Wrapf(
			types.ErrRevShareSafetyViolation,
			"rev share safety violation",
		)
	}

	return nil, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
