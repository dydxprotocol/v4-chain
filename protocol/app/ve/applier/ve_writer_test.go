package ve_writer_test

import (
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	veapplier "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/applier"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vemath "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	vecache "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/vecache"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWritePricesToStoreAndMaybeCache(t *testing.T) {
	testHeight := int64(101)

	// Note that the difference between the intitial prices and the input prices
	// has to exceed the minimum price change for the price to be updated in the store.
	tests := map[string]struct {
		initialPrices         map[uint32]*pricestypes.MarketPrice
		inputPrices           map[string]voteweighted.AggregatorPricePair
		marketParams          []pricestypes.MarketParam
		round                 int32
		writeToCache          bool
		expectedError         error
		expectedPrices        map[uint32]*pricestypes.MarketPrice
		expectedCachedUpdates vecache.PriceUpdates
	}{
		"Valid prices, write to cache": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1500000), PnlPrice: big.NewInt(1500000)},
				"ETH-USD": {SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
				{Id: 1, Pair: "ETH-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedCachedUpdates: vecache.PriceUpdates{
				{MarketId: 0, SpotPrice: big.NewInt(1500000), PnlPrice: big.NewInt(1500000)},
				{MarketId: 1, SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
		},
		"Valid prices, don't write to cache": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1500000), PnlPrice: big.NewInt(1500000)},
				"ETH-USD": {SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
				{Id: 1, Pair: "ETH-USD"},
			},
			round:         1,
			writeToCache:  false,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedCachedUpdates: nil,
		},
		"Spot price change too small": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1000001), PnlPrice: big.NewInt(1500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1500000},
			},
			expectedCachedUpdates: vecache.PriceUpdates{
				{MarketId: 0, SpotPrice: nil, PnlPrice: big.NewInt(1500000)},
			},
		},
		"PnL price change too small": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1500000), PnlPrice: big.NewInt(1000001)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1000000},
			},
			expectedCachedUpdates: vecache.PriceUpdates{
				{MarketId: 0, SpotPrice: big.NewInt(1500000), PnlPrice: nil},
			},
		},
		"Spot and pnl price change too small": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1000001), PnlPrice: big.NewInt(1000001)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			expectedCachedUpdates: nil,
		},
		"Spot price change too small for BTC-USD, valid for ETH-USD": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1000001), PnlPrice: big.NewInt(1500000)},
				"ETH-USD": {SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
				{Id: 1, Pair: "ETH-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedCachedUpdates: vecache.PriceUpdates{
				{MarketId: 0, SpotPrice: nil, PnlPrice: big.NewInt(1500000)},
				{MarketId: 1, SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
		},
		"PnL price change too small for ETH-USD, valid for BTC-USD": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1500000), PnlPrice: big.NewInt(1500000)},
				"ETH-USD": {SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2000001)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
				{Id: 1, Pair: "ETH-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2000000},
			},
			expectedCachedUpdates: vecache.PriceUpdates{
				{MarketId: 0, SpotPrice: big.NewInt(1500000), PnlPrice: big.NewInt(1500000)},
				{MarketId: 1, SpotPrice: big.NewInt(2500000), PnlPrice: nil},
			},
		},
		"Spot and PnL price change too small for BTC-USD, valid for ETH-USD": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(1000001), PnlPrice: big.NewInt(1000001)},
				"ETH-USD": {SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
				{Id: 1, Pair: "ETH-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedCachedUpdates: vecache.PriceUpdates{
				{MarketId: 1, SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
		},
		"Negative spot price": { // Note that pnl price cannot be negative if spot is not negative
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(-1500000), PnlPrice: big.NewInt(1500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			expectedCachedUpdates: nil,
		},
		"Negative spot and pnl price": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(-1500000), PnlPrice: big.NewInt(-1500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			expectedCachedUpdates: nil,
		},
		"Negative spot price for BTC, valid prices for ETH": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(-1500000), PnlPrice: big.NewInt(1500000)},
				"ETH-USD": {SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
				{Id: 1, Pair: "ETH-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedCachedUpdates: vecache.PriceUpdates{
				{MarketId: 1, SpotPrice: big.NewInt(2500000), PnlPrice: big.NewInt(2500000)},
			},
		},
		"Negative spot and pnl price for BTC, negative spot price for ETH": {
			initialPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			inputPrices: map[string]voteweighted.AggregatorPricePair{
				"BTC-USD": {SpotPrice: big.NewInt(-1500000), PnlPrice: big.NewInt(-1500000)},
				"ETH-USD": {SpotPrice: big.NewInt(-2500000), PnlPrice: big.NewInt(2500000)},
			},
			marketParams: []pricestypes.MarketParam{
				{Id: 0, Pair: "BTC-USD"},
				{Id: 1, Pair: "ETH-USD"},
			},
			round:         1,
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			expectedCachedUpdates: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			voteCodec := vecodec.NewDefaultVoteExtensionCodec()
			extCodec := vecodec.NewDefaultExtendedCommitCodec()
			voteAggregator := &mocks.VoteAggregator{}
			ctx, _, pricesKeeper, _, _, _, _, ratelimitKeeper, _, _ := keepertest.SubaccountsKeepers(t, false)

			keepertest.CreateTestMarkets(t, ctx, pricesKeeper)

			for id, price := range tc.initialPrices {
				err := pricesKeeper.UpdateSpotAndPnlMarketPrices(ctx, &pricestypes.MarketPriceUpdate{
					MarketId:  uint32(id),
					SpotPrice: price.SpotPrice,
					PnlPrice:  price.PnlPrice,
				})
				require.NoError(t, err)
			}

			veApplier := veapplier.NewVEApplier(
				log.NewNopLogger(),
				voteAggregator,
				pricesKeeper,
				ratelimitKeeper,
				voteCodec,
				extCodec,
			)

			ctx = ctx.WithBlockHeight(testHeight)
			err := veApplier.WritePricesToStoreAndMaybeCache(ctx, tc.inputPrices, tc.round, tc.writeToCache)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			for id, expectedPrice := range tc.expectedPrices {
				actualPrice, err := pricesKeeper.GetMarketPrice(ctx, id)
				require.NoError(t, err)
				require.Equal(t, expectedPrice.Id, actualPrice.Id)
				require.Equal(t, expectedPrice.SpotPrice, actualPrice.SpotPrice)
				require.Equal(t, expectedPrice.PnlPrice, actualPrice.PnlPrice)
			}

			actualCachedUpdates := veApplier.GetCachedPrices()
			require.Equal(t, tc.expectedCachedUpdates, actualCachedUpdates)
		})
	}
}

