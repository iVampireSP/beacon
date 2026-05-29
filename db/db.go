package db

import (
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

// driverName is the database/sql driver registered by pgx/v5/stdlib.
const driverName = "pgx"

// GetDSN returns the PostgreSQL DSN (keyword/value form) from config.
func GetDSN(cfg Config) string {
	sslMode := cfg.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		sslMode,
	)
}

// NewDB opens a PostgreSQL connection pool with OTel SQL tracing.
// It returns a *sql.DB; the ORM (ent) binding is layered on top by the
// application, keeping this package driver-specific but ORM-agnostic.
func NewDB(cfg Config) (*sql.DB, error) {
	conn, err := otelsql.Open(driverName, GetDSN(cfg),
		otelsql.WithAttributes(
			semconv.DBSystemNamePostgreSQL,
			semconv.DBNamespace(cfg.Name),
		),
		otelsql.WithSpanOptions(otelsql.SpanOptions{
			DisableErrSkip: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	// Configure connection pool to prevent connection exhaustion.
	conn.SetMaxOpenConns(cfg.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	conn.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return conn, nil
}
