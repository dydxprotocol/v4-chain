package client_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	pricefeedclient "github.com/dydxprotocol/v4/daemons/pricefeed/client"
	pricefeedconstants "github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/handler"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/types"
	pricefeedtypes "github.com/dydxprotocol/v4/daemons/pricefeed/types"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/client"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/testutil/grpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	subTaskRunnerImpl = pricefeedclient.SubTaskRunnerImpl{}
)

// FakeSubTaskRunner acts as a dummy struct replacing `SubTaskRunner` that simply advances the
// counter for each task in a threadsafe manner and allows awaiting go-routine completion. This
// struct should only be used for testing.
type FakeSubTaskRunner struct {
	sync.WaitGroup
	sync.RWMutex
	UpdaterCallCount int
	EncoderCallCount int
	FetcherCallCount int
}

// StartPriceUpdater replaces `client.StartPriceUpdater` and advances `UpdaterCallCount` by one.
func (f *FakeSubTaskRunner) StartPriceUpdater(
	ctx context.Context,
	exchangeToMarketPrices *types.ExchangeToMarketPrices,
	priceFeedServiceClient api.PriceFeedServiceClient,
	loopDelayMs uint32,
	logger log.Logger,
) {
	// No need to lock/unlock since there is only one updater running and no risk of race-condition.
	f.UpdaterCallCount += 1
}

// StartPriceEncoder replaces `client.StartPriceEncoder`, marks the embedded waitgroup done and
// advances `EncoderCallCount` by one. This function will be called from a go-routine and is
// threadsafe.
func (f *FakeSubTaskRunner) StartPriceEncoder(
	exchangeFeedId types.ExchangeFeedId,
	exchangeToMarketPrices *types.ExchangeToMarketPrices,
	logger log.Logger,
	bCh chan *pricefeedclient.PriceFetcherSubtaskResponse,
) {
	f.Lock()
	defer f.Unlock()

	f.EncoderCallCount += 1
	f.Done()
}

// StartPriceFetcher replaces `client.StartPriceFetcher`, marks the embedded waitgroup done and
// advances `FetcherCallCount` by one. This function will be called from a go-routine and is
// threadsafe.
func (f *FakeSubTaskRunner) StartPriceFetcher(
	exchangeConfig types.ExchangeConfig,
	queryHandler handler.ExchangeQueryHandler,
	logger log.Logger,
	bCh chan *pricefeedclient.PriceFetcherSubtaskResponse,
) {
	f.Lock()
	defer f.Unlock()

	f.FetcherCallCount += 1
	f.Done()
}

const (
	maxBufferedChannelLength             = 2
	priceUpdaterLoopDelayMs              = 1000
	connectionFailsErrorMsg              = "Failed to create connection"
	closeConnectionFailsErrorMsg         = "Failed to close connection"
	missingMarketsExchangeFeedId         = 1000
	missingExchangeDetailsExchangeFeedId = 1001
	fiveKilobytes                        = 5 * 1024
)

var (
	staticExchangeStartupConfigLength = len(pricefeedconstants.StaticExchangeStartupConfig)
)

func TestFixedBufferSize(t *testing.T) {
	require.Equal(t, fiveKilobytes, pricefeedclient.FixedBufferSize)
}

