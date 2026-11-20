package streaming_test

import (
	"errors"
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

type failingSender struct{}

func (failingSender) Send(*clobtypes.StreamOrderbookUpdatesResponse) error {
	return errors.New("client is gone")
}

func TestSubscribe_RemovesSubscriptionWhenSendFails(t *testing.T) {
	manager := streaming.NewFullNodeStreamingManager(
		log.NewNopLogger(),
		uint32(time.Hour.Milliseconds()), // The ticker never fires during the test.
		10,
		10,
		0,
		storetypes.NewTransientStoreKey("streaming"),
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
	)
	defer manager.Stop()

	subscriptionEnded := make(chan error, 1)
	go func() {
		subscriptionEnded <- manager.Subscribe([]uint32{0}, nil, []uint32{0}, false, failingSender{})
	}()
	require.Eventually(t, func() bool { return manager.TracksMarketId(0) }, time.Second, time.Millisecond)

	manager.AddOrderUpdatesToCache(make([]clobtypes.StreamUpdate, 1), []uint32{0})
	manager.FlushStreamUpdates()

	select {
	case err := <-subscriptionEnded:
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("subscription did not end")
	}
	require.False(t, manager.TracksMarketId(0))
}
