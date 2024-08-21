package client

import (
	"context"
	"fmt"
	"time"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
		metrics.LiquidationDaemon,
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

// GetAllPerpetuals queries gRPC server and returns a list of perpetuals.
func (c *Client) GetAllPerpetuals(
	ctx context.Context,
	pageLimit uint64,
) (
	perpetuals []perptypes.Perpetual,
	err error,
) {
	defer metrics.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		metrics.DaemonGetAllPerpetualsLatency,
		time.Now(),
	)

	perpetuals = make([]perptypes.Perpetual, 0)

	var nextKey []byte
	for {
		perpetualsFromKey, next, err := getPerpetualsFromKey(
			ctx,
			c.PerpetualsQueryClient,
			nextKey,
			pageLimit,
		)

		if err != nil {
			return nil, err
		}

		perpetuals = append(perpetuals, perpetualsFromKey...)
		nextKey = next

		if len(nextKey) == 0 {
			break
		}
	}
	return perpetuals, nil
}

// GetAllLiquidityTiers queries gRPC server and returns a list of liquidityTiers.
func (c *Client) GetAllLiquidityTiers(
	ctx context.Context,
	pageLimit uint64,
) (
	liquidityTiers []perptypes.LiquidityTier,
	err error,
) {
	defer metrics.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		metrics.DaemonGetAllLiquidityTiersLatency,
		time.Now(),
	)

	liquidityTiers = make([]perptypes.LiquidityTier, 0)

	var nextKey []byte
	for {
		liquidityTiersFromKey, next, err := getLiquidityTiersFromKey(
			ctx,
			c.PerpetualsQueryClient,
			nextKey,
			pageLimit,
		)

		if err != nil {
			return nil, err
		}

		liquidityTiers = append(liquidityTiers, liquidityTiersFromKey...)
		nextKey = next

		if len(nextKey) == 0 {
			break
		}
	}
	return liquidityTiers, nil
}

// GetAllMarketPrices queries gRPC server and returns a list of market prices.
func (c *Client) GetAllMarketPrices(
	ctx context.Context,
	pageLimit uint64,
) (
	marketPrices []pricestypes.MarketPrice,
	err error,
) {
	defer metrics.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		metrics.DaemonGetAllMarketPricesLatency,
		time.Now(),
	)

	marketPrices = make([]pricestypes.MarketPrice, 0)

	var nextKey []byte
	for {
		marketPricesFromKey, next, err := getMarketPricesFromKey(
			ctx,
			c.PricesQueryClient,
			nextKey,
			pageLimit,
		)

		if err != nil {
			return nil, err
		}

		marketPrices = append(marketPrices, marketPricesFromKey...)
		nextKey = next

		if len(nextKey) == 0 {
			break
		}
	}
	return marketPrices, nil
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
	defer telemetry.ModuleMeasureSince(metrics.LiquidationDaemon, time.Now(), metrics.GetAllSubaccounts, metrics.Latency)
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
		metrics.LiquidationDaemon,
		float32(len(subaccounts)),
		metrics.GetAllSubaccounts,
		metrics.Count,
	)

	return subaccounts, nil
}

// SendLiquidatableSubaccountIds sends a list of unique and potentially liquidatable
// subaccount ids to a gRPC server via `LiquidateSubaccounts`.
func (c *Client) SendLiquidatableSubaccountIds(
	ctx context.Context,
	blockHeight uint32,
	liquidatableSubaccountIds []satypes.SubaccountId,
	negativeTncSubaccountIds []satypes.SubaccountId,
	openPositionInfoMap map[uint32]*clobtypes.SubaccountOpenPositionInfo,
	pageLimit uint64,
) error {
	defer telemetry.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		time.Now(),
		metrics.SendLiquidatableSubaccountIds,
		metrics.Latency,
	)

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(liquidatableSubaccountIds)),
		metrics.LiquidatableSubaccountIds,
		metrics.Count,
	)
	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(negativeTncSubaccountIds)),
		metrics.NegativeTncSubaccountIds,
		metrics.Count,
	)

	// Convert the map to a slice.
	// Note that sorting here is not strictly necessary but is done for safety and to avoid making
	// any assumptions on the server side.
	sortedPerpetualIds := lib.GetSortedKeys[lib.Sortable[uint32]](openPositionInfoMap)
	subaccountOpenPositionInfo := make([]clobtypes.SubaccountOpenPositionInfo, 0)
	for _, perpetualId := range sortedPerpetualIds {
		subaccountOpenPositionInfo = append(subaccountOpenPositionInfo, *openPositionInfoMap[perpetualId])
	}

	// Break this down to multiple requests if the number of subaccounts is too large.

	// Liquidatable subaccount ids.
	requests := GenerateLiquidateSubaccountsPaginatedRequests(
		liquidatableSubaccountIds,
		blockHeight,
		pageLimit,
	)

	// Negative TNC subaccount ids.
	requests = append(
		requests,
		GenerateNegativeTNCSubaccountsPaginatedRequests(
			negativeTncSubaccountIds,
			blockHeight,
			pageLimit,
		)...,
	)

	// Subaccount open position info.
	requests = append(
		requests,
		GenerateSubaccountOpenPositionPaginatedRequests(
			subaccountOpenPositionInfo,
			blockHeight,
			pageLimit,
		)...,
	)

	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(requests)),
		metrics.NumRequests,
		metrics.Count,
	)

	for _, req := range requests {
		if _, err := c.LiquidationServiceClient.LiquidateSubaccounts(ctx, req); err != nil {
			return errorsmod.Wrapf(
				err,
				"failed to send liquidatable subaccount ids to protocol at block height %d",
				blockHeight,
			)
		}
	}

	return nil
}

