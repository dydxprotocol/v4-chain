package ante_test

import (
	"fmt"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsign "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
	accountplustypes "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSigVerification(t *testing.T) {
	suite := testante.SetupTestSuite(t, true)
	suite.TxBankKeeper.EXPECT().DenomMetadata(
		gomock.Any(),
		gomock.Any(),
	).Return(&banktypes.QueryDenomMetadataResponse{}, nil).AnyTimes()

	// signing.SignMode_SIGN_MODE_TEXTUAL and signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON are disabled
	// from the test suite since they include the sequence number as part of the signature preventing
	// us from skipping sequence validation for certain message types.
	enabledSignModes := []signing.SignMode{signing.SignMode_SIGN_MODE_DIRECT}

	// Since TEXTUAL is not enabled by default, we create a custom TxConfig
	// here which includes it.
	txConfigOpts := authtx.ConfigOptions{
		TextualCoinMetadataQueryFn: txmodule.NewGRPCCoinMetadataQueryFn(suite.ClientCtx),
		EnabledSignModes:           enabledSignModes,
	}
	var err error
	suite.ClientCtx.TxConfig, err = authtx.NewTxConfigWithOptions(
		codec.NewProtoCodec(suite.EncCfg.InterfaceRegistry),
		txConfigOpts,
	)
	require.NoError(t, err)
	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	// make block height non-zero to ensure account numbers part of signBytes
	suite.Ctx = suite.Ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	priv2, _, addr2 := testdata.KeyTestPubAddr()
	priv3, _, addr3 := testdata.KeyTestPubAddr()

	// initial accountplus AccountState
	maxEjectedNonce := uint64(testante.TestBlockTime - 1000)
	var timestampNonces []uint64
	for i := range accountpluskeeper.MaxTimestampNonceArrSize {
		timestampNonces = append(timestampNonces, testante.TestBlockTime+uint64(i)+1000)
	}

	addrs := []sdk.AccAddress{addr1, addr2, addr3}

	msgs := make([]sdk.Msg, len(addrs))
	accs := make([]sdk.AccountI, len(addrs))
	accStates := make([]accountplustypes.AccountState, len(addrs))

	// set accounts and create msg for each address
	for i, addr := range addrs {
		acc := suite.AccountKeeper.NewAccountWithAddress(suite.Ctx, addr)
		require.NoError(t, acc.SetAccountNumber(uint64(i)+1000))
		suite.AccountKeeper.SetAccount(suite.Ctx, acc)
		msgs[i] = testdata.NewTestMsg(addr)
		accs[i] = acc

		accStates[i] = accountplustypes.AccountState{
			Address: addr.String(),
			TimestampNonceDetails: accountplustypes.TimestampNonceDetails{
				MaxEjectedNonce: maxEjectedNonce,
				TimestampNonces: timestampNonces,
			},
		}
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()

	spkd := sdkante.NewSetPubKeyDecorator(suite.AccountKeeper)
	txConfigOpts = authtx.ConfigOptions{
		TextualCoinMetadataQueryFn: txmodule.NewBankKeeperCoinMetadataQueryFn(suite.TxBankKeeper),
		EnabledSignModes:           enabledSignModes,
	}
	anteTxConfig, err := authtx.NewTxConfigWithOptions(
		codec.NewProtoCodec(suite.EncCfg.InterfaceRegistry),
		txConfigOpts,
	)
	require.NoError(t, err)
	rpd := customante.NewReplayProtectionDecorator(
		suite.AccountKeeper,
		suite.AccountplusKeeper,
	)
	svd := customante.NewSigVerificationDecorator(
		suite.AccountKeeper,
		anteTxConfig.SignModeHandler(),
	)
	antehandler := sdk.ChainAnteDecorators(spkd, rpd, svd)
	defaultSignMode, err := authsign.APISignModeToInternal(anteTxConfig.SignModeHandler().DefaultMode())
	require.NoError(t, err)

	type testCase struct {
		name           string
		msgs           []sdk.Msg
		privs          []cryptotypes.PrivKey
		accNums        []uint64
		accSeqs        []uint64
		invalidSigs    bool // used for testing sigverify on RecheckTx
		recheck        bool
		shouldErr      bool
		expectedErrMsg string // supply empty string to ignore this check
		setAccState    bool   // used for ts nonce tests determine whether to use initial accountplus AccountState
	}

	testMsgs := make([]sdk.Msg, len(addrs))
	for i, addr := range addrs {
		testMsgs[i] = testdata.NewTestMsg(addr)
	}
	validSigs := false
	testCases := []testCase{
		{
			"no signers",
			testMsgs,
			[]cryptotypes.PrivKey{},
			[]uint64{},
			[]uint64{},
			validSigs,
			false,
			true,
			"",
			false,
		},
		{
			"not enough signers",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber()},
			[]uint64{0, 0},
			validSigs,
			false,
			true,
			"",
			false,
		},
		{
			"wrong order signers",
			testMsgs,
			[]cryptotypes.PrivKey{priv3, priv2, priv1},
			[]uint64{accs[2].GetAccountNumber(), accs[1].GetAccountNumber(), accs[0].GetAccountNumber()},
			[]uint64{0, 0, 0},
			validSigs,
			false,
			true,
			"",
			false,
		},
		{
			"wrong accnums",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{7, 8, 9},
			[]uint64{0, 0, 0},
			validSigs,
			false,
			true,
			"",
			false,
		},
		{
			"wrong sequences",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber(), accs[2].GetAccountNumber()},
			[]uint64{3, 4, 5},
			validSigs,
			false,
			true,
			"",
			false,
		},
		{
			"wrong sequences but skip validation - place order",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr1)},
			[]cryptotypes.PrivKey{priv1},
			[]uint64{accs[0].GetAccountNumber()},
			[]uint64{3},
			validSigs,
			false,
			false,
			"",
			false,
		},
		{
			"wrong sequences but skip validation - cancel order",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr2)},
			[]cryptotypes.PrivKey{priv2},
			[]uint64{accs[1].GetAccountNumber()},
			[]uint64{4},
			validSigs,
			false,
			false,
			"",
			false,
		},
		{
			"wrong sequences but skip validation - transfer",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr3)},
			[]cryptotypes.PrivKey{priv3},
			[]uint64{accs[2].GetAccountNumber()},
			[]uint64{5},
			validSigs,
			false,
			false,
			"",
			false,
		},
		{
			"wrong sequences - mixed messages",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr1), testdata.NewTestMsg(addr2)},
			[]cryptotypes.PrivKey{priv1, priv2},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber()},
			[]uint64{3, 4},
			validSigs,
			false,
			true,
			"",
			false,
		},
		{
			"valid tx",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber(), accs[2].GetAccountNumber()},
			[]uint64{0, 0, 0},
			validSigs,
			false,
			false,
			"",
			false,
		},
		{
			"no err on recheck",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{0, 0, 0},
			[]uint64{0, 0, 0},
			!validSigs,
			true,
			false,
			"",
			false,
		},
		{
			"invalid timestamp nonce",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber(), accs[2].GetAccountNumber()},
			[]uint64{
				testante.TestBlockTime - accountpluskeeper.MaxTimeInPastMs - 1,
				testante.TestBlockTime,
				testante.TestBlockTime,
			},
			validSigs,
			false,
			true,
			fmt.Sprintf(
				"timestamp nonce %d not within valid time window: incorrect account sequence",
				testante.TestBlockTime-accountpluskeeper.MaxTimeInPastMs-1),
			false,
		},
		{
			"can initialize AccountState",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber(), accs[2].GetAccountNumber()},
			[]uint64{
				testante.TestBlockTime, // any value within the window will work since no existing AccountState
				testante.TestBlockTime,
				testante.TestBlockTime,
			},
			validSigs,
			false,
			false,
			"",
			false,
		},
		{
			"reject timestamp nonce",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber(), accs[2].GetAccountNumber()},
			[]uint64{
				testante.TestBlockTime, // ts <= min(tsNonces)
				testante.TestBlockTime,
				testante.TestBlockTime,
			},
			validSigs,
			false,
			true,
			fmt.Sprintf("timestamp nonce %d rejected: incorrect account sequence", testante.TestBlockTime),
			true,
		},
		{
			"accept timestamp nonce",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber(), accs[2].GetAccountNumber()},
			[]uint64{
				testante.TestBlockTime + 5000, // ts > min(tsNonces)
				testante.TestBlockTime + 5000,
				testante.TestBlockTime + 5000,
			},
			validSigs,
			false,
			false,
			"",
			true,
		},
		{
			"timestamp nonce invalid sigs",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{accs[0].GetAccountNumber(), accs[1].GetAccountNumber(), accs[2].GetAccountNumber()},
			[]uint64{
				testante.TestBlockTime + 5000, // ts > min(tsNonces)
				testante.TestBlockTime + 5000,
				testante.TestBlockTime + 5000,
			},
			!validSigs,
			false,
			true,
			"",
			true,
		},
	}

	for i, tc := range testCases {
		for _, signMode := range enabledSignModes {
			t.Run(fmt.Sprintf("%s with %s", tc.name, signMode), func(t *testing.T) {
				suite.Ctx = suite.Ctx.WithIsReCheckTx(tc.recheck)
				suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder() // Create new txBuilder for each test

				require.NoError(t, suite.TxBuilder.SetMsgs(tc.msgs...))
				suite.TxBuilder.SetFeeAmount(feeAmount)
				suite.TxBuilder.SetGasLimit(gasLimit)

				// Set accountplus AccountStates
				if tc.setAccState {
					for _, acc := range accs {
						accState := accountplustypes.AccountState{
							Address: acc.GetAddress().String(),
							TimestampNonceDetails: accountplustypes.TimestampNonceDetails{
								MaxEjectedNonce: maxEjectedNonce,
								TimestampNonces: timestampNonces,
							},
						}
						suite.AccountplusKeeper.SetAccountState(suite.Ctx, acc.GetAddress(), accState)
					}
				}

				tx, err := suite.CreateTestTx(suite.Ctx, tc.privs, tc.accNums, tc.accSeqs, suite.Ctx.ChainID(), signMode)
				require.NoError(t, err)
				if tc.invalidSigs {
					txSigs, _ := tx.GetSignaturesV2()
					badSig, _ := tc.privs[0].Sign([]byte("unrelated message"))
					txSigs[0] = signing.SignatureV2{
						PubKey: tc.privs[0].PubKey(),
						Data: &signing.SingleSignatureData{
							SignMode:  defaultSignMode,
							Signature: badSig,
						},
						Sequence: tc.accSeqs[0],
					}
					require.NoError(t, suite.TxBuilder.SetSignatures(txSigs...))
					tx = suite.TxBuilder.GetTx()
				}

				txBytes, err := suite.ClientCtx.TxConfig.TxEncoder()(tx)
				require.NoError(t, err)
				byteCtx := suite.Ctx.WithTxBytes(txBytes)
				_, err = antehandler(byteCtx, tx, false)
				if tc.shouldErr {
					require.NotNil(t, err, "TestCase %d: %s did not error as expected", i, tc.name)
					if tc.expectedErrMsg != "" {
						require.Equal(t, tc.expectedErrMsg, err.Error())
					}
				} else {
					require.Nil(t, err, "TestCase %d: %s errored unexpectedly. Err: %v", i, tc.name, err)
				}
			})
		}
	}
}

