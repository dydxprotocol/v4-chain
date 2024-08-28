package server

import (
	"context"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/deleveraging/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types"
	deleveragingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/deleveraging"
)

// DeleveragingServer defines the fields required for deleveraging updates.
type DeleveragingServer struct {
	daemonDeleveragingInfo *deleveragingtypes.DaemonDeleveragingInfo
}

// WithDaemonDeleveragingInfo sets the `daemonDeleveragingInfo` field.
// This is updated by the deleveraging service with a list of subaccounts with open positions for each perp.
func (server *Server) WithDaemonDeleveragingInfo(
	daemonDeleveragingInfo *deleveragingtypes.DaemonDeleveragingInfo,
) *Server {
	server.daemonDeleveragingInfo = daemonDeleveragingInfo
	return server
}

// DeleverageSubaccounts stores the list of subaccount ids
// in a go-routine safe slice.
func (s *Server) DeleverageSubaccounts(
	ctx context.Context,
	req *api.DeleveragingSubaccountsRequest,
) (
	response *api.DeleveragingSubaccountsResponse,
	err error,
) {

	s.daemonDeleveragingInfo.UpdateSubaccountsWithPositions(req.SubaccountOpenPositionInfo)

	// Capture valid responses in metrics.
	s.reportValidResponse(types.DeleveragingDaemonServiceName)

	return &api.DeleveragingSubaccountsResponse{}, nil
}
