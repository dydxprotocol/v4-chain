package types_test

import (
	"fmt"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExchangeConfigJsonValidate_Mixed(t *testing.T) {
	tests := map[string]struct {
		exchangeConfigJson types.ExchangeConfigJson
		expectedErr        error
	}{
		"Valid": {
			exchangeConfigJson: types.ExchangeConfigJson{
				Exchanges: []types.ExchangeMarketConfigJson{
					{
						ExchangeName: "binance",
						Ticker:       "BTC-USDT",
					},
				},
			},
		},
		"Invalid - no exchanges": {
			exchangeConfigJson: types.ExchangeConfigJson{},
			expectedErr:        fmt.Errorf("exchanges cannot be empty"),
		},
		"Invalid - invalid exchange": {
			exchangeConfigJson: types.ExchangeConfigJson{
				Exchanges: []types.ExchangeMarketConfigJson{
					{
						ExchangeName: "not-a-real-exchange", // invalid
						Ticker:       "BTC-USDT",
					},
				},
			},
			expectedErr: fmt.Errorf("invalid exchange: exchange name 'not-a-real-exchange' is not valid"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.exchangeConfigJson.Validate(
				[]types.ExchangeId{"binance"},
				map[string]types.MarketId{},
			)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr.Error())
			}
		})
	}
}
