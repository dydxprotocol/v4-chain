package testutil

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client/constants/exchange_common"
)

const (
	ETHUSDC = "ETHUSDC"
	BTCUSDC = "BTCUSDC"
)

var (
	ExponentSymbolMap = map[string]int32{
		ETHUSDC: constants.StaticMarketPriceExponent[exchange_common.MARKET_ETH_USD],
	}

	MedianizationError = errors.New("Failed to get median")
)

func CreateResponseFromJson(m string) *http.Response {
	jsonBlob := bytes.NewReader([]byte(m))
	return &http.Response{Body: io.NopCloser(jsonBlob)}
}
