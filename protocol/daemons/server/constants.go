package server

import "time"

const (
	// DaemonStartupGracePeriod is the amount of time to wait for before a daemon is expected to start querying
	// the daemon server. This is used to ensure that spurious panics aren't produced due to the daemons waiting for
	// the cosmos grpc service to start. If cli tests are failing due to panics because it is taking the network
	// a long time to start the protocol, it's possible this value could be increased.
	DaemonStartupGracePeriod = 30 * time.Second

	liquidationsDaemonServiceName = "liquidations-daemon"
	pricefeedDaemonServiceName    = "pricefeed-daemon"
	bridgeDaemonServiceName       = "bridge-daemon"
)
