package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/stretchr/testify/require"
)

func TestMarketKey(t *testing.T) {
	result := types.MarketKey(uint32(2))
	require.Equal(t, "\x02\x00\x00\x00/", string(result))
}
