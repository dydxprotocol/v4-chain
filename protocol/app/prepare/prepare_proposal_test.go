package prepare_test

import (
	"errors"
	"fmt"
	"testing"

	"math/big"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/slinky/abci/strategies/aggregator"
	aggregatormock "github.com/dydxprotocol/slinky/abci/strategies/aggregator/mocks"
	"github.com/dydxprotocol/slinky/abci/strategies/codec"
	strategymock "github.com/dydxprotocol/slinky/abci/strategies/currencypair/mocks"
	slinkytestutils "github.com/dydxprotocol/slinky/abci/testutils"
	vetypes "github.com/dydxprotocol/slinky/abci/ve/types"
	oracletypes "github.com/dydxprotocol/slinky/pkg/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare/prices"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	clobtestutil "github.com/dydxprotocol/v4-chain/protocol/testutil/clob"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	ctx = sdk.Context{}

	failingTxEncoder = func(tx sdk.Tx) ([]byte, error) {
		return nil, errors.New("encoder failed")
	}
	emptyTxEncoder = func(tx sdk.Tx) ([]byte, error) {
		return []byte{}, nil
	}
	passingTxEncoderOne = func(tx sdk.Tx) ([]byte, error) {
		return []byte{1}, nil
	}
	passingTxEncoderTwo = func(tx sdk.Tx) ([]byte, error) {
		return []byte{1, 2}, nil
	}
	passingTxEncoderFour = func(tx sdk.Tx) ([]byte, error) {
		return []byte{1, 2, 3, 4}, nil
	}
)

