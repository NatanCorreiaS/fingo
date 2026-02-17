package dbsqlite

import (
	"database/sql"
	"errors"
	"testing"

	"natan/fingo/model"
	"natan/fingo/utils"
)

// Tests for goals CRUD operations using table-driven style.
// Relies on setupDB helper defined in users_sqlite_test.go.

func TestGoals_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T, db *sql.DB, userID int64)
	}{
		{
			name: "GetAllGoals on empty DB returns zero",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				goals, err := GetAllGoals(db)
				if err != nil {
					t.Fatalf("GetAllGoals() returned unexpected error: %v", err)
				}
				if got := len(goals); got != 0 {
					t.Fatalf("expected 0 goals, got %d", got)
				}
			},
		},
		{
			name: "CreateGoal inserts and GetAllGoals returns it",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Vacation",
					Desc:     "Trip to the beach",
					Price:    utils.Money(150000), // 1500.00
					Pros:     "relaxing",
					Cons:     "costly",
					UserID:   userID,
					Deadline: "2026-12-31",
				}
				ret, err := CreateGoal(g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				if ret == nil {
					t.Fatalf("CreateGoal() returned nil")
				}
				if ret.ID == 0 {
					t.Fatalf("expected inserted goal to have non-zero ID")
				}

				all, err := GetAllGoals(db)
				if err != nil {
					t.Fatalf("GetAllGoals() returned error after insert: %v", err)
				}
				if len(all) != 1 {
					t.Fatalf("expected 1 goal after insert, got %d", len(all))
				}
			},
		},
		{
			name: "CreateGoal accepts description with apostrophe",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Book",
					Desc:     "O'Reilly special edition",
					Price:    utils.Money(4999),
					Pros:     "informative",
					Cons:     "",
					UserID:   userID,
					Deadline: "2026-06-01",
				}
				ret, err := CreateGoal(g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error for apostrophe: %v", err)
				}
				got, err := GetGoalByID(ret.ID, db)
				if err != nil {
					t.Fatalf("GetGoalByID() returned error for inserted goal: %v", err)
				}
				if got.Desc != g.Desc {
					t.Errorf("description mismatch: want %q got %q", g.Desc, got.Desc)
				}
			},
		},
		{
			name: "GetGoalByID returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				_, err := GetGoalByID(999999, db)
				if err == nil {
					t.Fatalf("expected error when fetching non-existent goal, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows, got: %v", err)
				}
			},
		},
		{
			name: "GetGoalByID returns inserted goal",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Laptop",
					Desc:     "Work laptop",
					Price:    utils.Money(350000),
					Pros:     "fast",
					Cons:     "heavy",
					UserID:   userID,
					Deadline: "2026-09-30",
				}
				ret, err := CreateGoal(g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				got, err := GetGoalByID(ret.ID, db)
				if err != nil {
					t.Fatalf("GetGoalByID() returned error: %v", err)
				}
				if got.ID != ret.ID {
					t.Errorf("id mismatch: expected %d, got %d", ret.ID, got.ID)
				}
				if got.Name != g.Name {
					t.Errorf("name mismatch: expected %q, got %q", g.Name, got.Name)
				}
			},
		},
		{
			name: "UpdateGoalById updates fields and returns updated record",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Camera",
					Desc:     "Old camera",
					Price:    utils.Money(80000),
					Pros:     "portable",
					Cons:     "low-res",
					UserID:   userID,
					Deadline: "2026-04-01",
				}
				created, err := CreateGoal(g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}

				updated := model.Goal{
					Name:     "Camera Pro",
					Desc:     "Upgraded model",
					Price:    utils.Money(120000),
					Pros:     "high-res",
					Cons:     "expensive",
					UserID:   userID,
					Deadline: "2026-12-01",
				}

				got, err := UpdateGoalById(created.ID, updated, db)
				if err != nil {
					t.Fatalf("UpdateGoalById() returned error: %v", err)
				}
				if got.ID != created.ID {
					t.Errorf("id changed after update: expected %d, got %d", created.ID, got.ID)
				}
				if got.Name != updated.Name {
					t.Errorf("name not updated: expected %q, got %q", updated.Name, got.Name)
				}
				if got.Price != updated.Price {
					t.Errorf("price not updated: expected %v, got %v", updated.Price, got.Price)
				}
			},
		},
		{
			name: "UpdateGoalById returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				updated := model.Goal{
					Name:     "NonExistent",
					Desc:     "",
					Price:    utils.Money(0),
					Pros:     "",
					Cons:     "",
					UserID:   userID,
					Deadline: "2026-01-01",
				}
				_, err := UpdateGoalById(999999, updated, db)
				if err == nil {
					t.Fatalf("expected error when updating non-existent goal, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows when updating non-existent goal, got: %v", err)
				}
			},
		},
		{
			name: "DeleteGoalByID deletes existing goal and subsequent GetGoalByID returns ErrNoRows",
			testFn: func(t *testing.T, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "To Delete",
					Desc:     "temporary",
					Price:    utils.Money(1000),
					Pros:     "",
					Cons:     "",
					UserID:   userID,
					Deadline: "2026-02-02",
				}
				created, err := CreateGoal(g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				deleted, err := DeleteGoalByID(created.ID, db)
				if err != nil {
					t.Fatalf("DeleteGoalByID() returned error: %v", err)
				}
				if deleted != 1 {
					t.Fatalf("expected 1 row deleted, got %d", deleted)
				}
				_, err = GetGoalByID(created.ID, db)
				if err == nil {
					t.Fatalf("expected error when fetching deleted goal, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows after delete, got: %v", err)
				}
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := setupDB(t)
			defer teardown()

			// create a user to satisfy foreign key constraint
			u := model.User{
				UserName:       "goal-user",
				CurrentAmount:  utils.Money(0),
				MonthlyInputs:  utils.Money(0),
				MonthlyOutputs: utils.Money(0),
			}
			uRet, err := CreateUser(u, db)
			if err != nil {
				t.Fatalf("failed to create user for goals tests: %v", err)
			}

			tc.testFn(t, db, uRet.ID)
		})
	}
}