func TestStart(t *testing.T) {
	tests := map[string]struct {
		// parameters
		mockGrpcClient                  *mocks.GrpcClient
		exchangeFeedIdToStartupConfig   map[types.ExchangeFeedId]*types.ExchangeStartupConfig
		exchangeFeedIdToMarkets         map[types.ExchangeFeedId][]types.MarketId
		exchangeFeedIdToExchangeDetails map[types.ExchangeFeedId]types.ExchangeQueryDetails
		// expectations
		expectedError         error
		expectCloseConnection bool
		// This should equal the length of the `exchangeFeedIdToStartupConfig` passed into
		// `client.Start`.
		expectedNumExchangeTasks int
	}{
		"Connection Fails": {
			mockGrpcClient: grpc.GenerateMockGrpcClientWithReturns(
				errors.New(connectionFailsErrorMsg),
				nil,
				false,
			),
			expectedError: errors.New(connectionFailsErrorMsg),
		},
		"Connection Succeeds": {
			mockGrpcClient:                grpc.GenerateMockGrpcClientWithReturns(nil, nil, true),
			exchangeFeedIdToStartupConfig: pricefeedconstants.StaticExchangeStartupConfig,
			expectCloseConnection:         true,
			expectedNumExchangeTasks:      staticExchangeStartupConfigLength,
		},
		"Empty exchange startup config": {
			mockGrpcClient:                grpc.GenerateMockGrpcClientWithReturns(nil, nil, true),
			exchangeFeedIdToStartupConfig: map[types.ExchangeFeedId]*types.ExchangeStartupConfig{},
			expectedError:                 errors.New("exchangeFeedIds must not be empty"),
			expectCloseConnection:         true,
		},
		"Invalid: exchange feed id in startup config does not have corresponding markets": {
			mockGrpcClient: grpc.GenerateMockGrpcClientWithReturns(nil, nil, true),
			exchangeFeedIdToStartupConfig: map[types.ExchangeFeedId]*types.ExchangeStartupConfig{
				missingMarketsExchangeFeedId: {},
			},
			expectedError: fmt.Errorf(
				"no exchange information exists for exchangeFeedId: %v",
				missingMarketsExchangeFeedId,
			),
			expectCloseConnection: true,
		},
		"Invalid: exchange feed id in startup config does not have corresponding exchange query details": {
			mockGrpcClient: grpc.GenerateMockGrpcClientWithReturns(nil, nil, true),
			exchangeFeedIdToStartupConfig: map[types.ExchangeFeedId]*types.ExchangeStartupConfig{
				missingExchangeDetailsExchangeFeedId: {},
			},
			exchangeFeedIdToMarkets: map[types.ExchangeFeedId][]types.MarketId{
				missingExchangeDetailsExchangeFeedId: {1},
			},
			expectedError: fmt.Errorf(
				"no exchange details exists for exchangeFeedId: %v",
				missingExchangeDetailsExchangeFeedId,
			),
			expectCloseConnection: true,
		},
		"Close connection fails": {
			mockGrpcClient: grpc.GenerateMockGrpcClientWithReturns(
				nil,
				errors.New(closeConnectionFailsErrorMsg),
				true,
			),
			exchangeFeedIdToStartupConfig:   pricefeedconstants.StaticExchangeStartupConfig,
			exchangeFeedIdToMarkets:         pricefeedconstants.StaticExchangeMarkets,
			exchangeFeedIdToExchangeDetails: pricefeedconstants.StaticExchangeDetails,
			expectedError:                   errors.New(closeConnectionFailsErrorMsg),
			expectCloseConnection:           true,
			expectedNumExchangeTasks:        staticExchangeStartupConfigLength,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			faketaskRunner := FakeSubTaskRunner{
				UpdaterCallCount: 0,
				EncoderCallCount: 0,
				FetcherCallCount: 0,
			}

			// Wait for each encoder and fetcher call to complete.
			faketaskRunner.WaitGroup.Add(tc.expectedNumExchangeTasks * 2)

			// Run Start.
			if tc.expectedError != nil {
				require.EqualError(
					t,
					pricefeedclient.Start(
						grpc.Ctx,
						grpc.SocketPath,
						log.NewNopLogger(),
						tc.mockGrpcClient,
						tc.exchangeFeedIdToStartupConfig,
						tc.exchangeFeedIdToMarkets,
						tc.exchangeFeedIdToExchangeDetails,
						priceUpdaterLoopDelayMs,
						&faketaskRunner,
					),
					tc.expectedError.Error(),
				)
			} else {
				require.NoError(
					t,
					pricefeedclient.Start(
						grpc.Ctx,
						grpc.SocketPath,
						log.NewNopLogger(),
						tc.mockGrpcClient,
						pricefeedconstants.StaticExchangeStartupConfig,
						pricefeedconstants.StaticExchangeMarkets,
						pricefeedconstants.StaticExchangeDetails,
						priceUpdaterLoopDelayMs,
						&faketaskRunner,
					),
				)
			}

			// Wait for encoder and fetcher go-routines to complete and thenv erify each subtask was
			// called the expected amount.
			faketaskRunner.Wait()
			require.Equal(t, tc.expectedNumExchangeTasks, faketaskRunner.EncoderCallCount)
			require.Equal(t, tc.expectedNumExchangeTasks, faketaskRunner.FetcherCallCount)
			if tc.expectedNumExchangeTasks > 0 {
				require.Equal(t, 1, faketaskRunner.UpdaterCallCount)
			} else {
				require.Equal(t, 0, faketaskRunner.UpdaterCallCount)
			}

			// Verify new connection and close connection calls.
			tc.mockGrpcClient.AssertCalled(t, "NewGrpcConnection", grpc.Ctx, grpc.SocketPath)
			if tc.expectCloseConnection {
				tc.mockGrpcClient.AssertCalled(t, "CloseConnection", grpc.ClientConn)
			} else {
				tc.mockGrpcClient.AssertNotCalled(t, "CloseConnection", grpc.ClientConn)
			}
		})
	}
}

