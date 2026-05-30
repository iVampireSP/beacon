// Package console holds the Foundation console kernel, mirroring
// Illuminate\Foundation\Console\Kernel. It is distinct from the Artisan engine in
// the top-level console package (imported here as `artisan`, exactly as Laravel
// aliases Illuminate\Console\Application as Artisan).
package console

import (
	artisan "github.com/iVampireSP/beacon/console"
	"github.com/iVampireSP/beacon/container"
)

// Kernel builds the Artisan with every registered command and runs it — the
// command-dispatch half of the console entry point. Bootstrapping (load config,
// register + boot providers) and termination (Shutdown) are orchestrated by
// foundation.Application.HandleCommand around this call, mirroring Laravel's
// handleCommand → kernel.handle → terminate. The kernel needs only the container
// (to resolve commands), so it stays decoupled from the application.
type Kernel struct {
	container        *container.Container
	commandFactories []any
}

// NewKernel creates a console kernel over the DI container.
func NewKernel(c *container.Container) *Kernel {
	return &Kernel{container: c}
}

// RegisterCommands records command constructors to be built when the kernel
// runs. Providers call this through support.ServiceProvider.AddCommand (forwarded
// by foundation.Application.RegisterCommands).
func (k *Kernel) RegisterCommands(constructors ...any) {
	k.commandFactories = append(k.commandFactories, constructors...)
}

// Handle builds the Artisan with every registered command and executes it. It is
// called after the application has been bootstrapped, so all provider-pushed
// commands are already registered. The returned error is cobra's (already
// rendered to stderr); HandleCommand maps it to the process exit code.
func (k *Kernel) Handle(use, short string) error {
	app := artisan.NewApplication(k.container, use, short)
	for _, ctor := range k.commandFactories {
		if err := app.Add(ctor); err != nil {
			return err
		}
	}
	return app.Run()
}
