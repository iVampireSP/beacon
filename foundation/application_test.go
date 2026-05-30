package foundation

import (
	"strings"
	"testing"
)

type recordProvider struct{}

func (recordProvider) Register() {}
func (recordProvider) Boot()     {}

func TestBootRunsAllProviders(t *testing.T) {
	app := newApplication()
	app.Register(recordProvider{}, recordProvider{})

	if err := app.Boot(); err != nil {
		t.Fatalf("Boot returned error: %v", err)
	}
	if got := len(app.Providers()); got != 2 {
		t.Fatalf("expected 2 registered providers, got %d", got)
	}
}

func TestAddCollectsContributions(t *testing.T) {
	app := newApplication()
	app.Add("a", 1)
	app.Add(struct{}{})
	if got := len(app.Contributions()); got != 3 {
		t.Fatalf("Contributions = %d, want 3", got)
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
