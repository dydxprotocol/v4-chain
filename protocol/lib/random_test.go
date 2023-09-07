package lib

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	exprand "golang.org/x/exp/rand"
)

func TestRandomBool(t *testing.T) {
	numIterations := 128
	bools := make([]bool, numIterations)

	for i := 0; i < numIterations; i++ {
		bools[i] = RandomBool()
	}

	require.Contains(t, bools, true)
	require.Contains(t, bools, false)
}

func TestRandomBytesBetween(t *testing.T) {
	tests := map[string]struct {
		start []byte
		end   []byte
	}{
		"start equals end": {
			start: []byte{7, 7},
			end:   []byte{7, 7},
		},
		"start is a prefix of end": {
			start: []byte{7, 7},
			end:   []byte{7, 7, 0},
		},
		"start is longer then end": {
			start: []byte{7, 7, 255},
			end:   []byte{7, 8},
		},
		"both are empty": {
			start: []byte{},
			end:   []byte{},
		},
		"start is empty": {
			start: []byte{},
			end:   []byte{0},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			rand := rand.New(rand.NewSource(53))
			result := RandomBytesBetween(tc.start, tc.end, rand)
			require.LessOrEqual(t, tc.start, result)
			require.GreaterOrEqual(t, tc.end, result)
		})
	}
}

func TestRandomBytesBetween_RandomlyGenerated(t *testing.T) {
	rand := rand.New(rand.NewSource(53))
	for i := 0; i < 100_000; i++ {
		start := make([]byte, rand.Intn(3))
		end := make([]byte, rand.Intn(3))
		rand.Read(start)
		rand.Read(end)
		if bytes.Compare(start, end) > 0 {
			require.Panics(t, func() {
				RandomBytesBetween(start, end, rand)
			})
		} else {
			result := RandomBytesBetween(start, end, rand)
			require.LessOrEqual(t, start, result)
			require.GreaterOrEqual(t, end, result)
		}
	}
}

func TestRandomBytesBetween_InvalidInputs(t *testing.T) {
	require.Panics(t, func() {
		RandomBytesBetween([]byte{}, []byte{}, nil)
	})
}

func TestWeightedRandomSample(t *testing.T) {
	source := exprand.NewSource(53)

	freq := make([]float64, 4)
	for i := 0; i < 1_000_000; i++ {
		weights := []float64{0.1, 0.2, 0.3, 0.4}
		sample, err := WeightedRandomSample(weights, 1, source)
		require.NoError(t, err)
		freq[sample[0]]++
	}

	exp := []float64{100_000, 200_000, 300_000, 400_000}

	// Check that this is within statistical expectations.
	require.Less(
		t,
		chi2(freq, exp),
		16.92, // p = 0.05 df = 9
	)
}

func TestWeightedRandomSample_NoReplacement(t *testing.T) {
	source := exprand.NewSource(53)

	for i := 0; i < 10_000; i++ {
		weights := []float64{0.1, 0.2, 0.3, 0.4}
		sample, err := WeightedRandomSample(weights, 4, source)
		require.NoError(t, err)
		require.ElementsMatch(t, sample, []int{0, 1, 2, 3})
	}
}

func TestWeightedRandomSample_ZeroWeight(t *testing.T) {
	source := exprand.NewSource(53)

	for i := 0; i < 10_000; i++ {
		weights := []float64{0.1, 0.2, 0.3, 0.0, 0.4}
		// Pick 4 elements.
		sample, err := WeightedRandomSample(weights, 4, source)
		require.NoError(t, err)
		require.ElementsMatch(t, sample, []int{0, 1, 2, 4})
	}
}

func TestWeightedRandomSample_ErrorNoRemainingElements(t *testing.T) {
	source := exprand.NewSource(53)
	weights := []float64{0.1, 0.2, 0.3, 0.0, 0.4}
	// Pick 5 elements.
	_, err := WeightedRandomSample(weights, 5, source)
	require.ErrorContains(t, err, "failed to take item from weighted")

	source = exprand.NewSource(53)
	weights = []float64{0.1, 0.2, 0.3, 0.0, 0.4}
	// Pick 6 elements.
	_, err = WeightedRandomSample(weights, 6, source)
	require.ErrorContains(t, err, "failed to take item from weighted")
}

func chi2(ob, ex []float64) (sum float64) {
	for i := range ob {
		x := ob[i] - ex[i]
		sum += (x * x) / ex[i]
	}

	return sum
}
