// Package transport provides the service runtime: a Kratos-style App that runs
// one or more servers under a single service identity, registers the instance
// with a service registry, and shuts everything down gracefully on an OS
// signal. It mirrors the lifecycle of Kratos's kratos.App.
package transport

import (
	"context"
	"net/url"
)

// Server is a long-running transport server managed by the App. Start should
// block until the server stops (returning nil on a graceful stop); Stop should
// drain it within the given context's deadline. The App runs every server's
// Start concurrently and calls Stop when the App is shutting down.
type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Endpointer is an optional Server capability: it reports the advertised
// endpoint the App records in the registered ServiceInstance, so consumers
// discover where to reach the server. Servers that don't implement it simply
// contribute no endpoint.
type Endpointer interface {
	Endpoint() (*url.URL, error)
}
