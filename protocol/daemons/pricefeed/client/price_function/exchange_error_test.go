package price_function_test

import (
	"errors"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExchangeError(t *testing.T) {
	error := price_function.NewExchangeError("exchange", "error")
	var exchangeError price_function.ExchangeError
	found := errors.As(error, &exchangeError)
	require.True(t, found)
	require.Equal(t, error, exchangeError)
	require.Equal(t, "exchange", exchangeError.GetExchangeId())
}
