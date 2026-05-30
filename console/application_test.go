package console

import (
	"testing"

	"github.com/iVampireSP/beacon/container"
	"github.com/spf13/cobra"
)

type buildDep struct{ name string }

// TestApplicationAddInjectsDependencies verifies the Artisan resolves a command
// constructor's parameters from the container (injecting *Application too) and
// registers the resulting command on the root.
func TestApplicationAddInjectsDependencies(t *testing.T) {
	app := container.NewApplication()
	if err := app.Singleton(func() *buildDep { return &buildDep{name: "svc"} }); err != nil {
		t.Fatalf("singleton: %v", err)
	}
	artisan := NewApplication(app, "foundation", "test root")

	var gotDep *buildDep
	var gotApp *container.Application
	err := artisan.Add(func(dep *buildDep, a *container.Application) *cobra.Command {
		gotDep = dep
		gotApp = a
		return &cobra.Command{Use: "demo"}
	})
	if err != nil {
		t.Fatalf("add: %v", err)
	}

	if gotDep == nil || gotDep.name != "svc" {
		t.Fatalf("dependency not injected: %+v", gotDep)
	}
	if gotApp != app {
		t.Fatal("*Application not injected into command constructor")
	}

	var found bool
	for _, c := range artisan.root.Commands() {
		if c.Use == "demo" {
			found = true
		}
	}
	if !found {
		t.Fatal("command was not registered on the root")
	}
}

// TestKernelRegisterCommandsAccumulates verifies pushed constructors accumulate
// for the kernel to hand to the Artisan on Run.
func TestKernelRegisterCommandsAccumulates(t *testing.T) {
	k := NewKernel(container.NewApplication())
	k.RegisterCommands(func() *cobra.Command { return &cobra.Command{Use: "a"} })
	k.RegisterCommands(
		func() *cobra.Command { return &cobra.Command{Use: "b"} },
		func() *cobra.Command { return &cobra.Command{Use: "c"} },
	)
	if len(k.commandFactories) != 3 {
		t.Fatalf("commandFactories = %d, want 3", len(k.commandFactories))
	}
}
