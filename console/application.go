package console

import (
	"github.com/iVampireSP/beacon/container"
	"github.com/iVampireSP/beacon/version"
	"github.com/spf13/cobra"
)

// Application is the console application that owns and runs the CLI commands —
// the analog of Laravel's Illuminate\Console\Application (Artisan). It resolves
// each command's constructor through the container (constructor injection) and
// assembles the cobra root. The Foundation console kernel bootstraps the
// framework and then delegates execution here.
type Application struct {
	container *container.Container
	root      *cobra.Command
}

// NewApplication creates the Artisan over the DI container, with the given root
// command identity.
func NewApplication(c *container.Container, use, short string) *Application {
	return &Application{
		container: c,
		root: &cobra.Command{
			Use:     use,
			Short:   short,
			Version: version.String(),
		},
	}
}

// Add resolves a command constructor through the container and registers the
// resulting *cobra.Command on the root. The constructor declares its
// dependencies as parameters; the foundation.Application is registered as a
// container singleton, so a command may inject it like any other dependency.
func (a *Application) Add(constructor any) error {
	cmd, err := a.build(constructor)
	if err != nil {
		return err
	}
	a.root.AddCommand(cmd)
	return nil
}

// Run executes the root command, dispatching to the selected subcommand.
func (a *Application) Run() error {
	return a.root.Execute()
}

// build invokes a command constructor with its parameters injected from the
// container and returns the *cobra.Command it produces. Each constructor runs in
// its own child scope so several constructors requesting the same dependency
// type don't collide.
func (a *Application) build(constructor any) (*cobra.Command, error) {
	var cmd *cobra.Command
	scope := a.container.Scope().Scope("command")
	if err := scope.Provide(constructor); err != nil {
		return nil, err
	}
	if err := scope.Invoke(func(c *cobra.Command) { cmd = c }); err != nil {
		return nil, err
	}
	return cmd, nil
}
