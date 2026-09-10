package prepare_test

import (
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare/prices"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/encoding"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestPrepareProposal_OverBudgetOperationsQueueYieldsNonEmptyProposal verifies that when the CLOB
// operations queue does not fit in the block's remaining byte budget, PrepareProposalHandler still
// returns a valid, non-empty proposal by truncating the operations tx to the longest fitting prefix,
// rather than collapsing to an empty proposal that ProcessProposal would reject. It asserts the
// result against CometBFT's own block-validity check.
func TestPrepareProposal_OverBudgetOperationsQueueYieldsNonEmptyProposal(t *testing.T) {
	txConfig := encoding.GetTestEncodingCfg().TxConfig

	mockPricesKeeper, mockPerpKeeper, mockBridgeKeeper := setUpFixedGroupMocks()

	// A large operations queue whose full encoded tx is far bigger than the budget granted below,
	// so the full queue genuinely does not fit.
	const numOps = 200
	ops := buildOperationsQueue(numOps)
	mockClobKeeper := &mocks.PrepareClobKeeper{}
	mockClobKeeper.On("GetOperations", mock.Anything, mock.Anything).
		Return(&clobtypes.MsgProposedOperations{OperationsQueue: ops})

	ctx, _, _, _, _, _, _ := keepertest.PricesKeepers(t)

	// Derive MaxTxBytes from actual encoder output rather than guessed constants.
	pricesTxResp, err := prepare.EncodeMarketPriceUpdates(txConfig, &pricestypes.MsgUpdateMarketPrices{})
	require.NoError(t, err)
	fundingTxResp, err := prepare.GetAddPremiumVotesTx(ctx, txConfig, mockPerpKeeper)
	require.NoError(t, err)
	bridgeTxResp, err := prepare.GetAcknowledgeBridgesTx(ctx, txConfig, mockBridgeKeeper)
	require.NoError(t, err)
	fullOpsTx, err := prepare.EncodeMsgsIntoTxBytes(
		txConfig,
		&clobtypes.MsgProposedOperations{OperationsQueue: ops},
	)
	require.NoError(t, err)

	fixedRawTotal := uint64(len(pricesTxResp.Tx) + len(fundingTxResp.Tx) + len(bridgeTxResp.Tx))
	fullOpsSize := uint64(cmttypes.ComputeProtoSizeForTxs([]cmttypes.Tx{fullOpsTx}))

	// Grant the operations group only half of what its full queue needs, so the full queue does not
	// fit but there is ample room for a non-trivial truncated prefix.
	opsBudget := fullOpsSize / 2
	// Slack so the fixed-group txs (and their small per-tx proto framing) always fit, keeping the
	// test isolated to the ops-truncation path.
	const slack = uint64(512)
	maxTxBytes := int64(fixedRawTotal + slack + opsBudget) //nolint:gosec

	require.Less(
		t, opsBudget, fullOpsSize,
		"test is only meaningful if the full operations queue does not fit in its budget",
	)

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

	// A valid block must contain at least the 4 mandatory txs (operations, bridges, funding, prices),
	// and must pass CometBFT's block-size validation.
	require.GreaterOrEqual(t, len(resp.Txs), 4)
	require.NoError(t, cmttypes.ToTxs(resp.Txs).Validate(maxTxBytes))
}
