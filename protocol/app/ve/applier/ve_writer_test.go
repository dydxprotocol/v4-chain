package ve_writer_test

import (
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	veapplier "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/applier"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	bigintcache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/bigintcache"
	pricecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/pricecache"
	vecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/vecache"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/memclob"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
)

func TestWritePricesToStoreAndMaybeCache(t *testing.T) {
	testHeight := int64(101)

	// Note that the difference between the intitial prices and the input prices
	// has to exceed the minimum price change for the price to be updated in the store.
	tests := map[string]struct {
		initialPrices             map[uint32]*pricestypes.MarketPrice
		inputPrices               map[string]voteweighted.AggregatorPricePair
		marketParams              []pricestypes.MarketParam
		writeToCache              bool
		expectedError             error
		expectedPrices            map[uint32]*pricestypes.MarketPrice
		expectedSpotCachedUpdates pricecache.PriceUpdates
		expectedPnlCachedUpdates  pricecache.PriceUpdates
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedSpotCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 0, Price: big.NewInt(1500000)},
				{MarketId: 1, Price: big.NewInt(2500000)},
			},
			expectedPnlCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 0, Price: big.NewInt(1500000)},
				{MarketId: 1, Price: big.NewInt(2500000)},
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
			writeToCache:  false,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedSpotCachedUpdates: nil,
			expectedPnlCachedUpdates:  nil,
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1500000},
			},
			expectedSpotCachedUpdates: nil,
			expectedPnlCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 0, Price: big.NewInt(1500000)},
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1000000},
			},
			expectedSpotCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 0, Price: big.NewInt(1500000)},
			},
			expectedPnlCachedUpdates: nil,
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			expectedSpotCachedUpdates: nil,
			expectedPnlCachedUpdates:  nil,
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedSpotCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 1, Price: big.NewInt(2500000)},
			},
			expectedPnlCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 0, Price: big.NewInt(1500000)},
				{MarketId: 1, Price: big.NewInt(2500000)},
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1500000, PnlPrice: 1500000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2000000},
			},
			expectedSpotCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 0, Price: big.NewInt(1500000)},
				{MarketId: 1, Price: big.NewInt(2500000)},
			},
			expectedPnlCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 0, Price: big.NewInt(1500000)},
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedSpotCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 1, Price: big.NewInt(2500000)},
			},
			expectedPnlCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 1, Price: big.NewInt(2500000)},
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			expectedSpotCachedUpdates: nil,
			expectedPnlCachedUpdates:  nil,
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
			},
			expectedSpotCachedUpdates: nil,
			expectedPnlCachedUpdates:  nil,
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2500000, PnlPrice: 2500000},
			},
			expectedSpotCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 1, Price: big.NewInt(2500000)},
			},
			expectedPnlCachedUpdates: pricecache.PriceUpdates{
				{MarketId: 1, Price: big.NewInt(2500000)},
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
			writeToCache:  true,
			expectedError: nil,
			expectedPrices: map[uint32]*pricestypes.MarketPrice{
				0: {Id: 0, SpotPrice: 1000000, PnlPrice: 1000000},
				1: {Id: 1, SpotPrice: 2000000, PnlPrice: 2000000},
			},
			expectedSpotCachedUpdates: nil,
			expectedPnlCachedUpdates:  nil,
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

			spotPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
			pnlPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
			sDaiConversionRateCache := bigintcache.BigIntCacheImpl{}
			veCache := vecache.NewVECache()

			veApplier := veapplier.NewVEApplier(
				log.NewNopLogger(),
				voteAggregator,
				pricesKeeper,
				ratelimitKeeper,
				voteCodec,
				extCodec,
				&spotPriceUpdateCache,
				&pnlPriceUpdateCache,
				&sDaiConversionRateCache,
				veCache,
			)

			ctx = ctx.WithBlockHeight(testHeight)
			err := veApplier.WritePricesToStoreAndMaybeCache(ctx, tc.inputPrices, []byte{}, tc.writeToCache)

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

			actualSpotCachedUpdates := veApplier.GetSpotPriceUpdateCache().GetPriceUpdates()
			require.Equal(t, tc.expectedSpotCachedUpdates, actualSpotCachedUpdates)

			actualPnlCachedUpdates := veApplier.GetPnlPriceUpdateCache().GetPriceUpdates()
			require.Equal(t, tc.expectedPnlCachedUpdates, actualPnlCachedUpdates)
		})
	}
}

