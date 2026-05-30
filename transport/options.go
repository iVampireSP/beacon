package transport

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/iVampireSP/beacon/registry"
)

// Option configures an App. Mirrors Kratos's functional-option style.
type Option func(*options)

type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	servers          []Server

	beforeStart []func(context.Context) error
	beforeStop  []func(context.Context) error
	afterStart  []func(context.Context) error
	afterStop   []func(context.Context) error
}

// ID sets the unique instance ID. When unset, New generates a random one.
func ID(id string) Option { return func(o *options) { o.id = id } }

// Name sets the service name shared by all instances of the service.
func Name(name string) Option { return func(o *options) { o.name = name } }

// Version sets the running service version.
func Version(version string) Option { return func(o *options) { o.version = version } }

// Metadata sets arbitrary instance attributes recorded in the registry.
func Metadata(md map[string]string) Option { return func(o *options) { o.metadata = md } }

// Endpoint sets the advertised endpoints explicitly. When unset, the App
// collects endpoints from servers implementing Endpointer.
func Endpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

// Context sets the base context the App derives its lifecycle from.
func Context(ctx context.Context) Option { return func(o *options) { o.ctx = ctx } }

// Signal sets the OS signals that trigger graceful shutdown. When unset, New
// uses SIGTERM, SIGQUIT and SIGINT.
func Signal(sigs ...os.Signal) Option { return func(o *options) { o.sigs = sigs } }

// Registrar sets the service registry the instance registers with on start and
// deregisters from on stop. When unset, the App runs without registration.
func Registrar(r registry.Registrar) Option { return func(o *options) { o.registrar = r } }

// RegistrarTimeout bounds each (de)registration call.
func RegistrarTimeout(d time.Duration) Option {
	return func(o *options) { o.registrarTimeout = d }
}

// StopTimeout bounds how long each server is given to stop gracefully.
func StopTimeout(d time.Duration) Option { return func(o *options) { o.stopTimeout = d } }

// Servers adds servers for the App to run. May be called more than once.
func Servers(srv ...Server) Option {
	return func(o *options) { o.servers = append(o.servers, srv...) }
}

// BeforeStart registers a hook run before any server starts. A returned error
// aborts startup.
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) { o.beforeStart = append(o.beforeStart, fn) }
}

// AfterStart registers a hook run after all servers have started and the
// instance is registered. A returned error aborts startup.
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) { o.afterStart = append(o.afterStart, fn) }
}

// BeforeStop registers a hook run at the start of shutdown, before deregistering.
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) { o.beforeStop = append(o.beforeStop, fn) }
}

// AfterStop registers a hook run after all servers have stopped.
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) { o.afterStop = append(o.afterStop, fn) }
}
