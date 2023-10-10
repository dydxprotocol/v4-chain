package prepare_test

import (
	"bytes"
	"testing"
	"time"

	gometrics "github.com/armon/go-metrics"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/flags"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
)

// TestFullNodePrepareProposalHandler test that the full-node PrepareProposal handler always returns
// an empty result.
func TestFullNodePrepareProposalHandler(t *testing.T) {
	defer gometrics.Shutdown()

	conf := gometrics.DefaultConfig("testService")
	sink := gometrics.NewInmemSink(time.Hour, time.Hour)
	_, err := gometrics.NewGlobal(conf, sink)
	require.NoError(t, err)

	var logBuffer bytes.Buffer
	appOpts := map[string]interface{}{
		flags.NonValidatingFullNodeFlag: true,
		testApp.LoggerInstanceForTest:   log.TestingLoggerWithOutput(&logBuffer),
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

		if metric, ok := metrics.Counters["testService.prepare_proposal.handler.error.count;detail=prepare_proposal_txs"]; ok {
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
	require.True(t, found)
	require.Contains(t, logBuffer.String(), "This validator may be incorrectly running in full-node mode!")
}
