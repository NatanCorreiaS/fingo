package service

import (
	"context"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
)

func GetGoalByID(ctx context.Context, id int64)(*model.Goal, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return nil, err
	}
	
	goal, err := dbsqlite.GetGoalByID(ctx, id, db)
	if err != nil{
		return nil, err
	}
	
	return goal, nil
}

func GetAllGoals(ctx context.Context)([]model.Goal, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return nil, err
	}
	
	goals, err := dbsqlite.GetAllGoals(ctx, db)
	if err != nil{
		return  nil, err
	}
	
	return goals, nil
}

func CreateGoal(ctx context.Context, goal model.Goal)(*model.Goal, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return nil, err
	}
	
	goalRec, err := dbsqlite.CreateGoal(ctx, goal, db)
	if err != nil{
		return nil, err
	}
	
	return goalRec, nil
}

func UpdateGoalByID(ctx context.Context, id int64, goal *model.GoalUpdate)(*model.Goal, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return nil, err
	}
	
	goalRec, err := dbsqlite.UpdateGoalPartialByID(ctx, id, goal, db)
	if err != nil{
		return nil, err
	}
	
	return goalRec, nil
}

func DeleteGoalByID(ctx context.Context, id int64)(int64, error){
	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil{
		return 0, err
	}
	
	rows, err := dbsqlite.DeleteGoalByID(ctx, id, db)
	if err != nil{
		return 0, err
	}
	
	return rows, nil
}