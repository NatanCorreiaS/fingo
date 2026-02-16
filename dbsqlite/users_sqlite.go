package dbsqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"natan/fingo/model"
)

func ReturnAllUsers(db *sql.DB) ([]model.User, error) {
	const query = `
	SELECT id, user_name, current_amount, monthly_inputs, monthly_outputs FROM users ORDER BY id;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute the query to return all users: %v", err)
	}
	defer rows.Close()

	var usersList []model.User

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.UserName, &user.CurrentAmount, &user.MonthlyInputs, &user.MonthlyOutputs); err != nil {
			return nil, fmt.Errorf("could not send the rows data to user struct: %v", err)
		}

		usersList = append(usersList, user)
	}

	return usersList, nil
}

func CreateUser(user model.User, db *sql.DB) (*model.User, error) {
	const insertStmt = `INSERT INTO users(user_name, current_amount, monthly_inputs, monthly_outputs) VALUES (?, ?, ?, ?);`

	res, err := db.Exec(insertStmt, user.UserName, user.CurrentAmount, user.MonthlyInputs, user.MonthlyOutputs)
	if err != nil {
		return nil, fmt.Errorf("could not execute insert into users table: %v", err)
	}

	// Optionally populate the inserted ID on the returned struct
	if id, err := res.LastInsertId(); err == nil {
		user.ID = id
	}

	return &user, nil
}

func GetUserByID(id int64, db *sql.DB) (*model.User, error) {
	const selectStmt = `SELECT id, user_name, current_amount, monthly_inputs, monthly_outputs FROM users WHERE id = ?`

	row := db.QueryRow(selectStmt, id)
	var user model.User
	if err := row.Scan(&user.ID, &user.UserName, &user.CurrentAmount, &user.MonthlyInputs, &user.MonthlyOutputs); err != nil {
		// Wrap the underlying error so callers can use errors.Is(err, sql.ErrNoRows)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("could not scan the row into user struct: %w", err)
	}

	return &user, nil
}

func DeleteUserByID(id int64, db *sql.DB) (int64, error) {
	const deleteStmt = `DELETE FROM users WHERE id = ?`

	res, err := db.Exec(deleteStmt, id)
	if err != nil {
		return 0, fmt.Errorf("could not execute the delete query for user: %v", err)
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not get rows affected for delete: %v", err)
	}

	return deleted, nil
}

func UpdateUserByID(id int64, user *model.User, db *sql.DB) (*model.User, error) {
	const updateStmt = `
	UPDATE users
	SET user_name = ?, current_amount = ?, monthly_inputs = ?, monthly_outputs = ?
	WHERE id = ?
	`

	// Execute the update using placeholders to avoid SQL injection
	res, err := db.Exec(updateStmt, user.UserName, user.CurrentAmount, user.MonthlyInputs, user.MonthlyOutputs, id)
	if err != nil {
		return nil, fmt.Errorf("could not execute update query: %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("could not get rows affected: %v", err)
	}
	if affected == 0 {
		// No rows updated: id not found
		return nil, sql.ErrNoRows
	}

	// Fetch the updated row to return the canonical values from DB
	const selectStmt = `SELECT id, user_name, current_amount, monthly_inputs, monthly_outputs FROM users WHERE id = ?`
	var u model.User
	row := db.QueryRow(selectStmt, id)
	if err := row.Scan(&u.ID, &u.UserName, &u.CurrentAmount, &u.MonthlyInputs, &u.MonthlyOutputs); err != nil {
		return nil, fmt.Errorf("could not fetch updated user: %v", err)
	}

	return &u, nil
}
