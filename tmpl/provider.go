package tmpl

import (
	"io/fs"

	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
)

// NewTemplateServiceProvider returns a provider that parses the templates under
// fsys/dir into the template set on Register — the analog of Laravel's
// ViewServiceProvider. It is an ordinary provider listed in the app's
// WithProviders, carrying its own embedded fs, so foundation never imports tmpl.
func NewTemplateServiceProvider(fsys fs.FS, dir string) contracts.ProviderConstructor {
	return func(contracts.Application) support.Provider {
		return &TemplateServiceProvider{fs: fsys, dir: dir}
	}
}

// TemplateServiceProvider parses templates on Register.
type TemplateServiceProvider struct {
	fs  fs.FS
	dir string
}

func (p *TemplateServiceProvider) Register() {
	MustInitWithFS(p.fs, p.dir)
}
func (p *TemplateServiceProvider) Boot() {}
