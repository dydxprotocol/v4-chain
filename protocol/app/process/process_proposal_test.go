package process_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/dydxprotocol/v4/app/process"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	testmsgs "github.com/dydxprotocol/v4/testutil/msgs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessProposalHandler_Error(t *testing.T) {
	acceptResponse := abci.ResponseProcessProposal{
		Status: abci.ResponseProcessProposal_ACCEPT,
	}
	rejectResponse := abci.ResponseProcessProposal{
		Status: abci.ResponseProcessProposal_REJECT,
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
		"Reject: decode fails": {
			txsBytes:         [][]byte{{1}, {2}},
			expectedResponse: rejectResponse,
		},
		"Reject: invalid price tx": {
			txsBytes: [][]byte{
				validOperationsTx,
				validAddFundingTx,
				invalidUpdatePriceTx, // invalid.
			},
			expectedResponse: rejectResponse,
		},
		"Error: place order type is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				constants.Msg_PlaceOrder_TxBtyes, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Error: cancel order type is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				constants.Msg_CancelOrder_TxBtyes, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Error: app-injected msg type is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				validUpdatePriceTx, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Error: internal msg type is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				testmsgs.MsgSoftwareUpgradeTxBytes, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Error: unsupported msg type is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				testmsgs.GovBetaMsgSubmitProposalTxBytes, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Error: nested msg type with unsupported inner is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				testmsgs.MsgSubmitProposalWithUnsupportedInnerTxBytes, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Error: nested msg type with app-injected inner is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				testmsgs.MsgSubmitProposalWithAppInjectedInnerTxBytes, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Error: nested msg type with double-nested inner is not allowed": {
			txsBytes: [][]byte{
				validOperationsTx,
				testmsgs.MsgSubmitProposalWithDoubleNestedInnerTxBytes, // invalid other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: rejectResponse,
		},
		"Accept: Valid txs": {
			txsBytes: [][]byte{
				validOperationsTx,
				validMultiMsgOtherTx,  // other txs.
				validSingleMsgOtherTx, // other txs.
				validAddFundingTx,
				validUpdatePriceTx,
			},
			expectedResponse: acceptResponse,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			mockContextHelper := mocks.ContextHelper{}
			mockContextHelper.On("Height", mock.Anything).Return(int64(2))
			ctx, pricesKeeper, _, indexPriceCache, mockTimeProvider := keepertest.PricesKeepers(t)
			keepertest.CreateTestMarketsAndExchangeFeeds(t, ctx, pricesKeeper)
			indexPriceCache.UpdatePrices(constants.AtTimeTSingleExchangePriceUpdate)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			handler := process.ProcessProposalHandler(
				&mockContextHelper,
				constants.TestEncodingCfg.TxConfig,
				pricesKeeper,
			)
			req := abci.RequestProcessProposal{Txs: tc.txsBytes}

			// Run.
			resp := handler(ctx, req)

			// Validate.
			require.Equal(t, tc.expectedResponse, resp)
		})
	}
}
