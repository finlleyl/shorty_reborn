package database

import (
	"fmt"

	"github.com/finlleyl/shorty_reborn/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewDB(cfg *config.Database) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	switch cfg.Driver {
	case "postgres":
		db, err = postgresDB(cfg)
	default:
		return nil, fmt.Errorf("driver not supported: %s", cfg.Driver)
	}
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func postgresDB(cfg *config.Database) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
	
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	db.SetConnMaxLifetime(cfg.Timeout)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	schema := []string{
		`CREATE TABLE IF NOT EXISTS url (
			id SERIAL PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`,
	}

	for _, stmt := range schema {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
