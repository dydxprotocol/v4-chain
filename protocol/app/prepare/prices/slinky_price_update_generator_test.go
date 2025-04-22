package prices_test

import (
	"fmt"
	"math/big"
	"testing"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/slinky/abci/strategies/aggregator"
	aggregatormock "github.com/dydxprotocol/slinky/abci/strategies/aggregator/mocks"
	codecmock "github.com/dydxprotocol/slinky/abci/strategies/codec/mocks"
	strategymock "github.com/dydxprotocol/slinky/abci/strategies/currencypair/mocks"
	"github.com/dydxprotocol/slinky/abci/testutils"
	vetypes "github.com/dydxprotocol/slinky/abci/ve/types"
	oracletypes "github.com/dydxprotocol/slinky/pkg/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/prepare/prices"
	"github.com/stretchr/testify/suite"
)

type SlinkyPriceUpdateGeneratorSuite struct {
	suite.Suite

	spug *prices.SlinkyPriceUpdateGenerator

	cps *strategymock.CurrencyPairStrategy

	veCodec *codecmock.VoteExtensionCodec

	extCommitCodec *codecmock.ExtendedCommitCodec

	va *aggregatormock.VoteAggregator
}

func TestSlinkyPriceUpdateGeneratorSuite(t *testing.T) {
	suite.Run(t, new(SlinkyPriceUpdateGeneratorSuite))
}

func (suite *SlinkyPriceUpdateGeneratorSuite) SetupTest() {
	// setup mocks
	suite.veCodec = codecmock.NewVoteExtensionCodec(suite.T())
	suite.extCommitCodec = codecmock.NewExtendedCommitCodec(suite.T())
	suite.va = aggregatormock.NewVoteAggregator(suite.T())
	suite.cps = strategymock.NewCurrencyPairStrategy(suite.T())
	suite.spug = prices.NewSlinkyPriceUpdateGenerator(
		suite.va,
		suite.extCommitCodec,
		suite.veCodec,
		suite.cps,
	)
}

// Test that if vote-extensions aren't enabled price-update-generator returns an empty update
func (suite *SlinkyPriceUpdateGeneratorSuite) TestWithVoteExtensionsNotEnabled() {
	// setup
	ctx := testutils.UpdateContextWithVEHeight(testutils.CreateBaseSDKContext(suite.T()), 5)
	// ctx.BlockHeight() < ctx.ConsensusParams.VoteExtensionsEnableHeight (VEs are not enabled rn)
	ctx = ctx.WithBlockHeight(4)

	// expect
	msg, err := suite.spug.GetValidMarketPriceUpdates(ctx, []byte{}) // 2nd argument shld be irrelevant

	// assert
	suite.NoError(err)
	suite.NotNil(msg)
	suite.Empty(msg.MarketPriceUpdates) // no updates
}

// Test that if aggregating oracle votes fails, we fail
func (suite *SlinkyPriceUpdateGeneratorSuite) TestVoteExtensionAggregationFails() {
	// setup
	ctx := testutils.UpdateContextWithVEHeight(testutils.CreateBaseSDKContext(suite.T()), 5)
	ctx = ctx.WithBlockHeight(6) // VEs enabled

	// create vote-extensions
	validator := []byte("validator")
	// we j mock what the actual wire-transmitted bz are for this vote-extension
	voteExtensionBz := []byte("vote-extension")
	extCommitBz := []byte("ext-commit") // '' for ext-commit
	extCommit := cmtabci.ExtendedCommitInfo{
		Votes: []cmtabci.ExtendedVoteInfo{
			{
				Validator: cmtabci.Validator{
					Address: validator,
				},
				VoteExtension: voteExtensionBz,
			},
		},
	}

	// mock codecs
	suite.extCommitCodec.On("Decode", extCommitBz).Return(extCommit, nil)

	ve := vetypes.OracleVoteExtension{
		Prices: map[uint64][]byte{
			1: []byte("price"),
		},
	}
	suite.veCodec.On("Decode", voteExtensionBz).Return(ve, nil)

	// expect an error from the vote-extension aggregator
	suite.va.On("AggregateOracleVotes", ctx, []aggregator.Vote{
		{
			ConsAddress:         sdk.ConsAddress(validator),
			OracleVoteExtension: ve,
		},
	}).Return(nil, fmt.Errorf("error in aggregation"))

	// execute
	msg, err := suite.spug.GetValidMarketPriceUpdates(ctx, extCommitBz)
	suite.Nil(msg)
	suite.Error(err, "error in aggregation")
}

