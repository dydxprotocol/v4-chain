package preblocker_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	preblocker "github.com/StreamFinance-Protocol/stream-chain/protocol/app/preblocker"
	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	veaggregator "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/aggregator"
	veapplier "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/applier"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	bigintcache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/bigintcache"
	pricecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/pricecache"
	vecache "github.com/StreamFinance-Protocol/stream-chain/protocol/caches/vecache"
	valutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	pricefeedtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/pricefeed"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	pricestest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/prices"
	vetesting "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type PreBlockTestSuite struct {
	suite.Suite

	validator         sdk.ConsAddress
	ctx               sdk.Context
	marketParamPrices []pricestypes.MarketParamPrice
	pricesKeeper      *pk.Keeper
	ratelimitKeeper   *ratelimitkeeper.Keeper
	daemonPriceCache  *pricefeedtypes.MarketToExchangePrices
	veApplier         *veapplier.VEApplier
	handler           *preblocker.PreBlockHandler
	valStore          *mocks.ValidatorStore
	voteCodec         vecodec.VoteExtensionCodec
	extCodec          vecodec.ExtendedCommitCodec
	logger            log.Logger
}

func TestPreBlockTestSuite(t *testing.T) {
	suite.Run(t, new(PreBlockTestSuite))
}

func (s *PreBlockTestSuite) SetupTest() {
	s.validator = constants.AliceConsAddress

	ctx, _, pricesKeeper, _, _, _, _, ratelimitKeeper, _, _ := keepertest.SubaccountsKeepers(s.T(), true)

	mockTimeProvider := &mocks.TimeProvider{}
	mockTimeProvider.On("Now").Return(constants.TimeT)
	s.ctx = ctx
	s.pricesKeeper = pricesKeeper
	s.ratelimitKeeper = ratelimitKeeper
	s.daemonPriceCache = pricesKeeper.DaemonPriceCache

	s.voteCodec = vecodec.NewDefaultVoteExtensionCodec()
	s.extCodec = vecodec.NewDefaultExtendedCommitCodec()

	s.logger = log.NewTestLogger(s.T())

	mValStore := &mocks.ValidatorStore{}
	s.valStore = mValStore

	pricesAggregatorFn := voteweighted.MedianPrices(
		s.logger,
		s.valStore,
		voteweighted.DefaultPowerThreshold,
	)

	conversionRateAggregatorFn := voteweighted.MedianConversionRate(
		s.logger,
		s.valStore,
		voteweighted.DefaultPowerThreshold,
	)

	aggregator := veaggregator.NewVeAggregator(
		s.logger,
		*s.pricesKeeper,
		*s.ratelimitKeeper,
		pricesAggregatorFn,
		conversionRateAggregatorFn,
	)

	spotPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
	pnlPriceUpdateCache := pricecache.PriceUpdatesCacheImpl{}
	sDaiConversionRateCache := bigintcache.BigIntCacheImpl{}
	vecache := vecache.NewVECache()
	s.veApplier = veapplier.NewVEApplier(
		s.logger,
		aggregator,
		*s.pricesKeeper,
		*s.ratelimitKeeper,
		s.voteCodec,
		s.extCodec,
		&spotPriceUpdateCache,
		&pnlPriceUpdateCache,
		&sDaiConversionRateCache,
		vecache,
	)

	s.marketParamPrices = s.setMarketPrices()

	s.createTestMarkets()
}

