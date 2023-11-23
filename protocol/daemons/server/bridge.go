package server

import (
	"context"
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

// AddBridgeEvents stores any bridge events recognized by the daemon
// in a go-routine safe slice.
func (s *Server) AddBridgeEvents(
	ctx context.Context,
	req *api.AddBridgeEventsRequest,
) (
	*api.AddBridgeEventsResponse,
	error,
) {
	// If the daemon is unable to report a response, there is either an error in the registration of
	// this daemon, or another one. In either case, the protocol should panic.
	// TODO(CORE-582): Re-enable this check once the bridge daemon is fixed in local / CI environments.
	//if err := s.reportResponse(types.BridgeDaemonServiceName); err != nil {
	//	panic(err)
	//}
	if err := s.bridgeEventManager.AddBridgeEvents(req.BridgeEvents); err != nil {
		return nil, err
	}
	return &api.AddBridgeEventsResponse{}, nil
}
