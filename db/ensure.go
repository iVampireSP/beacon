package db

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

// EnsureDatabase makes sure the configured database (cfg.Name) exists. When it
// is missing and the connecting role may create databases, it asks for
// confirmation on stdin (default yes) and creates it. It returns true when it
// created the database.
//
// It connects to the maintenance "postgres" database on the same server; the
// target database itself is never opened here (the ORM opens it lazily after).
func EnsureDatabase(ctx context.Context, cfg Config) (bool, error) {
	adminCfg := cfg
	adminCfg.Name = "postgres"

	admin, err := sql.Open(driverName, GetDSN(adminCfg))
	if err != nil {
		return false, err
	}
	defer func() { _ = admin.Close() }()

	var exists bool
	if err := admin.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.Name,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("check database %q exists: %w", cfg.Name, err)
	}
	if exists {
		return false, nil
	}

	var mayCreate bool
	if err := admin.QueryRowContext(ctx,
		"SELECT rolcreatedb OR rolsuper FROM pg_roles WHERE rolname = current_user",
	).Scan(&mayCreate); err != nil {
		return false, fmt.Errorf("check create-database privilege: %w", err)
	}
	if !mayCreate {
		return false, fmt.Errorf("database %q does not exist and the current role cannot create databases", cfg.Name)
	}

	if !confirmCreateDatabase(cfg.Name) {
		return false, fmt.Errorf("database %q does not exist", cfg.Name)
	}

	if _, err := admin.ExecContext(ctx, "CREATE DATABASE "+quoteIdentifier(cfg.Name)); err != nil {
		return false, fmt.Errorf("create database %q: %w", cfg.Name, err)
	}
	return true, nil
}

// confirmCreateDatabase prompts to create the database, defaulting to yes.
func confirmCreateDatabase(name string) bool {
	fmt.Printf("Database %q does not exist. Create it? [Y/n]: ", name)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "", "y", "yes":
		return true
	default:
		return false
	}
}

// quoteIdentifier quotes a PostgreSQL identifier for safe use in DDL, where bind
// parameters are not allowed.
func quoteIdentifier(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}
