package keystore

import (
	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

type KeyStoreServiceProvider struct {
	app contracts.Application
}

func NewKeyStoreServiceProvider(app contracts.Application) support.Provider {
	return &KeyStoreServiceProvider{app: app}
}

func (p *KeyStoreServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewKeyStore)
}

func (p *KeyStoreServiceProvider) Boot() {}

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
