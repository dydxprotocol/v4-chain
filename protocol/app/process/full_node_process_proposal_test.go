package process_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/dydxprotocol/v4/app/process"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
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
				validAddFundingTx,
				invalidUpdatePriceTx, // invalid.
			},
		},
		"Valid txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validMultiMsgOtherTx,  // other txs.
				validSingleMsgOtherTx, // other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			mockContextHelper := mocks.ContextHelper{}
			mockContextHelper.On("Height", mock.Anything).Return(int64(2))
			ctx, pricesKeeper, _, indexPriceCache, _, mockTimeProvider := keepertest.PricesKeepers(t)
			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)
			indexPriceCache.UpdatePrices(constants.AtTimeTSingleExchangePriceUpdate)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			handler := process.FullNodeProcessProposalHandler()
			req := abci.RequestProcessProposal{Txs: tc.txsBytes}

			// Run.
			resp := handler(ctx, req)

			// Validate.
			require.Equal(t, acceptResponse, resp)
		})
	}
}
