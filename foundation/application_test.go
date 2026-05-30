package foundation

import (
	"strings"
	"testing"
)

// capable is a local capability interface used to test ProvidersImplementing.
type capable interface{ Capability() }

type capableProvider struct{}

func (capableProvider) Register()   {}
func (capableProvider) Boot()       {}
func (capableProvider) Capability() {}

type plainProvider struct{}

func (plainProvider) Register() {}
func (plainProvider) Boot()     {}

func TestBootRunsAllProviders(t *testing.T) {
	app := newApplication()
	app.Register(capableProvider{}, plainProvider{})

	if err := app.Boot(); err != nil {
		t.Fatalf("Boot returned error: %v", err)
	}
	if got := len(app.Providers()); got != 2 {
		t.Fatalf("expected 2 registered providers, got %d", got)
	}
}

func TestProvidersImplementingFiltersByCapability(t *testing.T) {
	app := newApplication()
	app.Register(capableProvider{}, plainProvider{})

	if got := ProvidersImplementing[capable](app); len(got) != 1 {
		t.Fatalf("expected exactly 1 capable provider, got %d", len(got))
	}
}

func TestBootFiresCallbacks(t *testing.T) {
	app := newApplication()
	var order []string
	app.Booting(func(*Application) { order = append(order, "booting") })
	app.Booted(func(*Application) { order = append(order, "booted") })

	if err := app.Boot(); err != nil {
		t.Fatalf("Boot: %v", err)
	}
	if got := strings.Join(order, ","); got != "booting,booted" {
		t.Fatalf("callback order = %q, want booting,booted", got)
	}

	// A Booted callback registered after boot must run immediately.
	app.Booted(func(*Application) { order = append(order, "immediate") })
	if order[len(order)-1] != "immediate" {
		t.Fatal("Booted callback registered post-boot did not run immediately")
	}
}
