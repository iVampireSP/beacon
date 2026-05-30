package console

import (
	"github.com/iVampireSP/foundation/container"
	"github.com/iVampireSP/foundation/version"
	"github.com/spf13/cobra"
)

// Application is the console application that owns and runs the CLI commands —
// the analog of Laravel's Illuminate\Console\Application (Artisan). It resolves
// each command's constructor through the container (constructor injection) and
// assembles the cobra root. It is the command-running counterpart to the
// service-running transport.App; the Kernel bootstraps the framework and then
// delegates execution here.
type Application struct {
	app  *container.Application
	root *cobra.Command
}

// NewApplication creates the Artisan over the container, with the given root
// command identity.
func NewApplication(app *container.Application, use, short string) *Application {
	return &Application{
		app: app,
		root: &cobra.Command{
			Use:     use,
			Short:   short,
			Version: version.String(),
		},
	}
}

// Add resolves a command constructor through the container and registers the
// resulting *cobra.Command on the root. The constructor declares its
// dependencies as parameters, and may take *container.Application to reach the
// container directly.
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
// container and returns the *cobra.Command it produces. Each constructor runs
// in its own child scope so several constructors providing *container.Application
// (or the same dependency type) don't collide.
func (a *Application) build(constructor any) (*cobra.Command, error) {
	var cmd *cobra.Command
	scope := a.app.Scope().Scope("command")
	if err := scope.Provide(func() *container.Application { return a.app }); err != nil {
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
