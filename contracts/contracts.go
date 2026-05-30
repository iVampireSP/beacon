// Package contracts holds the framework's core interfaces, mirroring
// Illuminate\Contracts. They invert dependencies so service providers, the
// Foundation kernels, and the bootstrappers depend on small, focused interfaces
// instead of the concrete foundation.Application — which keeps the package graph
// acyclic and the components decoupled, mockable, and maintainable. Unlike
// Laravel's ~32-method Contracts\Foundation\Application, this is sized to exactly
// what the framework uses (interface segregation, Go-idiomatic).
package contracts

import "github.com/iVampireSP/beacon/support"

// Application is the application surface that service providers and the
// bootstrappers depend on — the decoupled analog of
// Illuminate\Contracts\Foundation\Application. The concrete foundation.Application
// satisfies it. Commands that need generics (e.g. ProvidersImplementing) take the
// concrete type instead; everything else depends on this interface.
type Application interface {
	// Singleton registers a constructor whose result is built once and cached.
	Singleton(constructor any) error
	// Invoke resolves the function's parameters from the container and calls it,
	// returning the function's error (if any). Used by a provider's Boot to act
	// on resolved services.
	Invoke(fn any) error
	// Register registers (and immediately Register()s) service providers.
	Register(providers ...support.Provider)
	// RegisterCommands records console-command constructors (a provider pushes
	// commands here via support.ServiceProvider.AddCommand).
	RegisterCommands(constructors ...any)
	// Boot runs every registered provider's Boot phase (once).
	Boot() error
	// OnShutdown registers a cleanup callback, run in reverse order on Shutdown.
	OnShutdown(fn func() error)
	// Shutdown runs the cleanup callbacks.
	Shutdown() error
}

// ProviderConstructor builds a service provider from the application — the Go
// analog of a class reference in Laravel's bootstrap/providers.php `return [...]`
// (Go can't new up a class from a name, so the list holds constructors). The
// RegisterProviders bootstrapper calls each with the application: new Provider($app).
type ProviderConstructor = func(Application) support.Provider

// ConsoleKernel is the console entry point — the analog of
// Illuminate\Contracts\Console\Kernel. foundation.Application.HandleCommand
// delegates to it so the application need not import the concrete kernel.
type ConsoleKernel interface {
	// Handle bootstraps the application and runs the CLI under the given root
	// command identity, returning the command's error (nil on success).
	Handle(use, short string) error
	// RegisterCommands records command constructors to build when Handle runs.
	RegisterCommands(constructors ...any)
}
