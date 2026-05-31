package db

import (
	"database/sql"
	"io/fs"
	"strings"

	"github.com/iVampireSP/beacon/contracts"
	"github.com/iVampireSP/beacon/support"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

// MigrationSource is where versioned .sql migrations are read from: an embedded
// filesystem and the subpath within it holding the files. NewMigrationServiceProvider
// binds it as a container singleton so the migrate command resolves it FROM THE
// CONTAINER — no beacon globals (mirrors Laravel, where the migrator is a
// container service, not global state). goose's own classic API is process-global,
// which we accept (it's swappable); the FS source stays injected, so beacon holds
// no migration globals.
type MigrationSource struct {
	FS  fs.FS
	Dir string
}

// NewMigrationServiceProvider returns a provider that binds the migration source
// (so the migrate command resolves it from the container — no globals) and pushes
// the migrate command tree — the analog of Laravel's MigrationServiceProvider. It
// is an ordinary provider listed in the app's WithProviders, carrying its own
// embedded fs, so foundation never imports db.
func NewMigrationServiceProvider(fsys fs.FS, dir string) contracts.ProviderConstructor {
	return func(app contracts.Application) support.Provider {
		return &MigrationServiceProvider{app: app, fs: fsys, dir: dir}
	}
}

// MigrationServiceProvider binds the MigrationSource and registers the migrate
// command (up/down/status/fresh).
type MigrationServiceProvider struct {
	app contracts.Application
	fs  fs.FS
	dir string
}

func (p *MigrationServiceProvider) Register() {
	p.app.Singleton(func() MigrationSource {
		return MigrationSource{FS: p.fs, Dir: p.dir}
	})
	p.app.Add(NewMigrate)
}

func (p *MigrationServiceProvider) Boot() {}

// Migrate provides database migration commands. It holds no global state: the
// pool and migration source are injected per run in PersistentPreRunE.
type Migrate struct {
	app      contracts.Application
	db       *sql.DB
	hasFiles bool
}

// NewMigrate builds the migrate command tree. app is injected by the container;
// the db pool and migration source are resolved lazily in PersistentPreRunE so
// unrelated commands don't touch the database.
func NewMigrate(app contracts.Application) *cobra.Command {
	return (&Migrate{app: app}).Command()
}

// Command returns the migrate command tree.
func (m *Migrate) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate", Short: "Database migrations",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return m.app.Invoke(func(d *sql.DB, src MigrationSource) error {
				m.db = d
				sub, err := fs.Sub(src.FS, src.Dir)
				if err != nil {
					return err
				}
				m.hasFiles = hasSQL(sub)
				// Point goose at the DI-injected FS, per run, right before use —
				// not as a process-wide startup side-effect. goose's classic API is
				// itself global, which we accept (it's swappable); the FS source is
				// ours and stays injected, so beacon holds no migration globals.
				goose.SetBaseFS(sub)
				goose.SetTableName("migrations")
				return goose.SetDialect("postgres")
			})
		},
	}

	up := &cobra.Command{Use: "up [N]", Short: "Apply migrations",
		RunE: func(c *cobra.Command, _ []string) error { return m.Up(c) }}
	down := &cobra.Command{Use: "down [N]", Short: "Rollback migrations (default: 1 step)",
		RunE: func(c *cobra.Command, _ []string) error { return m.Down(c) }}
	down.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	status := &cobra.Command{Use: "status", Short: "Show migration status",
		RunE: func(c *cobra.Command, _ []string) error { return m.Status(c) }}

	fresh := &cobra.Command{Use: "fresh", Short: "Drop all tables and re-run migrations (DANGEROUS)",
		RunE: func(c *cobra.Command, _ []string) error { return m.Fresh(c) }}
	fresh.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	cmd.AddCommand(up, down, status, fresh)

	return cmd
}

// hasSQL reports whether the (already sub-rooted) migrations filesystem holds any
// .sql files.
func hasSQL(fsys fs.FS) bool {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			return true
		}
	}
	return false
}