func TestPrepareProposalHandler(t *testing.T) {
	msgSendTxBytesLen := int64(len(constants.Msg_Send_TxBytes))
	msgSendAndTransferTxBytesLen := int64(len(constants.Msg_SendAndTransfer_TxBytes))

	tests := map[string]struct {
		txs      [][]byte
		maxBytes int64

		pricesResp    *pricestypes.MsgUpdateMarketPrices
		pricesEncoder sdk.TxEncoder

		fundingResp    *perpetualtypes.MsgAddPremiumVotes
		fundingEncoder sdk.TxEncoder

		clobResp    *clobtypes.MsgProposedOperations
		clobEncoder sdk.TxEncoder

		bridgeResp    *bridgetypes.MsgAcknowledgeBridges
		bridgeEncoder sdk.TxEncoder

		expectedTxs [][]byte
	}{
		"Error: newPrepareProposalTransactions fails": {
			maxBytes:    0,          // <= 0 value throws error.
			expectedTxs: [][]byte{}, // error returns empty result.
		},

		// Prices related.
		"Error: GetPricesTx returns err": {
			maxBytes: 1,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: failingTxEncoder, // encoder fails and returns err.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: GetPricesTx returns empty": {
			maxBytes: 1,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: emptyTxEncoder, // encoder returns empty.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: SetPricesTx returns err": {
			maxBytes: 1,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderTwo, // encoder returns two bytes, which exceeds max.

			expectedTxs: [][]byte{}, // error returns empty result.
		},

		// Funding related.
		"Error: GetFundingTx returns err": {
			maxBytes: 2,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: failingTxEncoder, // encoder fails and returns err.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: GetFundingTx returns empty": {
			maxBytes: 2,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: emptyTxEncoder, // encoder returns empty.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: SetFundingTx returns err": {
			maxBytes: 1, // only upto 1 byte, not enough space for funding tx bytes.

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne, // takes up 1 byte.

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne, // takes up another 1 byte, so exceeds max.

			expectedTxs: [][]byte{}, // error returns empty result.
		},

		// Bridge related.
		"Error: GetAcknowledgeBridgesTx returns err": {
			maxBytes: 3,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: failingTxEncoder, // encoder fails and returns err.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: GetAcknowledgeBridgesTx returns empty": {
			maxBytes: 3,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: emptyTxEncoder, // encoder returns empty.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: SetAcknowledgeBridgesTx returns err": {
			maxBytes: 2,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne, // takes up 1 byte

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne, // takes up 1 byte

			bridgeResp:    constants.MsgAcknowledgeBridges_Id0_Height0,
			bridgeEncoder: passingTxEncoderOne, // takes up another 1 byte, so exceeds max.

			expectedTxs: [][]byte{}, // error returns empty result.
		},

		// Operations related.
		"Error: GetOperationsTx returns err": {
			maxBytes: 4,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: passingTxEncoderOne,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: failingTxEncoder, // encoder fails and returns err.

			expectedTxs: [][]byte{}, // error returns empty result.
		},
		"Error: GetOperationsTx returns empty": {
			maxBytes: 4,

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: passingTxEncoderOne,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: emptyTxEncoder, // encoder returns empty.

			// GetProposedOperationsTx errors; the handler falls back to an empty operations tx, so the
			// proposal keeps the fixed bridge/funding/prices txs instead of collapsing to empty.
			expectedTxs: [][]byte{{1}, {1}, {1}},
		},
		"Error: SetOperationsTx returns err": {
			maxBytes: 3, // only upto 3 bytes, not enough space for the order tx.

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderOne, // takes up 1 byte.

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderOne, // takes up another 1 byte.

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: passingTxEncoderOne, // takes up another 1 byte.

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderOne, // takes up another 1, so exceeds max.

			expectedTxs: [][]byte{}, // error returns empty result.
		},

		// "Others" related.
		"Error: AddOtherTxs return error": {
			maxBytes: 17,
			txs:      [][]byte{{}},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			// AddOtherTxs failure no longer empties the proposal; the fixed txs are kept.
			expectedTxs: [][]byte{{1, 2, 3, 4}, {1, 2, 3, 4}, {1, 2, 3, 4}, {1, 2, 3, 4}},
		},
		"Error: AddOtherTxs (additional) return error": {
			maxBytes: 19,
			txs:      [][]byte{{9, 8}, {9}, {}, {}},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			// The additional AddOtherTxs failure no longer empties the proposal; the fixed txs are kept.
			expectedTxs: [][]byte{{1, 2, 3, 4}, {1, 2, 3, 4}, {1, 2, 3, 4}, {1, 2, 3, 4}},
		},
		"Valid: Not all Others than can fit": {
			maxBytes: int64(16) + msgSendTxBytesLen + 1,
			txs: [][]byte{
				constants.Msg_Send_TxBytes,
				constants.Msg_Send_TxBytes, // not included due to maxBytes.
				constants.Msg_Send_TxBytes, // not included due to maxBytes.
			},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			expectedTxs: [][]byte{
				{1, 2, 3, 4},               // order.
				constants.Msg_Send_TxBytes, // others.
				{1, 2, 3, 4},               // bridge.
				{1, 2, 3, 4},               // funding.
				{1, 2, 3, 4},               // prices.
			},
		},
		"Valid: Additional Others fit": {
			maxBytes: int64(16) + msgSendTxBytesLen + msgSendAndTransferTxBytesLen,
			txs: [][]byte{
				constants.Msg_Send_TxBytes,
				constants.Msg_SendAndTransfer_TxBytes,
				constants.Msg_Send_TxBytes, // not included due to maxBytes.
			},

			pricesResp:    &pricestypes.MsgUpdateMarketPrices{},
			pricesEncoder: passingTxEncoderFour,

			fundingResp:    &perpetualtypes.MsgAddPremiumVotes{},
			fundingEncoder: passingTxEncoderFour,

			bridgeResp:    &bridgetypes.MsgAcknowledgeBridges{},
			bridgeEncoder: passingTxEncoderFour,

			clobResp:    &clobtypes.MsgProposedOperations{},
			clobEncoder: passingTxEncoderFour,

			expectedTxs: [][]byte{
				{1, 2, 3, 4},                          // order.
				constants.Msg_Send_TxBytes,            // others.
				constants.Msg_SendAndTransfer_TxBytes, // additional others.
				{1, 2, 3, 4},                          // bridge.
				{1, 2, 3, 4},                          // funding.
				{1, 2, 3, 4},                          // prices.
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockTxConfig := createMockTxConfig(
				nil,
				[]sdk.TxEncoder{
					tc.pricesEncoder,
					tc.fundingEncoder,
					tc.bridgeEncoder,
					tc.clobEncoder,
					// GetProposedOperationsTx re-encodes an empty ops tx when the full ops tx
					// doesn't fit within the (mocked, fixed-size) available bytes.
					tc.clobEncoder,
					// On a GetProposedOperationsTx error, the handler falls back to encoding an empty
					// operations tx rather than returning an empty proposal.
					tc.clobEncoder,
				},
			)

			mockPricesKeeper := mocks.PreparePricesKeeper{}
			mockPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).
				Return(tc.pricesResp)

			mockPerpKeeper := mocks.PreparePerpetualsKeeper{}
			mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything).
				Return(tc.fundingResp)

			mockBridgeKeeper := mocks.PrepareBridgeKeeper{}
			mockBridgeKeeper.On("GetAcknowledgeBridges", mock.Anything, mock.Anything).
				Return(tc.bridgeResp)

			mockClobKeeper := mocks.PrepareClobKeeper{}
			mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
				Return(tc.clobResp)

			ctx, _, _, _, _, _, _ := keepertest.PricesKeepers(t)

			handler := prepare.PrepareProposalHandler(
				mockTxConfig,
				&mockBridgeKeeper,
				&mockClobKeeper,
				&mockPerpKeeper,
				prices.NewDefaultPriceUpdateGenerator(&mockPricesKeeper),
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
				constants.MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes, // bridge.
				constants.ValidMsgAddPremiumVotesTxBytes,               // funding.
				constants.ValidMsgUpdateMarketPricesTxBytes,            // prices.
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
				constants.ValidEmptyMsgProposedOperationsTxBytes,       // order.
				constants.Msg_SendAndTransfer_TxBytes,                  // others.
				constants.Msg_Send_TxBytes,                             // others.
				constants.MsgAcknowledgeBridges_Ids0_1_Height0_TxBytes, // bridge.
				constants.ValidMsgAddPremiumVotesTxBytes,               // funding.
				constants.ValidMsgUpdateMarketPricesTxBytes,            // prices.
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

			mockBridgeKeeper := mocks.PrepareBridgeKeeper{}
			mockBridgeKeeper.On("GetAcknowledgeBridges", mock.Anything, mock.Anything).
				Return(constants.MsgAcknowledgeBridges_Ids0_1_Height0)

			ctx, _, _, _, _, _, _ := keepertest.PricesKeepers(t)

			handler := prepare.PrepareProposalHandler(
				encodingCfg.TxConfig,
				&mockBridgeKeeper,
				&mockClobKeeper,
				&mockPerpKeeper,
				prices.NewDefaultPriceUpdateGenerator(&mockPricesKeeper),
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

// setUpFixedGroupMocks wires up mocks for the "fixed" tx groups (prices, funding, bridge) so that
// their encoded sizes are deterministic and can be measured ahead of time via
// `EncodeMarketPriceUpdates`/`GetAddPremiumVotesTx`/`GetAcknowledgeBridgesTx`.
func setUpFixedGroupMocks() (*mocks.PreparePricesKeeper, *mocks.PreparePerpetualsKeeper, *mocks.PrepareBridgeKeeper) {
	mockPricesKeeper := &mocks.PreparePricesKeeper{}
	mockPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).
		Return(constants.ValidMsgUpdateMarketPrices)

	mockPerpKeeper := &mocks.PreparePerpetualsKeeper{}
	mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything).
		Return(constants.ValidMsgAddPremiumVotes)

	mockBridgeKeeper := &mocks.PrepareBridgeKeeper{}
	mockBridgeKeeper.On("GetAcknowledgeBridges", mock.Anything, mock.Anything).
		Return(constants.MsgAcknowledgeBridges_Ids0_1_Height0)

	return mockPricesKeeper, mockPerpKeeper, mockBridgeKeeper
}

// fixedGroupSize returns the combined encoded size, using the real tx encoder, of the prices,
// funding, and bridge txs that `PrepareProposalHandler` always includes ahead of the operations tx.
func fixedGroupSize(
	t *testing.T,
	txConfig client.TxConfig,
	mockPerpKeeper *mocks.PreparePerpetualsKeeper,
	mockBridgeKeeper *mocks.PrepareBridgeKeeper,
) uint64 {
	pricesTxResp, err := prepare.EncodeMarketPriceUpdates(txConfig, constants.ValidMsgUpdateMarketPrices)
	require.NoError(t, err)

	fundingTxResp, err := prepare.GetAddPremiumVotesTx(sdk.Context{}, txConfig, mockPerpKeeper)
	require.NoError(t, err)

	bridgeTxResp, err := prepare.GetAcknowledgeBridgesTx(sdk.Context{}, txConfig, mockBridgeKeeper)
	require.NoError(t, err)

	return uint64(len(pricesTxResp.Tx)) + uint64(len(fundingTxResp.Tx)) + uint64(len(bridgeTxResp.Tx))
}

func TestPrepareProposalHandler_OperationsTxTruncatedToPrefix(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()
	txConfig := encodingCfg.TxConfig

	const numOps = 50
	ops := buildOperationsQueue(numOps)

	mockPricesKeeper, mockPerpKeeper, mockBridgeKeeper := setUpFixedGroupMocks()
	fixedSize := fixedGroupSize(t, txConfig, mockPerpKeeper, mockBridgeKeeper)

	emptyOpsTx, err := prepare.EncodeMsgsIntoTxBytes(txConfig, &clobtypes.MsgProposedOperations{})
	require.NoError(t, err)
	fullOpsTx, err := prepare.EncodeMsgsIntoTxBytes(txConfig, &clobtypes.MsgProposedOperations{OperationsQueue: ops})
	require.NoError(t, err)
	require.Less(t, uint64(len(emptyOpsTx)), uint64(len(fullOpsTx)))

	// Target a specific, non-trivial prefix length so the truncation boundary is predictable, with
	// enough slack above it to absorb the cometbft per-tx proto-encoding overhead that isn't
	// accounted for by the raw byte-length budget `PrepareProposalTxs` tracks internally.
	const targetNumOps = numOps / 2
	targetOpsTx, err := prepare.EncodeMsgsIntoTxBytes(
		txConfig,
		&clobtypes.MsgProposedOperations{OperationsQueue: ops[:targetNumOps]},
	)
	require.NoError(t, err)
	nextOpsTx, err := prepare.EncodeMsgsIntoTxBytes(
		txConfig,
		&clobtypes.MsgProposedOperations{OperationsQueue: ops[:targetNumOps+1]},
	)
	require.NoError(t, err)

	const slack = uint64(60)
	require.Greater(
		t,
		uint64(len(nextOpsTx))-uint64(len(targetOpsTx)),
		slack,
		"slack must be smaller than the per-operation size increment for the truncation boundary "+
			"to land exactly on targetNumOps",
	)
	opsBudget := uint64(len(targetOpsTx)) + slack
	maxTxBytes := int64(fixedSize + opsBudget) //nolint:gosec

	mockClobKeeper := &mocks.PrepareClobKeeper{}
	mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
		Return(&clobtypes.MsgProposedOperations{OperationsQueue: ops})

	ctx, _, _, _, _, _, _ := keepertest.PricesKeepers(t)

	handler := prepare.PrepareProposalHandler(
		txConfig,
		mockBridgeKeeper,
		mockClobKeeper,
		mockPerpKeeper,
		prices.NewDefaultPriceUpdateGenerator(mockPricesKeeper),
	)

	req := abci.RequestPrepareProposal{
		Txs:        [][]byte{},
		MaxTxBytes: maxTxBytes,
	}

	resp, err := handler(ctx, &req)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(resp.Txs), 4)

	require.NoError(t, cmttypes.ToTxs(resp.Txs).Validate(maxTxBytes))

	opsTx, err := txConfig.TxDecoder()(resp.Txs[0])
	require.NoError(t, err)
	msgs := opsTx.GetMsgs()
	require.Len(t, msgs, 1)
	proposedOps, ok := msgs[0].(*clobtypes.MsgProposedOperations)
	require.True(t, ok)

	require.NotEmpty(t, proposedOps.OperationsQueue)
	require.Less(t, len(proposedOps.OperationsQueue), numOps)
	require.Equal(t, ops[:len(proposedOps.OperationsQueue)], proposedOps.OperationsQueue)
	require.Equal(t, targetNumOps, len(proposedOps.OperationsQueue))
}

func TestPrepareProposalHandler_OperationsTxFallsBackToEmpty(t *testing.T) {
	encodingCfg := encoding.GetTestEncodingCfg()
	txConfig := encodingCfg.TxConfig

	const numOps = 50
	ops := buildOperationsQueue(numOps)

	mockPricesKeeper, mockPerpKeeper, mockBridgeKeeper := setUpFixedGroupMocks()
	fixedSize := fixedGroupSize(t, txConfig, mockPerpKeeper, mockBridgeKeeper)

	emptyOpsTx, err := prepare.EncodeMsgsIntoTxBytes(txConfig, &clobtypes.MsgProposedOperations{})
	require.NoError(t, err)
	oneOpTx, err := prepare.EncodeMsgsIntoTxBytes(
		txConfig,
		&clobtypes.MsgProposedOperations{OperationsQueue: ops[:1]},
	)
	require.NoError(t, err)

	emptySize := uint64(len(emptyOpsTx))
	oneOpSize := uint64(len(oneOpTx))
	require.Greater(t, oneOpSize, emptySize)

	// Budget for operations sits strictly between the empty-tx size and the smallest non-empty
	// (1-op) tx size, so only the empty ops tx fits.
	opsBudget := emptySize + (oneOpSize-emptySize)/2
	require.Greater(t, opsBudget, emptySize)
	require.Less(t, opsBudget, oneOpSize)

	maxTxBytes := int64(fixedSize + opsBudget) //nolint:gosec

	mockClobKeeper := &mocks.PrepareClobKeeper{}
	mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
		Return(&clobtypes.MsgProposedOperations{OperationsQueue: ops})

	ctx, _, _, _, _, _, _ := keepertest.PricesKeepers(t)

	handler := prepare.PrepareProposalHandler(
		txConfig,
		mockBridgeKeeper,
		mockClobKeeper,
		mockPerpKeeper,
		prices.NewDefaultPriceUpdateGenerator(mockPricesKeeper),
	)

	req := abci.RequestPrepareProposal{
		Txs:        [][]byte{},
		MaxTxBytes: maxTxBytes,
	}

	resp, err := handler(ctx, &req)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(resp.Txs), 4)

	require.NoError(t, cmttypes.ToTxs(resp.Txs).Validate(maxTxBytes))

	opsTx, err := txConfig.TxDecoder()(resp.Txs[0])
	require.NoError(t, err)
	msgs := opsTx.GetMsgs()
	require.Len(t, msgs, 1)
	proposedOps, ok := msgs[0].(*clobtypes.MsgProposedOperations)
	require.True(t, ok)

	require.Empty(t, proposedOps.OperationsQueue)
	require.Equal(t, emptyOpsTx, resp.Txs[0])
}

