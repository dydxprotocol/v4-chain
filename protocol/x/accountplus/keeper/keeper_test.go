package keeper_test

import (
	"encoding/hex"
	"encoding/json"
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
		[]types.Authenticator{
			authenticator.NewSignatureVerification(s.tApp.App.AccountKeeper),
			authenticator.NewMessageFilter(),
			authenticator.NewClobPairIdFilter(),
			authenticator.NewSubaccountFilter(),
			authenticator.NewAllOf(s.tApp.App.AuthenticatorManager),
			authenticator.NewAnyOf(s.tApp.App.AuthenticatorManager),
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

func (s *KeeperTestSuite) TestKeeper_Set_Get_GetAllAccountState() {
	ctx := s.Ctx

	accountState1 := types.AccountState{
		Address: "address1",
		TimestampNonceDetails: types.TimestampNonceDetails{
			TimestampNonces: []uint64{1, 2, 3},
			MaxEjectedNonce: 0,
		},
	}

	accountState2 := types.AccountState{
		Address: "address2",
		TimestampNonceDetails: types.TimestampNonceDetails{
			TimestampNonces: []uint64{1, 2, 3},
			MaxEjectedNonce: 0,
		},
	}

	// SetAccountState
	s.tApp.App.AccountPlusKeeper.SetAccountState(
		ctx,
		sdk.AccAddress([]byte(accountState1.Address)),
		accountState1,
	)

	// GetAccountState
	_, found := s.tApp.App.AccountPlusKeeper.GetAccountState(ctx, sdk.AccAddress([]byte(accountState1.Address)))
	s.Require().True(found, "Account state not found")

	// GetAllAccountStates
	accountStates, err := s.tApp.App.AccountPlusKeeper.GetAllAccountStates(ctx)
	s.Require().NoError(err, "Should not error when getting all account states")
	s.Require().Equal(len(accountStates), 1, "Incorrect number of AccountStates retrieved")
	s.Require().Equal(accountStates[0], accountState1, "Incorrect AccountState retrieved")

	// Add one more AccountState and check GetAllAccountStates
	s.tApp.App.AccountPlusKeeper.SetAccountState(
		ctx,
		sdk.AccAddress([]byte(accountState2.Address)),
		accountState2,
	)

	accountStates, err = s.tApp.App.AccountPlusKeeper.GetAllAccountStates(ctx)
	s.Require().NoError(err, "Should not error when getting all account states")
	s.Require().Equal(len(accountStates), 2, "Incorrect number of AccountStates retrieved")
	s.Require().Contains(accountStates, accountState1, "Retrieved AccountStates does not contain accountState1")
	s.Require().Contains(accountStates, accountState2, "Retrieved AccountStates does not contain accountState2")
}

func (s *KeeperTestSuite) TestKeeper_AddAuthenticator() {
	ctx := s.Ctx

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	// SignatureVerification should succeed
	id, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		priv.PubKey().Bytes(),
	)
	s.Require().NoError(err, "Should successfully add a SignatureVerification")
	s.Require().Equal(id, uint64(0), "Adding authenticator returning incorrect id")

	// MessageFilter should now fail because it doesn't require signature verification
	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"MessageFilter",
		[]byte("/cosmos.bank.v1beta1.MsgSend"),
	)
	s.Require().Error(err, "Should fail to add a MessageFilter as it doesn't require signature verification")
	s.Require().Contains(err.Error(), "authenticator tree does not require signature verification")

	// Invalid public key should fail
	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"SignatureVerification",
		[]byte("BrokenBytes"),
	)
	s.Require().Error(err, "Should have failed as OnAuthenticatorAdded fails")

	// After resetting authenticator manager, authenticator types should not be registered
	s.tApp.App.AuthenticatorManager.ResetAuthenticators()
	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"MessageFilter",
		[]byte("/cosmos.bank.v1beta1.MsgSend"),
	)
	s.Require().Error(err, "Authenticator not registered so should fail")
	s.Require().Contains(err.Error(), "authenticator type MessageFilter is not registered")
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

