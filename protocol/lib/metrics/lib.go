package metrics

import (
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// main entrypoint for logging in the v4 protocol
// type Label struct {
// 	Name  string
// 	Value string
// }

type Label = gometrics.Label

// IncrCounter provides a wrapper functionality for emitting a counter
// metric with global labels (if any) along with the provided labels.
func IncrCounterWithLabels(key string, val float32, labels ...Label) {
	telemetry.IncrCounterWithLabels([]string{key}, val, labels)
}

// IncrCounter provides a wrapper functionality for emitting a counter
// metric with global labels (if any) along with the provided labels.
func IncrCounter(key string, val float32) {
	telemetry.IncrCounterWithLabels([]string{key}, val, []gometrics.Label{})
}

// Gauge provides a wrapper functionality for emitting a counter
// metric with global labels (if any) along with the provided labels.
func SetGaugeWithLabels(key string, val float32, labels ...gometrics.Label) {
	telemetry.SetGaugeWithLabels([]string{key}, val, labels)
}

// Gauge provides a wrapper functionality for emitting a counter
// metric with global labels (if any) along with the provided labels.
func SetGauge(key string, val float32) {
	telemetry.SetGaugeWithLabels([]string{key}, val, []gometrics.Label{})
}

// Histogram provides a wrapper functionality for emitting a counter
// metric with global labels (if any) along with the provided labels.
func AddSampleWithLabels(key string, val float32, labels ...gometrics.Label) {
	// TODO why the f is this a differnet library
	gometrics.AddSampleWithLabels(
		[]string{key},
		val,
		labels,
	)
}

// Histogram provides a wrapper functionality for emitting a counter
// metric with global labels (if any) along with the provided labels.
func AddSample(key string, val float32) {
	// TODO why the f is this a differnet library
	gometrics.AddSampleWithLabels(
		[]string{key},
		val,
		[]gometrics.Label{},
	)
}

func ModuleMeasureSince(module string, key string, start time.Time) {
	telemetry.ModuleMeasureSince(
		module,
		start,
		key,
	)
}

// ModuleMeasureSinceWithLabels provides a short hand method for emitting a time measure
// metric for a module with a given set of keys and labels.
// NOTE: global labels are not included in this metric.
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
