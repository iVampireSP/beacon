package httpserver

import (
	"time"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type HTTPServerServiceProvider struct {
	app contracts.Application
}

func NewHTTPServerServiceProvider(app contracts.Application) support.Provider {
	return &HTTPServerServiceProvider{app: app}
}

func (p *HTTPServerServiceProvider) Register() {
	p.app.Singleton(NewDefaultMetricsConfig)
}

func (p *HTTPServerServiceProvider) Boot() {}

func NewDefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:         config.Bool("metrics.enabled", true),
		Host:            config.String("metrics.host", "0.0.0.0"),
		Port:            config.Int("metrics.port", 9090),
		ShutdownTimeout: 30 * time.Second,
	}
}
