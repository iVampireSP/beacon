package db

import (
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

// Fresh drops all tables in the public schema and re-runs all migrations.
func (m *Migrate) Fresh(cmd *cobra.Command) error {
	force, _ := cmd.Flags().GetBool("force")

	if !confirmDangerousOperation("drop ALL tables and rebuild database", force) {
		fmt.Println("Operation cancelled.")
		return nil
	}

	db := m.db

	rows, err := db.Query(`SELECT tablename FROM pg_tables WHERE schemaname = 'public'`)
	if err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			rows.Close()
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, table)
	}
	rows.Close()

	// CASCADE handles foreign-key dependencies (no session-level FK toggle in PostgreSQL).
	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf(`DROP TABLE IF EXISTS %q CASCADE`, table)); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("[OK] Database rebuilt")
	return printVersion(db)
}
