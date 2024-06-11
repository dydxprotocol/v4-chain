package perpetuals_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/perpetuals"
	"github.com/stretchr/testify/require"
)

func TestMustHumanSizeToBaseQuantums(t *testing.T) {
	tests := []struct {
		humanSize        string
		atomicResolution int32
		expected         uint64
		expectPanic      bool
	}{
		{"1.123", -8, 112_300_000, false},
		{"0.55", -9, 550_000_000, false},
		{"0.00000001", -8, 1, false},
		{"235", -1, 2350, false},
		{"235", 1, 23, false},
		{"1", -10, 10_000_000_000, false},
		{"0.0000000001", -10, 1, false},
		{"abc", -8, 0, true}, // Invalid humanSize
	}

	for _, test := range tests {
		if test.expectPanic {
			require.Panics(t, func() {
				perpetuals.MustHumanSizeToBaseQuantums(test.humanSize, test.atomicResolution)
			}, "For humanSize %v and atomicResolution %v, expected a panic", test.humanSize, test.atomicResolution)
		} else {
			result := perpetuals.MustHumanSizeToBaseQuantums(test.humanSize, test.atomicResolution)
			require.Equal(t,
				test.expected,
				result,
				"expected = %v, result =%v, for humanSize %v and atomicResolution %v",
				test.expected,
				result,
				test.humanSize,
				test.atomicResolution,
			)
		}
	}
}
