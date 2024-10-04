package ve_test

import (
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/vecache"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"

	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCleanAndValidateExtCommitInfo(t *testing.T) {
	// Create signed vote infos
	validVote, err := vetesting.CreateSignedExtendedVoteInfo(
		vetesting.NewDefaultSignedVeInfo(
			constants.AliceConsAddress,
			constants.ValidSingleVEPrice,
			"1000000000000000000000000000",
		),
	)
	require.NoError(t, err)

	invalidMarketVote, err := vetesting.CreateSignedExtendedVoteInfo(
		vetesting.NewDefaultSignedVeInfo(
			constants.AliceConsAddress,
			constants.ValidVEPrices,
			"1000000000000000000000000000",
		),
	)
	require.NoError(t, err)

	// Create a pruned version of the vote
	prunedVote := validVote
	prunedVote.BlockIdFlag = cmtproto.BlockIDFlagAbsent
	prunedVote.ExtensionSignature = nil
	prunedVote.VoteExtension = nil

	tests := map[string]struct {
		setupMocks    func(*mocks.PreBlockExecPricesKeeper, *mocks.VoteExtensionRateLimitKeeper)
		extCommitInfo cometabci.ExtendedCommitInfo
		expectedInfo  cometabci.ExtendedCommitInfo
		expectedError error
		blockHeight   int64
	}{
		"Valid ExtCommitInfo": {
			setupMocks: func(pricesKeeper *mocks.PreBlockExecPricesKeeper, ratelimitKeeper *mocks.VoteExtensionRateLimitKeeper) {
				pricesKeeper.On("GetAllMarketParams", mock.Anything).Return([]types.MarketParam{
					{Id: 0, Pair: constants.BtcUsdPair},
					{Id: 1, Pair: constants.EthUsdPair},
				})
				ratelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(big.NewInt(0), false)
				ratelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(big.NewInt(1), false)
			},
			extCommitInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{validVote},
			},
			expectedInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{validVote},
			},
			expectedError: nil,
			blockHeight:   100,
		},
		"Invalid market in VE": {
			setupMocks: func(pricesKeeper *mocks.PreBlockExecPricesKeeper, ratelimitKeeper *mocks.VoteExtensionRateLimitKeeper) {
				pricesKeeper.On("GetAllMarketParams", mock.Anything).Return([]types.MarketParam{
					{Id: 0, Pair: constants.BtcUsdPair},
					{Id: 1, Pair: constants.EthUsdPair},
				})
			},
			extCommitInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{invalidMarketVote},
			},
			expectedInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{prunedVote},
			},
			expectedError: nil,
			blockHeight:   100,
		},
		"Invalid sDai conversion rate height": {
			setupMocks: func(pricesKeeper *mocks.PreBlockExecPricesKeeper, ratelimitKeeper *mocks.VoteExtensionRateLimitKeeper) {
				pricesKeeper.On("GetAllMarketParams", mock.Anything).Return([]types.MarketParam{
					{Id: 0, Pair: constants.BtcUsdPair},
					{Id: 1, Pair: constants.EthUsdPair},
				})
				ratelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(big.NewInt(200), true)
			},
			extCommitInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{validVote},
			},
			expectedInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{prunedVote},
			},
			expectedError: nil,
			blockHeight:   100,
		},
		"Nil vote extension": {
			setupMocks: func(pricesKeeper *mocks.PreBlockExecPricesKeeper, ratelimitKeeper *mocks.VoteExtensionRateLimitKeeper) {
			},
			extCommitInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{prunedVote},
			},
			expectedInfo: cometabci.ExtendedCommitInfo{
				Round: 1,
				Votes: []cometabci.ExtendedVoteInfo{prunedVote},
			},
			expectedError: nil,
			blockHeight:   100,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = ctx.WithBlockHeight(tc.blockHeight)
			voteCodec := vecodec.NewDefaultVoteExtensionCodec()

			pricesKeeper := &mocks.PreBlockExecPricesKeeper{}
			ratelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}

			if tc.setupMocks != nil {
				tc.setupMocks(pricesKeeper, ratelimitKeeper)
			}
			veCache := vecache.NewVECache()

			result, err := ve.CleanAndValidateExtCommitInfo(
				ctx,
				tc.extCommitInfo,
				voteCodec,
				pricesKeeper,
				ratelimitKeeper,
				veCache,
			)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedInfo, result)

			pricesKeeper.AssertExpectations(t)
			ratelimitKeeper.AssertExpectations(t)
		})
	}
}
