package dbsqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "embed"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

const createMonthlyAdjustmentsTableSQL = `
CREATE TABLE IF NOT EXISTS monthly_adjustments_log(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	year_month TEXT NOT NULL UNIQUE,
	applied_at TEXT DEFAULT CURRENT_TIMESTAMP
);`

// Compiler directive below
//
//go:embed schema.sql
var schemaSQL string

// CheckAndCreate verifies if the database file exists and creates it with the schema if not.
// For existing databases, it ensures the monthly_adjustments_log table exists (migration).
func CheckAndCreate() error {
	if _, err := os.Stat("fingo.db"); err == nil {
		log.Printf("fingo.db already detected! Skipping database creation...")
		// Ensure the monthly_adjustments_log table exists for existing databases
		if err := EnsureMonthlyAdjustmentsTable(); err != nil {
			return fmt.Errorf("failed to ensure monthly_adjustments_log table: %w", err)
		}
		return nil
	}

	if err := createDatabase(); err != nil {
		return err
	}

	return nil
}

// EnsureMonthlyAdjustmentsTable creates the monthly_adjustments_log table if it doesn't already exist.
// This acts as a migration for databases that were created before the monthly adjustment feature.
func EnsureMonthlyAdjustmentsTable() error {
	db, err := GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(createMonthlyAdjustmentsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create monthly_adjustments_log table: %w", err)
	}

	log.Println("monthly_adjustments_log table ensured.")
	return nil
}

// GetDatabaseConnection opens and returns a connection to the SQLite database.
func GetDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "file:fingo.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open the database :%v", err)
	}

	return db, nil
}

func createDatabase() error {

	db, err := GetDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(schemaSQL)

	if err != nil {
		return fmt.Errorf("Failed to create tables: %v", err)
	}

	log.Println("Tables successfully created in the database...")
	return nil

}