func (s *KeeperTestSuite) TestAddAuthenticator_SignatureVerificationRequired() {
	ctx := s.Ctx

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	tests := []struct {
		name                    string
		authenticatorType       string
		config                  []byte
		expectError             bool
		expectedErrorMsg        string
		expectSignatureRequired bool
	}{
		{
			name:                    "SignatureVerification authenticator requires signature",
			authenticatorType:       "SignatureVerification",
			config:                  priv.PubKey().Bytes(),
			expectError:             false,
			expectSignatureRequired: true,
		},
		{
			name:                    "MessageFilter authenticator does not require signature",
			authenticatorType:       "MessageFilter",
			config:                  []byte("/cosmos.bank.v1beta1.MsgSend"),
			expectError:             true,
			expectedErrorMsg:        "authenticator tree does not require signature verification",
			expectSignatureRequired: false,
		},
		{
			name:                    "ClobPairIdFilter authenticator does not require signature",
			authenticatorType:       "ClobPairIdFilter",
			config:                  []byte("0,1,2"),
			expectError:             true,
			expectedErrorMsg:        "authenticator tree does not require signature verification",
			expectSignatureRequired: false,
		},
		{
			name:                    "SubaccountFilter authenticator does not require signature",
			authenticatorType:       "SubaccountFilter",
			config:                  []byte("0,1"),
			expectError:             true,
			expectedErrorMsg:        "authenticator tree does not require signature verification",
			expectSignatureRequired: false,
		},
		{
			name:                    "AllOf with only filters fails",
			authenticatorType:       "AllOf",
			config:                  s.createAllOfConfig(false, false),
			expectError:             true,
			expectedErrorMsg:        "authenticator tree does not require signature verification",
			expectSignatureRequired: false,
		},
		{
			name:                    "AllOf with at least one SignatureVerification succeeds",
			authenticatorType:       "AllOf",
			config:                  s.createAllOfConfig(true, false),
			expectError:             false,
			expectSignatureRequired: true,
		},
		{
			name:                    "AllOf with all SignatureVerification succeeds",
			authenticatorType:       "AllOf",
			config:                  s.createAllOfConfig(true, true),
			expectError:             false,
			expectSignatureRequired: true,
		},
		{
			name:                    "AnyOf with only filters fails",
			authenticatorType:       "AnyOf",
			config:                  s.createAnyOfConfig(false, false),
			expectError:             true,
			expectedErrorMsg:        "authenticator tree does not require signature verification",
			expectSignatureRequired: false,
		},
		{
			name:                    "AnyOf with some SignatureVerification fails",
			authenticatorType:       "AnyOf",
			config:                  s.createAnyOfConfig(true, false),
			expectError:             true,
			expectedErrorMsg:        "authenticator tree does not require signature verification",
			expectSignatureRequired: false,
		},
		{
			name:                    "AnyOf with all SignatureVerification succeeds",
			authenticatorType:       "AnyOf",
			config:                  s.createAnyOfConfig(true, true),
			expectError:             false,
			expectSignatureRequired: true,
		},
		{
			name:                    "Nested AllOf with signature verification in inner succeeds",
			authenticatorType:       "AllOf",
			config:                  s.createNestedAllOfConfig(),
			expectError:             false,
			expectSignatureRequired: true,
		},
		{
			name:                    "Nested AnyOf with not all paths having signature fails",
			authenticatorType:       "AnyOf",
			config:                  s.createNestedAnyOfConfig(),
			expectError:             true,
			expectedErrorMsg:        "authenticator tree does not require signature verification",
			expectSignatureRequired: false,
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			id, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
				ctx,
				accAddress,
				tc.authenticatorType,
				tc.config,
			)

			if tc.expectError {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.expectedErrorMsg)
			} else {
				s.Require().NoError(err)
				s.Require().GreaterOrEqual(id, uint64(0))

				// Verify the authenticator was added
				authenticator, err := s.tApp.App.AccountPlusKeeper.GetSelectedAuthenticatorData(
					ctx,
					accAddress,
					id,
				)
				s.Require().NoError(err)
				s.Require().Equal(tc.authenticatorType, authenticator.Type)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAddAuthenticator_ComplexNestedStructures() {
	ctx := s.Ctx

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	// Test: AllOf containing AnyOf where all paths have signature verification
	innerAnyOfConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		},
		{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		},
	}
	innerAnyOfBytes, err := json.Marshal(innerAnyOfConfig)
	s.Require().NoError(err)

	allOfConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "AnyOf",
			Config: innerAnyOfBytes,
		},
		{
			Type:   "MessageFilter",
			Config: []byte("/dydxprotocol.clob.MsgPlaceOrder"),
		},
	}
	allOfBytes, err := json.Marshal(allOfConfig)
	s.Require().NoError(err)

	// This should succeed because the AnyOf has all paths with signature verification
	id, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"AllOf",
		allOfBytes,
	)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(id, uint64(0))

	// Test: AnyOf containing AllOf where at least one sub-authenticator has signature
	innerAllOfConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		},
		{
			Type:   "MessageFilter",
			Config: []byte("/dydxprotocol.clob.MsgPlaceOrder"),
		},
	}
	innerAllOfBytes, err := json.Marshal(innerAllOfConfig)
	s.Require().NoError(err)

	anyOfConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "AllOf",
			Config: innerAllOfBytes,
		},
		{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		},
	}
	anyOfBytes, err := json.Marshal(anyOfConfig)
	s.Require().NoError(err)

	// This should succeed because all paths in AnyOf lead to signature verification
	id, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"AnyOf",
		anyOfBytes,
	)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(id, uint64(0))
}