func TestSlinkyPrepareProposalHandler(t *testing.T) {
	// test an empty UpdateMarketPrices tx is inserted if ves are not enabled
	t.Run("ves not enabled", func(t *testing.T) {
		// mocks
		mockPerpKeeper := mocks.PreparePerpetualsKeeper{}
		mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything).
			Return(constants.ValidMsgAddPremiumVotes)

		mockClobKeeper := mocks.PrepareClobKeeper{}
		mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
			Return(constants.ValidEmptyMsgProposedOperations)

		mockBridgeKeeper := mocks.PrepareBridgeKeeper{}
		mockBridgeKeeper.On("GetAcknowledgeBridges", mock.Anything, mock.Anything).
			Return(constants.MsgAcknowledgeBridges_Ids0_1_Height0)

		ctx := slinkytestutils.CreateBaseSDKContext(t)
		ctx = slinkytestutils.UpdateContextWithVEHeight(ctx, 3)
		ctx = ctx.WithBlockHeight(2) // disable vote-extensions

		gen := prices.NewSlinkyPriceUpdateGenerator(nil, nil, nil, nil) // ignore all fields, should immediately return

		handler := prepare.PrepareProposalHandler(
			encoding.GetTestEncodingCfg().TxConfig,
			&mockBridgeKeeper,
			&mockClobKeeper,
			&mockPerpKeeper,
			gen,
		)

		txs := [][]byte{}

		resp, err := handler(ctx, &abci.RequestPrepareProposal{Txs: txs, MaxTxBytes: 100_000})
		require.NoError(t, err)

		// expect all txs to have been inserted
		// check that the last tx is a valid update-market-prices tx
		marketPricesTx, err := encoding.GetTestEncodingCfg().TxConfig.TxDecoder()(resp.Txs[len(resp.Txs)-1])
		require.NoError(t, err)
		require.Len(t, marketPricesTx.GetMsgs(), 1)

		// expect the message to be an UpdateMarketPrices message w/ no markets
		updateMarketPricesMsg := marketPricesTx.GetMsgs()[0].(*pricestypes.MsgUpdateMarketPrices)
		require.Len(t, updateMarketPricesMsg.MarketPriceUpdates, 0)
	})

	// test that a valid UpdateMarketPricesTx is inserted if ves are enabled, and a valid ExtendedCommitInfo is present
	t.Run("ves enabled", func(t *testing.T) {
		// mocks
		mockPerpKeeper := mocks.PreparePerpetualsKeeper{}
		mockPerpKeeper.On("GetAddPremiumVotes", mock.Anything).
			Return(constants.ValidMsgAddPremiumVotes)

		mockClobKeeper := mocks.PrepareClobKeeper{}
		mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
			Return(constants.ValidEmptyMsgProposedOperations)

		mockBridgeKeeper := mocks.PrepareBridgeKeeper{}
		mockBridgeKeeper.On("GetAcknowledgeBridges", mock.Anything, mock.Anything).
			Return(constants.MsgAcknowledgeBridges_Ids0_1_Height0)

		ctx := slinkytestutils.CreateBaseSDKContext(t)
		ctx = slinkytestutils.UpdateContextWithVEHeight(ctx, 3)
		ctx = ctx.WithBlockHeight(4) // enable vote-extensions
		ctx = ctx.WithLogger(log.NewTestLogger(t))

		cpMock := strategymock.NewCurrencyPairStrategy(t)
		aggMock := aggregatormock.NewVoteAggregator(t)
		extCommitCodec := codec.NewDefaultExtendedCommitCodec()
		veCodec := codec.NewDefaultVoteExtensionCodec()
		gen := prices.NewSlinkyPriceUpdateGenerator(
			aggMock,
			extCommitCodec,
			veCodec,
			cpMock,
		)

		// mock dependencies
		validator1 := []byte("validator1")
		validator2 := []byte("validator2")

		validator1ve := vetypes.OracleVoteExtension{
			Prices: map[uint64][]byte{
				0: []byte("99"),
				1: []byte("100"),
			},
		}
		validator2ve := vetypes.OracleVoteExtension{
			Prices: map[uint64][]byte{
				0: []byte("99"),
				1: []byte("100"),
			},
		}

		validator1veBz, err := veCodec.Encode(validator1ve)
		require.NoError(t, err)

		validator2veBz, err := veCodec.Encode(validator2ve)
		require.NoError(t, err)

		// setup extendedCommit
		extCommit := abci.ExtendedCommitInfo{
			Votes: []abci.ExtendedVoteInfo{
				{
					Validator: abci.Validator{
						Address: validator1,
					},
					VoteExtension: validator1veBz,
				},
				{
					Validator: abci.Validator{
						Address: validator2,
					},
					VoteExtension: validator2veBz,
				},
			},
		}
		extCommitBz, err := extCommitCodec.Encode(extCommit)
		require.NoError(t, err)

		mogBtc := oracletypes.NewCurrencyPair("MOG", "BTC")
		tiaPepe := oracletypes.NewCurrencyPair("TIA", "PEPE")

		aggMock.On("AggregateOracleVotes", ctx, []aggregator.Vote{
			{
				ConsAddress:         validator1,
				OracleVoteExtension: validator1ve,
			},
			{
				ConsAddress:         validator2,
				OracleVoteExtension: validator2ve,
			},
		}).Return(map[oracletypes.CurrencyPair]*big.Int{
			mogBtc:  big.NewInt(100),
			tiaPepe: big.NewInt(99),
		}, nil)

		cpMock.On("ID", ctx, mogBtc).Return(uint64(0), nil)
		cpMock.On("ID", ctx, tiaPepe).Return(uint64(1), nil)

		handler := prepare.PrepareProposalHandler(
			encoding.GetTestEncodingCfg().TxConfig,
			&mockBridgeKeeper,
			&mockClobKeeper,
			&mockPerpKeeper,
			gen,
		)

		txs := [][]byte{
			extCommitBz, // extended commit should be first
			constants.Msg_Send_TxBytes,
		}

		resp, err := handler(ctx, &abci.RequestPrepareProposal{Txs: txs, MaxTxBytes: 100_000})
		require.NoError(t, err)

		// expect all txs to have been inserted
		// check that the last tx is a valid update-market-prices tx
		marketPricesTx, err := encoding.GetTestEncodingCfg().TxConfig.TxDecoder()(resp.Txs[len(resp.Txs)-1])
		require.NoError(t, err)
		require.Len(t, marketPricesTx.GetMsgs(), 1)

		// expect the message to be an UpdateMarketPrices message w/ no markets
		updateMarketPricesMsg := marketPricesTx.GetMsgs()[0].(*pricestypes.MsgUpdateMarketPrices)
		require.Len(t, updateMarketPricesMsg.MarketPriceUpdates, 2)

		expectedPrices := map[uint64]uint64{
			0: 100,
			1: 99,
		}

		for _, update := range updateMarketPricesMsg.MarketPriceUpdates {
			expectedPrice, ok := expectedPrices[uint64(update.MarketId)]
			require.True(t, ok)
			require.Equal(t, expectedPrice, update.Price)
		}
	})
}

