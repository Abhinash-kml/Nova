package entitlements_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Global pool shared across all tests in this package
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 1. Start Testcontainers Postgres
	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	// Ensure the container terminates after all tests complete
	defer func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	// 2 & 3. Execute specific migration files in dependency order
	// List parent tables before child tables to satisfy foreign keys
	migrationsDir := "../../migrations/postgres"
	specificMigrations := []string{
		"000001_users.up.sql",        // Parent table
		"000009_player_items.up.sql", // Child table with FK to users
		"000016_entitlements.up.sql",
	}

	if err := runSpecificMigrations(connStr, migrationsDir, specificMigrations); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	// 4. Initialize shared pgxpool for all tests
	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to create pgxpool: %v", err)
	}
	defer testPool.Close()

	// Run all tests in the package
	exitCode := m.Run()
	os.Exit(exitCode)
}

// runSpecificMigrations filters the directory to run only the specified .up.sql files in sequence
func runSpecificMigrations(dbURL, dirPath string, targetFiles []string) error {
	// Create a temporary in-memory/temp directory or sub-FS containing only requested files
	tempDir, err := os.MkdirTemp("", "test-migrations-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Re-index files sequentially (000001_, 000002_...) so golang-migrate runs them in exact slice order
	for i, fileName := range targetFiles {
		srcPath := filepath.Join(dirPath, fileName)
		fmt.Println("sourcepath:", srcPath)
		content, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", srcPath, err)
		}

		// Strip original prefix and apply sequential version number
		destName := fmt.Sprintf("%06d_%s", i+1, fileName)
		destPath := filepath.Join(tempDir, destName)
		if err := os.WriteFile(destPath, content, 0600); err != nil {
			return fmt.Errorf("write temp migration %s: %w", destPath, err)
		}
	}

	sourceDriver, err := iofs.New(os.DirFS(tempDir), ".")
	if err != nil {
		return fmt.Errorf("create iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dbURL)
	if err != nil {
		return fmt.Errorf("init migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}

// --- Example Usage Across Multiple Tests ---

// cleanupTables resets state between tests without violating foreign keys
func cleanupTables(t *testing.T, ctx context.Context) {
	t.Helper()
	_, err := testPool.Exec(ctx, `TRUNCATE TABLE users, orders RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

func Test_Integration_Create(t *testing.T) {
	query := `
		INSERT INTO 
			entitlements(
				entitlement_id,
				user_id,
				item_id,
				source,
				is_active
			) VALUES (
				$1,
				$2,
				$3,
				'purchase',
				$4 
			) RETURNING entitlement_id;
	`

	entitlementID, _ := uuid.Parse("4ccba445-a0a7-41f1-bc6b-228937fb6f4b")
	userID, _ := uuid.Parse("40ef63e4-7e79-4e4e-bd7f-49a9ae1ce84a")
	itemID, _ := uuid.Parse("1d9fe0fb-c42c-4686-b224-52653727150f")
	isActive := true

	var gotID uuid.UUID
	err := testPool.QueryRow(context.Background(), query, entitlementID, userID, itemID, isActive).Scan(&gotID)
	assert.Nil(t, err)
	assert.Equal(t, entitlementID, gotID)
}
