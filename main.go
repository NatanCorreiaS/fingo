package main

import (
	"fmt"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	// "net/http"
)

func main() {
	// http.HandleFunc("/", teste)
	if err := dbsqlite.CheckAndCreate(); err != nil {
		log.Fatal(err)
	}
	user1 := model.User{UserName: "Natan", CurrentAmount: 100, MonthlyInputs: 0, MonthlyOutputs: 0}
	user2 := model.User{UserName: "Adriana", CurrentAmount: 1000, MonthlyInputs: 300, MonthlyOutputs: 100}

	db, err := dbsqlite.GetDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	u1, err := dbsqlite.CreateUser(user1, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*u1)
	u2, err := dbsqlite.CreateUser(user2, db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*u2)

	users, err := dbsqlite.ReturnAllUsers(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)

	user, err := dbsqlite.GetUserByID(1, db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Usuario por id: ", *user)

	row, err := dbsqlite.DeleteUserByID(2, db)
	if err != nil {
		log.Fatal(err)
	}

	if row == 1 {
		fmt.Println("usuario deletado")
	} else {
		fmt.Println("usuario não foi deletado")
	}

	users, _ = dbsqlite.ReturnAllUsers(db)
	fmt.Println(users)

	user, err = dbsqlite.UpdateUserByID(1, &model.User{UserName: "Sigma", CurrentAmount: 1000, MonthlyInputs: 800, MonthlyOutputs: 300}, db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User updated: ", *user)
	// log.Fatal(http.ListenAndServe(":9000", nil))

}
