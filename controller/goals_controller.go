package controller

import (
	"encoding/json"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/model"
	"natan/fingo/service"
	"net/http"
)

// GetGoalByIDHandler handles GET /goals/{id} and returns the goal with the given ID.
func GetGoalByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	idStr := r.PathValue("id")
	id, ok := GetID(idStr, w, r)
	if !ok {
		return
	}

	goal, err := service.GetGoalByID(ctx, id)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "goal not found"})
		return
	}

	writeJSON(w, http.StatusOK, *goal)
}

// GetAllGoalsHandler handles GET /goals and returns all goals.
func GetAllGoalsHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	goalsList, err := service.GetAllGoals(ctx)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when fetching goals"})
		return
	}

	writeJSON(w, http.StatusOK, goalsList)
}

// CreateGoalHandler handles POST /goals and creates a new goal from the request body.
func CreateGoalHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()
	var goal model.Goal

	if err := json.NewDecoder(r.Body).Decode(&goal); err != nil {
		log.Printf("could not decode into body: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	goalRec, err := service.CreateGoal(ctx, goal)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when creating goal"})
		return
	}
	writeJSON(w, http.StatusCreated, *goalRec)
}

// UpdateGoalByIDHandler handles PATCH /goals/{id} and applies a partial update to the given goal.
func UpdateGoalByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()
	var goalUpdate *model.GoalUpdate
	idStr := r.PathValue("id")
	id, ok := GetID(idStr, w, r)
	if !ok {
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&goalUpdate); err != nil {
		log.Printf("could not decode into body: %v", err)
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	goal, err := service.UpdateGoalByID(ctx, id, goalUpdate)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when updating goal"})
		return
	}

	writeJSON(w, http.StatusOK, *goal)
}

// DeleteGoalByIDHandler handles DELETE /goals/{id} and removes the goal with the given ID.
func DeleteGoalByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := dbsqlite.NewDBContext()
	defer cancel()

	idStr := r.PathValue("id")
	id, ok := GetID(idStr, w, r)
	if !ok {
		return
	}

	rows, err := service.DeleteGoalByID(ctx, id)
	if err != nil {
		log.Println(err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "problem when deleting goal"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]int64{"rows_affected": rows})
}
