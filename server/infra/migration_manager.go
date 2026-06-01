package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationManager struct {
	db            *sql.DB
	migrationpath string
	migrate       *migrate.Migrate
}

func NewMigrationManager(db *sql.DB, path string) (*MigrationManager, error) {
	// Create postgres driver from existing connection
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Create migrate instance with file source and postgres driver
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return &MigrationManager{
		db:            db,
		migrationpath: path,
		migrate:       m,
	}, nil
}

func (mm *MigrationManager) Up() error {
	log.Println("Applying all migrations...")

	if err := mm.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No pending migrations left to apply")
			return nil
		}
		return fmt.Errorf("migrating up failed: %w", err)
	}

	version, _, _ := mm.migrate.Version()
	log.Println("Successfully migrated to version %i", version)
	return nil
}

func (mm *MigrationManager) UpN(n int) error {
	log.Println("Applying next", n, "migrations...")

	if err := mm.migrate.Steps(n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No pending migrations to apply")
		}
		return fmt.Errorf("migration up %d stepts failed: %w", err)
	}

	return nil
}

func (mm *MigrationManager) Down(n int) error {
	log.Println("Reverting last", n, "migrations...")

	if err := mm.migrate.Steps(-n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No pending migrations to revert")
		}
		return fmt.Errorf("migration down %d stepts failed: %w", err)
	}

	return nil
}

func (mm *MigrationManager) Goto(version uint) error {
	log.Println("Migrating to version", version)

	if err := mm.migrate.Migrate(version); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Already at version", version)
		}
		return fmt.Errorf("migration to version %d failed: %w", version, err)
	}

	return nil
}

func (mm *MigrationManager) Version() (uint, bool, error) {

}

func (mm *MigrationManager) Force(version int) error {
	log.Printf("Forcing migration version to %d...", version)

	return mm.migrate.Force(version)
}

func (mm *MigrationManager) Close() error {
	srcErr, dbErr := mm.migrate.Close()
	if srcErr != nil {
		return srcErr
	}

	return dbErr
}

func (mm *MigrationManager) MigrateWithLock(ctx context.Context) error {

}