func (s *PreBlockTestSuite) TestPreBlocker() {
	s.Run("fail on nil request", func() {
		s.ctx = vetesting.GetVeEnabledCtx(s.ctx, 3)
		s.handler = preblocker.NewDaemonPreBlockHandler(
			s.logger,
			s.veApplier,
		)
		s.daemonPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)

		prePrices := s.getAllMarketPrices()

		_, err := s.handler.PreBlocker(s.ctx, nil)
		s.Require().Error(err)

		s.Require().True(s.ensurePricesEqualToCurrent(prePrices))
	})

	s.Run("skip when vote extensions are disabled", func() {
		s.ctx = vetesting.GetVeEnabledCtx(s.ctx, 1)

		s.handler = preblocker.NewDaemonPreBlockHandler(
			s.logger,
			s.veApplier,
		)

		s.daemonPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)

		prePrices := s.getAllMarketPrices()

		_, err := s.handler.PreBlocker(s.ctx, &cometabci.RequestFinalizeBlock{})
		s.Require().NoError(err)

		s.Require().True(s.ensurePricesEqualToCurrent(prePrices))
	})

	s.Run("ignore vote-extensions w/ prices for non-existent pairs", func() {
		s.ctx = vetesting.GetVeEnabledCtx(s.ctx, 4)

		s.handler = preblocker.NewDaemonPreBlockHandler(
			s.logger,
			s.veApplier,
		)

		s.daemonPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)

		priceBz, err := big.NewInt(1).GobEncode()
		s.Require().NoError(err)

		prices := []vetypes.PricePair{
			{
				MarketId:  10,
				SpotPrice: priceBz, // 10 is a nonexistent market
				PnlPrice:  priceBz, // 10 is a nonexistent market
			},
		}

		extCommitBz := s.getVoteExtensionsForValidatorsWithSamePrices(
			[]string{"alice", "bob"},
			prices,
			"",
		)

		s.mockValStoreAndTotalBondedTokensCall([]string{"alice", "bob"})

		prePrices := s.getAllMarketPrices()

		_, err = s.handler.PreBlocker(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})
		s.Require().NoError(err)

		s.Require().True(s.ensurePricesEqualToCurrent(prePrices))
	})

	s.Run("multiple markets to write prices for", func() {
		s.ctx = vetesting.GetVeEnabledCtx(s.ctx, 5)

		s.handler = preblocker.NewDaemonPreBlockHandler(
			s.logger,
			s.veApplier,
		)

		s.daemonPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)

		price1 := uint64(1)
		price2 := uint64(2)
		price3 := uint64(3)

		price1Bz, err := big.NewInt(int64(price1)).GobEncode()
		s.Require().NoError(err)
		price2Bz, err := big.NewInt(int64(price2)).GobEncode()
		s.Require().NoError(err)
		price3Bz, err := big.NewInt(int64(price3)).GobEncode()
		s.Require().NoError(err)

		prices := []vetypes.PricePair{
			{
				MarketId:  6,
				SpotPrice: price1Bz,
				PnlPrice:  price1Bz,
			},
			{
				MarketId:  7,
				SpotPrice: price2Bz,
				PnlPrice:  price2Bz,
			},
			{
				MarketId:  8,
				SpotPrice: price3Bz,
				PnlPrice:  price3Bz,
			},
		}

		extCommitBz := s.getVoteExtensionsForValidatorsWithSamePrices(
			[]string{"alice", "bob"},
			prices,
			"",
		)

		s.mockValStoreAndTotalBondedTokensCall([]string{"alice", "bob"})

		_, err = s.handler.PreBlocker(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})
		s.Require().NoError(err)

		s.arePriceUpdatesCorrect(
			map[uint32]ve.VEPricePair{
				6: {
					SpotPrice: price1,
					PnlPrice:  price1,
				},
				7: {
					SpotPrice: price2,
					PnlPrice:  price2,
				},
				8: {
					SpotPrice: price3,
					PnlPrice:  price3,
				},
			},
		)
	})

	s.Run("multiple markets with different spot and pnl prices", func() {
		s.ctx = vetesting.GetVeEnabledCtx(s.ctx, 5)

		s.handler = preblocker.NewDaemonPreBlockHandler(
			s.logger,
			s.veApplier,
		)

		s.daemonPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)

		spotPrice1 := uint64(1)
		pnlPrice1 := uint64(10)
		spotPrice2 := uint64(2)
		pnlPrice2 := uint64(20)
		spotPrice3 := uint64(3)
		pnlPrice3 := uint64(30)

		spotPrice1Bz, err := big.NewInt(int64(spotPrice1)).GobEncode()
		s.Require().NoError(err)
		pnlPrice1Bz, err := big.NewInt(int64(pnlPrice1)).GobEncode()
		s.Require().NoError(err)
		spotPrice2Bz, err := big.NewInt(int64(spotPrice2)).GobEncode()
		s.Require().NoError(err)
		pnlPrice2Bz, err := big.NewInt(int64(pnlPrice2)).GobEncode()
		s.Require().NoError(err)
		spotPrice3Bz, err := big.NewInt(int64(spotPrice3)).GobEncode()
		s.Require().NoError(err)
		pnlPrice3Bz, err := big.NewInt(int64(pnlPrice3)).GobEncode()
		s.Require().NoError(err)

		prices := []vetypes.PricePair{
			{
				MarketId:  6,
				SpotPrice: spotPrice1Bz,
				PnlPrice:  pnlPrice1Bz,
			},
			{
				MarketId:  7,
				SpotPrice: spotPrice2Bz,
				PnlPrice:  pnlPrice2Bz,
			},
			{
				MarketId:  8,
				SpotPrice: spotPrice3Bz,
				PnlPrice:  pnlPrice3Bz,
			},
		}

		extCommitBz := s.getVoteExtensionsForValidatorsWithSamePrices(
			[]string{"alice", "bob"},
			prices,
			"",
		)

		s.mockValStoreAndTotalBondedTokensCall([]string{"alice", "bob"})

		_, err = s.handler.PreBlocker(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{extCommitBz, {1, 2, 3, 4}, {1, 2, 3, 4}},
		})
		s.Require().NoError(err)

		s.arePriceUpdatesCorrect(
			map[uint32]ve.VEPricePair{
				6: {
					SpotPrice: spotPrice1,
					PnlPrice:  pnlPrice1,
				},
				7: {
					SpotPrice: spotPrice2,
					PnlPrice:  pnlPrice2,
				},
				8: {
					SpotPrice: spotPrice3,
					PnlPrice:  pnlPrice3,
				},
			},
		)
	})

	s.Run("throws error if can't get extCommitInfo", func() {
		s.ctx = vetesting.GetVeEnabledCtx(s.ctx, 6)

		s.handler = preblocker.NewDaemonPreBlockHandler(
			s.logger,
			s.veApplier,
		)

		s.daemonPriceCache.UpdatePrices(constants.MixedTimePriceUpdate)

		_, err := s.handler.PreBlocker(s.ctx, &cometabci.RequestFinalizeBlock{
			Txs: [][]byte{},
		})

		s.Require().EqualError(err, "error fetching extended-commit-info: proposal slice is too short, expected at least 1 elements but got 0")
	})
}

