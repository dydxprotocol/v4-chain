package server_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	pricefeedconstants "github.com/dydxprotocol/v4/daemons/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/api"
	"github.com/dydxprotocol/v4/daemons/server"
	pricefeedtypes "github.com/dydxprotocol/v4/daemons/server/types/pricefeed"
	"github.com/dydxprotocol/v4/daemons/types"
	"github.com/dydxprotocol/v4/mocks"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/stretchr/testify/require"
)

func TestUpdateMarketPrices_Valid(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}

	s := createServerWithMocks(
		mockGrpcServer,
		mockFileHandler,
	).WithPriceFeedMarketToExchangePrices(
		pricefeedtypes.NewMarketToExchangePrices(),
	)

	sendAndCheckPriceUpdate(
		t,
		s,
		&api.UpdateMarketPricesRequest{MarketPriceUpdates: constants.AtTimeTPriceUpdate},
		nil,
	)
}

func TestUpdateMarketPrices_NotInitialized(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}

	// Create a new server without initializing `MarketToExchange` field.
	s := createServerWithMocks(
		mockGrpcServer,
		mockFileHandler,
	)

	req := &api.UpdateMarketPricesRequest{MarketPriceUpdates: constants.AtTimeTPriceUpdate}
	require.PanicsWithError(
		t,
		sdkerrors.Wrapf(
			types.ErrServerNotInitializedCorrectly,
			"MarketToExchange not initialized",
		).Error(),
		func() {
			_, err := s.UpdateMarketPrices(
				context.TODO(),
				req,
			)
			require.NoError(t, err, "No error expected when sending Price Update")
		},
	)
}

func TestUpdateMarketPrices_InvalidEmptyRequest(t *testing.T) {
	mockGrpcServer := &mocks.GrpcServer{}
	mockFileHandler := &mocks.FileHandler{}

	s := createServerWithMocks(
		mockGrpcServer,
		mockFileHandler,
	).WithPriceFeedMarketToExchangePrices(
		pricefeedtypes.NewMarketToExchangePrices(),
	)

	sendAndCheckPriceUpdate(
		t,
		s,
		&api.UpdateMarketPricesRequest{MarketPriceUpdates: []*api.MarketPriceUpdate{}},
		types.ErrPriceFeedMarketPriceUpdateEmpty,
	)
}

func TestUpdateMarketPrices_InvalidExchangePrices(t *testing.T) {
	tests := map[string]struct {
		input         api.UpdateMarketPricesRequest
		expectedError error
	}{
		"UpdateMarketPricesRequest Has ExchangePrice With 0 Price": {
			input: api.UpdateMarketPricesRequest{MarketPriceUpdates: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId9,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeFeedId: constants.ExchangeFeedId1,
							Price:          constants.InvalidPrice,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			}},
			expectedError: fmt.Errorf(
				"Price is set to %d which is not a valid price",
				constants.InvalidPrice,
			),
		},
		"UpdateMarketPricesRequest Has ExchangePrice With Price Not Set": {
			input: api.UpdateMarketPricesRequest{MarketPriceUpdates: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId9,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeFeedId: constants.ExchangeFeedId1,
							LastUpdateTime: &constants.TimeT,
						},
					},
				},
			}},
			expectedError: fmt.Errorf(
				"Price is set to %d which is not a valid price",
				pricefeedconstants.DefaultPrice,
			),
		},
		"UpdateMarketPricesRequest Has ExchangePrice With LastUpdatedAt Not Set": {
			input: api.UpdateMarketPricesRequest{MarketPriceUpdates: []*api.MarketPriceUpdate{
				{
					MarketId: constants.MarketId9,
					ExchangePrices: []*api.ExchangePrice{
						{
							ExchangeFeedId: constants.ExchangeFeedId1,
							Price:          constants.Price1,
						},
					},
				},
			}},
			expectedError: errors.New("LastUpdateTime is not set"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mockGrpcServer := &mocks.GrpcServer{}
			mockFileHandler := &mocks.FileHandler{}

			s := createServerWithMocks(
				mockGrpcServer,
				mockFileHandler,
			).WithPriceFeedMarketToExchangePrices(
				pricefeedtypes.NewMarketToExchangePrices(),
			)
			expectedErr := sdkerrors.Wrapf(
				tc.expectedError,
				"ExchangePrice: %v and MarketId: %d",
				// Assumes first ExchangePrice is the one with a validation error.
				tc.input.MarketPriceUpdates[0].ExchangePrices[0],
				tc.input.MarketPriceUpdates[0].MarketId,
			)

			sendAndCheckPriceUpdate(t, s, &tc.input, expectedErr)
		})
	}
}

func sendAndCheckPriceUpdate(
	t *testing.T,
	s *server.Server,
	req *api.UpdateMarketPricesRequest,
	expectedErr error,
) {
	apiResponse, err := s.UpdateMarketPrices(
		context.TODO(),
		req,
	)

	if expectedErr != nil {
		require.EqualError(t, err, expectedErr.Error())
	} else {
		require.NoError(t, err, "No error expected when sending Price Update")
		require.Equal(
			t,
			&api.UpdateMarketPricesResponse{},
			apiResponse,
			"response should be a pointer to a valid UpdateMarketPricesResponse",
		)
	}
}
