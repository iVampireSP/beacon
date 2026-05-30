package console

import (
	"testing"

	"github.com/iVampireSP/foundation/container"
	"github.com/spf13/cobra"
)

type buildDep struct{ name string }

// TestKernelBuildInjectsDependencies verifies the kernel resolves a command
// constructor's parameters from the container (and injects *Application).
func TestKernelBuildInjectsDependencies(t *testing.T) {
	app := container.NewApplication()
	if err := app.Singleton(func() *buildDep { return &buildDep{name: "svc"} }); err != nil {
		t.Fatalf("singleton: %v", err)
	}
	k := NewKernel(app)

	var gotDep *buildDep
	var gotApp *container.Application
	cmd, err := k.build(func(dep *buildDep, a *container.Application) *cobra.Command {
		gotDep = dep
		gotApp = a
		return &cobra.Command{Use: "demo"}
	})
	if err != nil {
		t.Fatalf("build: %v", err)
	}

	if cmd == nil || cmd.Use != "demo" {
		t.Fatalf("unexpected command: %+v", cmd)
	}
	if gotDep == nil || gotDep.name != "svc" {
		t.Fatalf("dependency not injected: %+v", gotDep)
	}
	if gotApp != app {
		t.Fatal("*Application not injected into command constructor")
	}
}
