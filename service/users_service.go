package service

import (
	"context"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
)

func CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	u, err := dbsqlite.CreateUser(ctx, user, db)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	u, err := dbsqlite.GetUserByID(ctx, id, db)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func GetAllUsers(ctx context.Context) ([]model.User, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	users, err := dbsqlite.GetAllUsers(ctx, db)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func DeleteUserByID(ctx context.Context, id int64) (int64, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return 0, err
	}

	defer db.Close()

	rows, err := dbsqlite.DeleteUserByID(ctx, id, db)
	if err != nil {
		return rows, err
	}

	return rows, nil
}

func UpdateUserByID(ctx context.Context, id int64, user *model.UserUpdate) (*model.User, error) {
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		return nil, err
	}

	defer db.Close()
	
	u, err := dbsqlite.UpdateUserPartialByID(ctx, id, user, db)
	if err != nil{
		return nil, err
	}
	
	return u, nil
}
