package server

import "time"

const (
	DaemonStartupGracePeriod = 30 * time.Second

	liquidationsDaemonKey = "liquidations-daemon"
	pricefeedDaemonKey    = "pricefeed-daemon"
	bridgeDaemonKey       = "bridge-daemon"
)
