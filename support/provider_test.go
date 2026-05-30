package support_test

import (
	"testing"

	"github.com/iVampireSP/foundation/support"
)

// fakeApp captures RegisterCommands so the base can be tested without a real
// container.
type fakeApp struct{ registered []any }

func (f *fakeApp) RegisterCommands(ctors ...any) { f.registered = append(f.registered, ctors...) }

// demoProvider embeds the base, as a real provider would, and registers two
// command constructors in Register.
type demoProvider struct {
	support.ServiceProvider
}

func (p *demoProvider) Register() { p.AddCommand(func() {}, func() {}) }

func TestServiceProviderCommandsPushToApplication(t *testing.T) {
	app := &fakeApp{}
	p := &demoProvider{ServiceProvider: support.ServiceProvider{App: app}}

	p.Register()

	if len(app.registered) != 2 {
		t.Fatalf("want 2 registered command constructors, got %d", len(app.registered))
	}
}

// Compile-time check that the embedded base satisfies the Provider contract.
var _ support.Provider = (*demoProvider)(nil)
