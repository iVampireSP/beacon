package console

import (
	"testing"

	"github.com/spf13/cobra"
)

type fakeJob struct{}

func (fakeJob) JobName() string { return "x" }

// TestIsCommandConstructor verifies the kernel picks command constructors
// (funcs returning *cobra.Command) out of the mixed contribution bucket and
// ignores work-unit values.
func TestIsCommandConstructor(t *testing.T) {
	cases := []struct {
		name string
		item any
		want bool
	}{
		{"no-arg ctor", func() *cobra.Command { return nil }, true},
		{"dep ctor", func(int) *cobra.Command { return nil }, true},
		{"job value", fakeJob{}, false},
		{"plain func", func() {}, false},
		{"wrong return", func() int { return 0 }, false},
		{"nil", nil, false},
		{"string", "serve", false},
	}
	for _, c := range cases {
		if got := isCommandConstructor(c.item); got != c.want {
			t.Errorf("%s: isCommandConstructor = %v, want %v", c.name, got, c.want)
		}
	}
}
