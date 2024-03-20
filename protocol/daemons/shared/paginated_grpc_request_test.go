package shared_test

import (
	"testing"

	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/shared"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/stretchr/testify/require"
)

var (
	nextKey1 = []byte("next-key-1")
	nextKey2 = []byte("next-key-2")
)

// tests that AllPaginatedMarketParams makes queries with the given next-key + returns the
// aggregated responses
func TestPaginatedGRPCRequests(t *testing.T) {
	qc := mocks.NewQueryClient(t)

	marketParam1 := pricestypes.MarketParam{
		Id:   1,
		Pair: "BTC-USDC",
	}

	marketParam2 := pricestypes.MarketParam{
		Id:   2,
		Pair: "ETH-USDC",
	}

	marketParam3 := pricestypes.MarketParam{
		Id:   3,
		Pair: "LINK-USDC",
	}

	ctx := sdk.Context{}

	// ON ...
	// a query that returns a single-market param with a next-key-1 on receiving a page-request
	// with a limit of types.PaginatedRequestLimit
	initialPagination := &query.PageRequest{
		Limit: shared.PaginatedRequestLimit,
	}
	qc.On("AllMarketParams", ctx, &pricestypes.QueryAllMarketParamsRequest{
		Pagination: initialPagination,
	}).Return(
		&pricestypes.QueryAllMarketParamsResponse{
			MarketParams: []pricestypes.MarketParam{marketParam1},
			Pagination: &query.PageResponse{
				NextKey: nextKey1,
			},
		},
		nil,
	)

	// a query that returns a single-market param with a next-key-2 on receiving a page-request
	// with a limit of types.PaginatedRequestLimit with the next-key-1 return marketParam2 and next-key-2
	qc.On("AllMarketParams", ctx, &pricestypes.QueryAllMarketParamsRequest{
		Pagination: &query.PageRequest{
			Limit: shared.PaginatedRequestLimit,
			Key:   nextKey1,
		},
	}).Return(
		&pricestypes.QueryAllMarketParamsResponse{
			MarketParams: []pricestypes.MarketParam{marketParam2},
			Pagination: &query.PageResponse{
				NextKey: nextKey2,
			},
		},
		nil,
	)

	// a query that returns a single-market param with no next-key on receiving a page-request with a limit
	// of types.PaginatedRequestLimit with the next-key-2 return marketParam3 and no next-key
	qc.On("AllMarketParams", ctx, &pricestypes.QueryAllMarketParamsRequest{
		Pagination: &query.PageRequest{
			Limit: shared.PaginatedRequestLimit,
			Key:   nextKey2,
		},
	}).Return(
		&pricestypes.QueryAllMarketParamsResponse{
			MarketParams: []pricestypes.MarketParam{marketParam3},
			Pagination:   &query.PageResponse{},
		},
		nil,
	)

	// ... THEN
	// the function should return all the market-params
	marketParams, err := shared.AllPaginatedMarketParams(ctx, qc)
	require.NoError(t, err)

	require.Equal(t, []pricestypes.MarketParam{marketParam1, marketParam2, marketParam3}, marketParams)
}

func TestPaginatedGRPCRequestsWithError(t *testing.T) {
	qc := mocks.NewQueryClient(t)

	marketParam1 := pricestypes.MarketParam{
		Id:   1,
		Pair: "BTC-USDC",
	}

	ctx := sdk.Context{}

	// ON ...
	// a query that returns a single-market param with a next-key-1 on receiving a page-request
	// with a limit of types.PaginatedRequestLimit
	initialPagination := &query.PageRequest{
		Limit: shared.PaginatedRequestLimit,
	}
	qc.On("AllMarketParams", ctx, &pricestypes.QueryAllMarketParamsRequest{
		Pagination: initialPagination,
	}).Return(
		&pricestypes.QueryAllMarketParamsResponse{
			MarketParams: []pricestypes.MarketParam{marketParam1},
			Pagination: &query.PageResponse{
				NextKey: nextKey1,
			},
		},
		nil,
	)

	expErr := fmt.Errorf("error")
	// a query that returns an error on receiving a page-request with a limit of types.PaginatedRequestLimit + next-key-1
	qc.On("AllMarketParams", ctx, &pricestypes.QueryAllMarketParamsRequest{
		Pagination: &query.PageRequest{
			Limit: shared.PaginatedRequestLimit,
			Key:   nextKey1,
		},
	}).Return(
		nil,
		expErr,
	)

	// ... THEN
	// the function should return the error
	_, err := shared.AllPaginatedMarketParams(ctx, qc)
	require.EqualError(t, err, expErr.Error())
}
