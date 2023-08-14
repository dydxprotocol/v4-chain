package mev_telemetry

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	sdklog "cosmossdk.io/log"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var client = &http.Client{
	Timeout: 2 * time.Second,
}

func logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, "x/clob/mev_telemetry")
}

// SendDatapoints sends MEV metrics to an HTTP-based metric collection service
func SendDatapoints(ctx sdk.Context, address string, mevMetrics types.MevMetrics) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.MevSentDatapoints, metrics.Latency)

	defer func() {
		if r := recover(); r != nil {
			logger(ctx).Error(
				"panic when recording mev metrics",
				"panic",
				r,
				"stack trace",
				string(debug.Stack()),
			)
		}
	}()

	data, err := json.Marshal(mevMetrics)

	if err != nil {
		logger(ctx).Error("error marshalling mev metrics", "error", err)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	resp, err := client.Post(address, "application/json", bytes.NewBuffer(data))

	if err != nil {
		logger(ctx).Error("error sending mev metric", "error", err)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		logger(ctx).Error("error reading response", "error", err)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	if len(responseBody) == 0 {
		logger(ctx).Error("0-byte response from mev telemetry server")
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logger(ctx).Error("error sending mev datapoint", "error", "non-200 http status-code", "status_code", resp.StatusCode)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Success, metrics.Count)
}
