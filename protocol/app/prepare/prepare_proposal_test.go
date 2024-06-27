package prepare_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/prepare"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/encoding"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perpetualtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	ctx = sdktypes.Context{}

	failingTxEncoder = func(tx sdktypes.Tx) ([]byte, error) {
		return nil, errors.New("encoder failed")
	}
	emptyTxEncoder = func(tx sdktypes.Tx) ([]byte, error) {
		return []byte{}, nil
	}
	passingTxEncoderOne = func(tx sdktypes.Tx) ([]byte, error) {
		return []byte{1}, nil
	}
	passingTxEncoderTwo = func(tx sdktypes.Tx) ([]byte, error) {
		return []byte{1, 2}, nil
	}
	passingTxEncoderFour = func(tx sdktypes.Tx) ([]byte, error) {
		return []byte{1, 2, 3, 4}, nil
	}
)

func TestPrepareProposalHandler(t *testing.T) {
	msgSendTxBytesLen := int64(len(constants.Msg_Send_TxBytes))
	msgSendAndTransferTxBytesLen := int64(len(constants.Msg_SendAndTransfer_TxBytes))

	tests := map[string]struct {
		txs      [][]byte
		maxBytes int64

		fundingResp    *perpetualtypes.MsgAddPremiumVotes
		fundingEncoder sdktypes.TxEncoder

		clobResp    *clobtypes.MsgProposedOperations
		clobEncoder sdktypes.TxEncoder

		expectedTxs [][]byte
	}{
		"Error: newPrepareProposalTransactions fails": {
			maxBytes:    0,          // <= 0 value throws error.
			expectedTxs: [][]byte{}, // error returns empty result.
		},

		// Funding related.
		"Error: GetFundingTx returns err": {
			maxBytes: 1,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: failingTxEncoder, // encoder fails and returns err.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: GetFundingTx returns empty": {
			maxBytes: 1,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: emptyTxEncoder, // encoder returns empty.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: SetFundingTx returns err": {
			maxBytes: 0, // only upto 1 byte, not enough space for funding tx bytes.

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne, // takes up another 1 byte, so exceeds max.

			expectedTxs: [][]byte{}, // error returns empty result.
		},

		// "Others" related.
		"Error: AddOtherTxs return error": {
			maxBytes: 13,
			txs:      [][]byte{{}},

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: AddOtherTxs (additional) return error": {
			maxBytes: 15,
			txs:      [][]byte{{9, 8}, {9}, {}, {}},

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Valid: Not all Others than can fit": {
			maxBytes: int64(12) + msgSendTxBytesLen + 1,
			txs: [][]byte{
				constants.Msg_Send_TxBytes,
				constants.Msg_Send_TxBytes, // not included due to maxBytes.
				constants.Msg_Send_TxBytes, // not included due to maxBytes.
			},

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			expectedTxs: [][]byte{
				{1, 2, 3, 4},               // order.
				constants.Msg_Send_TxBytes, // others.
				{1, 2, 3, 4},               // funding.
				{1, 2, 3, 4},               // prices.
			},
		},
		"Valid: Additional Others fit": {
			maxBytes: int64(12) + msgSendTxBytesLen + msgSendAndTransferTxBytesLen,
			txs: [][]byte{
				constants.Msg_Send_TxBytes,
				constants.Msg_SendAndTransfer_TxBytes,
				constants.Msg_Send_TxBytes, // not included due to maxBytes.
			},

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			expectedTxs: [][]byte{
				{1, 2, 3, 4},                          // order.
				constants.Msg_Send_TxBytes,            // others.
				constants.Msg_SendAndTransfer_TxBytes, // additional others.
				{1, 2, 3, 4},                          // funding.
				{1, 2, 3, 4},                          // prices.
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockTxConfig := createMockTxConfig(
				nil,
				[]sdktypes.TxEncoder{
					tc.fundingEncoder,
					tc.clobEncoder,
				},
			)

			mockPerpKeeper := mocks.PreparePerpetualsKeeper{}
			mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything).
				Return(tc.fundingResp)

			mockClobKeeper := mocks.PrepareClobKeeper{}
			mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
				Return(tc.clobResp)

			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

			handler := prepare.PrepareProposalHandler(
				mockTxConfig,
				&mockClobKeeper,
				&mockPerpKeeper,
			)

			req := abci.RequestPrepareProposal{
				Txs:        tc.txs,
				MaxTxBytes: tc.maxBytes,
			}

			response, err := handler(ctx, &req)
			require.NoError(t, err)
			require.Equal(t, tc.expectedTxs, response.Txs)
		})
	}
}

