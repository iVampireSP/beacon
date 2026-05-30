package support

// Provider is the service provider contract: the application registers and
// boots these. Register() binds services into the container; Boot() runs after
// all providers have been registered. Mirrors the registrable role of Laravel's
// Illuminate\Support\ServiceProvider.
type Provider interface {
	Register()
	Boot()
}

// Kernel is the command-registrar capability a ServiceProvider needs to declare
// its console commands. It is declared here (a tiny interface) rather than
// imported from the foundation/console kernel so support stays a pure leaf and
// never imports upward — dependency inversion. *foundation.Application satisfies
// it (forwarding RegisterCommands to its console kernel), so a provider that
// holds the application can push commands through it.
type Kernel interface {
	RegisterCommands(constructors ...any)
}

// ServiceProvider is the embeddable base, mirroring Illuminate\Support\ServiceProvider.
// It holds the console kernel so a provider can declare its commands with
// AddCommand — the analog of Laravel's $this->commands([...]). Concrete providers
// embed it and override Register/Boot as needed.
type ServiceProvider struct {
	Kernel Kernel
}

// Register binds services into the container. The embedded default is a no-op;
// concrete providers override it.
func (ServiceProvider) Register() {}

// Boot runs after all providers are registered. The embedded default is a
// no-op; concrete providers override it.
func (ServiceProvider) Boot() {}

// AddCommand registers console commands by their constructors with the kernel,
// which resolves each through the container (constructor injection). This is the
// analog of Laravel's ServiceProvider::commands().
func (p ServiceProvider) AddCommand(commands ...any) {
	p.Kernel.RegisterCommands(commands...)
}
