package types_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExchangeMarketConfigJsonValidate_Mixed(t *testing.T) {
	tests := map[string]struct {
		exchangeMarketConfigJson types.ExchangeMarketConfigJson
		expectedErr              error
	}{
		"Valid": {
			exchangeMarketConfigJson: types.ExchangeMarketConfigJson{
				ExchangeName: "binance",
				Ticker:       "BTC-USDT",
			},
		},
		"Valid with adjust-by market": {
			exchangeMarketConfigJson: types.ExchangeMarketConfigJson{
				ExchangeName:   "binance",
				Ticker:         "BTC-USDT",
				AdjustByMarket: "ABC-USD",
			},
		},
		"Invalid - no exchange name": {
			exchangeMarketConfigJson: types.ExchangeMarketConfigJson{
				Ticker: "BTC-USDT",
			},
			expectedErr: fmt.Errorf("exchange name cannot be empty"),
		},
		"Invalid - invalid exchange name": {
			exchangeMarketConfigJson: types.ExchangeMarketConfigJson{
				ExchangeName: "not-a-real-exchange", // invalid
				Ticker:       "BTC-USDT",
			},
			expectedErr: fmt.Errorf("exchange name 'not-a-real-exchange' is not valid"),
		},
		"Invalid - no ticker": {
			exchangeMarketConfigJson: types.ExchangeMarketConfigJson{
				ExchangeName: "binance",
			},
			expectedErr: fmt.Errorf("ticker cannot be empty"),
		},
		"Invalid - adjust by market not valid": {
			exchangeMarketConfigJson: types.ExchangeMarketConfigJson{
				ExchangeName:   "binance",
				Ticker:         "BTC-USDT",
				AdjustByMarket: "XYZ-USD", // invalid
			},
			expectedErr: fmt.Errorf("adjustment market 'XYZ-USD' is not valid"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.exchangeMarketConfigJson.Validate(
				[]types.ExchangeId{"binance"},
				map[string]types.MarketId{
					"ABC-USD": 3,
				},
			)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