func GenerateLiquidateSubaccountsPaginatedRequests(
	ids []satypes.SubaccountId,
	blockHeight uint32,
	pageLimit uint64,
) []*api.LiquidateSubaccountsRequest {
	if len(ids) == 0 {
		return []*api.LiquidateSubaccountsRequest{
			{
				BlockHeight:               blockHeight,
				LiquidatableSubaccountIds: []satypes.SubaccountId{},
			},
		}
	}

	requests := make([]*api.LiquidateSubaccountsRequest, 0)
	for start := 0; start < len(ids); start += int(pageLimit) {
		end := lib.Min(start+int(pageLimit), len(ids))
		request := &api.LiquidateSubaccountsRequest{
			BlockHeight:               blockHeight,
			LiquidatableSubaccountIds: ids[start:end],
		}
		requests = append(requests, request)
	}
	return requests
}

func GenerateNegativeTNCSubaccountsPaginatedRequests(
	ids []satypes.SubaccountId,
	blockHeight uint32,
	pageLimit uint64,
) []*api.LiquidateSubaccountsRequest {
	if len(ids) == 0 {
		return []*api.LiquidateSubaccountsRequest{
			{
				BlockHeight:              blockHeight,
				NegativeTncSubaccountIds: []satypes.SubaccountId{},
			},
		}
	}

	requests := make([]*api.LiquidateSubaccountsRequest, 0)
	for start := 0; start < len(ids); start += int(pageLimit) {
		end := lib.Min(start+int(pageLimit), len(ids))
		request := &api.LiquidateSubaccountsRequest{
			BlockHeight:              blockHeight,
			NegativeTncSubaccountIds: ids[start:end],
		}
		requests = append(requests, request)
	}
	return requests
}

func GenerateSubaccountOpenPositionPaginatedRequests(
	subaccountOpenPositionInfo []clobtypes.SubaccountOpenPositionInfo,
	blockHeight uint32,
	pageLimit uint64,
) []*api.LiquidateSubaccountsRequest {
	if len(subaccountOpenPositionInfo) == 0 {
		return []*api.LiquidateSubaccountsRequest{
			{
				BlockHeight:                blockHeight,
				SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{},
			},
		}
	}

	requests := make([]*api.LiquidateSubaccountsRequest, 0)
	for _, info := range subaccountOpenPositionInfo {
		// Long positions.
		for start := 0; start < len(info.SubaccountsWithLongPosition); start += int(pageLimit) {
			end := lib.Min(start+int(pageLimit), len(info.SubaccountsWithLongPosition))
			request := &api.LiquidateSubaccountsRequest{
				BlockHeight: blockHeight,
				SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
					{
						PerpetualId:                 info.PerpetualId,
						SubaccountsWithLongPosition: info.SubaccountsWithLongPosition[start:end],
					},
				},
			}
			requests = append(requests, request)
		}

		// Short positions.
		for start := 0; start < len(info.SubaccountsWithShortPosition); start += int(pageLimit) {
			end := lib.Min(start+int(pageLimit), len(info.SubaccountsWithShortPosition))
			request := &api.LiquidateSubaccountsRequest{
				BlockHeight: blockHeight,
				SubaccountOpenPositionInfo: []clobtypes.SubaccountOpenPositionInfo{
					{
						PerpetualId:                  info.PerpetualId,
						SubaccountsWithShortPosition: info.SubaccountsWithShortPosition[start:end],
					},
				},
			}
			requests = append(requests, request)
		}
	}
	return requests
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

func getMarketPricesFromKey(
	ctx context.Context,
	client pricestypes.QueryClient,
	pageRequestKey []byte,
	limit uint64,
) (
	marketPrices []pricestypes.MarketPrice,
	nextKey []byte,
	err error,
) {
	defer metrics.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		metrics.DaemonGetMarketPricesPaginatedLatency,
		time.Now(),
	)

	query := &pricestypes.QueryAllMarketPricesRequest{
		Pagination: &query.PageRequest{
			Key:   pageRequestKey,
			Limit: limit,
		},
	}

	response, err := client.AllMarketPrices(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	if response.Pagination != nil {
		nextKey = response.Pagination.NextKey
	}
	return response.MarketPrices, nextKey, nil
}

func getPerpetualsFromKey(
	ctx context.Context,
	client perptypes.QueryClient,
	pageRequestKey []byte,
	limit uint64,
) (
	perpetuals []perptypes.Perpetual,
	nextKey []byte,
	err error,
) {
	defer metrics.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		metrics.DaemonGetPerpetualsPaginatedLatency,
		time.Now(),
	)

	query := &perptypes.QueryAllPerpetualsRequest{
		Pagination: &query.PageRequest{
			Key:   pageRequestKey,
			Limit: limit,
		},
	}

	response, err := client.AllPerpetuals(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	if response.Pagination != nil {
		nextKey = response.Pagination.NextKey
	}
	return response.Perpetual, nextKey, nil
}

func getLiquidityTiersFromKey(
	ctx context.Context,
	client perptypes.QueryClient,
	pageRequestKey []byte,
	limit uint64,
) (
	liquidityTiers []perptypes.LiquidityTier,
	nextKey []byte,
	err error,
) {
	defer metrics.ModuleMeasureSince(
		metrics.LiquidationDaemon,
		metrics.DaemonGetLiquidityTiersPaginatedLatency,
		time.Now(),
	)

	query := &perptypes.QueryAllLiquidityTiersRequest{
		Pagination: &query.PageRequest{
			Key:   pageRequestKey,
			Limit: limit,
		},
	}

	response, err := client.AllLiquidityTiers(ctx, query)
	if err != nil {
		return nil, nil, err
	}
	if response.Pagination != nil {
		nextKey = response.Pagination.NextKey
	}
	return response.LiquidityTiers, nextKey, nil
}
