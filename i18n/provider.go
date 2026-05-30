package i18n

import "github.com/iVampireSP/beacon/config"

// NewDefaultConfig returns an i18n Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		DefaultLocale:  config.String("app.locale", "en"),
		FallbackLocale: config.String("app.fallback_locale", "en"),
	}
}
