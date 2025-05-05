package process_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/slinky/abci/testutils"
	"github.com/stretchr/testify/suite"

	"github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

type SlinkyMarketPriceDecoderSuite struct {
	suite.Suite

	// mock price-update generator
	gen *mocks.PriceUpdateGenerator
	// mock UpdateMarketPriceTxDecoder
	decoder *mocks.UpdateMarketPriceTxDecoder
	// mock context
	ctx sdk.Context
}

func TestSlinkyMarketPriceDecoderSuite(t *testing.T) {
	suite.Run(t, new(SlinkyMarketPriceDecoderSuite))
}

func (suite *SlinkyMarketPriceDecoderSuite) SetupTest() {
	// setup context
	suite.ctx = testutils.CreateBaseSDKContext(suite.T())

	// setup mocks
	suite.gen = mocks.NewPriceUpdateGenerator(suite.T())
	suite.decoder = mocks.NewUpdateMarketPriceTxDecoder(suite.T())
}

// test that if vote-extensions are not enabled, the MsgUpdateMarketPrices proposed should be empty
func (suite *SlinkyMarketPriceDecoderSuite) TestVoteExtensionsNotEnabled() {
	suite.Run("test that a non-empty proposed market-price update fails", func() {
		// disable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(3)

		proposal := [][]byte{[]byte("test")}

		// mock decoder response that returns non-empty prices
		suite.decoder.On("DecodeUpdateMarketPricesTx",
			suite.ctx, proposal).Return(
			process.NewUpdateMarketPricesTx(suite.ctx, nil, &pricestypes.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
					{
						MarketId: 1, // propose non-empty prices
						Price:    100,
					},
				},
			}), nil).Once()

		// expect an error
		expectError := process.IncorrectNumberUpdatesError(0, 1)

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(tx)
		suite.EqualError(expectError, err.Error())
	})

	suite.Run("test that the proposed prices must be empty", func() {
		// disable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(3)

		proposal := [][]byte{[]byte("test")}

		// mock decoder response that returns non-empty prices
		suite.decoder.On("DecodeUpdateMarketPricesTx",
			suite.ctx, proposal).Return(
			process.NewUpdateMarketPricesTx(suite.ctx, nil, &pricestypes.MsgUpdateMarketPrices{}), nil).Once()

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(err)
		suite.Len(tx.GetMsg().(*pricestypes.MsgUpdateMarketPrices).MarketPriceUpdates, 0)
	})
}

// test that if vote-extensions are enabled
//   - missing extended commit -> failure
//   - price-update generator fails -> failure
func (suite *SlinkyMarketPriceDecoderSuite) TestVoteExtensionsEnabled() {
	suite.Run("test that missing extended commit -> failure", func() {
		// enable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(5)

		proposal := [][]byte{}

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(tx)
		suite.Error(err)
	})

	suite.Run("test that price-update generator fails -> failure", func() {
		// enable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(5)

		proposal := [][]byte{[]byte("test")}
		err := fmt.Errorf("error")

		suite.gen.On("GetValidMarketPriceUpdates",
			suite.ctx, proposal[constants.OracleInfoIndex]).Return(nil, err)

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(tx)
		suite.Error(err)
	})
}