func TestPriceEncoder_NoWrites(t *testing.T) {
	etmp, bChMap := generateBufferedChannelAndExchangeToMarketPrices(t, constants.Exchange1Exchange2Array)

	runPriceEncoderSequentially(
		t,
		constants.ExchangeFeedId1,
		etmp,
		bChMap[constants.ExchangeFeedId1],
		[]*types.MarketPriceTimestamp{},
	)

	require.Empty(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp)
	require.Empty(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp)
	require.Empty(t, bChMap[constants.ExchangeFeedId1])
	require.Empty(t, bChMap[constants.ExchangeFeedId2])
}

func TestPriceEncoder_DoNotWriteError(t *testing.T) {
	etmp, bChMap := generateBufferedChannelAndExchangeToMarketPrices(t, constants.Exchange1Exchange2Array)

	bCh := bChMap[constants.ExchangeFeedId1]

	bCh <- &pricefeedclient.PriceFetcherSubtaskResponse{nil, errors.New("Failed to query")}

	close(bCh)
	subTaskRunnerImpl.StartPriceEncoder(constants.ExchangeFeedId1, etmp, log.NewNopLogger(), bCh)

	require.Empty(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp)
	require.Empty(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp)
	require.Empty(t, bChMap[constants.ExchangeFeedId1])
	require.Empty(t, bChMap[constants.ExchangeFeedId2])
}

func TestPriceEncoder_WriteToOneMarket(t *testing.T) {
	etmp, bChMap := generateBufferedChannelAndExchangeToMarketPrices(t, constants.Exchange1Exchange2Array)

	runPriceEncoderSequentially(
		t,
		constants.ExchangeFeedId1,
		etmp,
		bChMap[constants.ExchangeFeedId1],
		[]*types.MarketPriceTimestamp{
			constants.Market9_TimeT_Price1,
		},
	)

	require.Len(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp, 1)
	require.Empty(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp)

	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price1,
			LastUpdateTime: constants.TimeT,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp[constants.MarketId9],
	)
}

func TestPriceEncoder_WriteToTwoMarkets(t *testing.T) {
	etmp, bChMap := generateBufferedChannelAndExchangeToMarketPrices(t, constants.Exchange1Exchange2Array)

	runPriceEncoderSequentially(
		t,
		constants.ExchangeFeedId1,
		etmp,
		bChMap[constants.ExchangeFeedId1],
		[]*types.MarketPriceTimestamp{
			constants.Market9_TimeT_Price1,
			constants.Market8_TimeTMinusThreshold_Price2,
		},
	)

	require.Len(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp, 2)
	require.Empty(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp)

	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price1,
			LastUpdateTime: constants.TimeT,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp[constants.MarketId9],
	)
	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price2,
			LastUpdateTime: constants.TimeTMinusThreshold,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp[constants.MarketId8],
	)
}

