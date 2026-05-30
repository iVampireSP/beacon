// Package transport provides the service runtime: a Kratos-style App that runs
// one or more servers under a single service identity and shuts everything down
// gracefully on an OS signal. It mirrors the run half of Kratos's kratos.App.
//
// Service registration/discovery is intentionally out of scope: the mesh (Istio)
// / Kubernetes does discovery + LB + mTLS in production, and peer addresses come
// from config locally — so there is no registry here.
package transport

import "context"

// Server is a long-running transport server managed by the App. Start should
// block until the server stops (returning nil on a graceful stop); Stop should
// drain it within the given context's deadline. The App runs every server's
// Start concurrently and calls Stop when the App is shutting down.
type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
