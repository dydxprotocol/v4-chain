package preblocker_test

import (
	"testing"

	"cosmossdk.io/log"

	preblocker "github.com/StreamFinance-Protocol/stream-chain/protocol/app/preblocker"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/suite"
)

type PreBlockTestSuite struct {
	suite.Suite

	validator       sdk.ConsAddress
	ctx             sdk.Context
	marketPairs     []pricestypes.MarketParam
	pricesKeeper    *pk.Keeper
	indexPriceCache *pricefeedtypes.MarketToExchangePrices
	priceApplier    veaggregator.PriceApplier
	handler         *preblocker.PreBlockHandler
	voteCodec       vecodec.VoteExtensionCodec
	extCodec        vecodec.ExtendedCommitCodec
	logger          log.Logger
}

func TestPreBlockTestSuite(t *testing.T) {
	suite.Run(t, new(PreBlockTestSuite))
}

func (s *PreBlockTestSuite) SetupTest() {
	s.validator = constants.AliceConsAddress

	ctx, pricesKeeper, _, indexPriceCahce, _, _ := keepertest.PricesKeepers(s.T())
	s.ctx = ctx
	s.pricesKeeper = pricesKeeper
	s.indexPriceCache = indexPriceCahce

	s.voteCodec = vecodec.NewDefaultVoteExtensionCodec()
	s.extCodec = vecodec.NewDefaultExtendedCommitCodec()

	s.logger = log.NewTestLogger(s.T())

	mCCVStore := &mocks.CCValidatorStore{}
	aggregationFn := voteweighted.Median(
		s.logger,
		mCCVStore,
		voteweighted.DefaultPowerThreshold,
	)

	aggregator := veaggregator.NewVeAggregator(
		s.logger,
		s.indexPriceCache,
		*s.pricesKeeper,
		aggregationFn,
	)

	s.priceApplier = veaggregator.NewPriceWriter(
		aggregator,
		*s.pricesKeeper,
		s.voteCodec,
		s.extCodec,
		s.logger,
	)

}

func (s *PreBlockTestSuite) TestPreBlocker() {

	s.Run("fail on nil request", func() {

		s.handler = preblocker.NewDaemonPreBlockHandler(
			s.logger,
			s.indexPriceCache,
			*s.pricesKeeper,
			s.priceApplier,
		)
		_, err := s.handler.PreBlocker(s.ctx, nil)
		s.Require().Error(err)

	})
}
