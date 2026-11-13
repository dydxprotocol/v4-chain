package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	"github.com/stretchr/testify/require"
)

func TestPendingKeyRoundTripAndRouteIsolation(t *testing.T) {
	key := types.PendingKey(types.Route_ROUTE_INJECTIVE, "channel-0", 7)
	route, channel, sequence, err := types.ParsePendingKey(key)
	require.NoError(t, err)
	require.Equal(t, types.Route_ROUTE_INJECTIVE, route)
	require.Equal(t, "channel-0", channel)
	require.Equal(t, uint64(7), sequence)

	other := types.PendingKey(types.Route_ROUTE_NOBLE, "channel-0", 7)
	require.NotEqual(t, key, other)
}

func TestParsePendingKeyRejectsMalformed(t *testing.T) {
	_, _, _, err := types.ParsePendingKey([]byte{1, 2, 3})
	require.Error(t, err)
}
