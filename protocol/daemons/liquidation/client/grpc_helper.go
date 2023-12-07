package client

import (
	"context"
	"time"

	gometrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// GetAllSubaccounts queries a gRPC server and returns a list of subaccounts and
// their balances and open positions.
func GetAllSubaccounts(
	daemon *Client,
	ctx context.Context,
	client satypes.QueryClient,
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
			client,
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
	daemon *Client,
	ctx context.Context,
	client clobtypes.QueryClient,
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
	response, err := client.AreSubaccountsLiquidatable(ctx, query)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}

// SendLiquidatableSubaccountIds sends a list of unique and potentially liquidatable
// subaccount ids to a gRPC server via `LiquidateSubaccounts`.
func SendLiquidatableSubaccountIds(
	ctx context.Context,
	client api.LiquidationServiceClient,
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

	if _, err := client.LiquidateSubaccounts(ctx, request); err != nil {
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