func TestWriteSDaiConversionRateToStoreAndMaybeCache(t *testing.T) {
	testHeight := int64(101)

	tests := map[string]struct {
		initialSDaiPrice            *big.Int
		initialLastBlockUpdated     *big.Int
		sDaiConversionRate          *big.Int
		round                       int32
		writeToCache                bool
		expectedError               error
		expectedSDaiPrice           *big.Int
		expectedLastBlockUpdated    *big.Int
		expectedCacheConversionRate *big.Int
		expectedCacheBlockHeight    *big.Int
	}{
		"Valid conversion rate": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          big.NewInt(1000000),
			round:                       1,
			writeToCache:                true,
			expectedError:               nil,
			expectedSDaiPrice:           big.NewInt(1000000),
			expectedLastBlockUpdated:    big.NewInt(101), // Assuming current block height is 101
			expectedCacheConversionRate: big.NewInt(1000000),
			expectedCacheBlockHeight:    big.NewInt(testHeight),
		},
		"Nil conversion rate": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          nil,
			round:                       1,
			writeToCache:                true,
			expectedError:               nil,
			expectedSDaiPrice:           big.NewInt(500000),
			expectedLastBlockUpdated:    big.NewInt(100),
			expectedCacheConversionRate: nil,
			expectedCacheBlockHeight:    nil,
		},
		"Negative conversion rate": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          big.NewInt(-1000000),
			round:                       1,
			writeToCache:                true,
			expectedError:               fmt.Errorf("sDAI conversion rate cannot be negative: -1000000"),
			expectedSDaiPrice:           big.NewInt(500000),
			expectedLastBlockUpdated:    big.NewInt(100),
			expectedCacheConversionRate: nil,
			expectedCacheBlockHeight:    nil,
		},
		"Zero conversion rate": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          big.NewInt(0),
			round:                       1,
			writeToCache:                true,
			expectedError:               fmt.Errorf("sDAI conversion rate cannot be zero"),
			expectedSDaiPrice:           big.NewInt(500000),
			expectedLastBlockUpdated:    big.NewInt(100),
			expectedCacheConversionRate: nil,
			expectedCacheBlockHeight:    nil,
		},
		"Valid conversion rate, don't write to cache": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          big.NewInt(2000000),
			round:                       2,
			writeToCache:                false,
			expectedError:               nil,
			expectedSDaiPrice:           big.NewInt(2000000),
			expectedLastBlockUpdated:    big.NewInt(testHeight),
			expectedCacheConversionRate: nil,
			expectedCacheBlockHeight:    nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			voteCodec := vecodec.NewDefaultVoteExtensionCodec()
			extCodec := vecodec.NewDefaultExtendedCommitCodec()
			voteAggregator := &mocks.VoteAggregator{}
			ctx, _, pricesKeeper, _, _, _, _, ratelimitKeeper, _, _ := keepertest.SubaccountsKeepers(t, false)

			ratelimitKeeper.SetSDAIPrice(ctx, tc.initialSDaiPrice)
			ratelimitKeeper.SetSDAILastBlockUpdated(ctx, tc.initialLastBlockUpdated)

			veApplier := veapplier.NewVEApplier(
				log.NewNopLogger(),
				voteAggregator,
				pricesKeeper,
				ratelimitKeeper,
				voteCodec,
				extCodec,
			)

			ctx = ctx.WithBlockHeight(testHeight)
			err := veApplier.WriteSDaiConversionRateToStoreAndMaybeCache(ctx, tc.sDaiConversionRate, tc.round, tc.writeToCache)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}

			// Check ratelimitKeeper state
			actualSDaiPrice, found := ratelimitKeeper.GetSDAIPrice(ctx)
			require.True(t, found)
			require.Equal(t, tc.expectedSDaiPrice, actualSDaiPrice)

			actualLastBlockUpdated, found := ratelimitKeeper.GetSDAILastBlockUpdated(ctx)
			require.True(t, found)
			require.Equal(t, tc.expectedLastBlockUpdated, actualLastBlockUpdated)

			actualCacheConversionRate, actualCacheBlockHeight := veApplier.GetCachedSDaiConversionRate()
			require.Equal(t, tc.expectedCacheConversionRate, actualCacheConversionRate)
			require.Equal(t, tc.expectedCacheBlockHeight, actualCacheBlockHeight)
		})
	}
}

