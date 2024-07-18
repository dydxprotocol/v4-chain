package ws

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	port             uint16
	server           *http.Server
}

func NewWebsocketServer(
	streamingManager types.FullNodeStreamingManager,
	cdc codec.JSONCodec,
	logger log.Logger,
	port uint16,
) *WebsocketServer {
	return &WebsocketServer{
		streamingManager: streamingManager,
		cdc:              cdc,
		logger:           logger.With(log.ModuleKey, "full-node-streaming"),
		port:             port,
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

	// Parse clobPairIds from query parameters
	clobPairIds, err := parseClobPairIds(r)
	if err != nil {
		ws.logger.Error(
			"Error parsing clobPairIds",
			"err", err,
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	websocketMessageSender := &WebsocketMessageSender{
		cdc:  ws.cdc,
		conn: conn,
	}

	ws.logger.Info(
		fmt.Sprintf("Recieved websocket streaming request for clob pair ids: %+v", clobPairIds),
	)

	err = ws.streamingManager.Subscribe(
		clobPairIds,
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

// parseClobPairIds is a helper function to parse the clobPairIds from the query parameters.
func parseClobPairIds(r *http.Request) ([]uint32, error) {
	clobPairIdsParam := r.URL.Query().Get("clobPairIds")
	if clobPairIdsParam == "" {
		return nil, fmt.Errorf("missing clobPairIds parameter")
	}

	idStrs := strings.Split(clobPairIdsParam, ",")
	clobPairIds := make([]uint32, len(idStrs))
	for i, idStr := range idStrs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid clobPairId: %s", idStr)
		}
		clobPairIds[i] = uint32(id)
	}

	return clobPairIds, nil
}

// Start the websocket server in a separate goroutine.
func (ws *WebsocketServer) Start() {
	go func() {
		http.HandleFunc("/ws", ws.Handler)
		addr := fmt.Sprintf(":%d", ws.port)
		ws.logger.Info("Starting websocket server on address " + addr)

		server := &http.Server{Addr: addr}
		ws.server = server
		err := server.ListenAndServe()
		if err != nil {
			ws.logger.Error(
				"Http websocket server error",
				"err", err,
			)
		}
		ws.logger.Info("Shutting down websocket server")
	}()
}

func (ws *WebsocketServer) Shutdown() {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownRelease()
	ws.server.Shutdown(shutdownCtx)
}
