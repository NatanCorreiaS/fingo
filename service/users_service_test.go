package service

import (
	"context"
	"testing"

	"natan/fingo/model"
)

var ctxTest = context.Background()

func strPtr(s string) *string {
	return &s
}

// ---- Users service tests ----

func TestCreateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       model.User
		wantErr     bool
		wantName    string
		wantHasUser bool
	}{
		{
			name: "valid_user",
			input: model.User{
				UserName: "User Service Test",
			},
			wantErr:     false,
			wantName:    "User Service Test",
			wantHasUser: true,
		},
		{
			name: "empty_username_still_ok",
			input: model.User{
				UserName: "",
			},
			wantErr:     false,
			wantHasUser: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := CreateUser(ctxTest, tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("CreateUser() expected error, got nil; user=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("CreateUser() unexpected error: %v", err)
			}

			if tc.wantHasUser && got == nil {
				t.Fatalf("CreateUser() returned nil user without error")
			}
			if !tc.wantHasUser && got != nil {
				t.Fatalf("CreateUser() returned non-nil user when not expected: %+v", got)
			}

			if got != nil && got.ID == 0 {
				t.Errorf("CreateUser() expected non-zero ID, got %d", got.ID)
			}

			if tc.wantName != "" && got != nil && got.UserName != tc.wantName {
				t.Errorf("CreateUser() username mismatch: got=%q want=%q", got.UserName, tc.wantName)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	t.Parallel()

	base, err := CreateUser(ctxTest, model.User{
		UserName: "GetUserByID Base",
	})
	if err != nil {
		t.Fatalf("failed to create base user: %v", err)
	}

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
			got, err := GetUserByID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetUserByID() expected error, got nil; user=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetUserByID() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("GetUserByID() returned nil without error")
			}
			if got.ID != tc.id {
				t.Errorf("GetUserByID() ID mismatch: got=%d want=%d", got.ID, tc.id)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	t.Parallel()

	created, err := CreateUser(ctxTest, model.User{
		UserName: "GetAllUsers Base",
	})
	if err != nil {
		t.Fatalf("failed to create base user: %v", err)
	}
	_ = created

	tests := []struct {
		name        string
		wantMinSize int
	}{
		{
			name:        "at_least_one_user",
			wantMinSize: 1,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			users, err := GetAllUsers(ctxTest)
			if err != nil {
				t.Fatalf("GetAllUsers() unexpected error: %v", err)
			}
			if len(users) < tc.wantMinSize {
				t.Errorf("GetAllUsers() length < %d, got %d", tc.wantMinSize, len(users))
			}
		})
	}
}

func TestUpdateUserByID(t *testing.T) {
	t.Parallel()

	base, err := CreateUser(ctxTest, model.User{
		UserName: "UpdateUser Base",
	})
	if err != nil {
		t.Fatalf("failed to create base user: %v", err)
	}

	tests := []struct {
		name     string
		id       int64
		update   *model.UserUpdate
		wantErr  bool
		wantName string
	}{
		{
			name: "valid_update",
			id:   base.ID,
			update: &model.UserUpdate{
				UserName: strPtr("updated_user"),
			},
			wantErr:  false,
			wantName: "updated_user",
		},
		{
			name: "non_existing_id",
			id:   base.ID + 999999,
			update: &model.UserUpdate{
				UserName: strPtr("nonexistent"),
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
			got, err := UpdateUserByID(ctxTest, tc.id, tc.update)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("UpdateUserByID() expected error, got nil; user=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("UpdateUserByID() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("UpdateUserByID() returned nil without error")
			}
			if tc.wantName != "" && got.UserName != tc.wantName {
				t.Errorf("UpdateUserByID() username mismatch: got=%q want=%q", got.UserName, tc.wantName)
			}
		})
	}
}

func TestDeleteUserByID(t *testing.T) {
	t.Parallel()

	base, err := CreateUser(ctxTest, model.User{
		UserName: "DeleteUser Base",
	})
	if err != nil {
		t.Fatalf("failed to create base user: %v", err)
	}

	tests := []struct {
		name     string
		id       int64
		wantErr  bool
		wantRows int64
	}{
		{
			name:     "delete_existing_user",
			id:       base.ID,
			wantErr:  false,
			wantRows: 1,
		},
		{
			name:     "delete_non_existing_user",
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
			rows, err := DeleteUserByID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("DeleteUserByID() expected error, got nil; rows=%d", rows)
				}
				return
			}

			if err != nil {
				t.Fatalf("DeleteUserByID() unexpected error: %v", err)
			}
			if rows != tc.wantRows {
				t.Errorf("DeleteUserByID() rows mismatch: got=%d want=%d", rows, tc.wantRows)
			}
		})
	}
}

// ---- GetAllTransactionsByUserID service tests ----

