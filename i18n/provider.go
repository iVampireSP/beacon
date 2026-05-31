package i18n

import (
	"io/fs"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// NewDefaultConfig returns an i18n Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		DefaultLocale:  config.String("app.locale", "en"),
		FallbackLocale: config.String("app.fallback_locale", "en"),
	}
}

// NewI18nServiceProvider returns a provider that loads the i18n message catalogs
// from fsys/dir on Register — the analog of Laravel's TranslationServiceProvider.
// It is an ordinary provider listed in the app's WithProviders, carrying its own
// embedded fs, so foundation never imports i18n.
func NewI18nServiceProvider(fsys fs.FS, dir string) contracts.ProviderConstructor {
	return func(contracts.Application) support.Provider {
		return &I18nServiceProvider{fs: fsys, dir: dir}
	}
}

// I18nServiceProvider loads translation catalogs on Register.
type I18nServiceProvider struct {
	fs  fs.FS
	dir string
}

func (p *I18nServiceProvider) Register() { MustInitWithFS(NewDefaultConfig(), p.fs, p.dir) }
func (p *I18nServiceProvider) Boot()     {}
