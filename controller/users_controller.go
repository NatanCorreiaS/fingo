package controller

import (
	"encoding/json"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	"natan/fingo/service"
	"net/http"
)

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	id, ok := GetID(r.PathValue("id"), w, r)
	if !ok {
		return
	}

	user, err := service.GetUserByID(ctx, id)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	writeJSON(w, http.StatusOK, *user)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	usersList, err := service.GetAllUsers(ctx)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem fetching users"})
		return
	}
	writeJSON(w, http.StatusOK, usersList)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("could not decode request body: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	userRec, err := service.CreateUser(ctx, user)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "problem when creating user"})
		return
	}
	writeJSON(w, http.StatusCreated, *userRec)
}

func UpdateUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	id, ok := GetID(r.PathValue("id"), w, r)
	if !ok {
		return
	}

	var userUpdate *model.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&userUpdate); err != nil {
		log.Printf("could not decode request body: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	user, err := service.UpdateUserByID(ctx, id, userUpdate)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when updating user"})
		return
	}
	writeJSON(w, http.StatusOK, *user)
}

func DeleteUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	id, ok := GetID(r.PathValue("id"), w, r)
	if !ok {
		return
	}

	rows, err := service.DeleteUserByID(ctx, id)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when deleting user"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"rows_affected": rows})
}
