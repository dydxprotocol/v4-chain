package aggregator_test

import (
	"fmt"
	"testing"

	"cosmossdk.io/log"
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
	"github.com/stretchr/testify/require"
)

func TestPriceWriter(t *testing.T) {
	voteCodec := vecodec.NewDefaultVoteExtensionCodec()
	extCodec := vecodec.NewDefaultExtendedCommitCodec()

	voteAggregator := &mocks.VoteAggregator{}

	ctx, pricesKeeper, _, _, _, _ := keepertest.PricesKeepers(t)

	pricesApplier := veaggregator.NewPriceWriter(
		voteAggregator,
		*pricesKeeper,
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

		vote1, err := vetesting.CreateExtendedVoteInfo(
			constants.AliceConsAddress,
			prices,
			voteCodec,
		)
		require.NoError(t, err)

		_, extCommitInfoBz, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{vote1},
			extCodec,
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

}