// test that if price / ID conversion fails we fail
func (suite *SlinkyPriceUpdateGeneratorSuite) TestCurrencyPairConversionFails() {
	// setup
	ctx := testutils.UpdateContextWithVEHeight(testutils.CreateBaseSDKContext(suite.T()), 5)
	ctx = ctx.WithBlockHeight(6) // VEs enabled

	// create vote-extensions
	validator := []byte("validator")
	// we j mock what the actual wire-transmitted bz are for this vote-extension
	voteExtensionBz := []byte("vote-extension")
	extCommitBz := []byte("ext-commit") // '' for ext-commit
	extCommit := cmtabci.ExtendedCommitInfo{
		Votes: []cmtabci.ExtendedVoteInfo{
			{
				Validator: cmtabci.Validator{
					Address: validator,
				},
				VoteExtension: voteExtensionBz,
			},
		},
	}

	// mock codecs
	suite.extCommitCodec.On("Decode", extCommitBz).Return(extCommit, nil)

	ve := vetypes.OracleVoteExtension{
		Prices: map[uint64][]byte{
			1: []byte("price"),
		},
	}
	suite.veCodec.On("Decode", voteExtensionBz).Return(ve, nil)

	mogBtc := oracletypes.NewCurrencyPair("MOG", "BTC")
	// expect an error from the vote-extension aggregator
	suite.va.On("AggregateOracleVotes", ctx, []aggregator.Vote{
		{
			ConsAddress:         sdk.ConsAddress(validator),
			OracleVoteExtension: ve,
		},
	}).Return(map[oracletypes.CurrencyPair]*big.Int{
		mogBtc: big.NewInt(1),
	}, nil)

	// expect an error from the currency-pair strategy
	suite.cps.On("ID", ctx, mogBtc).Return(uint64(0), fmt.Errorf("error in currency-pair conversion"))

	// execute
	msg, err := suite.spug.GetValidMarketPriceUpdates(ctx, extCommitBz)
	suite.Nil(msg)
	suite.Error(err, "error in currency-pair conversion")
}

// test that the MsgUpdateMarketPricesTx is generated correctly
func (suite *SlinkyPriceUpdateGeneratorSuite) TestValidMarketPriceUpdate() {
	// setup
	ctx := testutils.UpdateContextWithVEHeight(testutils.CreateBaseSDKContext(suite.T()), 5)
	ctx = ctx.WithBlockHeight(6) // VEs enabled

	// create vote-extensions
	validator := []byte("validator")
	// we j mock what the actual wire-transmitted bz are for this vote-extension
	voteExtensionBz := []byte("vote-extension")
	extCommitBz := []byte("ext-commit") // '' for ext-commit
	extCommit := cmtabci.ExtendedCommitInfo{
		Votes: []cmtabci.ExtendedVoteInfo{
			{
				Validator: cmtabci.Validator{
					Address: validator,
				},
				VoteExtension: voteExtensionBz,
			},
		},
	}

	// mock codecs
	suite.extCommitCodec.On("Decode", extCommitBz).Return(extCommit, nil)

	ve := vetypes.OracleVoteExtension{
		Prices: map[uint64][]byte{
			1: []byte("price"),
		},
	}
	suite.veCodec.On("Decode", voteExtensionBz).Return(ve, nil)

	mogBtc := oracletypes.NewCurrencyPair("MOG", "BTC")
	pepeEth := oracletypes.NewCurrencyPair("PEPE", "ETH")
	// expect an error from the vote-extension aggregator
	suite.va.On("AggregateOracleVotes", ctx, []aggregator.Vote{
		{
			ConsAddress:         sdk.ConsAddress(validator),
			OracleVoteExtension: ve,
		},
	}).Return(map[oracletypes.CurrencyPair]*big.Int{
		mogBtc:  big.NewInt(1),
		pepeEth: big.NewInt(2),
	}, nil)

	// expect no error from currency-pair strategies
	suite.cps.On("ID", ctx, mogBtc).Return(uint64(0), nil)
	suite.cps.On("ID", ctx, pepeEth).Return(uint64(1), nil)

	// execute
	msg, err := suite.spug.GetValidMarketPriceUpdates(ctx, extCommitBz)
	suite.NoError(err)

	// check the message
	suite.Len(msg.MarketPriceUpdates, 2)
	expectedPrices := map[uint64]uint64{
		0: 1,
		1: 2,
	}
	for _, mpu := range msg.MarketPriceUpdates {
		price, ok := expectedPrices[uint64(mpu.MarketId)]
		suite.True(ok)

		suite.Equal(price, mpu.Price)
	}
}
