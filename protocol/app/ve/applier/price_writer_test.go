package price_writer_test

import (
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	pricewriter "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/applier"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPriceWriter(t *testing.T) {
	voteCodec := vecodec.NewDefaultVoteExtensionCodec()
	extCodec := vecodec.NewDefaultExtendedCommitCodec()

	voteAggregator := &mocks.VoteAggregator{}

	ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)

	pricesKeeper := &mocks.PriceApplierPricesKeeper{}

	pricesKeeper.On("PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything).Return(nil)

	pricesApplier := pricewriter.NewPriceApplier(
		log.NewNopLogger(),
		voteAggregator,
		pricesKeeper,
		voteCodec,
		extCodec,
	)

	t.Run("if extracting oracle votes fails, fail", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(1)
		err := pricesApplier.ApplyPricesFromVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{[]byte("garbage"), {1, 2, 3, 4}, {1, 2, 3, 4}},
		})

		priceUpdates := pricesApplier.GetCachedPrices()

		cachedPrices := make(map[string]*big.Int)
		for _, priceUpdate := range priceUpdates.MarketPriceUpdates {
			marketId := priceUpdate.GetMarketId()
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = big.NewInt(int64(priceUpdate.GetPrice()))
		}

		require.Error(t, err)
		require.Equal(t, cachedPrices, make(map[string]*big.Int))
	})

	t.Run("if vote aggregation fails, fail", func(t *testing.T) {
		prices := map[uint32][]byte{
			1: []byte("price1"),
		}

		_, extCommitInfoBz, err := vetesting.CreateSingleValidatorExtendedCommitInfo(
			constants.AliceConsAddress,
			prices,
		)
		require.NoError(t, err)

		ctx := sdk.Context{}

		// fail vote aggregation
		voteAggregator.On("AggregateDaemonVEIntoFinalPrices", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(nil, fmt.Errorf("fail")).Once()

		err = pricesApplier.ApplyPricesFromVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})

		priceUpdates := pricesApplier.GetCachedPrices()

		cachedPrices := make(map[string]*big.Int)
		for _, priceUpdate := range priceUpdates.MarketPriceUpdates {
			marketId := priceUpdate.GetMarketId()
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = big.NewInt(int64(priceUpdate.GetPrice()))
		}

		require.Error(t, err)
		require.Equal(t, cachedPrices, make(map[string]*big.Int))
	})

	t.Run("ignore negative prices", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(2)

		priceBz := big.NewInt(-100).Bytes()

		prices := map[uint32][]byte{
			1: priceBz,
		}

		_, extCommitInfoBz, err := vetesting.CreateSingleValidatorExtendedCommitInfo(
			constants.AliceConsAddress,
			prices,
		)
		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVEIntoFinalPrices", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(-100),
		}, nil)

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

		err = pricesApplier.ApplyPricesFromVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})

		require.NoError(t, err)
	})

	t.Run("update prices in state", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(3)

		price1Bz := big.NewInt(100).Bytes()
		price2Bz := big.NewInt(200).Bytes()

		prices1 := map[uint32][]byte{
			1: price1Bz,
		}

		prices2 := map[uint32][]byte{
			1: price2Bz,
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

		voteAggregator.On("AggregateDaemonVEIntoFinalPrices", ctx, []aggregator.Vote{
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
		}).Return(map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(150),
		}, nil)

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

		pricesKeeper.On("UpdateMarketPrice", ctx, mock.Anything).Return(nil)

		err = pricesApplier.ApplyPricesFromVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: []cometabci.VoteInfo{},
			},
		})

		priceUpdates := pricesApplier.GetCachedPrices()

		cachedPrices := make(map[string]*big.Int)
		for _, priceUpdate := range priceUpdates.MarketPriceUpdates {
			marketId := priceUpdate.GetMarketId()
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = big.NewInt(int64(priceUpdate.GetPrice()))
		}

		require.NoError(t, err)
		require.Equal(t, map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(150),
		}, cachedPrices)
	})

	t.Run("doesn't update prices for same round and height", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(4)

		price1Bz := big.NewInt(100).Bytes()
		price2Bz := big.NewInt(200).Bytes()

		prices1 := map[uint32][]byte{
			1: price1Bz,
		}

		prices2 := map[uint32][]byte{
			1: price2Bz,
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

		voteAggregator.On("AggregateDaemonVEIntoFinalPrices", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices1,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(100),
		}, nil).Once()

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

		pricesKeeper.On("UpdateMarketPrice", ctx, mock.Anything).Return(nil).Twice()

		// First call
		err = pricesApplier.ApplyPricesFromVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz1, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: []cometabci.VoteInfo{},
			},
		})
		require.NoError(t, err)

		priceUpdates := pricesApplier.GetCachedPrices()
		fmt.Println("priceUpdates 1", priceUpdates)
		cachedPrices := make(map[string]*big.Int)
		for _, priceUpdate := range priceUpdates.MarketPriceUpdates {
			marketId := priceUpdate.GetMarketId()
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = big.NewInt(int64(priceUpdate.GetPrice()))
		}

		require.Equal(t, map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(100),
		}, cachedPrices)

		voteAggregator.On("AggregateDaemonVEIntoFinalPrices", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices2,
				},
				ConsAddress: constants.BobConsAddress,
			},
		}).Return(map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(200),
		}, nil).Once()

		// Second call with the same round and height
		err = pricesApplier.ApplyPricesFromVE(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz2, {1, 2, 3, 4}, {1, 2, 3, 4}},
			DecidedLastCommit: cometabci.CommitInfo{
				Round: 1,
				Votes: []cometabci.VoteInfo{},
			},
		})
		require.NoError(t, err)

		priceUpdates = pricesApplier.GetCachedPrices()
		cachedPrices = make(map[string]*big.Int)
		for _, priceUpdate := range priceUpdates.MarketPriceUpdates {
			marketId := priceUpdate.GetMarketId()
			pair, exists := pricesKeeper.GetMarketParam(ctx, marketId)
			if !exists {
				continue
			}
			cachedPrices[pair.Pair] = big.NewInt(int64(priceUpdate.GetPrice()))
		}

		// Ensure the cached prices are still the same as the first call
		require.Equal(t, map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(100),
		}, cachedPrices)
	})
}
