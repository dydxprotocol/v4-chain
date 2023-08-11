package metrics_test

import (
	"testing"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
	pricefeedmetrics "github.com/dydxprotocol/v4/daemons/pricefeed/metrics"
	"github.com/dydxprotocol/v4/lib/metrics"
	"github.com/stretchr/testify/require"
)

const (
	INVALID_ID = 10000000
)

func TestGetLabelForMarketIdSuccess(t *testing.T) {
	require.Equal(
		t,
		metrics.GetLabelForStringValue(metrics.MarketId, "BTCUSD"),
		pricefeedmetrics.GetLabelForMarketId(exchange_common.MARKET_BTC_USD),
	)
}

func TestGetLabelForMarketIdFailure(t *testing.T) {
	require.Equal(
		t,
		metrics.GetLabelForStringValue(metrics.MarketId, pricefeedmetrics.INVALID),
		pricefeedmetrics.GetLabelForMarketId(INVALID_ID),
	)
}

func TestGetLabelForExchangeId(t *testing.T) {
	require.Equal(
		t,
		metrics.GetLabelForStringValue(metrics.ExchangeId, exchange_common.EXCHANGE_NAME_BINANCEUS),
		pricefeedmetrics.GetLabelForExchangeId(exchange_common.EXCHANGE_ID_BINANCE_US),
	)
}
