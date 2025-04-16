package client_test

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	slinkytypes "github.com/dydxprotocol/slinky/pkg/types"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/slinky/client"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestMarketPairFetcher(t *testing.T) {
	logger := log.NewTestLogger(t)
	queryClient := mocks.NewQueryClient(t)
	fetcher := client.MarketPairFetcherImpl{
		Logger:            logger,
		PricesQueryClient: queryClient,
	}
	asset0 := "FOO"
	asset1 := "BAR"
	pair0 := types.MarketParam{Id: 0, Pair: fmt.Sprintf("%s-%s", asset0, asset1)}
	pair1 := types.MarketParam{Id: 1, Pair: fmt.Sprintf("%s-%s", asset1, asset0)}
	invalidPair := types.MarketParam{Id: 2, Pair: "foobar"}

	t.Run("caches and returns valid pairs", func(t *testing.T) {
		queryClient.
			On("AllMarketParams", mock.Anything, mock.Anything).
			Return(
				&types.QueryAllMarketParamsResponse{
					MarketParams: []types.MarketParam{
						pair0,
						pair1,
					}},
				nil,
			).Once()
		err := fetcher.FetchIdMappings(context.Background())
		require.NoError(t, err)
		id, err := fetcher.GetIDForPair(slinkytypes.CurrencyPair{Base: asset0, Quote: asset1})
		require.NoError(t, err)
		require.Equal(t, pair0.Id, id)
		id, err = fetcher.GetIDForPair(slinkytypes.CurrencyPair{Base: asset1, Quote: asset0})
		require.NoError(t, err)
		require.Equal(t, pair1.Id, id)
	})

	t.Run("errors on fetch non-cached pair", func(t *testing.T) {
		queryClient.
			On("AllMarketParams", mock.Anything, mock.Anything).
			Return(
				&types.QueryAllMarketParamsResponse{
					MarketParams: []types.MarketParam{}},
				nil,
			).Once()
		err := fetcher.FetchIdMappings(context.Background())
		require.NoError(t, err)
		_, err = fetcher.GetIDForPair(slinkytypes.CurrencyPair{Base: asset0, Quote: asset1})
		require.Error(t, err, fmt.Errorf("pair %s/%s not found in compatMappings", asset0, asset1))
	})

	t.Run("fails on fetching invalid pairs", func(t *testing.T) {
		queryClient.
			On("AllMarketParams", mock.Anything, mock.Anything).
			Return(
				&types.QueryAllMarketParamsResponse{
					MarketParams: []types.MarketParam{
						invalidPair,
					}},
				nil,
			).Once()
		err := fetcher.FetchIdMappings(context.Background())
		require.Error(t, err, "incorrectly formatted CurrencyPair: foobar")
	})

	t.Run("fails on prices query error", func(t *testing.T) {
		queryClient.
			On("AllMarketParams", mock.Anything, mock.Anything).
			Return(
				&types.QueryAllMarketParamsResponse{},
				fmt.Errorf("test error"),
			).Once()
		err := fetcher.FetchIdMappings(context.Background())
		require.Error(t, err, "test error")
	})
}
