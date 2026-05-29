package command

import (
	"database/sql"

	"github.com/iVampireSP/foundation/console"
	"github.com/iVampireSP/foundation/container"
	"github.com/iVampireSP/foundation/db"
)

// DBServiceProvider wires the PostgreSQL connection pool (*sql.DB) and owns the
// `migrate` command. The ORM (ent) client is layered on top app-side; this
// provider stays ORM-agnostic. It lives alongside the command to avoid a
// db -> db/command import cycle.
type DBServiceProvider struct {
	app *container.Application
}

func NewDBServiceProvider(app *container.Application) *DBServiceProvider {
	return &DBServiceProvider{app: app}
}

func (p *DBServiceProvider) Register() {
	p.app.Singleton(db.NewDefaultConfig)
	p.app.Singleton(db.NewDB)
	p.app.OnShutdown(func() error {
		if c, err := container.Make[*sql.DB](p.app.Container); err == nil && c != nil {
			return c.Close()
		}
		return nil
	})
}

func (p *DBServiceProvider) Boot() {}

// Commands implements console.CommandProvider.
func (p *DBServiceProvider) Commands() []console.ConsoleCommand {
	return []console.ConsoleCommand{NewMigrate(p.app)}
}
