package gate_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/price_function/gate"
	"github.com/stretchr/testify/require"
)

func TestGateUrl(t *testing.T) {
	require.Equal(t, "https://api.gateio.ws/api/v4/spot/tickers", gate.GateDetails.Url)
}

func TestGateIsMultiMarket(t *testing.T) {
	require.True(t, gate.GateDetails.IsMultiMarket)
}
