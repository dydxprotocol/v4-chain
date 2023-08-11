package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPriceSmoothingPpm(t *testing.T) {
	require.Equal(t, PriceSmoothingPpm, uint32(300_000))
}
