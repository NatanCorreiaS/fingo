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
	// log.Fatal(http.ListenAndServe(":9000", nil))

}
