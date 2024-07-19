package ante_test

import (
	"testing"

	testante "github.com/dydxprotocol/v4-chain/protocol/testutil/ante"
	"github.com/stretchr/testify/require"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	customante "github.com/dydxprotocol/v4-chain/protocol/app/ante"
	accountpluskeeper "github.com/dydxprotocol/v4-chain/protocol/x/accountplus/keeper"
)

// Modified from cosmossdk test for IncrementSequenceDecorator
func TestDydxIncrementSequenceDecorator(t *testing.T) {
	suite := testante.SetupTestSuite(t, true)
	suite.TxBuilder = suite.ClientCtx.TxConfig.NewTxBuilder()

	priv, _, addr := testdata.KeyTestPubAddr()
	acc := suite.AccountKeeper.NewAccountWithAddress(suite.Ctx, addr)
	require.NoError(t, acc.SetAccountNumber(uint64(50)))
	suite.AccountKeeper.SetAccount(suite.Ctx, acc)

	msgs := []sdk.Msg{testdata.NewTestMsg(addr)}
	require.NoError(t, suite.TxBuilder.SetMsgs(msgs...))
	privs := []cryptotypes.PrivKey{priv}
	accNums := []uint64{suite.AccountKeeper.GetAccount(suite.Ctx, addr).GetAccountNumber()}
	feeAmount := testdata.NewTestFeeAmount()
	gasLimit := testdata.NewTestGasLimit()
	suite.TxBuilder.SetFeeAmount(feeAmount)
	suite.TxBuilder.SetGasLimit(gasLimit)

	isd := customante.NewDydxIncrementSequenceDecorator(suite.AccountKeeper)
	antehandler := sdk.ChainAnteDecorators(isd)

	testCases := []struct {
		ctx      sdk.Context
		simulate bool
		// This value need not be valid (accountSeq + 1). Validity is handed in customante.NewSigVerificationDecorator
		signatureSeq uint64
		expectedSeq  uint64
	}{
		// tests from cosmossdk checking incrementing seqence
		{suite.Ctx.WithIsReCheckTx(true), false, 0, 1},
		{suite.Ctx.WithIsCheckTx(true).WithIsReCheckTx(false), false, 0, 2},
		{suite.Ctx.WithIsReCheckTx(true), false, 0, 3},
		{suite.Ctx.WithIsReCheckTx(true), false, 0, 4},
		{suite.Ctx.WithIsReCheckTx(true), true, 0, 5},

		// tests checking that tx with timestamp nonces will not increment sequence
		{suite.Ctx.WithIsReCheckTx(true), true, accountpluskeeper.TimestampNonceSequenceCutoff, 5},
		{suite.Ctx.WithIsReCheckTx(true), true, accountpluskeeper.TimestampNonceSequenceCutoff + 100000, 5},
	}

	for i, tc := range testCases {
		accSeqs := []uint64{tc.signatureSeq}
		tx, err := suite.CreateTestTx(
			suite.Ctx,
			privs,
			accNums,
			accSeqs,
			suite.Ctx.ChainID(),
			signing.SignMode_SIGN_MODE_DIRECT,
		)
		require.NoError(t, err)

		_, err = antehandler(tc.ctx, tx, tc.simulate)
		require.NoError(t, err, "unexpected error; tc #%d, %v", i, tc)
		require.Equal(t, tc.expectedSeq, suite.AccountKeeper.GetAccount(suite.Ctx, addr).GetSequence())
	}
}