func TestVEWriter(t *testing.T) {
	voteCodec := vecodec.NewDefaultVoteExtensionCodec()
	extCodec := vecodec.NewDefaultExtendedCommitCodec()

	voteAggregator := &mocks.VoteAggregator{}

	ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

	pricesKeeper := &mocks.VEApplierPricesKeeper{}
	pricesKeeper.On("PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything).Return(true, true)

	ratelimitKeeper := &mocks.VEApplierRatelimitKeeper{}
	ratelimitKeeper.On("SetSDAIPrice", mock.Anything, mock.Anything).Return()
	ratelimitKeeper.On("SetSDAILastBlockUpdated", mock.Anything, mock.Anything).Return()

	veApplier := veapplier.NewVEApplier(
		log.NewNopLogger(),
		voteAggregator,
		pricesKeeper,
		ratelimitKeeper,
		voteCodec,
		extCodec,
	)

	t.Run("if extracting oracle votes fails, fail", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		err := veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{[]byte("garbage"), {1, 2, 3, 4}, {1, 2, 3, 4}},
		}, true)

		priceUpdates := veApplier.GetCachedPrices()

		cachedPrices := make(map[string]ve.VEPricePair)
		for _, priceUpdate := range priceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = ve.VEPricePair{
				SpotPrice: priceUpdate.SpotPrice.Uint64(),
				PnlPrice:  priceUpdate.PnlPrice.Uint64(),
			}
		}

		require.Error(t, err)
		require.Equal(t, cachedPrices, make(map[string]ve.VEPricePair))
	})

	t.Run("if vote aggregation fails, fail", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(2)
		prices := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: []byte("price1"),
				PnlPrice:  []byte("price1"),
			},
		}

		_, extCommitInfoBz, err := vetesting.CreateSingleValidatorExtendedCommitInfo(
			constants.AliceConsAddress,
			prices,
			"",
		)
		require.NoError(t, err)

		// fail vote aggregation
		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(nil, nil, fmt.Errorf("fail")).Once()

		err = veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		}, true)

		priceUpdates := veApplier.GetCachedPrices()

		cachedPrices := make(map[string]ve.VEPricePair)
		for _, priceUpdate := range priceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = ve.VEPricePair{
				SpotPrice: priceUpdate.SpotPrice.Uint64(),
				PnlPrice:  priceUpdate.PnlPrice.Uint64(),
			}
		}

		require.Error(t, err)
		require.Equal(t, cachedPrices, make(map[string]ve.VEPricePair))
	})

	t.Run("ignore negative prices", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(3)

		priceBz := big.NewInt(-100).Bytes()

		prices := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: priceBz,
				PnlPrice:  priceBz,
			},
		}

		_, extCommitInfoBz, err := vetesting.CreateSingleValidatorExtendedCommitInfo(
			constants.AliceConsAddress,
			prices,
			"",
		)
		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]vemath.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(-100),
				PnlPrice:  big.NewInt(-100),
			},
		}, nil, nil)

		pricesKeeper.On("GetAllMarketParams", ctx).Return(
			[]pricestypes.MarketParam{
				{
					Id:   1,
					Pair: constants.BtcUsdPair,
				},
			},
		)

		pricesKeeper.On("GetMarketParam", ctx, uint32(1)).Return(
			pricestypes.MarketParam{
				Id:   1,
				Pair: constants.BtcUsdPair,
			},
			true,
		)

		err = veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		}, true)

		require.NoError(t, err)
	})

	t.Run("update prices in state", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(4)

		price1Bz := big.NewInt(100).Bytes()
		price2Bz := big.NewInt(200).Bytes()

		prices1 := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: price1Bz,
				PnlPrice:  price1Bz,
			},
		}

		prices2 := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: price2Bz,
				PnlPrice:  price2Bz,
			},
		}

		vote1, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceConsAddress,
				prices1,
				"",
			),
		)
		require.NoError(t, err)

		vote2, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobConsAddress,
				prices2,
				"",
			),
		)
		require.NoError(t, err)

		_, extCommitInfoBz, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote1, vote2},
		)
		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices1,
				},
				ConsAddress: constants.AliceConsAddress,
			},
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices2,
				},
				ConsAddress: constants.BobConsAddress,
			},
		}).Return(map[string]vemath.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(150),
				PnlPrice:  big.NewInt(150),
			},
		}, nil, nil)

		pricesKeeper.On("GetAllMarketParams", ctx).Return(
			[]pricestypes.MarketParam{
				{
					Id:   1,
					Pair: constants.BtcUsdPair,
				},
			},
		)

		pricesKeeper.On("GetMarketParam", ctx, uint32(1)).Return(
			pricestypes.MarketParam{
				Id:   1,
				Pair: constants.BtcUsdPair,
			},
			true,
		)

		pricesKeeper.On("UpdateSpotAndPnlMarketPrices", ctx, mock.Anything).Return(nil)

		err = veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: []cometabci.VoteInfo{},
			},
		}, true)

		priceUpdates := veApplier.GetCachedPrices()

		cachedPrices := make(map[string]ve.VEPricePair)
		for _, priceUpdate := range priceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = ve.VEPricePair{
				SpotPrice: priceUpdate.SpotPrice.Uint64(),
				PnlPrice:  priceUpdate.PnlPrice.Uint64(),
			}
		}

		require.NoError(t, err)
		require.Equal(t, map[string]ve.VEPricePair{
			constants.BtcUsdPair: {
				SpotPrice: 150,
				PnlPrice:  150,
			},
		}, cachedPrices)
	})

	t.Run("doesn't update prices for same round and height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(5)

		price1Bz := big.NewInt(100).Bytes()
		price2Bz := big.NewInt(200).Bytes()

		prices1 := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: price1Bz,
				PnlPrice:  price1Bz,
			},
		}

		prices2 := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: price2Bz,
				PnlPrice:  price2Bz,
			},
		}

		vote1, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceConsAddress,
				prices1,
				"",
			),
		)
		require.NoError(t, err)

		vote2, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobConsAddress,
				prices2,
				"",
			),
		)
		require.NoError(t, err)

		_, extCommitInfoBz1, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote1},
		)
		require.NoError(t, err)

		_, extCommitInfoBz2, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote2},
		)
		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices1,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]vemath.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(100),
				PnlPrice:  big.NewInt(100),
			},
		}, nil, nil).Once()

		pricesKeeper.On("GetAllMarketParams", ctx).Return(
			[]pricestypes.MarketParam{
				{
					Id:   1,
					Pair: constants.BtcUsdPair,
				},
			},
		).Twice()

		pricesKeeper.On("GetMarketParam", ctx, uint32(1)).Return(
			pricestypes.MarketParam{
				Id:   1,
				Pair: constants.BtcUsdPair,
			},
			true,
		).Twice()

		pricesKeeper.On("UpdateSpotAndPnlMarketPrices", ctx, mock.Anything).Return(nil).Twice()

		// First call
		err = veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: []cometabci.VoteInfo{},
			},
		}, true)
		require.NoError(t, err)

		priceUpdates := veApplier.GetCachedPrices()
		cachedPrices := make(map[string]ve.VEPricePair)
		for _, priceUpdate := range priceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = ve.VEPricePair{
				SpotPrice: priceUpdate.SpotPrice.Uint64(),
				PnlPrice:  priceUpdate.PnlPrice.Uint64(),
			}
		}

		require.Equal(t, map[string]ve.VEPricePair{
			constants.BtcUsdPair: {
				SpotPrice: 100,
				PnlPrice:  100,
			},
		}, cachedPrices)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices2,
				},
				ConsAddress: constants.BobConsAddress,
			},
		}).Return(map[string]vemath.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(200),
				PnlPrice:  big.NewInt(200),
			},
		}, nil, nil).Once()

		// Second call with the same round and height
		err = veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz2, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: []cometabci.VoteInfo{},
			},
		}, true)
		require.NoError(t, err)

		priceUpdates = veApplier.GetCachedPrices()
		cachedPrices = make(map[string]ve.VEPricePair)
		for _, priceUpdate := range priceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = ve.VEPricePair{
				SpotPrice: priceUpdate.SpotPrice.Uint64(),
				PnlPrice:  priceUpdate.PnlPrice.Uint64(),
			}
		}

		// Ensure the cached prices are still the same as the first call
		require.Equal(t, map[string]ve.VEPricePair{
			constants.BtcUsdPair: {
				SpotPrice: 100,
				PnlPrice:  100,
			},
		}, cachedPrices)
	})
	t.Run("correctly updates prices in cache", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(5)

		price1Bz := big.NewInt(100).Bytes()
		price2Bz := big.NewInt(200).Bytes()

		prices1 := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: price1Bz,
				PnlPrice:  price1Bz,
			},
		}

		prices2 := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: price2Bz,
				PnlPrice:  price2Bz,
			},
		}

		vote1, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceConsAddress,
				prices1,
				"",
			),
		)
		require.NoError(t, err)

		vote2, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobConsAddress,
				prices2,
				"",
			),
		)
		require.NoError(t, err)

		_, extCommitInfoBz1, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote1},
		)
		require.NoError(t, err)

		_, extCommitInfoBz2, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote2},
		)
		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices1,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]vemath.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(100),
				PnlPrice:  big.NewInt(100),
			},
		}, nil, nil).Once()

		pricesKeeper.On("GetAllMarketParams", ctx).Return(
			[]pricestypes.MarketParam{
				{
					Id:   1,
					Pair: constants.BtcUsdPair,
				},
			},
		).Twice()

		pricesKeeper.On("GetMarketParam", ctx, uint32(1)).Return(
			pricestypes.MarketParam{
				Id:   1,
				Pair: constants.BtcUsdPair,
			},
			true,
		).Twice()

		pricesKeeper.On("UpdateSpotAndPnlMarketPrices", ctx, mock.Anything).Return(nil).Twice()

		// First call
		err = veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: []cometabci.VoteInfo{},
			},
		}, true)
		require.NoError(t, err)

		priceUpdates := veApplier.GetCachedPrices()
		cachedPrices := make(map[string]ve.VEPricePair)
		for _, priceUpdate := range priceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = ve.VEPricePair{
				SpotPrice: priceUpdate.SpotPrice.Uint64(),
				PnlPrice:  priceUpdate.PnlPrice.Uint64(),
			}
		}

		require.Equal(t, map[string]ve.VEPricePair{
			constants.BtcUsdPair: {
				SpotPrice: 100,
				PnlPrice:  100,
			},
		}, cachedPrices)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices2,
				},
				ConsAddress: constants.BobConsAddress,
			},
		}).Return(map[string]vemath.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(200),
				PnlPrice:  big.NewInt(200),
			},
		}, nil, nil).Once()

		// Second call with different round
		err = veApplier.ApplyVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz2, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 2,
				Votes: []cometabci.VoteInfo{},
			},
		}, true)
		require.NoError(t, err)

		priceUpdates = veApplier.GetCachedPrices()
		cachedPrices = make(map[string]ve.VEPricePair)
		for _, priceUpdate := range priceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = ve.VEPricePair{
				SpotPrice: priceUpdate.SpotPrice.Uint64(),
				PnlPrice:  priceUpdate.PnlPrice.Uint64(),
			}
		}

		// Ensure the cached prices change
		require.Equal(t, map[string]ve.VEPricePair{
			constants.BtcUsdPair: {
				SpotPrice: 200,
				PnlPrice:  200,
			},
		}, cachedPrices)
	})
}
