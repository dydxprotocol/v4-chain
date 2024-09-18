package keeper_test

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

func (s *KeeperTestSuite) TestMsgServer_AddAuthenticator() {
	msgServer := keeper.NewMsgServerImpl(s.tApp.App.AccountPlusKeeper)
	ctx := s.Ctx

	// Ensure the SigVerificationAuthenticator type is registered
	s.Require().True(
		s.tApp.App.AuthenticatorManager.IsAuthenticatorTypeRegistered(
			authenticator.SignatureVerification{}.Type(),
		),
	)

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	// Create a test message
	msg := &types.MsgAddAuthenticator{
		Sender:            accAddress.String(),
		AuthenticatorType: authenticator.SignatureVerification{}.Type(),
		Data:              priv.PubKey().Bytes(),
	}

	resp, err := msgServer.AddAuthenticator(ctx, msg)
	s.Require().NoError(err)
	s.Require().True(resp.Success)
}

func (s *KeeperTestSuite) TestMsgServer_AddAuthenticatorFail() {
	msgServer := keeper.NewMsgServerImpl(s.tApp.App.AccountPlusKeeper)
	ctx := s.Ctx

	// Ensure the SigVerificationAuthenticator type is registered
	s.Require().True(
		s.tApp.App.AuthenticatorManager.IsAuthenticatorTypeRegistered(
			authenticator.SignatureVerification{}.Type(),
		),
	)

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	// Create a test message
	msg := &types.MsgAddAuthenticator{
		Sender:            accAddress.String(),
		AuthenticatorType: authenticator.SignatureVerification{}.Type(),
		Data:              priv.PubKey().Bytes(),
	}

	msg.AuthenticatorType = "PassKeyAuthenticator"
	_, err := msgServer.AddAuthenticator(ctx, msg)
	s.Require().Error(err)
}

func (s *KeeperTestSuite) TestMsgServer_RemoveAuthenticator() {
	msgServer := keeper.NewMsgServerImpl(s.tApp.App.AccountPlusKeeper)
	ctx := s.Ctx

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	// Create a test message
	addMsg := &types.MsgAddAuthenticator{
		Sender:            accAddress.String(),
		AuthenticatorType: authenticator.SignatureVerification{}.Type(),
		Data:              priv.PubKey().Bytes(),
	}
	_, err := msgServer.AddAuthenticator(ctx, addMsg)
	s.Require().NoError(err)

	// Now attempt to remove it
	removeMsg := &types.MsgRemoveAuthenticator{
		Sender: accAddress.String(),
		Id:     0,
	}

	resp, err := msgServer.RemoveAuthenticator(ctx, removeMsg)
	s.Require().NoError(err)
	s.Require().True(resp.Success)
}

func (s *KeeperTestSuite) TestMsgServer_SetActiveState() {
	ak := s.tApp.App.AccountPlusKeeper
	msgServer := keeper.NewMsgServerImpl(ak)
	ctx := s.Ctx

	// Set up account
	key := "0dd4d1506e18a5712080708c338eb51ecf2afdceae01e8162e890b126ac190fe"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	unauthorizedAccAddress := sdk.AccAddress(priv.PubKey().Address())

	// activated by default
	initialParams := s.tApp.App.AccountPlusKeeper.GetParams(ctx)
	s.Require().True(initialParams.IsSmartAccountActive)

	// deactivate by unauthorized account
	_, err := msgServer.SetActiveState(
		ctx,
		&types.MsgSetActiveState{
			Authority: unauthorizedAccAddress.String(),
			Active:    false,
		})

	s.Require().Error(err)
	s.Require().Equal(
		err.Error(),
		"dydx1jns2dl6u55vy72g7r76l6cse7kmlj4me87xj2j is not recognized as a valid authority "+
			"for setting smart account active state: unauthorized",
	)

	// deactivate
	_, err = msgServer.SetActiveState(
		ctx,
		&types.MsgSetActiveState{
			Authority: lib.GovModuleAddress.String(),
			Active:    false,
		})

	s.Require().NoError(err)

	// active state should be false
	params := ak.GetParams(ctx)
	s.Require().False(params.IsSmartAccountActive)

	// reactivate by gov
	_, err = msgServer.SetActiveState(
		ctx,
		&types.MsgSetActiveState{
			Authority: lib.GovModuleAddress.String(),
			Active:    true,
		})
	s.Require().NoError(err)

	// active state should be true
	params = ak.GetParams(ctx)
	s.Require().True(params.IsSmartAccountActive)
}

func (s *KeeperTestSuite) TestMsgServer_SmartAccountsNotActive() {
	msgServer := keeper.NewMsgServerImpl(s.tApp.App.AccountPlusKeeper)
	ctx := s.Ctx

	s.tApp.App.AccountPlusKeeper.SetParams(s.Ctx, types.Params{IsSmartAccountActive: false})

	// Create a test message
	msg := &types.MsgAddAuthenticator{
		Sender:            "",
		AuthenticatorType: authenticator.SignatureVerification{}.Type(),
		Data:              []byte(""),
	}

	_, err := msgServer.AddAuthenticator(ctx, msg)
	s.Require().Error(err)
	s.Require().Equal(err.Error(), "smart account authentication flow is not active: unauthorized")

	removeMsg := &types.MsgRemoveAuthenticator{
		Sender: "",
		Id:     1,
	}

	_, err = msgServer.RemoveAuthenticator(ctx, removeMsg)
	s.Require().Error(err)
	s.Require().Equal(err.Error(), "smart account authentication flow is not active: unauthorized")
}
