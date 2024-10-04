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
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
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
	// Parse subaccountIds from query parameters
	subaccountIds, err := parseSubaccountIds(r)
	if err != nil {
		ws.logger.Error(
			"Error parsing subaccountIds",
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
		fmt.Sprintf("Received websocket streaming request for clob pair ids: %+v", clobPairIds),
	)

	err = ws.streamingManager.Subscribe(
		clobPairIds,
		subaccountIds,
		websocketMessageSender,
	)
	if err != nil {
		ws.logger.Error(
			"Ending handler for websocket connection",
			"err", err,
		)
		return
	}
}

// parseSubaccountIds is a helper function to parse the subaccountIds from the query parameters.
func parseSubaccountIds(r *http.Request) ([]*satypes.SubaccountId, error) {
	subaccountIdsParam := r.URL.Query().Get("subaccountIds")
	if subaccountIdsParam == "" {
		return []*satypes.SubaccountId{}, nil
	}
	idStrs := strings.Split(subaccountIdsParam, ",")
	subaccountIds := make([]*satypes.SubaccountId, 0)
	for _, idStr := range idStrs {
		parts := strings.Split(idStr, "/")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid subaccountId format: %s, expected subaccount_id format: owner/number", idStr)
		}

		number, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid subaccount number: %s, expected subaccount_id format: owner/number", parts[1])
		}

		subaccountIds = append(subaccountIds, &satypes.SubaccountId{
			Owner:  parts[0],
			Number: uint32(number),
		})
	}

	return subaccountIds, nil
}

// parseClobPairIds is a helper function to parse the clobPairIds from the query parameters.
func parseClobPairIds(r *http.Request) ([]uint32, error) {
	clobPairIdsParam := r.URL.Query().Get("clobPairIds")
	if clobPairIdsParam == "" {
		return []uint32{}, nil
	}
	idStrs := strings.Split(clobPairIdsParam, ",")
	clobPairIds := make([]uint32, 0)
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid clobPairId: %s", idStr)
		}
		clobPairIds = append(clobPairIds, uint32(id))
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
	err := ws.server.Shutdown(shutdownCtx)
	if err != nil {
		ws.logger.Error("Failed to shutdown websocket server", "err", err)
	}
}
