package ante

import (
	"testing"
	"time"

	storetypes "cosmossdk.io/store/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txtestutil "github.com/cosmos/cosmos-sdk/x/auth/tx/testutil"
	"github.com/cosmos/cosmos-sdk/x/bank"
	v4module "github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	accountplustypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"

	perpetualskeeper "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/keeper"

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

// Use nice number to make debugging easier
const TestBlockTime = uint64(1600000000000) // sep 13 2020 12:26:40

// TestAccount represents an account used in the tests in x/auth/ante.
type TestAccount struct {
	acc  sdk.AccountI
	priv cryptotypes.PrivKey
}

// AnteTestSuite is a test suite to be used with ante handler tests.
type AnteTestSuite struct {
	AnteHandler       sdk.AnteHandler
	Ctx               sdk.Context
	ClientCtx         client.Context
	TxBuilder         client.TxBuilder
	AccountKeeper     keeper.AccountKeeper
	AccountplusKeeper accountpluskeeper.Keeper
	BankKeeper        *authtestutil.MockBankKeeper
	TxBankKeeper      *txtestutil.MockBankKeeper
	FeeGrantKeeper    *antetestutil.MockFeegrantKeeper
	PerpetualsKeeper  perpetualskeeper.Keeper
	EncCfg            moduletestutil.TestEncodingConfig
}

// SetupTest setups a new test, with new app, context, and anteHandler.
func SetupTestSuite(t testing.TB, isCheckTx bool) *AnteTestSuite {
	suite := &AnteTestSuite{}
	ctrl := gomock.NewController(t)
	suite.BankKeeper = authtestutil.NewMockBankKeeper(ctrl)
	suite.TxBankKeeper = txtestutil.NewMockBankKeeper(ctrl)
	suite.FeeGrantKeeper = antetestutil.NewMockFeegrantKeeper(ctrl)

	keys := map[string]*storetypes.KVStoreKey{
		types.StoreKey:            storetypes.NewKVStoreKey(types.StoreKey),
		accountplustypes.StoreKey: storetypes.NewKVStoreKey(accountplustypes.StoreKey),
	}
	transKeys := map[string]*storetypes.TransientStoreKey{
		"transient_test": storetypes.NewTransientStoreKey("transient_test"),
	}
	memKeys := map[string]*storetypes.MemoryStoreKey{}

	testCtx := testutil.DefaultContextWithKeys(keys, transKeys, memKeys)
	suite.Ctx = testCtx.WithIsCheckTx(isCheckTx).WithBlockHeight(1).WithBlockTime(time.UnixMilli(int64(TestBlockTime)))

	suite.EncCfg = MakeTestEncodingConfig(auth.AppModuleBasic{}, bank.AppModuleBasic{})

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
		runtime.NewKVStoreService(keys[types.StoreKey]),
		types.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		types.NewModuleAddress("gov").String(),
	)
	suite.AccountKeeper.GetModuleAccount(suite.Ctx, types.FeeCollectorName)
	err := suite.AccountKeeper.Params.Set(suite.Ctx, types.DefaultParams())
	require.NoError(t, err)

	// Initialize accountplus keeper
	suite.AccountplusKeeper = *accountpluskeeper.NewKeeper(
		suite.EncCfg.Codec,
		keys[accountplustypes.StoreKey],
		authenticator.NewAuthenticatorManager(),
		[]string{
			lib.GovModuleAddress.String(),
		},
	)

	// We're using TestMsg encoding in some tests, so register it here.
	suite.EncCfg.Amino.RegisterConcrete(&testdata.TestMsg{}, "testdata.TestMsg", nil)
	testdata.RegisterInterfaces(suite.EncCfg.InterfaceRegistry)

	suite.ClientCtx = client.Context{}.
		WithTxConfig(suite.EncCfg.TxConfig).
		WithClient(clitestutil.NewMockCometRPC(abci.ResponseQuery{}))

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   suite.AccountKeeper,
			BankKeeper:      suite.BankKeeper,
			FeegrantKeeper:  suite.FeeGrantKeeper,
			SignModeHandler: suite.EncCfg.TxConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)

	require.NoError(t, err)
	suite.AnteHandler = anteHandler

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