func (s *PreBlockTestSuite) getAllMarketPrices() []pricestypes.MarketPrice {
	return s.pricesKeeper.GetAllMarketPrices(s.ctx)
}

func (s *PreBlockTestSuite) arePriceUpdatesCorrect(
	correctPrices map[uint32]ve.VEPricePair,
) bool {
	prices := s.getAllMarketPrices()
	if len(prices) != len(correctPrices) {
		return false
	}

	for _, price := range prices {
		correctPrice, ok := correctPrices[price.Id]
		if !ok {
			return false
		}

		if price.SpotPrice != correctPrice.SpotPrice {
			return false
		}

		if price.PnlPrice != correctPrice.PnlPrice {
			return false
		}
	}
	return true
}

func (s *PreBlockTestSuite) ensurePricesEqualToCurrent(before []pricestypes.MarketPrice) bool {
	currPrices := s.getAllMarketPrices()
	if len(before) != len(currPrices) {
		return false
	}

	for i, price := range before {
		if price.Id != currPrices[i].Id {
			return false
		}
		if price.Exponent != currPrices[i].Exponent {
			return false
		}
		if price.SpotPrice != currPrices[i].SpotPrice {
			return false
		}

		if price.PnlPrice != currPrices[i].PnlPrice {
			return false
		}
	}
	return true
}

