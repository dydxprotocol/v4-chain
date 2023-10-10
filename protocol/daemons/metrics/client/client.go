package client

import (
	"context"
	"runtime/debug"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

var (
	// 30 minutes
	METRICS_DAEMON_LOOP_DELAY_MS uint32 = 30 * 60 * 1000
)

// Start begins a job that periodically:
// 1) Emits metrics about app version and git commit.
// This job should never panic or block the application from running.
func Start(
	ctx context.Context,
	logger log.Logger,
) {
	ticker := time.NewTicker(time.Duration(METRICS_DAEMON_LOOP_DELAY_MS) * time.Millisecond)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		RunMetricsDaemonTaskLoop(
			ctx,
			logger,
		)
	}
}

// RunMetricsDaemonTaskLoop contains the logic to emit metrics about the application running.
func RunMetricsDaemonTaskLoop(
	ctx context.Context,
	logger log.Logger,
) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error(
				"panic when reporting metrics from metrics daemon",
				"panic",
				r,
				"stack",
				string(debug.Stack()),
			)
		}
	}()

	// Report out app version and git commit.
	version := version.NewInfo()
	telemetry.SetGaugeWithLabels(
		[]string{metrics.AppInfo},
		1,
		[]gometrics.Label{
			metrics.GetLabelForStringValue(metrics.AppVersion, version.Version),
			metrics.GetLabelForStringValue(metrics.GitCommit, version.GitCommit),
		},
	)
	logger.Info(
		"App version",
		metrics.AppVersion,
		version.Version,
		metrics.GitCommit,
		version.GitCommit,
	)
}
