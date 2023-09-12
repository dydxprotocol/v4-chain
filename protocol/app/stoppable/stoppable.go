package stoppable

var (
	servicesRequiringCleanup = make(map[string]Stoppable)
)

// Stoppable is an interface for objects that can be stopped. A global map of these objects is maintained in the
// protocol/app package. This map is used to stop all running services that aren't cleaned up by the Network test
// object for our cli test suite.
type Stoppable interface {
	Stop()
}

// RegisterServiceForCleanup registers a service for cleanup.
func RegisterServiceForCleanup(serviceName string, service Stoppable) {
	servicesRequiringCleanup[serviceName] = service
}

// StopServices stops all services that were registered for cleanup.
func StopServices() {
	for _, service := range servicesRequiringCleanup {
		service.Stop()
	}
}
