package shared

import (
	"context"
	"github.com/cosmos/cosmos-sdk/types/query"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

const (
	// PaginatedRequestLimit is the maximum number of entries that can be returned in a paginated request
	PaginatedRequestLimit = 10000
)

// AllPaginatedMarketParams returns all MarketParams from the prices module, paginating through the results and
// returning a list of all the aggregated market-params.
func AllPaginatedMarketParams(ctx context.Context, client pricestypes.QueryClient) ([]pricestypes.MarketParam, error) {
	mps := make([]pricestypes.MarketParam, 0)

	pq := func(ctx context.Context, req *query.PageRequest) (ResponseWithPagination, error) {
		resp, err := client.AllMarketParams(ctx, &pricestypes.QueryAllMarketParamsRequest{
			Pagination: req,
		})
		if err != nil {
			return nil, err
		}

		mps = append(mps, resp.MarketParams...)
		return resp, nil
	}

	if err := HandlePaginatedQuery(ctx, pq, &query.PageRequest{
		Limit: PaginatedRequestLimit,
	}); err != nil {
		return nil, err
	}

	return mps, nil
}

// ResponseWithPagination represents a response-type from a cosmos-module's GRPC service for entries that are paginated
type ResponseWithPagination interface {
	GetPagination() *query.PageResponse
}

// PaginatedQuery is a function type that represents a paginated query to a cosmos-module's GRPC service
type PaginatedQuery func(ctx context.Context, req *query.PageRequest) (ResponseWithPagination, error)

func HandlePaginatedQuery(ctx context.Context, pq PaginatedQuery, initialPagination *query.PageRequest) error {
	for {
		// make the query
		resp, err := pq(ctx, initialPagination)
		if err != nil {
			return err
		}

		// break if there is no next-key
		if resp.GetPagination() == nil || len(resp.GetPagination().NextKey) == 0 {
			return nil
		}

		// otherwise, update the next-key and continue
		initialPagination.Key = resp.GetPagination().NextKey
	}
}
