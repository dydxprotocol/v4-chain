package stoppable

import (
	"sync"
	"testing"
)

var (
	servicesRequiringCleanup = make(map[string][]Stoppable)
	lock                     sync.Mutex
)

// Stoppable is an interface for objects that can be stopped. A global map of these objects is maintained in the
// protocol/app package. This map is used to stop all running services that aren't cleaned up by the Network test
// object for our cli test suite.
type Stoppable interface {
	Stop()
}

// RegisterServiceForTestCleanup registers a service for cleanup.
func RegisterServiceForTestCleanup(serviceName string, service Stoppable) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := servicesRequiringCleanup[serviceName]; !ok {
		servicesRequiringCleanup[serviceName] = []Stoppable{}
	}
	servicesRequiringCleanup[serviceName] = append(servicesRequiringCleanup[serviceName], service)
}

// StopServices stops all services that were registered for cleanup.
func StopServices(t *testing.T, testUuid string) {
	lock.Lock()
	defer lock.Unlock()

	t.Log("Stopping services for test", "uuid", testUuid)
	if services, ok := servicesRequiringCleanup[testUuid]; ok {
		for _, service := range services {
			service.Stop()
		}
		delete(servicesRequiringCleanup, testUuid)
	}
}
