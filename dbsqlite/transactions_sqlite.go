package dbsqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"natan/fingo/model"
)

func GetAllTransactions(db *sql.DB) ([]model.Transaction, error) {
	const query = "SELECT id, description, amount, is_debt, created_at, user_id FROM transactions"
	var transactionsList []model.Transaction

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not execute the query to return all transactions: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Desc, &transaction.Amount, &transaction.IsDebt, &transaction.CreatedAt, &transaction.UserID); err != nil {
			return nil, fmt.Errorf("could not send the rows data to transaction struct: %v", err)
		}
		transactionsList = append(transactionsList, transaction)
	}

	return transactionsList, nil
}

func CreateTransaction(transaction model.Transaction, db *sql.DB) (*model.Transaction, error) {
	const createStmt = "INSERT INTO transactions(description, amount, is_debt, user_id)VALUES(?,?,?,?)"

	res, err := db.Exec(createStmt, transaction.Desc, transaction.Amount, transaction.IsDebt, transaction.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not execute insert into transaction table: %v", err)
	}

	if id, err := res.LastInsertId(); err == nil {
		transaction.ID = id
	}

	return &transaction, nil
}

func GetTransactionByID(id int64, db *sql.DB) (*model.Transaction, error) {
	const selectStmt = "SELECT id, description, amount, is_debt, created_at, user_id FROM transactions WHERE id = ?"

	var transaction model.Transaction
	row := db.QueryRow(selectStmt, id)
	if err := row.Scan(&transaction.ID, &transaction.Desc, &transaction.Amount, &transaction.CreatedAt, &transaction.IsDebt, &transaction.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("could not scan the row into transaction struct: %w", err)
	}

	return &transaction, nil
}

func DeleteTransactionByID(id int64, db *sql.DB) (int64, error) {
	const deleteStmt = "DELETE FROM transactions WHERE Id = ?"

	res, err := db.Exec(deleteStmt, id)
	if err != nil {
		return 0, fmt.Errorf("could not execute the delete query for transaction: %v", err)
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not get rows affected for delete: %v", err)
	}

	return deleted, nil
}

func UpdateTransactionByID(id int64, transaction *model.Transaction, db *sql.DB) (*model.Transaction, error) {
	const updateStmt = "UPDATE transactions SET description = ?, amount = ?, is_debt = ? WHERE id = ?"

	res, err := db.Exec(updateStmt, transaction.Desc, transaction.Amount, transaction.IsDebt, id)

	if err != nil {
		return nil, fmt.Errorf("could not execute update query: %v", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("could not get rows affected: %v", err)
	}

	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	t, err := GetTransactionByID(id, db)
	if err != nil {
		return nil, err
	}

	return t, nil

}
