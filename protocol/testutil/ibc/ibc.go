package ibc

import (
	"fmt"

	"github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
)

// Given a denom trace string (e.g. `transfer/channel-141/uosmo`),
// output the IBC denom.
// For more information on IBC denom, see:
// https://interchainacademy.cosmos.network/tutorials/5-ibc-dev/#how-are-ibc-denoms-derived
func DenomTraceToIBCDenom(denomTraceStr string) (
	ibcDenom string,
	err error,
) {
	// Parse the denom trace
	denomTrace := types.ParseDenomTrace(denomTraceStr)

	if denomTrace.Path == "" || denomTrace.BaseDenom == "" {
		return "", fmt.Errorf("invalid denom trace '%+v' is parsed into empty path or base denom", denomTraceStr)
	}

	// Get the full denom path
	return denomTrace.IBCDenom(), nil
}
