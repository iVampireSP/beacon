package container

import "github.com/spf13/cobra"

// BuildCommand invokes constructor with its parameters injected from the
// container and returns the *cobra.Command it produces. The constructor declares
// its dependencies as parameters (constructor injection) and may also take
// *Application to reach the container directly.
//
// This makes *Application satisfy support.App, so support.Base can build
// commands without importing this package — inverting the dependency and
// avoiding an import cycle (container already imports support).
func (app *Application) BuildCommand(constructor any) (*cobra.Command, error) {
	var cmd *cobra.Command
	// Isolated child scope: exposes *Application and the constructor, so several
	// constructors that return *cobra.Command do not collide.
	scope := app.Scope().Scope("command")
	if err := scope.Provide(func() *Application { return app }); err != nil {
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