func TestPrepareProposalHandler_OtherTxs(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()

	tests := map[string]struct {
		txs [][]byte

		expectedTxs [][]byte
	}{
		"Valid: all others txs contain disallow msgs": {
			txs: [][]byte{
				multiMsgsTxHasDisallowOnlyTxBytes,  // filtered out.
				multiMsgsTxHasDisallowMixedTxBytes, // filtered out.
			},
			expectedTxs: [][]byte{
				constants.ValidEmptyMsgProposedOperationsTxBytes, // order.
				// no other txs.
				constants.ValidMsgAddPremiumVotesTxBytes,    // funding.
				constants.ValidMsgUpdateMarketPricesTxBytes, // prices.
			},
		},
		"Valid: some others txs contain disallow msgs": {
			txs: [][]byte{
				multiMsgsTxHasDisallowMixedTxBytes, // filtered out.
				constants.Msg_SendAndTransfer_TxBytes,
				multiMsgsTxHasDisallowOnlyTxBytes, // filtered out.
				constants.Msg_Send_TxBytes,
				constants.ValidMsgAddPremiumVotesTxBytes, // filtered out.
			},
			expectedTxs: [][]byte{
				constants.ValidEmptyMsgProposedOperationsTxBytes, // order.
				constants.Msg_SendAndTransfer_TxBytes,            // others.
				constants.Msg_Send_TxBytes,                       // others.
				constants.ValidMsgAddPremiumVotesTxBytes,         // funding.
				constants.ValidMsgUpdateMarketPricesTxBytes,      // prices.
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockPricesKeeper := mocks.PreparePricesKeeper{}
			mockPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).
				Return(constants.ValidMsgUpdateMarketPrices)

			mockPerpKeeper := mocks.PreparePerpetualsKeeper{}
			mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything).
				Return(constants.ValidMsgAddPremiumVotes)

			mockClobKeeper := mocks.PrepareClobKeeper{}
			mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
				Return(constants.ValidEmptyMsgProposedOperations)

			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

			handler := prepare.PrepareProposalHandler(
				encodingCfg.TxConfig,
				&mockClobKeeper,
				&mockPerpKeeper,
			)

			req := abci.RequestPrepareProposal{
				Txs:        tc.txs,
				MaxTxBytes: 100_000, // something large.
			}

			response, err := handler(ctx, &req)
			require.NoError(t, err)
			require.Equal(t, tc.expectedTxs, response.Txs)
		})
	}
}

func TestGetAddPremiumVotesTx(t *testing.T) {
	tests := map[string]struct {
		keeperResp *perpetualtypes.MsgAddPremiumVotes
		txEncoder  sdktypes.TxEncoder

		expectedTx       []byte
		expectedNumVotes int
		expectedErr      error
	}{
		"nil message fails": {
			keeperResp: nil,

			expectedErr: fmt.Errorf("MsgAddPremiumVotes cannot be nil"),
		},
		"empty message": {
			keeperResp: &perpetualtypes.MsgAddPremiumVotes{}, // empty
			txEncoder:  passingTxEncoderOne,

			expectedTx:       []byte{1},
			expectedNumVotes: 0,
		},
		"empty tx": {
			keeperResp: &perpetualtypes.MsgAddPremiumVotes{},
			txEncoder:  emptyTxEncoder, // returns empty tx.

			expectedErr: fmt.Errorf("Invalid tx: []"),
		},
		"valid message, but encoding fails": {
			keeperResp: &perpetualtypes.MsgAddPremiumVotes{}, // empty
			txEncoder:  failingTxEncoder,

			expectedErr: fmt.Errorf("encoder failed"),
		},
		"valid message": {
			keeperResp: &perpetualtypes.MsgAddPremiumVotes{
				Votes: []perpetualtypes.FundingPremium{{}, {}},
			},
			txEncoder: passingTxEncoderOne,

			expectedTx:       []byte{1},
			expectedNumVotes: 2,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockTxConfig := createMockTxConfig(nil, []sdktypes.TxEncoder{tc.txEncoder})
			mockPerpKeeper := mocks.PreparePerpetualsKeeper{}
			mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything).
				Return(tc.keeperResp)

			resp, err := prepare.GetAddPremiumVotesTx(ctx, mockTxConfig, &mockPerpKeeper)
			if tc.expectedErr != nil {
				require.Equal(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTx, resp.Tx)
			require.Equal(t, tc.expectedNumVotes, resp.NumVotes)
		})
	}
}

