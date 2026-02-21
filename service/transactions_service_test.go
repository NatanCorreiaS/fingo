package service

import (
	"testing"

	"natan/fingo/model"
)

// createTransactionForTests creates a basic transaction using the service layer.
// It creates a user first to satisfy the foreign key constraint.
func createTransactionForTests(t *testing.T) *model.Transaction {
	t.Helper()

	user, err := CreateUser(ctxTest, model.User{
		UserName: "tx-user-service",
	})
	if err != nil {
		t.Fatalf("failed to create user for transaction tests: %v", err)
	}

	tx, err := CreateTransaction(ctxTest, model.Transaction{
		Desc:   "Transaction Service Test",
		Amount: 0,
		IsDebt: false,
		UserID: user.ID,
	})
	if err != nil {
		t.Fatalf("failed to create transaction for tests: %v", err)
	}
	if tx == nil {
		t.Fatalf("createTransactionForTests returned nil transaction without error")
	}
	return tx
}

func TestGetTransactionByID(t *testing.T) {
	base := createTransactionForTests(t)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "existing_id",
			id:      base.ID,
			wantErr: false,
		},
		{
			name:    "non_existing_id",
			id:      base.ID + 999999,
			wantErr: true,
		},
		{
			name:    "zero_id",
			id:      0,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetTransactionByID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetTransactionByID() expected error, got nil; tx=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetTransactionByID() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("GetTransactionByID() returned nil without error")
			}
			if got.ID != tc.id {
				t.Errorf("GetTransactionByID() ID mismatch: got=%d want=%d", got.ID, tc.id)
			}
		})
	}
}

func TestGetAllTransactions(t *testing.T) {
	// Ensure at least one transaction exists
	_ = createTransactionForTests(t)

	tests := []struct {
		name        string
		wantMinSize int
	}{
		{
			name:        "at_least_one_transaction",
			wantMinSize: 1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			txs, err := GetAllTransactions(ctxTest)
			if err != nil {
				t.Fatalf("GetAllTransactions() unexpected error: %v", err)
			}
			if len(txs) < tc.wantMinSize {
				t.Errorf("GetAllTransactions() length < %d, got %d", tc.wantMinSize, len(txs))
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	// Create user once for all subtests
	user, err := CreateUser(ctxTest, model.User{
		UserName: "tx-create-service",
	})
	if err != nil {
		t.Fatalf("failed to create user for CreateTransaction tests: %v", err)
	}

	tests := []struct {
		name     string
		input    model.Transaction
		wantErr  bool
		wantDesc string
	}{
		{
			name: "valid_transaction",
			input: model.Transaction{
				Desc:   "Create Service Test",
				Amount: 100,
				IsDebt: false,
				UserID: user.ID,
			},
			wantErr:  false,
			wantDesc: "Create Service Test",
		},
		{
			name: "zero_amount",
			input: model.Transaction{
				Desc:   "Zero Amount",
				Amount: 0,
				IsDebt: false,
				UserID: user.ID,
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := CreateTransaction(ctxTest, tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("CreateTransaction() expected error, got nil; tx=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("CreateTransaction() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("CreateTransaction() returned nil without error")
			}
			if got.ID == 0 {
				t.Errorf("CreateTransaction() expected non-zero ID, got %d", got.ID)
			}
			if tc.wantDesc != "" && got.Desc != tc.wantDesc {
				t.Errorf("CreateTransaction() desc mismatch: got=%q want=%q", got.Desc, tc.wantDesc)
			}
		})
	}
}

func TestUpdateTransactionByID(t *testing.T) {
	base := createTransactionForTests(t)

	tests := []struct {
		name       string
		id         int64
		update     *model.TransactionUpdate
		wantErr    bool
		wantDetail string
	}{
		{
			name: "valid_update",
			id:   base.ID,
			update: &model.TransactionUpdate{
				Desc: strPtr("Updated Transaction Description"),
			},
			wantErr:    false,
			wantDetail: "Updated Transaction Description",
		},
		{
			name: "non_existing_id",
			id:   base.ID + 999999,
			update: &model.TransactionUpdate{
				Desc: strPtr("Won't Exist"),
			},
			wantErr: true,
		},
		{
			name:    "nil_update",
			id:      base.ID,
			update:  nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := UpdateTransactionByID(ctxTest, tc.id, tc.update)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("UpdateTransactionByID() expected error, got nil; tx=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("UpdateTransactionByID() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("UpdateTransactionByID() returned nil without error")
			}
			if tc.wantDetail != "" && got.Desc != tc.wantDetail {
				t.Errorf("UpdateTransactionByID() description mismatch: got=%q want=%q", got.Desc, tc.wantDetail)
			}
		})
	}
}

func TestDeleteTransactionByID(t *testing.T) {
	base := createTransactionForTests(t)

	tests := []struct {
		name     string
		id       int64
		wantErr  bool
		wantRows int64
	}{
		{
			name:     "delete_existing_transaction",
			id:       base.ID,
			wantErr:  false,
			wantRows: 1,
		},
		{
			name:     "delete_non_existing_transaction",
			id:       base.ID + 999999,
			wantErr:  false,
			wantRows: 0,
		},
		{
			name:     "delete_zero_id",
			id:       0,
			wantErr:  false,
			wantRows: 0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			rows, err := DeleteTransactionByID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("DeleteTransactionByID() expected error, got nil; rows=%d", rows)
				}
				return
			}

			if err != nil {
				t.Fatalf("DeleteTransactionByID() unexpected error: %v", err)
			}
			if rows != tc.wantRows {
				t.Errorf("DeleteTransactionByID() rows mismatch: got=%d want=%d", rows, tc.wantRows)
			}
		})
	}
}
