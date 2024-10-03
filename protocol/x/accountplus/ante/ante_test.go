package ante_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"testing"

	storetypes "cosmossdk.io/store/types"
	tmtypes "github.com/cometbft/cometbft/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"

	"github.com/stretchr/testify/suite"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/ante"
)

type AuthenticatorAnteSuite struct {
	suite.Suite

	tApp                   *testapp.TestApp
	Ctx                    sdk.Context
	EncodingConfig         app.EncodingConfig
	AuthenticatorDecorator ante.AuthenticatorDecorator
	TestKeys               []string
	TestAccAddress         []sdk.AccAddress
	TestPrivKeys           []*secp256k1.PrivKey
	HomeDir                string
}

func TestAuthenticatorAnteSuite(t *testing.T) {
	suite.Run(t, new(AuthenticatorAnteSuite))
}

func (s *AuthenticatorAnteSuite) SetupTest() {
	// Test data for authenticator signature verification
	TestKeys := []string{
		"6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159",
		"0dd4d1506e18a5712080708c338eb51ecf2afdceae01e8162e890b126ac190fe",
		"49006a359803f0602a7ec521df88bf5527579da79112bb71f285dd3e7d438033",
	}

	s.HomeDir = fmt.Sprintf("%d", rand.Int())

	// Set up test accounts
	accounts := make([]sdk.AccountI, 0)
	for _, key := range TestKeys {
		bz, _ := hex.DecodeString(key)
		priv := &secp256k1.PrivKey{Key: bz}

		// Add the test private keys to an array for later use
		s.TestPrivKeys = append(s.TestPrivKeys, priv)

		// Generate an account address from the public key
		accAddress := sdk.AccAddress(priv.PubKey().Address())
		accounts = append(
			accounts,
			authtypes.NewBaseAccount(accAddress, priv.PubKey(), 0, 0),
		)

		// Add the test accounts' addresses to an array for later use
		s.TestAccAddress = append(s.TestAccAddress, accAddress)
	}

	s.tApp = testapp.NewTestAppBuilder(s.T()).WithGenesisDocFn(func() (genesis tmtypes.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *authtypes.GenesisState) {
				for _, acct := range accounts {
					genesisState.Accounts = append(genesisState.Accounts, codectypes.UnsafePackAny(acct))
				}
			},
		)
		return genesis
	}).Build()
	s.Ctx = s.tApp.InitChain()

	s.EncodingConfig = app.GetEncodingConfig()
	s.AuthenticatorDecorator = ante.NewAuthenticatorDecorator(
		s.tApp.App.AppCodec(),
		&s.tApp.App.AccountPlusKeeper,
		s.tApp.App.AccountKeeper,
		s.EncodingConfig.TxConfig.SignModeHandler(),
	)
	s.Ctx = s.Ctx.WithGasMeter(storetypes.NewGasMeter(1_000_000))
}

func (s *AuthenticatorAnteSuite) TearDownTest() {
	os.RemoveAll(s.HomeDir)
}

// TestSignatureVerificationNoAuthenticatorInStore test a non-smart account signature verification
// with no authenticator in the store
func (s *AuthenticatorAnteSuite) TestSignatureVerificationNoAuthenticatorInStore() {
	bech32Prefix := config.Bech32PrefixAccAddr
	coins := sdk.Coins{sdk.NewInt64Coin(constants.TestNativeTokenDenom, 2500)}

	// Create a test messages for signing
	testMsg1 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[0]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}
	testMsg2 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}
	feeCoins := constants.TestFeeCoins_5Cents

	tx, err := testtx.GenTx(
		s.Ctx,
		s.EncodingConfig.TxConfig,
		[]sdk.Msg{
			testMsg1,
			testMsg2,
		},
		feeCoins,
		300000,
		s.Ctx.ChainID(),
		[]uint64{6, 6},
		[]uint64{0, 0},
		[]cryptotypes.PrivKey{
			s.TestPrivKeys[0],
			s.TestPrivKeys[1],
		},
		[]cryptotypes.PrivKey{
			s.TestPrivKeys[0],
			s.TestPrivKeys[1],
		},
		[]uint64{0, 0},
	)
	s.Require().NoError(err)

	anteHandler := sdk.ChainAnteDecorators(s.AuthenticatorDecorator)
	_, err = anteHandler(s.Ctx, tx, false)

	s.Require().Error(err, "Expected error when no authenticator is in the store")
}

