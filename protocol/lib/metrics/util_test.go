package metrics

import (
	"github.com/stretchr/testify/require"
	"testing"
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
