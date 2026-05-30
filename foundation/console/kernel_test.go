package console

import (
	"testing"

	"github.com/iVampireSP/beacon/container"
	"github.com/spf13/cobra"
)

// TestKernelRegisterCommandsAccumulates verifies pushed constructors accumulate
// for the kernel to hand to the Artisan on Handle.
func TestKernelRegisterCommandsAccumulates(t *testing.T) {
	k := NewKernel(container.NewContainer())
	k.RegisterCommands(func() *cobra.Command { return &cobra.Command{Use: "a"} })
	k.RegisterCommands(
		func() *cobra.Command { return &cobra.Command{Use: "b"} },
		func() *cobra.Command { return &cobra.Command{Use: "c"} },
	)
	if len(k.commandFactories) != 3 {
		t.Fatalf("commandFactories = %d, want 3", len(k.commandFactories))
	}
}
