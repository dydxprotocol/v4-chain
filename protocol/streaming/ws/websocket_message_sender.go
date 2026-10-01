package ws

import (
	"context"

	"github.com/gorilla/websocket"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.OutgoingMessageSender = (*WebsocketMessageSender)(nil)

type WebsocketMessageSender struct {
	cdc codec.JSONCodec

	conn *websocket.Conn

	// ctx is cancelled once the connection's read loop detects that the
	// client has disconnected.
	ctx context.Context
}

func (wms *WebsocketMessageSender) Send(
	response *clobtypes.StreamOrderbookUpdatesResponse,
) (err error) {
	responseJson, err := wms.cdc.MarshalJSON(response)
	if err != nil {
		return err
	}
	return wms.conn.WriteMessage(websocket.TextMessage, responseJson)
}

func (wms *WebsocketMessageSender) Context() context.Context {
	return wms.ctx
}
