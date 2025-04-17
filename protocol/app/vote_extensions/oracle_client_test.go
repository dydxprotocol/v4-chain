package vote_extensions

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/dydxprotocol/slinky/pkg/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

func TestStartStopNoop(t *testing.T) {
	cli := NewOraclePrices(nil)

	err := cli.Start(context.TODO())
	require.NoError(t, err)
	err = cli.Stop()
	require.NoError(t, err)
}

func TestValidPriceResponse(t *testing.T) {
	pk := mocks.NewPricesKeeper(t)
	cli := NewOraclePrices(pk)
	pk.On("GetValidMarketPriceUpdates", mock.Anything).
		Return(&types.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
				{MarketId: 2, Price: 1},
			},
		}).Once()
	pk.On("GetCurrencyPairFromID", mock.Anything, mock.Anything).
		Return(
			oracletypes.CurrencyPair{Base: "FOO", Quote: "BAR"},
			true,
		).Once()

	_, err := cli.Prices(sdk.Context{}.WithLogger(log.NewNopLogger()), nil)

	require.NoError(t, err)
	pk.AssertExpectations(t)
}

func TestNonSdkContextFails(t *testing.T) {
	pk := mocks.NewPricesKeeper(t)
	cli := NewOraclePrices(pk)

	_, err := cli.Prices(context.TODO(), nil)

	require.Error(t, err)
}

func TestEmptyUpdatesPasses(t *testing.T) {
	pk := mocks.NewPricesKeeper(t)
	cli := NewOraclePrices(pk)
	pk.On("GetValidMarketPriceUpdates", mock.Anything).
		Return(&types.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{},
		}).Once()

	_, err := cli.Prices(sdk.Context{}.WithLogger(log.NewNopLogger()), nil)

	require.NoError(t, err)
	pk.AssertExpectations(t)
}

func TestNilUpdatesPasses(t *testing.T) {
	pk := mocks.NewPricesKeeper(t)
	cli := NewOraclePrices(pk)
	pk.On("GetValidMarketPriceUpdates", mock.Anything).
		Return(nil).Once()

	_, err := cli.Prices(sdk.Context{}.WithLogger(log.NewNopLogger()), nil)

	require.NoError(t, err)
	pk.AssertExpectations(t)
}

func TestPairNotFoundNoOps(t *testing.T) {
	pk := mocks.NewPricesKeeper(t)
	cli := NewOraclePrices(pk)
	pk.On("GetValidMarketPriceUpdates", mock.Anything).
		Return(&types.MsgUpdateMarketPrices{
			MarketPriceUpdates: []*types.MsgUpdateMarketPrices_MarketPrice{
				{MarketId: 2, Price: 1},
			},
		}).Once()
	pk.On("GetCurrencyPairFromID", mock.Anything, mock.Anything).
		Return(
			oracletypes.CurrencyPair{Base: "", Quote: ""},
			false,
		).Once()

	_, err := cli.Prices(sdk.Context{}.WithLogger(log.NewNopLogger()), nil)

	require.NoError(t, err)
	pk.AssertExpectations(t)
}
