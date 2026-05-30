package support_test

import (
	"testing"

	"github.com/iVampireSP/beacon/support"
)

// fakeKernel captures RegisterCommands so the base can be tested without the
// real console kernel.
type fakeKernel struct{ registered []any }

func (k *fakeKernel) RegisterCommands(ctors ...any) { k.registered = append(k.registered, ctors...) }

// demoProvider embeds the base and declares commands in Register, as a real
// provider does.
type demoProvider struct {
	support.ServiceProvider
}

func (p *demoProvider) Register() { p.AddCommand(func() {}, func() {}) }

func TestServiceProviderAddCommandPushesToKernel(t *testing.T) {
	k := &fakeKernel{}
	p := &demoProvider{ServiceProvider: support.ServiceProvider{Kernel: k}}

	p.Register()

	if len(k.registered) != 2 {
		t.Fatalf("want 2 registered command constructors, got %d", len(k.registered))
	}
}

// Compile-time check that the embedded base satisfies the Provider contract.
var _ support.Provider = (*demoProvider)(nil)
