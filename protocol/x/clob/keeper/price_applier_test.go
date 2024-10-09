package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	abcicomet "github.com/cometbft/cometbft/abci/types"
)

func TestSetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(t *testing.T) {
	// Setup
	memClob := &mocks.MemClob{}
	memClob.On("SetClobKeeper", mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(
		t,
		memClob,
		&mocks.BankKeeper{},
		&mocks.IndexerEventManager{},
		nil,
	)
	mockVEApplier := &mocks.VEApplierClobInterface{}
	mockVoteAggregator := &mocks.VoteAggregator{}
	ks.ClobKeeper.VEApplier = mockVEApplier
	mockVEApplier.On("VoteAggregator").Return(mockVoteAggregator)

	t.Run("nil ExtendedCommitInfo", func(t *testing.T) {
		err := ks.ClobKeeper.SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ks.Ctx, nil)
		require.NoError(t, err)
		mockVoteAggregator.AssertNotCalled(t, "AggregateDaemonVEIntoFinalPricesAndConversionRate")
	})

	t.Run("empty votes", func(t *testing.T) {
		extCommitInfo := &abcicomet.ExtendedCommitInfo{
			Votes: []abcicomet.ExtendedVoteInfo{},
		}
		err := ks.ClobKeeper.SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ks.Ctx, extCommitInfo)
		require.NoError(t, err)
		mockVoteAggregator.AssertNotCalled(t, "AggregateDaemonVEIntoFinalPricesAndConversionRate")
	})

	t.Run("error fetching votes", func(t *testing.T) {
		extCommitInfo := &abcicomet.ExtendedCommitInfo{
			Votes: []abcicomet.ExtendedVoteInfo{
				{
					VoteExtension: []byte("invalid"),
				},
			},
		}
		err := ks.ClobKeeper.SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ks.Ctx, extCommitInfo)
		require.Error(t, err)
		mockVoteAggregator.AssertNotCalled(t, "AggregateDaemonVEIntoFinalPricesAndConversionRate")
	})

	t.Run("error aggregating votes", func(t *testing.T) {
		extCommitInfo := CreateValidExtendedCommitInfo(t)
		mockVoteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", mock.Anything, mock.Anything).
			Return(nil, nil, fmt.Errorf("aggregation error")).Once()

		err := ks.ClobKeeper.SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ks.Ctx, extCommitInfo)
		require.NoError(t, err) // The function should not return an error in this case
		mockVEApplier.AssertNotCalled(t, "WritePricesToStoreAndMaybeCache")
		mockVEApplier.AssertNotCalled(t, "WriteSDaiConversionRateToStoreAndMaybeCache")
	})

	t.Run("successful price and rate update", func(t *testing.T) {
		extCommitInfo := CreateValidExtendedCommitInfo(t)
		prices := map[string]voteweighted.AggregatorPricePair{
			"BTC-USD": {SpotPrice: big.NewInt(50000), PnlPrice: big.NewInt(50000)},
		}
		conversionRate := big.NewInt(1000000)

		mockVoteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", mock.Anything, mock.Anything).
			Return(prices, conversionRate, nil).Once()
		mockVEApplier.On("WritePricesToStoreAndMaybeCache", mock.Anything, prices, []byte{}, false).
			Return(nil).Once()
		mockVEApplier.On("WriteSDaiConversionRateToStoreAndMaybeCache", mock.Anything, conversionRate, []byte{}, false).
			Return(nil).Once()

		err := ks.ClobKeeper.SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ks.Ctx, extCommitInfo)
		require.NoError(t, err)
		mockVEApplier.AssertExpectations(t)
	})

	t.Run("error writing prices", func(t *testing.T) {
		extCommitInfo := CreateValidExtendedCommitInfo(t)
		prices := map[string]voteweighted.AggregatorPricePair{
			"BTC-USD": {SpotPrice: big.NewInt(50000), PnlPrice: big.NewInt(50000)},
		}
		conversionRate := big.NewInt(1000000)

		mockVoteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", mock.Anything, mock.Anything).
			Return(prices, conversionRate, nil).Once()
		mockVEApplier.On("WritePricesToStoreAndMaybeCache", mock.Anything, prices, []byte{}, false).
			Return(fmt.Errorf("price write error")).Once()

		err := ks.ClobKeeper.SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ks.Ctx, extCommitInfo)
		require.Error(t, err)
		mockVEApplier.AssertNotCalled(t, "WriteSDaiConversionRateToStoreAndMaybeCache")
	})

	t.Run("error writing conversion rate", func(t *testing.T) {
		extCommitInfo := CreateValidExtendedCommitInfo(t)
		prices := map[string]voteweighted.AggregatorPricePair{
			"BTC-USD": {SpotPrice: big.NewInt(50000), PnlPrice: big.NewInt(50000)},
		}
		conversionRate := big.NewInt(1000000)

		mockVoteAggregator.On("AggregateDaemonVEIntoFinalPricesAndConversionRate", mock.Anything, mock.Anything).
			Return(prices, conversionRate, nil).Once()
		mockVEApplier.On("WritePricesToStoreAndMaybeCache", mock.Anything, prices, []byte{}, false).
			Return(nil).Once()
		mockVEApplier.On("WriteSDaiConversionRateToStoreAndMaybeCache", mock.Anything, conversionRate, []byte{}, false).
			Return(fmt.Errorf("conversion rate write error")).Once()

		err := ks.ClobKeeper.SetNextBlocksPricesAndSDAIRateFromExtendedCommitInfo(ks.Ctx, extCommitInfo)
		require.Error(t, err)
	})
}

func CreateValidExtendedCommitInfo(t *testing.T) *abcicomet.ExtendedCommitInfo {
	valVoteInfo, err := vetesting.CreateSignedExtendedVoteInfo(
		vetesting.NewDefaultSignedVeInfo(
			constants.AliceConsAddress,
			constants.ValidVEPrices,
			"1000",
		),
	)
	require.NoError(t, err)

	return &abcicomet.ExtendedCommitInfo{
		Votes: []abcicomet.ExtendedVoteInfo{valVoteInfo},
	}
}
