package types_test

import (
	"errors"
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/client"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExchangeToMarketPrices_IsEmpty(t *testing.T) {
	exchangeToMarketPrices := getNewExchangeToMarketPricesAndCheckForError(
		t,
		constants.Exchange1Exchange2Array,
		nil,
	)

	require.Empty(
		t,
		exchangeToMarketPrices.ExchangeMarketPrices[constants.ExchangeId1].MarketToPriceTimestamp,
	)
	require.Empty(
		t,
		exchangeToMarketPrices.ExchangeMarketPrices[constants.ExchangeId2].MarketToPriceTimestamp,
	)

	exchangePrices := exchangeToMarketPrices.GetAllPrices()
	require.Len(t, exchangePrices, 2)
	require.Empty(t, exchangePrices[constants.ExchangeId1])
	require.Empty(t, exchangePrices[constants.ExchangeId2])
}

func TestNewExchangeToMarketPrices_InvalidWithNoExchangeIds(t *testing.T) {
	getNewExchangeToMarketPricesAndCheckForError(
		t,
		[]types.ExchangeId{},
		errors.New("exchangeIds must not be empty"),
	)
}

func TestNewExchangeToMarketPrices_InvalidWithDuplicateExchangeIds(t *testing.T) {
	getNewExchangeToMarketPricesAndCheckForError(
		t,
		[]types.ExchangeId{
			constants.ExchangeId1,
			constants.ExchangeId2,
			constants.ExchangeId1,
		},
		fmt.Errorf("exchangeId: '%v' appears twice in request", constants.ExchangeId1),
	)
}

func TestUpdatePrice_IsValid(t *testing.T) {
	exchangeToMarketPrices := getNewExchangeToMarketPricesAndCheckForError(
		t,
		constants.Exchange1Exchange2Array,
		nil,
	)

	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market9_TimeT_Price1,
		false,
	)

	mtpMap := exchangeToMarketPrices.GetAllPrices()
	require.Len(t, mtpMap[constants.ExchangeId1], 1)
	require.Empty(t, mtpMap[constants.ExchangeId2])

	marketPriceTimestamp := mtpMap[constants.ExchangeId1][0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateIsValid(t *testing.T) {
	exchangeToMarketPrices := getNewExchangeToMarketPricesAndCheckForError(
		t,
		constants.Exchange1Exchange2Array,
		nil,
	)

	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market9_TimeTMinusThreshold_Price2,
		false,
	)
	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market9_TimeT_Price1,
		false,
	)

	mtpMap := exchangeToMarketPrices.GetAllPrices()

	require.Len(t, mtpMap[constants.ExchangeId1], 1)

	marketPriceTimestamp := mtpMap[constants.ExchangeId1][0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_UpdateIsInvalid(t *testing.T) {
	exchangeToMarketPrices := getNewExchangeToMarketPricesAndCheckForError(
		t,
		constants.Exchange1Exchange2Array,
		nil,
	)

	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market9_TimeT_Price1,
		false,
	)
	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market9_TimeTMinusThreshold_Price2,
		false,
	)

	mtpMap := exchangeToMarketPrices.GetAllPrices()

	require.Len(t, mtpMap[constants.ExchangeId1], 1)

	marketPriceTimestamp := mtpMap[constants.ExchangeId1][0]
	require.Equal(t, constants.MarketId9, marketPriceTimestamp.MarketId)
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)
}

func TestUpdatePrice_IsValidForTwoMarkets(t *testing.T) {
	exchangeToMarketPrices := getNewExchangeToMarketPricesAndCheckForError(
		t,
		constants.Exchange1Exchange2Array,
		nil,
	)

	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market9_TimeT_Price1,
		false,
	)
	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market8_TimeTMinusThreshold_Price2,
		false,
	)

	mtpMap := exchangeToMarketPrices.GetAllPrices()

	require.Len(t, mtpMap[constants.ExchangeId1], 2)
	assert.ElementsMatch(
		t,
		[]types.MarketPriceTimestamp{
			{
				MarketId:      constants.MarketId9,
				Price:         constants.Price1,
				LastUpdatedAt: constants.TimeT,
			},
			{
				MarketId:      constants.MarketId8,
				Price:         constants.Price2,
				LastUpdatedAt: constants.TimeTMinusThreshold,
			},
		},
		mtpMap[constants.ExchangeId1],
	)
}

func TestUpdatePrice_IsValidForTwoExchanges(t *testing.T) {
	exchangeToMarketPrices := getNewExchangeToMarketPricesAndCheckForError(
		t,
		constants.Exchange1Exchange2Array,
		nil,
	)

	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId1,
		constants.Market9_TimeT_Price1,
		false,
	)
	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId2,
		constants.Market8_TimeTMinusThreshold_Price2,
		false,
	)

	mtpMap := exchangeToMarketPrices.GetAllPrices()

	require.Len(t, mtpMap[constants.ExchangeId1], 1)
	require.Len(t, mtpMap[constants.ExchangeId2], 1)

	marketPriceTimestamp := mtpMap[constants.ExchangeId1][0]
	require.Equal(t, constants.Price1, marketPriceTimestamp.Price)
	require.Equal(t, constants.TimeT, marketPriceTimestamp.LastUpdatedAt)

	marketPriceTimestamp2 := mtpMap[constants.ExchangeId2][0]
	require.Equal(t, constants.Price2, marketPriceTimestamp2.Price)
	require.Equal(t, constants.TimeTMinusThreshold, marketPriceTimestamp2.LastUpdatedAt)
}

