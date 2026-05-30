// Package http holds the Foundation HTTP kernel — the network entry point,
// mirroring Illuminate\Foundation\Http\Kernel. Where the console kernel runs CLI
// commands, the HTTP kernel runs the service's network (gRPC/HTTP) servers: it
// derives the service identity from config, wires tracing, and runs everything
// under the transport App lifecycle (signals, registry, graceful stop). It
// depends on the application through contracts.Application (not the concrete
// foundation.Application), so this package never imports foundation.
package http

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/tracing"
	"github.com/iVampireSP/beacon/transport"
	"github.com/iVampireSP/beacon/version"
)

// Kernel runs the service's transport servers — the network analog of the
// console kernel. Handle bootstraps (idempotently) then serves, mirroring
// Http\Kernel::handle's bootstrap-then-dispatch shape.
type Kernel struct {
	app  contracts.Application
	opts []transport.Option
}

// NewKernel creates an HTTP kernel bound to the application. Extra transport
// options (registry, custom signals, hooks) are applied on every Handle.
func NewKernel(app contracts.Application, opts ...transport.Option) *Kernel {
	return &Kernel{app: app, opts: opts}
}

// Handle boots the application (idempotent — the console kernel has normally
// booted it already) and runs the given servers under the transport App, using
// the service identity from config ("app.name") and the build version, until a
// shutdown signal. Shutdown of application resources is owned by the console
// kernel that invoked the serving command, so it is not repeated here.
func (k *Kernel) Handle(servers ...transport.Server) error {
	if err := k.app.Boot(); err != nil {
		return err
	}

	name := config.String("app.name", "beacon")

	tp, err := tracing.GetService(name)
	if err != nil {
		return err
	}
	defer tracing.ShutdownWithTimeout(tp)

	opts := append([]transport.Option{
		transport.Name(name),
		transport.Version(version.Version),
	}, k.opts...)
	opts = append(opts, transport.Servers(servers...))

	return transport.New(opts...).Run()
}
