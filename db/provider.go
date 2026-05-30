package db

import (
	"database/sql"
	"sync/atomic"

	"github.com/iVampireSP/beacon/contracts"
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
	app contracts.Application
}

func NewDatabaseServiceProvider(app contracts.Application) support.Provider {
	return &DatabaseServiceProvider{ServiceProvider: support.ServiceProvider{Registry: app}, app: app}
}

func (p *DatabaseServiceProvider) Register() {
	p.app.Singleton(NewDefaultConfig)

	// Capture the pool when (and only when) it is actually constructed, and
	// register the close hook NOW, during Register, so it runs in this
	// provider's shutdown slot. Shutdown is reverse-registration order, so this
	// foundational resource — registered early — closes late, after anything
	// that may use it. Registering the close lazily (on first resolve) or in a
	// Booted callback would append it last and thus close the DB first.
	var pool atomic.Pointer[sql.DB]
	p.app.Singleton(func(cfg Config) (*sql.DB, error) {
		db, err := NewDB(cfg)
		if err == nil {
			pool.Store(db)
		}
		return db, err
	})
	p.app.OnShutdown(func() error {
		if db := pool.Load(); db != nil {
			return db.Close()
		}
		return nil
	})

	p.Add(NewMigrate)
}
