package client

import (
	"context"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetPreviousBlockInfo queries a gRPC server using `QueryPreviousBlockInfoRequest`
// and returns the previous block height.
func GetPreviousBlockInfo(
	ctx context.Context,
	daemon *Client,
) (
	blockHeight uint32,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.GetPreviousBlockInfo,
		metrics.Latency,
	)

	query := &blocktimetypes.QueryPreviousBlockInfoRequest{}
	response, err := daemon.BlocktimeQueryClient.PreviousBlockInfo(ctx, query)
	if err != nil {
		return 0, err
	}

	return response.Info.Height, nil
}

// GetAllSubaccounts queries a gRPC server and returns a list of subaccounts and
// their balances and open positions.
func GetAllSubaccounts(
	ctx context.Context,
	daemon *Client,
	limit uint64,
) (
	subaccounts []satypes.Subaccount,
	err error,
) {
	defer telemetry.ModuleMeasureSince(metrics.LiquidationDaemon, time.Now(), metrics.GetAllSubaccounts, metrics.Latency)
	subaccounts = make([]satypes.Subaccount, 0)

	var nextKey []byte
	for {
		subaccountsFromKey, next, err := getSubaccountsFromKey(
			ctx,
			daemon.SubaccountQueryClient,
			limit,
			nextKey,
		)

		if err != nil {
			return nil, err
		}

		subaccounts = append(subaccounts, subaccountsFromKey...)
		nextKey = next

		if len(nextKey) == 0 {
			break
		}
	}

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(subaccounts)),
		metrics.GetAllSubaccounts,
		metrics.Count,
	)

	return subaccounts, nil
}

// CheckCollateralizationForSubaccounts queries a gRPC server using `AreSubaccountsLiquidatable`
// and returns a list of collateralization statuses for the given list of subaccount ids.
func CheckCollateralizationForSubaccounts(
	ctx context.Context,
	daemon *Client,
	subaccountIds []satypes.SubaccountId,
) (
	results []clobtypes.AreSubaccountsLiquidatableResponse_Result,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.CheckCollateralizationForSubaccounts,
		metrics.Latency,
	)

	query := &clobtypes.AreSubaccountsLiquidatableRequest{
		SubaccountIds: subaccountIds,
	}
	response, err := daemon.ClobQueryClient.AreSubaccountsLiquidatable(ctx, query)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}

// SendLiquidatableSubaccountIds sends a list of unique and potentially liquidatable
// subaccount ids to a gRPC server via `LiquidateSubaccounts`.
func SendLiquidatableSubaccountIds(
	ctx context.Context,
	daemon *Client,
	subaccountIds []satypes.SubaccountId,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.SendLiquidatableSubaccountIds,
		metrics.Latency,
	)

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(subaccountIds)),
		metrics.LiquidatableSubaccountIds,
		metrics.Count,
	)

	request := &api.LiquidateSubaccountsRequest{
		SubaccountIds: subaccountIds,
	}

	if _, err := daemon.LiquidationServiceClient.LiquidateSubaccounts(ctx, request); err != nil {
		return err
	}
	return nil
}

func getSubaccountsFromKey(
	ctx context.Context,
	client satypes.QueryClient,
	limit uint64,
	pageRequestKey []byte,
) (
	subaccounts []satypes.Subaccount,
	nextKey []byte,
	err error,
) {
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.LiquidationDaemon,
		[]string{metrics.GetSubaccountsFromKey, metrics.Latency},
		time.Now(),
		[]gometrics.Label{
			metrics.GetLabelForIntValue(metrics.PageLimit, int(limit)),
		},
	)

	query := &satypes.QueryAllSubaccountRequest{
		Pagination: &query.PageRequest{
			Key:   pageRequestKey,
			Limit: limit,
		},
	}

	response, err := client.SubaccountAll(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	if response.Pagination != nil {
		nextKey = response.Pagination.NextKey
	}
	return response.Subaccount, nextKey, nil
}
