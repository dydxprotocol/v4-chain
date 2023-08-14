package sim_helpers_test

import (
	"errors"
	"math"
	"math/rand"
	"testing"

	testutil_rand "github.com/dydxprotocol/v4-chain/protocol/testutil/rand"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/sim_helpers"
	"github.com/stretchr/testify/require"
)

func TestGetRandomBucketValue_Invalid(t *testing.T) {
	tests := map[string]struct {
		buckets       []int
		expectedPanic error
	}{
		"empty buckets": {
			buckets:       []int{},
			expectedPanic: errors.New("you must supply at least 3 or more bucket boundary values"),
		},
		"invalid num of buckets": {
			buckets:       []int{1, 2},
			expectedPanic: errors.New("you must supply at least 3 or more bucket boundary values"),
		},
		"non-increasing buckets": {
			buckets:       []int{-99, 0, 0},
			expectedPanic: errors.New("bucket boundary values must be strictly increasing"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := rand.NewSource(1)
			r := rand.New(s)

			require.PanicsWithError(
				t,
				tc.expectedPanic.Error(),
				func() { sim_helpers.GetRandomBucketValue(r, tc.buckets) },
			)
		})
	}
}

func TestRandPositiveUint64(t *testing.T) {
	for i := 0; i < 1_000; i++ {
		r := testutil_rand.NewRand()

		n := sim_helpers.RandPositiveUint64(r)
		require.Positive(t, n)
	}
}

func TestRandWithWeight(t *testing.T) {
	const NUM_ITERATIONS = 10_000
	const EPSILON = .5 // .5% delta epsilon allowed for each side.

	tests := map[string]struct {
		optionsWithWeights map[string]int
		expectedPanic      error
	}{
		"Regular operation": {
			optionsWithWeights: map[string]int{
				"five":      5,  // 5%
				"ten":       10, // 10%
				"thirty":    30, // 30%
				"fiftyfive": 55, // 55%
			},
		},
		"Negative weight": {
			optionsWithWeights: map[string]int{
				"five":   5,
				"ten":    10,
				"negOne": -1,
			},
			expectedPanic: errors.New("RandWithWeight weights cannot be negative"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			s := rand.NewSource(1)
			r := rand.New(s)

			if tc.expectedPanic != nil {
				require.PanicsWithError(
					t,
					tc.expectedPanic.Error(),
					func() { sim_helpers.RandWithWeight(r, tc.optionsWithWeights) },
				)
				return
			}
			counts := map[string]int{}
			totalWeightSum := 0
			for option, weight := range tc.optionsWithWeights {
				counts[option] = 0
				totalWeightSum += weight
			}

			for i := 0; i < NUM_ITERATIONS; i++ {
				result := sim_helpers.RandWithWeight(r, tc.optionsWithWeights)
				counts[result] += 1
			}

			for option, weight := range tc.optionsWithWeights {
				expectedPercent := float64(weight) * 100 / float64(totalWeightSum)
				actualPercent := float64(counts[option]) * 100 / float64(NUM_ITERATIONS)
				require.InEpsilonf(t, expectedPercent, actualPercent, float64(EPSILON),
					"Failed to generate random element within epsilon %v weight", EPSILON)
			}
		})
	}
}

func TestMakeRange(t *testing.T) {
	tests := map[string]struct {
		n             uint32
		expectedSlice []uint32
	}{
		"n == 0": {
			n:             0,
			expectedSlice: []uint32{},
		},
		"n > 0": {
			n:             5,
			expectedSlice: []uint32{0, 1, 2, 3, 4},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.ElementsMatch(t, tc.expectedSlice, sim_helpers.MakeRange(tc.n))
		})
	}
}

func TestGetMin(t *testing.T) {
	tests := map[string]struct {
		x        int
		y        int
		expected int
	}{
		"x < y": {
			x:        -10,
			y:        10,
			expected: -10,
		},
		"x == y": {
			x:        9,
			y:        9,
			expected: 9,
		},
		"x > y": {
			x:        10,
			y:        -10,
			expected: -10,
		},
		"max int64 does NOT overflow": {
			x:        math.MaxInt64,
			y:        math.MaxInt64,
			expected: math.MaxInt64,
		},
		"min int64 does NOT underflow": {
			x:        math.MinInt64,
			y:        math.MinInt64,
			expected: math.MinInt64,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.expected, sim_helpers.GetMin(tc.x, tc.y))
		})
	}
}

func TestPickGenesisParameter(t *testing.T) {
	tests := map[string]struct {
		params       sim_helpers.GenesisParameters[int]
		isReasonable bool
		expected     int
	}{
		"isReasonable = true": {
			params: sim_helpers.GenesisParameters[int]{
				Reasonable: 0,
				Valid:      10,
			},
			isReasonable: true,
			expected:     0,
		},
		"isReasonable = false": {
			params: sim_helpers.GenesisParameters[int]{
				Reasonable: 0,
				Valid:      10,
			},
			isReasonable: false,
			expected:     10,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(
				t,
				tc.expected,
				sim_helpers.PickGenesisParameter(
					tc.params,
					tc.isReasonable,
				),
			)
		})
	}
}
