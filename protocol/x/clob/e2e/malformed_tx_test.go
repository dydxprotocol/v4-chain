package clob_test

import (
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

// mutateAddDuplicateBodyBytes prepends an extra body_bytes (field 1) entry to an encoded TxRaw,
// producing bytes that repeat that field and must be rejected by the duplicate-field check.
func mutateAddDuplicateBodyBytes(txBytes []byte) []byte {
	dup := protowire.AppendTag(nil, 1, protowire.BytesType)
	dup = protowire.AppendBytes(dup, []byte("x"))
	return append(dup, txBytes...)
}

func TestCheckTx_RejectsDuplicateTxRawBody(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	validCheckTx := testapp.MustMakeCheckTx(
		ctx,
		tApp.App,
		testapp.MustMakeCheckTxOptions{AccAddressForSigning: constants.AliceAccAddress.String()},
		&PlaceOrder_Alice_Num0_Id0_Clob0_Buy5_Price10_GTB20,
	)

	// Control: the un-mutated tx is accepted by CheckTx.
	resp := tApp.CheckTx(validCheckTx)
	require.Conditionf(t, resp.IsOK, "Expected CheckTx to succeed. Response: %+v", resp)

	mutatedCheckTx := abcitypes.RequestCheckTx{
		Tx:   mutateAddDuplicateBodyBytes(validCheckTx.Tx),
		Type: abcitypes.CheckTxType_New,
	}
	mutatedResp := tApp.CheckTx(mutatedCheckTx)
	require.Equal(t, sdkerrors.ErrTxDecode.ABCICode(), mutatedResp.Code)
}

func TestProcessProposal_RejectsDuplicateTxRawBody(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	tApp.InitChain()
	tApp.AdvanceToBlock(2, testapp.AdvanceToBlockOptions{})

	proposal, err := tApp.PrepareProposal()
	require.NoError(t, err)
	require.NotEmpty(t, proposal.Txs)

	// Mutate the (valid) ProposedOperationsTx at index 0 to repeat its body_bytes field. Mutating an
	// otherwise-valid tx (rather than swapping in a tx of the wrong type) ensures that, absent the
	// duplicate-field check, the proposal would decode and validate successfully.
	proposal.Txs[0] = mutateAddDuplicateBodyBytes(proposal.Txs[0])

	tApp.AdvanceToBlock(3, testapp.AdvanceToBlockOptions{
		RequestProcessProposalTxsOverride: proposal.Txs,
		ValidateRespProcess: func(ctx sdktypes.Context, resp abcitypes.ResponseProcessProposal) (stopProcessing bool) {
			require.Equal(t, abcitypes.ResponseProcessProposal_REJECT, resp.Status)
			return true // stop processing this rejected block in the test harness.
		},
	})
}