func (s *PreBlockTestSuite) createTestMarkets() {
	keepertest.CreateTestPriceMarkets(
		s.T(),
		s.ctx,
		s.pricesKeeper,
		s.marketParamPrices,
	)
}

func (s *PreBlockTestSuite) getVoteExtension(
	prices []vetypes.PricePair,
	sdaiConversionRate string,
	val sdk.ConsAddress,
) cometabci.ExtendedVoteInfo {
	ve, err := vetesting.CreateSignedExtendedVoteInfo(
		vetesting.SignedVEInfo{
			Val:                val,
			Power:              500,
			Prices:             prices,
			SDaiConversionRate: sdaiConversionRate,
			Height:             3,
			Round:              0,
			ChainId:            "localdydxprotocol",
		},
	)
	s.Require().NoError(err)
	return ve
}

func (s *PreBlockTestSuite) getExtendedCommitInfoBz(
	votes []cometabci.ExtendedVoteInfo,
) []byte {
	_, extCommitBz, err := vetesting.CreateExtendedCommitInfo(
		votes,
	)
	s.Require().NoError(err)
	return extCommitBz
}

func (s *PreBlockTestSuite) setMarketPrices() []pricestypes.MarketParamPrice {
	return []pricestypes.MarketParamPrice{
		*pricestest.GenerateMarketParamPrice(
			pricestest.WithId(6),
			pricestest.WithMinExchanges(2),
		),
		*pricestest.GenerateMarketParamPrice(
			pricestest.WithId(7),
			pricestest.WithMinExchanges(2),
		),
		*pricestest.GenerateMarketParamPrice(
			pricestest.WithId(8),
			pricestest.WithMinExchanges(2),
			pricestest.WithExponent(-8),
		),
		*pricestest.GenerateMarketParamPrice(
			pricestest.WithId(9),
			pricestest.WithMinExchanges(2),
			pricestest.WithExponent(-9),
		),
	}
}

func (s *PreBlockTestSuite) getVoteExtensionsForValidatorsWithSamePrices(
	validators []string,
	prices []vetypes.PricePair,
	sdaiConversionRate string,
) []byte {
	var votes []cometabci.ExtendedVoteInfo
	for _, valName := range validators {
		ve := s.getVoteExtension(prices, sdaiConversionRate, s.getValidatorConsAddr(valName))
		votes = append(votes, ve)
	}
	return s.getExtendedCommitInfoBz(votes)
}

func (s *PreBlockTestSuite) mockValStoreAndTotalBondedTokensCall(validators []string) {

	for _, valName := range validators {
		s.buildAndMockValidator(valName, math.NewInt(1))
	}
	s.valStore.On("TotalBondedTokens", s.ctx).Return(valutils.ConvertPowerToTokens(int64(len(validators))), nil)
}

func (s *PreBlockTestSuite) buildAndMockValidator(name string, bondedTokens math.Int) stakingtypes.ValidatorI {
	val := stakingtypes.Validator{
		Tokens: bondedTokens,
		Status: stakingtypes.Bonded,
	}
	s.valStore.On("GetValidator", s.ctx, s.getValidatorValAddress(name)).Return(val, nil)
	return val
}

func (s *PreBlockTestSuite) getValidatorConsAddr(name string) sdk.ConsAddress {
	if name == "alice" {
		return constants.AliceConsAddress
	} else {
		return constants.BobConsAddress
	}
}

func (s *PreBlockTestSuite) getValidatorValAddress(name string) sdk.ValAddress {
	if name == "alice" {
		return constants.AliceValAddress
	} else {
		return constants.BobValAddress
	}
}
