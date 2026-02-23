package controller

import (
	"encoding/json"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	"natan/fingo/service"
	"net/http"
)

// GetUserByIDHandler handles GET /users/{id} and returns the user with the given ID.
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

// GetAllUsersHandler handles GET /users and returns all users.
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
func GetAllTransactionsByUserIDHandler(w http.ResponseWriter, r *http.Request){
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()
	
	id, ok := GetID(r.PathValue("id"),w ,r)
	if !ok{
		return
	}
	
	transactionsList, err := service.GetAllTransactionsByUserID(ctx, id)
	if err != nil{
		log.Println(err)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "transactions not found for user"})
		return
	}
	
	writeJSON(w, http.StatusOK, transactionsList)
}

func GetAllGoalsByUserIDHandler(w http.ResponseWriter, r *http.Request){
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()
	
	id, ok := GetID(r.PathValue("id"), w,r)
	if !ok{
		return
	}
	
	goalsList, err := service.GetAllGoalsByUserID(ctx, id)
	if err != nil{
		log.Println(err)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "goals not found for user"})
		return
	}
	writeJSON(w, http.StatusOK, goalsList)
}

// CreateUserHandler handles POST /users and creates a new user from the request body.
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

// UpdateUserByIDHandler handles PATCH /users/{id} and applies a partial update to the given user.
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

// DeleteUserByIDHandler handles DELETE /users/{id} and removes the user with the given ID.
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