func (t TestCaseArgs) WithAccountsInfo(accs []TestAccount) TestCaseArgs {
	newT := t
	for _, acc := range accs {
		newT.accNums = append(newT.accNums, acc.acc.GetAccountNumber())
		newT.accSeqs = append(newT.accSeqs, acc.acc.GetSequence())
		newT.privs = append(newT.privs, acc.priv)
	}
	return newT
}

// DeliverMsgs constructs a tx and runs it through the ante handler. This is used to set the context for a
// test case, for example to test for replay protection.
func (suite *AnteTestSuite) DeliverMsgs(
	t testing.TB,
	privs []cryptotypes.PrivKey,
	msgs []sdk.Msg,
	feeAmount sdk.Coins,
	gasLimit uint64,
	accNums, accSeqs []uint64,
	chainID string,
	simulate bool,
) (sdk.Context, error) {
	require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))
	suite.TxBuilder.SetFeeAmount(feeAmount)
	suite.TxBuilder.SetGasLimit(gasLimit)

	tx, txErr := suite.CreateTestTx(suite.Ctx, privs, accNums, accSeqs, chainID, signing.SignMode_SIGN_MODE_DIRECT)
	require.NoError(t, txErr)
	txBytes, err := suite.ClientCtx.TxConfig.TxEncoder()(tx)
	bytesCtx := suite.Ctx.WithTxBytes(txBytes)
	require.NoError(t, err)
	return suite.AnteHandler(bytesCtx, tx, simulate)
}

func (suite *AnteTestSuite) RunTestCase(t testing.TB, tc TestCase, args TestCaseArgs) {
	require.NoError(t, suite.TxBuilder.SetMsgs(args.msgs...))
	suite.TxBuilder.SetFeeAmount(args.feeAmount)
	suite.TxBuilder.SetGasLimit(args.gasLimit)

	// Theoretically speaking, ante handler unit tests should only test
	// ante handlers, but here we sometimes also test the tx creation
	// process.
	tx, txErr := suite.CreateTestTx(
		suite.Ctx,
		args.privs,
		args.accNums,
		args.accSeqs,
		args.chainID,
		signing.SignMode_SIGN_MODE_DIRECT,
	)
	txBytes, err := suite.ClientCtx.TxConfig.TxEncoder()(tx)
	require.NoError(t, err)
	bytesCtx := suite.Ctx.WithTxBytes(txBytes)
	newCtx, anteErr := suite.AnteHandler(bytesCtx, tx, tc.simulate)

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
	ctx sdk.Context, privs []cryptotypes.PrivKey,
	accNums, accSeqs []uint64,
	chainID string, signMode signing.SignMode,
) (xauthsigning.Tx, error) {
	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signMode,
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
			Address:       sdk.AccAddress(priv.PubKey().Address()).String(),
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
			PubKey:        priv.PubKey(),
		}
		sigV2, err := tx.SignWithPrivKey(
			ctx, signMode, signerData,
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

func MakeTestEncodingConfig(modules ...module.AppModuleBasic) moduletestutil.TestEncodingConfig {
	aminoCodec := codec.NewLegacyAmino()
	interfaceRegistry := v4module.InterfaceRegistry
	codec := codec.NewProtoCodec(interfaceRegistry)

	encCfg := moduletestutil.TestEncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Codec:             codec,
		TxConfig:          authtx.NewTxConfig(codec, authtx.DefaultSignModes),
		Amino:             aminoCodec,
	}

	mb := module.NewBasicManager(modules...)

	std.RegisterLegacyAminoCodec(encCfg.Amino)
	std.RegisterInterfaces(encCfg.InterfaceRegistry)
	mb.RegisterLegacyAminoCodec(encCfg.Amino)
	mb.RegisterInterfaces(encCfg.InterfaceRegistry)

	return encCfg
}
