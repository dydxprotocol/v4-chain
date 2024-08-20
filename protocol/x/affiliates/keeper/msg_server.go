package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

type msgServer struct {
	Keeper
}

// RegisterAffiliate implements types.MsgServer.
func (k msgServer) RegisterAffiliate(ctx context.Context, msg *types.MsgRegisterAffiliate) (*types.MsgRegisterAffiliateResponse, error) {
	return nil, nil
}

func (k msgServer) UpdateAffiliateTiers(ctx context.Context, msg *types.MsgUpdateAffiliateTiers) (*types.MsgUpdateAffiliateTiersResponse, error) {
	return nil, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
