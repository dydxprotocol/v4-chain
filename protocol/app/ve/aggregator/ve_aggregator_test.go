package aggregator_test

import (
	"testing"

	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	ethosutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ethos"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	voteCodec = vecodec.NewDefaultVoteExtensionCodec()
	extCodec  = vecodec.NewDefaultExtendedCommitCodec()
)

func SetupTest(t *testing.T, vals []string) (sdk.Context, veaggregator.VoteAggregator) {
	ctx, pk, _, daemonPriceCache, _, mTimeProvider := keepertest.PricesKeepers(t)
	mTimeProvider.On("Now").Return(constants.TimeT)

	keepertest.CreateTestMarkets(t, ctx, pk)

	mCCVStore := ethosutils.NewGetAllCCValidatorMockReturn(ctx, vals)

	aggregateFn := voteweighted.Median(
		ctx.Logger(),
		mCCVStore,
		voteweighted.DefaultPowerThreshold,
	)

	handler := veaggregator.NewVeAggregator(
		ctx.Logger(),
		daemonPriceCache,
		*pk,
		aggregateFn,
	)
	return ctx, handler
}
func TestVEAggregator(t *testing.T) {
	t.Run("no daemon data", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice"})
		_, commitBz, err := vetesting.CreateExtendedCommitInfo(nil)
		require.NoError(t, err)

		proposal := [][]byte{commitBz}

		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)

		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 0)
	})

	t.Run("Single daemon data", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice"})
		valVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 1)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
	})

	t.Run("Multiple price updates, single validator", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice"})

		valVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{valVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 3)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price6Big,
			PnlPrice:  constants.Price6Big,
		}, prices[constants.EthUsdPair])
		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price7Big,
			PnlPrice:  constants.Price7Big,
		}, prices[constants.SolUsdPair])
	})

	t.Run("single price update, from two validators", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob"})
		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		bobVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{aliceVoteInfo, bobVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 1)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
	})

	t.Run("multiple price updates, from two validators", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob"})

		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		bobVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{aliceVoteInfo, bobVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 3)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price6Big,
			PnlPrice:  constants.Price6Big,
		}, prices[constants.EthUsdPair])
		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price7Big,
			PnlPrice:  constants.Price7Big,
		}, prices[constants.SolUsdPair])
	})

	t.Run("single price update, from multiple validators", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob", "carl"})

		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		bobVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		carlVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.CarlEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{
				aliceVoteInfo,
				bobVoteInfo,
				carlVoteInfo,
			},
		)
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 1)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
	})

	t.Run("multiple price updates, from multiple validators", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob", "carl"})

		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		bobVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		carlVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.CarlEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{aliceVoteInfo, bobVoteInfo, carlVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 3)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price6Big,
			PnlPrice:  constants.Price6Big,
		}, prices[constants.EthUsdPair])
		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price7Big,
			PnlPrice:  constants.Price7Big,
		}, prices[constants.SolUsdPair])
	})

	t.Run("single price update from multiple validators but not enough voting power", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob", "carl"})
		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo(
			[]cometabci.ExtendedVoteInfo{
				aliceVoteInfo,
			},
		)
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 0)
	})

	t.Run("multiple price updates from multiple validators but not enough voting power", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob", "carl"})

		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{aliceVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 0)
	})

	t.Run("multiple prices from multiple validators but not enough voting power for some", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob", "carl"})

		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidVEPrice,
			),
		)
		require.NoError(t, err)

		bobVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)
		require.NoError(t, err)

		carlVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.CarlEthosConsAddress,
				constants.ValidSingleVEPrice,
			),
		)

		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{aliceVoteInfo, bobVoteInfo, carlVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 1)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
	})

	t.Run("continues when the validator's prices are malformed", func(t *testing.T) {
		ctx, handler := SetupTest(t, []string{"alice", "bob", "carl"})

		aliceVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.AliceEthosConsAddress,
				constants.ValidVEPricesWithOneInvalid,
			),
		)
		require.NoError(t, err)

		bobVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.BobEthosConsAddress,
				constants.ValidVEPricesWithOneInvalid,
			),
		)
		require.NoError(t, err)

		carlVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
			vetesting.NewDefaultSignedVeInfo(
				constants.CarlEthosConsAddress,
				constants.ValidVEPricesWithOneInvalid,
			),
		)
		require.NoError(t, err)

		// Create the extended commit info
		_, commitBz, err := vetesting.CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{aliceVoteInfo, bobVoteInfo, carlVoteInfo})
		require.NoError(t, err)

		proposal := [][]byte{commitBz}
		votes, err := veaggregator.GetDaemonVotesFromBlock(proposal, voteCodec, extCodec)
		require.NoError(t, err)
		prices, err := handler.AggregateDaemonVEIntoFinalPrices(ctx, votes)
		require.NoError(t, err)
		require.Len(t, prices, 2)

		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price5Big,
			PnlPrice:  constants.Price5Big,
		}, prices[constants.BtcUsdPair])
		require.Equal(t, voteweighted.AggregatorPricePair{
			SpotPrice: constants.Price6Big,
			PnlPrice:  constants.Price6Big,
		}, prices[constants.EthUsdPair])
	})
}
