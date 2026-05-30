package console

import (
	"github.com/iVampireSP/foundation/container"
)

// Kernel is the console kernel — the bootstrap orchestrator for the CLI entry
// point, mirroring Laravel's Illuminate\Foundation\Console\Kernel. It boots the
// container and then delegates command execution to a console.Application (the
// Artisan). Service providers push their command constructors here via
// support.ServiceProvider.AddCommand (so *Kernel satisfies support.Kernel); on
// Run the Kernel hands them to the Artisan, which resolves each through the
// container. The DI container never deals with commands.
//
// This is the command-side counterpart of the service runtime: console.Kernel
// + console.Application run CLI commands, just as transport.App runs servers.
type Kernel struct {
	app              *container.Application
	commandFactories []any
}

// NewKernel creates a console kernel over the application.
func NewKernel(app *container.Application) *Kernel {
	return &Kernel{app: app}
}

// RegisterCommands records command constructors to be built when the kernel
// runs. Providers call this through support.ServiceProvider.AddCommand; *Kernel
// satisfies support.Kernel.
func (k *Kernel) RegisterCommands(constructors ...any) {
	k.commandFactories = append(k.commandFactories, constructors...)
}

// Run boots the application, builds the Artisan with every registered command,
// executes it, and then shuts the application down. Mirrors the Laravel Kernel's
// handle()/terminate(): bootstrap, delegate to Artisan, clean up. It returns the
// command error (or shutdown error) for the entry point to map to an exit code.
func (k *Kernel) Run(use, short string) error {
	if err := k.app.Boot(); err != nil {
		return err
	}

	artisan := NewApplication(k.app, use, short)
	for _, ctor := range k.commandFactories {
		if err := artisan.Add(ctor); err != nil {
			return err
		}
	}

	err := artisan.Run()
	if shutdownErr := k.app.Shutdown(); shutdownErr != nil && err == nil {
		err = shutdownErr
	}
	return err
}
