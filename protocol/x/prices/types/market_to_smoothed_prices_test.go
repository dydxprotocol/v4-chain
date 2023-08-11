package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMarketToSmoothedPrices_IsEmpty(t *testing.T) {
	mtsp := NewMarketToSmoothedPrices()
	require.Empty(t, mtsp)
}
