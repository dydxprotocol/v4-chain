package server

import (
	"context"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types"
	bdtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
)

// sDAIServer defines the fields required for sDAI conversion rate updates.
type SDAIServer struct {
	sDAIEventManager *bdtypes.SDAIEventManager
}

// WithsDAIEventManager sets the `sDAIEventManager` field.
// This is updated by the sDAI service with a new conversion rate
func (server *Server) WithsDAIEventManager(
	sDAIEventManager *bdtypes.SDAIEventManager,
) *Server {
	server.sDAIEventManager = sDAIEventManager
	return server
}

// AddsDAIEvents stores any conversion rate recognized by the daemon
// in a go-routine safe slice.
func (s *Server) AddsDAIEvent(
	ctx context.Context,
	req *api.AddsDAIEventsRequest,
) (
	response *api.AddsDAIEventsResponse,
	err error,
) {
	if err = s.sDAIEventManager.AddsDAIEvent(req); err != nil {
		return nil, err
	}

	// Capture valid responses in metrics.
	s.reportValidResponse(types.SDAIDaemonServiceName)

	return &api.AddsDAIEventsResponse{}, nil
}