func TestPriceEncoder_WriteToOneMarketTwice(t *testing.T) {
	etmp, bChMap := generateBufferedChannelAndExchangeToMarketPrices(t, constants.Exchange1Exchange2Array)

	runPriceEncoderSequentially(
		t,
		constants.ExchangeFeedId1,
		etmp,
		bChMap[constants.ExchangeFeedId1],
		[]*types.MarketPriceTimestamp{
			constants.Market9_TimeTMinusThreshold_Price2,
			constants.Market9_TimeT_Price1,
		},
	)

	require.Len(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp, 1)
	require.Empty(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp)

	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price1,
			LastUpdateTime: constants.TimeT,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp[constants.MarketId9],
	)
}

func TestPriceEncoder_WriteToTwoExchanges(t *testing.T) {
	etmp, bChMap := generateBufferedChannelAndExchangeToMarketPrices(t, constants.Exchange1Exchange2Array)

	runPriceEncoderSequentially(
		t,
		constants.ExchangeFeedId1,
		etmp,
		bChMap[constants.ExchangeFeedId1],
		[]*types.MarketPriceTimestamp{
			constants.Market9_TimeT_Price1,
		},
	)

	runPriceEncoderSequentially(
		t,
		constants.ExchangeFeedId2,
		etmp,
		bChMap[constants.ExchangeFeedId2],
		[]*types.MarketPriceTimestamp{
			constants.Market8_TimeTMinusThreshold_Price2,
		},
	)

	require.Len(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp, 1)
	require.Len(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp, 1)

	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price1,
			LastUpdateTime: constants.TimeT,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp[constants.MarketId9],
	)
	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price2,
			LastUpdateTime: constants.TimeTMinusThreshold,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp[constants.MarketId8],
	)
}

func TestPriceEncoder_WriteToTwoExchangesConcurrentlyWithManyUpdates(t *testing.T) {
	etmp, bChMap := generateBufferedChannelAndExchangeToMarketPrices(t, constants.Exchange1Exchange2Array)

	largeMarketWrite := []*types.MarketPriceTimestamp{
		constants.Market8_TimeTMinusThreshold_Price1,
		constants.Market8_TimeTMinusThreshold_Price2,
		constants.Market8_TimeTMinusThreshold_Price3,
		constants.Market9_TimeTMinusThreshold_Price1,
		constants.Market9_TimeTMinusThreshold_Price2,
		constants.Market9_TimeTMinusThreshold_Price3,
		constants.Market8_TimeT_Price3,
		constants.Market9_TimeT_Price1,
		constants.Market9_TimeT_Price2,
		constants.Market9_TimeT_Price3,
		constants.Market9_TimeTPlusThreshold_Price1,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runPriceEncoderConcurrently(
			t,
			constants.ExchangeFeedId1,
			etmp,
			bChMap[constants.ExchangeFeedId1],
			largeMarketWrite,
		)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runPriceEncoderConcurrently(
			t,
			constants.ExchangeFeedId2,
			etmp,
			bChMap[constants.ExchangeFeedId2],
			largeMarketWrite,
		)
	}()

	wg.Wait()

	require.Len(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp, 2)
	require.Len(t, etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp, 2)

	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price1,
			LastUpdateTime: constants.TimeTPlusThreshold,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp[constants.MarketId9],
	)
	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price3,
			LastUpdateTime: constants.TimeT,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId1].MarketToPriceTimestamp[constants.MarketId8],
	)

	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price1,
			LastUpdateTime: constants.TimeTPlusThreshold,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp[constants.MarketId9],
	)
	require.Equal(
		t,
		&pricefeedtypes.PriceTimestamp{
			Price:          constants.Price3,
			LastUpdateTime: constants.TimeT,
		},
		etmp.ExchangeMarketPrices[constants.ExchangeFeedId2].MarketToPriceTimestamp[constants.MarketId8],
	)
}

