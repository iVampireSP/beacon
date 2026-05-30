package support

// Provider is the service provider contract: the application registers and
// boots these. Register() binds services into the container; Boot() runs after
// all providers have been registered. Mirrors the registrable role of Laravel's
// Illuminate\Support\ServiceProvider.
type Provider interface {
	Register()
	Boot()
}

// Kernel is the console kernel capability a ServiceProvider needs to register
// its commands. It is declared here — not imported from the console package — to
// invert the dependency: console imports container which imports support, so
// support must never import console. *console.Kernel satisfies it.
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
