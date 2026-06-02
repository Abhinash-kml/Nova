package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

type MigrationManager struct {
	db             *sql.DB
	migrationpath  string
	startVersion   uint
	currentVersion uint
	migrate        *migrate.Migrate
	logger         *zap.Logger
}

func NewMigrationManager(db *sql.DB, path string, l *zap.Logger) (*MigrationManager, error) {
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

	version, _, _ := m.Version()

	return &MigrationManager{
		db:             db,
		migrationpath:  path,
		migrate:        m,
		startVersion:   version,
		currentVersion: version,
		logger:         l,
	}, nil
}

func (mm *MigrationManager) Up() error {
	mm.logger.Info("Applying all migrations...")

	if err := mm.migrate.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			mm.logger.Info("No pending migrations left to apply")
			return nil
		}
		return fmt.Errorf("migrating up failed: %w", err)
	}

	version, _, _ := mm.migrate.Version()
	mm.logger.Info("Successfully migrated", zap.Uint("version", version))
	return nil
}

func (mm *MigrationManager) UpN(n int) error {
	mm.logger.Sugar().Infof("Applying next", n, "migrations...")

	if err := mm.migrate.Steps(n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			mm.logger.Info("No pending migrations to apply")
		}
		return fmt.Errorf("migration up %d stepts failed: %w", n, err)
	}

	return nil
}

func (mm *MigrationManager) Down(n int) error {
	mm.logger.Sugar().Infof("Reverting last", n, "migrations...")

	if err := mm.migrate.Steps(-n); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			mm.logger.Info("No pending migrations to revert")
		}
		return fmt.Errorf("migration down %d stepts failed: %w", n, err)
	}

	return nil
}

func (mm *MigrationManager) Goto(version uint) error {
	mm.logger.Info("Migrating to version", zap.Uint("version", version))

	if err := mm.migrate.Migrate(version); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			mm.logger.Info("Already at version", zap.Uint("version", version))
		}
		return fmt.Errorf("migration to version %d failed: %w", version, err)
	}

	return nil
}

func (mm *MigrationManager) Version() (uint, bool, error) {
	version, dirty, err := mm.migrate.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return version, dirty, nil
}

func (mm *MigrationManager) Force(version int) error {
	mm.logger.Info("Forcing migration to version...", zap.Int("version", version))

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
	const lockId = 1245 //  TODO: Add configs for these

	// Acquire advisory lock
	var locked bool
	err := mm.db.QueryRowContext(ctx, "SELECT pq_try_advisory_lock($1);", lockId).Scan(&locked)
	if err != nil {
		return fmt.Errorf("failed to acquire advisory lock: %w", err)
	}

	if !locked {
		mm.logger.Info("Another migration is in progress, skipping...")
		return nil
	}

	// Ensure lock is released when done
	defer func() {
		mm.db.ExecContext(ctx, "SELECT pq_advisory_unlock($1);", lockId)
	}()

	// Run migrations
	return mm.Up()
}

func (mm *MigrationManager) MigrateWithRollback() error {
	defer mm.migrate.Close()

	mm.logger.Info("Starting migration with version", zap.Uint("version", mm.startVersion))

	var count int
	// Apply all migrations one at a time
	for {
		// Apply migration
		err := mm.migrate.Steps(1)
		if errors.Is(err, migrate.ErrNoChange) {
			count++
			mm.logger.Info("Successfully applied migration", zap.Int("count", count))
			break
		}

		if err != nil {
			mm.logger.Info("Migration failed", zap.Int("count", count), zap.Error(err))
			return mm.rollbackToStart()
		}

		// Update current version
		mm.currentVersion, _, _ = mm.migrate.Version()
		mm.logger.Info("Successfully applied migrations", zap.Uint("version", mm.currentVersion), zap.Int("count", count))
	}

	return nil
}

func (mm *MigrationManager) rollbackToStart() error {
	if mm.currentVersion == mm.startVersion {
		return fmt.Errorf("Migration failed, no rollback needed")
	}

	mm.logger.Info("Rolling back migrations", zap.Uint("from", mm.currentVersion), zap.Uint("to", mm.startVersion))

	// Calculate steps to rollback
	steps := int(mm.startVersion) - int(mm.currentVersion)

	if err := mm.migrate.Steps(steps); err != nil {
		return fmt.Errorf("Rollback failed: %w (manual intervention needed)", err)
	}

	mm.logger.Info("Successfully rolled back to version", zap.Uint("version", mm.startVersion))
	return errors.New("Migration failed, rolled back to previous version")
}
