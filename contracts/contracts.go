// Package contracts holds the framework's core interfaces, mirroring
// Illuminate\Contracts. They invert dependencies so the Foundation integration
// layer and its kernels don't form import cycles: a Foundation kernel depends on
// the contracts.Application interface (not the concrete foundation.Application),
// and foundation.Application drives the console via the contracts.ConsoleKernel
// interface (not the concrete kernel). This package imports nothing — it is a
// pure leaf of interfaces.
package contracts

// Application is the bootstrap-able application a kernel drives — the analog of
// Illuminate\Contracts\Foundation\Application. The concrete foundation.Application
// satisfies it; declaring it here lets the kernels in foundation/console and
// foundation/http depend on the interface instead of importing the foundation
// package (which would cycle).
type Application interface {
	// Boot runs every registered service provider's Boot phase (once).
	Boot() error
	// Shutdown runs cleanup callbacks in reverse registration order.
	Shutdown() error
}

// ConsoleKernel is the console entry point — the analog of
// Illuminate\Contracts\Console\Kernel. foundation.Application.HandleCommand
// delegates to it (via this interface) so the application need not import the
// concrete foundation/console kernel. Providers push command constructors to it
// through support.ServiceProvider.AddCommand.
type ConsoleKernel interface {
	// Handle bootstraps the application and runs the CLI under the given root
	// command identity, returning the command's error (nil on success).
	Handle(use, short string) error
	// RegisterCommands records command constructors to build when Handle runs.
	RegisterCommands(constructors ...any)
}
