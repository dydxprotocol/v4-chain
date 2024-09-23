package process_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/process"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestFullNodeProcessProposalHandler validates that the TestFullNodeProcessProposalHandler
// always returns ResponseProcessProposal_ACCEPT.
func TestFullNodeProcessProposalHandler(t *testing.T) {
	acceptResponse := abci.ResponseProcessProposal{
		Status: abci.ResponseProcessProposal_ACCEPT,
	}

	// Valid operations tx.
	validOperationsTx := constants.ValidEmptyMsgProposedOperationsTxBytes

	// Valid add funding tx.
	validAddFundingTx := constants.ValidMsgAddPremiumVotesTxBytes

	// Valid "other" single msg tx.
	validSingleMsgOtherTx := constants.Msg_Send_TxBytes

	// Valid "other" multi msgs tx.
	validMultiMsgOtherTx := constants.Msg_SendAndTransfer_TxBytes

	tests := map[string]struct {
		txsBytes [][]byte

		expectedResponse abci.ResponseProcessProposal
	}{
		"Bad txsBytes": {
			txsBytes: [][]byte{{1}, {2}},
		},
		"Invalid transactions": {
			txsBytes: [][]byte{
				validOperationsTx,
				validAddFundingTx,
			},
		},
		"Valid txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validMultiMsgOtherTx,  // other txs.
				validSingleMsgOtherTx, // other txs.
				validAddFundingTx,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, pricesKeeper, _, daemonPriceCache, _, mockTimeProvider := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			daemonPriceCache.UpdatePrices(constants.AtTimeTSingleExchangePriceUpdate)

			mockClobKeeper := &mocks.ProcessClobKeeper{}
			mockClobKeeper.On("RecordMevMetricsIsEnabled").Return(true)
			mockClobKeeper.On("RecordMevMetrics", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			handler := process.FullNodeProcessProposalHandler(
				constants.TestEncodingCfg.TxConfig,
				mockClobKeeper,
				&mocks.ProcessPerpetualKeeper{},
				pricesKeeper,
			)
			req := abci.RequestProcessProposal{Txs: tc.txsBytes}

			// Run.
			resp, err := handler(ctx, &req)
			require.NoError(t, err)

			// Validate.
			require.Equal(t, acceptResponse, *resp)
		})
	}
}