func TestPriceUpdater_Mixed(t *testing.T) {
	tests := map[string]struct {
		// parameters
		exchangeAndMarketPrices []*client.ExchangeFeedIdMarketPriceTimestamp
		priceUpdateError        error

		// expectations
		expectedMarketPriceUpdate []*api.MarketPriceUpdate
	}{
		"Update throws": {
			// Throws error due to mock so that we can simulate fail state.
			exchangeAndMarketPrices: []*client.ExchangeFeedIdMarketPriceTimestamp{
				constants.ExchangeId1_Market9_TimeT_Price1,
			},
			priceUpdateError: errors.New("failed to send price update"),
		},
		"No exchange market prices, does not call `UpdateMarketPrices`": {
			exchangeAndMarketPrices: []*client.ExchangeFeedIdMarketPriceTimestamp{},
		},
		"One market for one exchange": {
			exchangeAndMarketPrices: []*client.ExchangeFeedIdMarketPriceTimestamp{
				constants.ExchangeId1_Market9_TimeT_Price1,
			},
			expectedMarketPriceUpdate: constants.Market9_SingleExchange_AtTimeUpdate,
		},
		"Three markets at timeT": {
			exchangeAndMarketPrices: []*client.ExchangeFeedIdMarketPriceTimestamp{
				constants.ExchangeId1_Market9_TimeT_Price1,
				constants.ExchangeId2_Market9_TimeT_Price2,
				constants.ExchangeId2_Market8_TimeT_Price2,
				constants.ExchangeId3_Market8_TimeT_Price3,
				constants.ExchangeId1_Market7_TimeT_Price1,
				constants.ExchangeId3_Market7_TimeT_Price3,
			},
			expectedMarketPriceUpdate: constants.AtTimeTPriceUpdate,
		},
		"Three markets at mixed time": {
			exchangeAndMarketPrices: []*client.ExchangeFeedIdMarketPriceTimestamp{
				constants.ExchangeId1_Market9_TimeT_Price1,
				constants.ExchangeId2_Market9_TimeT_Price2,
				constants.ExchangeId3_Market9_TimeT_Price3,
				constants.ExchangeId1_Market8_BeforeTimeT_Price3,
				constants.ExchangeId2_Market8_TimeT_Price2,
				constants.ExchangeId3_Market8_TimeT_Price3,
				constants.ExchangeId2_Market7_BeforeTimeT_Price1,
				constants.ExchangeId1_Market7_BeforeTimeT_Price3,
				constants.ExchangeId3_Market7_TimeT_Price3,
			},
			expectedMarketPriceUpdate: constants.MixedTimePriceUpdate,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Create `ExchangeFeedIdMarketPriceTimestamp` and populate it with market-price updates.
			etmp, _ := types.NewExchangeToMarketPrices(
				[]types.ExchangeFeedId{
					constants.ExchangeFeedId1,
					constants.ExchangeFeedId2,
					constants.ExchangeFeedId3,
				},
			)
			for _, exchangeAndMarketPrice := range tc.exchangeAndMarketPrices {
				etmp.UpdatePrice(
					exchangeAndMarketPrice.ExchangeFeedId,
					exchangeAndMarketPrice.MarketPriceTimestamp,
				)
			}

			// Create a mock `PriceFeedServiceClient` and run `RunPriceUpdaterTaskLoop`.
			mockPriceFeedClient := generateMockPriceFeedServiceClient()
			mockPriceFeedClient.On("UpdateMarketPrices", grpc.Ctx, mock.Anything).
				Return(nil, tc.priceUpdateError)

			err := pricefeedclient.RunPriceUpdaterTaskLoop(
				grpc.Ctx,
				etmp,
				mockPriceFeedClient,
				log.NewNopLogger(),
			)
			require.Equal(
				t,
				tc.priceUpdateError,
				err,
			)

			// We sort the `expectedUpdates` as ordering is not guaranteed.
			// We then verify `UpdateMarketPrices` was called with an update that, when sorted, matches
			// the sorted `expectedUpdates`.
			expectedUpdates := tc.expectedMarketPriceUpdate
			sortMarketPriceUpdateByMarketIdDescending(expectedUpdates)

			if tc.expectedMarketPriceUpdate != nil {
				mockPriceFeedClient.AssertCalled(
					t,
					"UpdateMarketPrices",
					grpc.Ctx,
					mock.MatchedBy(func(i interface{}) bool {
						param := i.(*api.UpdateMarketPricesRequest)
						updates := param.MarketPriceUpdates
						sortMarketPriceUpdateByMarketIdDescending(updates)

						for i, update := range updates {
							prices := update.ExchangePrices
							require.ElementsMatch(
								t,
								expectedUpdates[i].ExchangePrices,
								prices,
							)
						}
						return true
					}),
				)
			} else {
				mockPriceFeedClient.AssertNotCalled(t, "UpdateMarketPrices")
			}
		})
	}
}

