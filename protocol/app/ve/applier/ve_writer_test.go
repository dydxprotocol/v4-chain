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
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPriceWriter(t *testing.T) {
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
			),
		)
		require.NoError(t, err)

		vote2, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobConsAddress,
				prices2,
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
			),
		)
		require.NoError(t, err)

		vote2, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobConsAddress,
				prices2,
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
		fmt.Println("priceUpdates 1", priceUpdates)
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
}
