package dbsqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"natan/fingo/model"
	"strings"
)

// Retrieves all goals from the database with context support
func GetAllGoals(ctx context.Context, db *sql.DB) ([]model.Goal, error) {
	const query = "SELECT id, name, description, price, pros, cons, user_id, created_at, deadline FROM goals"

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not execute the query to return all goals: %w", err)
	}
	defer rows.Close()

	var goalsList []model.Goal

	for rows.Next() {
		var goal model.Goal
		if err = rows.Scan(&goal.ID, &goal.Name, &goal.Desc, &goal.Price, &goal.Pros, &goal.Cons, &goal.UserID, &goal.CreatedAt, &goal.Deadline); err != nil {
			return nil, fmt.Errorf("could not scan the data into goal struct: %w", err)
		}
		goalsList = append(goalsList, goal)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return goalsList, nil
}

// Retrieves a goal by its ID with context support
func GetGoalByID(ctx context.Context, id int64, db *sql.DB) (*model.Goal, error) {
	const selectStmt = "SELECT id, name, description, price, pros, cons, user_id, created_at, deadline FROM goals WHERE id = ?"

	row := db.QueryRowContext(ctx, selectStmt, id)
	var goal model.Goal
	if err := row.Scan(&goal.ID, &goal.Name, &goal.Desc, &goal.Price, &goal.Pros, &goal.Cons, &goal.UserID, &goal.CreatedAt, &goal.Deadline); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("goal not found: %w", err)
		}
		return nil, fmt.Errorf("could not scan the row into goal struct: %w", err)
	}

	return &goal, nil
}

// Inserts a new goal into the database with context support
func CreateGoal(ctx context.Context, goal model.Goal, db *sql.DB) (*model.Goal, error) {
	const createStmt = "INSERT INTO goals(name, description, price, pros, cons, user_id, deadline)VALUES(?,?,?,?,?,?,?)"

	res, err := db.ExecContext(ctx, createStmt, goal.Name, goal.Desc, goal.Price, goal.Pros, goal.Cons, goal.UserID, goal.Deadline)
	if err != nil {
		return nil, fmt.Errorf("could not execute insert into goals table: %w", err)
	}

	if id, err := res.LastInsertId(); err == nil {
		goal.ID = id
	}

	return &goal, nil
}

// Deletes a goal by its ID with context support
func DeleteGoalByID(ctx context.Context, id int64, db *sql.DB) (int64, error) {
	const deleteStmt = "DELETE FROM goals WHERE id = ?"

	res, err := db.ExecContext(ctx, deleteStmt, id)
	if err != nil {
		return 0, fmt.Errorf("could not execute delete query for goal: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not get affected rows for delete: %w", err)
	}

	return rows, nil
}

// Updates a goal by its ID with only the provided fields.
// Fields not provided (nil) remain unchanged, preserving existing values.
func UpdateGoalPartialByID(ctx context.Context, id int64, update *model.GoalUpdate, db *sql.DB) (*model.Goal, error) {
	if update == nil {
		return nil, fmt.Errorf("update data cannot be nil")
	}

	// Verify that the goal exists
	_, err := GetGoalByID(ctx, id, db)
	if err != nil {
		return nil, err
	}

	// Build dynamic UPDATE statement with only provided fields
	var setParts []string
	var args []interface{}

	if update.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *update.Name)
	}

	if update.Desc != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *update.Desc)
	}

	if update.Price != nil {
		setParts = append(setParts, "price = ?")
		args = append(args, *update.Price)
	}

	if update.Pros != nil {
		setParts = append(setParts, "pros = ?")
		args = append(args, *update.Pros)
	}

	if update.Cons != nil {
		setParts = append(setParts, "cons = ?")
		args = append(args, *update.Cons)
	}

	if update.Deadline != nil {
		setParts = append(setParts, "deadline = ?")
		args = append(args, *update.Deadline)
	}

	// If no fields are provided for update, return the current goal as-is
	if len(setParts) == 0 {
		return GetGoalByID(ctx, id, db)
	}

	// Construct and execute the dynamic UPDATE statement
	updateStmt := fmt.Sprintf("UPDATE goals SET %s WHERE id = ?", strings.Join(setParts, ", "))
	args = append(args, id)

	res, err := db.ExecContext(ctx, updateStmt, args...)
	if err != nil {
		return nil, fmt.Errorf("could not execute partial update query: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("could not get rows affected: %w", err)
	}

	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	// Retrieve and return the updated goal
	return GetGoalByID(ctx, id, db)
}
