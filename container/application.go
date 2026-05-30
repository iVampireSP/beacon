package container

import (
	"errors"

	"github.com/iVampireSP/foundation/support"
)

// Application manages the lifecycle of service providers over the DI container.
// It follows Laravel's Application pattern: Register → Boot. It is a PURE DI +
// provider-lifecycle manager — it knows nothing about cobra, console commands,
// or the microservice runtime; those live in the console kernel and the
// transport package respectively.
type Application struct {
	*Container
	providers   []support.Provider
	shutdownFns []func() error
	booted      bool
}

// NewApplication creates a new application with a fresh container.
func NewApplication() *Application {
	return &Application{
		Container: NewContainer(),
	}
}

// Register hands service providers to the application. Each provider's
// Register() is invoked immediately so it can bind dependencies into the service
// container; the provider is then retained for Boot and capability collection.
//
// Providers are passed as constructed instances (a *XxxServiceProvider satisfies
// support.Provider directly), mirroring Laravel's
// $app->register(new Provider($app)). No adapter or reflection is involved.
func (app *Application) Register(providers ...support.Provider) {
	for _, p := range providers {
		p.Register()
		app.providers = append(app.providers, p)
	}
}

// Boot calls Boot() on all registered providers. This is Phase 2 of the
// lifecycle, called after all providers are registered. Command collection is
// NOT done here — the console kernel collects commands when it runs.
func (app *Application) Boot() error {
	if app.booted {
		return nil
	}
	for _, p := range app.providers {
		p.Boot()
	}
	app.booted = true
	return nil
}

// Providers returns the registered service providers in registration order.
func (app *Application) Providers() []support.Provider {
	return app.providers
}

// ProvidersImplementing returns all registered providers that satisfy the
// capability interface T (e.g. console.CommandProvider, job.HandlerProvider,
// bus.ListenerProvider, schedule.CronProvider). It lets a subsystem gather
// contributions from every provider without a central registration list.
func ProvidersImplementing[T any](app *Application) []T {
	var out []T
	for _, p := range app.providers {
		if t, ok := any(p).(T); ok {
			out = append(out, t)
		}
	}
	return out
}

// OnShutdown registers a cleanup callback to be called during shutdown.
// Callbacks are executed in reverse registration order.
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
