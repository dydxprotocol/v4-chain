package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
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

// AddAuthenticator allows the addition of various types of authenticators to an account.
// This method serves as a versatile function for adding diverse authenticator types
// to an account, making it highly adaptable for different use cases.
func (m msgServer) AddAuthenticator(
	goCtx context.Context,
	msg *types.MsgAddAuthenticator,
) (*types.MsgAddAuthenticatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isSmartAccountActive := m.GetIsSmartAccountActive(ctx)
	if !isSmartAccountActive {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "smart account authentication flow is not active")
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// Finally, add the authenticator to the store.
	_, err = m.Keeper.AddAuthenticator(ctx, sender, msg.AuthenticatorType, msg.Data)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddAuthenticatorResponse{
		Success: true,
	}, nil
}

// RemoveAuthenticator removes an authenticator from the store. The message specifies a sender address and an index.
func (m msgServer) RemoveAuthenticator(
	goCtx context.Context,
	msg *types.MsgRemoveAuthenticator,
) (*types.MsgRemoveAuthenticatorResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isSmartAccountActive := m.GetIsSmartAccountActive(ctx)
	if !isSmartAccountActive {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "smart account authentication flow is not active")
	}

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid sender address")
	}

	// At this point, we assume that verification has occurred on the account, and we
	// proceed to remove the authenticator from the store.
	err = m.Keeper.RemoveAuthenticator(ctx, sender, msg.Id)
	if err != nil {
		return nil, err
	}

	return &types.MsgRemoveAuthenticatorResponse{
		Success: true,
	}, nil
}

// SetActiveState sets the active state of the smart account authentication flow.
func (m msgServer) SetActiveState(
	goCtx context.Context,
	msg *types.MsgSetActiveState,
) (*types.MsgSetActiveStateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !m.Keeper.HasAuthority(msg.GetAuthority()) {
		return nil, errorsmod.Wrapf(
			sdkerrors.ErrUnauthorized,
			"%v is not recognized as a valid authority for setting smart account active state",
			msg.GetAuthority(),
		)
	}

	// Set the active state of the authenticator
	m.Keeper.SetActiveState(ctx, msg.Active)

	return &types.MsgSetActiveStateResponse{}, nil
}
