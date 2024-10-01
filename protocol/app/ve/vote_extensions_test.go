package ve_test

import (
	"fmt"
	"math"
	"math/big"

	"cosmossdk.io/log"

	"testing"

	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	keepertest "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	vetestutils "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ve"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"

	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
)

type TestExtendedVoteTC struct {
	expectedResponse  *vetypes.DaemonVoteExtension
	pricesKeeper      func() *mocks.PreBlockExecPricesKeeper
	ratelimitKeeper   func() *mocks.VoteExtensionRateLimitKeeper
	sdaiEventManager  *sdaiservertypes.SDAIEventManager
	perpKeeper        func() *mocks.ExtendVotePerpetualsKeeper
	clobKeeper        func() *mocks.ExtendVoteClobKeeper
	extendVoteRequest func() *cometabci.RequestExtendVote
	expectedError     bool
}

func TestExtendVoteHandler(t *testing.T) {
	tests := map[string]TestExtendedVoteTC{
		"nil request returns error, no conversion rate available": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(true),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				return mClobKeeper
			},
			extendVoteRequest: func() *cometabci.RequestExtendVote {
				return nil
			},
			expectedResponse: nil,
			expectedError:    true,
		},
		"price daemon returns no prices and conversion rate daemon returns no conversion rate, error thrown": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					[]*pricestypes.MarketSpotPriceUpdate{},
				)

				mPricesKeeper.On("GetMarketParam", mock.Anything, mock.Anything).Return(
					&pricestypes.MarketParam{},
					false,
				)

				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(true),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).Return(
					nil, fmt.Errorf("error"),
				)
				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					nil, false,
				)
				return mClobKeeper
			},
			expectedError: true,
		},
		"price daemon returns no prices but sdai daemon returns conversion rate, no error": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					[]*pricestypes.MarketSpotPriceUpdate{},
				)

				mPricesKeeper.On("GetMarketParam", mock.Anything, mock.Anything).Return(
					&pricestypes.MarketParam{},
					false,
				)

				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).Return(
					nil, fmt.Errorf("error"),
				)
				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					nil, false,
				)
				return mClobKeeper
			},
			expectedError: false,
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices:             []vetypes.PricePair{},
				SDaiConversionRate: sdaiservertypes.TestSDAIEventRequest.ConversionRate,
			},
		},
		"oracle service returns single price, but no conversion rate available": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mpricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mpricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					constants.ValidSingleSpotMarketPriceUpdate,
				)
				mpricesKeeper.On("GetMarketParam", mock.Anything, mock.Anything).Return(
					constants.TestSingleMarketParam,
					true,
				)
				mpricesKeeper.On("GetSmoothedSpotPrice", mock.Anything).Return(
					constants.Price5,
					true,
				)

				return mpricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(true),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, uint32(0)).Return(
					perptypes.Perpetual{
						LastFundingRate: dtypes.NewInt(int64(constants.FundingRate1)),
					},
					nil,
				)

				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetSingleMarketClobMetadata", mock.Anything, mock.Anything).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price5)),
					},
				)
				return mClobKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
				},
				SDaiConversionRate: "",
			},
		},
		"oracle service returns single price and conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mpricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mpricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					constants.ValidSingleSpotMarketPriceUpdate,
				)
				mpricesKeeper.On("GetMarketParam", mock.Anything, mock.Anything).Return(
					constants.TestSingleMarketParam,
					true,
				)
				mpricesKeeper.On("GetSmoothedSpotPrice", mock.Anything).Return(
					constants.Price5,
					true,
				)

				return mpricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, uint32(0)).Return(
					perptypes.Perpetual{
						LastFundingRate: dtypes.NewInt(int64(constants.FundingRate1)),
					},
					nil,
				)

				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetSingleMarketClobMetadata", mock.Anything, mock.Anything).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price5)),
					},
				)
				return mClobKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
				},
				SDaiConversionRate: sdaiservertypes.TestSDAIEventRequest.ConversionRate,
			},
		},
		"oracle service returns multiple prices and no conversion rate available": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					constants.ValidMultiMarketSpotPriceUpdates,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId0).Return(
					constants.TestMarketParams[0],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId1).Return(
					constants.TestMarketParams[1],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId2).Return(
					constants.TestMarketParams[2],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId3).Return(
					constants.TestMarketParams[3],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId4).Return(
					constants.TestMarketParams[4],
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId0).Return(
					constants.Price5,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId1).Return(
					constants.Price6,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId2).Return(
					constants.Price7,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId3).Return(
					constants.Price4,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId4).Return(
					constants.Price3,
					true,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(true),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).Return(
					perptypes.Perpetual{
						LastFundingRate: dtypes.NewInt(int64(0)),
					},
					nil,
				)
				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId1,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId2,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId3,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId4,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId1,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId2,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId3,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId4,
						SubticksPerTick: 100_000,
					},
					true,
				)

				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price5)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId1,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price6)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId2,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price7)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId3,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price4)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId4,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price3)),
					},
				)
				return mClobKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
					{
						MarketId:  constants.MarketId2,
						SpotPrice: constants.Price7Bytes,
						PnlPrice:  constants.Price7Bytes,
					},
					{
						MarketId:  constants.MarketId3,
						SpotPrice: constants.Price4Bytes,
						PnlPrice:  constants.Price4Bytes,
					},
					{
						MarketId:  constants.MarketId4,
						SpotPrice: constants.Price3Bytes,
						PnlPrice:  constants.Price3Bytes,
					},
				},
				SDaiConversionRate: "",
			},
		},
		"oracle service returns multiple prices and conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					constants.ValidMultiMarketSpotPriceUpdates,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId0).Return(
					constants.TestMarketParams[0],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId1).Return(
					constants.TestMarketParams[1],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId2).Return(
					constants.TestMarketParams[2],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId3).Return(
					constants.TestMarketParams[3],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId4).Return(
					constants.TestMarketParams[4],
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId0).Return(
					constants.Price5,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId1).Return(
					constants.Price6,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId2).Return(
					constants.Price7,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId3).Return(
					constants.Price4,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId4).Return(
					constants.Price3,
					true,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).Return(
					perptypes.Perpetual{
						LastFundingRate: dtypes.NewInt(int64(0)),
					},
					nil,
				)
				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId1,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId2,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId3,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId4,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price5)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId1,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price6)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId2,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price7)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId3,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price4)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId4,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price3)),
					},
				)
				return mClobKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
					{
						MarketId:  constants.MarketId2,
						SpotPrice: constants.Price7Bytes,
						PnlPrice:  constants.Price7Bytes,
					},
					{
						MarketId:  constants.MarketId3,
						SpotPrice: constants.Price4Bytes,
						PnlPrice:  constants.Price4Bytes,
					},
					{
						MarketId:  constants.MarketId4,
						SpotPrice: constants.Price3Bytes,
						PnlPrice:  constants.Price3Bytes,
					},
				},
				SDaiConversionRate: sdaiservertypes.TestSDAIEventRequest.ConversionRate,
			},
		},
		"getting prices panics": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Panic("panic")
				return mPricesKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: nil,
			},
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				return mPerpKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(),
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				return mClobKeeper
			},
			expectedError: true,
		},
		"oracle service returns multiple prices, but no conversion rate because last set height is not old enough": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					constants.ValidMultiMarketSpotPriceUpdates,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId0).Return(
					constants.TestMarketParams[0],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId1).Return(
					constants.TestMarketParams[1],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId2).Return(
					constants.TestMarketParams[2],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId3).Return(
					constants.TestMarketParams[3],
					true,
				)
				mPricesKeeper.On("GetMarketParam", mock.Anything, constants.MarketId4).Return(
					constants.TestMarketParams[4],
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId0).Return(
					constants.Price5,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId1).Return(
					constants.Price6,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId2).Return(
					constants.Price7,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId3).Return(
					constants.Price4,
					true,
				)
				mPricesKeeper.On("GetSmoothedSpotPrice", constants.MarketId4).Return(
					constants.Price3,
					true,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(big.NewInt(1), true) // Note: assumes that test runs with low block height and offset is large
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).Return(
					perptypes.Perpetual{
						LastFundingRate: dtypes.NewInt(int64(0)),
					},
					nil,
				)
				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId1,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId2,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId3,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					clobtypes.ClobPair{
						Id:              constants.MarketId4,
						SubticksPerTick: 100_000,
					},
					true,
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId0,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price5)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId1,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price6)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId2,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price7)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId3,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price4)),
					},
				)
				mClobKeeper.On(
					"GetSingleMarketClobMetadata",
					mock.Anything,
					clobtypes.ClobPair{
						Id:              constants.MarketId4,
						SubticksPerTick: 100_000,
					},
				).Return(
					clobtypes.ClobMetadata{
						MidPrice: clobtypes.Subticks(getSubticksFromPrice(constants.Price3)),
					},
				)
				return mClobKeeper
			},
			expectedResponse: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
					{
						MarketId:  constants.MarketId2,
						SpotPrice: constants.Price7Bytes,
						PnlPrice:  constants.Price7Bytes,
					},
					{
						MarketId:  constants.MarketId3,
						SpotPrice: constants.Price4Bytes,
						PnlPrice:  constants.Price4Bytes,
					},
					{
						MarketId:  constants.MarketId4,
						SpotPrice: constants.Price3Bytes,
						PnlPrice:  constants.Price3Bytes,
					},
				},
				SDaiConversionRate: "",
			},
		},
		"price daemon returns no prices and conversion rate daemon returns conversion rate but block height is too new, error thrown": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(
					[]*pricestypes.MarketSpotPriceUpdate{},
				)

				mPricesKeeper.On("GetMarketParam", mock.Anything, mock.Anything).Return(
					&pricestypes.MarketParam{},
					false,
				)

				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
					Return(big.NewInt(1), true) // Note: assumes that test runs with low block height and offset is large
				return mRatelimitKeeper
			},
			sdaiEventManager: sdaiservertypes.SetupMockEventManager(),
			perpKeeper: func() *mocks.ExtendVotePerpetualsKeeper {
				mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
				mPerpKeeper.On("GetPerpetual", mock.Anything, mock.Anything).Return(
					nil, fmt.Errorf("error"),
				)
				return mPerpKeeper
			},
			clobKeeper: func() *mocks.ExtendVoteClobKeeper {
				mClobKeeper := &mocks.ExtendVoteClobKeeper{}
				mClobKeeper.On("GetClobPair", mock.Anything, mock.Anything).Return(
					nil, false,
				)
				return mClobKeeper
			},
			expectedError: true,
		},
	}
	var round int64 = 3
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = vetestutils.GetVeEnabledCtx(ctx, round)
			votecodec := vecodec.NewDefaultVoteExtensionCodec()

			mVEApplier := &mocks.VEApplierInterface{}

			h := ve.NewVoteExtensionHandler(
				log.NewTestLogger(t),
				votecodec,
				tc.pricesKeeper(),
				tc.perpKeeper(),
				tc.clobKeeper(),
				tc.ratelimitKeeper(),
				tc.sdaiEventManager,
				mVEApplier,
			)

			req := &cometabci.RequestExtendVote{}
			if tc.extendVoteRequest != nil {
				req = tc.extendVoteRequest()
			}

			mVEApplier.On("ApplyVE", mock.Anything, mock.Anything, mock.Anything).Return(nil)

			resp, err := h.ExtendVoteHandler()(ctx, req)
			if !tc.expectedError {
				require.NoError(t, err)
				if resp == nil && tc.expectedResponse == nil {
					return
				}
				require.NotNil(t, resp)
				ext, err := votecodec.Decode(resp.VoteExtension)
				require.NoError(t, err)
				require.Equal(t, len(tc.expectedResponse.Prices), len(ext.Prices))
				expectedPriceMap := make(map[uint32]vetypes.PricePair)
				for _, expectedPricePair := range tc.expectedResponse.Prices {
					expectedPriceMap[expectedPricePair.MarketId] = expectedPricePair
				}
				for _, actualPricePair := range ext.Prices {
					expectedPricePair, exists := expectedPriceMap[actualPricePair.MarketId]
					require.True(t, exists, "MarketId %d not found in expected prices", actualPricePair.MarketId)
					require.Equal(t, expectedPricePair.PnlPrice, actualPricePair.PnlPrice)
					require.Equal(t, expectedPricePair.SpotPrice, actualPricePair.SpotPrice)
				}
				require.Equal(t, tc.expectedResponse.SDaiConversionRate, ext.SDaiConversionRate)
			} else {
				require.Error(t, err)
			}
		})
		round++
	}
}

