package ante_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestSigVerification(t *testing.T) {
	suite := testante.SetupTestSuite(t, true)
	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	// make block height non-zero to ensure account numbers part of signBytes
	suite.Ctx = suite.Ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()
	priv2, _, addr2 := testdata.KeyTestPubAddr()
	priv3, _, addr3 := testdata.KeyTestPubAddr()

	addrs := []sdk.AccAddress{addr1, addr2, addr3}

	// set accounts and create msg for each address
	for i, addr := range addrs {
		acc := suite.AccountKeeper.NewAccountWithAddress(suite.Ctx, addr)
		require.NoError(t, acc.SetAccountNumber(uint64(i)))
		suite.AccountKeeper.SetAccount(suite.Ctx, acc)
	}

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()

	spkd := sdkante.NewSetPubKeyDecorator(suite.AccountKeeper)
	svd := customante.NewSigVerificationDecorator(
		suite.AccountKeeper,
		suite.ClientCtx.TxConfig.SignModeHandler(),
	)
	antehandler := sdk.ChainAnteDecorators(spkd, svd)

	type testCase struct {
		name        string
		msgs        []sdk.Msg
		privs       []cryptotypes.PrivKey
		accNums     []uint64
		accSeqs     []uint64
		invalidSigs bool
		recheck     bool
		shouldErr   bool
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
		},
		{
			"not enough signers",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2},
			[]uint64{0, 1},
			[]uint64{0, 0},
			validSigs,
			false,
			true,
		},
		{
			"wrong order signers",
			testMsgs,
			[]cryptotypes.PrivKey{priv3, priv2, priv1},
			[]uint64{2, 1, 0},
			[]uint64{0, 0, 0},
			validSigs,
			false,
			true,
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
		},
		{
			"wrong sequences",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{0, 1, 2},
			[]uint64{3, 4, 5},
			validSigs,
			false,
			true,
		},
		{
			"wrong sequences but skip validation - place order",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr1)},
			[]cryptotypes.PrivKey{priv1},
			[]uint64{0},
			[]uint64{3},
			validSigs,
			false,
			false,
		},
		{
			"wrong sequences but skip validation - cancel order",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr2)},
			[]cryptotypes.PrivKey{priv2},
			[]uint64{1},
			[]uint64{4},
			validSigs,
			false,
			false,
		},
		{
			"wrong sequences but skip validation - transfer",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr3)},
			[]cryptotypes.PrivKey{priv3},
			[]uint64{2},
			[]uint64{5},
			validSigs,
			false,
			false,
		},
		{
			"wrong sequences - mixed messages",
			[]sdk.Msg{newPlaceOrderMessageForAddr(addr1), testdata.NewTestMsg(addr2)},
			[]cryptotypes.PrivKey{priv1, priv2},
			[]uint64{0, 1},
			[]uint64{3, 4},
			validSigs,
			false,
			true,
		},
		{
			"valid tx",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{0, 1, 2},
			[]uint64{0, 0, 0},
			validSigs,
			false,
			false,
		},
		{
			"no err on recheck",
			testMsgs,
			[]cryptotypes.PrivKey{priv1, priv2, priv3},
			[]uint64{0, 0, 0}, // account numbers are incorrect, but skips the check
			[]uint64{0, 0, 0},
			validSigs,
			true,
			false,
		},
	}
	for i, tc := range testCases {
		suite.Ctx = suite.Ctx.WithIsReCheckTx(tc.recheck)
		suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder() // Create new TxBuilder for each test

		require.NoError(t, suite.TxBuilder.SetMsgs(tc.msgs...))
		suite.TxBuilder.SetFeeAmount(feeAmount)
		suite.TxBuilder.SetGasLimit(gasLimit)

		tx, err := suite.CreateTestTx(tc.privs, tc.accNums, tc.accSeqs, suite.Ctx.ChainID())
		require.NoError(t, err)

		if tc.invalidSigs {
			txSigs, _ := tx.GetSignaturesV2()
			badSig, _ := tc.privs[0].Sign([]byte("unrelated message"))
			txSigs[0] = signing.SignatureV2{
				PubKey: tc.privs[0].PubKey(),
				Data: &signing.SingleSignatureData{
					SignMode:  suite.ClientCtx.TxConfig.SignModeHandler().DefaultMode(),
					Signature: badSig,
				},
				Sequence: tc.accSeqs[0],
			}
			err := suite.TxBuilder.SetSignatures(txSigs...)
			require.NoError(t, err)
			tx = suite.TxBuilder.GetTx()
		}

		_, err = antehandler(suite.Ctx, tx, false)
		if tc.shouldErr {
			require.NotNil(t, err, "TestCase %d: %s did not error as expected", i, tc.name)
		} else {
			require.Nil(t, err, "TestCase %d: %s errored unexpectedly. Err: %v", i, tc.name, err)
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

func runSigDecorators(t *testing.T, params types.Params, _ bool, privs ...cryptotypes.PrivKey) (sdk.Gas, error) {
	suite := testante.SetupTestSuite(t, true)
	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	// Make block-height non-zero to include accNum in SignBytes
	suite.Ctx = suite.Ctx.WithBlockHeight(1)
	err := suite.AccountKeeper.SetParams(suite.Ctx, params)
	require.NoError(t, err)

	msgs := make([]sdk.Msg, len(privs))
	accNums := make([]uint64, len(privs))
	accSeqs := make([]uint64, len(privs))
	// set accounts and create msg for each address
	for i, priv := range privs {
		addr := sdk.AccAddress(priv.PubKey().Address())
		acc := suite.AccountKeeper.NewAccountWithAddress(suite.Ctx, addr)
		require.NoError(t, acc.SetAccountNumber(uint64(i)))
		suite.AccountKeeper.SetAccount(suite.Ctx, acc)
		msgs[i] = testdata.NewTestMsg(addr)
		accNums[i] = uint64(i)
		accSeqs[i] = uint64(0)
	}
	require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))

	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.TxBuilder.SetFeeAmount(feeAmount)
	suite.TxBuilder.SetGasLimit(gasLimit)

	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.Ctx.ChainID())
	require.NoError(t, err)

	spkd := sdkante.NewSetPubKeyDecorator(suite.AccountKeeper)
	svgc := sdkante.NewSigGasConsumeDecorator(suite.AccountKeeper, sdkante.DefaultSigVerificationGasConsumer)
	svd := sdkante.NewSigVerificationDecorator(suite.AccountKeeper, suite.ClientCtx.TxConfig.SignModeHandler())
	antehandler := sdk.ChainAnteDecorators(spkd, svgc, svd)

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
