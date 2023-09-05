package ante_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	kmultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256r1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	sdkante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestSetPubKey(t *testing.T) {
	suite := testante.SetupTestSuite(t, true)
	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	// keys and addresses
	priv1, pub1, addr1 := testdata.KeyTestPubAddr()
	priv2, pub2, addr2 := testdata.KeyTestPubAddr()
	priv3, pub3, addr3 := testdata.KeyTestPubAddrSecp256R1(require.New(t))

	addrs := []sdk.AccAddress{addr1, addr2, addr3}
	pubs := []cryptotypes.PubKey{pub1, pub2, pub3}

	msgs := make([]sdk.Msg, len(addrs))
	// set accounts and create msg for each address
	for i, addr := range addrs {
		acc := suite.AccountKeeper.NewAccountWithAddress(suite.Ctx, addr)
		require.NoError(t, acc.SetAccountNumber(uint64(i)))
		suite.AccountKeeper.SetAccount(suite.Ctx, acc)
		msgs[i] = testdata.NewTestMsg(addr)
	}
	require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))
	suite.TxBuilder.SetFeeAmount(testdata.NewTestFeeAmount())
	suite.TxBuilder.SetGasLimit(testdata.NewTestGasLimit())

	privs, accNums, accSeqs := []cryptotypes.PrivKey{priv1, priv2, priv3}, []uint64{0, 1, 2}, []uint64{0, 0, 0}
	tx, err := suite.CreateTestTx(privs, accNums, accSeqs, suite.Ctx.ChainID())
	require.NoError(t, err)

	spkd := sdkante.NewSetPubKeyDecorator(suite.AccountKeeper)
	antehandler := sdk.ChainAnteDecorators(spkd)

	ctx, err := antehandler(suite.Ctx, tx, false)
	require.NoError(t, err)

	// Require that all accounts have pubkey set after Decorator runs
	for i, addr := range addrs {
		pk, err := suite.AccountKeeper.GetPubKey(ctx, addr)
		require.NoError(t, err, "Error on retrieving pubkey from account")
		require.True(t, pubs[i].Equals(pk),
			"Wrong Pubkey retrieved from AccountKeeper, idx=%d\nexpected=%s\n     got=%s", i, pubs[i], pk)
	}
}

func TestConsumeSignatureVerificationGas(t *testing.T) {
	suite := testante.SetupTestSuite(t, true)
	params := types.DefaultParams()
	msg := []byte{1, 2, 3, 4}

	p := types.DefaultParams()
	skR1, _ := secp256r1.GenPrivKey()
	pkSet1, sigSet1 := generatePubKeysAndSignatures(5, msg, false)
	multisigKey1 := kmultisig.NewLegacyAminoPubKey(2, pkSet1)
	multisignature1 := multisig.NewMultisig(len(pkSet1))
	expectedCost1 := expectedGasCostByKeys(pkSet1)
	for i := 0; i < len(pkSet1); i++ {
		// Ignore "SA1019: legacytx.StdSignature is deprecated: StdSignature represents a sig" error for testing.
		//nolint:staticcheck
		stdSig := legacytx.StdSignature{PubKey: pkSet1[i], Signature: sigSet1[i]}
		sigV2, err := legacytx.StdSignatureToSignatureV2(suite.ClientCtx.LegacyAmino, stdSig)
		require.NoError(t, err)
		err = multisig.AddSignatureV2(multisignature1, sigV2, pkSet1)
		require.NoError(t, err)
	}

	type args struct {
		meter  sdk.GasMeter
		sig    signing.SignatureData
		pubkey cryptotypes.PubKey
		params types.Params
	}
	tests := []struct {
		name        string
		args        args
		gasConsumed uint64
		shouldErr   bool
	}{
		{
			"PubKeyEd25519",
			args{
				sdk.NewInfiniteGasMeter(),
				nil,
				ed25519.GenPrivKey().PubKey(),
				params,
			},
			p.SigVerifyCostED25519,
			true,
		},
		{
			"PubKeySecp256k1",
			args{
				sdk.NewInfiniteGasMeter(),
				nil,
				secp256k1.GenPrivKey().PubKey(),
				params,
			},
			p.SigVerifyCostSecp256k1,
			false,
		},
		{
			"PubKeySecp256r1",
			args{
				sdk.NewInfiniteGasMeter(),
				nil,
				skR1.PubKey(),
				params,
			},
			p.SigVerifyCostSecp256r1(),
			false,
		},
		{
			"Multisig",
			args{
				sdk.NewInfiniteGasMeter(),
				multisignature1,
				multisigKey1,
				params,
			},
			expectedCost1,
			false,
		},
		{
			"unknown key",
			args{
				sdk.NewInfiniteGasMeter(),
				nil,
				nil,
				params,
			},
			0,
			true,
		},
	}
	for _, tt := range tests {
		sigV2 := signing.SignatureV2{
			PubKey:   tt.args.pubkey,
			Data:     tt.args.sig,
			Sequence: 0, // Arbitrary account sequence
		}
		err := sdkante.DefaultSigVerificationGasConsumer(tt.args.meter, sigV2, tt.args.params)

		if tt.shouldErr {
			require.NotNil(t, err)
		} else {
			require.Nil(t, err)
			require.Equal(
				t,
				tt.gasConsumed,
				tt.args.meter.GasConsumed(),
				fmt.Sprintf("%d != %d", tt.gasConsumed, tt.args.meter.GasConsumed()),
			)
		}
	}
}

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
	err := suite.AccountKeeper.Params.Set(suite.Ctx, params)
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