func TestGetUpdateMarketPricesTx(t *testing.T) {
	tests := map[string]struct {
		keeperResp *pricestypes.MsgUpdateMarketPrices
		txEncoder  sdk.TxEncoder

		expectedTx         []byte
		expectedNumMarkets int
		expectedErr        error
	}{
		"nil message fails": {
			keeperResp: nil,

			expectedErr: fmt.Errorf("MsgUpdateMarketPrices cannot be nil"),
		},
		"empty message": {
			keeperResp: &pricestypes.MsgUpdateMarketPrices{}, // empty
			txEncoder:  passingTxEncoderOne,

			expectedTx:         []byte{1},
			expectedNumMarkets: 0,
		},
		"empty tx": {
			keeperResp: &pricestypes.MsgUpdateMarketPrices{},
			txEncoder:  emptyTxEncoder, // returns empty tx.

			expectedErr: fmt.Errorf("Invalid tx: []"),
		},
		"valid message, but encoding fails": {
			keeperResp: &pricestypes.MsgUpdateMarketPrices{}, // empty
			txEncoder:  failingTxEncoder,

			expectedErr: fmt.Errorf("encoder failed"),
		},
		"valid message": {
			keeperResp: &pricestypes.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{{}, {}, {}},
			},
			txEncoder: passingTxEncoderOne,

			expectedTx:         []byte{1},
			expectedNumMarkets: 3,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockTxConfig := createMockTxConfig(nil, []sdk.TxEncoder{tc.txEncoder})
			mockPricesKeeper := mocks.PreparePricesKeeper{}
			mockPricesKeeper.On("GetValidMarketPriceUpdates", mock.Anything).
				Return(tc.keeperResp)

			resp, err := getMarketPriceUpdates(prices.NewDefaultPriceUpdateGenerator(&mockPricesKeeper), mockTxConfig)

			if tc.expectedErr != nil {
				require.Equal(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTx, resp.Tx)
			require.Equal(t, tc.expectedNumMarkets, resp.NumMarkets)
		})
	}
}

