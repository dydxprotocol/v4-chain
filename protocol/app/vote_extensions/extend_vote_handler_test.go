package vote_extensions

import (
	"fmt"
	"testing"

	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/app/process"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
)

func TestExtendVoteHandlerDeecodeMarketPricesFailure(t *testing.T) {
	slinkyEvh := mocks.NewExtendVoteHandler(t)
	pricesTxDecoder := mocks.NewUpdateMarketPriceTxDecoder(t)
	pricesKeeper := mocks.NewPricesKeeper(t)
	evh := ExtendVoteHandler{
		SlinkyExtendVoteHandler: slinkyEvh.Execute,
		PricesTxDecoder:         pricesTxDecoder,
		PricesKeeper:            pricesKeeper,
	}

	pricesTxDecoder.On("DecodeUpdateMarketPricesTx", mock.Anything, mock.Anything).Return(
		nil, fmt.Errorf("foobar"))
	_, err := evh.ExtendVoteHandler()(sdk.Context{}, &cometabci.RequestExtendVote{Txs: make([][]byte, 0)})

	require.ErrorContains(t, err, "DecodeMarketPricesTx failure foobar")
	pricesTxDecoder.AssertExpectations(t)
	pricesKeeper.AssertExpectations(t)
	slinkyEvh.AssertExpectations(t)
}

func TestExtendVoteHandlerUpdatePricesTxValidateFailure(t *testing.T) {
	slinkyEvh := mocks.NewExtendVoteHandler(t)
	pricesTxDecoder := mocks.NewUpdateMarketPriceTxDecoder(t)
	pricesKeeper := mocks.NewPricesKeeper(t)
	evh := ExtendVoteHandler{
		SlinkyExtendVoteHandler: slinkyEvh.Execute,
		PricesTxDecoder:         pricesTxDecoder,
		PricesKeeper:            pricesKeeper,
	}

	pricesTxDecoder.On("DecodeUpdateMarketPricesTx", mock.Anything, mock.Anything).Return(
		process.NewUpdateMarketPricesTx(sdk.Context{}, pricesKeeper, constants.InvalidMsgUpdateMarketPricesStateless),
		nil)
	_, err := evh.ExtendVoteHandler()(sdk.Context{}, &cometabci.RequestExtendVote{Txs: make([][]byte, 0)})

	require.ErrorContains(t, err, "updatePricesTx.Validate failure")
	pricesTxDecoder.AssertExpectations(t)
	pricesKeeper.AssertExpectations(t)
	slinkyEvh.AssertExpectations(t)
}

func TestExtendVoteHandlerUpdateMarketPricesError(t *testing.T) {
	slinkyEvh := mocks.NewExtendVoteHandler(t)
	pricesTxDecoder := mocks.NewUpdateMarketPriceTxDecoder(t)
	pricesKeeper := mocks.NewPricesKeeper(t)
	evh := ExtendVoteHandler{
		SlinkyExtendVoteHandler: slinkyEvh.Execute,
		PricesTxDecoder:         pricesTxDecoder,
		PricesKeeper:            pricesKeeper,
	}

	pricesTxDecoder.On("DecodeUpdateMarketPricesTx", mock.Anything, mock.Anything).Return(
		process.NewUpdateMarketPricesTx(sdk.Context{}, pricesKeeper, constants.EmptyMsgUpdateMarketPrices),
		nil)
	pricesKeeper.On("PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	pricesKeeper.On("UpdateMarketPrices", mock.Anything, mock.Anything).
		Return(fmt.Errorf(""))
	_, err := evh.ExtendVoteHandler()(sdk.Context{}, &cometabci.RequestExtendVote{Txs: make([][]byte, 0)})

	require.ErrorContains(t, err, "failed to update market prices in extend vote handler pre-slinky invocation")
	pricesTxDecoder.AssertExpectations(t)
	pricesKeeper.AssertExpectations(t)
	slinkyEvh.AssertExpectations(t)
}

func TestExtendVoteHandlerSlinkyFailure(t *testing.T) {
	slinkyEvh := mocks.NewExtendVoteHandler(t)
	pricesTxDecoder := mocks.NewUpdateMarketPriceTxDecoder(t)
	pricesKeeper := mocks.NewPricesKeeper(t)
	evh := ExtendVoteHandler{
		SlinkyExtendVoteHandler: slinkyEvh.Execute,
		PricesTxDecoder:         pricesTxDecoder,
		PricesKeeper:            pricesKeeper,
	}

	pricesTxDecoder.On("DecodeUpdateMarketPricesTx", mock.Anything, mock.Anything).Return(
		process.NewUpdateMarketPricesTx(sdk.Context{}, pricesKeeper, constants.EmptyMsgUpdateMarketPrices),
		nil)
	pricesKeeper.On("PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	pricesKeeper.On("UpdateMarketPrices", mock.Anything, mock.Anything).Return(nil)
	slinkyEvh.On("Execute", mock.Anything, mock.Anything).
		Return(&cometabci.ResponseExtendVote{}, fmt.Errorf("slinky failure"))
	_, err := evh.ExtendVoteHandler()(sdk.Context{}, &cometabci.RequestExtendVote{Txs: make([][]byte, 0)})

	require.ErrorContains(t, err, "slinky failure")
	pricesTxDecoder.AssertExpectations(t)
	pricesKeeper.AssertExpectations(t)
	slinkyEvh.AssertExpectations(t)
}

func TestExtendVoteHandlerSlinkySuccess(t *testing.T) {
	slinkyEvh := mocks.NewExtendVoteHandler(t)
	pricesTxDecoder := mocks.NewUpdateMarketPriceTxDecoder(t)
	pricesKeeper := mocks.NewPricesKeeper(t)
	evh := ExtendVoteHandler{
		SlinkyExtendVoteHandler: slinkyEvh.Execute,
		PricesTxDecoder:         pricesTxDecoder,
		PricesKeeper:            pricesKeeper,
	}

	pricesTxDecoder.On("DecodeUpdateMarketPricesTx", mock.Anything, mock.Anything).Return(
		process.NewUpdateMarketPricesTx(sdk.Context{}, pricesKeeper, constants.EmptyMsgUpdateMarketPrices),
		nil)
	pricesKeeper.On("PerformStatefulPriceUpdateValidation", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	pricesKeeper.On("UpdateMarketPrices", mock.Anything, mock.Anything).Return(nil)
	slinkyEvh.On("Execute", mock.Anything, mock.Anything).
		Return(&cometabci.ResponseExtendVote{}, nil)
	_, err := evh.ExtendVoteHandler()(sdk.Context{}, &cometabci.RequestExtendVote{Txs: make([][]byte, 0)})

	require.NoError(t, err)
	pricesTxDecoder.AssertExpectations(t)
	pricesKeeper.AssertExpectations(t)
	slinkyEvh.AssertExpectations(t)
}