func (s *KeeperTestSuite) TestAddAuthenticator_EdgeCases() {
	ctx := s.Ctx

	// Set up account
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	// Test: Empty AllOf config (should fail during initialization)
	emptyConfig := []types.SubAuthenticatorInitData{}
	emptyConfigBytes, err := json.Marshal(emptyConfig)
	s.Require().NoError(err)

	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"AllOf",
		emptyConfigBytes,
	)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "no sub-authenticators provided")

	// Test: AllOf with only one authenticator (should fail)
	singleConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		},
	}
	singleConfigBytes, err := json.Marshal(singleConfig)
	s.Require().NoError(err)

	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"AllOf",
		singleConfigBytes,
	)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "no sub-authenticators provided")

	// Test: Invalid authenticator type
	_, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		ctx,
		accAddress,
		"NonExistentAuthenticator",
		[]byte("config"),
	)
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "authenticator type NonExistentAuthenticator is not registered")
}

// Helper functions to create various authenticator configurations
func (s *KeeperTestSuite) createAllOfConfig(firstHasSig, secondHasSig bool) []byte {
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}

	config := []types.SubAuthenticatorInitData{}

	if firstHasSig {
		config = append(config, types.SubAuthenticatorInitData{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		})
	} else {
		config = append(config, types.SubAuthenticatorInitData{
			Type:   "MessageFilter",
			Config: []byte("/cosmos.bank.v1beta1.MsgSend"),
		})
	}

	if secondHasSig {
		config = append(config, types.SubAuthenticatorInitData{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		})
	} else {
		config = append(config, types.SubAuthenticatorInitData{
			Type:   "ClobPairIdFilter",
			Config: []byte("0,1"),
		})
	}

	configBytes, err := json.Marshal(config)
	s.Require().NoError(err)
	return configBytes
}

func (s *KeeperTestSuite) createAnyOfConfig(firstHasSig, secondHasSig bool) []byte {
	// Same structure as AllOf but different logic applies
	return s.createAllOfConfig(firstHasSig, secondHasSig)
}

func (s *KeeperTestSuite) createNestedAllOfConfig() []byte {
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}

	// Create inner AllOf with SignatureVerification
	innerConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		},
		{
			Type:   "MessageFilter",
			Config: []byte("/dydxprotocol.clob.MsgPlaceOrder"),
		},
	}
	innerConfigBytes, err := json.Marshal(innerConfig)
	s.Require().NoError(err)

	// Create outer AllOf
	outerConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "AllOf",
			Config: innerConfigBytes,
		},
		{
			Type:   "SubaccountFilter",
			Config: []byte("0"),
		},
	}
	outerConfigBytes, err := json.Marshal(outerConfig)
	s.Require().NoError(err)
	return outerConfigBytes
}

func (s *KeeperTestSuite) createNestedAnyOfConfig() []byte {
	key := "6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159"
	bz, _ := hex.DecodeString(key)
	priv := &secp256k1.PrivKey{Key: bz}

	// Create inner AnyOf with one path having SignatureVerification
	innerConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "SignatureVerification",
			Config: priv.PubKey().Bytes(),
		},
		{
			Type:   "MessageFilter",
			Config: []byte("/dydxprotocol.clob.MsgPlaceOrder"),
		},
	}
	innerConfigBytes, err := json.Marshal(innerConfig)
	s.Require().NoError(err)

	// Create outer AnyOf - this will fail because inner AnyOf doesn't have all paths with signature
	outerConfig := []types.SubAuthenticatorInitData{
		{
			Type:   "AnyOf",
			Config: innerConfigBytes,
		},
		{
			Type:   "SubaccountFilter",
			Config: []byte("0"),
		},
	}
	outerConfigBytes, err := json.Marshal(outerConfig)
	s.Require().NoError(err)
	return outerConfigBytes
}
