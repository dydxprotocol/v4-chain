package lib

import (
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/stretchr/testify/require"
)

func TestTelemetryMetrics(t *testing.T) {
	m, err := telemetry.New(telemetry.Config{
		Enabled:                 true,
		EnableHostname:          false,
		ServiceName:             "test",
		PrometheusRetentionTime: 60,
		EnableHostnameLabel:     false,
	})
	require.NoError(t, err)
	require.NotNil(t, m)

	// IncrCounter
	// Increment the counter five times.
	for i := 0; i < 5; i++ {
		telemetry.IncrCounter(1.1, "dummy", "counter")
	}

	response, err := m.Gather(telemetry.FormatPrometheus)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(response.Metrics), "test_dummy_counter 5.5"))

	// SetGauge with three different values (1.2, 2.2, 3.2).
	telemetry.SetGauge(1.2, "dummy", "gauge")
	telemetry.SetGauge(2.2, "dummy", "gauge")
	telemetry.SetGauge(3.2, "dummy", "gauge")

	response, err = m.Gather(telemetry.FormatPrometheus)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(response.Metrics), "test_dummy_gauge 3.2"))

	// MeasureSince
	// Measure 1 second interval three times.
	for i := 0; i < 3; i++ {
		telemetry.MeasureSince(time.Now().Add(-1*time.Second), "dummy", "latency")
	}

	response, err = m.Gather(telemetry.FormatPrometheus)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(response.Metrics), "test_dummy_latency{quantile=\"0.99\"} 1000"))
	require.True(t, strings.Contains(string(response.Metrics), "test_dummy_latency_sum 3000"))
	require.True(t, strings.Contains(string(response.Metrics), "test_dummy_latency_count 3"))
}
