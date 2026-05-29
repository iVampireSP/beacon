package console

import "github.com/spf13/cobra"

// ConsoleCommand defines a command that can be registered with the application.
// This mirrors Laravel's Illuminate\Console\Command.
type ConsoleCommand interface {
	Command() *cobra.Command
}

// CommandProvider is an optional capability interface for service providers.
// A provider that implements it declares the console commands it owns; the
// application collects them during Boot. This mirrors Laravel's
// ServiceProvider $commands, and lets each module own its own commands instead
// of registering them in a central place.
type CommandProvider interface {
	Commands() []ConsoleCommand
}
