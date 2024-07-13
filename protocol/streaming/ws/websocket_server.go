package ws

import (
	"net/http"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default
	},
}

type WebsocketServer struct {
	streamingManager types.FullNodeStreamingManager
	cdc              codec.JSONCodec
	logger           log.Logger
}

func NewWebsocketServer(
	streamingManager types.FullNodeStreamingManager,
	cdc codec.JSONCodec,
	logger log.Logger,
) *WebsocketServer {
	return &WebsocketServer{
		streamingManager: streamingManager,
		cdc:              cdc,
		logger:           logger,
	}
}

func (ws *WebsocketServer) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.logger.Error(
			"Error upgrading websocket connection",
			"error", err,
		)
		return
	}
	defer conn.Close()

	websocketMessageSender := &WebsocketMessageSender{
		cdc:  ws.cdc,
		conn: conn,
	}

	// TODO: remove this
	conn.WriteMessage(websocket.TextMessage, []byte("Connected to server..."))

	err = ws.streamingManager.Subscribe(
		[]uint32{0, 1}, // TODO: Get clobPairIds from request
		websocketMessageSender,
	)
	if err != nil {
		ws.logger.Error(
			"Error subscribing to stream",
			"error", err,
		)
		return
	}
}

// Start the websocket server in a separate goroutine.
func (ws *WebsocketServer) Start() {
	// TODO: use app flags to control the port.
	go func() {
		http.HandleFunc("/ws", ws.Handler)
		err := http.ListenAndServe(":4321", nil)
		if err != nil {
			ws.logger.Error(
				"Error starting websocket server",
				"error", err,
			)
		}
	}()
}
