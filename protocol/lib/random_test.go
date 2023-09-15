package lib_test

import (
	"bytes"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

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
		"no shared bytes": {
			start: []byte{1, 7, 255},
			end:   []byte{7, 8, 125},
		},
		"start is longer then end": {
			start: []byte{7, 7, 255},
			end:   []byte{7, 8},
		},
		"start is shorter then end": {
			start: []byte{7, 7},
			end:   []byte{7, 8, 255},
		},
		"both are the same length": {
			start: []byte{1, 2, 3},
			end:   []byte{3, 2, 1},
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
			result := lib.RandomBytesBetween(tc.start, tc.end, rand)
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
				lib.RandomBytesBetween(start, end, rand)
			})
		} else {
			result := lib.RandomBytesBetween(start, end, rand)
			require.LessOrEqual(t, start, result)
			require.GreaterOrEqual(t, end, result)
		}
	}
}

func TestRandomBytesBetween_InvalidInputs(t *testing.T) {
	require.Panics(t, func() {
		lib.RandomBytesBetween([]byte{}, []byte{}, nil)
	})
}
