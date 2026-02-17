package dbsqlite

import (
	"context"
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
		testFn func(t *testing.T, ctx context.Context, db *sql.DB, userID int64)
	}{
		{
			name: "GetAllTransactions on empty DB returns zero",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				txns, err := GetAllTransactions(ctx, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Salary",
					Amount: utils.Money(500000),
					IsDebt: false,
					UserID: userID,
				}
				ret, err := CreateTransaction(ctx, tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}
				if ret == nil {
					t.Fatalf("CreateTransaction returned nil")
				}
				if ret.ID == 0 {
					t.Fatalf("expected inserted transaction to have non-zero ID")
				}

				list, err := GetAllTransactions(ctx, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "O'Reilly bonus",
					Amount: utils.Money(12345),
					IsDebt: false,
					UserID: userID,
				}
				ret, err := CreateTransaction(ctx, tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() error for apostrophe: %v", err)
				}
				got, err := GetTransactionByID(ctx, ret.ID, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				_, err := GetTransactionByID(ctx, 999999, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Groceries",
					Amount: utils.Money(2500),
					IsDebt: false,
					UserID: userID,
				}
				ret, err := CreateTransaction(ctx, tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}
				got, err := GetTransactionByID(ctx, ret.ID, db)
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
			name: "UpdateTransactionPartialByID updates only provided fields",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Old Desc",
					Amount: utils.Money(1000),
					IsDebt: true,
					UserID: userID,
				}
				created, err := CreateTransaction(ctx, tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}

				newDesc := "New Desc"
				update := &model.TransactionUpdate{
					Desc: &newDesc,
				}

				got, err := UpdateTransactionPartialByID(ctx, created.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateTransactionPartialByID() returned error: %v", err)
				}
				if got.ID != created.ID {
					t.Errorf("id changed after update: expected %d, got %d", created.ID, got.ID)
				}
				if got.Desc != newDesc {
					t.Errorf("description not updated: expected %q, got %q", newDesc, got.Desc)
				}
				// Verify non-updated fields remain unchanged
				if got.Amount != tr.Amount {
					t.Errorf("amount should not change: expected %v, got %v", tr.Amount, got.Amount)
				}
				if got.IsDebt != tr.IsDebt {
					t.Errorf("is_debt should not change: expected %v, got %v", tr.IsDebt, got.IsDebt)
				}
			},
		},
		{
			name: "UpdateTransactionPartialByID updates multiple fields",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Initial",
					Amount: utils.Money(1000),
					IsDebt: true,
					UserID: userID,
				}
				created, err := CreateTransaction(ctx, tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}

				newDesc := "Updated"
				newAmount := utils.Money(2000)
				update := &model.TransactionUpdate{
					Desc:   &newDesc,
					Amount: &newAmount,
				}

				got, err := UpdateTransactionPartialByID(ctx, created.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateTransactionPartialByID() returned error: %v", err)
				}
				if got.Desc != newDesc {
					t.Errorf("description not updated: expected %q, got %q", newDesc, got.Desc)
				}
				if got.Amount != newAmount {
					t.Errorf("amount not updated: expected %v, got %v", newAmount, got.Amount)
				}
				if got.IsDebt != tr.IsDebt {
					t.Errorf("is_debt should not change: expected %v, got %v", tr.IsDebt, got.IsDebt)
				}
			},
		},
		{
			name: "UpdateTransactionPartialByID with no fields returns transaction unchanged",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "Original",
					Amount: utils.Money(5000),
					IsDebt: false,
					UserID: userID,
				}
				created, err := CreateTransaction(ctx, tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}

				update := &model.TransactionUpdate{}
				got, err := UpdateTransactionPartialByID(ctx, created.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateTransactionPartialByID() returned error: %v", err)
				}

				if got.Desc != tr.Desc {
					t.Errorf("description changed unexpectedly: expected %q, got %q", tr.Desc, got.Desc)
				}
				if got.Amount != tr.Amount {
					t.Errorf("amount changed unexpectedly: expected %v, got %v", tr.Amount, got.Amount)
				}
				if got.IsDebt != tr.IsDebt {
					t.Errorf("is_debt changed unexpectedly: expected %v, got %v", tr.IsDebt, got.IsDebt)
				}
			},
		},
		{
			name: "UpdateTransactionPartialByID returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				newDesc := "Ghost"
				update := &model.TransactionUpdate{
					Desc: &newDesc,
				}
				_, err := UpdateTransactionPartialByID(ctx, 999999, update, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				tr := model.Transaction{
					Desc:   "To be deleted",
					Amount: utils.Money(777),
					IsDebt: true,
					UserID: userID,
				}
				created, err := CreateTransaction(ctx, tr, db)
				if err != nil {
					t.Fatalf("CreateTransaction() returned error: %v", err)
				}
				deleted, err := DeleteTransactionByID(ctx, created.ID, db)
				if err != nil {
					t.Fatalf("DeleteTransactionByID() returned error: %v", err)
				}
				if deleted != 1 {
					t.Fatalf("expected 1 row deleted, got %d", deleted)
				}
				_, err = GetTransactionByID(ctx, created.ID, db)
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := setupDB(t)
			defer teardown()
			ctx := context.Background()

			user := model.User{
				UserName:       "tx-user",
				CurrentAmount:  utils.Money(0),
				MonthlyInputs:  utils.Money(0),
				MonthlyOutputs: utils.Money(0),
			}
			uRet, err := CreateUser(ctx, user, db)
			if err != nil {
				t.Fatalf("failed to create user for transactions tests: %v", err)
			}

			tc.testFn(t, ctx, db, uRet.ID)
		})
	}
}
