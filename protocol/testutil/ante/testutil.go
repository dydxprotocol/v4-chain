package ante

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	antetestutil "github.com/cosmos/cosmos-sdk/x/auth/ante/testutil"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtestutil "github.com/cosmos/cosmos-sdk/x/auth/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// TestAccount represents an account used in the tests in x/auth/ante.
type TestAccount struct {
	acc  types.AccountI
	priv cryptotypes.PrivKey
}

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	AnteHandler    sdk.AnteHandler
	Ctx            sdk.Context
	ClientCtx      client.Context
	TxBuilder      client.TxBuilder
	AccountKeeper  keeper.AccountKeeper
	BankKeeper     *authtestutil.MockBankKeeper
	FeeGrantKeeper *antetestutil.MockFeegrantKeeper
	EncCfg         moduletestutil.TestEncodingConfig
}

// SetupTest setups a new test, with new app, context, and AnteHandler.
func SetupTestSuite(t *testing.T, isCheckTx bool) *AnteTestSuite {
	suite := &AnteTestSuite{}
	ctrl := gomock.NewController(t)
	suite.BankKeeper = authtestutil.NewMockBankKeeper(ctrl)

	suite.FeeGrantKeeper = antetestutil.NewMockFeegrantKeeper(ctrl)

	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	suite.Ctx = testCtx.Ctx.WithIsCheckTx(isCheckTx).WithBlockHeight(1)
	// app.BaseApp.NewContext(isCheckTx, tmproto.Header{}).WithBlockHeight(1)
	suite.EncCfg = moduletestutil.MakeTestEncodingConfig(auth.AppModuleBasic{})

	maccPerms := map[string][]string{
		"fee_collector":          nil,
		"mint":                   {"minter"},
		"bonded_tokens_pool":     {"burner", "staking"},
		"not_bonded_tokens_pool": {"burner", "staking"},
		"multiPerm":              {"burner", "minter", "staking"},
		"random":                 {"random"},
	}

	suite.AccountKeeper = keeper.NewAccountKeeper(
		suite.EncCfg.Codec,
		key,
		types.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32MainPrefix,
		types.NewModuleAddress("gov").String(),
	)
	suite.AccountKeeper.GetModuleAccount(suite.Ctx, types.FeeCollectorName)
	err := suite.AccountKeeper.SetParams(suite.Ctx, types.DefaultParams())
	require.NoError(t, err)

	// We're using TestMsg encoding in some tests, so register it here.
	suite.EncCfg.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(suite.EncCfg.InterfaceRegistry)

	suite.ClientCtx = client.Context{}.
		WithTxConfig(suite.EncCfg.TxConfig)

	AnteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   suite.AccountKeeper,
			BankKeeper:      suite.BankKeeper,
			FeegrantKeeper:  suite.FeeGrantKeeper,
			SignModeHandler: suite.EncCfg.TxConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)

	require.NoError(t, err)
	suite.AnteHandler = AnteHandler

	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	return suite
}

func (suite *AnteTestSuite) CreateTestAccounts(numAccs int) []TestAccount {
	var accounts []TestAccount

	for i := 0; i < numAccs; i++ {
		priv, _, addr := testdata.KeyTestPubAddr()
		acc := suite.AccountKeeper.NewAccountWithAddress(suite.Ctx, addr)
		if err := acc.SetAccountNumber(uint64(i)); err != nil {
			panic(err)
		}
		suite.AccountKeeper.SetAccount(suite.Ctx, acc)
		accounts = append(accounts, TestAccount{acc, priv})
	}

	return accounts
}

// TestCase represents a test case used in test tables.
type TestCase struct {
	// Comment out the follow `desc` and `malleate` fields since they are not used.
	// desc     string
	// malleate func(*AnteTestSuite) TestCaseArgs
	simulate bool
	expPass  bool
	expErr   error
}

type TestCaseArgs struct {
	chainID   string
	accNums   []uint64
	accSeqs   []uint64
	feeAmount sdk.Coins
	gasLimit  uint64
	msgs      []sdk.Msg
	privs     []cryptotypes.PrivKey
}

// DeliverMsgs constructs a tx and runs it through the ante handler. This is used to set the context for a test case,
// for example to test for replay protection.
func (suite *AnteTestSuite) DeliverMsgs(
	t *testing.T,
	privs []cryptotypes.PrivKey,
	msgs []sdk.Msg,
	feeAmount sdk.Coins,
	gasLimit uint64,
	accNums,
	accSeqs []uint64,
	chainID string,
	simulate bool,
) (sdk.Context, error) {
	require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))
	suite.TxBuilder.SetFeeAmount(feeAmount)
	suite.TxBuilder.SetGasLimit(gasLimit)

	tx, txErr := suite.CreateTestTx(privs, accNums, accSeqs, chainID)
	require.NoError(t, txErr)
	return suite.AnteHandler(suite.Ctx, tx, simulate)
}

func (suite *AnteTestSuite) RunTestCase(t *testing.T, tc TestCase, args TestCaseArgs) {
	require.NoError(t, suite.TxBuilder.SetMsgs(args.msgs...))
	suite.TxBuilder.SetFeeAmount(args.feeAmount)
	suite.TxBuilder.SetGasLimit(args.gasLimit)

	// Theoretically speaking, ante handler unit tests should only test
	// ante handlers, but here we sometimes also test the tx creation
	// process.
	tx, txErr := suite.CreateTestTx(args.privs, args.accNums, args.accSeqs, args.chainID)
	newCtx, anteErr := suite.AnteHandler(suite.Ctx, tx, tc.simulate)

	if tc.expPass {
		require.NoError(t, txErr)
		require.NoError(t, anteErr)
		require.NotNil(t, newCtx)

		suite.Ctx = newCtx
	} else {
		switch {
		case txErr != nil:
			require.Error(t, txErr)
			require.ErrorIs(t, txErr, tc.expErr)

		case anteErr != nil:
			require.Error(t, anteErr)
			require.ErrorIs(t, anteErr, tc.expErr)

		default:
			t.Fatal("expected one of txErr, anteErr to be an error")
		}
	}
}

// CreateTestTx is a helper function to create a tx given multiple inputs.
func (suite *AnteTestSuite) CreateTestTx(
	privs []cryptotypes.PrivKey,
	accNums []uint64,
	accSeqs []uint64,
	chainID string,
) (xauthsigning.Tx, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  suite.ClientCtx.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := suite.TxBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
			suite.ClientCtx.TxConfig.SignModeHandler().DefaultMode(), signerData,
			suite.TxBuilder, priv, suite.ClientCtx.TxConfig, accSeqs[i])
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = suite.TxBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return suite.TxBuilder.GetTx(), nil
}
