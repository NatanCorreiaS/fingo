package dbsqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"natan/fingo/model"
)

func GetAllGoals(db *sql.DB) ([]model.Goal, error) {
	const query = "SELECT id, name, description, price, pros, cons, user_id, created_at, deadline FROM goals"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute the query to return all goals: %v", err)
	}
	defer rows.Close()

	var goalsList []model.Goal

	for rows.Next() {
		var goal model.Goal
		// fixed scan order: id, name, description, price, pros, cons, user_id, created_at, deadline
		if err = rows.Scan(&goal.ID, &goal.Name, &goal.Desc, &goal.Price, &goal.Pros, &goal.Cons, &goal.UserID, &goal.CreatedAt, &goal.Deadline); err != nil {
			return nil, fmt.Errorf("could not scan the data into goal struct: %v", err)
		}
		goalsList = append(goalsList, goal)
	}

	return goalsList, nil
}

func GetGoalByID(id int64, db *sql.DB) (*model.Goal, error) {
	const selectStmt = "SELECT id, name, description, price, pros, cons, user_id, created_at, deadline FROM goals WHERE id = ?"

	row := db.QueryRow(selectStmt, id)
	var goal model.Goal
	// fixed scan order matching the select statement
	if err := row.Scan(&goal.ID, &goal.Name, &goal.Desc, &goal.Price, &goal.Pros, &goal.Cons, &goal.UserID, &goal.CreatedAt, &goal.Deadline); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("goal not found: %w", err)
		}
		return nil, fmt.Errorf("could not scan the row into goal struct: %v", err)
	}

	return &goal, nil
}

func CreateGoal(goal model.Goal, db *sql.DB) (*model.Goal, error) {
	const createStmt = "INSERT INTO goals(name, description, price, pros, cons, user_id, deadline)VALUES(?,?,?,?,?,?,?)"

	res, err := db.Exec(createStmt, goal.Name, goal.Desc, goal.Price, goal.Pros, goal.Cons, goal.UserID, goal.Deadline)
	if err != nil {
		return nil, fmt.Errorf("could not execute insert into goals table: %v", err)
	}

	if id, err := res.LastInsertId(); err == nil {
		goal.ID = id
	}

	return &goal, nil
}

func DeleteGoalByID(id int64, db *sql.DB) (int64, error) {
	const deleteStmt = "DELETE FROM goals WHERE id = ?"

	res, err := db.Exec(deleteStmt, id)
	if err != nil {
		return 0, fmt.Errorf("could not execute delete query for goal: %v", err)
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not get affected rows for delete: %v", err)
	}

	return deleted, nil
}

func UpdateGoalById(id int64, goal model.Goal, db *sql.DB) (*model.Goal, error) {
	const updateStmt = "UPDATE goals SET name = ?, description = ?, price = ?, pros = ?, cons = ?, deadline = ? WHERE id = ?"
	res, err := db.Exec(updateStmt, goal.Name, goal.Desc, goal.Price, goal.Pros, goal.Cons, goal.Deadline, id)
	if err != nil {
		return nil, fmt.Errorf("could not execute the update query: %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("could not get rows affected: %v", err)
	}

	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	g, err := GetGoalByID(id, db)
	if err != nil {
		return nil, err
	}

	return g, nil
}
