package support_test

import (
	"testing"

	"github.com/iVampireSP/beacon/support"
)

// fakeRegistry captures contributions so the base can be tested without the
// real application.
type fakeRegistry struct{ added []any }

func (r *fakeRegistry) Add(contributions ...any) { r.added = append(r.added, contributions...) }

// demoProvider embeds the base and declares contributions in Register, as a real
// provider does.
type demoProvider struct {
	support.ServiceProvider
}

func (p *demoProvider) Register() { p.Add(func() {}, "a-job") }

func TestServiceProviderAddPushesToRegistry(t *testing.T) {
	r := &fakeRegistry{}
	p := &demoProvider{ServiceProvider: support.ServiceProvider{Registry: r}}

	p.Register()

	if len(r.added) != 2 {
		t.Fatalf("want 2 contributions pushed to the registry, got %d", len(r.added))
	}
}

// Compile-time check that the embedded base satisfies the Provider contract.
var _ support.Provider = (*demoProvider)(nil)