func getMarketPriceUpdates(
	gen prices.PriceUpdateGenerator, txConfig client.TxConfig) (prepare.PricesTxResponse, error) {
	msg, err := gen.GetValidMarketPriceUpdates(sdk.Context{}, nil)
	if err != nil {
		return prepare.PricesTxResponse{}, err
	}

	return prepare.EncodeMarketPriceUpdates(txConfig, msg)
}

func TestGetAcknowledgeBridgesTx(t *testing.T) {
	tests := map[string]struct {
		keeperResp *bridgetypes.MsgAcknowledgeBridges
		txEncoder  sdk.TxEncoder

		expectedTx         []byte
		expectedNumBridges int
		expectedErr        error
	}{
		"empty list of msgs": {
			keeperResp: &bridgetypes.MsgAcknowledgeBridges{},
			txEncoder:  passingTxEncoderOne,

			expectedTx:         []byte{1},
			expectedNumBridges: 0,
		},
		"empty tx": {
			keeperResp: &bridgetypes.MsgAcknowledgeBridges{},
			txEncoder:  emptyTxEncoder, // returns empty tx.

			expectedErr: fmt.Errorf("Invalid tx: []"),
		},
		"valid messages, but encoding fails": {
			keeperResp: &bridgetypes.MsgAcknowledgeBridges{},
			txEncoder:  failingTxEncoder,

			expectedErr: fmt.Errorf("encoder failed"),
		},
		"1 bridge event": {
			keeperResp: constants.MsgAcknowledgeBridges_Id0_Height0,
			txEncoder:  passingTxEncoderTwo,

			expectedTx:         []byte{1, 2},
			expectedNumBridges: 1,
		},
		"2 bridge events": {
			keeperResp: constants.MsgAcknowledgeBridges_Ids0_1_Height0,
			txEncoder:  passingTxEncoderFour,

			expectedTx:         []byte{1, 2, 3, 4},
			expectedNumBridges: 2,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockTxConfig := createMockTxConfig(nil, []sdk.TxEncoder{tc.txEncoder})
			mockBridgeKeeper := mocks.PrepareBridgeKeeper{}
			mockBridgeKeeper.On("GetAcknowledgeBridges", mock.Anything, mock.Anything).
				Return(tc.keeperResp)

			resp, err := prepare.GetAcknowledgeBridgesTx(ctx, mockTxConfig, &mockBridgeKeeper)
			if tc.expectedErr != nil {
				require.Equal(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTx, resp.Tx)
			require.Equal(t, tc.expectedNumBridges, resp.NumBridges)
		})
	}
}

