package process_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
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

	// Valid acknowledge bridges tx.
	validAcknowledgeBridgesTx := constants.MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes

	// Valid add funding tx.
	validAddFundingTx := constants.ValidMsgAddPremiumVotesTxBytes

	// Valid update price tx.
	validUpdatePriceTx := constants.ValidMsgUpdateMarketPricesTxBytes

	// Valid "other" single msg tx.
	validSingleMsgOtherTx := constants.Msg_Send_TxBytes

	// Valid "other" multi msgs tx.
	validMultiMsgOtherTx := constants.Msg_SendAndTransfer_TxBytes

	// Invalid update price tx.
	invalidUpdatePriceTx := constants.InvalidMsgUpdateMarketPricesStatelessTxBytes

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
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				invalidUpdatePriceTx, // invalid.
			},
		},
		"Valid txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validMultiMsgOtherTx,  // other txs.
				validSingleMsgOtherTx, // other txs.
				validAcknowledgeBridgesTx,
				validAddFundingTx,
				validUpdatePriceTx,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			_, bridgeKeeper, _, _, _, _, _ := keepertest.BridgeKeepers(t)

			ctx, pricesKeeper, _, indexPriceCache, mockTimeProvider, _, _ := keepertest.PricesKeepers(t)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			indexPriceCache.UpdatePrices(constants.AtTimeTSingleExchangePriceUpdate)

			mockClobKeeper := &mocks.ProcessClobKeeper{}
			mockClobKeeper.On("RecordMevMetricsIsEnabled").Return(true)
			mockClobKeeper.On("RecordMevMetrics", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			handler := process.FullNodeProcessProposalHandler(
				constants.TestEncodingCfg.TxConfig,
				bridgeKeeper,
				mockClobKeeper,
				&mocks.ProcessStakingKeeper{},
				&mocks.ProcessPerpetualKeeper{},
				process.NewDefaultUpdateMarketPriceTxDecoder(pricesKeeper, constants.TestEncodingCfg.TxConfig.TxDecoder()),
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
