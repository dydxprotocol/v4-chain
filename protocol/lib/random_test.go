package lib

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
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
