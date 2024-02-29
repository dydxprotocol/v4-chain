package e2e_test

import (
	"context"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	blocktime "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	"github.com/stretchr/testify/require"
	"sync"
	"sync/atomic"
	"testing"
)

// There is some basic validation in this test to ensure that numbers are monotonically increasing
// but the meaningful validation will come from Go's ability to perform data race detection during testing
// when using the `-race` flag.
func TestParallelQuery(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	tApp.InitChain()

	// In parallel:
	//   - advance the block
	//   - get the app/version
	//   - load block info from store directly
	//   - perform a custom gRPC query to get the block info

	// We specifically use an atomic to ensure that we aren't providing any synchronization between the threads
	// maximizing any data races that could exist. The wait group is only used to synchronize the testing thread
	// when the other 4 threads are done.
	blockLimitReached := atomic.Bool{}
	blockLimitReached.Store(false)
	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		defer func() {
			blockLimitReached.Store(true)
		}()
		for i := uint32(2); i < 50; i++ {
			tApp.AdvanceToBlock(i, testapp.AdvanceToBlockOptions{})
		}
	}()

	version := make([]*abcitypes.ResponseQuery, 0)
	versionErrs := make([]error, 0)
	go func() {
		defer wg.Done()
		for !blockLimitReached.Load() {
			resp, err := tApp.App.Query(context.Background(), &abcitypes.RequestQuery{
				Path: "app/version",
			})
			version = append(version, resp)
			versionErrs = append(versionErrs, err)
		}
	}()

	store := make([]*abcitypes.ResponseQuery, 0)
	storeErrs := make([]error, 0)
	go func() {
		defer wg.Done()
		for !blockLimitReached.Load() {
			resp, err := tApp.App.Query(context.Background(),
				&abcitypes.RequestQuery{
					Path: "store/blocktime/key",
					Data: []byte(blocktime.PreviousBlockInfoKey),
				},
			)
			store = append(store, resp)
			storeErrs = append(storeErrs, err)
		}
	}()

	grpc := make([]*abcitypes.ResponseQuery, 0)
	grpcErrs := make([]error, 0)
	blocktimeRequest := blocktime.QueryPreviousBlockInfoRequest{}
	blockTimeRequestBytes := tApp.App.AppCodec().MustMarshal(&blocktimeRequest)
	go func() {
		defer wg.Done()
		for !blockLimitReached.Load() {
			resp, err := tApp.App.Query(
				context.Background(),
				&abcitypes.RequestQuery{
					Path: "/dydxprotocol.blocktime.Query/PreviousBlockInfo",
					Data: blockTimeRequestBytes,
				},
			)
			grpc = append(grpc, resp)
			grpcErrs = append(grpcErrs, err)
		}
	}()

	wg.Wait()

	previousVersionHeight := int64(0)
	for i := range version {
		require.NoError(t, versionErrs[i])
		require.NotNil(t, version[i])
		require.GreaterOrEqual(t, version[i].Height, previousVersionHeight)
		previousVersionHeight = version[i].Height
	}

	previousStoreHeight := uint32(0)
	for i := range store {
		require.NoError(t, storeErrs[i])
		require.NotNil(t, store[i])

		var blockInfo blocktime.BlockInfo
		tApp.App.AppCodec().MustUnmarshal(store[i].Value, &blockInfo)
		require.GreaterOrEqual(t, blockInfo.Height, previousStoreHeight)
		require.Equal(t, store[i].Height, int64(blockInfo.Height))
		previousStoreHeight = blockInfo.Height
	}

	previousGrpcHeight := uint32(0)
	for i := range grpc {
		require.NoError(t, grpcErrs[i])
		require.NotNil(t, grpc[i])

		var query blocktime.QueryPreviousBlockInfoResponse
		tApp.App.AppCodec().MustUnmarshal(grpc[i].Value, &query)
		require.GreaterOrEqual(t, query.Info.Height, previousGrpcHeight)
		require.Equal(t, grpc[i].Height, int64(query.Info.Height))
		previousGrpcHeight = query.Info.Height
	}
}
