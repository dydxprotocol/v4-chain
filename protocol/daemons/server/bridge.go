package server

import (
	"context"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	bdtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/bridge"
)

// BridgeServer defines the fields required for bridge event updates.
type BridgeServer struct {
	bridgeEventManager *bdtypes.BridgeEventManager
}

// WithBridgeEventManager sets the `bridgeEventManager` field.
// This is updated by the bridge service with a list of recognized bridge events.
func (server *Server) WithBridgeEventManager(
	bridgeEventManager *bdtypes.BridgeEventManager,
) *Server {
	server.bridgeEventManager = bridgeEventManager
	return server
}

// AddBridgeEvents stores any bridge events recognized by the daemon
// in a go-routine safe slice.
func (s *Server) AddBridgeEvents(
	ctx context.Context,
	req *api.AddBridgeEventsRequest,
) (
	response *api.AddBridgeEventsResponse,
	err error,
) {
	// Capture valid responses in metrics.
	defer func() {
		if err == nil {
			s.reportValidResponse(types.PricefeedDaemonServiceName)
		}
	}()

	s.reportValidResponse(types.BridgeDaemonServiceName)
	if err = s.bridgeEventManager.AddBridgeEvents(req.BridgeEvents); err != nil {
		return nil, err
	}
	return &api.AddBridgeEventsResponse{}, nil
}
