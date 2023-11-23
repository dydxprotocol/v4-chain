package server

import (
	"context"
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
	s.liquidatableSubaccountIds.UpdateSubaccountIds(req.SubaccountIds)
	return &api.LiquidateSubaccountsResponse{}, nil
}
