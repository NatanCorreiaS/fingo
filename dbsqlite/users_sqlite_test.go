package dbsqlite

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"natan/fingo/model"
	"natan/fingo/utils"
)

// setupDB prepares a clean database and returns an open connection and a teardown function.
func setupDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	_ = os.Remove("fingo.db")

	if err := createDatabase(); err != nil {
		t.Fatalf("failed to create database for test setup: %v", err)
	}

	db, err := GetDatabaseConnection()
	if err != nil {
		t.Fatalf("GetDatabaseConnection() returned error during setup: %v", err)
	}

	teardown := func() {
		_ = db.Close()
		_ = os.Remove("fingo.db")
	}

	return db, teardown
}

func TestUsersSQLite_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T, db *sql.DB)
	}{
		{
			name: "ReturnAllUsers on empty DB returns zero users",
			testFn: func(t *testing.T, db *sql.DB) {
				users, err := ReturnAllUsers(db)
				if err != nil {
					t.Fatalf("ReturnAllUsers() returned unexpected error on empty DB: %v", err)
				}
				if got := len(users); got != 0 {
					t.Fatalf("expected 0 users, got %d", got)
				}
			},
		},
		{
			name: "CreateUser inserts user and ReturnAllUsers returns it",
			testFn: func(t *testing.T, db *sql.DB) {
				u := model.User{
					UserName:       "Alice",
					CurrentAmount:  utils.Money(1000),
					MonthlyInputs:  utils.Money(2000),
					MonthlyOutputs: utils.Money(500),
				}

				uRet, err := CreateUser(u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}
				if uRet == nil {
					t.Fatalf("CreateUser() returned nil user")
				}
				if uRet.UserName != u.UserName {
					t.Errorf("returned username mismatch: expected %q, got %q", u.UserName, uRet.UserName)
				}
				// verify via ReturnAllUsers as well
				users, err := ReturnAllUsers(db)
				if err != nil {
					t.Fatalf("ReturnAllUsers() returned error after insert: %v", err)
				}
				if len(users) != 1 {
					t.Fatalf("expected 1 user, got %d", len(users))
				}
			},
		},
		{
			name: "CreateUser accepts username with apostrophe (safe parameterization)",
			testFn: func(t *testing.T, db *sql.DB) {
				u := model.User{
					UserName:       "O'Neil",
					CurrentAmount:  utils.Money(0),
					MonthlyInputs:  utils.Money(0),
					MonthlyOutputs: utils.Money(0),
				}

				uRet, err := CreateUser(u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned unexpected error for username with apostrophe: %v", err)
				}
				if uRet == nil {
					t.Fatalf("CreateUser() returned nil user")
				}

				// verify retrieval
				got, err := GetUserByID(uRet.ID, db)
				if err != nil {
					t.Fatalf("GetUserByID() returned error for inserted user: %v", err)
				}
				if got.UserName != u.UserName {
					t.Errorf("username mismatch: expected %q, got %q", u.UserName, got.UserName)
				}
			},
		},
		{
			name: "GetUserByID returns error for non-existent id",
			testFn: func(t *testing.T, db *sql.DB) {
				_, err := GetUserByID(9999, db)
				if err == nil {
					t.Fatalf("expected error when fetching non-existent user, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					// GetUserByID returns sql.ErrNoRows when not found
					t.Fatalf("expected sql.ErrNoRows, got: %v", err)
				}
			},
		},
		{
			name: "GetUserByID returns inserted user",
			testFn: func(t *testing.T, db *sql.DB) {
				u := model.User{
					UserName:       "Bob",
					CurrentAmount:  utils.Money(1500),
					MonthlyInputs:  utils.Money(2500),
					MonthlyOutputs: utils.Money(300),
				}
				uRet, err := CreateUser(u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}
				got, err := GetUserByID(uRet.ID, db)
				if err != nil {
					t.Fatalf("GetUserByID() returned error: %v", err)
				}
				if got.ID != uRet.ID {
					t.Errorf("id mismatch: expected %d, got %d", uRet.ID, got.ID)
				}
				if got.UserName != u.UserName {
					t.Errorf("username mismatch: expected %q, got %q", u.UserName, got.UserName)
				}
			},
		},
		{
			name: "DeleteUserByID deletes existing user and subsequent GetUserByID returns ErrNoRows",
			testFn: func(t *testing.T, db *sql.DB) {
				u := model.User{
					UserName:       "Charlie",
					CurrentAmount:  utils.Money(500),
					MonthlyInputs:  utils.Money(1000),
					MonthlyOutputs: utils.Money(200),
				}
				uRet, err := CreateUser(u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}
				deleted, err := DeleteUserByID(uRet.ID, db)
				if err != nil {
					t.Fatalf("DeleteUserByID() returned error: %v", err)
				}
				if deleted != 1 {
					t.Fatalf("expected 1 row deleted, got %d", deleted)
				}
				_, err = GetUserByID(uRet.ID, db)
				if err == nil {
					t.Fatalf("expected error when fetching deleted user, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows after delete, got: %v", err)
				}
			},
		},
		{
			name: "UpdateUserByID updates fields and returns updated record",
			testFn: func(t *testing.T, db *sql.DB) {
				u := model.User{
					UserName:       "Dana",
					CurrentAmount:  utils.Money(800),
					MonthlyInputs:  utils.Money(1200),
					MonthlyOutputs: utils.Money(100),
				}
				uRet, err := CreateUser(u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}

				updated := model.User{
					UserName:       "Dana Updated",
					CurrentAmount:  utils.Money(900),
					MonthlyInputs:  utils.Money(1300),
					MonthlyOutputs: utils.Money(150),
				}

				got, err := UpdateUserByID(uRet.ID, &updated, db)
				if err != nil {
					t.Fatalf("UpdateUserByID() returned error: %v", err)
				}
				if got.ID != uRet.ID {
					t.Errorf("id changed after update: expected %d, got %d", uRet.ID, got.ID)
				}
				if got.UserName != updated.UserName {
					t.Errorf("username not updated: expected %q, got %q", updated.UserName, got.UserName)
				}
				if got.CurrentAmount != updated.CurrentAmount {
					t.Errorf("current amount not updated: expected %v, got %v", updated.CurrentAmount, got.CurrentAmount)
				}
			},
		},
		{
			name: "UpdateUserByID returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, db *sql.DB) {
				updated := model.User{
					UserName:       "NonExistent",
					CurrentAmount:  utils.Money(0),
					MonthlyInputs:  utils.Money(0),
					MonthlyOutputs: utils.Money(0),
				}
				_, err := UpdateUserByID(9999, &updated, db)
				if err == nil {
					t.Fatalf("expected error when updating non-existent user, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows when updating non-existent user, got: %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := setupDB(t)
			defer teardown()
			tc.testFn(t, db)
		})
	}
}
