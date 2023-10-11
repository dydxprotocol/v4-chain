package types

type DaemonClient interface {
	HealthCheckable

	// TODO(CORE-29): gracefully shut down daemons (uncomment following line and implement)
	// Stoppable
}
