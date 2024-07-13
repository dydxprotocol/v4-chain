package ws

import (
	"github.com/gorilla/websocket"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var _ types.OutgoingMessageSender = (*WebsocketMessageSender)(nil)

type WebsocketMessageSender struct {
	cdc codec.JSONCodec

	conn *websocket.Conn
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
