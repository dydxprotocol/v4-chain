package server

import (
	"context"
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/bridge/api"
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

// ExpectBridgeDaemon registers the bridge daemon with the server. This is required
// in order to ensure that the daemon service is called at least once during every
// maximumAcceptableUpdateDelay duration. It will cause the protocol to panic if the daemon does not
// respond regularly.
func (server *Server) ExpectBridgeDaemon(maximumAcceptableUpdateDelay time.Duration) {
	server.registerDaemon(bridgeDaemonServiceName, maximumAcceptableUpdateDelay)
}

// AddBridgeEvents stores any bridge events recognized by the daemon
// in a go-routine safe slice.
func (s *Server) AddBridgeEvents(
	ctx context.Context,
	req *api.AddBridgeEventsRequest,
) (
	*api.AddBridgeEventsResponse,
	error,
) {
	_ = s.reportResponse(bridgeDaemonServiceName)
	if err := s.bridgeEventManager.AddBridgeEvents(req.BridgeEvents); err != nil {
		return nil, err
	}
	return &api.AddBridgeEventsResponse{}, nil
}
