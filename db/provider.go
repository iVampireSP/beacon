package db

import (
	"database/sql"

	"github.com/iVampireSP/beacon/container"
	"github.com/iVampireSP/beacon/support"
)

// DatabaseServiceProvider wires the PostgreSQL connection pool (*sql.DB) and
// owns the `migrate` command, mirroring Laravel's
// Illuminate\Database\DatabaseServiceProvider: the provider and its commands
// live together in the module root.
//
// The ORM (ent) client is layered on top app-side; this provider stays
// ORM-agnostic.
type DatabaseServiceProvider struct {
	support.ServiceProvider
	app *container.Application
}

func NewDatabaseServiceProvider(app *container.Application, kernel support.Kernel) *DatabaseServiceProvider {
	return &DatabaseServiceProvider{ServiceProvider: support.ServiceProvider{Kernel: kernel}, app: app}
}

func (p *DatabaseServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)
	p.app.Singleton(NewDB)
	p.app.OnShutdown(func() error {
		if c, err := container.Make[*sql.DB](p.app.Container); err == nil && c != nil {
			return c.Close()
		}
		return nil
	})
	p.AddCommand(NewMigrate)
}
