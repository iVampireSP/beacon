// Package console holds the Foundation console kernel — the bootstrap
// orchestrator for the CLI entry point, mirroring
// Illuminate\Foundation\Console\Kernel. It is distinct from the Artisan engine in
// the top-level console package (imported here as `artisan`, exactly as Laravel
// aliases Illuminate\Console\Application as Artisan).
package console

import (
	"errors"

	artisan "github.com/iVampireSP/beacon/console"
	"github.com/iVampireSP/beacon/container"
	"github.com/iVampireSP/beacon/contracts"
)

// Kernel boots the application and delegates command execution to the Artisan.
// It depends on the application through the contracts.Application interface (not
// the concrete foundation.Application) so this package never imports foundation —
// that inversion is what keeps the foundation ↔ foundation/console graph acyclic.
// Service providers push their command constructors here via
// support.ServiceProvider.AddCommand; the kernel hands them to the Artisan, which
// resolves each through the container.
type Kernel struct {
	app              contracts.Application
	container        *container.Container
	commandFactories []any
}

// NewKernel creates a console kernel over the application and its container.
func NewKernel(app contracts.Application, c *container.Container) *Kernel {
	return &Kernel{app: app, container: c}
}

// RegisterCommands records command constructors to be built when the kernel
// runs. Providers call this through support.ServiceProvider.AddCommand.
func (k *Kernel) RegisterCommands(constructors ...any) {
	k.commandFactories = append(k.commandFactories, constructors...)
}

// Handle boots the application, builds the Artisan with every registered
// command, executes it, then shuts the application down. Mirrors the Laravel
// Kernel's handle()/terminate(): bootstrap, delegate to Artisan, clean up.
//
// Providers register their shutdown hooks during Register (before Handle), so
// Shutdown is deferred to run on every exit path — including a failed boot or a
// command-build error — and its error is joined with the command's.
func (k *Kernel) Handle(use, short string) (err error) {
	defer func() {
		if shutdownErr := k.app.Shutdown(); shutdownErr != nil {
			err = errors.Join(err, shutdownErr)
		}
	}()

	if err = k.app.Boot(); err != nil {
		return err
	}

	app := artisan.NewApplication(k.container, use, short)
	for _, ctor := range k.commandFactories {
		if err = app.Add(ctor); err != nil {
			return err
		}
	}

	return app.Run()
}