type TestVerifyExtendedVoteTC struct {
	getReq           func() *cometabci.RequestVerifyVoteExtension
	pricesKeeper     func() *mocks.PreBlockExecPricesKeeper
	ratelimitKeeper  func() *mocks.VoteExtensionRateLimitKeeper
	expectedResponse *cometabci.ResponseVerifyVoteExtension
	expectedError    bool
}

func TestVerifyVoteHandler(t *testing.T) {
	votecodec := vecodec.NewDefaultVoteExtensionCodec()
	tests := map[string]TestVerifyExtendedVoteTC{
		"nil request returns error": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return nil
			},
			expectedResponse: nil,
			expectedError:    true,
		},
		"empty vote extension": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"empty vote extension - 1 valid price": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetMaxPairs", mock.Anything).Return(1)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"malformed bytes returns error": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: []byte("malformed"),
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"vote extension with no prices and no conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					[]pricestypes.MarketParam{}, // two prices
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := []vetypes.PricePair{}

				extBz, err := vetestutils.CreateVoteExtensionBytes(
					prices,
					"",
				)
				require.NoError(t, err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},

			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - single valid price with no conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidSingleVEPrice,
					"",
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - no price with conversion rate 1": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					[]vetypes.PricePair{},
					sdaiservertypes.TestSDAIEventRequest.ConversionRate,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - no price with conversion rate 2": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					nil,
					sdaiservertypes.TestSDAIEventRequest.ConversionRate,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - single valid price with conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidSingleVEPrice,
					sdaiservertypes.TestSDAIEventRequest.ConversionRate,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - multiple valid prices and no conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					"",
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - multiple valid prices and initial conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					sdaiservertypes.TestSDAIEventRequest.ConversionRate,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - multiple valid prices and valid conversion rate": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(big.NewInt(1), true)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(big.NewInt(1000000000000000), true)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					sdaiservertypes.TestSDAIEventRequest.ConversionRate,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"invalid vote extension - multiple valid prices - should fail": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams[1:3], // two prices
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					"",
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"invalid vote extension - vote extension with malformed prices": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					[]pricestypes.MarketParam{constants.TestMarketParams[0]},
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				prices := []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: make([]byte, 34),
						PnlPrice:  make([]byte, 34),
					},
				}

				extBz, err := vetestutils.CreateVoteExtensionBytes(
					prices,
					"",
				)
				require.NoError(t, err)

				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"invalid vote extension - multiple valid prices, no conversion rate, but ve height is wrong": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					"",
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        5,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_ACCEPT,
			},
			expectedError: false,
		},
		"valid vote extension - multiple valid prices but conversion rate height is too new": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				// Note: the below assumes a low block height for the test and larger delay
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(big.NewInt(5500), true)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					sdaiservertypes.TestSDAIEventRequest.ConversionRate,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"valid vote extension - multiple valid prices but conversion rate is not increasing": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(
					ratelimitkeeper.ConvertStringToBigIntWithPanicOnErr(sdaiservertypes.TestSDAIEventRequest.ConversionRate),
					true,
				)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					sdaiservertypes.TestSDAIEventRequest.ConversionRate,
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"valid vote extension - multiple valid prices but conversion rate is too large": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					"10000000000000000000000000000000000000000000000000000000000",
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"valid vote extension - multiple valid prices but conversion rate is zero": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					"0",
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
		"valid vote extension - multiple valid prices but conversion rate is negative": {
			pricesKeeper: func() *mocks.PreBlockExecPricesKeeper {
				mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
				mPricesKeeper.On("GetAllMarketParams", mock.Anything).Return(
					constants.TestMarketParams,
				)
				return mPricesKeeper
			},
			ratelimitKeeper: func() *mocks.VoteExtensionRateLimitKeeper {
				mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
				mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).Return(new(big.Int), false)
				mRatelimitKeeper.On("GetSDAIPrice", mock.Anything).Return(new(big.Int), false)
				return mRatelimitKeeper
			},
			getReq: func() *cometabci.RequestVerifyVoteExtension {
				extBz, err := vetestutils.CreateVoteExtensionBytes(
					constants.ValidVEPrices,
					"-1",
				)
				require.NoError(t, err)
				return &cometabci.RequestVerifyVoteExtension{
					VoteExtension: extBz,
					Height:        6000,
				}
			},
			expectedResponse: &cometabci.ResponseVerifyVoteExtension{
				Status: cometabci.ResponseVerifyVoteExtension_REJECT,
			},
			expectedError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = vetestutils.GetVeEnabledCtx(ctx, 6000)
			mVEApplier := &mocks.VEApplierInterface{}
			mClobKeeper := &mocks.ExtendVoteClobKeeper{}
			mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}
			mPricesKeeper := tc.pricesKeeper()
			mRatelimitKeeper := tc.ratelimitKeeper()
			sdaiEventManager := sdaiservertypes.SetupMockEventManager()

			handler := ve.NewVoteExtensionHandler(
				log.NewTestLogger(t),
				votecodec,
				mPricesKeeper,
				mPerpKeeper,
				mClobKeeper,
				mRatelimitKeeper,
				sdaiEventManager,
				mVEApplier,
			).VerifyVoteExtensionHandler()

			resp, err := handler(ctx, tc.getReq())
			require.Equal(t, tc.expectedResponse, resp)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetVEBytes(t *testing.T) {
	tests := map[string]struct {
		markets        []uint32
		daemonPrices   []*pricestypes.MarketSpotPriceUpdate
		smoothedPrices map[uint32]uint64
		midPrices      map[uint32]uint64
		fundingRates   map[uint32]int64
		expected       *vetypes.DaemonVoteExtension
		expectedError  bool
	}{
		"valid single price, no funding-smooth-or-mid": {
			markets: []uint32{constants.MarketId0},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: math.MaxInt64,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
				},
			},
			expectedError: false,
		},
		"valid single price with funding, no smooth or mid": {
			markets: []uint32{constants.MarketId0},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: 2000, // 0.2% in ppm
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
				},
			},
			expectedError: false,
		},
		"valid multiple prices, no funding, smooth or mid": {
			markets: []uint32{constants.MarketId0, constants.MarketId1},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price6),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
				constants.MarketId1: math.MaxUint64,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: math.MaxInt64,
				constants.MarketId1: math.MaxInt64,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
				constants.MarketId1: math.MaxUint64,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
				},
			},
			expectedError: false,
		},
		"valid multiple prices with funding, no smooth or mid": {
			markets: []uint32{constants.MarketId0, constants.MarketId1},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price6),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
				constants.MarketId1: math.MaxUint64,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: 2000,
				constants.MarketId1: 1000,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
				constants.MarketId1: math.MaxUint64,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
				},
			},
			expectedError: false,
		},
		"valid multiple prices with funding and smooth, no mid": {
			markets: []uint32{constants.MarketId0, constants.MarketId1},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price6),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5,
				constants.MarketId1: constants.Price6,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: 2000,
				constants.MarketId1: 1000,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
				constants.MarketId1: math.MaxUint64,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
				},
			},
			expectedError: false,
		},
		"valid multiple prices with smooth and mid, no funding": {
			markets: []uint32{constants.MarketId0, constants.MarketId1},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price6),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5,
				constants.MarketId1: constants.Price6,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: math.MaxInt64,
				constants.MarketId1: math.MaxInt64,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5,
				constants.MarketId1: constants.Price6,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
				},
			},
			expectedError: false,
		},
		"valid multiple prices with smooth, no funding or mid": {
			markets: []uint32{constants.MarketId0, constants.MarketId1},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price6),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5,
				constants.MarketId1: constants.Price6,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: math.MaxInt64,
				constants.MarketId1: math.MaxInt64,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
				constants.MarketId1: math.MaxUint64,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
				},
			},
			expectedError: false,
		},
		"valid multiple prices with mid, no funding or smooth": {
			markets: []uint32{constants.MarketId0, constants.MarketId1},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price6),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: math.MaxUint64,
				constants.MarketId1: math.MaxUint64,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: math.MaxInt64,
				constants.MarketId1: math.MaxInt64,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5,
				constants.MarketId1: constants.Price6,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: constants.Price5Bytes,
						PnlPrice:  constants.Price5Bytes,
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: constants.Price6Bytes,
						PnlPrice:  constants.Price6Bytes,
					},
				},
			},
			expectedError: false,
		},
		"single price with smooth, funding, and mid": {
			markets: []uint32{constants.MarketId0},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5 - 100,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: 2000, // 0.2% in ppm
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5In100_000SubticksPerTick - 200_000,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: getGobEncodedPriceBytes(500005),
						PnlPrice:  getGobEncodedPriceBytes(500003),
					},
				},
			},
			expectedError: false,
		},
		"multiple prices with smooth, funding, and mid": {
			markets: []uint32{constants.MarketId0, constants.MarketId1},
			daemonPrices: []*pricestypes.MarketSpotPriceUpdate{
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId0, constants.Price5),
				pricestypes.NewMarketSpotPriceUpdate(constants.MarketId1, constants.Price6),
			},
			smoothedPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5 - 100,
				constants.MarketId1: constants.Price6 + 100,
			},
			fundingRates: map[uint32]int64{
				constants.MarketId0: 2000, // 0.2% in ppm
				constants.MarketId1: 500,
			},
			midPrices: map[uint32]uint64{
				constants.MarketId0: constants.Price5In100_000SubticksPerTick - 200_000,
				constants.MarketId1: constants.Price6In100_000SubticksPerTick + 200_000,
			},
			expected: &vetypes.DaemonVoteExtension{
				Prices: []vetypes.PricePair{
					{
						MarketId:  constants.MarketId0,
						SpotPrice: getGobEncodedPriceBytes(500005),
						PnlPrice:  getGobEncodedPriceBytes(500003),
					},
					{
						MarketId:  constants.MarketId1,
						SpotPrice: getGobEncodedPriceBytes(60006),
						PnlPrice:  getGobEncodedPriceBytes(60036),
					},
				},
			},
			expectedError: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			votecodec := vecodec.NewDefaultVoteExtensionCodec()
			mVEApplier := &mocks.VEApplierInterface{}
			mPricesKeeper := &mocks.PreBlockExecPricesKeeper{}
			mClobKeeper := &mocks.ExtendVoteClobKeeper{}
			mPerpKeeper := &mocks.ExtendVotePerpetualsKeeper{}

			mPricesKeeper.On("GetValidMarketSpotPriceUpdates", mock.Anything).Return(tc.daemonPrices)

			for _, market := range tc.markets {
				mPricesKeeper.On("GetMarketParam", mock.Anything, market).Return(
					pricestypes.MarketParam{
						Id:                 market,
						Pair:               constants.IdsToPairs[market],
						MinExchanges:       1,
						MinPriceChangePpm:  50,
						ExchangeConfigJson: constants.TestMarketExchangeConfigs[market],
					},
					true,
				)

				if tc.smoothedPrices[market] == math.MaxUint64 {
					mPricesKeeper.On("GetSmoothedSpotPrice", market).Return(
						uint64(0),
						false,
					)
				} else {
					mPricesKeeper.On("GetSmoothedSpotPrice", market).Return(
						tc.smoothedPrices[market],
						true,
					)
				}

				if tc.fundingRates[market] == math.MaxInt64 {
					mPerpKeeper.On("GetPerpetual", mock.Anything, market).Return(
						perptypes.Perpetual{
							LastFundingRate: dtypes.NewInt(0),
						},
						fmt.Errorf("error"),
					)
				} else {
					mPerpKeeper.On("GetPerpetual", mock.Anything, market).Return(
						perptypes.Perpetual{
							LastFundingRate: dtypes.NewInt(tc.fundingRates[market]),
						},
						nil,
					)
				}

				mClobKeeper.On("GetClobPair", mock.Anything, clobtypes.ClobPairId(market)).Return(
					clobtypes.ClobPair{
						Id:              market,
						SubticksPerTick: 100_000,
					},
					true,
				)

				if tc.midPrices[market] == math.MaxUint64 {
					mClobKeeper.On(
						"GetSingleMarketClobMetadata",
						mock.Anything,
						clobtypes.ClobPair{
							Id:              market,
							SubticksPerTick: 100_000,
						},
					).Return(
						clobtypes.ClobMetadata{
							MidPrice: 0,
						},
					)
				} else {
					mClobKeeper.On(
						"GetSingleMarketClobMetadata",
						mock.Anything,
						clobtypes.ClobPair{
							Id:              market,
							SubticksPerTick: 100_000,
						},
					).Return(
						clobtypes.ClobMetadata{
							MidPrice: clobtypes.Subticks(tc.midPrices[market]),
						},
					)
				}
			}

			sDaIEventManager := sdaiservertypes.SetupMockEventManager()
			mRatelimitKeeper := &mocks.VoteExtensionRateLimitKeeper{}
			mRatelimitKeeper.On("GetSDAILastBlockUpdated", mock.Anything).
				Return(new(big.Int), false)

			h := ve.NewVoteExtensionHandler(
				log.NewTestLogger(t),
				votecodec,
				mPricesKeeper,
				mPerpKeeper,
				mClobKeeper,
				mRatelimitKeeper,
				sDaIEventManager,
				mVEApplier,
			)

			ctx, _, _, _, _, _ := keepertest.PricesKeepers(t)
			ctx = vetestutils.GetVeEnabledCtx(ctx, 3)

			var expectedVEBytes []byte
			if tc.expected != nil {
				var err error
				expectedVEBytes, err = votecodec.Encode(*tc.expected)
				require.NoError(t, err)
			}

			veBytes, err := h.GetVEBytes(ctx)

			if tc.expectedError {
				require.Error(t, err)
				require.Nil(t, veBytes)
			} else {
				require.NoError(t, err)
				// Decode both expected and actual vote extensions
				expectedVE, err := votecodec.Decode(expectedVEBytes)
				require.NoError(t, err)
				actualVE, err := votecodec.Decode(veBytes)
				require.NoError(t, err)

				require.Equal(t, len(expectedVE.Prices), len(actualVE.Prices))

				expectedPriceMap := make(map[uint32]vetypes.PricePair)
				for _, expectedPricePair := range expectedVE.Prices {
					expectedPriceMap[expectedPricePair.MarketId] = expectedPricePair
				}

				for _, actualPricePair := range actualVE.Prices {
					expectedPricePair, exists := expectedPriceMap[actualPricePair.MarketId]
					require.True(t, exists, "MarketId %d not found in expected prices", actualPricePair.MarketId)
					require.Equal(t, expectedPricePair.PnlPrice, actualPricePair.PnlPrice)
					require.Equal(t, expectedPricePair.SpotPrice, actualPricePair.SpotPrice)
				}
			}
		})
	}
}

func getGobEncodedPriceBytes(
	price int64,
) []byte {
	bytes, err := big.NewInt(price).GobEncode()
	if err != nil {
		return []byte{}
	}
	return bytes
}

func getSubticksFromPrice(price uint64) clobtypes.Subticks {
	return clobtypes.Subticks(price * 100_000)
}
