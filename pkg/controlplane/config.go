package controlplane

import (
	"sync"

	pb "github.com/SimonePesci/gomesh/api/proto"
)

// Routing configuration managed by the control plane
// Version is incremented with each update
type ConfigStore struct {
	mu sync.RWMutex // Protects concurrent access to the config store
	version int64 // Version number of the current config
	routes []*pb.Route // List of routing rules: use pointer to avoid copying the whole slice
}

// Create a new config store with default route
func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		// no need to initialize the mutex, it's zero-valued and ready to use
		version: 1,
		routes: []*pb.Route{
			// Default route is the test backend we have
			{
				Path: "/",
				Backend: "localhost:3000",
				AuthRequired: false,
				TimeoutMs: 5000,
			},
		},
	}
}

// Get the current config version (safe to call from multiple goroutines)
func (cs *ConfigStore) GetConfig() *pb.ConfigUpdate {

	// Read-only lock the config store
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	// Return the current config version and routes (a copy to avoid modifying the original slice)
	return &pb.ConfigUpdate{
		Version: cs.version,
		Routes: cs.routes,
	}
}

// Update the config store
func (cs * ConfigStore) UpdateConfig(routes []*pb.Route) *pb.ConfigUpdate {
	
	// Lock the config store
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Increment version number
	cs.version++

	// Update the config store
	cs.routes = routes

	// Return a new struct for simplicity (could have just returned void also)
	return &pb.ConfigUpdate{
		Version: cs.version,
		Routes: cs.routes,
	}
}

// Add a new route to the config store
func (cs *ConfigStore) AddRoute(route *pb.Route) *pb.ConfigUpdate {
	// Lock the config store
	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Increment version number
	cs.version++

	// Add the new route to the config store
	cs.routes = append(cs.routes, route)

	// Return a new struct for simplicity (could have just returned void also)
	return &pb.ConfigUpdate{
		Version: cs.version,
		Routes: cs.routes,
	}
}
