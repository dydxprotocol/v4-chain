package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

var _ types.QueryServer = Keeper{}

// GetAuthenticators returns all authenticators for an account.
func (k Keeper) GetAuthenticators(
	ctx context.Context,
	request *types.GetAuthenticatorsRequest,
) (*types.GetAuthenticatorsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	acc, err := sdk.AccAddressFromBech32(request.Account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	authenticators, err := k.GetAuthenticatorDataForAccount(sdkCtx, acc)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.GetAuthenticatorsResponse{AccountAuthenticators: authenticators}, nil
}

// GetAuthenticator returns a specific authenticator for an account given its authenticator id.
func (k Keeper) GetAuthenticator(
	ctx context.Context,
	request *types.GetAuthenticatorRequest,
) (*types.GetAuthenticatorResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	acc, err := sdk.AccAddressFromBech32(request.Account)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	authenticator, err := k.GetSelectedAuthenticatorData(sdkCtx, acc, int(request.AuthenticatorId))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.GetAuthenticatorResponse{AccountAuthenticator: authenticator}, nil
}

// GetParams returns the parameters for the accountplus module.
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
