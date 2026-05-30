// Package registry defines the service registration and discovery contracts.
//
// The transport App registers its running instance through a Registrar so a
// service registry (e.g. Consul) can advertise it, and consumers resolve peers
// through a Discovery — all without the runtime depending on any concrete
// registry implementation. Implementations live in their own packages and are
// wired in as a registry.Registrar option on the App.
package registry

import "context"

// Registrar registers and deregisters a service instance with a registry.
type Registrar interface {
	// Register adds the instance to the registry.
	Register(ctx context.Context, ins *ServiceInstance) error
	// Deregister removes the instance from the registry.
	Deregister(ctx context.Context, ins *ServiceInstance) error
}

// Discovery resolves and watches the instances of a service by name (the
// consumer side of a registry).
type Discovery interface {
	// GetService returns the current instances of the named service.
	GetService(ctx context.Context, name string) ([]*ServiceInstance, error)
	// Watch streams instance-list updates for the named service.
	Watch(ctx context.Context, name string) (Watcher, error)
}

// Watcher streams instance-list changes for a service until Stop is called.
type Watcher interface {
	// Next blocks until the instance list changes, then returns the new list.
	Next() ([]*ServiceInstance, error)
	// Stop ends the watch and releases its resources.
	Stop() error
}

// ServiceInstance is a single running instance of a service.
type ServiceInstance struct {
	// ID uniquely identifies this instance within the service.
	ID string `json:"id"`
	// Name is the service name shared by all instances of the service.
	Name string `json:"name"`
	// Version is the running version of the service.
	Version string `json:"version"`
	// Metadata carries arbitrary instance attributes (region, weight, ...).
	Metadata map[string]string `json:"metadata"`
	// Endpoints are the addresses this instance can be reached at, each a URL
	// of the form "scheme://host:port" (e.g. "grpc://127.0.0.1:9000").
	Endpoints []string `json:"endpoints"`
}
