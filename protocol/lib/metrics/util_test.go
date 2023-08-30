package metrics

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncrCountMetricWithLabelsDoesntPanic(t *testing.T) {
	require.NotPanics(t, func() {
		IncrCountMetricWithLabels("module", "metric", NewBinaryStringLabel("label", true))
	})
}

func TestNewBinaryStringLabel(t *testing.T) {
	tests := map[string]struct {
		name               string
		condition          bool
		expectedLabelValue string
	}{
		"positive condition": {
			name:               "labelname",
			condition:          true,
			expectedLabelValue: Yes,
		},
		"negative condition": {
			name:               "labelname",
			condition:          false,
			expectedLabelValue: No,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			label := NewBinaryStringLabel(tc.name, tc.condition)
			require.Equal(t, tc.name, label.Name)
			require.Equal(t, tc.expectedLabelValue, label.Value)
		})
	}
}

func TestGetMetricValueFromBigInt(t *testing.T) {
	tests := map[string]struct {
		input    *big.Int
		expected float32
	}{
		"zero": {
			input:    big.NewInt(0),
			expected: float32(0),
		},
		"positive": {
			input:    big.NewInt(1234),
			expected: float32(1234),
		},
		"negative": {
			input:    big.NewInt(-1234),
			expected: float32(-1234),
		},
		"overflow": {
			input:    new(big.Int).SetUint64(math.MaxUint64),
			expected: float32(1.8446744e+19),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, GetMetricValueFromBigInt(tc.input))
		})
	}
}
