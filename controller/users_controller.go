package controller

import (
	"encoding/json"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	"natan/fingo/service"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("could not encode response: %v", err)
	}
}

func GetID(idStr string, w http.ResponseWriter, r *http.Request) (int64, bool) {
	if idStr == "" {
		log.Println("could not get id in User Handler")
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "empty id"})
		return 0, false
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("could not convert the path variable to id: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

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
	writeJSON(w, http.StatusOK, user)
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "error when creating user"})
		return
	}
	writeJSON(w, http.StatusCreated, userRec)
}

func UpdateUserByID(w http.ResponseWriter, r *http.Request) {
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
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error when updating user"})
		return
	}
	writeJSON(w, http.StatusOK, user)
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	id, ok := GetID(r.PathValue("id"), w, r)
	if !ok {
		return
	}

	rows, err := service.DeleteUserByID(ctx, id)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "error when deleting user"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]int64{"rows_affected": rows})
}