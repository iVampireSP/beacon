package console

import (
	"os"

	"github.com/iVampireSP/foundation/container"
	"github.com/iVampireSP/foundation/version"
	"github.com/spf13/cobra"
)

// Kernel is the console application — the single home of all CLI/cobra handling,
// mirroring Laravel's Artisan (Illuminate\Console\Application). Service providers
// push their command constructors here via support.ServiceProvider.AddCommand
// (so *Kernel satisfies support.Kernel); the kernel resolves each through the
// container and assembles the cobra root. The DI container never deals with
// commands.
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

// build invokes a command constructor with its parameters injected from the
// container and returns the *cobra.Command it produces. The constructor may
// declare any container-managed dependency, and may take *container.Application
// to reach the container directly.
func (k *Kernel) build(constructor any) (*cobra.Command, error) {
	var cmd *cobra.Command
	scope := k.app.Scope().Scope("command")
	if err := scope.Provide(func() *container.Application { return k.app }); err != nil {
		return nil, err
	}
	if err := scope.Provide(constructor); err != nil {
		return nil, err
	}
	if err := scope.Invoke(func(c *cobra.Command) { cmd = c }); err != nil {
		return nil, err
	}
	return cmd, nil
}

// Run boots the application, builds every registered command through the
// container, assembles the cobra root, and executes it. The shutdown hooks run
// after the command completes.
func (k *Kernel) Run(use, short string) error {
	if err := k.app.Boot(); err != nil {
		return err
	}

	root := &cobra.Command{
		Use:     use,
		Short:   short,
		Version: version.String(),
	}
	for _, ctor := range k.commandFactories {
		cmd, err := k.build(ctor)
		if err != nil {
			return err
		}
		root.AddCommand(cmd)
	}

	root.PersistentPostRunE = func(_ *cobra.Command, _ []string) error {
		return k.app.Shutdown()
	}

	if err := root.Execute(); err != nil {
		_ = k.app.Shutdown()
		os.Exit(1)
	}
	return nil
}
