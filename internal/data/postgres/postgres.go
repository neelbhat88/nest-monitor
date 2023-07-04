package postgres

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	migrateiofs "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type MigrationSource interface {
	GetMigrations() PostgresMigrations
}

type DatabaseConfig struct {
	Port     int64  `env:"PGPORT" env-default:"5432"`
	Host     string `env:"PGHOST" env-default:"localhost"`
	Name     string `env:"PGNAME" env-default:"postgres"`
	User     string `env:"PGUSER" env-default:"user"`
	Password string `env:"PGPASS" env-default:"password"`
	SSLMode  bool   `env:"PGSSLENABLED" env-default:"true"`
}

type PostgresMigrations struct {
	SchemaName     string
	MigrationFiles fs.FS
	Path           string
}

func InitializeDB(params DatabaseConfig, migrationSource MigrationSource) (*sqlx.DB, error) {
	db, err := ConnectPostgres(params)
	if err != nil {
		log.Error().Err(err).Msg("InitializeDB failed to ConnectPostgres")
		return db, err
	}

	if migrationSource != nil {
		migrationParams := (migrationSource).GetMigrations()
		err = RunPostgresMigrations(db, migrationParams)
		if err != nil {
			log.Error().Err(err).Msg("RunPostgresMigrationsFailed")

			return db, err
		}
	}

	return db, nil
}

func ConnectPostgres(params DatabaseConfig) (*sqlx.DB, error) {
	validParams, err := validParams(params)
	if err != nil {
		return nil, err
	}
	sslmode := ""
	if !validParams.SSLMode {
		sslmode = "sslmode=disable"
	}
	pgconn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s %s",
		validParams.Host, validParams.Port, validParams.Name, validParams.User, validParams.Password,
		sslmode)
	db, err := sqlx.Connect("postgres", pgconn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func RunPostgresMigrations(db *sqlx.DB, migrationParams PostgresMigrations) error {

	m, err := getMigrations(db, migrationParams.SchemaName, migrationParams.MigrationFiles, migrationParams.Path)
	if err != nil {
		return err
	}

	// Run the migrations ignore error for none needed
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func DropPostgresMigrations(db *sqlx.DB, migrationParams PostgresMigrations) error {
	m, err := getMigrations(db, migrationParams.SchemaName, migrationParams.MigrationFiles, migrationParams.Path)
	if err != nil {
		log.Error().Err(err).Msg("failed to get migrations")
		return err
	}

	version, isDirty, err := m.Version()
	if err != nil {
		return err
	}

	if isDirty {
		m.Force(int(version - 1))
	}

	err = m.Down()
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to migrate down because: %v", err))
		return err
	}

	sourceErr, dbErr := m.Close()

	if sourceErr != nil {
		return fmt.Errorf("migrate Close failed with source err: %v", sourceErr)
	}

	if dbErr != nil {
		return fmt.Errorf("migrate Close failed with database err: %v", dbErr)
	}
	return nil
}

// ForcePostgresMigrations should only be used when necessary (manual process)
func ForcePostgresMigrations(db *sqlx.DB, migrationParams PostgresMigrations, targetVersion int) error {
	m, err := getMigrations(db, migrationParams.SchemaName, migrationParams.MigrationFiles, migrationParams.Path)
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to get migrations because: %v", err))
		return err
	}

	err = m.Force(targetVersion)
	if err != nil {
		return err
	}

	sourceErr, dbErr := m.Close()

	if sourceErr != nil {
		return fmt.Errorf("migrate Close failed with source err: %v", sourceErr)
	}

	if dbErr != nil {
		return fmt.Errorf("migrate Close failed with database err: %v", dbErr)
	}

	return nil
}

// DownPostgresMigrations should only be used when necessary (manual process)
func DownPostgresMigrations(db *sqlx.DB, migrationParams PostgresMigrations) error {
	m, err := getMigrations(db, migrationParams.SchemaName, migrationParams.MigrationFiles, migrationParams.Path)
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to get migrations because: %v", err))
		return err
	}

	err = m.Down()
	if err != nil {
		log.Error().Err(fmt.Errorf("failed to migrate down because: %v", err))
		return err
	}

	sourceErr, dbErr := m.Close()

	if sourceErr != nil {
		return fmt.Errorf("migrate Close failed with source err: %v", sourceErr)
	}

	if dbErr != nil {
		return fmt.Errorf("migrate Close failed with database err: %v", dbErr)
	}

	return nil
}

func validParams(params DatabaseConfig) (DatabaseConfig, error) {
	validParams := params
	if validParams.Host == "" {
		validParams.Host = "localhost"
	}
	if validParams.Port == 0 {
		validParams.Port = 5432
	}
	if validParams.User == "" {
		return validParams, errors.New("No username provided")
	}
	if validParams.Password == "" {
		return validParams, errors.New("No password provided")
	}
	if validParams.Name == "" {
		return validParams, errors.New("No database name provided")
	}
	return validParams, nil
}

func getMigrations(db *sqlx.DB, schemaName string, migrationsFS fs.FS, path string) (*migrate.Migrate, error) {
	config := postgres.Config{}

	// Make sure the domains schema exists in the database if one is provided for partitioning
	if schemaName != "" {
		schemaSQL := "CREATE SCHEMA IF NOT EXISTS " + schemaName + ";"
		_, err := db.Exec(schemaSQL)
		if err != nil {
			return nil, err
		}
		config.SchemaName = schemaName
	}

	// Creates a driver to run migrations against the existing database
	driver, err := postgres.WithInstance(db.DB, &config)
	if err != nil {
		return nil, err
	}

	// Source driver to read files via the io FS interface
	mfs, err := migrateiofs.New(migrationsFS, path)
	if err != nil {
		return nil, err
	}

	// Create the migration instance
	m, err := migrate.NewWithInstance("iofs", mfs, "postgres", driver)
	if err != nil {
		return nil, err
	}
	return m, nil
}
