package types

import "context"

// HealthCheckable is a common interface for services that can be health checked.
type HealthCheckable interface {
	// HealthCheck returns an error if a service is unhealthy. This method takes
	// a context to allow the caller to use a context-based timeout with the call.
	// If the service is healthy, this method returns nil.
	HealthCheck(ctx context.Context) (err error)
}