func TestGetAddPremiumVotesTx(t *testing.T) {
	tests := map[string]struct {
		keeperResp *perpetualtypes.MsgAddPremiumVotes
		txEncoder  sdk.TxEncoder

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
			mockTxConfig := createMockTxConfig(nil, []sdk.TxEncoder{tc.txEncoder})
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
		txEncoder  sdk.TxEncoder

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
			mockTxConfig := createMockTxConfig(nil, []sdk.TxEncoder{tc.txEncoder})
			mockClobKeeper := mocks.PrepareClobKeeper{}
			mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).Return(tc.keeperResp)

			resp, err := prepare.GetProposedOperationsTx(ctx, mockTxConfig, &mockClobKeeper, 1000)
			if tc.expectedErr != nil {
				require.Equal(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expectedTx, resp.Tx)
		})
	}
}

// buildOperationsQueue returns a queue of `n` distinct, realistically-sized short-term order
// placement `OperationRaw`s, encoded with the real tx encoder so that byte size grows meaningfully
// (and deterministically) with the length of the queue.
func buildOperationsQueue(n int) []clobtypes.OperationRaw {
	ops := make([]clobtypes.OperationRaw, n)
	for i := 0; i < n; i++ {
		ops[i] = clobtestutil.NewShortTermOrderPlacementOperationRaw(clobtypes.Order{
			OrderId:      clobtypes.OrderId{SubaccountId: constants.Alice_Num0, ClientId: uint32(i), ClobPairId: 0},
			Side:         clobtypes.Order_SIDE_BUY,
			Quantums:     100_000_000,
			Subticks:     50_000_000,
			GoodTilOneof: &clobtypes.Order_GoodTilBlock{GoodTilBlock: 25},
		})
	}
	return ops
}

