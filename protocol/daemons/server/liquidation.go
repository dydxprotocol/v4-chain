package server

import (
	"context"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/liquidation/api"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types"
	liquidationtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/liquidations"
)

// LiquidationServer defines the fields required for liquidation updates.
type LiquidationServer struct {
	daemonLiquidationInfo *liquidationtypes.DaemonLiquidationInfo
}

// WithDaemonLiquidationInfo sets the `daemonLiquidationInfo` field.
// This is updated by the liquidation service with a list of potentially liquidatable
// subaccount ids to be processed by the `PerpetualLiquidationsKeeper`.
func (server *Server) WithDaemonLiquidationInfo(
	daemonLiquidationInfo *liquidationtypes.DaemonLiquidationInfo,
) *Server {
	server.daemonLiquidationInfo = daemonLiquidationInfo
	return server
}

// LiquidateSubaccounts stores the list of potentially liquidatable subaccount ids
// in a go-routine safe slice.
func (s *Server) LiquidateSubaccounts(
	ctx context.Context,
	req *api.LiquidateSubaccountsRequest,
) (
	response *api.LiquidateSubaccountsResponse,
	err error,
) {

	s.daemonLiquidationInfo.UpdateSubaccountsWithPositions(req.SubaccountOpenPositionInfo)

	// Capture valid responses in metrics.
	s.reportValidResponse(types.LiquidationsDaemonServiceName)

	return &api.LiquidateSubaccountsResponse{}, nil
}
