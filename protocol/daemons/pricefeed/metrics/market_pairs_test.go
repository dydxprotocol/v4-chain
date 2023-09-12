package metrics_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/metrics"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetMarketPairForTelemetry(t *testing.T) {
	tests := map[string]struct {
		marketId uint32
		expected string
	}{
		"present id": {
			marketId: 1,
			expected: "BTC-USD",
		},
		"absent id": {
			marketId: 99,
			expected: "INVALID",
		},
	}
	metrics.AddMarketPairForTelemetry(1, "BTC-USD")
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actual := metrics.GetMarketPairForTelemetry(tc.marketId)
			require.Equal(t, tc.expected, actual)
		})
	}
}
