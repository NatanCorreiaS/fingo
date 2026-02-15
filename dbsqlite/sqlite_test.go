package dbsqlite

import (
	"database/sql"
	"os"
	"testing"
)

func TestCreateSQLDatabase(t *testing.T) {
	os.Remove("fingo.db")
	tests := []struct {
		name        string
		tableExists bool
		wantError   bool
	}{
		{"Success: Without tables present", false, false},
		{"Error: With tables present", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.tableExists {
				// Criamos o banco e a tabela manualmente antes de rodar a função
				db, _ := sql.Open("sqlite3", "file:fingo.db")
				db.Exec("CREATE TABLE usuarios (id INTEGER PRIMARY KEY);")
				db.Close()
			}
			defer os.Remove("fingo.db")

			err := createDatabase()
			if (err != nil) != tt.wantError {
				t.Errorf("createSQLiteDatabase() error = %v, wantError = %v\n", err, tt.wantError)
			}

		})
	}
}

func TestCheckAndCreate(t *testing.T) {
	os.Remove("fingo.db")
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
			if tt.dbExists {
				os.WriteFile("fingo.db", nil, 0644)
			}
			defer os.Remove("fingo.db")

			if err := CheckAndCreate(); (err != nil) != tt.wantError {

				t.Errorf("CheckAndCreate() error = %v, wantError = %v\n", err, tt.wantError)
			}
		})
	}

}
