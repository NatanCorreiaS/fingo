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

// Compiler directive below
//
//go:embed schema.sql
var schemaSQL string

// CheckAndCreate verifies if the database file exists and creates it with the schema if not.
func CheckAndCreate() error {
	if _, err := os.Stat("fingo.db"); err == nil {
		log.Printf("fingo.db already detected! Skipping database creation...")
		return nil
	}

	if err := createDatabase(); err != nil {
		return err
	}

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