// TestGetProposedOperationsTx_MaxBytes covers the prefix-truncation behavior of
// `GetProposedOperationsTx` using the real tx encoder (mocked encoders always return fixed-size
// output regardless of the number of operations, so they cannot exercise truncation).
func TestGetProposedOperationsTx_MaxBytes(t *testing.T) {
	txConfig := encoding.GetTestEncodingCfg().TxConfig

	const numOps = 50
	ops := buildOperationsQueue(numOps)

	mockClobKeeper := mocks.PrepareClobKeeper{}
	mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
		Return(&clobtypes.MsgProposedOperations{OperationsQueue: ops})

	encodeN := func(t *testing.T, n int) []byte {
		tx, err := prepare.EncodeMsgsIntoTxBytes(
			txConfig,
			&clobtypes.MsgProposedOperations{OperationsQueue: ops[:n]},
		)
		require.NoError(t, err)
		return tx
	}
	decodeOps := func(t *testing.T, tx []byte) []clobtypes.OperationRaw {
		decodedTx, err := txConfig.TxDecoder()(tx)
		require.NoError(t, err)
		msgs := decodedTx.GetMsgs()
		require.Len(t, msgs, 1)
		msgOps, ok := msgs[0].(*clobtypes.MsgProposedOperations)
		require.True(t, ok)
		return msgOps.OperationsQueue
	}

	emptyTx := encodeN(t, 0)
	fullTx := encodeN(t, numOps)
	oneOpTx := encodeN(t, 1)
	emptySize := uint64(len(emptyTx))
	fullSize := uint64(len(fullTx))
	oneOpSize := uint64(len(oneOpTx))
	// Sanity check that operation size grows meaningfully with queue length, otherwise the
	// prefix/empty-fits cases below wouldn't be meaningful.
	require.Less(t, emptySize, oneOpSize)
	require.Less(t, oneOpSize, fullSize)

	t.Run("fits", func(t *testing.T) {
		maxBytes := fullSize + 1_000 // comfortably larger than the full ops tx.
		resp, err := prepare.GetProposedOperationsTx(ctx, txConfig, &mockClobKeeper, maxBytes)
		require.NoError(t, err)
		require.Equal(t, numOps, resp.NumOperations)
		require.Equal(t, fullTx, resp.Tx)
		require.Equal(t, ops, decodeOps(t, resp.Tx))
	})

	t.Run("prefix", func(t *testing.T) {
		maxBytes := emptySize + (fullSize-emptySize)/2
		resp, err := prepare.GetProposedOperationsTx(ctx, txConfig, &mockClobKeeper, maxBytes)
		require.NoError(t, err)
		require.Greater(t, resp.NumOperations, 0)
		require.Less(t, resp.NumOperations, numOps)
		require.LessOrEqual(t, uint64(len(resp.Tx)), maxBytes)
		require.Equal(t, encodeN(t, resp.NumOperations), resp.Tx)
		require.Equal(t, ops[:resp.NumOperations], decodeOps(t, resp.Tx))
	})

	t.Run("empty-fits", func(t *testing.T) {
		// Strictly between the empty-tx size and the smallest non-empty (1-op) tx size, so that
		// only the empty ops tx fits.
		maxBytes := emptySize + (oneOpSize-emptySize)/2
		require.Greater(t, maxBytes, emptySize)
		require.Less(t, maxBytes, oneOpSize)

		resp, err := prepare.GetProposedOperationsTx(ctx, txConfig, &mockClobKeeper, maxBytes)
		require.NoError(t, err)
		require.Equal(t, 0, resp.NumOperations)
		require.Equal(t, emptyTx, resp.Tx)
	})

	t.Run("nothing-fits", func(t *testing.T) {
		_, err := prepare.GetProposedOperationsTx(ctx, txConfig, &mockClobKeeper, emptySize-1)
		require.Error(t, err)
	})
}

func TestEncodeMsgsIntoTxBytes(t *testing.T) {
	tests := map[string]struct {
		setMsgErr error
		txEncoder sdk.TxEncoder

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
			mockTxConfig := createMockTxConfig(tc.setMsgErr, []sdk.TxEncoder{tc.txEncoder})

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

func createMockTxConfig(setMsgsError error, allTxEncoders []sdk.TxEncoder) *mocks.TxConfig {
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
