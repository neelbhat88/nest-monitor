package postgres

import (
	"embed"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type NestMonitorDB struct {
	*sqlx.DB
}

func GetMigrations() PostgresMigrations {
	return PostgresMigrations{
		SchemaName:     "public",
		MigrationFiles: migrationFiles,
		Path:           "migrations",
	}
}

func NewPostgresRepository(db *sqlx.DB) NestMonitorDB {
	postgresRepository := NestMonitorDB{
		DB: db,
	}

	return postgresRepository
}
