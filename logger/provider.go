package logger

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type LoggerServiceProvider struct {
	app contracts.Application
}

func NewLoggerServiceProvider(app contracts.Application) support.Provider {
	return &LoggerServiceProvider{app: app}
}

func (p *LoggerServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewLogger)
}

func (p *LoggerServiceProvider) Boot() {}

// NewDefaultConfig returns a logger Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		Level: config.String("log.level", "info"),
		Debug: config.Bool("app.debug", false),
	}
}
