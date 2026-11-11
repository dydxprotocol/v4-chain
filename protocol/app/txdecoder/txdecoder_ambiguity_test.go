package txdecoder_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

// TestDuplicateTxRawFields_DecoderDivergence documents why a TxRaw that repeats body_bytes (field 1)
// is ambiguous: the same wire bytes decode to two different transactions under the two decoders that
// both run in production, because they interpret field 1 as different proto field types.
//
//   - authtx.DefaultTxDecoder (the SDK's tx decoder that admits/executes txs) decodes the raw bytes
//     as TxRaw, whose body_bytes is a scalar bytes field. Per the protobuf wire format, a repeated
//     scalar field is last-wins: only the second body_bytes entry survives, so only its messages are
//     decoded.
//
//   - CometBFT's mempool classifier (the dydxprotocol/cometbft fork's IsShortTermClobOrderTransaction
//     in mempool/dydx_helpers.go) instead unmarshals the raw bytes as cosmos-sdk/types/tx.Tx, whose
//     field 1 (Body *TxBody) is a singular embedded-message field. Per the protobuf spec, repeated
//     occurrences of a singular embedded message are merged, not overwritten: each occurrence calls
//     TxBody.Unmarshal again on the same *TxBody, and TxBody.Messages keeps appending. Both bodies'
//     messages survive.
//
// Because IsShortTermClobOrderTransaction requires exactly one message to classify a tx as a
// fee-free short-term clob order, the merged (2-message) view fails that check and bins the tx as
// "other", while the app's own decoder would execute it as a bare MsgPlaceOrder -- a classification
// divergence that leads to honest ProcessProposal rejection.
//
// This test asserts the wire-level divergence using the raw SDK/CometBFT decoder types directly (not
// the app's rejecting decoder), so it passes on the base revision by design. The end-to-end
// regression that the app rejects such bytes is in x/clob/e2e/malformed_tx_test.go.
func TestDuplicateTxRawFields_DecoderDivergence(t *testing.T) {
	msgSend := &banktypes.MsgSend{
		FromAddress: constants.AliceAccAddress.String(),
		ToAddress:   constants.BobAccAddress.String(),
		Amount:      sdktypes.NewCoins(sdktypes.NewInt64Coin(constants.Usdc.Denom, 1)),
	}
	msgPlaceOrder := &clobtypes.MsgPlaceOrder{
		Order: clobtypes.Order{
			OrderId: clobtypes.OrderId{
				SubaccountId: constants.Alice_Num0,
				ClientId:     0,
				ClobPairId:   0,
			},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     5,
			Subticks:     10,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 20},
		},
	}

	anySend, err := codectypes.NewAnyWithValue(msgSend)
	require.NoError(t, err)
	anyPlaceOrder, err := codectypes.NewAnyWithValue(msgPlaceOrder)
	require.NoError(t, err)

	body1 := &txtypes.TxBody{Messages: []*codectypes.Any{anySend}}
	body2 := &txtypes.TxBody{Messages: []*codectypes.Any{anyPlaceOrder}}

	body1Bytes, err := body1.Marshal()
	require.NoError(t, err)
	body2Bytes, err := body2.Marshal()
	require.NoError(t, err)

	// Build TxRaw wire bytes with TWO field-1 (body_bytes) entries, followed by field 2
	// (auth_info_bytes, an empty-but-valid AuthInfo). Tags are in non-decreasing order (1,1,2), so
	// the SDK's ADR-027 check -- which only rejects *decreasing* tag order -- accepts these bytes.
	// No signatures are needed to decode.
	var txRawBytes []byte
	txRawBytes = protowire.AppendTag(txRawBytes, 1, protowire.BytesType)
	txRawBytes = protowire.AppendBytes(txRawBytes, body1Bytes)
	txRawBytes = protowire.AppendTag(txRawBytes, 1, protowire.BytesType)
	txRawBytes = protowire.AppendBytes(txRawBytes, body2Bytes)
	txRawBytes = protowire.AppendTag(txRawBytes, 2, protowire.BytesType)
	txRawBytes = protowire.AppendBytes(txRawBytes, nil)

	// Decode path A: TxRaw.body_bytes is a scalar `bytes` field -- last-wins.
	cdc := encoding.GetTestEncodingCfg().Codec
	decodedTx, err := authtx.DefaultTxDecoder(cdc)(txRawBytes)
	require.NoError(t, err)
	require.Len(t, decodedTx.GetMsgs(), 1)
	_, isPlaceOrder := decodedTx.GetMsgs()[0].(*clobtypes.MsgPlaceOrder)
	require.True(t, isPlaceOrder, "gogoproto last-wins decode should keep only body2's MsgPlaceOrder")

	// Decode path B: Tx.Body is a singular embedded-message field -- merged. This is the exact type
	// CometBFT's mempool classifier unmarshals the same tx bytes into.
	var mergedTx txtypes.Tx
	err = mergedTx.Unmarshal(txRawBytes)
	require.NoError(t, err)
	require.Len(t, mergedTx.Body.Messages, 2)

	// The two decoders admit the SAME wire bytes as transactions with DIFFERENT message sets: one
	// message (path A) versus two merged messages (path B).
	require.NotEqual(t, len(decodedTx.GetMsgs()), len(mergedTx.Body.Messages))
}
