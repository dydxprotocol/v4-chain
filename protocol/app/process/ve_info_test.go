package process_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/process"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	prepareutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestVEInjectionHandling(t *testing.T) {
	// Valid order tx.
	validOperationsTx := constants.ValidEmptyMsgProposedOperationsTxBytes

	// Valid add funding tx.
	validAddFundingTx := constants.ValidMsgAddPremiumVotesTxBytes

	tests := map[string]struct {
		txsBytes [][]byte

		expectedTxCount int
	}{
		"Valid: block with VE's": {
			txsBytes: [][]byte{
				{}, // empty for ve.
				validOperationsTx,
				validAddFundingTx,
			},

			expectedTxCount: 3,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup.
			ctx, pricesKeeper, _, daemonPriceCache, _, mockTimeProvider := keepertest.PricesKeepers(t)

			ctx = vetesting.GetVeEnabledCtx(ctx, 3)
			ctx = ctx.WithCometInfo(
				vetesting.NewBlockInfo(
					nil,
					nil,
					nil,
					abci.CommitInfo{
						Round: 3,
						Votes: []abci.VoteInfo{},
					},
				),
			)
			mockTimeProvider.On("Now").Return(constants.TimeT)
			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)
			daemonPriceCache.UpdatePrices(constants.AtTimeTSingleExchangePriceUpdate)

			mockClobKeeper := &mocks.ProcessClobKeeper{}
			mockClobKeeper.On("RecordMevMetricsIsEnabled").Return(true)
			mockClobKeeper.On("RecordMevMetrics", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			mockRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}

			mockVEApplier := &mocks.ProcessProposalVEApplier{}
			mockVEApplier.On("ApplyVE", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			handler := process.ProcessProposalHandler(
				constants.TestEncodingCfg.TxConfig,
				mockClobKeeper,
				&mocks.ProcessPerpetualKeeper{},
				pricesKeeper,
				mockRatelimitKeeper,
				vecodec.NewDefaultExtendedCommitCodec(),
				vecodec.NewDefaultVoteExtensionCodec(),
				mockVEApplier,
				prepareutils.NoOpValidateVoteExtensionsFn,
			)

			req := &abci.RequestProcessProposal{
				Txs: tc.txsBytes,
			}

			_, err := handler(ctx, req)
			require.NoError(t, err)

			require.Equal(t, tc.expectedTxCount, len(req.Txs))
		})
	}
}
