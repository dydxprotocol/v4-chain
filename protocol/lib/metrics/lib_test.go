package metrics_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	gometrics "github.com/hashicorp/go-metrics"

	"github.com/stretchr/testify/require"
)

func TestSetGaugeWithLabelsForExecMode(t *testing.T) {
	t.Cleanup(gometrics.Shutdown)
	context := sdk.Context{}
	conf := gometrics.DefaultConfig("testService")
	conf.EnableHostname = false
	sink := gometrics.NewInmemSink(time.Hour, time.Hour)
	_, err := gometrics.NewGlobal(conf, sink)
	require.NoError(t, err)

	context = context.WithExecMode(sdk.ExecModeFinalize)
	metrics.EmitTelemetryWithLabelsForExecMode(
		context,
		[]sdk.ExecMode{sdk.ExecModeFinalize},
		metrics.SetGaugeWithLabels,
		"testKey1",
		3.14,
		gometrics.Label{
			Name:  "testLabel",
			Value: "testLabelValue",
		},
	)

	metrics.EmitTelemetryWithLabelsForExecMode(
		context,
		[]sdk.ExecMode{sdk.ExecModeSimulate},
		metrics.SetGaugeWithLabels,
		"testKey2",
		3.14,
		gometrics.Label{
			Name:  "testLabel",
			Value: "testLabelValue",
		},
	)

	FinalizeModeKeyFound := false
	SimulateModeKeyFound := false
	for _, metrics := range sink.Data() {
		metrics.RLock()
		defer metrics.RUnlock()

		if metric, ok := metrics.Gauges["testService.testKey1;testLabel=testLabelValue"]; ok {
			require.Equal(t,
				[]gometrics.Label{{
					Name:  "testLabel",
					Value: "testLabelValue",
				}},
				metric.Labels)
			require.Equal(t, float32(3.14), metric.Value)
			FinalizeModeKeyFound = true
		}
		if _, ok := metrics.Gauges["testService.testKey2;testLabel=testLabelValue"]; ok {
			SimulateModeKeyFound = true
		}
	}
	require.True(t, FinalizeModeKeyFound)
	require.False(t, SimulateModeKeyFound)
}
