package dbsqlite

import (
	"database/sql"
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
	query := fmt.Sprintf("INSERT INTO users(user_name, current_amount, monthly_inputs, monthly_outputs)VALUES('%v', %v, %v, %v);", user.UserName, user.CurrentAmount, user.MonthlyInputs, user.MonthlyOutputs)

	_, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute insert into users table: %v", err)
	}

	return &user, nil

}
