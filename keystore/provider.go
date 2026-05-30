package keystore

import (
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
	p.app.Singleton(NewKeyStore)
}

func (p *ServiceProvider) Boot() {}

// NewDefaultConfig returns a KeyStore Config populated from the application
// config, unmarshalling the keystore subtree directly (the keys map has
// arbitrary names, so yaml tags do the work field-by-field extraction used to).
func NewDefaultConfig() Config {
	var cfg Config
	if err := config.Unmarshal("keystore", &cfg); err != nil {
		panic(err)
	}
	return cfg
}
