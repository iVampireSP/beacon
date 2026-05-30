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
	"github.com/iVampireSP/beacon/foundation/bootstrap"
	"github.com/iVampireSP/beacon/logger"
	"github.com/iVampireSP/beacon/support"
)

// Application satisfies the decoupled contracts.Application surface that
// providers and bootstrappers depend on.
var _ contracts.Application = (*Application)(nil)

// Application is the Foundation-layer application: the DI container plus the
// service-provider lifecycle plus the console entry point. Build one with
// Configure(...).Create(); it satisfies contracts.Application and, by forwarding
// RegisterCommands to its console kernel, support.Kernel — so a provider can push
// commands with just the application in hand.
type Application struct {
	*container.Container
	providers        []support.Provider
	bootstrappers    []bootstrap.Bootstrapper
	contributions    []any
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
	// Register the app under both its concrete type and the contracts.Application
	// interface, so a command constructor can inject either (the beacon runtimes
	// — worker/eventbus/scheduler — take the interface to stay off foundation).
	_ = app.Singleton(func() *Application { return app })
	_ = app.Singleton(func() contracts.Application { return app })
	return app
}

// Singleton registers a constructor whose result is built once and cached. It
// narrows the embedded *container.Container's Singleton to the signature in
// contracts.Application, so the contract need not leak the DI library's options.
func (app *Application) Singleton(constructor any) error {
	return app.Container.Singleton(constructor)
}

// Invoke resolves fn's parameters from the container and calls it, narrowing the
// embedded Container's Invoke to the contracts.Application signature.
func (app *Application) Invoke(fn any) error {
	return app.Container.Invoke(fn)
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

// Add records a provider's contributions (command constructors, job handlers,
// event listeners, cron jobs) in one bucket — *Application is the single
// support.Registry. Each runtime later claims the kinds it understands from
// Contributions(); the console kernel picks out the command constructors.
func (app *Application) Add(contributions ...any) {
	app.contributions = append(app.contributions, contributions...)
}

// Contributions returns everything providers have added, for the runtimes
// (worker/eventbus/scheduler) to filter by their concrete contract type.
func (app *Application) Contributions() []any {
	return app.contributions
}

// bootstrap runs the ordered bootstrappers (load config, register + boot
// providers, …) — the analog of Laravel's Application::bootstrapWith.
func (app *Application) bootstrap() error {
	for _, b := range app.bootstrappers {
		if err := b.Bootstrap(app); err != nil {
			return err
		}
	}
	return nil
}

// HandleCommand bootstraps the application, dispatches the CLI command, then
// terminates — returning the process exit code. It mirrors Laravel's
// $app->handleCommand() → kernel.handle()/terminate(): bootstrap errors are
// reported here (cobra hasn't run yet to render them); command errors are
// already rendered to stderr by cobra, so they map straight to exit code 1; and
// Shutdown (terminate) runs on every path via defer.
func (app *Application) HandleCommand(use, short string) (code int) {
	defer func() {
		if err := app.Shutdown(); err != nil {
			logger.Error("shutdown", "error", err)
			if code == 0 {
				code = 1
			}
		}
	}()

	if err := app.bootstrap(); err != nil {
		logger.Error("bootstrap failed", "error", err)
		return 1
	}
	if err := app.consoleKernel.Handle(use, short, app.contributions); err != nil {
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
