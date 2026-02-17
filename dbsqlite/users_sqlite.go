package dbsqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"natan/fingo/model"
)

// GetAllUsers retrieves all users from the database with context support
func GetAllUsers(ctx context.Context, db *sql.DB) ([]model.User, error) {
	const query = `
	SELECT id, user_name, current_amount, monthly_inputs, monthly_outputs FROM users ORDER BY id;
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not execute the query to return all users: %w", err)
	}
	defer rows.Close()

	var usersList []model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.CurrentAmount, &user.MonthlyInputs, &user.MonthlyOutputs); err != nil {
			return nil, fmt.Errorf("could not send the rows data to user struct: %w", err)
		}

		usersList = append(usersList, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return usersList, nil
}

// CreateUser inserts a new user into the database with context support
func CreateUser(ctx context.Context, user model.User, db *sql.DB) (*model.User, error) {
	const insertStmt = `INSERT INTO users(user_name, current_amount, monthly_inputs, monthly_outputs) VALUES (?, ?, ?, ?);`

	res, err := db.ExecContext(ctx, insertStmt, user.UserName, user.CurrentAmount, user.MonthlyInputs, user.MonthlyOutputs)
	if err != nil {
		return nil, fmt.Errorf("could not execute insert into users table: %w", err)
	}

	if id, err := res.LastInsertId(); err == nil {
		user.ID = id
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID from the database with context support
func GetUserByID(ctx context.Context, id int64, db *sql.DB) (*model.User, error) {
	const selectStmt = `SELECT id, user_name, current_amount, monthly_inputs, monthly_outputs FROM users WHERE id = ?`

	row := db.QueryRowContext(ctx, selectStmt, id)
	var user model.User
	if err := row.Scan(&user.ID, &user.UserName, &user.CurrentAmount, &user.MonthlyInputs, &user.MonthlyOutputs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("could not scan the row into user struct: %w", err)
	}

	return &user, nil
}

// DeleteUserByID deletes a user by ID from the database with context support
func DeleteUserByID(ctx context.Context, id int64, db *sql.DB) (int64, error) {
	const deleteStmt = `DELETE FROM users WHERE id = ?`

	res, err := db.ExecContext(ctx, deleteStmt, id)
	if err != nil {
		return 0, fmt.Errorf("could not execute the delete query for user: %w", err)
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not get rows affected for delete: %w", err)
	}

	return deleted, nil
}

// UpdateUserPartialByID updates a user by ID with partial user data.
// Only fields that are provided (non-nil) in the UserUpdate struct will be updated.
// This prevents overwriting existing values with zero/empty values.
func UpdateUserPartialByID(ctx context.Context, id int64, update *model.UserUpdate, db *sql.DB) (*model.User, error) {
	if update == nil {
		return nil, fmt.Errorf("update data cannot be nil")
	}

	// First, fetch the current user data to verify it exists
	_, err := GetUserByID(ctx, id, db)
	if err != nil {
		return nil, err
	}

	// Build dynamic UPDATE statement based on which fields are provided
	var setParts []string
	var args []interface{}

	if update.UserName != nil {
		setParts = append(setParts, "user_name = ?")
		args = append(args, *update.UserName)
	}

	if update.CurrentAmount != nil {
		setParts = append(setParts, "current_amount = ?")
		args = append(args, *update.CurrentAmount)
	}

	if update.MonthlyInputs != nil {
		setParts = append(setParts, "monthly_inputs = ?")
		args = append(args, *update.MonthlyInputs)
	}

	if update.MonthlyOutputs != nil {
		setParts = append(setParts, "monthly_outputs = ?")
		args = append(args, *update.MonthlyOutputs)
	}

	// If no fields to update, return current user as-is
	if len(setParts) == 0 {
		return GetUserByID(ctx, id, db)
	}

	// Build and execute the dynamic UPDATE statement
	updateStmt := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setParts, ", "))
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

	// Fetch and return the updated user
	return GetUserByID(ctx, id, db)
}
