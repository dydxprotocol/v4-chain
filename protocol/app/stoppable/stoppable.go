package stoppable

import (
	"sync"
	"testing"
)

var (
	servicesRequiringCleanup = make(map[string][]Stoppable)
	lock                     sync.Mutex
)

// Stoppable is an interface for objects that can be stopped. A global map of these objects is maintained here. This
// map is used to stop all running services that aren't cleaned up by the Network test object for our cli test suite.
// Services are organized by a uuid per test case, which is that test's GRPC address, since the network package chooses
// these to not overlap, and these are easily accessible to the protocol from the app.New method, where many services
// are started, and which does not have a reference to an sdk context.
type Stoppable interface {
	Stop()
}

// RegisterServiceForTestCleanup registers a service for cleanup. All services are organized by a uuid per test case.
func RegisterServiceForTestCleanup(testUuid string, service Stoppable) {
	lock.Lock()
	defer lock.Unlock()

	if _, ok := servicesRequiringCleanup[testUuid]; !ok {
		servicesRequiringCleanup[testUuid] = []Stoppable{}
	}
	servicesRequiringCleanup[testUuid] = append(servicesRequiringCleanup[testUuid], service)
}

// StopServices stops all services that were registered for cleanup for a given test, identified by uuid.
// It also removes the services from the global map.
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
