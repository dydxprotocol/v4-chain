package events_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/stretchr/testify/require"
)

func TestConstants(t *testing.T) {
	require.Equal(t, "order_fill", events.SubtypeOrderFill)
	require.Equal(t, "subaccount_update", events.SubtypeSubaccountUpdate)
	require.Equal(t, "transfer", events.SubtypeTransfer)
	require.Equal(t, "market", events.SubtypeMarket)
	require.Equal(t, "funding_values", events.SubtypeFundingValues)
}
