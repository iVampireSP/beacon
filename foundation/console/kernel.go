// Package console holds the Foundation console kernel, mirroring
// Illuminate\Foundation\Console\Kernel. It is distinct from the Artisan engine in
// the top-level console package (imported here as `artisan`, exactly as Laravel
// aliases Illuminate\Console\Application as Artisan).
package console

import (
	"reflect"

	artisan "github.com/iVampireSP/beacon/console"
	"github.com/iVampireSP/beacon/container"
	"github.com/spf13/cobra"
)

// cobraCommandType is *cobra.Command, used to pick command constructors out of
// the application's single contribution bucket.
var cobraCommandType = reflect.TypeOf((*cobra.Command)(nil))

// Kernel builds the Artisan from the command constructors among the
// application's contributions and runs it — the command-dispatch half of the
// console entry point. Bootstrapping and termination are orchestrated by
// foundation.Application.HandleCommand around this call. The kernel needs only
// the container (to resolve commands), so it stays decoupled from the application.
type Kernel struct {
	container *container.Container
}

// NewKernel creates a console kernel over the DI container.
func NewKernel(c *container.Container) *Kernel {
	return &Kernel{container: c}
}

// Handle builds every command constructor among contributions and executes the
// Artisan. A contribution is a command constructor when it is a function
// returning *cobra.Command; everything else (job handlers, listeners, cron jobs)
// is claimed by other runtimes and ignored here. The returned error is cobra's
// (already rendered); HandleCommand maps it to the process exit code.
func (k *Kernel) Handle(use, short string, contributions []any) error {
	app := artisan.NewApplication(k.container, use, short)
	for _, c := range contributions {
		if !isCommandConstructor(c) {
			continue
		}
		if err := app.Add(c); err != nil {
			return err
		}
	}
	return app.Run()
}

// isCommandConstructor reports whether c is a function returning *cobra.Command.
func isCommandConstructor(c any) bool {
	t := reflect.TypeOf(c)
	if t == nil || t.Kind() != reflect.Func || t.NumOut() != 1 {
		return false
	}
	return t.Out(0) == cobraCommandType
}