// TestSignatureVerificationWithAuthenticatorInStore test a non-smart account signature verification
// with a single authenticator in the store
func (s *AuthenticatorAnteSuite) TestSignatureVerificationWithAuthenticatorInStore() {
	bech32Prefix := config.Bech32PrefixAccAddr
	coins := sdk.Coins{sdk.NewInt64Coin(constants.TestNativeTokenDenom, 2500)}

	// Create a test messages for signing
	testMsg1 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[0]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}
	testMsg2 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}
	feeCoins := constants.TestFeeCoins_5Cents

	id, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		s.Ctx,
		s.TestAccAddress[0],
		"SignatureVerification",
		s.TestPrivKeys[0].PubKey().Bytes(),
	)
	s.Require().NoError(err)
	s.Require().Equal(id, uint64(0), "Adding authenticator returning incorrect id")

	id, err = s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		s.Ctx,
		s.TestAccAddress[1],
		"SignatureVerification",
		s.TestPrivKeys[1].PubKey().Bytes(),
	)
	s.Require().NoError(err)
	s.Require().Equal(id, uint64(1), "Adding authenticator returning incorrect id")

	s.tApp.App.AccountPlusKeeper.SetActiveState(s.Ctx, true)
	s.Require().True(
		s.tApp.App.AccountPlusKeeper.GetIsSmartAccountActive(s.Ctx),
		"Expected smart account to be active",
	)

	tx, err := testtx.GenTx(
		s.Ctx,
		s.EncodingConfig.TxConfig,
		[]sdk.Msg{
			testMsg1,
			testMsg2,
		},
		feeCoins,
		300000,
		s.Ctx.ChainID(),
		[]uint64{5, 6},
		[]uint64{0, 0},
		[]cryptotypes.PrivKey{
			s.TestPrivKeys[0],
			s.TestPrivKeys[1],
		},
		[]cryptotypes.PrivKey{
			s.TestPrivKeys[0],
			s.TestPrivKeys[1],
		},
		[]uint64{0, 1},
	)
	s.Require().NoError(err)

	anteHandler := sdk.ChainAnteDecorators(s.AuthenticatorDecorator)
	_, err = anteHandler(s.Ctx, tx, false)

	s.Require().NoError(err)
}

// TestFeePayerGasComsumption tests that the fee payer only gets charged gas for the transaction once.
func (s *AuthenticatorAnteSuite) TestFeePayerGasComsumption() {
	bech32Prefix := config.Bech32PrefixAccAddr
	coins := sdk.Coins{sdk.NewInt64Coin(constants.TestNativeTokenDenom, 2500)}
	feeCoins := constants.TestFeeCoins_5Cents

	specifiedGasLimit := uint64(300_000)

	// Create two messages to ensure that the fee payer code path is reached twice
	testMsg1 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[0]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}

	testMsg2 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[0]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}

	// Add a signature verification authenticator to the account
	sigId, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		s.Ctx,
		s.TestAccAddress[0],
		"SignatureVerification",
		s.TestPrivKeys[1].PubKey().Bytes(),
	)
	s.Require().NoError(err)
	s.Require().Equal(sigId, uint64(0), "Adding authenticator returning incorrect id")

	s.tApp.App.AccountPlusKeeper.SetActiveState(s.Ctx, true)
	s.Require().True(
		s.tApp.App.AccountPlusKeeper.GetIsSmartAccountActive(s.Ctx),
		"Expected smart account to be active",
	)

	tx, err := testtx.GenTx(
		s.Ctx,
		s.EncodingConfig.TxConfig,
		[]sdk.Msg{
			testMsg1,
			testMsg2,
		},
		feeCoins,
		specifiedGasLimit,
		s.Ctx.ChainID(),
		[]uint64{5, 5},
		[]uint64{0, 0},
		[]cryptotypes.PrivKey{
			s.TestPrivKeys[0],
		},
		[]cryptotypes.PrivKey{
			s.TestPrivKeys[1],
		},
		[]uint64{sigId, sigId},
	)
	s.Require().NoError(err)

	anteHandler := sdk.ChainAnteDecorators(s.AuthenticatorDecorator)
	_, err = anteHandler(s.Ctx, tx, false)
	s.Require().NoError(err)
}

