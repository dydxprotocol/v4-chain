package keeper

import (
	"context"

	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) AffiliateInfo(c context.Context,
	req *types.AffiliateInfoRequest) (*types.AffiliateInfoResponse, error) {
	return nil, nil
}

func (k Keeper) ReferredBy(ctx context.Context,
	req *types.ReferredByRequest) (*types.ReferredByResponse, error) {
	return nil, nil
}

func (k Keeper) AllAffiliateTiers(c context.Context,
	req *types.AllAffiliateTiersRequest) (*types.AllAffiliateTiersResponse, error) {
	return nil, nil
}
