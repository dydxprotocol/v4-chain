package ws

import (
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

type subscribeFunc func(sender types.OutgoingMessageSender) error

type fakeStreamingManager struct {
	types.FullNodeStreamingManager
	subscribe subscribeFunc
}

func (m fakeStreamingManager) Subscribe(
	_ []uint32,
	_ []*satypes.SubaccountId,
	_ []uint32,
	_ bool,
	sender types.OutgoingMessageSender,
) error {
	return m.subscribe(sender)
}

// connect serves a WebsocketServer whose subscription loop is subscribe and connects a client to
// it. The returned channel receives what subscribe returned.
func connect(
	t *testing.T,
	configure func(*WebsocketServer),
	subscribe subscribeFunc,
) (*websocket.Conn, <-chan error) {
	subscribeResult := make(chan error, 1)
	server := NewWebsocketServer(
		fakeStreamingManager{subscribe: func(sender types.OutgoingMessageSender) error {
			err := subscribe(sender)
			subscribeResult <- err
			return err
		}},
		codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
		log.NewNopLogger(),
		0,
	)
	configure(server)
	httpServer := httptest.NewServer(http.HandlerFunc(server.Handler))
	t.Cleanup(httpServer.Close)

	url := "ws" + strings.TrimPrefix(httpServer.URL, "http") + "?clobPairIds=0"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return conn, subscribeResult
}

func TestSend_TimesOutWhenClientStopsReading(t *testing.T) {
	_, subscribeResult := connect(
		t,
		func(server *WebsocketServer) { server.writeDeadline = 100 * time.Millisecond },
		func(sender types.OutgoingMessageSender) error {
			response := &clobtypes.StreamOrderbookUpdatesResponse{
				Updates: make([]clobtypes.StreamUpdate, 20_000),
			}
			for {
				if err := sender.Send(response); err != nil {
					return err
				}
			}
		},
	)

	select {
	case err := <-subscribeResult:
		var netErr net.Error
		require.ErrorAs(t, err, &netErr)
		require.True(t, netErr.Timeout())
	case <-time.After(10 * time.Second):
		t.Fatal("send to a client that stopped reading did not time out")
	}
}

func TestHandler_KeepsClientThatAnswersPingsConnected(t *testing.T) {
	const pongDeadline = 500 * time.Millisecond
	unsubscribe := make(chan struct{})
	defer close(unsubscribe)
	conn, _ := connect(
		t,
		func(server *WebsocketServer) {
			server.pingPeriod = 20 * time.Millisecond
			server.pongDeadline = pongDeadline
		},
		func(types.OutgoingMessageSender) error {
			<-unsubscribe
			return nil
		},
	)

	var pings atomic.Int32
	answerPing := conn.PingHandler()
	conn.SetPingHandler(func(data string) error {
		pings.Add(1)
		return answerPing(data)
	})
	readErr := make(chan error, 1)
	go func() {
		_, _, err := conn.ReadMessage()
		readErr <- err
	}()

	select {
	case err := <-readErr:
		t.Fatalf("connection closed: %v", err)
	case <-time.After(2 * pongDeadline):
	}
	require.Greater(t, pings.Load(), int32(1))
}

func TestHandler_ClosesConnectionWhenPongsStop(t *testing.T) {
	_, subscribeResult := connect(
		t,
		func(server *WebsocketServer) {
			server.pingPeriod = 20 * time.Millisecond
			server.pongDeadline = 100 * time.Millisecond
		},
		func(sender types.OutgoingMessageSender) error {
			for {
				if err := sender.Send(&clobtypes.StreamOrderbookUpdatesResponse{}); err != nil {
					return err
				}
				time.Sleep(10 * time.Millisecond)
			}
		},
	)

	// The client never reads, so it never answers a ping.
	select {
	case err := <-subscribeResult:
		require.ErrorIs(t, err, net.ErrClosed)
	case <-time.After(10 * time.Second):
		t.Fatal("connection to a client that does not answer pings was not closed")
	}
}
