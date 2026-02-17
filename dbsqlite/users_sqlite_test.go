package dbsqlite

import (
	"context"
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
		testFn func(t *testing.T, ctx context.Context, db *sql.DB)
	}{
		{
			name: "GetAllUsers on empty DB returns zero users",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				users, err := GetAllUsers(ctx, db)
				if err != nil {
					t.Fatalf("GetAllUsers() returned unexpected error on empty DB: %v", err)
				}
				if got := len(users); got != 0 {
					t.Fatalf("expected 0 users, got %d", got)
				}
			},
		},
		{
			name: "CreateUser inserts user and GetAllUsers returns it",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				u := model.User{
					UserName:       "Alice",
					CurrentAmount:  utils.Money(1000),
					MonthlyInputs:  utils.Money(2000),
					MonthlyOutputs: utils.Money(500),
				}

				uRet, err := CreateUser(ctx, u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}
				if uRet == nil {
					t.Fatalf("CreateUser() returned nil user")
				}
				if uRet.UserName != u.UserName {
					t.Errorf("returned username mismatch: expected %q, got %q", u.UserName, uRet.UserName)
				}
				// verify via GetAllUsers as well
				users, err := GetAllUsers(ctx, db)
				if err != nil {
					t.Fatalf("GetAllUsers() returned error after insert: %v", err)
				}
				if len(users) != 1 {
					t.Fatalf("expected 1 user, got %d", len(users))
				}
			},
		},
		{
			name: "CreateUser accepts username with apostrophe (safe parameterization)",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				u := model.User{
					UserName:       "O'Neil",
					CurrentAmount:  utils.Money(0),
					MonthlyInputs:  utils.Money(0),
					MonthlyOutputs: utils.Money(0),
				}

				uRet, err := CreateUser(ctx, u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned unexpected error for username with apostrophe: %v", err)
				}
				if uRet == nil {
					t.Fatalf("CreateUser() returned nil user")
				}

				// verify retrieval
				got, err := GetUserByID(ctx, uRet.ID, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				_, err := GetUserByID(ctx, 9999, db)
				if err == nil {
					t.Fatalf("expected error when fetching non-existent user, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows, got: %v", err)
				}
			},
		},
		{
			name: "GetUserByID returns inserted user",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				u := model.User{
					UserName:       "Bob",
					CurrentAmount:  utils.Money(1500),
					MonthlyInputs:  utils.Money(2500),
					MonthlyOutputs: utils.Money(300),
				}
				uRet, err := CreateUser(ctx, u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}
				got, err := GetUserByID(ctx, uRet.ID, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				u := model.User{
					UserName:       "Charlie",
					CurrentAmount:  utils.Money(500),
					MonthlyInputs:  utils.Money(1000),
					MonthlyOutputs: utils.Money(200),
				}
				uRet, err := CreateUser(ctx, u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}
				deleted, err := DeleteUserByID(ctx, uRet.ID, db)
				if err != nil {
					t.Fatalf("DeleteUserByID() returned error: %v", err)
				}
				if deleted != 1 {
					t.Fatalf("expected 1 row deleted, got %d", deleted)
				}
				_, err = GetUserByID(ctx, uRet.ID, db)
				if err == nil {
					t.Fatalf("expected error when fetching deleted user, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows after delete, got: %v", err)
				}
			},
		},
		{
			name: "UpdateUserPartialByID updates only provided fields",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				u := model.User{
					UserName:       "Eve",
					CurrentAmount:  utils.Money(1000),
					MonthlyInputs:  utils.Money(2000),
					MonthlyOutputs: utils.Money(500),
				}
				uRet, err := CreateUser(ctx, u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}

				// Update only username
				newName := "Eve Updated"
				update := &model.UserUpdate{
					UserName: &newName,
				}

				got, err := UpdateUserPartialByID(ctx, uRet.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateUserPartialByID() returned error: %v", err)
				}

				// Verify username was updated
				if got.UserName != newName {
					t.Errorf("username not updated: expected %q, got %q", newName, got.UserName)
				}

				// Verify other fields remain unchanged
				if got.CurrentAmount != u.CurrentAmount {
					t.Errorf("current amount should not change: expected %v, got %v", u.CurrentAmount, got.CurrentAmount)
				}
				if got.MonthlyInputs != u.MonthlyInputs {
					t.Errorf("monthly inputs should not change: expected %v, got %v", u.MonthlyInputs, got.MonthlyInputs)
				}
				if got.MonthlyOutputs != u.MonthlyOutputs {
					t.Errorf("monthly outputs should not change: expected %v, got %v", u.MonthlyOutputs, got.MonthlyOutputs)
				}
			},
		},
		{
			name: "UpdateUserPartialByID updates multiple fields",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				u := model.User{
					UserName:       "Frank",
					CurrentAmount:  utils.Money(1000),
					MonthlyInputs:  utils.Money(2000),
					MonthlyOutputs: utils.Money(500),
				}
				uRet, err := CreateUser(ctx, u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}

				// Update username and current amount
				newName := "Frank Updated"
				newAmount := utils.Money(1500)
				update := &model.UserUpdate{
					UserName:      &newName,
					CurrentAmount: &newAmount,
				}

				got, err := UpdateUserPartialByID(ctx, uRet.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateUserPartialByID() returned error: %v", err)
				}

				// Verify updated fields
				if got.UserName != newName {
					t.Errorf("username not updated: expected %q, got %q", newName, got.UserName)
				}
				if got.CurrentAmount != newAmount {
					t.Errorf("current amount not updated: expected %v, got %v", newAmount, got.CurrentAmount)
				}

				// Verify unchanged fields
				if got.MonthlyInputs != u.MonthlyInputs {
					t.Errorf("monthly inputs should not change: expected %v, got %v", u.MonthlyInputs, got.MonthlyInputs)
				}
				if got.MonthlyOutputs != u.MonthlyOutputs {
					t.Errorf("monthly outputs should not change: expected %v, got %v", u.MonthlyOutputs, got.MonthlyOutputs)
				}
			},
		},
		{
			name: "UpdateUserPartialByID with no fields provided returns current user",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				u := model.User{
					UserName:       "Grace",
					CurrentAmount:  utils.Money(1000),
					MonthlyInputs:  utils.Money(2000),
					MonthlyOutputs: utils.Money(500),
				}
				uRet, err := CreateUser(ctx, u, db)
				if err != nil {
					t.Fatalf("CreateUser() returned error: %v", err)
				}

				// Update with no fields
				update := &model.UserUpdate{}

				got, err := UpdateUserPartialByID(ctx, uRet.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateUserPartialByID() returned error: %v", err)
				}

				// Verify all fields remain the same
				if got.UserName != u.UserName {
					t.Errorf("username should not change: expected %q, got %q", u.UserName, got.UserName)
				}
				if got.CurrentAmount != u.CurrentAmount {
					t.Errorf("current amount should not change: expected %v, got %v", u.CurrentAmount, got.CurrentAmount)
				}
				if got.MonthlyInputs != u.MonthlyInputs {
					t.Errorf("monthly inputs should not change: expected %v, got %v", u.MonthlyInputs, got.MonthlyInputs)
				}
				if got.MonthlyOutputs != u.MonthlyOutputs {
					t.Errorf("monthly outputs should not change: expected %v, got %v", u.MonthlyOutputs, got.MonthlyOutputs)
				}
			},
		},
		{
			name: "UpdateUserPartialByID returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB) {
				newName := "NonExistent"
				update := &model.UserUpdate{
					UserName: &newName,
				}
				_, err := UpdateUserPartialByID(ctx, 9999, update, db)
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
			ctx := context.Background()
			tc.testFn(t, ctx, db)
		})
	}
}
