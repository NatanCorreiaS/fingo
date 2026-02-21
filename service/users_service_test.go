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

// ---- Helpers ----
