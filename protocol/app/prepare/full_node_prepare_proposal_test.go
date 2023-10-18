package prepare_test

import (
	"testing"
	"time"

	gometrics "github.com/armon/go-metrics"
	abci "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/flags"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	testlog "github.com/dydxprotocol/v4-chain/protocol/testutil/logger"
	"github.com/stretchr/testify/require"
)

// TestFullNodePrepareProposalHandler test that the full-node PrepareProposal handler always returns
// an empty result.
func TestFullNodePrepareProposalHandler(t *testing.T) {
	t.Cleanup(gometrics.Shutdown)

	conf := gometrics.DefaultConfig("service")
	sink := gometrics.NewInmemSink(time.Hour, time.Hour)
	_, err := gometrics.NewGlobal(conf, sink)
	require.NoError(t, err)

	logger, logBuffer := testlog.TestLogger()
	appOpts := map[string]interface{}{
		flags.NonValidatingFullNodeFlag: true,
		testlog.LoggerInstanceForTest:   logger,
	}
	tApp := testApp.NewTestAppBuilder().WithTesting(t).WithAppCreatorFn(testApp.DefaultTestAppCreatorFn(appOpts)).Build()

	found := false
	tApp.AdvanceToBlock(2, testApp.AdvanceToBlockOptions{
		BlockTime:                         time.Time{},
		RequestPrepareProposalTxsOverride: [][]byte{{9}, {9, 8}, {9, 8, 7}},
		ValidateRespPrepare: func(context sdktypes.Context, proposal abci.ResponsePrepareProposal) (haltChain bool) {
			require.Empty(t, proposal.Txs)
			return true
		},
	})

	for _, metrics := range sink.Data() {
		metrics.RLock()
		defer metrics.RUnlock()

		if metric, ok := metrics.Counters["service.prepare_proposal.handler.error.count;detail=prepare_proposal_txs"]; ok {
			require.Equal(t,
				[]gometrics.Label{{
					Name:  "detail",
					Value: "prepare_proposal_txs",
				}},
				metric.Labels)
			require.Equal(t, 1, metric.Count)
			require.Equal(t, float64(1), metric.Sum)
			found = true
		}
	}
	require.True(t, found, "Expected metric not found")
	require.Contains(
		t,
		logBuffer.String(),
		"This validator may be incorrectly running in full-node mode!",
		"Expected log message not found",
	)
}
