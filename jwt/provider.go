package jwt

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type JWTServiceProvider struct {
	app contracts.Application
}

func NewJWTServiceProvider(app contracts.Application) support.Provider {
	return &JWTServiceProvider{app: app}
}

func (p *JWTServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewJWT)
}

func (p *JWTServiceProvider) Boot() {}

// NewDefaultConfig returns a JWT Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		KeyName:   config.String("jwt.key", "rsa"),
		Issuer:    config.String("discovery.issuer", "http://localhost"),
		ExpiresIn: config.Int("jwt.expires_in", 86400),
	}
}
