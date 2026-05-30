package container

import (
	"errors"
	"os"

	"github.com/iVampireSP/foundation/console"
	"github.com/iVampireSP/foundation/support"
	"github.com/spf13/cobra"
)

// Application manages the full lifecycle of service providers and the DI container.
// It follows Laravel's Application pattern: Register → Boot → Run.
type Application struct {
	*Container
	providers   []support.ServiceProvider
	commands    []console.ConsoleCommand
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
// support.ServiceProvider directly), mirroring Laravel's
// $app->register(new Provider($app)). No adapter or reflection is involved.
//
// Usage:
//
//	app.Register(cache.NewServiceProvider(app), jwt.NewServiceProvider(app), ...)
func (app *Application) Register(providers ...support.ServiceProvider) {
	for _, p := range providers {
		p.Register()
		app.providers = append(app.providers, p)
	}
}

// Boot calls Boot() on all registered providers, then collects the console
// commands declared by any provider implementing console.CommandProvider.
// This is Phase 2 of the lifecycle, called after all providers are registered.
func (app *Application) Boot() error {
	if app.booted {
		return nil
	}
	for _, p := range app.providers {
		p.Boot()
	}
	// Collect commands declared via the CommandProvider capability, so each
	// module owns its own commands instead of a central registration list.
	for _, cp := range ProvidersImplementing[console.CommandProvider](app) {
		app.AddCommand(cp.Commands()...)
	}
	app.booted = true
	return nil
}

// Providers returns the registered service providers in registration order.
func (app *Application) Providers() []support.ServiceProvider {
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

// AddCommand registers console commands with the application.
// Providers call this in their Register() method.
func (app *Application) AddCommand(cmds ...console.ConsoleCommand) {
	app.commands = append(app.commands, cmds...)
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

// Run executes the full application lifecycle: Boot → collect commands → cobra.Execute.
func (app *Application) Run(use, short string) error {
	if err := app.Boot(); err != nil {
		return err
	}

	root := &cobra.Command{
		Use:   use,
		Short: short,
	}

	// Add all registered commands
	for _, cc := range app.commands {
		root.AddCommand(cc.Command())
	}

	// Shutdown hook
	root.PersistentPostRunE = func(_ *cobra.Command, _ []string) error {
		return app.Shutdown()
	}

	if err := root.Execute(); err != nil {
		_ = app.Shutdown()
		os.Exit(1)
	}
	return nil
}
