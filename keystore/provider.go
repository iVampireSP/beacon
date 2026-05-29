package keystore

import (
	"github.com/iVampireSP/foundation/config"
	"github.com/iVampireSP/foundation/container"
)

type ServiceProvider struct {
	app *container.Application
}

func NewServiceProvider(app *container.Application) *ServiceProvider {
	return &ServiceProvider{app: app}
}

func (p *ServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewKeyStore)
}

func (p *ServiceProvider) Boot() {}

// NewDefaultConfig returns a KeyStore Config populated from the application config.
func NewDefaultConfig() Config {
	keysRaw := config.Map("keystore.keys")
	keys := make(map[string]KeyConfig, len(keysRaw))
	for name, v := range keysRaw {
		kc, ok := v.(map[string]any)
		if !ok {
			continue
		}
		keys[name] = KeyConfig{
			Type:       stringVal(kc, "type"),
			PrivateKey: stringVal(kc, "private_key"),
			PublicKey:  stringVal(kc, "public_key"),
		}
	}
	return Config{Keys: keys}
}

func stringVal(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}
