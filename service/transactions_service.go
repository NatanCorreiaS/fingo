package service

import (
	"context"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	"natan/fingo/utils"
)

// moneyPtr returns a pointer to the given Money value.
func moneyPtr(m utils.Money) *utils.Money {
	return &m
}

// applyBalanceDelta adjusts a balance by the given amount.
// If isDebt is true, the amount is subtracted; otherwise it is added.
func applyBalanceDelta(current utils.Money, amount utils.Money, isDebt bool) utils.Money {
	if isDebt {
		return current - amount
	}
	return current + amount
}

// GetTransactionByID returns the transaction with the given ID.
func GetTransactionByID(ctx context.Context, id int64) (*model.Transaction, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return dbsqlite.GetTransactionByID(ctx, id, db)
}

// GetAllTransactions returns all transactions in the database.
func GetAllTransactions(ctx context.Context) ([]model.Transaction, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return dbsqlite.GetAllTransactions(ctx, db)
}

// CreateTransaction persists a new transaction and updates the owner's balance accordingly.
// Debts decrease the balance; credits increase it.
func CreateTransaction(ctx context.Context, transaction model.Transaction) (*model.Transaction, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	created, err := dbsqlite.CreateTransaction(ctx, transaction, db)
	if err != nil {
		return nil, err
	}

	user, err := dbsqlite.GetUserByID(ctx, created.UserID, db)
	if err != nil {
		return nil, err
	}

	newBalance := applyBalanceDelta(user.CurrentAmount, created.Amount, created.IsDebt)

	_, err = dbsqlite.UpdateUserPartialByID(ctx, created.UserID, &model.UserUpdate{
		CurrentAmount: moneyPtr(newBalance),
	}, db)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// UpdateTransactionByID applies a partial update to the transaction with the given ID.
// If IsDebt changes, the owner's balance is adjusted to reflect the new transaction type.
func UpdateTransactionByID(ctx context.Context, id int64, update *model.TransactionUpdate) (*model.Transaction, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	original, err := dbsqlite.GetTransactionByID(ctx, id, db)
	if err != nil {
		return nil, err
	}

	updated, err := dbsqlite.UpdateTransactionPartialByID(ctx, id, update, db)
	if err != nil {
		return nil, err
	}

	if update.IsDebt != nil && *update.IsDebt != original.IsDebt {
		user, err := dbsqlite.GetUserByID(ctx, updated.UserID, db)
		if err != nil {
			return nil, err
		}

		balanceWithoutOld := applyBalanceDelta(user.CurrentAmount, original.Amount, !original.IsDebt)
		newBalance := applyBalanceDelta(balanceWithoutOld, updated.Amount, updated.IsDebt)

		_, err = dbsqlite.UpdateUserPartialByID(ctx, updated.UserID, &model.UserUpdate{
			CurrentAmount: moneyPtr(newBalance),
		}, db)
		if err != nil {
			return nil, err
		}
	}

	return updated, nil
}

// DeleteTransactionByID removes the transaction with the given ID and reverts its effect on the owner's balance.
// Returns 0 with no error if the transaction does not exist.
func DeleteTransactionByID(ctx context.Context, id int64) (int64, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := dbsqlite.GetTransactionByID(ctx, id, db)
	if err != nil {
		return 0, nil
	}

	rows, err := dbsqlite.DeleteTransactionByID(ctx, id, db)
	if err != nil {
		return 0, err
	}

	if rows == 0 {
		return 0, nil
	}

	user, err := dbsqlite.GetUserByID(ctx, tx.UserID, db)
	if err != nil {
		return 0, err
	}

	newBalance := applyBalanceDelta(user.CurrentAmount, tx.Amount, !tx.IsDebt)

	_, err = dbsqlite.UpdateUserPartialByID(ctx, tx.UserID, &model.UserUpdate{
		CurrentAmount: moneyPtr(newBalance),
	}, db)
	if err != nil {
		return 0, err
	}

	return rows, nil
}
