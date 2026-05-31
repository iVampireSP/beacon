package transport

import (
	"context"
	"os"
	"time"
)

// Option configures an App. Mirrors Kratos's functional-option style.
type Option func(*options)

type options struct {
	id      string
	name    string
	version string

	ctx  context.Context
	sigs []os.Signal

	stopTimeout time.Duration
	servers     []TransportServer

	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// ID sets the unique instance ID. When unset, New generates a random one.
func ID(id string) Option { return func(o *options) { o.id = id } }

// Name sets the service name (used in startup/shutdown logs).
func Name(name string) Option { return func(o *options) { o.name = name } }

// Version sets the running service version.
func Version(version string) Option { return func(o *options) { o.version = version } }

// Context sets the base context the App derives its lifecycle from.
func Context(ctx context.Context) Option { return func(o *options) { o.ctx = ctx } }

// Signal sets the OS signals that trigger graceful shutdown. When unset, New
// uses SIGTERM, SIGQUIT and SIGINT.
func Signal(sigs ...os.Signal) Option { return func(o *options) { o.sigs = sigs } }

// StopTimeout bounds how long each server is given to stop gracefully.
func StopTimeout(d time.Duration) Option { return func(o *options) { o.stopTimeout = d } }

// Servers adds servers for the App to run. May be called more than once.
func Servers(srv ...TransportServer) Option {
	return func(o *options) { o.servers = append(o.servers, srv...) }
}

// BeforeStart registers a hook run before any server starts. A returned error
// aborts startup.
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) { o.beforeStart = append(o.beforeStart, fn) }
}

// AfterStart registers a hook run after all servers have started. A returned
// error aborts startup.
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) { o.afterStart = append(o.afterStart, fn) }
}

// BeforeStop registers a hook run at the start of shutdown.
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) { o.beforeStop = append(o.beforeStop, fn) }
}

// AfterStop registers a hook run after all servers have stopped.
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) { o.afterStop = append(o.afterStop, fn) }
}
