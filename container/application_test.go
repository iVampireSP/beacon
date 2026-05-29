package container

import (
	"testing"

	"github.com/iVampireSP/foundation/console"
	"github.com/spf13/cobra"
)

type fakeCommand struct{ use string }

func (c fakeCommand) Command() *cobra.Command { return &cobra.Command{Use: c.use} }

// commandProvider implements both ServiceProvider and console.CommandProvider.
type commandProvider struct{ cmds []console.ConsoleCommand }

func (p *commandProvider) Register()                          {}
func (p *commandProvider) Boot()                              {}
func (p *commandProvider) Commands() []console.ConsoleCommand { return p.cmds }

// plainProvider implements only ServiceProvider (no capabilities).
type plainProvider struct{}

func (plainProvider) Register() {}
func (plainProvider) Boot()     {}

func TestBootCollectsCommandsFromCommandProviders(t *testing.T) {
	app := NewApplication()
	app.Register(
		Adapt(func(*Application) *commandProvider {
			return &commandProvider{cmds: []console.ConsoleCommand{fakeCommand{"a"}, fakeCommand{"b"}}}
		}),
		Adapt(func(*Application) plainProvider { return plainProvider{} }),
	)

	if err := app.Boot(); err != nil {
		t.Fatalf("Boot returned error: %v", err)
	}
	if got := len(app.commands); got != 2 {
		t.Fatalf("expected 2 commands collected from CommandProvider, got %d", got)
	}
}

func TestProvidersImplementingFiltersByCapability(t *testing.T) {
	app := NewApplication()
	app.Register(
		Adapt(func(*Application) *commandProvider { return &commandProvider{} }),
		Adapt(func(*Application) plainProvider { return plainProvider{} }),
	)

	cps := ProvidersImplementing[console.CommandProvider](app)
	if len(cps) != 1 {
		t.Fatalf("expected exactly 1 CommandProvider, got %d", len(cps))
	}
	if len(app.Providers()) != 2 {
		t.Fatalf("expected 2 registered providers, got %d", len(app.Providers()))
	}
}

// Note: a constructor that does not return a support.ServiceProvider (e.g.
// func(*Application) int) is now rejected at COMPILE time by Adapt's type
// parameter constraint, so there is no runtime panic case left to test.
