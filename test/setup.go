package test

import (
	"calculationengine/constants"
	storage "calculationengine/store"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB initializes a test database for integration tests.
// In CI environment, it uses the POSTGRES_TEST_DSN environment variable.
// Locally, it connects to the main database (or create a separate test one).
func SetupTestDB(t *testing.T) *gorm.DB {
	// Try to get test database DSN from environment (set by CI or locally)
	dsn := os.Getenv("POSTGRES_TEST_DSN")
	if dsn == "" {
		// Fallback: use default test database config
		// In your CI, set this env var to point to test database
		dsn = "host=127.0.0.1 user=postgres password=postgres dbname=postgres port=55432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Set the global DB variable used by store package
	storage.DB = db

	// Run migrations for tests
	if err := storage.AutoMigrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

// TeardownTestDB cleans up after tests.
// Truncates all tables to reset state for next test.
func TeardownTestDB(t *testing.T, db *gorm.DB) {
	if db == nil {
		return
	}

	// Clear all data (truncate tables in order of dependencies)
	db.Exec("TRUNCATE TABLE formula_dependencies CASCADE")
	db.Exec("TRUNCATE TABLE formulas CASCADE")
	db.Exec("TRUNCATE TABLE products CASCADE")
	db.Exec("TRUNCATE TABLE category_attribute_assignments CASCADE")
	db.Exec("TRUNCATE TABLE categories CASCADE")
	db.Exec("TRUNCATE TABLE attributes CASCADE")
}

// InitializeConstants loads configuration for tests
func InitializeConstants(t *testing.T) {
	constants.Load()
}

// LogTest helper function to print test information with [TEST] prefix
func LogTest(t *testing.T, format string, args ...interface{}) {
	t.Logf("[TEST] "+format, args...)
}