// ----------------- Generate Mock Instances ----------------- //

func generateMockPriceFeedServiceClient() *mocks.QueryClient {
	mockPriceFeedServiceClient := &mocks.QueryClient{}

	return mockPriceFeedServiceClient
}

// ----------------- Helper Functions ----------------- //

func generateBufferedChannelAndExchangeToMarketPrices(
	t *testing.T,
	exchangeFeedIds []types.ExchangeFeedId,
) (
	*types.ExchangeToMarketPrices,
	map[types.ExchangeFeedId]chan *pricefeedclient.PriceFetcherSubtaskResponse,
) {
	etmp, _ := types.NewExchangeToMarketPrices(exchangeFeedIds)

	exhangeIdToBufferedChannel :=
		map[types.ExchangeFeedId]chan *pricefeedclient.PriceFetcherSubtaskResponse{}
	for _, exchangeFeedId := range exchangeFeedIds {
		bCh := make(chan *pricefeedclient.PriceFetcherSubtaskResponse, maxBufferedChannelLength)
		exhangeIdToBufferedChannel[exchangeFeedId] = bCh
	}

	require.Len(t, etmp.ExchangeMarketPrices, len(exchangeFeedIds))
	return etmp, exhangeIdToBufferedChannel
}

func runPriceEncoderSequentially(
	t *testing.T,
	exchangeFeedId types.ExchangeFeedId,
	etmp *types.ExchangeToMarketPrices,
	bCh chan *pricefeedclient.PriceFetcherSubtaskResponse,
	writes []*types.MarketPriceTimestamp,
) {
	// Make sure there are not more write than the `bufferedChannel` can hold.
	require.True(t, len(writes) <= maxBufferedChannelLength)

	for _, write := range writes {
		bCh <- &pricefeedclient.PriceFetcherSubtaskResponse{write, nil}
	}

	close(bCh)
	subTaskRunnerImpl.StartPriceEncoder(exchangeFeedId, etmp, log.NewNopLogger(), bCh)
}

func runPriceEncoderConcurrently(
	t *testing.T,
	exchangeFeedId types.ExchangeFeedId,
	etmp *types.ExchangeToMarketPrices,
	bCh chan *pricefeedclient.PriceFetcherSubtaskResponse,
	writes []*types.MarketPriceTimestamp,
) {
	// Start a `waitGroup` for the `PriceEncoder` which will complete when the `bufferedChannel`
	// is empty and is closed.
	var priceEncoderWg sync.WaitGroup
	priceEncoderWg.Add(1)
	go func() {
		defer priceEncoderWg.Done()
		subTaskRunnerImpl.StartPriceEncoder(exchangeFeedId, etmp, log.NewNopLogger(), bCh)
	}()

	// Start a `waitGroup` for threads that will write to the `bufferedChannel`.
	var writeWg sync.WaitGroup
	for _, write := range writes {
		writeWg.Add(1)
		go func(write *types.MarketPriceTimestamp) {
			defer writeWg.Done()
			bCh <- &pricefeedclient.PriceFetcherSubtaskResponse{write, nil}
		}(write)
	}

	writeWg.Wait()
	close(bCh)
	priceEncoderWg.Wait()
}

func sortMarketPriceUpdateByMarketIdDescending(
	marketPriceUpdate []*api.MarketPriceUpdate,
) {
	sort.Slice(
		marketPriceUpdate,
		func(i, j int) bool {
			return marketPriceUpdate[i].MarketId > marketPriceUpdate[j].MarketId
		},
	)
}
