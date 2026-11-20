package streaming_test

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	streaming "github.com/dydxprotocol/v4-chain/protocol/streaming"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/stretchr/testify/require"
)

type chanSender chan *clobtypes.StreamOrderbookUpdatesResponse

func (s chanSender) Send(response *clobtypes.StreamOrderbookUpdatesResponse) error {
	s <- response
	return nil
}

func TestAddOrderUpdatesToCache_FlushesWhenFull(t *testing.T) {
	const maxUpdatesInCache = 2
	manager := streaming.NewFullNodeStreamingManager(
		log.NewNopLogger(),
		uint32(time.Hour.Milliseconds()), // The ticker never fires during the test.
		maxUpdatesInCache,
		10,
		0,
		storetypes.NewTransientStoreKey("streaming"),
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
	)
	defer manager.Stop()

	sender := make(chanSender, 1)
	subscriptionEnded := make(chan error, 1)
	go func() {
		subscriptionEnded <- manager.Subscribe([]uint32{0}, nil, []uint32{0}, false, sender)
	}()
	require.Eventually(t, func() bool { return manager.TracksMarketId(0) }, time.Second, time.Millisecond)

	updates := make([]clobtypes.StreamUpdate, maxUpdatesInCache+1)
	manager.AddOrderUpdatesToCache(updates, make([]uint32, len(updates)))

	select {
	case response := <-sender:
		require.Len(t, response.Updates, len(updates))
	case err := <-subscriptionEnded:
		t.Fatalf("subscription ended: %v", err)
	case <-time.After(time.Second):
		t.Fatal("full cache was not flushed")
	}
	require.True(t, manager.TracksMarketId(0))
}
