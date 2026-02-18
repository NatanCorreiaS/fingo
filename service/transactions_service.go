package service

import (
	"context"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
)

func GetTransactionByID(ctx context.Context, id int64) (*model.Transaction, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}

	transaction, err := dbsqlite.GetTransactionByID(ctx, id, db)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func GetAllTransactions(ctx context.Context) ([]model.Transaction, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	
	transactions, err := dbsqlite.GetAllTransactions(ctx, db)
	if err != nil{
		return nil, err
	}
	
	return transactions, nil
}

func CreateTransaction(ctx context.Context, transaction model.Transaction)(*model.Transaction, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return nil, err
	}
	
	transactionRec, err := dbsqlite.CreateTransaction(ctx, transaction, db)
	if err != nil{
		return nil, err
	}
	
	return transactionRec, nil
}

func UpdateTransactionByID(ctx context.Context, id int64, transaction *model.TransactionUpdate)(*model.Transaction, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return nil, err
	}
	
	transactionRec, err := dbsqlite.UpdateTransactionPartialByID(ctx, id, transaction, db)
	if err != nil{
		return nil, err
	}
	
	return transactionRec, nil
}

func DeleteTransactionByID(ctx context.Context, id int64)(int64, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return 0, err
	}
	
	rows, err := dbsqlite.DeleteTransactionByID(ctx, id, db)
	if err != nil{
		return 0, err
	}
	
	return rows, nil
}