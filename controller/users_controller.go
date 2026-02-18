package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	"natan/fingo/service"
	"net/http"
	"strconv"
)

func GetID(idStr string, w http.ResponseWriter, r *http.Request) int64 {
	if idStr == "" {
		log.Println("could not get id in User Handler")
		http.Error(w, "Empty id!", http.StatusBadRequest)
		return 0
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("could not convert the path variable to id: %v", err)
		http.Error(w, "invalid id!", http.StatusBadRequest)
		return 0
	}

	return id
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	idStr := r.PathValue("id")
	id := GetID(idStr, w, r)
	if id == 0 {
		return
	}
	user, err := service.GetUserByID(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid id!", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User: %v", *user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()
	usersList, err := service.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		http.Error(w, "Problem to fetch all users", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Users: %v", usersList)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("could not decode request body: %v", err)
		http.Error(w, "invalid body!", http.StatusBadRequest)
		return
	}

	userRec, err := service.CreateUser(ctx, user)
	if err != nil {
		log.Println(err)
		http.Error(w, "error when creating user!", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User created: %v", *userRec)
}

func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	var userUpdate *model.UserUpdate

	idStr := r.PathValue("id")
	id := GetID(idStr, w, r)

	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {

		log.Printf("could not decode request body: %v", err)
		http.Error(w, "invalid body!", http.StatusBadRequest)
		return
	}

	user, err := service.UpdateUserByID(ctx, id, userUpdate)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error when updating user!", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User updated: %v", *user)
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	idStr := r.PathValue("id")
	id := GetID(idStr, w, r)

	if id == 0 {
		return
	}

	rows, err := service.DeleteUserByID(ctx, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error when deleting user!", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Deleted, rows affected: %v", rows)
}
