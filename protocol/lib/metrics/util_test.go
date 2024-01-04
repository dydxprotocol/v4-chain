package metrics_test

import (
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	gometrics "github.com/hashicorp/go-metrics"

	"github.com/stretchr/testify/require"
)

func TestIncrCountMetricWithLabels(t *testing.T) {
	t.Cleanup(gometrics.Shutdown)

	conf := gometrics.DefaultConfig("testService")
	sink := gometrics.NewInmemSink(time.Hour, time.Hour)
	_, err := gometrics.NewGlobal(conf, sink)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		metrics.IncrCountMetricWithLabels(
			"testModule",
			"testMetric",
			gometrics.Label{
				Name:  "testLabel",
				Value: "testLabelValue",
			},
		)
	}

	found := false
	for _, metrics := range sink.Data() {
		metrics.RLock()
		defer metrics.RUnlock()

		if metric, ok := metrics.Counters["testService.testModule.testMetric.count;testLabel=testLabelValue"]; ok {
			require.Equal(t,
				[]gometrics.Label{{
					Name:  "testLabel",
					Value: "testLabelValue",
				}},
				metric.Labels)
			require.Equal(t, 3, metric.Count)
			require.Equal(t, float64(3), metric.Sum)
			found = true
		}
	}
	require.True(t, found)
}

func TestIncrCountMetricWithLabelsDoesntPanic(t *testing.T) {
	require.NotPanics(t, func() {
		metrics.IncrCountMetricWithLabels("module", "metric", metrics.GetLabelForBoolValue("label", true))
	})
}

func TestGetLabelForBoolValue(t *testing.T) {
	tests := map[string]struct {
		name               string
		condition          bool
		expectedLabelValue string
	}{
		"true": {
			name:               "labelname",
			condition:          true,
			expectedLabelValue: "true",
		},
		"false": {
			name:               "labelname",
			condition:          false,
			expectedLabelValue: "false",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			label := metrics.GetLabelForBoolValue(tc.name, tc.condition)
			require.Equal(t, tc.name, label.Name)
			require.Equal(t, tc.expectedLabelValue, label.Value)
		})
	}
}

func TestGetLabelForIntValue(t *testing.T) {
	tests := map[string]struct {
		name               string
		value              int
		expectedLabelValue string
	}{
		"min": {
			name:               "labelname",
			value:              math.MinInt,
			expectedLabelValue: "-9223372036854775808",
		},
		"negative": {
			name:               "labelname",
			value:              -1,
			expectedLabelValue: "-1",
		},
		"zero": {
			name:               "labelname",
			value:              0,
			expectedLabelValue: "0",
		},
		"positive": {
			name:               "labelname",
			value:              1,
			expectedLabelValue: "1",
		},
		"max": {
			name:               "labelname",
			value:              math.MaxInt,
			expectedLabelValue: "9223372036854775807",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			label := metrics.GetLabelForIntValue(tc.name, tc.value)
			require.Equal(t, tc.name, label.Name)
			require.Equal(t, tc.expectedLabelValue, label.Value)
		})
	}
}

func TestGetLabelForStringValue(t *testing.T) {
	tests := map[string]struct {
		name               string
		value              string
		expectedLabelValue string
	}{
		"empty": {
			name:               "labelname",
			value:              "",
			expectedLabelValue: "",
		},
		"short string": {
			name:               "labelname",
			value:              "abc",
			expectedLabelValue: "abc",
		},
		"long string": {
			name:               "labelname",
			value:              "abc def 2389209 lsdf ;'sdf';s#2",
			expectedLabelValue: "abc def 2389209 lsdf ;'sdf';s#2",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			label := metrics.GetLabelForStringValue(tc.name, tc.value)
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
		"underflow": {
			input:    big.NewInt(math.MinInt),
			expected: float32(-9.223372e+18),
		},
		"overflow": {
			input:    new(big.Int).SetUint64(math.MaxUint64),
			expected: float32(1.8446744e+19),
		},
		"overflow: 1234567 * 1e24": {
			input:    big_testutil.Int64MulPow10(1234567, 24), // 1234567 * 1e24
			expected: float32(1.234567e+30),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, metrics.GetMetricValueFromBigInt(tc.input))
		})
	}
}

func TestModuleMeasureSinceWithLabels(t *testing.T) {
	t.Cleanup(gometrics.Shutdown)

	conf := gometrics.DefaultConfig("testService")
	sink := gometrics.NewInmemSink(time.Hour, time.Hour)
	_, err := gometrics.NewGlobal(conf, sink)
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		metrics.ModuleMeasureSinceWithLabels(
			"testModule",
			[]string{"testKey1", "testKey2"},
			time.Now(),
			[]gometrics.Label{{
				Name:  "testLabel",
				Value: "testLabelValue",
			}},
		)
	}

	found := false
	for _, metrics := range sink.Data() {
		metrics.RLock()
		defer metrics.RUnlock()

		if metric, ok := metrics.Samples["testService.testKey1.testKey2;module=testModule;testLabel=testLabelValue"]; ok {
			require.Equal(t,
				[]gometrics.Label{
					{
						Name:  "module",
						Value: "testModule",
					},
					{
						Name:  "testLabel",
						Value: "testLabelValue",
					},
				},
				metric.Labels)
			require.Equal(t, 3, metric.Count)
			// Since we can't inject time into gometrics we can't calculate the exact expected timing sample
			// so we bound the value between 0 and 3 seconds assuming that hardware isn't so overloaded that
			// executing this test takes longer than 3 seconds.
			require.Less(t, 0.0, metric.Sum)
			require.Greater(t, 3.0, metric.Sum)
			found = true
		}
	}
	require.True(t, found)
}