func TestGetProposedOperationsTx(t *testing.T) {
	tests := map[string]struct {
		keeperResp *clobtypes.MsgProposedOperations
		txEncoder  sdktypes.TxEncoder

		expectedTx               []byte
		expectedNumPlaceOrders   int
		expectedNumMatchedOrders int
		expectedErr              error
	}{
		"nil message fails": {
			keeperResp: nil,

			expectedErr: fmt.Errorf("MsgProposedOperations cannot be nil"),
		},
		"empty message": {
			keeperResp: &clobtypes.MsgProposedOperations{}, // empty
			txEncoder:  passingTxEncoderOne,

			expectedTx:               []byte{1},
			expectedNumMatchedOrders: 0,
			expectedNumPlaceOrders:   0,
		},
		"empty tx": {
			keeperResp: &clobtypes.MsgProposedOperations{},
			txEncoder:  emptyTxEncoder, // returns empty tx.

			expectedErr: fmt.Errorf("Invalid tx: []"),
		},
		"valid message, but encoding fails": {
			keeperResp: &clobtypes.MsgProposedOperations{}, // empty
			txEncoder:  failingTxEncoder,

			expectedErr: fmt.Errorf("encoder failed"),
		},
		"valid message": {
			keeperResp: &clobtypes.MsgProposedOperations{
				OperationsQueue: []clobtypes.OperationRaw{{}, {}},
			},
			txEncoder: passingTxEncoderOne,

			expectedTx:               []byte{1},
			expectedNumPlaceOrders:   2,
			expectedNumMatchedOrders: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockTxConfig := createMockTxConfig(nil, []sdktypes.TxEncoder{tc.txEncoder})
			mockClobKeeper := mocks.PrepareClobKeeper{}
			mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).Return(tc.keeperResp)

			resp, err := prepare.GetProposedOperationsTx(ctx, mockTxConfig, &mockClobKeeper)
			if tc.expectedErr != nil {
				require.Equal(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTx, resp.Tx)
		})
	}
}

func TestEncodeMsgsIntoTxBytes(t *testing.T) {
	tests := map[string]struct {
		setMsgErr error
		txEncoder sdktypes.TxEncoder

		expectedTx  []byte
		expectedErr error
	}{
		"set message fails": {
			setMsgErr:   errors.New("unexpected SetMsgs error"),
			expectedErr: errors.New("unexpected SetMsgs error"),
		},
		"tx encoder fails": {
			setMsgErr:   nil,
			txEncoder:   failingTxEncoder,
			expectedErr: errors.New("encoder failed"),
		},
		"valid": {
			setMsgErr:  nil,
			txEncoder:  passingTxEncoderOne,
			expectedTx: []byte{1},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockTxConfig := createMockTxConfig(tc.setMsgErr, []sdktypes.TxEncoder{tc.txEncoder})

			tx, err := prepare.EncodeMsgsIntoTxBytes(mockTxConfig, &clobtypes.MsgProposedOperations{})

			if tc.expectedErr != nil {
				require.Equal(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTx, tx)
		})
	}
}

func createMockTxConfig(setMsgsError error, allTxEncoders []sdktypes.TxEncoder) *mocks.TxConfig {
	mockTxConfig := mocks.TxConfig{}
	mockTxBuilder := mocks.TxBuilder{}

	mockTxConfig.On("NewTxBuilder").Return(&mockTxBuilder)
	mockTxBuilder.On("SetMsgs", mock.Anything).Return(setMsgsError)
	mockTxBuilder.On("GetTx").Return(nil) // doesn't really matter, since encoder is mocked.

	for _, txEncoder := range allTxEncoders {
		mockTxConfig.On("TxEncoder").Return(txEncoder).Once()
	}

	mockTxConfig.On("TxDecoder").Return(encoding.GetTestEncodingCfg().TxConfig.TxDecoder())

	return &mockTxConfig
}
