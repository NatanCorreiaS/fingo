package service

import (
	"context"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
)

// GetGoalByID returns the goal with the given ID.
func GetGoalByID(ctx context.Context, id int64) (*model.Goal, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return dbsqlite.GetGoalByID(ctx, id, db)
}

// GetAllGoals returns all goals in the database.
func GetAllGoals(ctx context.Context) ([]model.Goal, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return dbsqlite.GetAllGoals(ctx, db)
}

// CreateGoal persists a new goal and returns the created record.
func CreateGoal(ctx context.Context, goal model.Goal) (*model.Goal, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return dbsqlite.CreateGoal(ctx, goal, db)
}

// UpdateGoalByID applies a partial update to the goal with the given ID and returns the updated record.
func UpdateGoalByID(ctx context.Context, id int64, goal *model.GoalUpdate) (*model.Goal, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return dbsqlite.UpdateGoalPartialByID(ctx, id, goal, db)
}

// DeleteGoalByID removes the goal with the given ID and returns the number of affected rows.
func DeleteGoalByID(ctx context.Context, id int64) (int64, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	return dbsqlite.DeleteGoalByID(ctx, id, db)
}
