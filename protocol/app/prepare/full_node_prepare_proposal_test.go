package prepare_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/app/prepare"
	"github.com/dydxprotocol/v4/mocks"
	keepertest "github.com/dydxprotocol/v4/testutil/keeper"
	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	perpetualtypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestFullNodePrepareProposalHandler test that the full-node PrepareProposal handler always returns
// an empty result.
func TestFullNodePrepareProposalHandler(t *testing.T) {
	tests := map[string]struct {
		txs      [][]byte
		maxBytes int64

		pricesResp    *pricestypes.MsgUpdateMarketPrices
		pricesEncoder sdktypes.TxEncoder

		fundingResp    *perpetualtypes.MsgAddPremiumVotes
		fundingEncoder sdktypes.TxEncoder

		clobResp    *clobtypes.MsgProposedOperations
		clobEncoder sdktypes.TxEncoder
	}{
		"Error: newPrepareProposalTransactions fails": {
			maxBytes: 0, // <= 0 value throws error.
		},

		// Prices related.
		"Error: GetPricesTx returns err": {
			maxBytes: 1,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: failingTxEncoder, // encoder fails and returns err.
		},
		"Error: GetPricesTx returns empty": {
			maxBytes: 1,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: emptyTxEncoder, // encoder returns empty.
		},
		"Error: SetPricesTx returns err": {
			maxBytes: 1,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderTwo, // encoder returns two bytes, which exceeds max.
		},

		// Funding related.
		"Error: GetFundingTx returns err": {
			maxBytes: 2,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: failingTxEncoder, // encoder fails and returns err.
		},
		"Error: GetFundingTx returns empty": {
			maxBytes: 2,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: emptyTxEncoder, // encoder returns empty.
		},
		"Error: SetFundingTx returns err": {
			maxBytes: 1, // only upto 1 byte, not enough space for funding tx bytes.

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne, // takes up 1 byte.

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne, // takes up another 1 byte, so exceeds max.
		},

		// Operations related.
		"Error: GetOperationsTx returns err": {
			maxBytes: 3,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: failingTxEncoder, // encoder fails and returns err.
		},
		"Error: GetOperationsTx returns empty": {
			maxBytes: 3,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: emptyTxEncoder, // encoder returns empty.
		},
		"Error: SetOperationsTx returns err": {
			maxBytes: 2, // only upto 2 bytes, not enough space for the operation tx.

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne, // takes up 1 byte.

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne, // takes up another 1 byte.

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderOne, // takes up another 1, so exceeds max.
		},

		// "Others" related.
		"Error: Others takes up all space, no space for order tx": {
			maxBytes: 12,
			txs:      [][]byte{{9}}, // "Other" order takes up 1 byte.

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour, // takes up 4 byte.

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour, // takes up 4 bytes

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour, // takes another 4 bytes, but exceeds max.
		},
		"Error: AddOtherTxs return error": {
			maxBytes: 13,
			txs:      [][]byte{{}},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,
		},
		"Error: AddOtherTxs (additional) return error": {
			maxBytes: 15,
			txs:      [][]byte{{9, 8}, {9}, {}, {}},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,
		},
		"Valid: Not all Others than can fit": {
			maxBytes: 13,
			txs:      [][]byte{{9}, {9}, {9}},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,
		},
		"Valid: Additional Others fit": {
			maxBytes: 15,
			txs:      [][]byte{{9}, {9, 8}, {9, 8, 7}},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockContextHelper := mocks.ContextHelper{}
			mockContextHelper.On("Height", mock.Anything).Return(int64(1))

			mockPricesKeeper := mocks.PreparePricesKeeper{}
			mockPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).
				Return(tc.pricesResp)

			mockPerpKeeper := mocks.PreparePerpetualsKeeper{}
			mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything, mock.Anything).
				Return(tc.fundingResp)

			mockClobKeeper := mocks.PrepareClobKeeper{}
			mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
				Return(tc.clobResp)

			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

			handler := prepare.FullNodePrepareProposalHandler()

			req := abci.RequestPrepareProposal{
				Txs:        tc.txs,
				MaxTxBytes: tc.maxBytes,
			}

			response := handler(ctx, req)
			require.Equal(t, [][]byte{}, response.Txs)
		})
	}
}
