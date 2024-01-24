package mev_telemetry

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/log"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

// SendDatapoints sends MEV metrics to an HTTP-based metric collection service
func SendDatapoints(ctx sdk.Context, addresses []string, mevMetrics types.MevMetrics) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), metrics.MevSentDatapoints, metrics.Latency)

	defer func() {
		if r := recover(); r != nil {
			log.ErrorLog(
				ctx,
				"panic when recording mev metrics",
				"error", r,
				log.StackTrace, string(debug.Stack()),
			)
		}
	}()

	for _, address := range addresses {
		sendDatapointsToTelemetryService(ctx, address, mevMetrics)
	}
}

func sendDatapointsToTelemetryService(ctx sdk.Context, address string, mevMetrics types.MevMetrics) {
	data, err := json.Marshal(mevMetrics)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"error marshalling mev metrics",
			err,
		)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	resp, err := client.Post(address, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"error sending mev metric",
			err,
			log.Address, address,
		)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.ErrorLogWithError(
			ctx,
			"error reading response",
			err,
			log.Address, address,
		)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	if len(responseBody) == 0 {
		log.ErrorLog(
			ctx,
			"0-byte response from mev telemetry server",
			log.Address, address,
		)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.ErrorLog(
			ctx,
			"error sending mev datapoint, non 200 http status code",
			log.Address, address,
			log.StatusCode, resp.StatusCode,
		)
		telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Error, metrics.Count)
		return
	}

	telemetry.IncrCounter(1, types.ModuleName, metrics.MevSentDatapoints, metrics.Success, metrics.Count)
}
