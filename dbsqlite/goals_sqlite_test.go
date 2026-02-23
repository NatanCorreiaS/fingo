package dbsqlite

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"natan/fingo/model"
	"natan/fingo/utils"
)

func TestGoals_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		testFn func(t *testing.T, ctx context.Context, db *sql.DB, userID int64)
	}{
		{
			name: "GetAllGoals on empty DB returns zero",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				goals, err := GetAllGoals(ctx, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Vacation",
					Desc:     "Trip to the beach",
					Price:    utils.Money(150000),
					Pros:     "relaxing",
					Cons:     "costly",
					UserID:   userID,
					Deadline: "2026-12-31",
				}
				ret, err := CreateGoal(ctx, g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				if ret == nil {
					t.Fatalf("CreateGoal() returned nil")
				}
				if ret.ID == 0 {
					t.Fatalf("expected inserted goal to have non-zero ID")
				}

				all, err := GetAllGoals(ctx, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Book",
					Desc:     "O'Reilly special edition",
					Price:    utils.Money(4999),
					Pros:     "informative",
					Cons:     "",
					UserID:   userID,
					Deadline: "2026-06-01",
				}
				ret, err := CreateGoal(ctx, g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error for apostrophe: %v", err)
				}
				got, err := GetGoalByID(ctx, ret.ID, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				_, err := GetGoalByID(ctx, 999999, db)
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
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Laptop",
					Desc:     "Work laptop",
					Price:    utils.Money(350000),
					Pros:     "fast",
					Cons:     "heavy",
					UserID:   userID,
					Deadline: "2026-09-30",
				}
				ret, err := CreateGoal(ctx, g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				got, err := GetGoalByID(ctx, ret.ID, db)
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
			name: "UpdateGoalPartialByID updates only provided fields",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Camera",
					Desc:     "Old camera",
					Price:    utils.Money(80000),
					Pros:     "portable",
					Cons:     "low-res",
					UserID:   userID,
					Deadline: "2026-04-01",
				}
				created, err := CreateGoal(ctx, g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}

				newName := "Camera Pro"
				update := &model.GoalUpdate{
					Name: &newName,
				}

				got, err := UpdateGoalPartialByID(ctx, created.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateGoalPartialByID() returned error: %v", err)
				}
				if got.ID != created.ID {
					t.Errorf("id changed after update: expected %d, got %d", created.ID, got.ID)
				}
				if got.Name != newName {
					t.Errorf("name not updated: expected %q, got %q", newName, got.Name)
				}
				if got.Desc != g.Desc {
					t.Errorf("description should not change: expected %q, got %q", g.Desc, got.Desc)
				}
			},
		},
		{
			name: "UpdateGoalPartialByID updates multiple fields",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Phone",
					Desc:     "Old phone",
					Price:    utils.Money(50000),
					Pros:     "works",
					Cons:     "slow",
					UserID:   userID,
					Deadline: "2026-03-01",
				}
				created, err := CreateGoal(ctx, g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}

				newName := "Phone Pro"
				newPrice := utils.Money(120000)
				update := &model.GoalUpdate{
					Name:  &newName,
					Price: &newPrice,
				}

				got, err := UpdateGoalPartialByID(ctx, created.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateGoalPartialByID() returned error: %v", err)
				}
				if got.Name != newName {
					t.Errorf("name not updated: expected %q, got %q", newName, got.Name)
				}
				if got.Price != newPrice {
					t.Errorf("price not updated: expected %v, got %v", newPrice, got.Price)
				}
				if got.Desc != g.Desc {
					t.Errorf("description should not change: expected %q, got %q", g.Desc, got.Desc)
				}
			},
		},
		{
			name: "UpdateGoalPartialByID with no fields returns goal unchanged",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "Tablet",
					Desc:     "New tablet",
					Price:    utils.Money(60000),
					Pros:     "portable",
					Cons:     "expensive",
					UserID:   userID,
					Deadline: "2026-05-15",
				}
				created, err := CreateGoal(ctx, g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}

				update := &model.GoalUpdate{}

				got, err := UpdateGoalPartialByID(ctx, created.ID, update, db)
				if err != nil {
					t.Fatalf("UpdateGoalPartialByID() returned error: %v", err)
				}

				if got.Name != g.Name {
					t.Errorf("name changed unexpectedly: expected %q, got %q", g.Name, got.Name)
				}
				if got.Price != g.Price {
					t.Errorf("price changed unexpectedly: expected %v, got %v", g.Price, got.Price)
				}
			},
		},
		{
			name: "UpdateGoalPartialByID returns ErrNoRows for non-existent id",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				newName := "NonExistent"
				update := &model.GoalUpdate{
					Name: &newName,
				}
				_, err := UpdateGoalPartialByID(ctx, 999999, update, db)
				if err == nil {
					t.Fatalf("expected error when updating non-existent goal, got nil")
				}
				if !errors.Is(err, sql.ErrNoRows) {
					t.Fatalf("expected sql.ErrNoRows when updating non-existent goal, got: %v", err)
				}
			},
		},
		{
			name: "GetAllGoalsByUserID returns goals for specific user",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g1 := model.Goal{
					Name:     "Goal A",
					Desc:     "First goal",
					Price:    utils.Money(10000),
					Pros:     "good",
					Cons:     "none",
					UserID:   userID,
					Deadline: "2026-01-01",
				}
				g2 := model.Goal{
					Name:     "Goal B",
					Desc:     "Second goal",
					Price:    utils.Money(20000),
					Pros:     "great",
					Cons:     "pricey",
					UserID:   userID,
					Deadline: "2026-06-01",
				}
				if _, err := CreateGoal(ctx, g1, db); err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				if _, err := CreateGoal(ctx, g2, db); err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}

				list, err := GetAllGoalsByUserID(ctx, userID, db)
				if err != nil {
					t.Fatalf("GetAllGoalsByUserID() returned error: %v", err)
				}
				if len(list) != 2 {
					t.Fatalf("expected 2 goals for user, got %d", len(list))
				}
				for _, g := range list {
					if g.UserID != userID {
						t.Errorf("expected UserID=%d, got %d", userID, g.UserID)
					}
				}
			},
		},
		{
			name: "GetAllGoalsByUserID returns empty for user with no goals",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				list, err := GetAllGoalsByUserID(ctx, userID, db)
				if err != nil {
					t.Fatalf("GetAllGoalsByUserID() returned error: %v", err)
				}
				if len(list) != 0 {
					t.Fatalf("expected 0 goals for user with no goals, got %d", len(list))
				}
			},
		},
		{
			name: "GetAllGoalsByUserID returns only goals for the given user",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				// Create a second user
				otherUser := model.User{
					UserName:       "other-goal-user",
					CurrentAmount:  utils.Money(0),
					MonthlyInputs:  utils.Money(0),
					MonthlyOutputs: utils.Money(0),
				}
				otherRet, err := CreateUser(ctx, otherUser, db)
				if err != nil {
					t.Fatalf("failed to create other user: %v", err)
				}

				// Create goals for both users
				if _, err := CreateGoal(ctx, model.Goal{
					Name: "Main user goal", Desc: "desc", Price: utils.Money(100),
					Pros: "", Cons: "", UserID: userID, Deadline: "2026-01-01",
				}, db); err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				if _, err := CreateGoal(ctx, model.Goal{
					Name: "Other user goal", Desc: "desc", Price: utils.Money(200),
					Pros: "", Cons: "", UserID: otherRet.ID, Deadline: "2026-01-01",
				}, db); err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}

				// Query for main user only
				list, err := GetAllGoalsByUserID(ctx, userID, db)
				if err != nil {
					t.Fatalf("GetAllGoalsByUserID() returned error: %v", err)
				}
				if len(list) != 1 {
					t.Fatalf("expected 1 goal for main user, got %d", len(list))
				}
				if list[0].UserID != userID {
					t.Errorf("expected UserID=%d, got %d", userID, list[0].UserID)
				}
				if list[0].Name != "Main user goal" {
					t.Errorf("expected name=%q, got %q", "Main user goal", list[0].Name)
				}
			},
		},
		{
			name: "GetAllGoalsByUserID returns empty for non-existent user id",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				list, err := GetAllGoalsByUserID(ctx, 999999, db)
				if err != nil {
					t.Fatalf("GetAllGoalsByUserID() returned error: %v", err)
				}
				if len(list) != 0 {
					t.Fatalf("expected 0 goals for non-existent user, got %d", len(list))
				}
			},
		},
		{
			name: "DeleteGoalByID deletes existing goal and subsequent GetGoalByID returns ErrNoRows",
			testFn: func(t *testing.T, ctx context.Context, db *sql.DB, userID int64) {
				g := model.Goal{
					Name:     "To Delete",
					Desc:     "temporary",
					Price:    utils.Money(1000),
					Pros:     "",
					Cons:     "",
					UserID:   userID,
					Deadline: "2026-02-02",
				}
				created, err := CreateGoal(ctx, g, db)
				if err != nil {
					t.Fatalf("CreateGoal() returned error: %v", err)
				}
				deleted, err := DeleteGoalByID(ctx, created.ID, db)
				if err != nil {
					t.Fatalf("DeleteGoalByID() returned error: %v", err)
				}
				if deleted != 1 {
					t.Fatalf("expected 1 row deleted, got %d", deleted)
				}
				_, err = GetGoalByID(ctx, created.ID, db)
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			db, teardown := setupDB(t)
			defer teardown()
			ctx := context.Background()

			u := model.User{
				UserName:       "goal-user",
				CurrentAmount:  utils.Money(0),
				MonthlyInputs:  utils.Money(0),
				MonthlyOutputs: utils.Money(0),
			}
			uRet, err := CreateUser(ctx, u, db)
			if err != nil {
				t.Fatalf("failed to create user for goals tests: %v", err)
			}

			tc.testFn(t, ctx, db, uRet.ID)
		})
	}
}
