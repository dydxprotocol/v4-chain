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
	if err = s.bridgeEventManager.AddBridgeEvents(req.BridgeEvents); err != nil {
		return nil, err
	}

	// Capture valid responses in metrics.
	s.reportValidResponse(types.BridgeDaemonServiceName)

	return &api.AddBridgeEventsResponse{}, nil
}
