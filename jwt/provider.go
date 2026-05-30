package jwt

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type ServiceProvider struct {
	app contracts.Application
}

func NewServiceProvider(app contracts.Application) support.Provider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewJWT)
}

func (p *ServiceProvider) Boot() {}

// NewDefaultConfig returns a JWT Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		KeyName:   config.String("jwt.key", "rsa"),
		Issuer:    config.String("discovery.issuer", "http://localhost"),
		ExpiresIn: config.Int("jwt.expires_in", 86400),
	}
}
