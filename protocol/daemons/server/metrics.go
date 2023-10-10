package server

import (
	"time"

	"github.com/dydxprotocol/v4-chain/protocol/daemons/server/types"
)

// ExpectMetricsDaemon registers the periodic metrics daemon with the server. This is required
// in order to ensure that the daemon service is called at least once during every
// maximumAcceptableUpdateDelay duration. It will cause the protocol to panic if the daemon does not
// respond within maximumAcceptableUpdateDelay duration.
func (server *Server) ExpectMetricsDaemon(maximumAcceptableUpdateDelay time.Duration) {
	server.registerDaemon(types.MetricsDaemonServiceName, maximumAcceptableUpdateDelay)
}
