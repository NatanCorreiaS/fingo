package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// writeJSON writes a JSON-encoded response with the given status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("could not encode response: %v", err)
	}
}

// GetID parses and returns the ID path variable as int64.
// Writes a 400 response and returns false if the value is missing or not a valid integer.
func GetID(idStr string, w http.ResponseWriter, r *http.Request) (int64, bool) {
	if idStr == "" {
		log.Println("could not get id in the URI")
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
