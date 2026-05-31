// Package bootstrap holds the ordered application bootstrappers, mirroring
// Illuminate\Foundation\Bootstrap. Each Bootstrapper is one focused, testable
// step of bringing the application up; the application runs them in order before
// dispatching a command. They depend only on the contracts.Application interface
// (never the concrete foundation.Application), so the boot sequence is decoupled.
//
// Resource loading is NOT here: config is a privileged provider (registered first
// via Builder.WithConfig), and translations/templates/migrations are ordinary
// providers in the app's WithProviders — mirroring Laravel, where config is the
// LoadConfiguration bootstrapper and the rest are service providers.
package bootstrap

import (
	"github.com/iVampireSP/beacon/contracts"
)

// Bootstrapper performs one ordered step of application bootstrapping. The
// framework owns this fixed sequence (RegisterProviders → BootProviders), like
// Laravel's Kernel::$bootstrappers — apps extend via providers, not bootstrappers.
type Bootstrapper interface {
	Bootstrap(app contracts.Application) error
}

// RegisterProviders registers the application's service providers — the analog of
// Illuminate\Foundation\Bootstrap\RegisterProviders. It calls each constructor
// with the application (new Provider($app)) and registers the result.
type RegisterProviders struct {
	Providers []contracts.ProviderConstructor
}

func (b RegisterProviders) Bootstrap(app contracts.Application) error {
	for _, newProvider := range b.Providers {
		app.Register(newProvider(app))
	}
	return nil
}

// BootProviders boots all registered providers — the analog of
// Illuminate\Foundation\Bootstrap\BootProviders.
type BootProviders struct{}

func (BootProviders) Bootstrap(app contracts.Application) error {
	return app.Boot()
}
