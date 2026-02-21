package service

import (
	"testing"

	"natan/fingo/model"
)

func TestGoalsService_GetGoalByID(t *testing.T) {
	user, err := CreateUser(ctxTest, model.User{
		UserName: "goal-user-service",
	})
	if err != nil {
		t.Fatalf("failed to create user for goals tests: %v", err)
	}

	base, err := CreateGoal(ctxTest, model.Goal{
		Name:     "GetGoalByID Base",
		Price:    0,
		UserID:   user.ID,
		Deadline: "",
	})
	if err != nil {
		t.Fatalf("failed to create base goal: %v", err)
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
			got, err := GetGoalByID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("GetGoalByID() expected error, got nil; goal=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("GetGoalByID() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("GetGoalByID() returned nil without error")
			}
			if got.ID != tc.id {
				t.Errorf("GetGoalByID() ID mismatch: got=%d want=%d", got.ID, tc.id)
			}
		})
	}
}

func TestGoalsService_GetAllGoals(t *testing.T) {
	user, err := CreateUser(ctxTest, model.User{
		UserName: "goal-user-service-all",
	})
	if err != nil {
		t.Fatalf("failed to create user for GetAllGoals tests: %v", err)
	}

	_, err = CreateGoal(ctxTest, model.Goal{
		Name:     "GetAllGoals Base",
		Price:    0,
		UserID:   user.ID,
		Deadline: "",
	})
	if err != nil {
		t.Fatalf("failed to create base goal: %v", err)
	}

	tests := []struct {
		name        string
		wantMinSize int
	}{
		{
			name:        "at_least_one_goal",
			wantMinSize: 1,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			goals, err := GetAllGoals(ctxTest)
			if err != nil {
				t.Fatalf("GetAllGoals() unexpected error: %v", err)
			}
			if len(goals) < tc.wantMinSize {
				t.Errorf("GetAllGoals() length < %d, got %d", tc.wantMinSize, len(goals))
			}
		})
	}
}

func TestGoalsService_UpdateGoalByID(t *testing.T) {
	user, err := CreateUser(ctxTest, model.User{
		UserName: "goal-user-service-update",
	})
	if err != nil {
		t.Fatalf("failed to create user for UpdateGoal tests: %v", err)
	}

	base, err := CreateGoal(ctxTest, model.Goal{
		Name:     "UpdateGoal Base",
		Price:    0,
		UserID:   user.ID,
		Deadline: "",
	})
	if err != nil {
		t.Fatalf("failed to create base goal: %v", err)
	}

	tests := []struct {
		name        string
		id          int64
		update      *model.GoalUpdate
		wantErr     bool
		wantNewName string
	}{
		{
			name: "valid_update",
			id:   base.ID,
			update: &model.GoalUpdate{
				Name: strPtr("Updated Goal Name"),
			},
			wantErr:     false,
			wantNewName: "Updated Goal Name",
		},
		{
			name: "non_existing_id",
			id:   base.ID + 999999,
			update: &model.GoalUpdate{
				Name: strPtr("Won't Exist"),
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
			got, err := UpdateGoalByID(ctxTest, tc.id, tc.update)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("UpdateGoalByID() expected error, got nil; goal=%+v", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("UpdateGoalByID() unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("UpdateGoalByID() returned nil without error")
			}
			if tc.wantNewName != "" && got.Name != tc.wantNewName {
				t.Errorf("UpdateGoalByID() name mismatch: got=%q want=%q", got.Name, tc.wantNewName)
			}
		})
	}
}

func TestGoalsService_DeleteGoalByID(t *testing.T) {
	user, err := CreateUser(ctxTest, model.User{
		UserName: "goal-user-service-delete",
	})
	if err != nil {
		t.Fatalf("failed to create user for DeleteGoal tests: %v", err)
	}

	base, err := CreateGoal(ctxTest, model.Goal{
		Name:     "DeleteGoal Base",
		Price:    0,
		UserID:   user.ID,
		Deadline: "",
	})
	if err != nil {
		t.Fatalf("failed to create base goal: %v", err)
	}

	tests := []struct {
		name     string
		id       int64
		wantErr  bool
		wantRows int64
	}{
		{
			name:     "delete_existing_goal",
			id:       base.ID,
			wantErr:  false,
			wantRows: 1,
		},
		{
			name:     "delete_non_existing_goal",
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
			rows, err := DeleteGoalByID(ctxTest, tc.id)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("DeleteGoalByID() expected error, got nil; rows=%d", rows)
				}
				return
			}

			if err != nil {
				t.Fatalf("DeleteGoalByID() unexpected error: %v", err)
			}
			if rows != tc.wantRows {
				t.Errorf("DeleteGoalByID() rows mismatch: got=%d want=%d", rows, tc.wantRows)
			}
		})
	}
}
