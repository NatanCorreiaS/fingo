package dbsqlite

import (
	"database/sql"
	"os"
	"testing"
)

func TestCreateSQLDatabase(t *testing.T) {
	tests := []struct {
		name        string
		tableExists bool
		wantError   bool
	}{
		{"Success: Without tables present", false, false},
		{"Error: With users table already present", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ensure clean state for each subtest
			_ = os.Remove("fingo.db")
			if tt.tableExists {
				// create the database file and a users table to cause schema creation to fail
				db, err := sql.Open("sqlite3", "file:fingo.db")
				if err != nil {
					t.Fatalf("failed to open sqlite for setup: %v", err)
				}
				_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY);")
				if err != nil {
					db.Close()
					t.Fatalf("failed to create users table for setup: %v", err)
				}
				db.Close()
			}
			defer os.Remove("fingo.db")

			err := createDatabase()
			if (err != nil) != tt.wantError {
				t.Errorf("createDatabase() error = %v, wantError = %v\n", err, tt.wantError)
			}
		})
	}
}

func TestCheckAndCreate(t *testing.T) {
	tests := []struct {
		name      string
		dbExists  bool
		wantError bool
	}{
		{"Success: without fingo.db present", false, false},
		{"Success: with fingo.db present", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Remove("fingo.db")
			if tt.dbExists {
				if err := os.WriteFile("fingo.db", nil, 0644); err != nil {
					t.Fatalf("could not create dummy fingo.db: %v", err)
				}
			}
			defer os.Remove("fingo.db")

			if err := CheckAndCreate(); (err != nil) != tt.wantError {
				t.Errorf("CheckAndCreate() error = %v, wantError = %v\n", err, tt.wantError)
			}
		})
	}
}

func TestGetDatabaseConnection(t *testing.T) {
	// ensure the database file exists so connection targets an existing file
	_ = os.Remove("fingo.db")
	if err := os.WriteFile("fingo.db", nil, 0644); err != nil {
		t.Fatalf("could not create fingo.db for TestGetDatabaseConnection: %v", err)
	}
	defer os.Remove("fingo.db")

	db, err := GetDatabaseConnection()
	if err != nil {
		t.Fatalf("GetDatabaseConnection() returned error: %v", err)
	}
	defer db.Close()

	// verify we can Ping the DB to ensure connection is usable
	if err := db.Ping(); err != nil {
		t.Fatalf("db.Ping() failed: %v", err)
	}
}