func (s *AuthenticatorAnteSuite) TestSpecificAuthenticator() {
	bech32Prefix := config.Bech32PrefixAccAddr
	coins := sdk.Coins{sdk.NewInt64Coin(constants.TestNativeTokenDenom, 2500)}
	feeCoins := constants.TestFeeCoins_5Cents

	// Create a test messages for signing
	testMsg1 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}

	sig1Id, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		s.Ctx,
		s.TestAccAddress[1],
		"SignatureVerification",
		s.TestPrivKeys[0].PubKey().Bytes(),
	)
	s.Require().NoError(err)
	s.Require().Equal(sig1Id, uint64(0), "Adding authenticator returning incorrect id")

	sig2Id, err := s.tApp.App.AccountPlusKeeper.AddAuthenticator(
		s.Ctx,
		s.TestAccAddress[1],
		"SignatureVerification",
		s.TestPrivKeys[1].PubKey().Bytes(),
	)
	s.Require().NoError(err)
	s.Require().Equal(sig2Id, uint64(1), "Adding authenticator returning incorrect id")

	s.tApp.App.AccountPlusKeeper.SetActiveState(s.Ctx, true)
	s.Require().True(
		s.tApp.App.AccountPlusKeeper.GetIsSmartAccountActive(s.Ctx),
		"Expected smart account to be active",
	)

	testCases := map[string]struct {
		name                  string
		senderKey             cryptotypes.PrivKey
		signKey               cryptotypes.PrivKey
		selectedAuthenticator []uint64
		shouldPass            bool
	}{
		"Correct authenticator 0": {
			senderKey:             s.TestPrivKeys[0],
			signKey:               s.TestPrivKeys[0],
			selectedAuthenticator: []uint64{sig1Id},
			shouldPass:            true,
		},
		"Correct authenticator 1": {
			senderKey:             s.TestPrivKeys[0],
			signKey:               s.TestPrivKeys[1],
			selectedAuthenticator: []uint64{sig2Id},
			shouldPass:            true,
		},
		"Incorrect authenticator 0": {
			senderKey:             s.TestPrivKeys[0],
			signKey:               s.TestPrivKeys[0],
			selectedAuthenticator: []uint64{sig2Id},
			shouldPass:            false,
		},
		"Incorrect authenticator 1": {
			senderKey:             s.TestPrivKeys[0],
			signKey:               s.TestPrivKeys[1],
			selectedAuthenticator: []uint64{sig1Id},
			shouldPass:            false,
		},
		"Not Specified for 0": {
			senderKey:             s.TestPrivKeys[0],
			signKey:               s.TestPrivKeys[0],
			selectedAuthenticator: []uint64{},
			shouldPass:            false,
		},
		"Not Specified for 1": {
			senderKey:             s.TestPrivKeys[0],
			signKey:               s.TestPrivKeys[1],
			selectedAuthenticator: []uint64{},
			shouldPass:            false,
		},
		"Bad selection": {
			senderKey:             s.TestPrivKeys[0],
			signKey:               s.TestPrivKeys[0],
			selectedAuthenticator: []uint64{3},
			shouldPass:            false,
		},
	}

	for name, tc := range testCases {
		s.Run(name, func() {
			tx, err := testtx.GenTx(
				s.Ctx,
				s.EncodingConfig.TxConfig,
				[]sdk.Msg{
					testMsg1,
				},
				feeCoins,
				300000,
				s.Ctx.ChainID(),
				[]uint64{6},
				[]uint64{0},
				[]cryptotypes.PrivKey{
					tc.senderKey,
				},
				[]cryptotypes.PrivKey{
					tc.signKey,
				},
				tc.selectedAuthenticator,
			)
			s.Require().NoError(err)

			anteHandler := sdk.ChainAnteDecorators(s.AuthenticatorDecorator)
			_, err = anteHandler(s.Ctx.WithGasMeter(storetypes.NewGasMeter(300000)), tx, false)

			if tc.shouldPass {
				s.Require().NoError(err, "Expected to pass but got error")
			} else {
				s.Require().Error(err, "Expected to fail but got no error")
			}
		})
	}
}
