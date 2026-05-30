// Package foundation is the integration layer that ties the container, service
// providers, and console together into a runnable application — the analog of
// Illuminate\Foundation. foundation.Application mirrors
// Illuminate\Foundation\Application (which extends Illuminate\Container\Container);
// here it embeds *container.Container for the same effect.
package foundation

import (
	"errors"

	"github.com/iVampireSP/beacon/container"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// Application is the Foundation-layer application: the DI container plus the
// service-provider lifecycle plus the console entry point. Build one with
// Configure(...).Create(); it satisfies contracts.Application and, by forwarding
// RegisterCommands to its console kernel, support.Kernel — so a provider can push
// commands with just the application in hand.
type Application struct {
	*container.Container
	providers        []support.Provider
	bootingCallbacks []func(*Application)
	bootedCallbacks  []func(*Application)
	shutdownFns      []func() error
	consoleKernel    contracts.ConsoleKernel
	booted           bool
}

// newApplication creates an application over a fresh container and registers the
// application into its own container, so commands and services can inject
// *foundation.Application (mirroring Laravel binding the app as a singleton).
func newApplication() *Application {
	app := &Application{Container: container.NewContainer()}
	_ = app.Singleton(func() *Application { return app })
	return app
}

// Register hands service providers to the application. Each provider's Register()
// is invoked immediately so it can bind dependencies into the container; the
// provider is then retained for Boot. Mirrors $app->register(new Provider($app)).
func (app *Application) Register(providers ...support.Provider) {
	for _, p := range providers {
		p.Register()
		app.providers = append(app.providers, p)
	}
}

// Boot fires the booting callbacks, calls Boot() on all registered providers
// (once), then fires the booted callbacks — mirroring Laravel's boot().
func (app *Application) Boot() error {
	if app.booted {
		return nil
	}
	app.fireAppCallbacks(app.bootingCallbacks)
	for _, p := range app.providers {
		p.Boot()
	}
	app.booted = true
	app.fireAppCallbacks(app.bootedCallbacks)
	return nil
}

// Booting registers a callback to run just before the providers are booted.
func (app *Application) Booting(callback func(*Application)) {
	app.bootingCallbacks = append(app.bootingCallbacks, callback)
}

// Booted registers a callback to run just after all providers have booted. If
// the application is already booted, the callback runs immediately — matching
// Laravel's $app->booted().
func (app *Application) Booted(callback func(*Application)) {
	app.bootedCallbacks = append(app.bootedCallbacks, callback)
	if app.booted {
		callback(app)
	}
}

func (app *Application) fireAppCallbacks(callbacks []func(*Application)) {
	for _, cb := range callbacks {
		cb(app)
	}
}

// Providers returns the registered service providers in registration order.
func (app *Application) Providers() []support.Provider {
	return app.providers
}

// RegisterCommands forwards command constructors to the console kernel, so
// *Application satisfies support.Kernel and a provider's AddCommand reaches the
// kernel through the application.
func (app *Application) RegisterCommands(constructors ...any) {
	app.consoleKernel.RegisterCommands(constructors...)
}

// HandleCommand bootstraps and runs the CLI under the given root identity,
// returning the process exit code. Mirrors Laravel's $app->handleCommand().
func (app *Application) HandleCommand(use, short string) int {
	if err := app.consoleKernel.Handle(use, short); err != nil {
		return 1
	}
	return 0
}

// OnShutdown registers a cleanup callback, run during Shutdown in reverse order.
func (app *Application) OnShutdown(fn func() error) {
	app.shutdownFns = append(app.shutdownFns, fn)
}

// Shutdown executes all registered cleanup callbacks in reverse order.
func (app *Application) Shutdown() error {
	var errs []error
	for i := len(app.shutdownFns) - 1; i >= 0; i-- {
		if err := app.shutdownFns[i](); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// ProvidersImplementing returns all registered providers that satisfy capability
// interface T (e.g. job.HandlerProvider, bus.ListenerProvider,
// schedule.CronProvider), letting a subsystem gather contributions from every
// provider without a central registration list.
func ProvidersImplementing[T any](app *Application) []T {
	var out []T
	for _, p := range app.providers {
		if t, ok := any(p).(T); ok {
			out = append(out, t)
		}
	}
	return out
}
