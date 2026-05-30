package foundation

import "testing"

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
