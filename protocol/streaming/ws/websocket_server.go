package ws

import (
	"context"
	"fmt"
	"math"
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

const (
	CLOB_PAIR_IDS_QUERY_PARAM = "clobPairIds"
	MARKET_IDS_QUERY_PARAM    = "marketIds"

	CLOSE_DEADLINE = 5 * time.Second
	WRITE_DEADLINE = 10 * time.Second
	PING_PERIOD    = 30 * time.Second
	PONG_DEADLINE  = 60 * time.Second
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

	writeDeadline time.Duration
	pingPeriod    time.Duration
	pongDeadline  time.Duration
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
		writeDeadline:    WRITE_DEADLINE,
		pingPeriod:       PING_PERIOD,
		pongDeadline:     PONG_DEADLINE,
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

	// Set ws max message size to 10 mb.
	conn.SetReadLimit(10 * 1024 * 1024)

	// Parse clobPairIds from query parameters
	clobPairIds, err := parseUint32(r, CLOB_PAIR_IDS_QUERY_PARAM)
	if err != nil {
		ws.logger.Error("Error parsing clobPairIds", "err", err)
		if err := sendCloseWithReason(conn, websocket.CloseUnsupportedData, err.Error()); err != nil {
			ws.logger.Error("Error sending close message", "err", err)
		}
		return
	}

	// Parse marketIds from query parameters
	marketIds, err := parseUint32(r, MARKET_IDS_QUERY_PARAM)
	if err != nil {
		ws.logger.Error("Error parsing marketIds", "err", err)
		if err := sendCloseWithReason(conn, websocket.CloseUnsupportedData, err.Error()); err != nil {
			ws.logger.Error("Error sending close message", "err", err)
		}
		return
	}

	// Parse subaccountIds from query parameters
	subaccountIds, err := parseSubaccountIds(r)
	if err != nil {
		ws.logger.Error("Error parsing subaccountIds", "err", err)
		if err := sendCloseWithReason(conn, websocket.CloseUnsupportedData, err.Error()); err != nil {
			ws.logger.Error("Error sending close message", "err", err)
		}
		return
	}

	// Parse filterOrdersBySubaccountId from query parameters
	filterOrdersBySubaccountId, err := parseFilterOrdersBySubaccountId(r)
	if err != nil {
		ws.logger.Error("Error parsing filterOrdersBySubaccountId", "err", err)
		if err := sendCloseWithReason(conn, websocket.CloseUnsupportedData, err.Error()); err != nil {
			ws.logger.Error("Error sending close message", "err", err)
		}
		return
	}

	websocketMessageSender := &WebsocketMessageSender{
		cdc:           ws.cdc,
		conn:          conn,
		writeDeadline: ws.writeDeadline,
	}

	done := make(chan struct{})
	defer close(done)
	go ws.readPongs(conn)
	go ws.sendPings(conn, done)

	ws.logger.Info(
		fmt.Sprintf("Received websocket streaming request for clob pair ids: %+v", clobPairIds),
	)

	err = ws.streamingManager.Subscribe(
		clobPairIds,
		subaccountIds,
		marketIds,
		filterOrdersBySubaccountId,
		websocketMessageSender,
	)
	if err != nil {
		ws.logger.Error(
			"Ending handler for websocket connection",
			"err", err,
		)
		if err := sendCloseWithReason(conn, websocket.CloseInternalServerErr, err.Error()); err != nil {
			ws.logger.Error("Error sending close message", "err", err)
		}
		return
	}
}

// readPongs closes the connection once the client stops answering pings. Clients are not expected
// to send anything, but gorilla only processes pongs and close frames while the connection is read.
func (ws *WebsocketServer) readPongs(conn *websocket.Conn) {
	defer conn.Close()
	extendReadDeadline := func(string) error {
		return conn.SetReadDeadline(time.Now().Add(ws.pongDeadline))
	}
	conn.SetPongHandler(extendReadDeadline)
	if err := extendReadDeadline(""); err != nil {
		return
	}
	for {
		if _, _, err := conn.NextReader(); err != nil {
			return
		}
	}
}

func (ws *WebsocketServer) sendPings(conn *websocket.Conn, done <-chan struct{}) {
	ticker := time.NewTicker(ws.pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			deadline := time.Now().Add(ws.writeDeadline)
			if err := conn.WriteControl(websocket.PingMessage, nil, deadline); err != nil {
				return
			}
		case <-done:
			return
		}
	}
}

func sendCloseWithReason(conn *websocket.Conn, closeCode int, reason string) error {
	closeMessage := websocket.FormatCloseMessage(closeCode, reason)
	// Set a write deadline to avoid blocking indefinitely
	if err := conn.SetWriteDeadline(time.Now().Add(CLOSE_DEADLINE)); err != nil {
		return err
	}
	return conn.WriteControl(
		websocket.CloseMessage,
		closeMessage,
		time.Now().Add(CLOSE_DEADLINE),
	)
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

		if number < 0 || number > math.MaxInt32 {
			return nil, fmt.Errorf("invalid subaccount number: %s", parts[1])
		}

		subaccountIds = append(subaccountIds, &satypes.SubaccountId{
			Owner:  parts[0],
			Number: uint32(number),
		})
	}

	return subaccountIds, nil
}

// parseFilterOrdersBySubaccountId is a helper function to parse the filterOrdersBySubaccountId flag
// from the query parameters.
func parseFilterOrdersBySubaccountId(r *http.Request) (bool, error) {
	token := r.URL.Query().Get("filterOrdersBySubaccountId")
	if token == "" {
		return false, nil
	}
	value, err := strconv.ParseBool(token)
	if err != nil {
		return false, fmt.Errorf("invalid filterOrdersBySubaccountId: %s", token)
	}
	return value, nil
}

// parseUint32 is a helper function to parse the uint32 from the query parameters.
func parseUint32(r *http.Request, queryParam string) ([]uint32, error) {
	param := r.URL.Query().Get(queryParam)
	if param == "" {
		return []uint32{}, nil
	}
	idStrs := strings.Split(param, ",")
	ids := make([]uint32, 0)
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid %s: %s", queryParam, idStr)
		}
		if id < 0 || id > math.MaxInt32 {
			return nil, fmt.Errorf("invalid %s: %s", queryParam, idStr)
		}
		ids = append(ids, uint32(id))
	}

	return ids, nil
}

func (ws *WebsocketServer) mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.Handler)
	return mux
}

// Start the websocket server in a separate goroutine.
func (ws *WebsocketServer) Start() {
	addr := fmt.Sprintf(":%d", ws.port)
	ws.server = &http.Server{Addr: addr, Handler: ws.mux()}
	go func() {
		ws.logger.Info("Starting websocket server on address " + addr)
		err := ws.server.ListenAndServe()
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
