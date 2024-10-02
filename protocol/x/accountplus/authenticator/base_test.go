package authenticator_test

import (
	"encoding/hex"
	"fmt"
	"math/rand"

	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	testtx "github.com/dydxprotocol/v4-chain/protocol/testutil/tx"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	smartaccounttypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/suite"
)

type BaseAuthenticatorSuite struct {
	suite.Suite
	tApp                         *testapp.TestApp
	Ctx                          sdk.Context
	EncodingConfig               app.EncodingConfig
	SigVerificationAuthenticator authenticator.SignatureVerification
	TestKeys                     []string
	TestAccAddress               []sdk.AccAddress
	TestPrivKeys                 []*secp256k1.PrivKey
	HomeDir                      string
}

func (s *BaseAuthenticatorSuite) SetupKeys() {
	// Test data for authenticator signature verification
	TestKeys := []string{
		"6cf5103c60c939a5f38e383b52239c5296c968579eec1c68a47d70fbf1d19159",
		"0dd4d1506e18a5712080708c338eb51ecf2afdceae01e8162e890b126ac190fe",
		"49006a359803f0602a7ec521df88bf5527579da79112bb71f285dd3e7d438033",
	}
	accounts := make([]sdk.AccountI, 0)
	// Set up test accounts
	for _, key := range TestKeys {
		bz, _ := hex.DecodeString(key)
		priv := &secp256k1.PrivKey{Key: bz}

		// add the test private keys to array for later use
		s.TestPrivKeys = append(s.TestPrivKeys, priv)

		accAddress := sdk.AccAddress(priv.PubKey().Address())
		accounts = append(
			accounts,
			authtypes.NewBaseAccount(accAddress, priv.PubKey(), 0, 0),
		)

		// add the test accounts to array for later use
		s.TestAccAddress = append(s.TestAccAddress, accAddress)
	}

	s.HomeDir = fmt.Sprintf("%d", rand.Int())
	s.tApp = testapp.NewTestAppBuilder(s.T()).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
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
	s.Ctx = s.Ctx.WithGasMeter(storetypes.NewGasMeter(1_000_000))

	s.EncodingConfig = app.GetEncodingConfig()
}

func (s *BaseAuthenticatorSuite) GenSimpleTx(msgs []sdk.Msg, signers []cryptotypes.PrivKey) (sdk.Tx, error) {
	txconfig := app.GetEncodingConfig().TxConfig
	feeCoins := constants.TestFeeCoins_5Cents
	var accNums []uint64
	var accSeqs []uint64

	ak := s.tApp.App.AccountKeeper

	for _, signer := range signers {
		var account sdk.AccountI
		if ak.HasAccount(s.Ctx, sdk.AccAddress(signer.PubKey().Address())) {
			account = ak.GetAccount(s.Ctx, sdk.AccAddress(signer.PubKey().Address()))
		} else {
			account = authtypes.NewBaseAccount(
				sdk.AccAddress(signer.PubKey().Address()),
				signer.PubKey(),
				ak.NextAccountNumber(s.Ctx),
				0,
			)
		}
		accNums = append(accNums, account.GetAccountNumber())
		accSeqs = append(accSeqs, account.GetSequence())
	}

	tx, err := testtx.GenTx(
		s.Ctx,
		txconfig,
		msgs,
		feeCoins,
		300000,
		"",
		accNums,
		accSeqs,
		signers,
		signers,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *BaseAuthenticatorSuite) GenSimpleTxWithSelectedAuthenticators(
	msgs []sdk.Msg,
	signers []cryptotypes.PrivKey,
	selectedAuthenticators []uint64,
) (sdk.Tx, error) {
	txconfig := app.GetEncodingConfig().TxConfig
	feeCoins := constants.TestFeeCoins_5Cents
	var accNums []uint64
	var accSeqs []uint64

	ak := s.tApp.App.AccountKeeper

	for _, signer := range signers {
		account := ak.GetAccount(s.Ctx, sdk.AccAddress(signer.PubKey().Address()))
		accNums = append(accNums, account.GetAccountNumber())
		accSeqs = append(accSeqs, account.GetSequence())
	}

	baseTxBuilder, err := testtx.MakeTxBuilder(
		s.Ctx,
		txconfig,
		msgs,
		feeCoins,
		300000,
		"",
		accNums,
		accSeqs,
		signers,
		signers,
		nil,
	)
	if err != nil {
		return nil, err
	}

	txBuilder, ok := baseTxBuilder.(authtx.ExtensionOptionsTxBuilder)
	if !ok {
		return nil, fmt.Errorf("expected authtx.ExtensionOptionsTxBuilder, got %T", baseTxBuilder)
	}
	if len(selectedAuthenticators) > 0 {
		value, err := codectypes.NewAnyWithValue(&smartaccounttypes.TxExtension{
			SelectedAuthenticators: selectedAuthenticators,
		})
		if err != nil {
			return nil, err
		}
		txBuilder.SetNonCriticalExtensionOptions(value)
	}

	tx := txBuilder.GetTx()
	return tx, nil
}

// FundAcc funds target address with specified amount.
func (s *BaseAuthenticatorSuite) FundAcc(acc sdk.AccAddress, amounts sdk.Coins) {
	err := testutil.FundAccount(s.Ctx, s.tApp.App.BankKeeper, acc, amounts)
	s.Require().NoError(err)
}