func TestNewExchangeToMarketPrices_UpdateIsInvalidForInvalidExchange(t *testing.T) {
	exchangeToMarketPrices := getNewExchangeToMarketPricesAndCheckForError(
		t,
		constants.Exchange1Exchange2Array,
		nil,
	)

	updatePriceAndCheckForPanic(
		t,
		exchangeToMarketPrices,
		constants.ExchangeId3,
		constants.Market8_TimeTMinusThreshold_Price2,
		true,
	)
}

func TestGetIndexPrice_Mixed(t *testing.T) {
	tests := map[string]struct {
		initialPrices []*client.ExchangeIdMarketPriceTimestamp
		market        types.MarketId
		cutoffTime    time.Time

		expectedMedianPrice         uint64
		expectedNumPricesMedianized int
	}{
		"invalid: no prices": {
			market:                      constants.MarketId9,
			cutoffTime:                  constants.TimeT,
			expectedMedianPrice:         0,
			expectedNumPricesMedianized: 0,
		},
		"invalid: no prices for market": {
			initialPrices: []*client.ExchangeIdMarketPriceTimestamp{
				constants.ExchangeId2_Market8_TimeT_Price2, // Valid timestamp, wrong market.
			},
			market:                      constants.MarketId9,
			cutoffTime:                  constants.TimeTMinus1,
			expectedMedianPrice:         0,
			expectedNumPricesMedianized: 0,
		},
		"valid: 1 price": {
			initialPrices: []*client.ExchangeIdMarketPriceTimestamp{
				constants.ExchangeId2_Market9_TimeT_Price2,
			},
			market:                      constants.MarketId9,
			cutoffTime:                  constants.TimeTMinus1,
			expectedMedianPrice:         constants.Price2,
			expectedNumPricesMedianized: 1,
		},
		"valid: 1 current price, 1 stale price": {
			initialPrices: []*client.ExchangeIdMarketPriceTimestamp{
				constants.ExchangeId2_Market8_TimeT_Price2,       // Valid timestamp, same market.
				constants.ExchangeId1_Market8_BeforeTimeT_Price3, // Stale timestamp, same market.
				constants.ExchangeId1_Market9_TimeT_Price1,       // Valid timestamp, different market.
			},
			market:                      constants.MarketId8,
			cutoffTime:                  constants.TimeTMinus1,
			expectedMedianPrice:         constants.Price2,
			expectedNumPricesMedianized: 1,
		},
		"valid: multiple prices": {
			initialPrices: []*client.ExchangeIdMarketPriceTimestamp{
				constants.ExchangeId1_Market9_TimeT_Price1,
				constants.ExchangeId2_Market9_TimeT_Price2,
				constants.ExchangeId3_Market9_TimeT_Price3,
			},
			market:                      constants.MarketId9,
			cutoffTime:                  constants.TimeTMinus1,
			expectedMedianPrice:         constants.Price2,
			expectedNumPricesMedianized: 3,
		},
	}

	testExchanges := []types.ExchangeId{constants.ExchangeId1, constants.ExchangeId2, constants.ExchangeId3}

	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			// Setup.
			etmp := getNewExchangeToMarketPricesAndCheckForError(t, testExchanges, nil)

			// Update prices with initial prices.
			for _, exchangeMarketPriceTimestamp := range tc.initialPrices {
				exchange := exchangeMarketPriceTimestamp.ExchangeId
				marketPriceTimestamp := exchangeMarketPriceTimestamp.MarketPriceTimestamp
				etmp.UpdatePrice(exchange, marketPriceTimestamp)
			}

			// Execute.
			resolver := lib.Median[uint64]
			medianPrice, numPricesMedianized := etmp.GetIndexPrice(tc.market, tc.cutoffTime, resolver)

			// Assert.
			require.Equal(t, tc.expectedMedianPrice, medianPrice)
			require.Equal(t, tc.expectedNumPricesMedianized, numPricesMedianized)
		})
	}
}

func updatePriceAndCheckForPanic(
	t *testing.T,
	exchangeToMarketPrices types.ExchangeToMarketPrices,
	exchangeId types.ExchangeId,
	marketPriceTimestamp *types.MarketPriceTimestamp,
	panics bool,
) {
	if panics {
		require.Panics(
			t,
			func() {
				exchangeToMarketPrices.UpdatePrice(
					exchangeId,
					marketPriceTimestamp,
				)
			},
		)
	} else {
		require.NotPanics(
			t,
			func() {
				exchangeToMarketPrices.UpdatePrice(
					exchangeId,
					marketPriceTimestamp,
				)
			},
		)
	}
}

func getNewExchangeToMarketPricesAndCheckForError(
	t *testing.T,
	exchangeIds []types.ExchangeId,
	expectedErr error,
) *types.ExchangeToMarketPricesImpl {
	exchangeToMarketPrices, err := types.NewExchangeToMarketPrices(exchangeIds)

	if expectedErr != nil {
		require.EqualError(t, err, expectedErr.Error())
		return nil
	}

	return exchangeToMarketPrices.(*types.ExchangeToMarketPricesImpl)
}
