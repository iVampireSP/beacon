package tracing

import (
	"time"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/foundation"
)

type ServiceProvider struct {
	app *foundation.Application
}

func NewServiceProvider(app *foundation.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewTracing)
}

func (p *ServiceProvider) Boot() {
	// Initialize the global tracing state from config so tracing.GetService
	// works for every command. With tracing.enabled=false (the default) this
	// records the "disabled" sentinel and GetService becomes a no-op; enabling
	// it wires the OTLP exporters.
	_ = p.app.Invoke(func(cfg Config) { _, _ = NewTracing(cfg) })
}

// NewDefaultConfig returns a Tracing Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		Enabled:       config.Bool("tracing.enabled", false),
		Endpoint:      config.String("tracing.endpoint", "localhost:4317"),
		SampleRatio:   config.GetFloat("tracing.sample_ratio", 1.0),
		QueryURL:      config.String("tracing.query_url", "http://jaeger-query:16686"),
		QueryUsername: config.String("tracing.query_username"),
		QueryPassword: config.String("tracing.query_password"),
		QueryTimeout:  config.Duration("tracing.query_timeout", 15*time.Second),
	}
}
