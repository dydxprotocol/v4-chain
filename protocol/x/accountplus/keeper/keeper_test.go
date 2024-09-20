package keeper_test

import (
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"

	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/testutils"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

// TODO: add explicit unit tests for Get and Set funcs
// https://linear.app/dydx/issue/OTE-633/add-explicit-unit-tests-for-get-and-set-methods-for-accountplus-keeper

type KeeperTestSuite struct {
	suite.Suite

	tApp *testapp.TestApp
	Ctx  sdk.Context
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.tApp = testapp.NewTestAppBuilder(s.T()).Build()
	s.Ctx = s.tApp.InitChain()

	s.tApp.App.AuthenticatorManager.ResetAuthenticators()
	s.tApp.App.AuthenticatorManager.InitializeAuthenticators(
		[]authenticator.Authenticator{
			authenticator.SignatureVerification{},
			authenticator.MessageFilter{},
			testutils.TestingAuthenticator{
				Approve:        testutils.Always,
				GasConsumption: 10,
				Confirm:        testutils.Always,
			},
		},
	)
	s.tApp.App.AccountPlusKeeper.SetParams(
		s.Ctx,
		types.Params{
			IsSmartAccountActive: true,
		},
	)
}

func (s *KeeperTestSuite) TestKeeper_AddAuthenticator() {
	ctx := s.Ctx

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	id, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		priv.PubKey().Bytes(),
	)
	s.Require().NoError(err, "Should successfully add a SignatureVerification")
	s.Require().Equal(id, uint64(0), "Adding authenticator returning incorrect id")

	id, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"MessageFilter",
		[]byte("/cosmos.bank.v1beta1.MsgSend"),
	)
	s.Require().NoError(err, "Should successfully add a MessageFilter")
	s.Require().Equal(id, uint64(1), "Adding authenticator returning incorrect id")

	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		[]byte("BrokenBytes"),
	)
	s.Require().Error(err, "Should have failed as OnAuthenticatorAdded fails")

	s.tApp.App.AuthenticatorManager.ResetAuthenticators()
	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"MessageFilter",
		[]byte("/cosmos.bank.v1beta1.MsgSend"),
	)
	s.Require().Error(err, "Authenticator not registered so should fail")
}

func (s *KeeperTestSuite) TestKeeper_GetAndSetAuthenticatorId() {
	ctx := s.Ctx

	authenticatorId := s.tApp.App.AccountPlusKeeper.InitializeOrGetNextAuthenticatorId(ctx)
	s.Require().Equal(uint64(0), authenticatorId, "Initialize/Get authenticator id returned incorrect id")

	authenticatorId = s.tApp.App.AccountPlusKeeper.InitializeOrGetNextAuthenticatorId(ctx)
	s.Require().Equal(uint64(0), authenticatorId, "Initialize/Get authenticator id returned incorrect id")

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	_, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		priv.PubKey().Bytes(),
	)
	s.Require().NoError(err, "Should successfully add a SignatureVerification")

	authenticatorId = s.tApp.App.AccountPlusKeeper.InitializeOrGetNextAuthenticatorId(ctx)
	s.Require().Equal(authenticatorId, uint64(1), "Initialize/Get authenticator id returned incorrect id")

	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		priv.PubKey().Bytes(),
	)
	s.Require().NoError(err, "Should successfully add a MessageFilter")

	authenticatorId = s.tApp.App.AccountPlusKeeper.InitializeOrGetNextAuthenticatorId(ctx)
	s.Require().Equal(authenticatorId, uint64(2), "Initialize/Get authenticator id returned incorrect id")
}

func (s *KeeperTestSuite) TestKeeper_GetAuthenticatorDataForAccount() {
	ctx := s.Ctx

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	_, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		priv.PubKey().Bytes(),
	)
	s.Require().NoError(err, "Should successfully add a SignatureVerification")

	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		priv.PubKey().Bytes(),
	)
	s.Require().NoError(err, "Should successfully add a MessageFilter")

	authenticators, err := s.tApp.App.AccountPlusKeeper.GetAuthenticatorDataForAccount(ctx, accAddress)
	s.Require().NoError(err)
	s.Require().Equal(len(authenticators), 2, "Getting authenticators returning incorrect data")
}
