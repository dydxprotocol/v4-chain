package metrics

import (
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// This file provides a main entrypoint for logging in the v4 protocol.
// TODO(CLOB-1013) Drop both metrics libraries above for a library
// that supports float64 (i.e hashicorp go-metrics)

type Label = gometrics.Label

// IncrCounterWithLabels provides a wrapper functionality for emitting a counter
// metric with global labels (if any) along with the provided labels.
func IncrCounterWithLabels(key string, val float32, labels ...Label) {
	telemetry.IncrCounterWithLabels([]string{key}, val, labels)
}

// IncrCounter provides a wrapper functionality for emitting a counter
// metric with global labels (if any).
func IncrCounter(key string, val float32) {
	telemetry.IncrCounterWithLabels([]string{key}, val, []gometrics.Label{})
}

// SetGaugeWithLabels provides a wrapper functionality for emitting a gauge
// metric with global labels (if any) along with the provided labels.
func SetGaugeWithLabels(key string, val float32, labels ...gometrics.Label) {
	telemetry.SetGaugeWithLabels([]string{key}, val, labels)
}

// SetGauge provides a wrapper functionality for emitting a gauge
// metric with global labels (if any).
func SetGauge(key string, val float32) {
	telemetry.SetGaugeWithLabels([]string{key}, val, []gometrics.Label{})
}

// AddSampleWithLabels provides a wrapper functionality for emitting a sample
// metric with the provided labels.
func AddSampleWithLabels(key string, val float32, labels ...gometrics.Label) {
	gometrics.AddSampleWithLabels(
		[]string{key},
		val,
		labels,
	)
}

// AddSample provides a wrapper functionality for emitting a sample
// metric.
func AddSample(key string, val float32) {
	gometrics.AddSampleWithLabels(
		[]string{key},
		val,
		[]gometrics.Label{},
	)
}

// ModuleMeasureSince provides a wrapper functionality for emitting a time measure
// metric with global labels (if any).
// Please try to use `AddSample` instead.
// TODO(CLOB-1022) Roll our own calculations for timing on top of AddSample instead
// of using MeasureSince.
func ModuleMeasureSince(module string, key string, start time.Time) {
	telemetry.ModuleMeasureSince(
		module,
		start,
		key,
	)
}

// ModuleMeasureSinceWithLabels provides a short hand method for emitting a time measure
// metric for a module with labels. Global labels are not included in this metric.
// Please try to use `AddSampleWithLabels` instead.
// TODO(CLOB-1022) Roll our own calculations for timing on top of AddSample instead
// of using MeasureSince.
func ModuleMeasureSinceWithLabels(
	module string,
	keys []string,
	start time.Time,
	labels []gometrics.Label,
) {
	gometrics.MeasureSinceWithLabels(
		keys,
		start.UTC(),
		append(
			[]gometrics.Label{telemetry.NewLabel(telemetry.MetricLabelNameModule, module)},
			labels...,
		),
	)
}
