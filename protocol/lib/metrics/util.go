package metrics

import (
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	gometrics "github.com/hashicorp/go-metrics"
)

// IncrCountMetricWithLabels increases a count metric from a module with the provided labels by a count of 1.
func IncrCountMetricWithLabels(module string, metric string, labels ...gometrics.Label) {
	telemetry.IncrCounterWithLabels(
		[]string{module, metric, Count},
		1,
		labels,
	)
}

// IncrSuccessOrErrorCounter increments either the success or error counter for a given handler
// based on whether the given error is nil or not. This function is intended to be called in a
// defer block at the top of any function which returns an error.
func IncrSuccessOrErrorCounter(err error, module string, handler string, callback string, labels ...gometrics.Label) {
	successOrError := Success
	if err != nil {
		successOrError = Error
	}

	telemetry.IncrCounterWithLabels(
		[]string{
			module,
			handler,
			successOrError,
			Count,
		},
		1,
		append(
			[]gometrics.Label{
				GetLabelForStringValue(Callback, callback),
			},
			labels...,
		),
	)
}

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

// GetMetricValueFromBigInt returns a telemetry value (float32) from an integer value.
// Any rounding information is ignored, so this function should only be used for metrics.
func GetMetricValueFromBigInt(i *big.Int) float32 {
	r, _ := new(big.Float).SetInt(i).Float32()
	return r
}

// GetCallbackMetricFromCtx determines the callback metric based on the context. Note that DeliverTx is implied
// if the context is not CheckTx or ReCheckTx. This function is unable to account for other callbacks like
// PrepareCheckState or EndBlocker.
func GetCallbackMetricFromCtx(ctx sdk.Context) string {
	if ctx.IsReCheckTx() {
		return ReCheckTx
	}
	if ctx.IsCheckTx() {
		return CheckTx
	}

	return DeliverTx
}
