package config

import (
	"io/fs"

	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// NewConfigServiceProvider returns a provider that loads configuration from
// fsys/dir into the config store. It is PRIVILEGED: the foundation registers it
// before the default providers (via Builder.WithConfig), so config is populated
// before any provider's Register/Boot reads it — the model goravel uses (config
// as base provider #1) and Laravel's LoadConfiguration bootstrapper. It carries
// its own embedded fs, so the app declares the source without foundation knowing
// the layout.
func NewConfigServiceProvider(fsys fs.FS, dir string) contracts.ProviderConstructor {
	return func(contracts.Application) support.Provider {
		return &ConfigServiceProvider{fs: fsys, dir: dir}
	}
}

// ConfigServiceProvider loads YAML configuration on Register.
type ConfigServiceProvider struct {
	fs  fs.FS
	dir string
}

func (p *ConfigServiceProvider) Register() { MustInitWithFS(p.fs, p.dir) }
func (p *ConfigServiceProvider) Boot()     {}
