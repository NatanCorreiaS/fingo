package dbsqlite

import (
	"database/sql"
	"os"
	"testing"
)

// TestDatabaseLifecycle is a table-driven test that covers database creation,
// schema initialization and the connection helper in fewer, idiomatic subtests.
func TestDatabaseLifecycle(t *testing.T) {
	tests := []struct {
		name             string
		createUsersTable bool // if true, create a users table before calling createDatabase
		preCreateDB      bool // if true, create an empty fingo.db before running CheckAndCreate
		wantCreateError  bool
	}{
		{"createDatabase: no existing tables", false, false, false},
		{"createDatabase: users table already present -> error", true, false, true},
		{"CheckAndCreate: fingo.db already present", false, true, false},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// Ensure clean environment for each subtest
			_ = os.Remove("fingo.db")

			// Optionally create a DB file with an existing users table to trigger schema creation error
			if tc.createUsersTable {
				db, err := sql.Open("sqlite3", "file:fingo.db")
				if err != nil {
					t.Fatalf("setup: failed to open sqlite for setup: %v", err)
				}
				// create a users table so schema init in createDatabase fails
				if _, err := db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY);"); err != nil {
					_ = db.Close()
					t.Fatalf("setup: failed to create users table: %v", err)
				}
				_ = db.Close()
			}

			// Optionally create an empty file to simulate existing database
			if tc.preCreateDB {
				if err := os.WriteFile("fingo.db", nil, 0o644); err != nil {
					t.Fatalf("setup: could not create dummy fingo.db: %v", err)
				}
			}

			// Clean up after subtest
			defer func() { _ = os.Remove("fingo.db") }()

			// If we pre-created the DB file, exercise CheckAndCreate and GetDatabaseConnection
			if tc.preCreateDB {
				if err := CheckAndCreate(); err != nil {
					t.Fatalf("CheckAndCreate() error = %v", err)
				}

				db, err := GetDatabaseConnection()
				if err != nil {
					t.Fatalf("GetDatabaseConnection() returned error: %v", err)
				}
				defer db.Close()

				if err := db.Ping(); err != nil {
					t.Fatalf("db.Ping() failed: %v", err)
				}
				return
			}

			// Otherwise test createDatabase behavior (success or expected failure)
			err := createDatabase()
			if (err != nil) != tc.wantCreateError {
				t.Fatalf("createDatabase() error = %v, wantError = %v", err, tc.wantCreateError)
			}
		})
	}
}
