// Package bootstrap holds the ordered application bootstrappers, mirroring
// Illuminate\Foundation\Bootstrap. Each Bootstrapper is one focused, testable
// step of bringing the application up; the application runs them in order before
// dispatching a command. They depend only on the contracts.Application interface
// (never the concrete foundation.Application), so the boot sequence is decoupled
// and this package stays free of an import cycle with foundation.
package bootstrap

import (
	"io/fs"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/i18n"
	"github.com/iVampireSP/beacon/support"
	"github.com/iVampireSP/beacon/tmpl"
)

// Bootstrapper performs one ordered step of application bootstrapping. New steps
// (exception handling, facades, …) slot into the ordered list without touching
// the others.
type Bootstrapper interface {
	Bootstrap(app contracts.Application) error
}

// LoadConfiguration loads the application configuration from an embedded (or any)
// filesystem — the analog of Illuminate\Foundation\Bootstrap\LoadConfiguration.
// It must run first: later steps (and the default locale) read from config.
type LoadConfiguration struct {
	FS  fs.FS
	Dir string
}

func (b LoadConfiguration) Bootstrap(contracts.Application) error {
	return config.InitWithFS(b.FS, b.Dir)
}

// LoadTranslations loads the i18n message catalogs. It must run after
// LoadConfiguration, since the default locale comes from config.
type LoadTranslations struct {
	FS  fs.FS
	Dir string
}

func (b LoadTranslations) Bootstrap(contracts.Application) error {
	return i18n.InitWithFS(i18n.NewDefaultConfig(), b.FS, b.Dir)
}

// LoadTemplates parses the templates under Dir in the given filesystem.
type LoadTemplates struct {
	FS  fs.FS
	Dir string
}

func (b LoadTemplates) Bootstrap(contracts.Application) error {
	return tmpl.InitWithFS(b.FS, b.Dir)
}

// RegisterProviders registers the application's service providers — the analog of
// Illuminate\Foundation\Bootstrap\RegisterProviders. The factory receives the
// application so each provider can take it (new Provider($app)).
type RegisterProviders struct {
	Factory func(contracts.Application) []support.Provider
}

func (b RegisterProviders) Bootstrap(app contracts.Application) error {
	if b.Factory != nil {
		app.Register(b.Factory(app)...)
	}
	return nil
}

// BootProviders boots all registered providers — the analog of
// Illuminate\Foundation\Bootstrap\BootProviders.
type BootProviders struct{}

func (BootProviders) Bootstrap(app contracts.Application) error {
	return app.Boot()
}