func TestWriteSDaiConversionRateToStoreAndMaybeCache(t *testing.T) {
	testHeight := int64(6000)

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
	}{
		"10^27 is valid but not written to state but is written to cache": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil),
			round:                       1,
			writeToCache:                true,
			expectedError:               nil,
			expectedSDaiPrice:           big.NewInt(500000),
			expectedLastBlockUpdated:    big.NewInt(100),
			expectedCacheConversionRate: new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil),
		},
		"Valid conversion rate": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          new(big.Int).Exp(big.NewInt(11), big.NewInt(27), nil),
			round:                       1,
			writeToCache:                true,
			expectedError:               nil,
			expectedSDaiPrice:           new(big.Int).Exp(big.NewInt(11), big.NewInt(27), nil),
			expectedLastBlockUpdated:    big.NewInt(testHeight),
			expectedCacheConversionRate: new(big.Int).Exp(big.NewInt(11), big.NewInt(27), nil),
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
		},
		"Negative conversion rate": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          big.NewInt(-1000000),
			round:                       1,
			writeToCache:                true,
			expectedError:               fmt.Errorf("invalid sDAI conversion rate: -1000000"),
			expectedSDaiPrice:           big.NewInt(500000),
			expectedLastBlockUpdated:    big.NewInt(100),
			expectedCacheConversionRate: nil,
		},
		"Zero conversion rate": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          big.NewInt(0),
			round:                       1,
			writeToCache:                true,
			expectedError:               fmt.Errorf("invalid sDAI conversion rate: 0"),
			expectedSDaiPrice:           big.NewInt(500000),
			expectedLastBlockUpdated:    big.NewInt(100),
			expectedCacheConversionRate: nil,
		},
		"Valid conversion rate, don't write to cache": {
			initialSDaiPrice:            big.NewInt(500000),
			initialLastBlockUpdated:     big.NewInt(100),
			sDaiConversionRate:          new(big.Int).Exp(big.NewInt(11), big.NewInt(27), nil),
			round:                       2,
			writeToCache:                false,
			expectedError:               nil,
			expectedSDaiPrice:           new(big.Int).Exp(big.NewInt(11), big.NewInt(27), nil),
			expectedLastBlockUpdated:    big.NewInt(testHeight),
			expectedCacheConversionRate: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			voteCodec := vecodec.NewDefaultVoteExtensionCodec()
			extCodec := vecodec.NewDefaultExtendedCommitCodec()
			voteAggregator := &mocks.VoteAggregator{}
			ctx, _, pricesKeeper, _, _, bankKeeper, _, ratelimitKeeper, _, _ := keepertest.SubaccountsKeepers(t, false)

			ratelimitKeeper.SetSDAIPrice(ctx, tc.initialSDaiPrice)
			ratelimitKeeper.SetSDAILastBlockUpdated(ctx, tc.initialLastBlockUpdated)

			spotPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
			pnlPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
			sDaiConversionRateCache := bigintcache.BigIntCacheImpl{}
			veCache := vecache.NewVECache()

			veApplier := veapplier.NewVEApplier(
				log.NewNopLogger(),
				voteAggregator,
				pricesKeeper,
				ratelimitKeeper,
				voteCodec,
				extCodec,
				&spotPriceUpdateCache,
				&pnlPriceUpdateCache,
				&sDaiConversionRateCache,
				veCache,
			)

			ctx = ctx.WithBlockHeight(testHeight)

			// Set test tdai balance
			initialTestAmountTDai := sdkmath.NewIntFromBigInt(big.NewInt(1))
			tDaiToMintCoins := sdk.NewCoins(sdk.NewCoin(types.TDaiDenom, initialTestAmountTDai))
			err := bankKeeper.MintCoins(ctx, types.TDaiPoolAccount, tDaiToMintCoins)
			require.NoError(t, err)

			initialTestAmountSDai := sdkmath.NewIntFromBigInt(ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr("100000000000000000000000000000000000000000000"))
			sDaiToMintCoins := sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, initialTestAmountSDai))
			err = bankKeeper.MintCoins(ctx, types.TDaiPoolAccount, sDaiToMintCoins)
			require.NoError(t, err)

			err = veApplier.WriteSDaiConversionRateToStoreAndMaybeCache(ctx, tc.sDaiConversionRate, []byte{}, tc.writeToCache)

			if tc.expectedError != nil {
				require.EqualError(t, err, tc.expectedError.Error())
				actualTDaiBalance := bankKeeper.GetSupply(ctx, types.TDaiDenom)
				require.Equal(t, 0, actualTDaiBalance.Amount.BigInt().Cmp(initialTestAmountTDai.BigInt()))
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

			actualCacheConversionRate := veApplier.GetCachedSDaiConversionRate()
			require.Equal(t, tc.expectedCacheConversionRate, actualCacheConversionRate)
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
	ratelimitKeeper.On("ProcessNewSDaiConversionRateUpdate", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	spotPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
	pnlPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
	sDaiConversionRateCache := bigintcache.BigIntCacheImpl{}
	veCache := vecache.NewVECache()

	veApplier := veapplier.NewVEApplier(
		log.NewNopLogger(),
		voteAggregator,
		pricesKeeper,
		ratelimitKeeper,
		voteCodec,
		extCodec,
		&spotPriceUpdateCache,
		&pnlPriceUpdateCache,
		&sDaiConversionRateCache,
		veCache,
	)

	t.Run("if extracting oracle votes fails, fail", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		err := veApplier.ApplyVE(
			ctx,
			[][]byte{[]byte("garbage"), {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)

		spotPriceUpdates := veApplier.GetCachedSpotPrices()
		pnlPriceUpdates := veApplier.GetCachedPnlPrices()

		cachedSpotPrices := make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices := make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}
		require.Error(t, err)
		require.Equal(t, cachedSpotPrices, make(map[string]uint64))
		require.Equal(t, cachedPnlPrices, make(map[string]uint64))
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

		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)

		spotPriceUpdates := veApplier.GetCachedSpotPrices()
		pnlPriceUpdates := veApplier.GetCachedPnlPrices()

		cachedSpotPrices := make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices := make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}
		require.Error(t, err)
		require.Equal(t, cachedSpotPrices, make(map[string]uint64))
		require.Equal(t, cachedPnlPrices, make(map[string]uint64))
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
		}).Return(map[string]voteweighted.AggregatorPricePair{
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

		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)

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
		}).Return(map[string]voteweighted.AggregatorPricePair{
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

		pricesKeeper.On("UpdateSpotPrice", ctx, mock.Anything).Return(nil)
		pricesKeeper.On("UpdatePnlPrice", ctx, mock.Anything).Return(nil)

		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)

		spotPriceUpdates := veApplier.GetCachedSpotPrices()
		pnlPriceUpdates := veApplier.GetCachedPnlPrices()

		cachedSpotPrices := make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices := make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		require.NoError(t, err)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 150}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 150}, cachedPnlPrices)

		pricesKeeper.AssertCalled(t, "PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything)
	})

	t.Run("doesn't update prices for cached values", func(t *testing.T) {
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

		_, extCommitInfoBz1, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote1},
		)
		require.NoError(t, err)

		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices1,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]voteweighted.AggregatorPricePair{
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
		).Times(4)

		pricesKeeper.On("UpdateSpotPrice", ctx, mock.Anything).Return(nil).Twice()
		pricesKeeper.On("UpdatePnlPrice", ctx, mock.Anything).Return(nil).Twice()

		// First call
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates := veApplier.GetCachedSpotPrices()
		pnlPriceUpdates := veApplier.GetCachedPnlPrices()

		cachedSpotPrices := make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices := make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedPnlPrices)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices2,
				},
				ConsAddress: constants.BobConsAddress,
			},
		}).Return(map[string]voteweighted.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(200),
				PnlPrice:  big.NewInt(200),
			},
		}, nil, nil).Once()

		// Second call with the same round and height
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates = veApplier.GetCachedSpotPrices()
		pnlPriceUpdates = veApplier.GetCachedPnlPrices()

		cachedSpotPrices = make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices = make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		// Ensure the cached prices are still the same as the first call
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedPnlPrices)

		pricesKeeper.AssertCalled(t, "PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything)
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
		}).Return(map[string]voteweighted.AggregatorPricePair{
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
		).Times(4)

		pricesKeeper.On("UpdateSpotPrice", ctx, mock.Anything).Return(nil).Twice()
		pricesKeeper.On("UpdatePnlPrice", ctx, mock.Anything).Return(nil).Twice()

		// First call
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates := veApplier.GetCachedSpotPrices()
		pnlPriceUpdates := veApplier.GetCachedPnlPrices()

		cachedSpotPrices := make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices := make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedPnlPrices)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices2,
				},
				ConsAddress: constants.BobConsAddress,
			},
		}).Return(map[string]voteweighted.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(200),
				PnlPrice:  big.NewInt(200),
			},
		}, nil, nil).Once()

		// Second call with different round
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz2, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates = veApplier.GetCachedSpotPrices()
		pnlPriceUpdates = veApplier.GetCachedPnlPrices()

		cachedSpotPrices = make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices = make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		// Ensure the cached prices change
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 200}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 200}, cachedPnlPrices)

		pricesKeeper.AssertCalled(t, "PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything)
	})

	t.Run("Different txHash in cache when trying to write to store sdai rate", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(5)

		rate1, ok := big.NewInt(0).SetString("123456789000000000000000000000", 10)
		require.True(t, ok)

		rate2, ok := big.NewInt(0).SetString("111111111100000000000000000000", 10)
		require.True(t, ok)

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
				"123456789000000000000000000000",
			),
		)
		require.NoError(t, err)

		vote3, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceConsAddress,
				prices1,
				"111111111100000000000000000000",
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

		_, extCommitInfoBz3, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote3},
		)
		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices:             prices1,
					SDaiConversionRate: "123456789000000000000000000000",
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]voteweighted.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(100),
				PnlPrice:  big.NewInt(100),
			},
		}, rate1, nil).Once()

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices:             prices1,
					SDaiConversionRate: "111111111100000000000000000000",
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]voteweighted.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(100),
				PnlPrice:  big.NewInt(100),
			},
		}, rate2, nil).Once()

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
		).Times(4)

		pricesKeeper.On("UpdateSpotPrice", ctx, mock.Anything).Return(nil).Times(3)
		pricesKeeper.On("UpdatePnlPrice", ctx, mock.Anything).Return(nil).Times(3)

		// First call
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates := veApplier.GetCachedSpotPrices()
		pnlPriceUpdates := veApplier.GetCachedPnlPrices()

		cachedSpotPrices := make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices := make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedPnlPrices)

		cachedSDAIRate := veApplier.GetCachedSDaiConversionRate()
		require.Equal(t, rate1, cachedSDAIRate)

		// Second call should be using the cache and so not change the sdai rate
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true)
		require.NoError(t, err)

		cachedSDAIRate = veApplier.GetCachedSDaiConversionRate()
		require.Equal(t, rate1, cachedSDAIRate)

		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz3, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		cachedSDAIRate = veApplier.GetCachedSDaiConversionRate()
		require.Equal(t, rate2, cachedSDAIRate)

		ctx = ctx.WithBlockHeight(6)

		pricesKeeper.On("UpdateSpotPrice", ctx, mock.Anything).Return(nil).Times(2)
		pricesKeeper.On("UpdatePnlPrice", ctx, mock.Anything).Return(nil).Times(2)

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
		).Times(2)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices2,
				},
				ConsAddress: constants.BobConsAddress,
			},
		}).Return(map[string]voteweighted.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(200),
				PnlPrice:  big.NewInt(200),
			},
		}, nil, nil).Once()

		// Second call with different round
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz2, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates = veApplier.GetCachedSpotPrices()
		pnlPriceUpdates = veApplier.GetCachedPnlPrices()

		cachedSpotPrices = make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices = make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		// Ensure the cached prices change
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 200}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 200}, cachedPnlPrices)

		// Second call
		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz2, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)
		ratelimitKeeper.AssertNumberOfCalls(t, "ProcessNewSDaiConversionRateUpdate", 3)
	})

	t.Run("correctly uses cache for prices and doesn't recompute", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(6)
		numCallsToAggregateVEFunc := 0
		for _, call := range voteAggregator.Calls {
			if call.Method == "AggregateDaemonVEIntoFinalPricesAndConversionRate" {
				numCallsToAggregateVEFunc++
			}
		}
		// clear cache
		spotPriceUpdatesTest := veApplier.GetCachedSpotPrices()
		pnlPriceUpdatesTest := veApplier.GetCachedPnlPrices()
		fmt.Println("spotPriceUpdatesTest", spotPriceUpdatesTest)
		fmt.Println("pnlPriceUpdatesTest", pnlPriceUpdatesTest)
		rate1, ok := big.NewInt(0).SetString("123456789000000000000000000000", 10)
		require.True(t, ok)

		require.True(t, ok)

		price1Bz := big.NewInt(100).Bytes()

		prices1 := []vetypes.PricePair{
			{
				MarketId:  1,
				SpotPrice: price1Bz,
				PnlPrice:  price1Bz,
			},
		}

		vote1, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceConsAddress,
				prices1,
				"123456789000000000000000000000",
			),
		)
		require.NoError(t, err)

		_, extCommitInfoBz1, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote1},
		)
		require.NoError(t, err)

		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices:             prices1,
					SDaiConversionRate: "123456789000000000000000000000",
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]voteweighted.AggregatorPricePair{
			constants.BtcUsdPair: {
				SpotPrice: big.NewInt(100),
				PnlPrice:  big.NewInt(100),
			},
		}, rate1, nil).Once()

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
		).Times(4)

		pricesKeeper.On("UpdateSpotPrice", ctx, mock.Anything).Return(nil).Times(3)
		pricesKeeper.On("UpdatePnlPrice", ctx, mock.Anything).Return(nil).Times(3)

		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates := veApplier.GetCachedSpotPrices()
		pnlPriceUpdates := veApplier.GetCachedPnlPrices()

		cachedSpotPrices := make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices := make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedPnlPrices)

		cachedSDAIRate := veApplier.GetCachedSDaiConversionRate()
		require.Equal(t, rate1, cachedSDAIRate)

		err = veApplier.ApplyVE(
			ctx,
			[][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			true,
		)
		require.NoError(t, err)

		spotPriceUpdates = veApplier.GetCachedSpotPrices()
		pnlPriceUpdates = veApplier.GetCachedPnlPrices()

		cachedSpotPrices = make(map[string]uint64)
		for _, priceUpdate := range spotPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedSpotPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		cachedPnlPrices = make(map[string]uint64)
		for _, priceUpdate := range pnlPriceUpdates {
			marketId := priceUpdate.MarketId
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPnlPrices[pair.Pair] = priceUpdate.Price.Uint64()
		}

		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedSpotPrices)
		require.Equal(t, map[string]uint64{constants.BtcUsdPair: 100}, cachedPnlPrices)
		voteAggregator.AssertNumberOfCalls(t, "AggregateDaemonVEIntoFinalPricesAndConversionRate", numCallsToAggregateVEFunc+1)
	})
}

