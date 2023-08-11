package metrics

import (
	"strconv"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// GetLabelForBoolValue returns a telemetry label for a given label and bool value.
func GetLabelForBoolValue(labelName string, labelValue bool) gometrics.Label {
	return GetLabelForStringValue(labelName, strconv.FormatBool(labelValue))
}

// GetLabelForIntValue returns a telemetry label for a given label and int value.
func GetLabelForIntValue(labelName string, labelValue int) gometrics.Label {
	return GetLabelForStringValue(labelName, strconv.Itoa(labelValue))
}

// GetLabelForStringValue returns a telemetry label for a given label and string value.
func GetLabelForStringValue(labelName string, labelValue string) gometrics.Label {
	return telemetry.NewLabel(labelName, labelValue)
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
