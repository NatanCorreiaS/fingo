package dbsqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"natan/fingo/model"
	"strings"
)

// GetAllTransactions retrieves all transactions from the database
func GetAllTransactions(ctx context.Context, db *sql.DB) ([]model.Transaction, error) {
	const query = "SELECT id, description, amount, is_debt, created_at, user_id FROM transactions"
	var transactionsList []model.Transaction

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not execute the query to return all transactions: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Desc, &transaction.Amount, &transaction.IsDebt, &transaction.CreatedAt, &transaction.UserID); err != nil {
			return nil, fmt.Errorf("could not send the rows data to transaction struct: %w", err)
		}
		transactionsList = append(transactionsList, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactionsList, nil
}

func GetAllTransactionsByUserID(ctx context.Context, id int64, db *sql.DB) ([]model.Transaction, error) {
	const query = "SELECT id, description, amount, is_debt, created_at, user_id FROM transactions WHERE user_id = ?"
	var transactionsList []model.Transaction

	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("could not execute the query for transactions using user_id: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var transaction model.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Desc, &transaction.Amount, &transaction.IsDebt, &transaction.CreatedAt, &transaction.UserID); err != nil {
			return nil, fmt.Errorf("could not send the rows data to transaction struct: %w", err)
		}
		transactionsList = append(transactionsList, transaction)
	}
	if err := rows.Err(); err != nil{
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	
	return transactionsList, nil
}

// CreateTransaction inserts a new transaction into the database
func CreateTransaction(ctx context.Context, transaction model.Transaction, db *sql.DB) (*model.Transaction, error) {
	const createStmt = "INSERT INTO transactions(description, amount, is_debt, user_id)VALUES(?,?,?,?)"

	res, err := db.ExecContext(ctx, createStmt, transaction.Desc, transaction.Amount, transaction.IsDebt, transaction.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not execute insert into transaction table: %w", err)
	}

	if id, err := res.LastInsertId(); err == nil {
		transaction.ID = id
	}

	return &transaction, nil
}

// GetTransactionByID retrieves a transaction by its ID
func GetTransactionByID(ctx context.Context, id int64, db *sql.DB) (*model.Transaction, error) {
	const selectStmt = "SELECT id, description, amount, is_debt, created_at, user_id FROM transactions WHERE id = ?"

	var transaction model.Transaction
	row := db.QueryRowContext(ctx, selectStmt, id)
	if err := row.Scan(&transaction.ID, &transaction.Desc, &transaction.Amount, &transaction.IsDebt, &transaction.CreatedAt, &transaction.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("transaction not found: %w", err)
		}
		return nil, fmt.Errorf("could not scan the row into transaction struct: %w", err)
	}

	return &transaction, nil
}

// DeleteTransactionByID deletes a transaction by its ID
func DeleteTransactionByID(ctx context.Context, id int64, db *sql.DB) (int64, error) {
	const deleteStmt = "DELETE FROM transactions WHERE Id = ?"

	res, err := db.ExecContext(ctx, deleteStmt, id)
	if err != nil {
		return 0, fmt.Errorf("could not execute the delete query for transaction: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("could not get rows affected for delete: %w", err)
	}

	return rows, nil
}

// UpdateTransactionPartialByID updates a transaction by its ID with only the provided fields
// Fields not provided (nil) are not updated, preserving existing values
func UpdateTransactionPartialByID(ctx context.Context, id int64, update *model.TransactionUpdate, db *sql.DB) (*model.Transaction, error) {
	if update == nil {
		return nil, fmt.Errorf("update data cannot be nil")
	}

	// Verify that the transaction exists
	_, err := GetTransactionByID(ctx, id, db)
	if err != nil {
		return nil, err
	}

	// Dynamically build the UPDATE statement with only the provided fields
	var setParts []string
	var args []interface{}

	if update.Desc != nil {
		setParts = append(setParts, "description = ?")
		args = append(args, *update.Desc)
	}

	if update.Amount != nil {
		setParts = append(setParts, "amount = ?")
		args = append(args, *update.Amount)
	}

	if update.IsDebt != nil {
		setParts = append(setParts, "is_debt = ?")
		args = append(args, *update.IsDebt)
	}

	// If no fields are provided, return the current transaction without modifications
	if len(setParts) == 0 {
		return GetTransactionByID(ctx, id, db)
	}

	// Build and execute the dynamic UPDATE statement
	updateStmt := fmt.Sprintf("UPDATE transactions SET %s WHERE id = ?", strings.Join(setParts, ", "))
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

	// Fetch and return the updated transaction
	return GetTransactionByID(ctx, id, db)
}