func TestCacheSeenExtendedVotes(t *testing.T) {
	tests := map[string]struct {
		setup                func(*testing.T, *cometabci.RequestCommit)
		expectedCacheEntries map[string]struct{}
	}{
		"Nil ExtendedCommitInfo": {
			setup: func(t *testing.T, req *cometabci.RequestCommit) {
				req.ExtendedCommitInfo = nil
			},
			expectedCacheEntries: map[string]struct{}{},
		},
		"Single validator vote": {
			setup: func(t *testing.T, req *cometabci.RequestCommit) {
				voteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
					vetesting.NewDefaultSignedVeInfo(
						constants.AliceConsAddress,
						constants.ValidVEPrices,
						"1000000000000000000000000000",
					),
				)
				require.NoError(t, err)
				req.ExtendedCommitInfo = &cometabci.ExtendedCommitInfo{
					Votes: []cometabci.ExtendedVoteInfo{voteInfo},
				}
			},
			expectedCacheEntries: map[string]struct{}{
				sdk.ConsAddress(constants.AliceConsAddress).String(): {},
			},
		},
		"Multiple validator votes": {
			setup: func(t *testing.T, req *cometabci.RequestCommit) {
				aliceVote, err := vetesting.CreateSignedExtendedVoteInfo(
					vetesting.NewDefaultSignedVeInfo(
						constants.AliceConsAddress,
						constants.ValidVEPrices,
						"1000000000000000000000000000",
					),
				)
				require.NoError(t, err)
				bobVote, err := vetesting.CreateSignedExtendedVoteInfo(
					vetesting.NewDefaultSignedVeInfo(
						constants.BobConsAddress,
						constants.ValidVEPrices,
						"1000000000000000000000000001",
					),
				)
				require.NoError(t, err)
				req.ExtendedCommitInfo = &cometabci.ExtendedCommitInfo{
					Votes: []cometabci.ExtendedVoteInfo{aliceVote, bobVote},
				}
			},
			expectedCacheEntries: map[string]struct{}{
				sdk.ConsAddress(constants.AliceConsAddress).String(): {},
				sdk.ConsAddress(constants.BobConsAddress).String():   {},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			memClob := memclob.NewMemClobPriceTimePriority(false)
			mockIndexerEventManager := &mocks.IndexerEventManager{}
			ks := keepertest.NewClobKeepersTestContext(t, memClob, &mocks.BankKeeper{}, mockIndexerEventManager, nil)

			req := &cometabci.RequestCommit{}
			tc.setup(t, req)

			// Execute
			err := ks.ClobKeeper.VEApplier.CacheSeenExtendedVotes(ks.Ctx, req)
			require.NoError(t, err)

			// Assert
			consAddresses := ks.ClobKeeper.VEApplier.GetVECache().GetSeenVotesInCache()
			require.Equal(t, tc.expectedCacheEntries, consAddresses)
		})
	}
}
