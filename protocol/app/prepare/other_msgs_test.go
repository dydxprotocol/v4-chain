package prepare_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/app/prepare"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	"github.com/stretchr/testify/require"
)

var (
	_ = constants.TestTxBuilder.SetMsgs(
		constants.ValidMsgAddPremiumVotes,
		constants.ValidMsgUpdateMarketPrices,
	)
	multiMsgsTxHasDisallowOnlyTxBytes, _ = constants.TestEncodingCfg.TxConfig.TxEncoder()(
		constants.TestTxBuilder.GetTx())

	_ = constants.TestTxBuilder.SetMsgs(
		constants.Msg_Send,
		constants.ValidMsgUpdateMarketPrices,
	)
	multiMsgsTxHasDisallowMixedTxBytes, _ = constants.TestEncodingCfg.TxConfig.TxEncoder()(
		constants.TestTxBuilder.GetTx())
)

func TestGetGroupMsgOther(t *testing.T) {
	tests := map[string]struct {
		txs      [][]byte
		maxBytes uint64

		expectedTxsInclude   [][]byte
		expectedTxsRemainder [][]byte
	}{
		"nil available txs": {
			txs:      nil,
			maxBytes: 10,

			expectedTxsInclude:   nil,
			expectedTxsRemainder: nil,
		},
		"empty available txs": {
			txs:      [][]byte{},
			maxBytes: 10,

			expectedTxsInclude:   nil,
			expectedTxsRemainder: nil,
		},
		"no tx fits under max bytes": {
			txs:      [][]byte{{1, 2}, {3, 4}},
			maxBytes: 1,

			expectedTxsInclude:   nil,
			expectedTxsRemainder: [][]byte{{1, 2}, {3, 4}},
		},
		"valid: subset under max": {
			txs:      [][]byte{{1, 2}, {3, 4}},
			maxBytes: 3,

			expectedTxsInclude:   [][]byte{{1, 2}},
			expectedTxsRemainder: [][]byte{{3, 4}},
		},
		"valid: all under max": {
			txs:      [][]byte{{1, 2}, {3, 4}},
			maxBytes: 4,

			expectedTxsInclude:   [][]byte{{1, 2}, {3, 4}},
			expectedTxsRemainder: nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			include, remainder := prepare.GetGroupMsgOther(tc.txs, tc.maxBytes)

			require.Equal(t, tc.expectedTxsInclude, include)
			require.Equal(t, tc.expectedTxsRemainder, remainder)
		})
	}
}

func TestRemoveDisallowMsgs(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()

	tests := map[string]struct {
		txs         [][]byte
		expectedTxs [][]byte
	}{
		"Nil": {
			txs:         nil,
			expectedTxs: nil,
		},
		"Empty": {
			txs:         [][]byte{},
			expectedTxs: nil,
		},
		"Single Tx, Single Msg Tx, Disallowed Msg": {
			txs:         [][]byte{constants.ValidMsgUpdateMarketPricesTxBytes},
			expectedTxs: nil,
		},
		"Single Tx, Single Msg Tx, Allowed Msg": {
			txs:         [][]byte{constants.Msg_Send_TxBytes},
			expectedTxs: [][]byte{constants.Msg_Send_TxBytes},
		},
		"Single Tx, Multi Msgs Tx, Disallowed Msg": {
			txs:         [][]byte{multiMsgsTxHasDisallowMixedTxBytes},
			expectedTxs: nil,
		},
		"Single Tx, Multi Msgs Tx, Allowed Msg": {
			txs:         [][]byte{constants.Msg_SendAndTransfer_TxBytes},
			expectedTxs: [][]byte{constants.Msg_SendAndTransfer_TxBytes},
		},
		"Multi Tx, Single Msg Tx, Disallowed Msg": {
			txs: [][]byte{
				constants.ValidMsgUpdateMarketPricesTxBytes,
				constants.ValidMsgAddPremiumVotesTxBytes,
			},
			expectedTxs: nil,
		},
		"Multi Tx, Single Msg Tx, Allowed Msg": {
			txs: [][]byte{
				constants.Msg_Send_TxBytes,
				constants.Msg_Send_TxBytes,
			},
			expectedTxs: [][]byte{
				constants.Msg_Send_TxBytes,
				constants.Msg_Send_TxBytes,
			},
		},
		"Multi Tx, Multi Msg Tx, Disallowed Msg": {
			txs: [][]byte{
				multiMsgsTxHasDisallowOnlyTxBytes,
				multiMsgsTxHasDisallowMixedTxBytes,
			},
			expectedTxs: nil,
		},
		"Multi Tx, Multi Msg Tx, Allowed Msg": {
			txs: [][]byte{
				constants.Msg_SendAndTransfer_TxBytes,
				constants.Msg_SendAndTransfer_TxBytes,
			},
			expectedTxs: [][]byte{
				constants.Msg_SendAndTransfer_TxBytes,
				constants.Msg_SendAndTransfer_TxBytes,
			},
		},
		"Multi Tx, Mixed Msgs Tx, Mixed": {
			txs: [][]byte{
				multiMsgsTxHasDisallowMixedTxBytes, // filtered out.
				constants.Msg_SendAndTransfer_TxBytes,
				multiMsgsTxHasDisallowOnlyTxBytes, // filtered out.
				constants.Msg_Send_TxBytes,
				constants.ValidMsgAddPremiumVotesTxBytes, // filtered out.
			},
			expectedTxs: [][]byte{
				constants.Msg_SendAndTransfer_TxBytes,
				constants.Msg_Send_TxBytes,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testApp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			txs := prepare.RemoveDisallowMsgs(ctx, encodingCfg.TxConfig.TxDecoder(), tc.txs)
			require.Equal(t, tc.expectedTxs, txs)
		})
	}
}
