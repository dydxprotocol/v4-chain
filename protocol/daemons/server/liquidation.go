package server

import (
	"context"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/liquidation/api"
	liquidationtypes "github.com/dydxprotocol/v4-chain/protocol/daemons/server/types/liquidations"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
)

// LiquidationServer defines the fields required for liquidation updates.
type LiquidationServer struct {
	liquidatableSubaccountIds *liquidationtypes.LiquidatableSubaccountIds
}

// WithLiquidatableSubaccountIds sets the `liquidatableSubaccountIds` field.
// This is updated by the liquidation service with a list of potentially liquidatable
// subaccount ids to be processed by the `PerpetualLiquidationsKeeper`.
func (server *Server) WithLiquidatableSubaccountIds(
	liquidatableSubaccountIds *liquidationtypes.LiquidatableSubaccountIds,
) *Server {
	server.liquidatableSubaccountIds = liquidatableSubaccountIds
	return server
}

// ExpectLiquidationsDaemon registers the liquidations daemon with the server. This is required
// in order to ensure that the daemon service is called at least once during every
// maximumAcceptableUpdateDelay duration. It will cause the protocol to panic if the daemon does not
// respond within maximumAcceptableUpdateDelay duration.
func (server *Server) ExpectLiquidationsDaemon(maximumAcceptableUpdateDelay time.Duration) {
	server.registerDaemon(types.LiquidationsDaemonServiceName, maximumAcceptableUpdateDelay)
}

// LiquidateSubaccounts stores the list of potentially liquidatable subaccount ids
// in a go-routine safe slice.
func (s *Server) LiquidateSubaccounts(
	ctx context.Context,
	req *api.LiquidateSubaccountsRequest,
) (*api.LiquidateSubaccountsResponse, error) {
	telemetry.ModuleSetGauge(
		metrics.LiquidationDaemon,
		float32(len(req.SubaccountIds)),
		metrics.LiquidatableSubaccountIds,
		metrics.Received,
		metrics.Count,
	)
	// If the daemon is unable to report a response, there is either an error in the registration of
	// this daemon, or another one. In either case, the protocol should panic.
	if err := s.reportResponse(types.LiquidationsDaemonServiceName); err != nil {
		s.logger.Error("Failed to report liquidations response to update monitor", "error", err)
	}

	s.liquidatableSubaccountIds.UpdateSubaccountIds(req.SubaccountIds)
	return &api.LiquidateSubaccountsResponse{}, nil
}
