package foundation

import (
	fconsole "github.com/iVampireSP/beacon/foundation/console"
	"github.com/iVampireSP/beacon/support"
)

// Builder configures and creates an Application — the analog of Laravel's
// Illuminate\Foundation\Configuration\ApplicationBuilder, used as
// foundation.Configure(...).WithProviders(...).WithCommands(...).Create(),
// mirroring Application::configure(...)->with...()->create().
type Builder struct {
	providerFactory func(*Application) []support.Provider
	commandCtors    []any
}

// Configure starts building an application. The CLI root identity (use, short)
// is supplied later to HandleCommand, the way Laravel's artisan entry passes the
// input to handleCommand.
func Configure() *Builder {
	return &Builder{}
}

// WithProviders sets the factory that builds the service providers. The factory
// receives the application so each provider can take it (the analog of Laravel's
// new Provider($app)). Returns the builder for chaining.
func (b *Builder) WithProviders(factory func(*Application) []support.Provider) *Builder {
	b.providerFactory = factory
	return b
}

// WithCommands registers console-command constructors directly on the kernel, in
// addition to any a provider pushes via AddCommand. Returns the builder.
func (b *Builder) WithCommands(constructors ...any) *Builder {
	b.commandCtors = append(b.commandCtors, constructors...)
	return b
}

// Create assembles the application: it wires the console kernel, registers the
// builder's commands, then registers the providers (which may push their own
// commands). Mirrors ApplicationBuilder::create().
func (b *Builder) Create() *Application {
	app := newApplication()
	app.consoleKernel = fconsole.NewKernel(app, app.Container)
	app.RegisterCommands(b.commandCtors...)
	if b.providerFactory != nil {
		app.Register(b.providerFactory(app)...)
	}
	return app
}
