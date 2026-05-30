package db

import (
	"time"

	"github.com/iVampireSP/beacon/config"
)

// NewDefaultConfig returns a db Config populated from the application config.
func NewDefaultConfig() Config {
	return Config{
		Host:            config.String("database.app.host", "localhost"),
		Port:            config.Int("database.app.port", 5432),
		User:            config.String("database.app.user", "postgres"),
		Password:        config.String("database.app.password"),
		Name:            config.String("database.app.name", "cloud"),
		SSLMode:         config.String("database.app.sslmode", "disable"),
		MaxOpenConns:    config.Int("database.app.max_open_conns", 25),
		MaxIdleConns:    config.Int("database.app.max_idle_conns", 5),
		ConnMaxLifetime: time.Duration(config.Int("database.app.conn_max_lifetime_seconds", 300)) * time.Second,
		ConnMaxIdleTime: time.Duration(config.Int("database.app.conn_max_idle_time_seconds", 60)) * time.Second,
	}
}
