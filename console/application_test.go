package console

import (
	"testing"

	"github.com/iVampireSP/beacon/container"
	"github.com/spf13/cobra"
)

type buildDep struct{ name string }

// TestApplicationAddInjectsDependencies verifies the Artisan resolves a command
// constructor's parameters from the container (constructor injection) and
// registers the resulting command on the root.
func TestApplicationAddInjectsDependencies(t *testing.T) {
	c := container.NewContainer()
	if err := c.Singleton(func() *buildDep { return &buildDep{name: "svc"} }); err != nil {
		t.Fatalf("singleton: %v", err)
	}
	artisan := NewArtisan(c, "beacon", "test root")

	var gotDep *buildDep
	err := artisan.Add(func(dep *buildDep) *cobra.Command {
		gotDep = dep
		return &cobra.Command{Use: "demo"}
	})
	if err != nil {
		t.Fatalf("add: %v", err)
	}

	if gotDep == nil || gotDep.name != "svc" {
		t.Fatalf("dependency not injected: %+v", gotDep)
	}

	var found bool
	for _, cmd := range artisan.root.Commands() {
		if cmd.Use == "demo" {
			found = true
		}
	}
	if !found {
		t.Fatal("command was not registered on the root")
	}
}
