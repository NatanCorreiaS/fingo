package dbsqlite

import (
	"database/sql"
	"errors"
	"testing"

	"natan/fingo/model"
	"natan/fingo/utils"
)

// Tests for transaction CRUD operations in a table-driven manner.
// These rely on the existing setupDB helper defined in users_sqlite_test.go.

func TestTransactions_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T, db *sql.DB, userID int64)
	}{
		{
			name: "GetAllTransactions on empty DB returns zero",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				txns, err := GetAllTransactions(db)
				if err != nil {
					t.Fatalf("GetAllTransactions() returned unexpected error: %v", err)
				}
				if got := len(txns); got != 0 {
					t.Fatalf("expected 0 transactions, got %d", got)
				}
			},
		},
		{
			name: "CreateTransaction inserts and GetAllTransactions returns it",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Salary",
					Amount: utils.Money(500000), // 5000.00
					IsDebt: false,
					UserID: userID,
				}
				ret, err := CreateTransaction(tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}
				if ret == nil {
					t.Fatalf("CreateTransaction returned nil")
				}
				if ret.ID == 0 {
					t.Fatalf("expected inserted transaction to have non-zero ID")
				}

				list, err := GetAllTransactions(db)
				if err != nil {
					t.Fatalf("GetAllTransactions() returned error after insert: %v", err)
				}
				if len(list) != 1 {
					t.Fatalf("expected 1 transaction after insert, got %d", len(list))
				}
			},
		},
		{
			name: "CreateTransaction accepts description with apostrophe",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "O'Reilly bonus",
					Amount: utils.Money(12345),
					IsDebt: false,
					UserID: userID,
				}
				ret, err := CreateTransaction(tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() error for apostrophe: %v", err)
				}
				got, err := GetTransactionByID(ret.ID, db)
				if err != nil {
					t.Fatalf("GetTransactionByID() returned error for inserted transaction: %v", err)
				}
				if got.Desc != tr.Desc {
					t.Errorf("description mismatch: want %q got %q", tr.Desc, got.Desc)
				}
			},
		},
		{
			name: "GetTransactionByID returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				_, err := GetTransactionByID(999999, db)
				if err == nil {
					t.Fatalf("expected error when fetching non-existent transaction, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows, got: %v", err)
				}
			},
		},
		{
			name: "GetTransactionByID returns inserted transaction",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Groceries",
					Amount: utils.Money(2500),
					IsDebt: false,
					UserID: userID,
				}
				ret, err := CreateTransaction(tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}
				got, err := GetTransactionByID(ret.ID, db)
				if err != nil {
					t.Fatalf("GetTransactionByID() returned error: %v", err)
				}
				if got.ID != ret.ID {
					t.Errorf("id mismatch: expected %d, got %d", ret.ID, got.ID)
				}
				if got.Desc != tr.Desc {
					t.Errorf("description mismatch: expected %q, got %q", tr.Desc, got.Desc)
				}
			},
		},
		{
			name: "UpdateTransactionByID updates fields and returns updated record",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Old Desc",
					Amount: utils.Money(1000),
					IsDebt: true,
					UserID: userID,
				}
				created, err := CreateTransaction(tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}

				updated := model.Transaction{
					Desc:   "New Desc",
					Amount: utils.Money(2000),
					IsDebt: false,
					UserID: userID,
				}

				got, err := UpdateTransactionByID(created.ID, &updated, db)
				if err != nil {
					t.Fatalf("UpdateTransactionByID() returned error: %v", err)
				}
				if got.ID != created.ID {
					t.Errorf("id changed after update: expected %d, got %d", created.ID, got.ID)
				}
				if got.Desc != updated.Desc {
					t.Errorf("description not updated: expected %q, got %q", updated.Desc, got.Desc)
				}
				if got.Amount != updated.Amount {
					t.Errorf("amount not updated: expected %v, got %v", updated.Amount, got.Amount)
				}
				if got.IsDebt != updated.IsDebt {
					t.Errorf("is_debt not updated: expected %v, got %v", updated.IsDebt, got.IsDebt)
				}
			},
		},
		{
			name: "UpdateTransactionByID returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				updated := model.Transaction{
					Desc:   "Doesn't matter",
					Amount: utils.Money(0),
					IsDebt: false,
					UserID: userID,
				}
				_, err := UpdateTransactionByID(999999, &updated, db)
				if err == nil {
					t.Fatalf("expected error when updating non-existent transaction, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows when updating non-existent transaction, got: %v", err)
				}
			},
		},
		{
			name: "DeleteTransactionByID deletes existing transaction and subsequent GetTransactionByID returns ErrNoRows",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "To be deleted",
					Amount: utils.Money(777),
					IsDebt: true,
					UserID: userID,
				}
				created, err := CreateTransaction(tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}
				deleted, err := DeleteTransactionByID(created.ID, db)
				if err != nil {
					t.Fatalf("DeleteTransactionByID() returned error: %v", err)
				}
				if deleted != 1 {
					t.Fatalf("expected 1 row deleted, got %d", deleted)
				}
				_, err = GetTransactionByID(created.ID, db)
				if err == nil {
					t.Fatalf("expected error when fetching deleted transaction, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows after delete, got: %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range var
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := setupDB(t)
			defer teardown()

			// create a user to satisfy the foreign key constraint
			user := model.User{
				UserName:       "tx-user",
				CurrentAmount:  utils.Money(0),
				MonthlyInputs:  utils.Money(0),
				MonthlyOutputs: utils.Money(0),
			}
			uRet, err := CreateUser(user, db)
			if err != nil {
				t.Fatalf("failed to create user for transactions tests: %v", err)
			}

			tc.testFn(t, db, uRet.ID)
		})
	}
}
