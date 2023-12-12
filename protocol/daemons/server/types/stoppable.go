package types

// Stoppable is an interface for a service that can be stopped.
// This is used to stop services registered with the health monitor.
type Stoppable interface {
	Stop()
}