func (suite *SlinkyMarketPriceDecoderSuite) TestMarketPriceUpdateValidation_WithVoteExtensionsEnabled() {
	suite.Run("if DecodeUpdatemarketPricesTx fails on underlying decoder - fail", func() {
		// enable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(5)

		proposal := [][]byte{[]byte("test")}

		suite.gen.On("GetValidMarketPriceUpdates",
			suite.ctx, proposal[constants.OracleInfoIndex]).Return(&pricestypes.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				{
					MarketId: 1,
					Price:    100,
				},
			},
		}, nil)

		suite.decoder.On("DecodeUpdateMarketPricesTx",
			suite.ctx, proposal).Return(nil, fmt.Errorf("error"))

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(tx)
		suite.Error(err)
	})

	suite.Run("if DecodeUpdateMarketPricesTx returns conflicting updates (missing market-id) - fail", func() {
		// enable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(5)

		proposal := [][]byte{[]byte("test")}

		expectedMsg := &pricestypes.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				{
					MarketId: 1,
					Price:    100,
				},
			},
		}

		suite.gen.On("GetValidMarketPriceUpdates",
			suite.ctx, proposal[constants.OracleInfoIndex]).Return(expectedMsg, nil)

		suite.decoder.On("DecodeUpdateMarketPricesTx",
			suite.ctx, proposal).Return(
			process.NewUpdateMarketPricesTx(suite.ctx, nil, &pricestypes.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
					{
						MarketId: 2, // propose non-empty prices
						Price:    100,
					},
				},
			}), nil)

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(tx)
		suite.Error(err, process.MissingPriceUpdateForMarket(expectedMsg.MarketPriceUpdates[0].MarketId).Error())
	})

	suite.Run("if DecodeUpdateMarketPricesTx returns conflicting updates (incorrect price for market-id) - fail", func() {
		// enable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(5)

		proposal := [][]byte{[]byte("test")}

		expectedMsg := &pricestypes.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				{
					MarketId: 1,
					Price:    100,
				},
			},
		}

		suite.gen.On("GetValidMarketPriceUpdates",
			suite.ctx, proposal[constants.OracleInfoIndex]).Return(expectedMsg, nil)

		suite.decoder.On("DecodeUpdateMarketPricesTx",
			suite.ctx, proposal).Return(
			process.NewUpdateMarketPricesTx(suite.ctx, nil, &pricestypes.MsgUpdateMarketPrices{
				MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
					{
						MarketId: 1, // propose non-empty prices
						Price:    101,
					},
				},
			}), nil)

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(tx)
		suite.Error(err, process.IncorrectPriceUpdateForMarket(1, 100, 101))
	})

	suite.Run("if DecodeUpdateMarketPricesTx returns msgs that don't pass ValidateBasic - fail", func() {
		// enable VEs
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(5)

		proposal := [][]byte{[]byte("test")}

		expectedMsg := &pricestypes.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				{
					MarketId: 2, // incorrect order
					Price:    100,
				},
				{
					MarketId: 1,
					Price:    100,
				},
			},
		}

		suite.gen.On("GetValidMarketPriceUpdates",
			suite.ctx, proposal[constants.OracleInfoIndex]).Return(expectedMsg, nil)

		suite.decoder.On("DecodeUpdateMarketPricesTx",
			suite.ctx, proposal).Return(process.NewUpdateMarketPricesTx(suite.ctx, nil, expectedMsg), nil)

		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
		tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
		suite.Nil(tx)
		suite.Error(err)
	})
}

// test happy path
func (suite *SlinkyMarketPriceDecoderSuite) TestHappyPath_VoteExtensionsEnabled() {
	// enable VEs
	suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
	suite.ctx = suite.ctx.WithBlockHeight(5)

	proposal := [][]byte{[]byte("test")}

	expectedMsg := &pricestypes.MsgUpdateMarketPrices{
		MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
			{
				MarketId: 1,
				Price:    100,
			},
		},
	}

	suite.gen.On("GetValidMarketPriceUpdates",
		suite.ctx, proposal[constants.OracleInfoIndex]).Return(expectedMsg, nil)

	suite.decoder.On("DecodeUpdateMarketPricesTx",
		suite.ctx, proposal).Return(
		process.NewUpdateMarketPricesTx(suite.ctx, nil, &pricestypes.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*pricestypes.MsgUpdateMarketPrices_MarketPrice{
				{
					MarketId: 1, // propose non-empty prices
					Price:    100,
				},
			},
		}), nil)

	decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)
	tx, err := decoder.DecodeUpdateMarketPricesTx(suite.ctx, proposal)
	suite.NoError(err)
	suite.NotNil(tx)

	msg := tx.GetMsg().(*pricestypes.MsgUpdateMarketPrices)
	suite.Len(msg.MarketPriceUpdates, 1)
	suite.Equal(expectedMsg.MarketPriceUpdates[0], msg.MarketPriceUpdates[0])
}

func (suite *SlinkyMarketPriceDecoderSuite) TestGetTxOffset() {
	suite.Run("TxOffset is 0 if VE is not enabled", func() {
		decoder := process.NewSlinkyMarketPriceDecoder(suite.decoder, suite.gen)

		suite.ctx = testutils.CreateBaseSDKContext(suite.T())
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(2)

		offset := decoder.GetTxOffset(suite.ctx)
		suite.Equal(0, offset)
	})

	suite.Run("TxOffset is constants.OracleVEInjectedTxs if VE is enabled", func() {
		decoder := process.NewSlinkyMarketPriceDecoder(
			process.NewDefaultUpdateMarketPriceTxDecoder(nil, nil), nil) // ignore deps

		suite.ctx = testutils.CreateBaseSDKContext(suite.T())
		suite.ctx = testutils.UpdateContextWithVEHeight(suite.ctx, 4)
		suite.ctx = suite.ctx.WithBlockHeight(5)

		offset := decoder.GetTxOffset(suite.ctx)
		suite.Equal(constants.OracleVEInjectedTxs, offset)
	})
}