func TestSigIntegration(t *testing.T) {
	// generate private keys
	privs := []cryptotypes.PrivKey{
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
		secp256k1.GenPrivKey(),
	}

	params := types.DefaultParams()
	initialSigCost := params.SigVerifyCostSecp256k1
	initialCost, err := runSigDecorators(t, params, false, privs...)
	require.Nil(t, err)

	params.SigVerifyCostSecp256k1 *= 2
	doubleCost, err := runSigDecorators(t, params, false, privs...)
	require.Nil(t, err)

	require.Equal(t, initialSigCost*uint64(len(privs)), doubleCost-initialCost)
}

func runSigDecorators(t *testing.T, params types.Params, _ bool, privs ...cryptotypes.PrivKey) (storetypes.Gas, error) {
	suite := testante.SetupTestSuite(t, true)
	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	// Make block-height non-zero to include accNum in SignBytes
	suite.Ctx = suite.Ctx.WithBlockHeight(1)
	err := suite.AccountKeeper.Params.Set(suite.Ctx, params)
	require.NoError(t, err)

	msgs := make([]sdk.Msg, len(privs))
	accNums := make([]uint64, len(privs))
	accSeqs := make([]uint64, len(privs))
	// set accounts and create msg for each address
	// set initial accountplus AccountState
	for i, priv := range privs {
		addr := sdk.AccAddress(priv.PubKey().Address())
		acc := suite.AccountKeeper.NewAccountWithAddress(suite.Ctx, addr)
		require.NoError(t, acc.SetAccountNumber(uint64(i)+1000))
		suite.AccountKeeper.SetAccount(suite.Ctx, acc)
		msgs[i] = testdata.NewTestMsg(addr)
		accNums[i] = acc.GetAccountNumber()
		accSeqs[i] = uint64(0)
	}
	require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.TxBuilder.SetFeeAmount(feeAmount)
	suite.TxBuilder.SetGasLimit(gasLimit)

	tx, err := suite.CreateTestTx(
		suite.Ctx,
		privs,
		accNums,
		accSeqs,
		suite.Ctx.ChainID(),
		signing.SignMode_SIGN_MODE_DIRECT,
	)
	require.NoError(t, err)

	spkd := sdkante.NewSetPubKeyDecorator(suite.AccountKeeper)
	svgc := sdkante.NewSigGasConsumeDecorator(suite.AccountKeeper, sdkante.DefaultSigVerificationGasConsumer)
	rpd := customante.NewReplayProtectionDecorator(
		suite.AccountKeeper,
		suite.AccountplusKeeper,
	)
	svd := customante.NewSigVerificationDecorator(
		suite.AccountKeeper,
		suite.ClientCtx.TxConfig.SignModeHandler(),
	)
	antehandler := sdk.ChainAnteDecorators(spkd, svgc, rpd, svd)

	txBytes, err := suite.ClientCtx.TxConfig.TxEncoder()(tx)
	require.NoError(t, err)
	suite.Ctx = suite.Ctx.WithTxBytes(txBytes)

	// Determine gas consumption of antehandler with default params
	before := suite.Ctx.GasMeter().GasConsumed()
	ctx, err := antehandler(suite.Ctx, tx, false)
	after := ctx.GasMeter().GasConsumed()

	return after - before, err
}

func newPlaceOrderMessageForAddr(addr sdk.AccAddress) sdk.Msg {
	return &clobtypes.MsgPlaceOrder{
		Order: clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: satypes.SubaccountId{
					Owner: addr.String(),
				},
			},
		},
	}
}
