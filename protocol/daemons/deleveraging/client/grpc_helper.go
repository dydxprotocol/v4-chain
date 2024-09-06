package client

import (
	"context"
	"fmt"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/deleveraging/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib/metrics"
	blocktimetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime/types"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	gometrics "github.com/hashicorp/go-metrics"
	"google.golang.org/grpc/metadata"
)

// GetPreviousBlockInfo queries a gRPC server using `QueryPreviousBlockInfoRequest`
// and returns the previous block height.
func (c *Client) GetPreviousBlockInfo(
	ctx context.Context,
) (
	blockHeight uint32,
	err error,
) {
	defer metrics.ModuleMeasureSince(
		metrics.DeleveragingDaemon,
		metrics.DaemonGetPreviousBlockInfoLatency,
		time.Now(),
	)

	query := &blocktimetypes.QueryPreviousBlockInfoRequest{}
	response, err := c.BlocktimeQueryClient.PreviousBlockInfo(ctx, query)
	if err != nil {
		return 0, err
	}

	return response.Info.Height, nil
}

// GetAllSubaccounts queries a gRPC server and returns a list of subaccounts and
// their balances and open positions.
func (c *Client) GetAllSubaccounts(
	ctx context.Context,
	pageLimit uint64,
) (
	subaccounts []satypes.Subaccount,
	err error,
) {
	defer telemetry.ModuleMeasureSince(metrics.DeleveragingDaemon, time.Now(), metrics.GetAllSubaccounts, metrics.Latency)
	subaccounts = make([]satypes.Subaccount, 0)

	var nextKey []byte
	for {
		subaccountsFromKey, next, err := getSubaccountsFromKey(
			ctx,
			c.SubaccountQueryClient,
			nextKey,
			pageLimit,
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
		metrics.DeleveragingDaemon,
		float32(len(subaccounts)),
		metrics.GetAllSubaccounts,
		metrics.Count,
	)

	return subaccounts, nil
}

// SendDeleveragingSubaccountIds sends a list of subaccounts with open positions for each perp to a gRPC server via `UpdateSubaccountsListForDeleveragingDaemon`.
func (c *Client) SendDeleveragingSubaccountIds(
	ctx context.Context,
	openPositionInfoMap map[uint32]*clobtypes.SubaccountOpenPositionInfo,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.DeleveragingDaemon,
		time.Now(),
		metrics.SendDeleveragingSubaccountIds,
		metrics.Latency,
	)

	// Convert the map to a slice.
	// Note that sorting here is not strictly necessary but is done for safety and to avoid making
	// any assumptions on the server side.
	sortedPerpetualIds := lib.GetSortedKeys[lib.Sortable[uint32]](openPositionInfoMap)
	subaccountOpenPositionInfo := make([]clobtypes.SubaccountOpenPositionInfo, 0)
	for _, perpetualId := range sortedPerpetualIds {
		subaccountOpenPositionInfo = append(subaccountOpenPositionInfo, *openPositionInfoMap[perpetualId])
	}

	request := &api.UpdateSubaccountsListForDeleveragingDaemonRequest{
		SubaccountOpenPositionInfo: subaccountOpenPositionInfo,
	}

	if _, err := c.DeleveragingServiceClient.UpdateSubaccountsListForDeleveragingDaemon(ctx, request); err != nil {
		return err
	}
	return nil
}

func newContextWithQueryBlockHeight(
	ctx context.Context,
	blockHeight uint32,
) context.Context {
	return metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs(
			grpc.GRPCBlockHeightHeader,
			fmt.Sprintf("%d", blockHeight),
		),
	)
}

func getSubaccountsFromKey(
	ctx context.Context,
	client satypes.QueryClient,
	pageRequestKey []byte,
	limit uint64,
) (
	subaccounts []satypes.Subaccount,
	nextKey []byte,
	err error,
) {
	defer metrics.ModuleMeasureSinceWithLabels(
		metrics.DeleveragingDaemon,
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