func TestGetAllTransactionsByUserID(t *testing.T) {
	t.Parallel()

	user, err := CreateUser(ctxTest, model.User{
		UserName: "tx-by-user-service",
	})
	if err != nil {
		t.Fatalf("failed to create user for GetAllTransactionsByUserID tests: %v", err)
	}

	_, err = CreateTransaction(ctxTest, model.Transaction{
		Desc:   "TX for user 1",
		Amount: 100,
		IsDebt: false,
		UserID: user.ID,
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	_, err = CreateTransaction(ctxTest, model.Transaction{
		Desc:   "TX for user 2",
		Amount: 200,
		IsDebt: true,
		UserID: user.ID,
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	tests := []struct {
		name        string
		id          int64
		wantMinSize int
		wantErr     bool
	}{
		{
			name:        "existing_user_with_transactions",
			id:          user.ID,
			wantMinSize: 2,
			wantErr:     false,
		},
		{
			name:        "non_existing_user_returns_empty",
			id:          user.ID + 999999,
			wantMinSize: 0,
			wantErr:     false,
		},
		{
			name:        "zero_id_returns_empty",
			id:          0,
			wantMinSize: 0,
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			txs, err := GetAllTransactionsByUserID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetAllTransactionsByUserID() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetAllTransactionsByUserID() unexpected error: %v", err)
			}
			if len(txs) < tc.wantMinSize {
				t.Errorf("GetAllTransactionsByUserID() length < %d, got %d", tc.wantMinSize, len(txs))
			}
			for _, tx := range txs {
				if tx.UserID != tc.id {
					t.Errorf("GetAllTransactionsByUserID() returned transaction with UserID=%d, want %d", tx.UserID, tc.id)
				}
			}
		})
	}
}

func TestGetAllTransactionsByUserID_Isolation(t *testing.T) {
	t.Parallel()

	user1, err := CreateUser(ctxTest, model.User{UserName: "tx-isolation-user1"})
	if err != nil {
		t.Fatalf("failed to create user1: %v", err)
	}
	user2, err := CreateUser(ctxTest, model.User{UserName: "tx-isolation-user2"})
	if err != nil {
		t.Fatalf("failed to create user2: %v", err)
	}

	_, err = CreateTransaction(ctxTest, model.Transaction{
		Desc: "User1 TX", Amount: 500, IsDebt: false, UserID: user1.ID,
	})
	if err != nil {
		t.Fatalf("failed to create transaction for user1: %v", err)
	}
	_, err = CreateTransaction(ctxTest, model.Transaction{
		Desc: "User2 TX", Amount: 300, IsDebt: true, UserID: user2.ID,
	})
	if err != nil {
		t.Fatalf("failed to create transaction for user2: %v", err)
	}

	txs, err := GetAllTransactionsByUserID(ctxTest, user1.ID)
	if err != nil {
		t.Fatalf("GetAllTransactionsByUserID() unexpected error: %v", err)
	}
	for _, tx := range txs {
		if tx.UserID != user1.ID {
			t.Errorf("expected all transactions to belong to user1 (ID=%d), got UserID=%d", user1.ID, tx.UserID)
		}
	}
}

// ---- GetAllGoalsByUserID service tests ----

func TestGetAllGoalsByUserID(t *testing.T) {
	t.Parallel()

	user, err := CreateUser(ctxTest, model.User{
		UserName: "goals-by-user-service",
	})
	if err != nil {
		t.Fatalf("failed to create user for GetAllGoalsByUserID tests: %v", err)
	}

	_, err = CreateGoal(ctxTest, model.Goal{
		Name:     "Goal A",
		Price:    100,
		UserID:   user.ID,
		Deadline: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("failed to create goal: %v", err)
	}

	_, err = CreateGoal(ctxTest, model.Goal{
		Name:     "Goal B",
		Price:    200,
		UserID:   user.ID,
		Deadline: "2026-06-01",
	})
	if err != nil {
		t.Fatalf("failed to create goal: %v", err)
	}

	tests := []struct {
		name        string
		id          int64
		wantMinSize int
		wantErr     bool
	}{
		{
			name:        "existing_user_with_goals",
			id:          user.ID,
			wantMinSize: 2,
			wantErr:     false,
		},
		{
			name:        "non_existing_user_returns_empty",
			id:          user.ID + 999999,
			wantMinSize: 0,
			wantErr:     false,
		},
		{
			name:        "zero_id_returns_empty",
			id:          0,
			wantMinSize: 0,
			wantErr:     false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			goals, err := GetAllGoalsByUserID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetAllGoalsByUserID() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetAllGoalsByUserID() unexpected error: %v", err)
			}
			if len(goals) < tc.wantMinSize {
				t.Errorf("GetAllGoalsByUserID() length < %d, got %d", tc.wantMinSize, len(goals))
			}
			for _, g := range goals {
				if g.UserID != tc.id {
					t.Errorf("GetAllGoalsByUserID() returned goal with UserID=%d, want %d", g.UserID, tc.id)
				}
			}
		})
	}
}

func TestGetAllGoalsByUserID_Isolation(t *testing.T) {
	t.Parallel()

	user1, err := CreateUser(ctxTest, model.User{UserName: "goal-isolation-user1"})
	if err != nil {
		t.Fatalf("failed to create user1: %v", err)
	}
	user2, err := CreateUser(ctxTest, model.User{UserName: "goal-isolation-user2"})
	if err != nil {
		t.Fatalf("failed to create user2: %v", err)
	}

	_, err = CreateGoal(ctxTest, model.Goal{
		Name: "User1 Goal", Price: 500, UserID: user1.ID, Deadline: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("failed to create goal for user1: %v", err)
	}
	_, err = CreateGoal(ctxTest, model.Goal{
		Name: "User2 Goal", Price: 300, UserID: user2.ID, Deadline: "2026-01-01",
	})
	if err != nil {
		t.Fatalf("failed to create goal for user2: %v", err)
	}

	goals, err := GetAllGoalsByUserID(ctxTest, user1.ID)
	if err != nil {
		t.Fatalf("GetAllGoalsByUserID() unexpected error: %v", err)
	}
	for _, g := range goals {
		if g.UserID != user1.ID {
			t.Errorf("expected all goals to belong to user1 (ID=%d), got UserID=%d", user1.ID, g.UserID)
		}
	}
}

// ---- Helpers ----
