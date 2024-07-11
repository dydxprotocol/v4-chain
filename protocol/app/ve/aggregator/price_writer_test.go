package aggregator_test

import (
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
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

	pricesApplier := veaggregator.NewPriceWriter(
		voteAggregator,
		pricesKeeper,
		voteCodec,
		extCodec,
		log.NewNopLogger(),
	)

	t.Run("if extracting oracle votes fails, fail", func(t *testing.T) {
		prices, err := pricesApplier.ApplyPricesFromVoteExtensions(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{[]byte("garbage"), {1, 2, 3, 4}, {1, 2, 3, 4}},
		})

		require.Error(t, err)
		require.Nil(t, prices)
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
		voteAggregator.On("AggregateDaemonVE", ctx, []aggregator.Vote{
			{
				DaemonVoteExtension: vetypes.DaemonVoteExtension{
					Prices: prices,
				},
				ConsAddress: constants.AliceConsAddress,
			},
		}).Return(nil, fmt.Errorf("fail")).Once()

		returnedPrices, err := pricesApplier.ApplyPricesFromVoteExtensions(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})

		require.Error(t, err)
		require.Nil(t, returnedPrices)
	})

	t.Run("ignore negative prices", func(t *testing.T) {
		priceBz := big.NewInt(-100).Bytes()

		prices := map[uint32][]byte{
			1: priceBz,
		}

		_, extCommitInfoBz, err := vetesting.CreateSingleValidatorExtendedCommitInfo(
			constants.AliceConsAddress,
			prices,
		)
		require.NoError(t, err)

		voteAggregator.On("AggregateDaemonVE", ctx, []aggregator.Vote{
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

		_, err = pricesApplier.ApplyPricesFromVoteExtensions(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})

		require.NoError(t, err)

	})

	t.Run("update prices in state", func(t *testing.T) {
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

		voteAggregator.On("AggregateDaemonVE", ctx, []aggregator.Vote{
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

		pricesKeeper.On("UpdateMarketPrice", ctx, mock.Anything).Return(nil)

		prices, err := pricesApplier.ApplyPricesFromVoteExtensions(ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitInfoBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})

		require.NoError(t, err)
		require.Equal(t, map[string]*big.Int{
			constants.BtcUsdPair: big.NewInt(150),
		}, prices)

	})
}
